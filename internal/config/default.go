package config

import "time"

type ctxKey string

const (
	UidKey ctxKey = "uid"
)

const (
	DefaultPage      = 1
	DefaultSize      = 40
	DefaultCacheTime = time.Hour
)

const ErrorSpanTag = "error"
