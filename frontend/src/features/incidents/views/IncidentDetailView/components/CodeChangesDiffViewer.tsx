import React from "react";
import { FileCode } from "lucide-react";
import type { CodeChange } from "../../../types/remediation.types";
import styles from "./CodeChangesDiffViewer.module.css";

interface CodeChangesDiffViewerProps {
  changes?: CodeChange[];
}

export function CodeChangesDiffViewer({ changes }: CodeChangesDiffViewerProps) {
  if (!changes || changes.length === 0) {
    return null;
  }

  return (
    <div className={styles.container}>
      {changes.map((change, index) => {
        const lines = change.patch ? change.patch.split("\n") : [];
        const changeClass =
          change.change_type === "added"
            ? styles.changeAdded
            : change.change_type === "deleted"
            ? styles.changeDeleted
            : styles.changeModified;

        return (
          <div key={index} className={styles.fileBlock}>
            <div className={styles.fileHeader}>
              <span className={styles.filePath}>
                <FileCode size={13} />
                {change.file_path}
              </span>
              <span className={`${styles.changeTypeTag} ${changeClass}`}>
                {change.change_type}
              </span>
            </div>

            {lines.length > 0 && (
              <div className={styles.diffContent}>
                {lines.map((line, lineIdx) => {
                  const isMeta = line.startsWith("@@");
                  const isAdded = line.startsWith("+") && !isMeta;
                  const isRemoved = line.startsWith("-") && !isMeta;

                  let lineStyle = styles.diffLine;
                  if (isMeta) lineStyle += ` ${styles.diffMeta}`;
                  else if (isAdded) lineStyle += ` ${styles.diffAdded}`;
                  else if (isRemoved) lineStyle += ` ${styles.diffRemoved}`;
                  else lineStyle += ` ${styles.diffUnchanged}`;

                  return (
                    <div key={lineIdx} className={lineStyle}>
                      {line}
                    </div>
                  );
                })}
              </div>
            )}
          </div>
        );
      })}
    </div>
  );
}
