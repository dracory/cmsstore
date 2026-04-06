package frontend

import (
	"fmt"
	"time"

	"github.com/jellydator/ttlcache/v3"
)

type LanguageKey struct{}

func init() {}

func initCache() *ttlcache.Cache[string, any] {
	fmt.Println("InMemCache Initialized")

	inMemCache := ttlcache.New[string, any](
		ttlcache.WithTTL[string, any](30 * time.Minute),
	)

	go inMemCache.Start()

	return inMemCache
}
