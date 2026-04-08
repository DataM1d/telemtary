export enum MetricType {
  CPU = 0,
  MEMORY = 1,
  DISK = 2,
  NETWORK = 3,
}

export interface Metric {
  id: number;
  value: number;
  type: MetricType;
  timestamp: number;
}

export interface MetricBatch {
  metrics: Metric[];
}

export interface TelemetryState {
  metrics: Metric[];
  isConnected: boolean;
  error: string | null;
}