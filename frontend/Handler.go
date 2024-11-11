package frontend

import (
	"log/slog"
	"net/http"
	"regexp"
	"strings"

	"github.com/gouniverse/cms/types"
	"github.com/gouniverse/cmsstore"
	"github.com/gouniverse/hb"
	"github.com/gouniverse/shortcode"
	"github.com/gouniverse/ui"
	"github.com/gouniverse/utils"
	"github.com/samber/lo"
)

type LanguageKey struct{}

type Config struct {
	BlockEditorRenderer func(blocks []ui.BlockInterface) string
	Logger              *slog.Logger
	Shortcodes          []cmsstore.ShortcodeInterface
	Store               cmsstore.StoreInterface
}

func New(config Config) frontend {
	return frontend{
		blockEditorRenderer: config.BlockEditorRenderer,
		logger:              config.Logger,
		shortcodes:          config.Shortcodes,
		store:               config.Store,
	}
}

type frontend struct {
	blockEditorRenderer func(blocks []ui.BlockInterface) string
	logger              *slog.Logger
	shortcodes          []cmsstore.ShortcodeInterface
	store               cmsstore.StoreInterface
}

// Handler is the main handler for the CMS frontend.
//
// It handles the routing of the request to the appropriate page.
//
// If the URI ends with ".ico", it will return a blank response, as the browsers
// (at least Chrome and Firefox) will always request the favicon even if
// it's not present in the HTML.
func (frontend *frontend) Handler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(frontend.StringHandler(w, r)))
}

// FrontendHandlerRenderAsString is the same as FrontendHandler but returns a string
// instead of writing to the http.ResponseWriter.
//
// It handles the routing of the request to the appropriate page.
//
// If the URI ends with ".ico", it will return a blank response, as the browsers
// (at least Chrome and Firefox) will always request the favicon even if
// it's not present in the HTML.
//
// If the translations are enabled, it will use the language from the request context.
// If the language is not valid, it will use the default language for the translations.
func (frontend *frontend) StringHandler(w http.ResponseWriter, r *http.Request) string {
	domain := r.Host

	uri := r.RequestURI

	if strings.HasSuffix(uri, ".ico") {
		return ""
	}

	languageAny := r.Context().Value(LanguageKey{})
	language := utils.ToString(languageAny)

	// if fr.translationsEnabled {
	// 	isValidLanguage := lo.Contains(lo.Keys(cms.translationLanguages), language)

	// 	if !isValidLanguage {
	// 		language = cms.translationLanguageDefault
	// 	}
	// }

	site, err := frontend.store.SiteFindByDomainName(domain)

	if err != nil {
		frontend.logger.Error(`At StringHandler`, "error", err.Error())
		return `Domain not supported: ` + domain
	}

	if site == nil {
		return `Domain not supported: ` + domain
	}

	return frontend.PageRenderHtmlByAlias(r, site.ID(), r.URL.Path, language)
}

// PageRenderHtmlByAlias builds the HTML of a page based on its alias
func (frontend *frontend) PageRenderHtmlByAlias(r *http.Request, siteID string, alias string, language string) string {
	page, err := frontend.PageFindBySiteAndAlias(siteID, alias)

	if err != nil {
		frontend.logger.Error(`At PageRenderHtmlByAlias`, "error", err.Error())
		return hb.NewDiv().
			Text(`Page with alias '`).Text(alias).Text(`' not found`).
			ToHTML()
	}

	if page == nil {
		return hb.NewDiv().
			Text(`Page with alias '`).Text(alias).Text(`' not found`).
			ToHTML()
	}

	pageContent := page.Content()
	pageTitle := page.Title()
	pageMetaKeywords := page.MetaKeywords()
	pageMetaDescription := page.MetaDescription()
	pageMetaRobots := page.MetaRobots()
	pageCanonicalURL := page.CanonicalUrl()
	pageEditor := page.Editor()
	pageTemplateID := page.TemplateID()

	if pageEditor == types.WEBPAGE_EDITOR_BLOCKEDITOR {
		if frontend.blockEditorRenderer == nil {
			return "Block editor not configured"
		}

		if !utils.IsJSON(pageContent) {
			return "Malformed block content"
		}

		blocks, err := ui.BlocksFromJson(pageContent)

		if err != nil {
			return "Error parsing block content"
		}

		pageContent = frontend.blockEditorRenderer(blocks)
	}

	if pageTemplateID == "" {
		return pageContent
	}

	finalContent := lo.If(pageTemplateID == "", pageContent).ElseF(func() string {
		template, err := frontend.store.TemplateFindByID(pageTemplateID)
		if err != nil {
			frontend.logger.Error(`At PageRenderHtmlByAlias`, "error", err.Error())
			return pageContent
		}

		if template == nil {
			return pageContent
		}

		return template.Content()
	})

	html, err := frontend.renderContentToHtml(r, finalContent, struct {
		PageContent         string
		PageCanonicalURL    string
		PageMetaDescription string
		PageMetaKeywords    string
		PageMetaRobots      string
		PageTitle           string
		Language            string
	}{
		PageContent:         pageContent,
		PageCanonicalURL:    pageCanonicalURL,
		PageMetaDescription: pageMetaDescription,
		PageMetaKeywords:    pageMetaKeywords,
		PageMetaRobots:      pageMetaRobots,
		PageTitle:           pageTitle,
	})

	if err != nil {
		frontend.logger.Error(`At PageRenderHtmlByAlias`, "error", err.Error())
		return hb.NewDiv().Text(`error occurred`).ToHTML()
	}

	return html
}

