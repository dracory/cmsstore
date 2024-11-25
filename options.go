package cmsstore

import (
	"context"
	"database/sql"
)

type Options struct {
	params map[string]any
}

type Option func(*Options)

func WithTransaction(tx *sql.Tx) Option {
	return func(o *Options) {
		if o.params == nil {
			o.params = make(map[string]any)
		}
		o.params["tx"] = tx
	}
}

func WithDryRun(dryRun bool) Option {
	return func(o *Options) {
		o.params["dryRun"] = dryRun
	}
}

func WithContext(ctx context.Context) Option {
	return func(o *Options) {
		o.params["ctx"] = ctx
	}
}
