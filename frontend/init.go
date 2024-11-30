package frontend

import (
	"time"

	"github.com/jellydator/ttlcache/v3"
	"github.com/mingrammer/cfmt"
)

type LanguageKey struct{}

func init() {}

func initCache() *ttlcache.Cache[string, any] {
	cfmt.Successln("InMemCache Initialized")

	inMemCache := ttlcache.New[string, any](
		ttlcache.WithTTL[string, any](30 * time.Minute),
	)

	go inMemCache.Start()

	return inMemCache
}
