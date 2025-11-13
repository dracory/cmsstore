package cmsstore

import (
    "testing"

    "github.com/dracory/sb"
    "github.com/dromara/carbon/v2"
    "github.com/stretchr/testify/require"
)

func TestNewSiteDefaults(t *testing.T) {
    site := NewSite()

    require.NotEmpty(t, site.ID(), "ID should be generated")
    require.NotEmpty(t, site.CreatedAt(), "CreatedAt should be set")
    require.NotEmpty(t, site.UpdatedAt(), "UpdatedAt should be set")
    require.Equal(t, TEMPLATE_STATUS_DRAFT, site.Status(), "Status should default to draft")
    require.Equal(t, sb.MAX_DATETIME, site.SoftDeletedAt(), "SoftDeletedAt should default to max datetime")
    require.False(t, site.IsSoftDeleted(), "New site should not be soft deleted")

    domainNames, err := site.DomainNames()
    require.NoError(t, err)
    require.Empty(t, domainNames, "Domain names should default to empty slice")

    metas, err := site.Metas()
    require.NoError(t, err)
    require.Empty(t, metas, "Metas should default to empty map")

    createdCarbon := site.CreatedAtCarbon()
    require.NotNil(t, createdCarbon)
    require.Equal(t, site.CreatedAt(), createdCarbon.ToDateTimeString(carbon.UTC))

    updatedCarbon := site.UpdatedAtCarbon()
    require.NotNil(t, updatedCarbon)
    require.Equal(t, site.UpdatedAt(), updatedCarbon.ToDateTimeString(carbon.UTC))

    softDeletedCarbon := site.SoftDeletedAtCarbon()
    require.NotNil(t, softDeletedCarbon)
    require.True(t, softDeletedCarbon.Gte(carbon.Now(carbon.UTC)), "SoftDeletedAt should be in the future by default")
}

func TestSiteDomainNamesRoundTrip(t *testing.T) {
    site := NewSite()

    expectedDomains := []string{"example.com", "www.example.com"}
    _, err := site.SetDomainNames(expectedDomains)
    require.NoError(t, err)

    domainNames, err := site.DomainNames()
    require.NoError(t, err)
    require.Equal(t, expectedDomains, domainNames)
}

func TestSiteMetasUpsertAndLookup(t *testing.T) {
    site := NewSite()

    err := site.SetMetas(map[string]string{"theme": "default"})
    require.NoError(t, err)
    require.Equal(t, "default", site.Meta("theme"))

    err = site.UpsertMetas(map[string]string{"theme": "custom", "currency": "USD"})
    require.NoError(t, err)
    require.Equal(t, "custom", site.Meta("theme"))
    require.Equal(t, "USD", site.Meta("currency"))
    require.Equal(t, "", site.Meta("missing"))
}

func TestSiteSoftDeleteBehaviour(t *testing.T) {
    site := NewSite()
    require.False(t, site.IsSoftDeleted())

    past := carbon.Now(carbon.UTC).SubHour()
    site.SetSoftDeletedAt(past.ToDateTimeString(carbon.UTC))

    require.True(t, site.IsSoftDeleted(), "Site should be marked as soft deleted when past timestamp is set")
    require.Equal(t, past.ToDateTimeString(carbon.UTC), site.SoftDeletedAt())
}
