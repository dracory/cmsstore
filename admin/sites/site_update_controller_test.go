package admin

import (
    "context"
    "io"
    "log/slog"
    "net/http"
    "net/url"
    "strings"
    "testing"

    "github.com/dracory/cmsstore"
    "github.com/dracory/cmsstore/admin/shared"
    "github.com/dracory/cmsstore/testutils"
    "github.com/dracory/test"
)

func initSiteUpdateHandler(store cmsstore.StoreInterface) func(w http.ResponseWriter, r *http.Request) string {
    ui := UI(shared.UiConfig{
        Layout: shared.Layout,
        Logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
        Store:  store,
    })

    return NewSiteUpdateController(ui).Handler
}

func Test_SiteUpdateController_RepeaterAdd_AppendsEmptyDomain(t *testing.T) {
    store, err := testutils.InitStore(":memory:")
    if err != nil {
        t.Fatalf("InitStore should succeed, got error: %v", err)
    }

    handler := initSiteUpdateHandler(store)

    site, err := testutils.SeedSite(store, testutils.SITE_01)
    if err != nil {
        t.Fatalf("Seeding site should succeed, got error: %v", err)
    }

    _, err = site.SetDomainNames([]string{"example.com"})
    if err != nil {
        t.Fatalf("Setting domain names should succeed, got error: %v", err)
    }

    if err := store.SiteUpdate(context.Background(), site); err != nil {
        t.Fatalf("Persisting site should succeed, got error: %v", err)
    }

    getValues := url.Values{
        "site_id": {site.ID()},
        "view":    {VIEW_SETTINGS},
    }

    postValues := url.Values{
        "action": {ACTION_REPEATER_ADD},
        "site_domain_names[0][site_domain_name]": {"example.com"},
    }

    body, response, err := test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
        GetValues:  getValues,
        PostValues: postValues,
    })

    if err != nil {
        t.Fatalf("CallStringEndpoint should succeed, got error: %v", err)
    }

    if response.StatusCode != http.StatusOK {
        t.Fatalf("Expected status %d, got %d", http.StatusOK, response.StatusCode)
    }

    if !strings.Contains(body, "name=\"site_domain_names[site_domain_name][]\" type=\"text\" value=\"example.com\"") {
        t.Fatalf("Expected existing domain to persist in response body, got: %s", body)
    }

    if strings.Count(body, "name=\"site_domain_names[site_domain_name][]\"") < 2 {
        t.Fatalf("Expected new empty domain input to be appended, got: %s", body)
    }

    if !strings.Contains(body, "name=\"site_domain_names[site_domain_name][]\" type=\"text\" value=\"\"") {
        t.Fatalf("Expected appended input to be blank, got: %s", body)
    }
}
