import React from "react";
import styles from "./LoadingState.module.css";
import { Loader2 } from "lucide-react";

export interface LoadingStateProps {
  message?: string;
  className?: string;
  fullHeight?: boolean;
}

export const LoadingState: React.FC<LoadingStateProps> = ({
  message = "Loading...",
  className = "",
  fullHeight = false,
}) => {
  return (
    <div className={`${styles.loadingState} ${fullHeight ? styles.fullHeight : ""} ${className}`.trim()}>
      <Loader2 className={styles.spinner} />
      {message && <p className={styles.message}>{message}</p>}
    </div>
  );
};
