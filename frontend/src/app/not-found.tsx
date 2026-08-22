import React from "react";
import Link from "next/link";
import { Compass, ArrowLeft } from "lucide-react";
import { Button } from "@/core/ui/primitives/Button";
import { Card, CardBody, CardFooter, CardHeader } from "@/core/ui/primitives/Card";
import { APP_ROUTES } from "@/core/routes/routes.config";

export default function NotFound() {
  return (
    <div
      style={{
        display: "flex",
        alignItems: "center",
        justifyContent: "center",
        minHeight: "70vh",
        padding: "24px",
      }}
    >
      <Card style={{ maxWidth: "480px", width: "100%" }}>
        <CardHeader style={{ display: "flex", alignItems: "center", gap: "12px" }}>
          <div
            style={{
              width: "32px",
              height: "32px",
              borderRadius: "4px",
              backgroundColor: "var(--surface-2)",
              color: "var(--text-primary)",
              display: "flex",
              alignItems: "center",
              justifyContent: "center",
            }}
          >
            <Compass size={20} />
          </div>
          <div>
            <h1 style={{ fontSize: "16px", fontWeight: 600, color: "var(--text-primary)" }}>
              Recurso No Encontrado
            </h1>
            <span
              style={{
                fontSize: "12px",
                fontFamily: "var(--font-mono)",
                color: "var(--text-dim)",
              }}
            >
              HTTP 404 / NOT_FOUND
            </span>
          </div>
        </CardHeader>

        <CardBody>
          <p style={{ fontSize: "14px", color: "var(--text-muted)", lineHeight: 1.5 }}>
            La ruta o recurso al que intentás acceder no existe o fue movido en el sistema de operaciones.
          </p>
        </CardBody>

        <CardFooter style={{ display: "flex", justifyContent: "flex-end" }}>
          <Link href={APP_ROUTES.OVERVIEW}>
            <Button variant="primary" size="md" leftIcon={<ArrowLeft size={16} />}>
              Volver al Dashboard
            </Button>
          </Link>
        </CardFooter>
      </Card>
    </div>
  );
}
