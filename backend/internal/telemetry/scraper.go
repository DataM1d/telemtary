package telemetry

import (
	"context"
	"encoding/json"
	"net/http"
	"telemetry-engine/internal/models" // This is already here, which is good
)

type Scraper struct {
	client *http.Client
}

func NewScraper() *Scraper {
	return &Scraper{
		client: &http.Client{},
	}
}

func (s *Scraper) ScrapeTarget(ctx context.Context, id int) (models.OptimizedMetric, error) {
	url := "http://localhost:8081/metrics"

	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
	resp, err := s.client.Do(req)
	if err != nil {
		return models.OptimizedMetric{}, err
	}
	defer resp.Body.Close()

	var data struct {
		ID    uint32  `json:"id"`
		Value float64 `json:"value"`
		Type  string  `json:"type"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return models.OptimizedMetric{}, err
	}

	var metricType uint8
	switch data.Type {
	case "CPU":
		metricType = models.CPU
	case "MEMORY":
		metricType = models.Memory
	default:
		metricType = models.CPU
	}

	return models.OptimizedMetric{
		ID:    data.ID,
		Value: data.Value,
		Type:  metricType,
	}, nil
}
