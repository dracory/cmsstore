package shared

import (
	"context"

	"github.com/dracory/cmsstore"
	"github.com/dracory/hb"
	"github.com/samber/lo"
)

func FilterDescriptionSite(ctx context.Context, store cmsstore.StoreInterface, siteID string) hb.TagInterface {
	if siteID == "" {
		return nil
	}

	siteList, err := CachedSiteList(ctx, store)
	if err != nil {
		siteList = []cmsstore.SiteInterface{}
	}

	site, isFound := lo.Find(siteList, func(site cmsstore.SiteInterface) bool {
		return site.ID() == siteID
	})

	siteName := lo.IfF(isFound, func() string { return site.Name() }).Else(siteID)

	return hb.Wrap().Child(hb.Span().
		Text("with site: ").
		Text(" ").
		Text(siteName))
}
