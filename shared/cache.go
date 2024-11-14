package shared

import (
	"errors"

	"github.com/golang-module/carbon/v2"
)

type CacheInterface interface {
	Has(key string) bool
	Set(key string, value any, expireSeconds int) error
	Get(key string) (any, error)
	Delete(key string)
}

func Cache() CacheInterface {
	return &cache{
		parameters: make(map[string]any),
		expires:    make(map[string]int64),
	}
}

type cache struct {
	parameters map[string]any
	expires    map[string]int64
}

var _ CacheInterface = (*cache)(nil)

func (c *cache) Has(key string) bool {
	c.expire()
	return c.parameters[key] != nil
}

func (c *cache) Set(key string, value any, expireSeconds int) error {
	c.parameters[key] = value
	c.expires[key] = carbon.Now().AddSeconds(expireSeconds).Timestamp()
	return nil
}

func (c *cache) Get(key string) (any, error) {
	c.expire()

	if !c.Has(key) {
		return nil, errors.New("key not found")
	}

	return c.parameters[key], nil
}

func (c *cache) Delete(key string) {
	delete(c.parameters, key)
	delete(c.expires, key)
}

func (c *cache) expire() {
	now := carbon.Now().Timestamp()
	for k, v := range c.expires {
		if v < now {
			c.Delete(k)
		}
	}
}
