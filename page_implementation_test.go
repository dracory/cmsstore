package cmsstore

import (
    "testing"

    "github.com/dracory/sb"
    "github.com/dromara/carbon/v2"
    "github.com/stretchr/testify/require"
)

func TestNewPageDefaults(t *testing.T) {
    page := NewPage()

    require.NotEmpty(t, page.ID(), "ID should be generated")
    require.NotEmpty(t, page.CreatedAt(), "CreatedAt should be set")
    require.NotEmpty(t, page.UpdatedAt(), "UpdatedAt should be set")
    require.Equal(t, PAGE_STATUS_DRAFT, page.Status(), "Status should default to draft")
    require.Equal(t, sb.MAX_DATETIME, page.SoftDeletedAt(), "SoftDeletedAt should default to max datetime")
    require.False(t, page.IsSoftDeleted(), "New page should not be soft deleted")

    require.Equal(t, "", page.Content())
    require.Equal(t, "", page.Editor())
    require.Equal(t, "", page.Handle())
    require.Equal(t, "", page.TemplateID())

    metas, err := page.Metas()
    require.NoError(t, err)
    require.Empty(t, metas)

    require.Empty(t, page.MiddlewaresBefore())
    require.Empty(t, page.MiddlewaresAfter())

    createdCarbon := page.CreatedAtCarbon()
    require.NotNil(t, createdCarbon)
    require.Equal(t, page.CreatedAt(), createdCarbon.ToDateTimeString(carbon.UTC))

    updatedCarbon := page.UpdatedAtCarbon()
    require.NotNil(t, updatedCarbon)
    require.Equal(t, page.UpdatedAt(), updatedCarbon.ToDateTimeString(carbon.UTC))

    softDeletedCarbon := page.SoftDeletedAtCarbon()
    require.NotNil(t, softDeletedCarbon)
    require.True(t, softDeletedCarbon.Gte(carbon.Now(carbon.UTC)))
}

func TestPageMetasUpsertAndLookup(t *testing.T) {
    page := NewPage()

    err := page.SetMetas(map[string]string{"layout": "default"})
    require.NoError(t, err)
    require.Equal(t, "default", page.Meta("layout"))

    err = page.UpsertMetas(map[string]string{"layout": "custom", "author": "cms"})
    require.NoError(t, err)
    require.Equal(t, "custom", page.Meta("layout"))
    require.Equal(t, "cms", page.Meta("author"))
    require.Equal(t, "", page.Meta("missing"))
}

func TestPageMiddlewaresRoundTrip(t *testing.T) {
    page := NewPage()

    before := []string{"auth", "cache"}
    after := []string{"compress"}

    page.SetMiddlewaresBefore(before)
    page.SetMiddlewaresAfter(after)

    require.Equal(t, before, page.MiddlewaresBefore())
    require.Equal(t, after, page.MiddlewaresAfter())
}

func TestPageSoftDeleteBehaviour(t *testing.T) {
    page := NewPage()
    require.False(t, page.IsSoftDeleted())

    past := carbon.Now(carbon.UTC).SubHour()
    page.SetSoftDeletedAt(past.ToDateTimeString(carbon.UTC))

    require.True(t, page.IsSoftDeleted())
    require.Equal(t, past.ToDateTimeString(carbon.UTC), page.SoftDeletedAt())
}
