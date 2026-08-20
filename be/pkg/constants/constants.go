package constants

// Environment constants
const (
	Production  = "production"
	Development = "development"
	Test        = "test"
	Debug       = "debug"

	PRODUCTION  = "production"
	DEVELOPMENT = "development"
	TEST        = "test"
	DEBUG       = "debug"
)

// Primitive utility constants (String & Integer)
const (
	ONE_STRING  = "1"
	ZERO_STRING = "0"

	ONE_INT    = 1
	ZERO_INT   = 0
	ONE_INT32  = int32(1)
	ZERO_INT32 = int32(0)
	ONE_INT64  = int64(1)
	ZERO_INT64 = int64(0)
)

// HTTP & Auth Context/Header/Cookie keys
const (
	RefreshCookieName = "gopa_refresh"
	RequestIDKey      = "request_id"
	IdentityKey       = "identity"
	RequestIDHeader   = "X-Request-ID"
)

// Example configuration defaults
const (
	ExampleJWTSecret   = "replace-with-a-long-random-secret"
	DefaultBcryptCost  = 12
	DefaultPageSize    = 20
	MaxPageSize        = 100
	DefaultHttpTimeout = 10
)
