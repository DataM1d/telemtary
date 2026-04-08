A high concurrency telemetry processing system built in Go, designed to scrape, aggregate, and stream real time server metrics. This engine utilizes a hybrid storage strategy, leveraging a custom Ring Buffer for low latency streaming and PostgreSQL for long term persistence.

Tech Stack:
    Backend: Go (Golang)
    Communication: WebSockets & Protocol Buffers (Protobuf)
    Database: PostgreSQL (via pgx/v5 with connection pooling)
    Concurrency: Worker Pool pattern with context propagation
    Data Structures: Thread safe Circular Ring Buffer

System Architecture:
The engine is built on a four stage pipeline:
    1. Scraper: A pool of workers concurrently scrapes JSON metrics from distributed mock endpoints.

    2. Ring Buffer: Metrics are stored in a memory optimized struct (16 bytes) in a thread safe circular buffer for 60fps UI updates.

    3. Aggregator: Batches data into Protobuf binary to minimize network overhead.

    4. Presistence: Metrics are flushed to PostgreSQL using Async Batch Inserts to prevent database latency from blocking the live stream.

Performance Features:
    Zero Waste Memory: Optimized structs avoid padding, reducing memory footprint by 33% compared to standard layouts.

    Batch Processing: Database writes are batched (100+ rows per transaction) to maximize throughput.

    Binary Streaming: Uses Protobuf instead of JSON for WebSocket communication, drastically reducing payload size and GC pressure and speeds up append operations by ~50x.

Getting Started:
    Prerequisites
    Go 1.21+

    PostgreSQL

    protoc compiler (for Protobuf modifications)

Installation
    Clone the Repo: 
    git clone https://github.com/your-username/telemetry-engine.git
    cd telemetry-engine

    Setup the Database:
    psql -d telemetry_db -f scripts/init.sql

    Run the Mock Server:
    go run cmd/mock-server/main.go

    Run the Engine:
    go run cmd/api/main.go

The high concurrency Go backbone is complete, featuring a worker pool scraper, Protobuf serialization, and a batched PostgreSQL persistence layer. Currently, I am working on systems engineering and interactive visualization by implementing the Three.js frontend. This phase focuses on decoding binary WebSocket streams in real time to drive a React Three Fiber dashboard. 

On the way: instanced mesh rendering for massive server clusters, custom GLSL shaders for health status effects, and full OpenTelemetry tracing for front to back observability.