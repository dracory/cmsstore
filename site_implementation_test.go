package cmsstore

import (
	"encoding/json"
	"testing"

	"github.com/dracory/sb"
	"github.com/dromara/carbon/v2"
	"github.com/stretchr/testify/require"
)

func TestNewSiteDefaults(t *testing.T) {
	site := NewSite()

	require.NotEmpty(t, site.ID())
	require.NotEmpty(t, site.CreatedAt())
	require.NotEmpty(t, site.UpdatedAt())
	require.Equal(t, TEMPLATE_STATUS_DRAFT, site.Status())
	require.Equal(t, sb.MAX_DATETIME, site.SoftDeletedAt())
	require.False(t, site.IsSoftDeleted())

	metas, err := site.Metas()
	require.NoError(t, err)
	require.Empty(t, metas)

	createdCarbon := site.CreatedAtCarbon()
	require.NotNil(t, createdCarbon)
	require.Equal(t, site.CreatedAt(), createdCarbon.ToDateTimeString(carbon.UTC))

	updatedCarbon := site.UpdatedAtCarbon()
	require.NotNil(t, updatedCarbon)
	require.Equal(t, site.UpdatedAt(), updatedCarbon.ToDateTimeString(carbon.UTC))

	softDeletedCarbon := site.SoftDeletedAtCarbon()
	require.NotNil(t, softDeletedCarbon)
	require.True(t, softDeletedCarbon.Gte(carbon.Now(carbon.UTC)))
}

func TestSiteGetterMethods(t *testing.T) {
	site := NewSite()

	// Test default values
	require.Equal(t, "", site.Handle())
	require.Equal(t, "", site.Memo())
	require.Equal(t, "", site.Name())
}

func TestSiteStatusMethods(t *testing.T) {
	site := NewSite()

	// Test default status (DRAFT)
	require.False(t, site.IsActive())
	require.False(t, site.IsInactive())

	// Test ACTIVE status
	site.SetStatus(PAGE_STATUS_ACTIVE)
	require.True(t, site.IsActive())
	require.False(t, site.IsInactive())

	// Test INACTIVE status
	site.SetStatus(PAGE_STATUS_INACTIVE)
	require.False(t, site.IsActive())
	require.True(t, site.IsInactive())

	// Test other status values
	site.SetStatus("unknown")
	require.False(t, site.IsActive())
	require.False(t, site.IsInactive())
}

func TestSiteSoftDeleteMethods(t *testing.T) {
	site := NewSite()
	require.False(t, site.IsSoftDeleted())

	// Test with future date
	future := carbon.Now(carbon.UTC).AddHour()
	site.SetSoftDeletedAt(future.ToDateTimeString(carbon.UTC))
	require.False(t, site.IsSoftDeleted())
	require.Equal(t, future.ToDateTimeString(carbon.UTC), site.SoftDeletedAt())

	// Test with past date
	past := carbon.Now(carbon.UTC).SubHour()
	site.SetSoftDeletedAt(past.ToDateTimeString(carbon.UTC))
	require.True(t, site.IsSoftDeleted())
	require.Equal(t, past.ToDateTimeString(carbon.UTC), site.SoftDeletedAt())
}

func TestSiteMetasMethods(t *testing.T) {
	site := NewSite()

	// Test empty metas
	metas, err := site.Metas()
	require.NoError(t, err)
	require.Empty(t, metas)

	// Test Meta lookup on empty metas
	require.Equal(t, "", site.Meta("nonexistent"))

	// Test SetMetas
	err = site.SetMetas(map[string]string{"layout": "main", "theme": "dark"})
	require.NoError(t, err)

	metas, err = site.Metas()
	require.NoError(t, err)
	require.Equal(t, "main", metas["layout"])
	require.Equal(t, "dark", metas["theme"])

	// Test Meta lookup
	require.Equal(t, "main", site.Meta("layout"))
	require.Equal(t, "dark", site.Meta("theme"))
	require.Equal(t, "", site.Meta("nonexistent"))

	// Test SetMeta
	err = site.SetMeta("newkey", "newvalue")
	require.NoError(t, err)
	require.Equal(t, "newvalue", site.Meta("newkey"))

	// Test UpsertMetas
	err = site.UpsertMetas(map[string]string{"layout": "sidebar", "color": "blue"})
	require.NoError(t, err)
	require.Equal(t, "sidebar", site.Meta("layout")) // Updated
	require.Equal(t, "dark", site.Meta("theme"))     // Preserved
	require.Equal(t, "newvalue", site.Meta("newkey")) // Preserved
	require.Equal(t, "blue", site.Meta("color"))      // Added
}

func TestSiteDomainNamesMethods(t *testing.T) {
	site := NewSite()

	// Test default domain names
	domainNames, err := site.DomainNames()
	require.NoError(t, err)
	require.Empty(t, domainNames)

	// Test SetDomainNames
	domains := []string{"example.com", "www.example.com"}
	site, err = site.SetDomainNames(domains)
	require.NoError(t, err)

	domainNames, err = site.DomainNames()
	require.NoError(t, err)
	require.Equal(t, domains, domainNames)

	// Test empty domain names
	site, err = site.SetDomainNames([]string{})
	require.NoError(t, err)

	domainNames, err = site.DomainNames()
	require.NoError(t, err)
	require.Empty(t, domainNames)
}

