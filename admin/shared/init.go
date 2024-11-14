package shared

import "github.com/gouniverse/cmsstore/shared"

var InMemCache shared.CacheInterface

func init() {
	InMemCache = shared.Cache()
}
