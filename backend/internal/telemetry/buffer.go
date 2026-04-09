package telemetry

import (
	"sync"
	"telemetry-engine/internal/models"
)

type RingBuffer struct {
	data   []models.OptimizedMetric
	size   int
	cursor int
	mu     sync.RWMutex
}

func NewRingBuffer(size int) *RingBuffer {
	return &RingBuffer{
		data: make([]models.OptimizedMetric, size),
		size: size,
	}
}

func (r *RingBuffer) Add(m models.OptimizedMetric) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.data[r.cursor] = m
	r.cursor = (r.cursor + 1) % r.size
}

func (r *RingBuffer) GetAll() []models.OptimizedMetric {
	r.mu.RLock()
	defer r.mu.RUnlock()

	out := make([]models.OptimizedMetric, len(r.data))
	copy(out, r.data)
	return out
}
