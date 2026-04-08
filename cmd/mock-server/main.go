package main

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"time"
)

type MockMetric struct {
	ID    int     `json:"id"`
	Value float64 `json:"value"`
	Type  string  `json:"type"`
}

func main() {
	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Duration(rand.Intn(50)) * time.Millisecond)

		data := MockMetric{
			ID:    1,
			Value: 40.0 + rand.Float64()*20.0,
			Type:  "cpu",
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(data)
	})

	println("Mock Server running on :8081...")
	http.ListenAndServe(":8081", nil)
}
