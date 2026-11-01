package shared

import (
	"log/slog"
	"net/http"

	"github.com/dracory/cmsstore"
	"github.com/dracory/hb"
	"github.com/dracory/req"
	"github.com/spf13/cast"
)

func AdminHeader(store cmsstore.StoreInterface, logger *slog.Logger, r *http.Request) hb.TagInterface {
	active := req.GetStringTrimmed(r, "path")
	if active == "" || active == "/" {
		active = PathHome
	}

	linkHome := hb.NewHyperlink().
		HTML("Dashboard").
		Href(URLR(r, PathHome, nil)).
		Class("nav-link").
		ClassIf(active == PathHome, "active")
	linkBlocks := hb.Hyperlink().
		HTML("Blocks ").
		Href(URLR(r, PathBlocksBlockManager, nil)).
		Class("nav-link").
		ClassIf(active == PathBlocksBlockManager, "active")
	linkMenus := hb.NewHyperlink().
		HTML("Menus ").
		Href(URLR(r, PathMenusMenuManager, nil)).
		Class("nav-link").
		ClassIf(active == PathMenusMenuManager, "active")
	linkPages := hb.Hyperlink().
		HTML("Pages ").
		Href(URLR(r, PathPagesPageManager, nil)).
		Class("nav-link").
		ClassIf(active == PathPagesPageManager, "active")
	linkTemplates := hb.Hyperlink().
		HTML("Templates ").
		Href(URLR(r, PathTemplatesTemplateManager, nil)).
		Class("nav-link").
		ClassIf(active == PathTemplatesTemplateManager, "active")
	linkSites := hb.Hyperlink().
		HTML("Sites ").
		Href(URLR(r, PathSitesSiteManager, nil)).
		Class("nav-link").
		ClassIf(active == PathSitesSiteManager, "active")
	// linkWidgets := hb.NewHyperlink().
	// 	HTML("Widgets ").
	// 	Href(endpoint + "?path=" + PathWidgetsWidgetManager).
	// 	Class("nav-link")
	// linkSettings := hb.NewHyperlink().
	// 	HTML("Settings").
	// 	Href(endpoint + "?path=" + PathSettingsSettingManager).
	// 	Class("nav-link")
	linkTranslations := hb.Hyperlink().
		HTML("Translations").
		Href(URLR(r, PathTranslationsTranslationManager, nil)).
		Class("nav-link").
		ClassIf(active == PathTranslationsTranslationManager, "active")
	mediaManagerURL := MediaManagerURL(r)

	templatesCount, err := store.TemplateCount(r.Context(), cmsstore.TemplateQuery())

	if err != nil {
		logger.Error(err.Error())
		templatesCount = -1
	}

	blocksCount, err := store.BlockCount(r.Context(), cmsstore.BlockQuery())

	if err != nil {
		logger.Error(err.Error())
		blocksCount = -1
	}

	menusCount, err := store.MenuCount(r.Context(), cmsstore.MenuQuery())

	if err != nil {
		logger.Error(err.Error())
		menusCount = -1
	}

	pagesCount, err := store.PageCount(r.Context(), cmsstore.PageQuery())

	if err != nil {
		logger.Error(err.Error())
		pagesCount = -1
	}

	sitesCount, err := store.SiteCount(r.Context(), cmsstore.SiteQuery())

	if err != nil {
		logger.Error(err.Error())
		sitesCount = -1
	}

	translationsCount, err := store.TranslationCount(r.Context(), cmsstore.TranslationQuery())

	if err != nil {
		logger.Error(err.Error())
		translationsCount = -1
	}

	ulNav := hb.NewUL().Class("nav nav-tabs mb-3")
	ulNav.AddChild(hb.NewLI().Class("nav-item").Child(linkHome))

	ulNav.Child(hb.LI().
		Class("nav-item").
		Child(linkSites.
			Child(hb.Span().
				Class("badge bg-secondary").
				HTML(cast.ToString(sitesCount)))))

	ulNav.Child(hb.LI().
		Class("nav-item").
		Child(linkTemplates.
			Child(hb.Span().
				Class("badge bg-secondary").
				HTML(cast.ToString(templatesCount)))))

	ulNav.Child(hb.
		LI().
		Class("nav-item").
		Child(linkPages.
			Child(hb.NewSpan().
				Class("badge bg-secondary").
				HTML(cast.ToString(pagesCount)))))

	if store.MenusEnabled() {
		ulNav.Child(hb.
			LI().
			Class("nav-item").
			Child(linkMenus.
				Child(hb.NewSpan().
					Class("badge bg-secondary").
					HTML(cast.ToString(menusCount)))))
	}

	ulNav.Child(hb.
		LI().
		Class("nav-item").
		Child(linkBlocks.
			Child(hb.NewSpan().
				Class("badge bg-secondary").
				HTML(cast.ToString(blocksCount)))))

	// if cms.widgetsEnabled {
	// 	ulNav.AddChild(hb.NewLI().Class("nav-item").AddChild(linkWidgets.AddChild(hb.NewSpan().Class("badge bg-secondary").HTML(strconv.FormatInt(widgetsCount, 10)))))
	// }

	if store.TranslationsEnabled() {
		ulNav.Child(hb.
			LI().
			Class("nav-item").
			Child(linkTranslations.
				Child(hb.NewSpan().
					Class("badge bg-secondary ms-1").
					HTML(cast.ToString(translationsCount)))))
	}

	if mediaManagerURL != "" {
		linkMedia := hb.Hyperlink().
			HTML("Media").
			Href(mediaManagerURL).
			Target("_blank").
			Class("nav-link")

		ulNav.Child(hb.LI().
			Class("nav-item").
			Child(linkMedia))
	}

	// if cms.settingsEnabled {
	// 	ulNav.AddChild(hb.NewLI().Class("nav-item").AddChild(linkSettings))
	// }
	// add Translations

	// for _, entity := range cms.customEntityList {
	// 	linkEntity := hb.NewHyperlink().HTML(entity.TypeLabel).Href(endpoint + "?path=entities/entity-manager&type=" + entity.Type).Class("nav-link")
	// 	ulNav.AddChild(hb.NewLI().Class("nav-item").Child(linkEntity))
	// }

	return ulNav
}
