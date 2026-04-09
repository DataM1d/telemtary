package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"telemetry-engine/internal/repository"
	"telemetry-engine/internal/telemetry"
	"time"

	"github.com/gorilla/websocket"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func main() {
	// 1. Logger Setup
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	// 2. Database Connection
	pool, err := pgxpool.New(context.Background(), "postgres://postgres:password@localhost:5432/telemetry_db")
	if err != nil {
		slog.Error("db connection failed", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	// 3. Initialize Engine Components
	hub := telemetry.NewHub()
	repo := &repository.PostgresRepo{Pool: pool}
	engine := telemetry.NewEngine(repo, hub)

	// 4. Routes
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		hub.Register(conn)
	})

	// 5. Lifecycle Management
	ctx, cancel := context.WithCancel(context.Background())
	go hub.Run()
	engine.Start(ctx)

	server := &http.Server{Addr: ":8080"}
	go func() {
		slog.Info("server starting", "port", 8080)
		server.ListenAndServe()
	}()

	// 6. Graceful Shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	slog.Info("shutting down...")
	cancel()
	engine.Stop()

	shutdownCtx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	server.Shutdown(shutdownCtx)
	slog.Info("goodbye")
}
