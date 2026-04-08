package main

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"telemetry-engine/internal/pb"
	"telemetry-engine/internal/repository"
	"telemetry-engine/internal/telemetry"
	"time"
	"unsafe"

	"github.com/gorilla/websocket"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/protobuf/proto"
)

var client = &http.Client{
	Timeout: 500 * time.Millisecond,
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func worker(id int, jobs <-chan telemetry.Task, results chan<- telemetry.OptimizedMetric, wg *sync.WaitGroup, ctx context.Context) {
	defer wg.Done()
	scraper := telemetry.NewScraper()
	for {
		select {
		case job, ok := <-jobs:
			if !ok {
				return
			}
			metric, err := scraper.ScrapeTarget(ctx, job.ID)
			if err != nil {
				continue
			}
			results <- metric
			fmt.Printf("Worker %d fetched metrics from server %d\n", id, job.ID)
		case <-ctx.Done():
			fmt.Printf("Worker %d shutting down\n", id)
			return
		}
	}
}

func startAggregator(ctx context.Context, results <-chan telemetry.OptimizedMetric, buffer *telemetry.RingBuffer, hub *telemetry.Hub, repo *repository.PostgresRepo) {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()
	dbBatch := make([]telemetry.OptimizedMetric, 0, 100)
	for {
		select {
		case res, ok := <-results:
			if !ok {
				return
			}
			buffer.Add(res)
			dbBatch = append(dbBatch, res)
			if len(dbBatch) >= 100 {
				go func(b []telemetry.OptimizedMetric) {
					if err := repo.StoreBatch(context.Background(), b); err != nil {
						fmt.Printf("Database sync error: %v\n", err)
					}
				}(dbBatch)
				dbBatch = make([]telemetry.OptimizedMetric, 0, 100)
			}
		case <-ticker.C:
			latest := buffer.GetAll()
			if len(latest) == 0 {
				continue
			}
			batch := &pb.MetricBatch{}
			for _, m := range latest {
				batch.Metrics = append(batch.Metrics, &pb.Metric{
					Id:        m.ID,
					Value:     m.Value,
					Type:      uint32(m.Type),
					Timestamp: time.Now().Unix(),
				})
			}
			out, err := proto.Marshal(batch)
			if err != nil {
				fmt.Println("Error marshaling", err)
				continue
			}
			hub.Broadcast(out)
		case <-ctx.Done():
			if len(dbBatch) > 0 {
				repo.StoreBatch(context.Background(), dbBatch)
			}
			fmt.Println("Aggregator shutting down...")
			return
		}
	}
}

func runEngine() {
	hub := telemetry.NewHub()
	go hub.Run()
	connStr := "postgres://postgres:password@localhost:5432/telemetry_db"
	pool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		fmt.Printf("Unable to connect to database: %v\n", err)
		return
	}
	defer pool.Close()
	repo := &repository.PostgresRepo{Pool: pool}
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		hub.Register(conn)
		fmt.Println("New 3D Dashboard connected")
	})
	go func() {
		fmt.Println("Telemetry WebSocket server starting on :8080/ws")
		http.ListenAndServe(":8080", nil)
	}()
	buffer := telemetry.NewRingBuffer(telemetry.MaxBufferSize)
	jobs := make(chan telemetry.Task, 100)
	results := make(chan telemetry.OptimizedMetric, 100)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	go startAggregator(ctx, results, buffer, hub, repo)
	var wg sync.WaitGroup
	for w := 1; w <= telemetry.WorkerCount; w++ {
		wg.Add(1)
		go worker(w, jobs, results, &wg, ctx)
	}
	for j := 1; j <= 500; j++ {
		jobs <- telemetry.Task{ID: j}
	}
	close(jobs)
	wg.Wait()
	close(results)
	finalData := buffer.GetAll()
	fmt.Printf("Engine stopped. Final Buffer Count: %d metrics.\n", len(finalData))
}

func sliceExperiment() {
	startJunior := time.Now()
	juniorSlice := []telemetry.OptimizedMetric{}
	for i := 0; i < 100000; i++ {
		juniorSlice = append(juniorSlice, telemetry.OptimizedMetric{ID: uint32(i)})
	}
	fmt.Printf("Junior Append Time: %v\n", time.Since(startJunior))
	startSenior := time.Now()
	seniorSlice := make([]telemetry.OptimizedMetric, 0, 100000)
	for i := 0; i < 100000; i++ {
		seniorSlice = append(seniorSlice, telemetry.OptimizedMetric{ID: uint32(i)})
	}
	fmt.Printf("Senior Append Time: %v\n", time.Since(startSenior))
}

func main() {
	fmt.Printf("Size of BadMetric: %d bytes\n", unsafe.Sizeof(telemetry.BadMetric{}))
	fmt.Printf("Size of OptimizedMetric: %d bytes\n", unsafe.Sizeof(telemetry.OptimizedMetric{}))
	sliceExperiment()
	runEngine()
}
