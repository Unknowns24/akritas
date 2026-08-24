import { ApiError } from "./api-error";
import { AppError } from "./app-error";

function stringFromRecord(value: Record<string, unknown>, key: string): string | undefined {
  const candidate = value[key];
  return typeof candidate === "string" && candidate.trim() ? candidate : undefined;
}

export function getErrorMessage(error: unknown, fallback: string): string {
  if (error instanceof ApiError || error instanceof AppError) {
    return error.userMessage;
  }

  if (error instanceof Error) {
    return error.message;
  }

  if (typeof error === "object" && error !== null) {
    const record = error as Record<string, unknown>;
    const nested = record.error;

    if (typeof nested === "object" && nested !== null) {
      const nestedRecord = nested as Record<string, unknown>;
      return (
        stringFromRecord(nestedRecord, "user_message") ||
        stringFromRecord(nestedRecord, "message") ||
        fallback
      );
    }

    return (
      stringFromRecord(record, "user_message") ||
      stringFromRecord(record, "message") ||
      fallback
    );
  }

  return fallback;
}

export function getErrorDetailsMessage(error: unknown): string {
  if (error instanceof ApiError) {
    return error.details
      .map((detail) => (detail.field ? `${detail.field}: ${detail.reason}` : detail.reason))
      .join(", ");
  }

  if (typeof error !== "object" || error === null) {
    return "";
  }

  const record = error as Record<string, unknown>;
  const maybeEnvelope = record.error;
  const source =
    typeof maybeEnvelope === "object" && maybeEnvelope !== null
      ? (maybeEnvelope as Record<string, unknown>)
      : record;
  const details = source.details;

  if (!Array.isArray(details)) {
    return "";
  }

  return details
    .map((detail) => {
      if (typeof detail !== "object" || detail === null) {
        return "";
      }

      const detailRecord = detail as Record<string, unknown>;
      const field = stringFromRecord(detailRecord, "field");
      const reason = stringFromRecord(detailRecord, "reason");

      if (!reason) {
        return "";
      }

      return field ? `${field}: ${reason}` : reason;
    })
    .filter(Boolean)
    .join(", ");
}

export function isApiNotFoundError(error: unknown): boolean {
  const message = getErrorMessage(error, "").toLowerCase();
  return message.includes("404") || message.includes("page not found") || message.includes("not found");
}

export function isApiUnauthorizedError(error: unknown): boolean {
  if (error instanceof ApiError) {
    return error.status === 401;
  }

  if (typeof error !== "object" || error === null) {
    return false;
  }

  const record = error as Record<string, unknown>;
  if (record.status === 401) {
    return true;
  }

  const source =
    typeof record.error === "object" && record.error !== null
      ? (record.error as Record<string, unknown>)
      : record;
  const code = stringFromRecord(source, "code");
  const message = getErrorMessage(error, "").toLowerCase();

  return code?.endsWith("U") === true || message.includes("unauthorized");
}
