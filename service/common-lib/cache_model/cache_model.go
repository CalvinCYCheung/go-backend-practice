package cachemodel

import "time"

type CacheModel[T any] struct {
	Data           T
	SoftExpireTime time.Time
	HardExpireTime time.Time
	Refreshing     bool
}
