package cmsstore

import (
    "testing"

    "github.com/dracory/sb"
    "github.com/dromara/carbon/v2"
    "github.com/stretchr/testify/require"
)

func TestNewMenuItemDefaults(t *testing.T) {
    menuItem := NewMenuItem()

    require.NotEmpty(t, menuItem.ID())
    require.NotEmpty(t, menuItem.CreatedAt())
    require.NotEmpty(t, menuItem.UpdatedAt())
    require.Equal(t, MENU_ITEM_STATUS_DRAFT, menuItem.Status())
    require.Equal(t, sb.MAX_DATETIME, menuItem.SoftDeletedAt())
    require.False(t, menuItem.IsSoftDeleted())

    metas, err := menuItem.Metas()
    require.NoError(t, err)
    require.Empty(t, metas)

    createdCarbon := menuItem.CreatedAtCarbon()
    require.NotNil(t, createdCarbon)
    require.Equal(t, menuItem.CreatedAt(), createdCarbon.ToDateTimeString(carbon.UTC))

    updatedCarbon := menuItem.UpdatedAtCarbon()
    require.NotNil(t, updatedCarbon)
    require.Equal(t, menuItem.UpdatedAt(), updatedCarbon.ToDateTimeString(carbon.UTC))

    softDeletedCarbon := menuItem.SoftDeletedAtCarbon()
    require.NotNil(t, softDeletedCarbon)
    require.True(t, softDeletedCarbon.Gte(carbon.Now(carbon.UTC)))
}

func TestMenuItemMetasUpsertAndLookup(t *testing.T) {
    menuItem := NewMenuItem()

    err := menuItem.SetMetas(map[string]string{"icon": "home"})
    require.NoError(t, err)
    require.Equal(t, "home", menuItem.Meta("icon"))

    err = menuItem.UpsertMetas(map[string]string{"icon": "dashboard", "section": "main"})
    require.NoError(t, err)
    require.Equal(t, "dashboard", menuItem.Meta("icon"))
    require.Equal(t, "main", menuItem.Meta("section"))
    require.Equal(t, "", menuItem.Meta("missing"))
}

func TestMenuItemSequenceConversions(t *testing.T) {
    menuItem := NewMenuItem()

    menuItem.SetSequence("10")
    require.Equal(t, "10", menuItem.Sequence())
    require.Equal(t, 10, menuItem.SequenceInt())

    menuItem.SetSequenceInt(25)
    require.Equal(t, "25", menuItem.Sequence())
    require.Equal(t, 25, menuItem.SequenceInt())
}

func TestMenuItemSoftDeleteBehaviour(t *testing.T) {
    menuItem := NewMenuItem()
    require.False(t, menuItem.IsSoftDeleted())

    past := carbon.Now(carbon.UTC).SubHour()
    menuItem.SetSoftDeletedAt(past.ToDateTimeString(carbon.UTC))

    require.True(t, menuItem.IsSoftDeleted())
    require.Equal(t, past.ToDateTimeString(carbon.UTC), menuItem.SoftDeletedAt())
}
