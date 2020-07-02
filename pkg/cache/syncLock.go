package cache

import "time"

//SyncLocker a sync locker interface
type SyncLocker interface {
	GetLock(key string, expire time.Duration) bool
	ReleaseLock(key string) error
}
