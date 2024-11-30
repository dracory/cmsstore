package shared

import (
	"time"

	"github.com/jellydator/ttlcache/v3"
)

var InMemCache *ttlcache.Cache[string, any]

func init() {
	InMemCache := ttlcache.New[string, any](
		ttlcache.WithTTL[string, any](30 * time.Minute),
	)

	go InMemCache.Start() // starts automatic expired item deletion
}
