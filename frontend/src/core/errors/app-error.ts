export class AppError extends Error {
  public readonly code: string;
  public readonly userMessage: string;

  constructor(message: string, code = "APP_ERROR", userMessage = "Ocurrió un error en la aplicación.") {
    super(message);
    this.name = "AppError";
    this.code = code;
    this.userMessage = userMessage;
  }
}

export class NetworkError extends AppError {
  constructor(message = "Network connection failed", originalError?: Error) {
    super(message, "NETWORK_ERROR", "No se pudo conectar con el servidor. Verificá tu conexión de red.");
    this.name = "NetworkError";
    if (originalError?.stack) {
      this.stack = originalError.stack;
    }
  }
}
