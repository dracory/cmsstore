package cmsstore

import (
	"testing"

	"github.com/dracory/sb"
	"github.com/dromara/carbon/v2"
	"github.com/stretchr/testify/require"
)

func TestNewTemplateDefaults(t *testing.T) {
	template := NewTemplate()

	require.NotEmpty(t, template.ID())
	require.NotEmpty(t, template.CreatedAt())
	require.NotEmpty(t, template.UpdatedAt())
	require.Equal(t, TEMPLATE_STATUS_DRAFT, template.Status())
	require.Equal(t, sb.MAX_DATETIME, template.SoftDeletedAt())
	require.False(t, template.IsSoftDeleted())

	metas, err := template.Metas()
	require.NoError(t, err)
	require.Empty(t, metas)

	createdCarbon := template.CreatedAtCarbon()
	require.NotNil(t, createdCarbon)
	require.Equal(t, template.CreatedAt(), createdCarbon.ToDateTimeString(carbon.UTC))

	updatedCarbon := template.UpdatedAtCarbon()
	require.NotNil(t, updatedCarbon)
	require.Equal(t, template.UpdatedAt(), updatedCarbon.ToDateTimeString(carbon.UTC))

	softDeletedCarbon := template.SoftDeletedAtCarbon()
	require.NotNil(t, softDeletedCarbon)
	require.True(t, softDeletedCarbon.Gte(carbon.Now(carbon.UTC)))
}

func TestTemplateGetterMethods(t *testing.T) {
	template := NewTemplate()

	// Test default values
	require.Equal(t, "", template.Content())
	require.Equal(t, "", template.Editor())
	require.Equal(t, "", template.Handle())
	require.Equal(t, "", template.Memo())
	require.Equal(t, "", template.Name())
	require.Equal(t, "", template.SiteID())
}

func TestTemplateStatusMethods(t *testing.T) {
	template := NewTemplate()

	// Test default status (DRAFT)
	require.False(t, template.IsActive())
	require.False(t, template.IsInactive())

	// Test ACTIVE status
	template.SetStatus(TEMPLATE_STATUS_ACTIVE)
	require.True(t, template.IsActive())
	require.False(t, template.IsInactive())

	// Test INACTIVE status
	template.SetStatus(TEMPLATE_STATUS_INACTIVE)
	require.False(t, template.IsActive())
	require.True(t, template.IsInactive())

	// Test other status values
	template.SetStatus("unknown")
	require.False(t, template.IsActive())
	require.False(t, template.IsInactive())
}

func TestTemplateSoftDeleteMethods(t *testing.T) {
	template := NewTemplate()
	require.False(t, template.IsSoftDeleted())

	// Test with future date
	future := carbon.Now(carbon.UTC).AddHour()
	template.SetSoftDeletedAt(future.ToDateTimeString(carbon.UTC))
	require.False(t, template.IsSoftDeleted())
	require.Equal(t, future.ToDateTimeString(carbon.UTC), template.SoftDeletedAt())

	// Test with past date
	past := carbon.Now(carbon.UTC).SubHour()
	template.SetSoftDeletedAt(past.ToDateTimeString(carbon.UTC))
	require.True(t, template.IsSoftDeleted())
	require.Equal(t, past.ToDateTimeString(carbon.UTC), template.SoftDeletedAt())
}

func TestTemplateMetasMethods(t *testing.T) {
	template := NewTemplate()

	// Test empty metas
	metas, err := template.Metas()
	require.NoError(t, err)
	require.Empty(t, metas)

	// Test Meta lookup on empty metas
	require.Equal(t, "", template.Meta("nonexistent"))

	// Test SetMetas
	err = template.SetMetas(map[string]string{"layout": "main", "theme": "dark"})
	require.NoError(t, err)

	metas, err = template.Metas()
	require.NoError(t, err)
	require.Equal(t, "main", metas["layout"])
	require.Equal(t, "dark", metas["theme"])

	// Test Meta lookup
	require.Equal(t, "main", template.Meta("layout"))
	require.Equal(t, "dark", template.Meta("theme"))
	require.Equal(t, "", template.Meta("nonexistent"))

	// Test SetMeta
	err = template.SetMeta("newkey", "newvalue")
	require.NoError(t, err)
	require.Equal(t, "newvalue", template.Meta("newkey"))

	// Test UpsertMetas
	err = template.UpsertMetas(map[string]string{"layout": "sidebar", "color": "blue"})
	require.NoError(t, err)
	require.Equal(t, "sidebar", template.Meta("layout")) // Updated
	require.Equal(t, "dark", template.Meta("theme"))     // Preserved
	require.Equal(t, "newvalue", template.Meta("newkey")) // Preserved
	require.Equal(t, "blue", template.Meta("color"))      // Added
}

