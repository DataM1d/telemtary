package telemetry

func NewRingBuffer(size int) *RingBuffer {
	return &RingBuffer{
		data: make([]OptimizedMetric, size),
		size: size,
	}
}

func (r *RingBuffer) Add(m OptimizedMetric) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.data[r.cursor] = m
	r.cursor = (r.cursor + 1) % r.size
}

func (r *RingBuffer) GetAll() []OptimizedMetric {
	r.mu.RLock()
	defer r.mu.RUnlock()

	out := make([]OptimizedMetric, len(r.data))
	copy(out, r.data)
	return out
}
