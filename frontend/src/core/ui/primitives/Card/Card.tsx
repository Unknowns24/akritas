import React from "react";
import styles from "./Card.module.css";

export interface CardProps extends React.HTMLAttributes<HTMLDivElement> {
  level?: 1 | 2;
  accent?: "none" | "indigo" | "error" | "warning" | "success";
}

export const Card: React.FC<CardProps> = ({
  children,
  level = 1,
  accent = "none",
  className = "",
  ...props
}) => {
  const levelClass = styles[`level-${level}`] || styles["level-1"];
  const accentClass = accent !== "none" ? styles[`accent-${accent}`] : "";

  return (
    <div className={`${styles.card} ${levelClass} ${accentClass} ${className}`.trim()} {...props}>
      {children}
    </div>
  );
};

export const CardHeader: React.FC<React.HTMLAttributes<HTMLDivElement>> = ({
  children,
  className = "",
  ...props
}) => (
  <div className={`${styles.header} ${className}`.trim()} {...props}>
    {children}
  </div>
);

export const CardBody: React.FC<React.HTMLAttributes<HTMLDivElement>> = ({
  children,
  className = "",
  ...props
}) => (
  <div className={`${styles.body} ${className}`.trim()} {...props}>
    {children}
  </div>
);

export const CardFooter: React.FC<React.HTMLAttributes<HTMLDivElement>> = ({
  children,
  className = "",
  ...props
}) => (
  <div className={`${styles.footer} ${className}`.trim()} {...props}>
    {children}
  </div>
);
