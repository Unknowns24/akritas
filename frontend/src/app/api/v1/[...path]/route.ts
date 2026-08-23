import http from "node:http";
import https from "node:https";
import { NextRequest } from "next/server";

export const dynamic = "force-dynamic";
export const runtime = "nodejs";

type RouteContext = {
  params: Promise<{
    path: string[];
  }>;
};

const DEFAULT_PROXY_TARGET = "https://localhost:8443";

const HOP_BY_HOP_HEADERS = new Set([
  "connection",
  "keep-alive",
  "proxy-authenticate",
  "proxy-authorization",
  "te",
  "trailer",
  "transfer-encoding",
  "upgrade",
]);

function getProxyTarget(): URL {
  const rawTarget = process.env.AKRITAS_API_PROXY_TARGET || DEFAULT_PROXY_TARGET;
  return new URL(rawTarget.replace(/\/$/, ""));
}

function allowsInsecureLocalTLS(url: URL): boolean {
  return (
    url.protocol === "https:" &&
    (url.hostname === "localhost" ||
      url.hostname === "127.0.0.1" ||
      url.hostname === "::1")
  );
}

function copyRequestHeaders(request: NextRequest): Record<string, string> {
  const headers: Record<string, string> = {};

  request.headers.forEach((value, key) => {
    const normalizedKey = key.toLowerCase();
    if (
      normalizedKey === "host" ||
      normalizedKey === "content-length" ||
      HOP_BY_HOP_HEADERS.has(normalizedKey)
    ) {
      return;
    }
    headers[key] = value;
  });

  return headers;
}

function copyResponseHeaders(headers: http.IncomingHttpHeaders): Headers {
  const responseHeaders = new Headers();

  for (const [key, value] of Object.entries(headers)) {
    const normalizedKey = key.toLowerCase();
    if (!value || HOP_BY_HOP_HEADERS.has(normalizedKey)) {
      continue;
    }

    if (Array.isArray(value)) {
      for (const item of value) {
        responseHeaders.append(key, item);
      }
      continue;
    }

    responseHeaders.set(key, value);
  }

  return responseHeaders;
}

async function getRequestBody(request: NextRequest): Promise<Buffer | undefined> {
  if (request.method === "GET" || request.method === "HEAD") {
    return undefined;
  }

  const body = await request.arrayBuffer();
  return Buffer.from(body);
}

async function proxy(request: NextRequest, context: RouteContext): Promise<Response> {
  const target = getProxyTarget();
  const { path } = await context.params;
  const destination = new URL(`${target.pathname.replace(/\/$/, "")}/api/v1/${path.join("/")}`, target);
  destination.search = request.nextUrl.search;

  const body = await getRequestBody(request);
  const transport = destination.protocol === "https:" ? https : http;

  return new Promise<Response>((resolve, reject) => {
    const proxyRequest = transport.request(
      destination,
      {
        method: request.method,
        headers: copyRequestHeaders(request),
        rejectUnauthorized: !allowsInsecureLocalTLS(destination),
      },
      (proxyResponse) => {
        const chunks: Buffer[] = [];

        proxyResponse.on("data", (chunk: Buffer) => {
          chunks.push(chunk);
        });

        proxyResponse.on("end", () => {
          resolve(
            new Response(Buffer.concat(chunks), {
              status: proxyResponse.statusCode || 502,
              headers: copyResponseHeaders(proxyResponse.headers),
            }),
          );
        });
      },
    );

    proxyRequest.on("error", reject);

    if (body) {
      proxyRequest.write(body);
    }

    proxyRequest.end();
  });
}

export async function GET(request: NextRequest, context: RouteContext) {
  return proxy(request, context);
}

export async function POST(request: NextRequest, context: RouteContext) {
  return proxy(request, context);
}

export async function PUT(request: NextRequest, context: RouteContext) {
  return proxy(request, context);
}

export async function PATCH(request: NextRequest, context: RouteContext) {
  return proxy(request, context);
}

export async function DELETE(request: NextRequest, context: RouteContext) {
  return proxy(request, context);
}

export async function OPTIONS(request: NextRequest, context: RouteContext) {
  return proxy(request, context);
}
