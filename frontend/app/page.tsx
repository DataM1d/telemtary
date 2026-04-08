"use client";
import { useTelemetry } from "./hooks/useTelemetry";

export default function Home() {
  const { metrics, isConnected } = useTelemetry();

  return (
    <main className="p-10 bg-black min-h-screen text-white font-mono">
      <div className="mb-4">
        Status: {isConnected ? "🟢 ONLINE" : "🔴 OFFLINE"}
      </div>
      <pre>{JSON.stringify(metrics, null, 2)}</pre>
    </main>
  );
}