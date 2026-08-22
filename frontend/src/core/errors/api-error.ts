export interface ErrorDetailDto {
  field: string;
  reason: string;
}

export interface ErrorDto {
  code: string;
  message: string;
  user_message: string;
  request_id: string;
  details?: ErrorDetailDto[];
}

export interface ErrorResponseDto {
  error: ErrorDto;
}

export class ApiError extends Error {
  public readonly code: string;
  public readonly userMessage: string;
  public readonly requestId: string;
  public readonly status: number;
  public readonly details: ErrorDetailDto[];

  constructor(status: number, errorDto: Partial<ErrorDto>, originalError?: Error) {
    const message = errorDto.message || "An unexpected API error occurred.";
    super(message);
    this.name = "ApiError";
    this.status = status;
    this.code = errorDto.code || "UNKNOWN_ERROR";
    this.userMessage = errorDto.user_message || "Ocurrió un error al procesar la solicitud.";
    this.requestId = errorDto.request_id || "req-unknown";
    this.details = errorDto.details || [];

    if (originalError?.stack) {
      this.stack = originalError.stack;
    }
  }

  public static fromErrorResponse(status: number, responseBody: unknown): ApiError {
    if (
      typeof responseBody === "object" &&
      responseBody !== null &&
      "error" in responseBody &&
      typeof (responseBody as Record<string, unknown>).error === "object"
    ) {
      const errorDto = (responseBody as ErrorResponseDto).error;
      return new ApiError(status, errorDto);
    }

    return new ApiError(status, {
      code: `HTTP_${status}`,
      message: `HTTP request failed with status ${status}`,
      user_message: "Error de comunicación con el servidor.",
      request_id: `req-${Date.now()}`,
      details: [],
    });
  }
}
