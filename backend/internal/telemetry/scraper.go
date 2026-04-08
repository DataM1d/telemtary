package telemetry

import (
	"context"
	"encoding/json"
	"net/http"
)

type Scraper struct {
	client *http.Client
}

func NewScraper() *Scraper {
	return &Scraper{
		client: &http.Client{},
	}
}

func (s *Scraper) ScrapeTarget(ctx context.Context, targetID int) (OptimizedMetric, error) {
	url := "http://localhost:8081/metrics"

	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
	resp, err := s.client.Do(req)
	if err != nil {
		return OptimizedMetric{}, err
	}
	defer resp.Body.Close()

	var data struct {
		ID    uint32  `json:"id"`
		Value float64 `json:"value"`
		Type  string  `json:"type"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return OptimizedMetric{}, err
	}

	return OptimizedMetric{
		ID:    data.ID,
		Value: data.Value,
		Type:  uint8(CPU),
	}, nil
}
