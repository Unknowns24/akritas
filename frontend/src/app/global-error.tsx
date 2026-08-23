"use client";

import React from "react";
import "./globals.css";

interface GlobalErrorProps {
  error: Error & { digest?: string };
  reset: () => void;
}

export default function GlobalError({ error, reset }: GlobalErrorProps) {
  return (
    <html lang="en">
      <body
        style={{
          backgroundColor: "#09090b",
          color: "#ffffff",
          fontFamily: "Inter, -apple-system, sans-serif",
          display: "flex",
          alignItems: "center",
          justifyContent: "center",
          minHeight: "100vh",
          margin: 0,
          padding: "24px",
        }}
      >
        <div
          style={{
            maxWidth: "480px",
            width: "100%",
            backgroundColor: "#18181b",
            border: "1px solid #27272a",
            borderRadius: "8px",
            padding: "24px",
            display: "flex",
            flexDirection: "column",
            gap: "16px",
            boxShadow: "0 8px 24px rgba(0, 0, 0, 0.5)",
          }}
        >
          <div style={{ display: "flex", alignItems: "center", gap: "12px" }}>
            <span style={{ fontSize: "24px" }}>⚠️</span>
            <div>
              <h1 style={{ fontSize: "18px", fontWeight: 600, margin: 0 }}>
                Error Crítico del Sistema
              </h1>
              <span
                style={{
                  fontSize: "12px",
                  fontFamily: "monospace",
                  color: "#8e9192",
                }}
              >
                {error.digest ? `ID: ${error.digest}` : "FATAL_APPLICATION_ERROR"}
              </span>
            </div>
          </div>

          <p style={{ fontSize: "14px", color: "#c4c7c8", lineHeight: 1.5, margin: 0 }}>
            Se produjo un error irrecuperable en la capa raíz de la aplicación.
          </p>

          <div style={{ display: "flex", justifyContent: "flex-end", marginTop: "8px" }}>
            <button
              type="button"
              onClick={reset}
              style={{
                backgroundColor: "#ffffff",
                color: "#09090b",
                border: "none",
                borderRadius: "4px",
                padding: "8px 16px",
                fontSize: "14px",
                fontWeight: 600,
                cursor: "pointer",
              }}
            >
              Reiniciar Aplicación
            </button>
          </div>
        </div>
      </body>
    </html>
  );
}
