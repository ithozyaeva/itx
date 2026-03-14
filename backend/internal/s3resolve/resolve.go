package s3resolve

import (
	"sync"
	"time"
)

// Resolver is a function that resolves S3 keys to presigned URLs.
// Set by utils.InitGlobalS3() at startup to break the import cycle.
var Resolver func(stored string) string

type cacheEntry struct {
	url       string
	expiresAt time.Time
}

var (
	cache    = make(map[string]cacheEntry)
	cacheMu  sync.RWMutex
	cacheTTL = 6 * time.Hour
)

// ResolveS3URL resolves a stored S3 key or old-style URL to a presigned URL.
// Results are cached to avoid repeated S3 presign API calls.
// Returns the input unchanged if no resolver is registered.
func ResolveS3URL(stored string) string {
	if Resolver == nil || stored == "" {
		return stored
	}

	cacheMu.RLock()
	if entry, ok := cache[stored]; ok && time.Now().Before(entry.expiresAt) {
		cacheMu.RUnlock()
		return entry.url
	}
	cacheMu.RUnlock()

	url := Resolver(stored)

	cacheMu.Lock()
	cache[stored] = cacheEntry{url: url, expiresAt: time.Now().Add(cacheTTL)}
	cacheMu.Unlock()

	return url
}
