package constants

// Internal-only infrastructure, broker, and cache constants

// Internal Context Keys
const (
	InternalRequestCtxKey = "internal_request_ctx"
	TraceSpanCtxKey       = "trace_span"
)

// Internal Cache Prefixes
const (
	CacheAuthRefreshPrefix = "auth:refresh:"
	CacheRateLimitPrefix   = "rate:"
	CachePomodoroStateKey  = "pomodoro:state:"
)

// Internal Queue & Exchange Names
const (
	QueueSRSNotifications = "gopa.queue.srs.notifications"
	QueuePomodoroHistory  = "gopa.queue.pomodoro.history"
	ExchangeEvents        = "gopa.exchange.events"
)
