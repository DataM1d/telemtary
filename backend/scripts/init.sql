CREATE TABLE IF NOT EXISTS telemetry_data (
    id SERIAL PRIMARY KEY,
    server_id INT NOT NULL,
    metric_value DOUBLE PRECISION NOT NULL,
    metric_type INT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_server_id ON telemetry_data(server_id); 

