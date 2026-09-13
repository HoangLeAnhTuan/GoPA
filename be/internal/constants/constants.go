package constants

// Internal-only infrastructure, broker, and cache constants

// Internal Context Keys
const (
	InternalRequestCtxKey = "internal_request_ctx"
	TraceSpanCtxKey       = "trace_span"
)

// Internal Cache Prefixes
const (
	CacheAuthRefreshPrefix     = "auth:refresh:"
	CacheRateLimitPrefix       = "rate:"
	CachePomodoroStateKey      = "pomodoro:state:"
	CachePomodoroChannelPrefix = "pomodoro:channel:"
)

// Internal Queue & Exchange Names
const (
	QueueSRSNotifications = "gopa.queue.srs.notifications"
	QueuePomodoroHistory  = "gopa.queue.pomodoro.history"
	QueueSRS              = "gopa.queue.srs"
	QueuePomodoro         = "gopa.queue.pomodoro"
	QueueFinance          = "gopa.queue.finance"
	QueueJournal          = "gopa.queue.journal"
	QueueDLQ              = "gopa.queue.dlq"
	ExchangeEvents        = "gopa.events"
	ExchangeDLX           = "gopa.events.dlx"
)
