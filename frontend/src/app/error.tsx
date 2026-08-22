"use client";

import React, { useEffect } from "react";
import Link from "next/link";
import { AlertOctagon, RotateCcw, Home } from "lucide-react";
import { Button } from "@/core/ui/primitives/Button";
import { Card, CardBody, CardFooter, CardHeader } from "@/core/ui/primitives/Card";
import { ApiError } from "@/core/errors/api-error";
import { APP_ROUTES } from "@/core/routes/routes.config";

interface ErrorProps {
  error: Error & { digest?: string };
  reset: () => void;
}

export default function ErrorBoundary({ error, reset }: ErrorProps) {
  useEffect(() => {
    // Log unexpected runtime error without exposing credentials
    console.error("[ErrorBoundary caught]:", error);
  }, [error]);

  const isApiError = error instanceof ApiError;
  const errorCode = isApiError ? error.code : error.digest || "RUNTIME_ERROR";
  const userMessage = isApiError
    ? error.userMessage
    : "Ocurrió un error inesperado al procesar la vista. Podés intentar nuevamente.";
  const requestId = isApiError ? error.requestId : null;

  return (
    <div
      style={{
        display: "flex",
        alignItems: "center",
        justifyContent: "center",
        minHeight: "60vh",
        padding: "24px",
      }}
    >
      <Card
        accent="error"
        style={{
          maxWidth: "520px",
          width: "100%",
        }}
      >
        <CardHeader style={{ display: "flex", alignItems: "center", gap: "12px" }}>
          <div
            style={{
              width: "32px",
              height: "32px",
              borderRadius: "4px",
              backgroundColor: "var(--status-error-bg)",
              color: "var(--status-error-light)",
              display: "flex",
              alignItems: "center",
              justifyContent: "center",
            }}
          >
            <AlertOctagon size={20} />
          </div>
          <div>
            <h2 style={{ fontSize: "16px", fontWeight: 600, color: "var(--text-primary)" }}>
              Error en la aplicación
            </h2>
            <span
              style={{
                fontSize: "12px",
                fontFamily: "var(--font-mono)",
                color: "var(--text-dim)",
              }}
            >
              CODE: {errorCode}
            </span>
          </div>
        </CardHeader>

        <CardBody style={{ display: "flex", flexDirection: "column", gap: "12px" }}>
          <p style={{ fontSize: "14px", color: "var(--text-secondary)", lineHeight: 1.5 }}>
            {userMessage}
          </p>

          {requestId && (
            <div
              style={{
                backgroundColor: "var(--surface-2)",
                padding: "8px 12px",
                borderRadius: "4px",
                fontSize: "12px",
                fontFamily: "var(--font-mono)",
                color: "var(--text-dim)",
              }}
            >
              Request ID: {requestId}
            </div>
          )}
        </CardBody>

        <CardFooter style={{ display: "flex", gap: "12px", justifyContent: "flex-end" }}>
          <Link href={APP_ROUTES.OVERVIEW}>
            <Button variant="ghost" size="md" leftIcon={<Home size={16} />}>
              Dashboard
            </Button>
          </Link>
          <Button variant="primary" size="md" leftIcon={<RotateCcw size={16} />} onClick={reset}>
            Reintentar
          </Button>
        </CardFooter>
      </Card>
    </div>
  );
}
