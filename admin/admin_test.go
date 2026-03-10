package admin

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/dracory/cmsstore/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

func newAdminForTest(t *testing.T, options AdminOptions) *admin {
	t.Helper()

	store, err := testutils.InitStore(":memory:")
	require.NoError(t, err)

	if options.Store == nil {
		options.Store = store
	}
	if options.Logger == nil {
		options.Logger = slog.New(slog.NewTextHandler(os.Stderr, nil))
	}

	a, err := New(options)
	require.NoError(t, err)
	return a
}

func TestAdminHandle_HomePath(t *testing.T) {
	a := newAdminForTest(t, AdminOptions{})

	req := httptest.NewRequest(http.MethodGet, "/admin?path=home", nil)
	w := httptest.NewRecorder()
	a.Handle(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Dashboard")
}

func TestAdminHandle_InvalidPathDefaultsToHome(t *testing.T) {
	a := newAdminForTest(t, AdminOptions{})

	req := httptest.NewRequest(http.MethodGet, "/admin?path=nonexistent", nil)
	w := httptest.NewRecorder()
	a.Handle(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Dashboard")
}

func TestAdminHandle_PageManagerPath(t *testing.T) {
	a := newAdminForTest(t, AdminOptions{})

	req := httptest.NewRequest(http.MethodGet, "/admin?path=/pages/page-manager", nil)
	w := httptest.NewRecorder()
	a.Handle(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Pages")
}

func TestAdminHandle_MediaLinkRenderedWhenConfigured(t *testing.T) {
	a := newAdminForTest(t, AdminOptions{
		MediaManagerURL: "/admin/media",
	})

	req := httptest.NewRequest(http.MethodGet, "/admin?path=home", nil)
	w := httptest.NewRecorder()
	a.Handle(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "href=\"/admin/media\"")
	assert.Contains(t, w.Body.String(), "target=\"_blank\"")
	assert.Contains(t, w.Body.String(), ">Media</a>")
}

func TestAdminHandle_MediaLinkNotRenderedWhenNotConfigured(t *testing.T) {
	a := newAdminForTest(t, AdminOptions{})

	req := httptest.NewRequest(http.MethodGet, "/admin?path=home", nil)
	w := httptest.NewRecorder()
	a.Handle(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotContains(t, w.Body.String(), ">Media</a>")
}

func TestAdminHandle_CMSContainerPaddingRendered(t *testing.T) {
	a := newAdminForTest(t, AdminOptions{
		PaddingTopPx:    9,
		PaddingRightPx:  8,
		PaddingBottomPx: 7,
		PaddingLeftPx:   6,
	})

	req := httptest.NewRequest(http.MethodGet, "/admin?path=home", nil)
	w := httptest.NewRecorder()
	a.Handle(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "class=\"cms\"")
	assert.Contains(t, w.Body.String(), "padding: 9px 8px 7px 6px;")
}
