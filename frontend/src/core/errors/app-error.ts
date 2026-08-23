export class AppError extends Error {
  public readonly code: string;
  public readonly userMessage: string;

  constructor(message: string, code = "APP_ERROR", userMessage = "An application error occurred.") {
    super(message);
    this.name = "AppError";
    this.code = code;
    this.userMessage = userMessage;
    Object.setPrototypeOf(this, AppError.prototype);
  }
}

export class NetworkError extends AppError {
  constructor(message = "Network connection failed", public readonly originalError?: Error) {
    super(message, "NETWORK_ERROR", "Could not connect to the server. Check your network connection.");
    this.name = "NetworkError";
    if (originalError?.stack) {
      this.stack = originalError.stack;
    }
  }
}
