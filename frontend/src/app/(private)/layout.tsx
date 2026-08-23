import React from "react";
import { AppShell } from "@/core/ui/layout/AppShell";

export default function PrivateLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return <AppShell>{children}</AppShell>;
}
