package frontend

import (
	"context"
	"time"
)

func (frontend *frontend) CacheHas(key string) bool {
	if !frontend.cacheEnabled {
		return false
	}

	if frontend.cache == nil {
		return false
	}

	return frontend.cache.Has(key)
}

func (frontend *frontend) CacheGet(key string) any {
	if !frontend.cacheEnabled {
		return nil
	}

	if frontend.cache == nil {
		return nil
	}

	item := frontend.cache.Get(key)

	if item == nil {
		return nil
	}

	return item.Value()
}

func (frontend *frontend) CacheSet(key string, value any, expireSeconds int) {
	if !frontend.cacheEnabled {
		return
	}

	if frontend.cache == nil {
		return
	}
	frontend.cache.Set(key, value, time.Duration(expireSeconds)*time.Second)
}

// warmUpCache periodically fetches the active sites and stores them in the cache
// to avoid an extra database query every time a request comes in to the frontend
// handler
func (frontend *frontend) warmUpCache() error {

	frontend.fetchActiveSites(context.Background())

	for range time.Tick(time.Second * 60) {
		frontend.warmUpCache()
	}
	return nil
}
