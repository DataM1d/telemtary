"use client";

import { Canvas } from "@react-three/fiber";
import { OrbitControls, Grid } from "@react-three/drei";
import { useTelemetry } from "../hooks/useTelemetry";
import ServerNode from "../components/ServerNode";

export default function TelemetryScene() {
  const { metrics, isConnected } = useTelemetry();

  return (
    <div className="h-screen w-screen bg-[#050505]">
      <Canvas camera={{ position: [15, 15, 15], fov: 50 }}>
        <color attach="background" args={["#050505"]} />
        
        <ambientLight intensity={0.5} />
        <pointLight position={[10, 10, 10]} intensity={1} />
        
        <Grid 
          infiniteGrid 
          fadeDistance={50} 
          fadeStrength={5} 
          sectionSize={3} 
          sectionColor="#333" 
          cellColor="#111" 
        />

        {metrics.map((metric) => (
      <ServerNode 
        key={metric.id} 
        metric={metric} 
        position={[
            (metric.id % 20) * 2 - 20,
              0, 
            Math.floor(metric.id / 20) * 2 - 10
            ]} 
          />
        ))}

        <OrbitControls makeDefault />
      </Canvas>

      <div className="absolute top-5 left-5 text-mono text-xs text-white bg-black/50 p-3 rounded border border-white/10">
        SYSTEM_STATUS: {isConnected ? "ONLINE" : "OFFLINE"} <br />
        ACTIVE_NODES: {metrics.length}
      </div>
    </div>
  );
}