"use client";

import { useRef } from "react";
import { useFrame } from "@react-three/fiber";
import * as THREE from "three";
import { Metric } from "../types/telemetry";

interface Props {
    metric: Metric;
    position: [number, number, number];
}

export default function ServerNode({ metric, position }: Props) {
    const meshRef = useRef<THREE.Mesh>(null);

    useFrame(() => {
        if (meshRef.current) {
            const targetScale = Math.max(0.5, metric.value / 10);
            meshRef.current.scale.y = THREE.MathUtils.lerp(
                meshRef.current.scale.y,
                targetScale,
                0.1
            );
        }
    });

    const color = new THREE.Color().setHSL(0.3 - (metric.value / 300), 1, 0.5);

    return (
        <mesh ref={meshRef} position={[position[0], position[1], position[2]]}>
            <boxGeometry args={[1, 1, 1]} />
            <meshStandardMaterial color={color} emissive={color} emissiveIntensity={0.5} />
         </mesh>
    );
}