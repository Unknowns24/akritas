import React from "react";
import { Skeleton } from "@/core/ui/primitives/Skeleton";

export default function Loading() {
  return (
    <div style={{ display: "flex", flexDirection: "column", gap: "24px", width: "100%" }}>
      {/* Header skeleton */}
      <div style={{ display: "flex", justifyContent: "space-between", alignItems: "center" }}>
        <div style={{ display: "flex", flexDirection: "column", gap: "8px" }}>
          <Skeleton width={220} height={28} />
          <Skeleton width={340} height={16} />
        </div>
        <Skeleton width={120} height={36} borderRadius={4} />
      </div>

      {/* Metric Cards Skeleton Grid */}
      <div
        style={{
          display: "grid",
          gridTemplateColumns: "repeat(auto-fit, minmax(240px, 1fr))",
          gap: "16px",
        }}
      >
        <Skeleton height={110} borderRadius={8} />
        <Skeleton height={110} borderRadius={8} />
        <Skeleton height={110} borderRadius={8} />
        <Skeleton height={110} borderRadius={8} />
      </div>

      {/* Main Content Area Skeleton */}
      <div
        style={{
          display: "grid",
          gridTemplateColumns: "2fr 1fr",
          gap: "24px",
          marginTop: "12px",
        }}
      >
        <Skeleton height={380} borderRadius={8} />
        <Skeleton height={380} borderRadius={8} />
      </div>
    </div>
  );
}
