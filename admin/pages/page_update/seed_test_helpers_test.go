package page_update

import (
	"github.com/dracory/cmsstore"
	"github.com/dracory/cmsstore/testutils"
)

func seedTestPage(store cmsstore.StoreInterface) (cmsstore.PageInterface, error) {
	return testutils.SeedPage(store, testutils.SITE_01, testutils.PAGE_01)
}
