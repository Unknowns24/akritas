import React from "react";
import { redirect } from "next/navigation";
import { AppShell } from "@/core/ui/layout/AppShell";
import { isApiUnauthorizedError } from "@/core/errors";
import { APP_ROUTES } from "@/core/routes/routes.config";
import { getCurrentSessionService } from "@/features/auth/services/auth.service";

export default async function PrivateLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  try {
    await getCurrentSessionService();
  } catch (error) {
    if (isApiUnauthorizedError(error)) {
      redirect(APP_ROUTES.AUTH.LOGIN);
    }

    throw error;
  }

  return <AppShell>{children}</AppShell>;
}
