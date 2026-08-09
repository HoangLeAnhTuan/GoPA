package migrations

import "embed"

// Files contains immutable, ordered SQL migrations embedded into the migration command.
//
//go:embed *.sql
var Files embed.FS
