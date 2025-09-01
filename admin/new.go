package admin

import (
	"errors"
	"log/slog"

	"github.com/dracory/blockeditor"
	"github.com/dracory/cmsstore"
	"github.com/dracory/cmsstore/admin/shared"
	"github.com/samber/lo"
)

type AdminOptions struct {
	// BlockEditorDefinitions is the block definitions to use for the block editor
	// these definitions must also be registered in the frontend
	BlockEditorDefinitions []blockeditor.BlockDefinition

	// FuncLayout is an optional function to use to render the admin interface inside
	// this is convinient when you want to use your own layout to wrap the admin
	// interface, i.e. completely replace the default layout with your own
	// admin panel with your own branding, logos, etc.
	FuncLayout func(title string, body string, options struct {
		Styles     []string
		StyleURLs  []string
		Scripts    []string
		ScriptURLs []string
	}) string

	// Logger is the logger to use to log any errors. Optional
	Logger *slog.Logger

	// Store is the cmsstore.StoreInterface to use by the admin panel
	Store cmsstore.StoreInterface

	AdminHomeURL string

	// flags holds a map of feature flags for internal use
	Flags map[string]bool
}

func New(options AdminOptions) (*admin, error) {
	if options.Store == nil {
		return nil, errors.New(shared.ERROR_STORE_IS_NIL)
	}

	if options.Logger == nil {
		return nil, errors.New(shared.ERROR_LOGGER_IS_NIL)
	}

	return &admin{
		blockEditorDefinitions: options.BlockEditorDefinitions,
		logger:                 options.Logger,
		store:                  options.Store,
		funcLayout:             options.FuncLayout,
		adminHomeURL:           options.AdminHomeURL,
		flags:                  lo.Ternary(options.Flags != nil, options.Flags, map[string]bool{}),
	}, nil
}
