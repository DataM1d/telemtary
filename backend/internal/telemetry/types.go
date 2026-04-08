package telemetry

import (
	"sync"
)

type MetricType uint8

const (
	CPU MetricType = iota
	Memory
	Network
	Disk
)

const (
	MaxBufferSize = 10000
	WorkerCount   = 5
)

type BadMetric struct {
	Active bool
	Value  float64
	ID     uint32
}

type OptimizedMetric struct {
	Value float64
	ID    uint32
	Type  uint8
}

type Task struct {
	ID int
}

type RingBuffer struct {
	data   []OptimizedMetric
	size   int
	cursor int
	mu     sync.RWMutex
}
