import { useEffect, useRef, useState } from "react";
import { useQueryClient, useQuery } from "@tanstack/react-query";
import { getAccessToken } from "../../auth/api/auth-api";
import { listPomodoroHistory, type PomodoroState, type PomodoroHistoryItem } from "../api/pomodoro-api";

export type WebSocketConnectionStatus = "connected" | "connecting" | "disconnected";

export function usePomodoroWs() {
  const queryClient = useQueryClient();
  const [status, setStatus] = useState<WebSocketConnectionStatus>("disconnected");
  const wsRef = useRef<WebSocket | null>(null);
  const reconnectTimeoutRef = useRef<number | null>(null);
  const backoffRef = useRef(1000);

  useEffect(() => {
    let unmounted = false;

    function connect() {
      const token = getAccessToken();
      if (!token) {
        setStatus("disconnected");
        return;
      }

      const apiBase = import.meta.env.VITE_API_BASE_URL ?? "http://localhost:8080/api/v1";
      const wsBase = apiBase.replace(/^http/, "ws");
      const wsUrl = `${wsBase}/ws/pomodoro?token=${encodeURIComponent(token)}`;

      setStatus("connecting");
      const ws = new WebSocket(wsUrl);
      wsRef.current = ws;

      ws.onopen = () => {
        if (unmounted) {
          ws.close();
          return;
        }
        setStatus("connected");
        backoffRef.current = 1000;
      };

      ws.onmessage = (event) => {
        try {
          const message = JSON.parse(event.data);
          if (message.type === "state") {
            const newState: PomodoroState | null = message.data ?? null;
            queryClient.setQueryData(["pomodoro"], newState);
          }
        } catch {
          // Ignore malformed message
        }
      };

      ws.onclose = () => {
        if (unmounted) return;
        setStatus("disconnected");
        wsRef.current = null;
        // Exponential backoff reconnect
        const nextDelay = Math.min(backoffRef.current * 1.5, 15000);
        backoffRef.current = nextDelay;
        reconnectTimeoutRef.current = window.setTimeout(() => {
          if (!unmounted) connect();
        }, nextDelay);
      };

      ws.onerror = () => {
        ws.close();
      };
    }

    connect();

    return () => {
      unmounted = true;
      if (reconnectTimeoutRef.current !== null) {
        window.clearTimeout(reconnectTimeoutRef.current);
      }
      if (wsRef.current) {
        wsRef.current.close();
      }
    };
  }, [queryClient]);

  return { status };
}

export function usePomodoroHistory(limit = 20) {
  return useQuery<PomodoroHistoryItem[]>({
    queryKey: ["pomodoro", "history", limit],
    queryFn: () => listPomodoroHistory(limit),
  });
}