func TestSiteCreatedAtMethods(t *testing.T) {
	site := NewSite()

	// Test default CreatedAt
	createdAt := site.CreatedAt()
	require.NotEmpty(t, createdAt)

	createdAtCarbon := site.CreatedAtCarbon()
	require.NotNil(t, createdAtCarbon)
	require.Equal(t, createdAt, createdAtCarbon.ToDateTimeString(carbon.UTC))

	// Test SetCreatedAt
	testDate := "2023-12-25 10:30:00"
	site.SetCreatedAt(testDate)
	require.Equal(t, testDate, site.CreatedAt())

	createdAtCarbon = site.CreatedAtCarbon()
	require.NotNil(t, createdAtCarbon)
	require.Equal(t, testDate, createdAtCarbon.ToDateTimeString(carbon.UTC))
}

func TestSiteUpdatedAtMethods(t *testing.T) {
	site := NewSite()

	// Test default UpdatedAt
	updatedAt := site.UpdatedAt()
	require.NotEmpty(t, updatedAt)

	updatedAtCarbon := site.UpdatedAtCarbon()
	require.NotNil(t, updatedAtCarbon)
	require.Equal(t, updatedAt, updatedAtCarbon.ToDateTimeString(carbon.UTC))

	// Test SetUpdatedAt
	testDate := "2023-12-25 15:45:00"
	site.SetUpdatedAt(testDate)
	require.Equal(t, testDate, site.UpdatedAt())

	updatedAtCarbon = site.UpdatedAtCarbon()
	require.NotNil(t, updatedAtCarbon)
	require.Equal(t, testDate, updatedAtCarbon.ToDateTimeString(carbon.UTC))
}

func TestSiteSoftDeletedAtMethods(t *testing.T) {
	site := NewSite()

	// Test default SoftDeletedAt
	softDeletedAt := site.SoftDeletedAt()
	require.Equal(t, sb.MAX_DATETIME, softDeletedAt)

	softDeletedAtCarbon := site.SoftDeletedAtCarbon()
	require.NotNil(t, softDeletedAtCarbon)
	require.Equal(t, softDeletedAt, softDeletedAtCarbon.ToDateTimeString(carbon.UTC))

	// Test SetSoftDeletedAt
	testDate := "2023-12-25 20:00:00"
	site.SetSoftDeletedAt(testDate)
	require.Equal(t, testDate, site.SoftDeletedAt())

	softDeletedAtCarbon = site.SoftDeletedAtCarbon()
	require.NotNil(t, softDeletedAtCarbon)
	require.Equal(t, testDate, softDeletedAtCarbon.ToDateTimeString(carbon.UTC))
}

func TestSiteIDMethods(t *testing.T) {
	site := NewSite()

	// Test default ID
	id := site.ID()
	require.NotEmpty(t, id)

	// Test SetID
	newID := "test-site-id-123"
	site.SetID(newID)
	require.Equal(t, newID, site.ID())
}

func TestSiteHandleMethods(t *testing.T) {
	site := NewSite()

	// Test default handle
	require.Equal(t, "", site.Handle())

	// Test SetHandle
	handle := "test-site-handle"
	site.SetHandle(handle)
	require.Equal(t, handle, site.Handle())
}

func TestSiteMemoMethods(t *testing.T) {
	site := NewSite()

	// Test default memo
	require.Equal(t, "", site.Memo())

	// Test SetMemo
	memo := "This is a site memo"
	site.SetMemo(memo)
	require.Equal(t, memo, site.Memo())
}

func TestSiteNameMethods(t *testing.T) {
	site := NewSite()

	// Test default name
	require.Equal(t, "", site.Name())

	// Test SetName
	name := "Test Site Name"
	site.SetName(name)
	require.Equal(t, name, site.Name())
}

func TestSiteStatusSettersAndGetters(t *testing.T) {
	site := NewSite()

	// Test default status
	require.Equal(t, TEMPLATE_STATUS_DRAFT, site.Status())

	// Test SetStatus
	site.SetStatus(PAGE_STATUS_ACTIVE)
	require.Equal(t, PAGE_STATUS_ACTIVE, site.Status())

	site.SetStatus(PAGE_STATUS_INACTIVE)
	require.Equal(t, PAGE_STATUS_INACTIVE, site.Status())

	site.SetStatus("custom-status")
	require.Equal(t, "custom-status", site.Status())
}

func TestSiteMarshalToVersioning(t *testing.T) {
	site := NewSite()
	site.SetHandle("test-handle")
	site.SetMemo("test-memo")
	site.SetName("Test Site")
	site.SetStatus(PAGE_STATUS_ACTIVE)

	versionedJSON, err := site.MarshalToVersioning()
	require.NoError(t, err)
	require.NotEmpty(t, versionedJSON)

	// Parse the JSON to verify it contains expected fields
	var versionedData map[string]string
	err = json.Unmarshal([]byte(versionedJSON), &versionedData)
	require.NoError(t, err)

	// Check that expected fields are present
	require.Equal(t, "test-handle", versionedData[COLUMN_HANDLE])
	require.Equal(t, "test-memo", versionedData[COLUMN_MEMO])
	require.Equal(t, "Test Site", versionedData[COLUMN_NAME])
	require.Equal(t, PAGE_STATUS_ACTIVE, versionedData[COLUMN_STATUS])

	// Check that timestamps and soft delete fields are excluded
	_, hasCreatedAt := versionedData[COLUMN_CREATED_AT]
	_, hasUpdatedAt := versionedData[COLUMN_UPDATED_AT]
	_, hasSoftDeletedAt := versionedData[COLUMN_SOFT_DELETED_AT]
	require.False(t, hasCreatedAt)
	require.False(t, hasUpdatedAt)
	require.False(t, hasSoftDeletedAt)
}