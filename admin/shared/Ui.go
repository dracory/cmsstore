package shared

// type Ui struct {
// 	adminHeader      hb.TagInterface
// 	adminBreadcrumbs func(endpoint string, breadcrumbs []Breadcrumb) hb.TagInterface
// 	endpoint         string
// 	layout           func(w http.ResponseWriter, r *http.Request, webpageTitle, webpageHtml string, options struct {
// 		Styles     []string
// 		StyleURLs  []string
// 		Scripts    []string
// 		ScriptURLs []string
// 	}) string
// 	logger *slog.Logger
// 	store  cmsstore.StoreInterface
// 	url    func(endpoint string, path string, params map[string]string) string
// }

// var _ UiInterface = Ui{}

// func (ui Ui) AdminHeader() hb.TagInterface {
// 	return ui.adminHeader
// }

// func (ui Ui) AdminBreadcrumbs(endpoint string, breadcrumbs []Breadcrumb) hb.TagInterface {
// 	return ui.adminBreadcrumbs(endpoint, breadcrumbs)
// }

// func (ui Ui) Endpoint() string {
// 	return ui.endpoint
// }

// func (ui Ui) Layout(w http.ResponseWriter, r *http.Request, webpageTitle, webpageHtml string, options struct {
// 	Styles     []string
// 	StyleURLs  []string
// 	Scripts    []string
// 	ScriptURLs []string
// }) string {
// 	return ui.layout(w, r, webpageTitle, webpageHtml, options)
// }

// func (ui Ui) Logger() *slog.Logger {
// 	return ui.logger
// }

// func (ui Ui) Store() cmsstore.StoreInterface {
// 	return ui.store
// }
