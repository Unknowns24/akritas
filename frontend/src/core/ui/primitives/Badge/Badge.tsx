import React from "react";
import styles from "./Badge.module.css";

export type BadgeVariant =
  | "neutral"
  | "error"
  | "warning"
  | "success"
  | "info"
  | "intel";

export interface BadgeProps extends React.HTMLAttributes<HTMLSpanElement> {
  variant?: BadgeVariant;
  size?: "sm" | "md";
  dot?: boolean;
}

export const Badge: React.FC<BadgeProps> = ({
  children,
  variant = "neutral",
  size = "sm",
  dot = false,
  className = "",
  ...props
}) => {
  const variantClass =
    styles[`variant-${variant}`] || styles["variant-neutral"];
  const sizeClass = styles[`size-${size}`] || styles["size-sm"];

  return (
    <span
      className={`${styles.badge} ${variantClass} ${sizeClass} ${className}`.trim()}
      {...props}
    >
      {dot && <span className={styles.dot} aria-hidden="true" />}
      {children}
    </span>
  );
};
