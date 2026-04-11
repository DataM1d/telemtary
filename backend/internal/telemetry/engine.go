package telemetry

import (
	"context"
	"log/slog"
	"sync"
	"telemetry-engine/internal/models"
	"telemetry-engine/internal/pb"
	"telemetry-engine/internal/repository"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"google.golang.org/protobuf/proto"
)

const (
	WorkerCount   = 5
	MaxBufferSize = 1000
)

type Task struct {
	ID int
}

type Engine struct {
	Repo    *repository.PostgresRepo
	Hub     *Hub
	Buffer  *RingBuffer
	dbChan  chan []models.OptimizedMetric
	jobs    chan Task
	results chan models.OptimizedMetric
	wg      sync.WaitGroup
}

var (
	opsProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "telemetry_engine_ops_total",
		Help: "The total number of processed metrics",
	})
	scrapeLatency = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "telemetry_engine_scrape_duration_seconds",
		Help:    "Latency of scraping targets in seconds",
		Buckets: prometheus.DefBuckets,
	})
	slicePool = sync.Pool{
		New: func() interface{} {
			return make([]models.OptimizedMetric, 0, 100)
		},
	}
)

func NewEngine(repo *repository.PostgresRepo, hub *Hub) *Engine {
	return &Engine{
		Repo:    repo,
		Hub:     hub,
		Buffer:  NewRingBuffer(MaxBufferSize),
		dbChan:  make(chan []models.OptimizedMetric, 50),
		jobs:    make(chan Task, 100),
		results: make(chan models.OptimizedMetric, 100),
	}
}

func (e *Engine) worker(ctx context.Context, id int) {
	defer e.wg.Done()
	slog.Info("worker initialized", "worker_id", id)
	scraper := NewScraper()

	for {
		select {
		case job, ok := <-e.jobs:
			if !ok {
				return
			}
			start := time.Now()
			metric, err := scraper.ScrapeTarget(ctx, job.ID)
			scrapeLatency.Observe(time.Since(start).Seconds())
			opsProcessed.Inc()

			if err != nil {
				slog.Error("scrape error", "worker_id", id, "target_id", job.ID, "error", err)
				continue
			}
			e.results <- metric
		case <-ctx.Done():
			return
		}
	}
}

func (e *Engine) dbWorker(ctx context.Context) {
	for {
		select {
		case batch, ok := <-e.dbChan:
			if !ok {
				return
			}
			if err := e.Repo.StoreBatch(ctx, batch); err != nil {
				slog.Error("database sync error", "error", err)
			}
			slicePool.Put(batch[:0])
		case <-ctx.Done():
			return
		}
	}
}

func (e *Engine) startAggregator(ctx context.Context) {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()
	dbBatch := slicePool.Get().([]models.OptimizedMetric)

	for {
		select {
		case res, ok := <-e.results:
			if !ok {
				return
			}
			e.Buffer.Add(res)
			dbBatch = append(dbBatch, res)

			if len(dbBatch) >= 100 {
				select {
				case e.dbChan <- dbBatch:
					dbBatch = slicePool.Get().([]models.OptimizedMetric)
				default:
					slog.Warn("db worker busy, dropping batch")
					dbBatch = dbBatch[:0]
				}
			}
		case <-ticker.C:
			latest := e.Buffer.GetAll()
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
			if out, err := proto.Marshal(batch); err == nil {
				e.Hub.Broadcast(out)
			}
		case <-ctx.Done():
			if len(dbBatch) > 0 {
				e.dbChan <- dbBatch
			}
			return
		}
	}
}

func (e *Engine) produceJobs(ctx context.Context) {
	for j := 1; ; j++ {
		select {
		case e.jobs <- Task{ID: j}:
			time.Sleep(10 * time.Millisecond)
		case <-ctx.Done():
			close(e.jobs)
			return
		}
	}
}

func (e *Engine) Start(ctx context.Context) {
	go e.dbWorker(ctx)
	go e.startAggregator(ctx)

	for w := 1; w <= WorkerCount; w++ {
		e.wg.Add(1)
		go e.worker(ctx, w)
	}

	go e.produceJobs(ctx)
}

func (e *Engine) Stop() {
	e.wg.Wait()
	close(e.dbChan)
}
