"use client";

import { useEffect, useState, useRef } from 'react';
import protobuf from 'protobufjs';
import { Metric, MetricBatch } from '../types/telemetry';

export function useTelemetry() {
  const [data, setData] = useState<Metric[]>([]);
  const [isConnected, setIsConnected] = useState(false);
  const protoRoot = useRef<protobuf.Root | null>(null);

  useEffect(() => {
    let socket: WebSocket | null = null;

    protobuf.load("/proto/metrics.proto", (err, root) => {
      if (err || !root) {
        console.error("Proto Load Error:", err);
        return;
      }
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

          setData(decoded.metrics || []);
        } catch (e) {
          console.error("Decoding error:", e);
        }
      };
    });

    return () => {
      if (socket) socket.close();
    };
  }, []);

  return { metrics: data, isConnected };
}