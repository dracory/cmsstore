// Package cmsstore provides the core functionality for the CMS store.
package cmsstore

import (
	"context"
	"database/sql"
)

// Options holds various configuration options for the CMS store operations.
type Options struct {
	params map[string]any
}

// Option is a function type that modifies an Options instance.
type Option func(*Options)

// WithTransaction sets the transaction for the CMS store operations.
func WithTransaction(tx *sql.Tx) Option {
	return func(o *Options) {
		if o.params == nil {
			o.params = make(map[string]any)
		}
		o.params["tx"] = tx
	}
}

// WithDryRun sets the dry run flag for the CMS store operations.
func WithDryRun(dryRun bool) Option {
	return func(o *Options) {
		o.params["dryRun"] = dryRun
	}
}

// WithContext sets the context for the CMS store operations.
func WithContext(ctx context.Context) Option {
	return func(o *Options) {
		o.params["ctx"] = ctx
	}
}
