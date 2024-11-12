package frontend

type LanguageKey struct{}

var inMemCache CacheInterface

func init() {
	inMemCache = Cache()
}
