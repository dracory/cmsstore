package cmsstore

import (
    "testing"

    "github.com/dracory/sb"
    "github.com/dromara/carbon/v2"
    "github.com/stretchr/testify/require"
)

func TestNewMenuDefaults(t *testing.T) {
    menu := NewMenu()

    require.NotEmpty(t, menu.ID())
    require.NotEmpty(t, menu.CreatedAt())
    require.NotEmpty(t, menu.UpdatedAt())
    require.Equal(t, TEMPLATE_STATUS_DRAFT, menu.Status())
    require.Equal(t, sb.MAX_DATETIME, menu.SoftDeletedAt())
    require.False(t, menu.IsSoftDeleted())

    metas, err := menu.Metas()
    require.NoError(t, err)
    require.Empty(t, metas)

    createdCarbon := menu.CreatedAtCarbon()
    require.NotNil(t, createdCarbon)
    require.Equal(t, menu.CreatedAt(), createdCarbon.ToDateTimeString(carbon.UTC))

    updatedCarbon := menu.UpdatedAtCarbon()
    require.NotNil(t, updatedCarbon)
    require.Equal(t, menu.UpdatedAt(), updatedCarbon.ToDateTimeString(carbon.UTC))

    softDeletedCarbon := menu.SoftDeletedAtCarbon()
    require.NotNil(t, softDeletedCarbon)
    require.True(t, softDeletedCarbon.Gte(carbon.Now(carbon.UTC)))
}

func TestMenuMetasUpsertAndLookup(t *testing.T) {
    menu := NewMenu()

    err := menu.SetMetas(map[string]string{"layout": "main"})
    require.NoError(t, err)
    require.Equal(t, "main", menu.Meta("layout"))

    err = menu.UpsertMetas(map[string]string{"layout": "sidebar", "theme": "dark"})
    require.NoError(t, err)
    require.Equal(t, "sidebar", menu.Meta("layout"))
    require.Equal(t, "dark", menu.Meta("theme"))
    require.Equal(t, "", menu.Meta("missing"))
}

func TestMenuSoftDeleteBehaviour(t *testing.T) {
    menu := NewMenu()
    require.False(t, menu.IsSoftDeleted())

    past := carbon.Now(carbon.UTC).SubHour()
    menu.SetSoftDeletedAt(past.ToDateTimeString(carbon.UTC))

    require.True(t, menu.IsSoftDeleted())
    require.Equal(t, past.ToDateTimeString(carbon.UTC), menu.SoftDeletedAt())
}
