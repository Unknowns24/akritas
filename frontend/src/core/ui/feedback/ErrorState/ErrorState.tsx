import React from "react";
import { AlertCircle, RefreshCw, ChevronRight } from "lucide-react";
import { Button } from "@/core/ui/primitives/Button";
import { ApiError } from "@/core/errors/api-error";
import styles from "./ErrorState.module.css";

export interface ErrorStateProps {
  error: ApiError | Error;
  onRetry?: () => void;
  className?: string;
  title?: string;
}

export const ErrorState: React.FC<ErrorStateProps> = ({
  error,
  onRetry,
  className = "",
  title = "Algo salió mal",
}) => {
  const isApiError = error instanceof ApiError;
  const userMessage = isApiError ? error.userMessage : error.message;
  
  return (
    <div className={`${styles.errorState} ${className}`.trim()}>
      <div className={styles.iconWrapper}>
        <AlertCircle size={32} />
      </div>
      
      <h3 className={styles.title}>{title}</h3>
      <p className={styles.message}>{userMessage}</p>
      
      {onRetry && (
        <Button 
          variant="secondary" 
          onClick={onRetry} 
          leftIcon={<RefreshCw size={16} />}
          className={styles.retryButton}
        >
          Reintentar
        </Button>
      )}

      {isApiError && (
        <details className={styles.detailsBlock}>
          <summary className={styles.detailsSummary}>
            Detalles avanzados
          </summary>
          <div className={styles.detailsContent}>
            <div className={styles.detailItem}>
              <strong>Request ID:</strong> {error.requestId}
            </div>
            <div className={styles.detailItem}>
              <strong>Código:</strong> {error.code}
            </div>
            {error.details && error.details.length > 0 && (
              <div className={styles.detailList}>
                <strong>Errores:</strong>
                <ul>
                  {error.details.map((detail, index) => (
                    <li key={index}>
                      <code>{detail.field}</code>: {detail.reason}
                    </li>
                  ))}
                </ul>
              </div>
            )}
            <div className={styles.detailItem}>
              <strong>Mensaje Técnico:</strong> {error.message}
            </div>
          </div>
        </details>
      )}
    </div>
  );
};