func TestTemplateCreatedAtMethods(t *testing.T) {
	template := NewTemplate()

	// Test default CreatedAt
	createdAt := template.CreatedAt()
	require.NotEmpty(t, createdAt)

	createdAtCarbon := template.CreatedAtCarbon()
	require.NotNil(t, createdAtCarbon)
	require.Equal(t, createdAt, createdAtCarbon.ToDateTimeString(carbon.UTC))

	// Test SetCreatedAt
	testDate := "2023-12-25 10:30:00"
	template.SetCreatedAt(testDate)
	require.Equal(t, testDate, template.CreatedAt())

	createdAtCarbon = template.CreatedAtCarbon()
	require.NotNil(t, createdAtCarbon)
	require.Equal(t, testDate, createdAtCarbon.ToDateTimeString(carbon.UTC))
}

func TestTemplateUpdatedAtMethods(t *testing.T) {
	template := NewTemplate()

	// Test default UpdatedAt
	updatedAt := template.UpdatedAt()
	require.NotEmpty(t, updatedAt)

	updatedAtCarbon := template.UpdatedAtCarbon()
	require.NotNil(t, updatedAtCarbon)
	require.Equal(t, updatedAt, updatedAtCarbon.ToDateTimeString(carbon.UTC))

	// Test SetUpdatedAt
	testDate := "2023-12-25 15:45:00"
	template.SetUpdatedAt(testDate)
	require.Equal(t, testDate, template.UpdatedAt())

	updatedAtCarbon = template.UpdatedAtCarbon()
	require.NotNil(t, updatedAtCarbon)
	require.Equal(t, testDate, updatedAtCarbon.ToDateTimeString(carbon.UTC))
}

func TestTemplateSoftDeletedAtMethods(t *testing.T) {
	template := NewTemplate()

	// Test default SoftDeletedAt
	softDeletedAt := template.SoftDeletedAt()
	require.Equal(t, sb.MAX_DATETIME, softDeletedAt)

	softDeletedAtCarbon := template.SoftDeletedAtCarbon()
	require.NotNil(t, softDeletedAtCarbon)
	require.Equal(t, softDeletedAt, softDeletedAtCarbon.ToDateTimeString(carbon.UTC))

	// Test SetSoftDeletedAt
	testDate := "2023-12-25 20:00:00"
	template.SetSoftDeletedAt(testDate)
	require.Equal(t, testDate, template.SoftDeletedAt())

	softDeletedAtCarbon = template.SoftDeletedAtCarbon()
	require.NotNil(t, softDeletedAtCarbon)
	require.Equal(t, testDate, softDeletedAtCarbon.ToDateTimeString(carbon.UTC))
}

func TestTemplateIDMethods(t *testing.T) {
	template := NewTemplate()

	// Test default ID
	id := template.ID()
	require.NotEmpty(t, id)

	// Test SetID
	newID := "test-template-id-123"
	template.SetID(newID)
	require.Equal(t, newID, template.ID())
}

func TestTemplateContentMethods(t *testing.T) {
	template := NewTemplate()

	// Test default content
	require.Equal(t, "", template.Content())

	// Test SetContent
	content := "This is template content"
	template.SetContent(content)
	require.Equal(t, content, template.Content())
}

func TestTemplateEditorMethods(t *testing.T) {
	template := NewTemplate()

	// Test default editor
	require.Equal(t, "", template.Editor())

	// Test SetEditor
	editor := "test-editor"
	template.SetEditor(editor)
	require.Equal(t, editor, template.Editor())
}

func TestTemplateHandleMethods(t *testing.T) {
	template := NewTemplate()

	// Test default handle
	require.Equal(t, "", template.Handle())

	// Test SetHandle
	handle := "test-template-handle"
	template.SetHandle(handle)
	require.Equal(t, handle, template.Handle())
}

func TestTemplateMemoMethods(t *testing.T) {
	template := NewTemplate()

	// Test default memo
	require.Equal(t, "", template.Memo())

	// Test SetMemo
	memo := "This is a template memo"
	template.SetMemo(memo)
	require.Equal(t, memo, template.Memo())
}

func TestTemplateNameMethods(t *testing.T) {
	template := NewTemplate()

	// Test default name
	require.Equal(t, "", template.Name())

	// Test SetName
	name := "Test Template Name"
	template.SetName(name)
	require.Equal(t, name, template.Name())
}

func TestTemplateSiteIDMethods(t *testing.T) {
	template := NewTemplate()

	// Test default site ID
	require.Equal(t, "", template.SiteID())

	// Test SetSiteID
	siteID := "test-site-id"
	template.SetSiteID(siteID)
	require.Equal(t, siteID, template.SiteID())
}

func TestTemplateStatusSettersAndGetters(t *testing.T) {
	template := NewTemplate()

	// Test default status
	require.Equal(t, TEMPLATE_STATUS_DRAFT, template.Status())

	// Test SetStatus
	template.SetStatus(TEMPLATE_STATUS_ACTIVE)
	require.Equal(t, TEMPLATE_STATUS_ACTIVE, template.Status())

	template.SetStatus(TEMPLATE_STATUS_INACTIVE)
	require.Equal(t, TEMPLATE_STATUS_INACTIVE, template.Status())

	template.SetStatus("custom-status")
	require.Equal(t, "custom-status", template.Status())
}