package shared

import (
	"log/slog"

	"github.com/gouniverse/cmsstore"
	"github.com/gouniverse/hb"
	"github.com/spf13/cast"
)

func AdminHeader(store cmsstore.StoreInterface, logger *slog.Logger, endpoint string) hb.TagInterface {
	linkHome := hb.NewHyperlink().
		HTML("Dashboard").
		Href(URL(endpoint, PathHome, nil)).
		Class("nav-link")
	linkBlocks := hb.Hyperlink().
		HTML("Blocks ").
		Href(URL(endpoint, PathBlocksBlockManager, nil)).
		Class("nav-link")
	// linkMenus := hb.NewHyperlink().
	// 	HTML("Menus ").
	// 	Href(endpoint + "?path=" + PathMenusMenuManager).
	// 	Class("nav-link")
	linkPages := hb.Hyperlink().
		HTML("Pages ").
		Href(URL(endpoint, PathPagesPageManager, nil)).
		Class("nav-link")
	linkTemplates := hb.Hyperlink().
		HTML("Templates ").
		Href(URL(endpoint, PathTemplatesTemplateManager, nil)).
		Class("nav-link")
	linkSites := hb.Hyperlink().
		HTML("Sites ").
		Href(URL(endpoint, PathSitesSiteManager, nil)).
		Class("nav-link")
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
		Href(URL(endpoint, PathTranslationsTranslationManager, nil)).
		Class("nav-link")

	templatesCount, err := store.TemplateCount(cmsstore.TemplateQuery())

	if err != nil {
		logger.Error(err.Error())
		templatesCount = -1
	}

	blocksCount, err := store.BlockCount(cmsstore.BlockQuery())

	if err != nil {
		logger.Error(err.Error())
		blocksCount = -1
	}

	pagesCount, err := store.PageCount(cmsstore.PageQuery())

	if err != nil {
		logger.Error(err.Error())
		pagesCount = -1
	}

	sitesCount, err := store.SiteCount(cmsstore.SiteQuery())

	if err != nil {
		logger.Error(err.Error())
		sitesCount = -1
	}

	translationsCount, err := store.TranslationCount(cmsstore.TranslationQuery())

	if err != nil {
		logger.Error(err.Error())
		translationsCount = -1
	}

	ulNav := hb.NewUL().Class("nav  nav-pills justify-content-center")
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

	// if cms.menusEnabled {
	// 	ulNav.AddChild(hb.NewLI().Class("nav-item").AddChild(linkMenus.AddChild(hb.NewSpan().Class("badge bg-secondary").HTML(strconv.FormatInt(menusCount, 10)))))
	// }

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

	ulNav.Child(hb.
		LI().
		Class("nav-item").
		Child(linkTranslations.
			Child(hb.NewSpan().
				Class("badge bg-secondary ms-1").
				HTML(cast.ToString(translationsCount)))))

	// if cms.settingsEnabled {
	// 	ulNav.AddChild(hb.NewLI().Class("nav-item").AddChild(linkSettings))
	// }
	// add Translations

	// for _, entity := range cms.customEntityList {
	// 	linkEntity := hb.NewHyperlink().HTML(entity.TypeLabel).Href(endpoint + "?path=entities/entity-manager&type=" + entity.Type).Class("nav-link")
	// 	ulNav.AddChild(hb.NewLI().Class("nav-item").Child(linkEntity))
	// }

	divCard := hb.NewDiv().Class("card card-default mt-3 mb-3")
	divCardBody := hb.NewDiv().Class("card-body").Style("padding: 2px;")
	return divCard.AddChild(divCardBody.AddChild(ulNav))
}