// renderContentToHtml renders the content to HTML
//
// This is done in the following steps (sequence is important):
// 1. replaces placeholders with values
// 2. renders the blocks
// 3. renders the shortcodes
// 3. renders the translations
// 4. returns the HTML
//
// Parameters:
// - r: the HTTP request
// - content: the content to render
// - options: the options for the rendering
//
// Returns:
// - html: the rendered HTML
// - err: the error, if any, or nil otherwise
func (frontend *frontend) renderContentToHtml(r *http.Request, content string, options struct {
	PageContent         string
	PageCanonicalURL    string
	PageMetaDescription string
	PageMetaKeywords    string
	PageMetaRobots      string
	PageTitle           string
	Language            string
}) (html string, err error) {
	replacements := map[string]string{
		"PageContent":         options.PageContent,
		"PageCanonicalUrl":    options.PageCanonicalURL,
		"PageMetaDescription": options.PageMetaDescription,
		"PageMetaKeywords":    options.PageMetaKeywords,
		"PageRobots":          options.PageMetaRobots,
		"PageTitle":           options.PageTitle,
	}

	for key, value := range replacements {
		content = strings.ReplaceAll(content, "[["+key+"]]", value)
		content = strings.ReplaceAll(content, "[[ "+key+" ]]", value)
	}

	content, err = frontend.ContentRenderBlocks(content)

	if err != nil {
		return "", err
	}

	content, err = frontend.ContentRenderShortcodes(r, content)

	if err != nil {
		return "", err
	}

	language := lo.If(options.Language == "", "en").Else(options.Language)

	content, err = frontend.ContentRenderTranslations(content, language)

	if err != nil {
		return "", err
	}

	return content, nil
}

// PageFindByAlias helper method to find a page by alias
//
// =====================================================================
//  1. It will attempt to find the page by the provided alias exactly
//     as provided
//  2. It will attempt to find the page with the alias prefixed with "/"
//     in case of error
//
// =====================================================================
func (frontend *frontend) PageFindBySiteAndAlias(siteID string, alias string) (cmsstore.PageInterface, error) {
	// Try to find by "alias"
	pages, err := frontend.store.PageList(cmsstore.PageQuery().
		SetSiteID(siteID).
		SetAlias(alias).
		SetLimit(1))

	if err != nil {
		return nil, err
	}

	if len(pages) > 0 {
		return pages[0], nil
	}

	// Try to find by "/alias"
	pages, err = frontend.store.PageList(cmsstore.PageQuery().
		SetSiteID(siteID).
		SetAlias("/" + alias).
		SetLimit(1))

	if err != nil {
		return nil, err
	}

	if len(pages) > 0 {
		return pages[0], nil
	}

	page, err := frontend.PageFindBySiteAndAliasWithPatterns(siteID, alias)

	if err != nil {
		return nil, err
	}

	if page != nil {
		return page, nil
	}

	return nil, nil
}

