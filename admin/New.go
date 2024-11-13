package admin

import (
	"errors"
	"log/slog"

	"github.com/gouniverse/blockeditor"
	"github.com/gouniverse/cmsstore"
	"github.com/gouniverse/cmsstore/admin/shared"
)

type AdminOptions struct {
	BlockEditorDefinitions []blockeditor.BlockDefinition
	FuncLayout             func(title string, body string, options struct {
		Styles     []string
		StyleURLs  []string
		Scripts    []string
		ScriptURLs []string
	}) string
	Logger       *slog.Logger
	Store        cmsstore.StoreInterface
	AdminHomeURL string
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
	}, nil
}
