import type { Metadata } from "next";
import { Inter, JetBrains_Mono } from "next/font/google";
import { PublicEnvScript } from "next-runtime-env";
import { AuthProvider } from "@/features/auth/components/AuthProvider/AuthProvider";
import "./globals.css";

const inter = Inter({
  variable: "--font-inter",
  subsets: ["latin"],
  display: "swap",
});

const jetbrainsMono = JetBrains_Mono({
  variable: "--font-mono",
  subsets: ["latin"],
  display: "swap",
});

export const metadata: Metadata = {
  title: "Akritas — Autonomous Incident Response Platform",
  description:
    "Autonomous production incident remediation platform. Continuous monitoring, root cause investigation, and validated pull requests.",
  icons: {
    icon: "/favicon.ico",
  },
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en" className={`${inter.variable} ${jetbrainsMono.variable}`}>
      <head>
        <PublicEnvScript />
      </head>
      <body>
        <AuthProvider>{children}</AuthProvider>
      </body>
    </html>
  );
}