// PageFindByAliasWithPatterns helper method to find a page by matching patterns
//
// =====================================================================
//
//	The following patterns are supported:
//	:any
//	:num
//	:all
//	:string
//	:number
//	:numeric
//	:alpha
//
// =====================================================================
func (frontend *frontend) PageFindBySiteAndAliasWithPatterns(siteID string, alias string) (cmsstore.PageInterface, error) {
	patterns := map[string]string{
		":any":     "([^/]+)",
		":num":     "([0-9]+)",
		":all":     "(.*)",
		":string":  "([a-zA-Z]+)",
		":number":  "([0-9]+)",
		":numeric": "([0-9-.]+)",
		":alpha":   "([a-zA-Z0-9-_]+)",
	}

	// can we optimize this to retrieve only the id and alias column?
	pages, err := frontend.store.PageList(cmsstore.PageQuery().
		SetSiteID(siteID))

	if err != nil {
		return nil, err
	}

	pageAliasMap := make(map[string]string, len(pages))

	for _, page := range pages {
		pageAliasMap[page.ID()] = page.Alias()
	}

	for pageID, pageAlias := range pageAliasMap {
		if !strings.Contains(pageAlias, ":") {
			continue
		}

		for pattern, replacement := range patterns {
			pageAlias = strings.ReplaceAll(pageAlias, pattern, replacement)
		}

		matcher := regexp.MustCompile("^" + pageAlias + "$")
		if matcher.MatchString(alias) {
			return frontend.store.PageFindByID(pageID)
		}
	}

	return nil, nil
}

// RenderBlocks renders the blocks in a string
func (frontend *frontend) ContentRenderBlocks(content string) (string, error) {
	blockIDs := ContentFindIdsByPatternPrefix(content, "BLOCK")

	var err error
	for _, blockID := range blockIDs {
		content, err = frontend.ContentRenderBlockByID(content, blockID)

		if err != nil {
			return content, err
		}
	}

	return content, nil
}

// ContentRenderTranslations renders the translations in a string
func (frontend *frontend) ContentRenderTranslations(content string, language string) (string, error) {
	translationIDs := ContentFindIdsByPatternPrefix(content, "TRANSLATION")

	var err error
	for _, translationID := range translationIDs {
		content, err = frontend.ContentRenderTranslationByIdOrHandle(content, translationID, language)

		if err != nil {
			return content, err
		}
	}

	return content, nil
}

// returns the IDs in the content who have the following format [[prefix_id]]
func ContentFindIdsByPatternPrefix(content, prefix string) []string {
	ids := []string{}

	re := regexp.MustCompilePOSIX("|\\[\\[" + prefix + "_(.*)\\]\\]|U")

	matches := re.FindAllStringSubmatch(content, -1)

	for _, match := range matches {
		if match[0] == "" {
			continue
		}
		ids = append(ids, match[1])
	}

	return ids
}

// ContentRenderBlockByID renders the block specified by the ID in a content
// if the blockID is empty or not found the initial content is returned
func (frontend *frontend) ContentRenderBlockByID(content string, blockID string) (string, error) {
	if blockID == "" {
		return content, nil
	}

	blockContent, err := frontend.findBlockContent(blockID)

	if err != nil {
		return content, err
	}

	if blockContent == "" {
		return content, nil
	}

	content = strings.ReplaceAll(content, "[[BLOCK_"+blockID+"]]", blockContent)
	content = strings.ReplaceAll(content, "[[ BLOCK_"+blockID+" ]]", blockContent)

	return content, nil
}

// ContentRenderShortcodes renders the shortcodes in a string
func (frontend *frontend) ContentRenderShortcodes(req *http.Request, content string) (string, error) {
	sh, err := shortcode.NewShortcode(shortcode.WithBrackets("<", ">"))

	if err != nil {
		return "", err
	}

	for _, shortcode := range frontend.shortcodes {
		content = sh.RenderWithRequest(req, content, shortcode.Alias(), shortcode.Render)
	}

	return content, nil
}

// ContentRenderTranslationByIdOrHandle renders the translation specified by the ID in a content
// if the blockID is empty or not found the initial content is returned
func (frontend *frontend) ContentRenderTranslationByIdOrHandle(content string, translationID string, language string) (string, error) {
	return content, nil

	// Will be implemented once translations are transferred

	// if translationID == "" {
	// 	return content, nil
	// }

	// translation, err := frontend.store.TranslationFindByIdOrHandle(translationID, language)

	// if err != nil {
	// 	return "", err
	// }

	// content = strings.ReplaceAll(content, "[[TRANSLATION_"+translationID+"]]", translation)
	// content = strings.ReplaceAll(content, "[[ TRANSLATION_"+translationID+" ]]", translation)

	// return content, nil
}

func (frontend *frontend) findBlockContent(blockID string) (string, error) {
	block, err := frontend.store.BlockFindByID(blockID)

	if err != nil {
		return "", err
	}

	if block == nil {
		return "", nil
	}

	if block.IsActive() {
		return block.Content(), nil
	}

	return "", nil
}
