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

func Test_Admin_Handle(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	require.NoError(t, err)

	options := AdminOptions{
		Store:  store,
		Logger: slog.New(slog.NewTextHandler(os.Stderr, nil)),
	}

	a, err := New(options)
	require.NoError(t, err)

	t.Run("home path", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/admin?path=home", nil)
		w := httptest.NewRecorder()
		a.Handle(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Dashboard")
	})

	t.Run("invalid path defaults to home", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/admin?path=nonexistent", nil)
		w := httptest.NewRecorder()
		a.Handle(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Dashboard")
	})

	t.Run("page manager path", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/admin?path=/pages/page-manager", nil)
		w := httptest.NewRecorder()
		a.Handle(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Pages")
	})
}
