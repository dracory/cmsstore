package frontend

import (
	"net/http"
)

// Mock middleware implementation
type MockMiddleware struct {
	identifier     string
	name           string
	description    string
	middlewareType string
	handler        func(http.Handler) http.Handler
}

func (m *MockMiddleware) Identifier() string                       { return m.identifier }
func (m *MockMiddleware) Name() string                             { return m.name }
func (m *MockMiddleware) Description() string                      { return m.description }
func (m *MockMiddleware) Type() string                             { return m.middlewareType }
func (m *MockMiddleware) Handler() func(http.Handler) http.Handler { return m.handler }

// Mock store to return middlewares
// type MockStore struct {
// 	shortcodes  []cmsstore.ShortcodeInterface
// 	middlewares []cmsstore.MiddlewareInterface
// }

// var _ cmsstore.StoreInterface = (*MockStore)(nil)

// func (ms *MockStore) BlockCount(ctx context.Context, query cmsstore.BlockQueryInterface) (int64, error) {
// 	return 0, nil
// }

// func (ms *MockStore) BlockCreate(ctx context.Context, block cmsstore.BlockInterface) error {
// 	return nil
// }

// func (ms *MockStore) BlockDelete(ctx context.Context, block cmsstore.BlockInterface) error {
// 	return nil
// }

// func (ms *MockStore) BlockDeleteByID(ctx context.Context, id string) error {
// 	return nil
// }

// func (ms *MockStore) BlockFindByHandle(ctx context.Context, handle string) (cmsstore.BlockInterface, error) {
// 	return nil, nil
// }

// func (ms *MockStore) BlockFindByID(ctx context.Context, blockID string) (cmsstore.BlockInterface, error) {
// 	return nil, nil
// }

// func (ms *MockStore) BlockList(ctx context.Context, query cmsstore.BlockQueryInterface) ([]cmsstore.BlockInterface, error) {
// 	return nil, nil
// }

// func (ms *MockStore) BlockSoftDelete(ctx context.Context, block cmsstore.BlockInterface) error {
// 	return nil
// }

// func (ms *MockStore) BlockSoftDeleteByID(ctx context.Context, id string) error {
// 	return nil
// }

// func (ms *MockStore) BlockUpdate(ctx context.Context, block cmsstore.BlockInterface) error {
// 	return nil
// }

// func (ms *MockStore) AutoMigrate(ctx context.Context, opts ...cmsstore.Option) error { return nil }

// func (ms *MockStore) AddShortcode(shortcode cmsstore.ShortcodeInterface) {
// 	ms.shortcodes = append(ms.shortcodes, shortcode)
// }
// func (ms *MockStore) AddShortcodes(shortcodes []cmsstore.ShortcodeInterface) {
// 	ms.shortcodes = append(ms.shortcodes, shortcodes...)
// }
// func (ms *MockStore) Shortcodes() []cmsstore.ShortcodeInterface { return ms.shortcodes }

// func (ms *MockStore) AddMiddleware(middleware cmsstore.MiddlewareInterface) {
// 	ms.middlewares = append(ms.middlewares, middleware)
// }

// func (ms *MockStore) AddMiddlewares(middlewares []cmsstore.MiddlewareInterface) {
// 	ms.middlewares = append(ms.middlewares, middlewares...)
// }

// func (ms *MockStore) Middlewares() []cmsstore.MiddlewareInterface { return ms.middlewares }
