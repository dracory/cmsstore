package admin

import (
	"log/slog"
	"os"
	"testing"

	"github.com/dracory/cmsstore/admin/shared"
	"github.com/dracory/cmsstore/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

func Test_NewAdmin(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	require.NoError(t, err)

	t.Run("valid options", func(t *testing.T) {
		options := AdminOptions{
			Store:  store,
			Logger: slog.New(slog.NewTextHandler(os.Stderr, nil)),
		}
		a, err := New(options)
		require.NoError(t, err)
		assert.NotNil(t, a)
	})

	t.Run("missing store", func(t *testing.T) {
		options := AdminOptions{
			Logger: slog.New(slog.NewTextHandler(os.Stderr, nil)),
		}
		a, err := New(options)
		assert.Error(t, err)
		assert.Nil(t, a)
		assert.Contains(t, err.Error(), shared.ERROR_STORE_IS_NIL)
	})

	t.Run("missing logger", func(t *testing.T) {
		options := AdminOptions{
			Store: store,
		}
		a, err := New(options)
		assert.Error(t, err)
		assert.Nil(t, a)
		assert.Contains(t, err.Error(), shared.ERROR_LOGGER_IS_NIL)
	})
}
