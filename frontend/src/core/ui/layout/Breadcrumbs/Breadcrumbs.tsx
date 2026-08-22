"use client";

import React from "react";
import Link from "next/link";
import { ChevronRight } from "lucide-react";
import styles from "./Breadcrumbs.module.css";

export interface BreadcrumbItem {
  label: string;
  href?: string;
}

export interface BreadcrumbsProps {
  items: BreadcrumbItem[];
  className?: string;
}

export const Breadcrumbs: React.FC<BreadcrumbsProps> = ({ items, className = "" }) => {
  return (
    <nav aria-label="Breadcrumb" className={`${styles.nav} ${className}`.trim()}>
      <ol className={styles.list}>
        {items.map((item, index) => {
          const isLast = index === items.length - 1;

          return (
            <li key={index} className={styles.item}>
              {index > 0 && (
                <span className={styles.separator} aria-hidden="true">
                  <ChevronRight size={14} />
                </span>
              )}
              {isLast || !item.href ? (
                <span className={styles.current} aria-current={isLast ? "page" : undefined}>
                  {item.label}
                </span>
              ) : (
                <Link href={item.href} className={styles.link}>
                  {item.label}
                </Link>
              )}
            </li>
          );
        })}
      </ol>
    </nav>
  );
};
