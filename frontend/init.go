package frontend

import "github.com/gouniverse/cmsstore/shared"

type LanguageKey struct{}

var inMemCache shared.CacheInterface

func init() {
	inMemCache = shared.Cache()
}
