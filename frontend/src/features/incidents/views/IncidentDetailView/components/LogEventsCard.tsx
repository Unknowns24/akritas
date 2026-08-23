"use client";

import React, { useEffect, useState, useCallback } from "react";
import { ChevronLeft, ChevronRight, List } from "lucide-react";
import { getErrorMessage } from "@/core/errors";
import { getIncidentLogEventsService } from "../../../services/get-incident-log-events.service";
import type { LogEvent } from "../../../services/get-incident-log-events.service";
import styles from "../IncidentDetailView.module.css";

export function LogEventsCard({ incidentId }: { incidentId: string }) {
  const [events, setEvents] = useState<LogEvent[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  // Pagination state
  const [nextCursor, setNextCursor] = useState<string>("");
  const [hasMore, setHasMore] = useState(false);
  const [cursorStack, setCursorStack] = useState<string[]>([]);
  
  const fetchEvents = useCallback(async (cursor?: string) => {
    setLoading(true);
    setError(null);
    try {
      const response = await getIncidentLogEventsService(incidentId, { limit: 10, cursor });
      setEvents(response.data || []);
      setNextCursor(response.paging?.next_cursor || "");
      setHasMore(response.paging?.has_more || false);
    } catch (error: unknown) {
      setError(getErrorMessage(error, "Failed to load log events"));
    } finally {
      setLoading(false);
    }
  }, [incidentId]);

  useEffect(() => {
    const timeoutId = window.setTimeout(() => {
      void fetchEvents();
    }, 0);

    return () => window.clearTimeout(timeoutId);
  }, [fetchEvents]);

  const handleNextPage = () => {
    if (hasMore && nextCursor) {
      setCursorStack((prev) => [...prev, nextCursor]);
      fetchEvents(nextCursor);
    }
  };

  const handlePrevPage = () => {
    setCursorStack((prev) => {
      const newStack = [...prev];
      newStack.pop(); // Remove current cursor
      const prevC = newStack.length > 0 ? newStack[newStack.length - 1] : undefined;
      fetchEvents(prevC);
      return newStack;
    });
  };

  if (loading && events.length === 0) {
    return (
      <div className={styles.card}>
        <div className={styles.cardHeader}>
          <List size={16} />
          Log Evidence
        </div>
        <div className={styles.loadingState}>Loading log events...</div>
      </div>
    );
  }

  if (error) {
    return (
      <div className={styles.card}>
        <div className={styles.cardHeader}>
          <List size={16} />
          Log Evidence
        </div>
        <div className={styles.errorState}>{error}</div>
      </div>
    );
  }

  if (events.length === 0) {
    return null;
  }

  return (
    <div className={styles.card}>
      <div className={styles.cardHeader}>
        <List size={16} />
        Log Evidence
      </div>
      
      <div className={styles.logEventsList}>
        {events.map((event) => (
          <div key={event.id} className={styles.logEventContainer}>
            <div className={styles.logEventHeader}>
              <span className={styles.logEventTimestamp}>
                {new Date(event.timestamp).toLocaleString()}
              </span>
              <span className={`${styles.logEventSeverity} ${styles[event.severity.toLowerCase()] || ''}`}>
                {event.severity}
              </span>
            </div>

            {/* Context Before */}
            {event.context_before && event.context_before.length > 0 && (
              <div className={styles.logContextSection}>
                {event.context_before.map((ctx, idx) => (
                  <div key={`before-${idx}`} className={styles.logContextLine}>
                    <span className={styles.logContextTime}>
                      {new Date(ctx.timestamp).toLocaleTimeString()}
                    </span>
                    <span className={styles.logContextMessage}>{ctx.message}</span>
                  </div>
                ))}
              </div>
            )}

            {/* Main Evidence */}
            <div className={styles.logEvidence}>
              <pre>{event.message}</pre>
            </div>

            {/* Context After */}
            {event.context_after && event.context_after.length > 0 && (
              <div className={styles.logContextSection}>
                {event.context_after.map((ctx, idx) => (
                  <div key={`after-${idx}`} className={styles.logContextLine}>
                    <span className={styles.logContextTime}>
                      {new Date(ctx.timestamp).toLocaleTimeString()}
                    </span>
                    <span className={styles.logContextMessage}>{ctx.message}</span>
                  </div>
                ))}
              </div>
            )}
          </div>
        ))}
      </div>

      <div className={styles.paginationControls}>
        <button
          className={styles.paginationButton}
          onClick={handlePrevPage}
          disabled={cursorStack.length === 0 || loading}
        >
          <ChevronLeft size={16} />
          Previous
        </button>
        <button
          className={styles.paginationButton}
          onClick={handleNextPage}
          disabled={!hasMore || loading}
        >
          Next
          <ChevronRight size={16} />
        </button>
      </div>
    </div>
  );
}
