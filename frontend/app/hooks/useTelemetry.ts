"use client";

import { useEffect, useState, useRef } from 'react';
import protobuf from 'protobufjs';
import { Metric, MetricBatch } from '../types/telemetry';

export function useTelemetry() {
  const [metricsMap, setMetricsMap] = useState<Record<number, Metric>>({});
  const [isConnected, setIsConnected] = useState(false);
  const protoRoot = useRef<protobuf.Root | null>(null);

  useEffect(() => {
    let socket: WebSocket | null = null;

    protobuf.load("/proto/metrics.proto", (err, root) => {
      if (err || !root) return;
      protoRoot.current = root;
      
      socket = new WebSocket("ws://localhost:8080/ws");
      socket.binaryType = "arraybuffer";

      socket.onopen = () => setIsConnected(true);
      socket.onclose = () => setIsConnected(false);

      socket.onmessage = (event) => {
        if (!protoRoot.current) return;

        try {
          const MetricBatchType = protoRoot.current.lookupType("telemetry.MetricBatch");
          const buffer = new Uint8Array(event.data);
          const message = MetricBatchType.decode(buffer);
          const decoded = MetricBatchType.toObject(message, { 
            longs: Number, 
            defaults: true 
          }) as MetricBatch;

          if (decoded.metrics && decoded.metrics.length > 0) {
            setMetricsMap((prev) => {
              const newMap = { ...prev };
              decoded.metrics.forEach((m) => {
                newMap[m.id] = m;
              });
              return newMap;
            });
          }
        } catch (e) {
          console.error("Decoding error:", e);
        }
      };
    });

    return () => { socket?.close(); };
  }, []);

  return { 
    metrics: Object.values(metricsMap), 
    isConnected 
  };
}