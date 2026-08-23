import React from "react";
import styles from "./Skeleton.module.css";

export interface SkeletonProps extends React.HTMLAttributes<HTMLDivElement> {
  width?: string | number;
  height?: string | number;
  borderRadius?: string | number;
}

export const Skeleton: React.FC<SkeletonProps> = ({
  width = "100%",
  height = "20px",
  borderRadius,
  style,
  className = "",
  ...props
}) => {
  const customStyle: React.CSSProperties = {
    width: typeof width === "number" ? `${width}px` : width,
    height: typeof height === "number" ? `${height}px` : height,
    borderRadius: borderRadius ? (typeof borderRadius === "number" ? `${borderRadius}px` : borderRadius) : undefined,
    ...style,
  };

  return (
    <div
      className={`${styles.skeleton} ${className}`.trim()}
      style={customStyle}
      aria-hidden="true"
      {...props}
    />
  );
};
