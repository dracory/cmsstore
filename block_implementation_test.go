package cmsstore

import (
	"testing"

	"github.com/dracory/sb"
	"github.com/dromara/carbon/v2"
	"github.com/stretchr/testify/require"
)

func TestNewBlockDefaults(t *testing.T) {
	block := NewBlock()

	require.NotEmpty(t, block.ID())
	require.NotEmpty(t, block.CreatedAt())
	require.NotEmpty(t, block.UpdatedAt())
	require.Equal(t, BLOCK_STATUS_DRAFT, block.Status())
	require.Equal(t, sb.MAX_DATETIME, block.SoftDeletedAt())
	require.False(t, block.IsSoftDeleted())

	metas, err := block.Metas()
	require.NoError(t, err)
	require.Empty(t, metas)

	createdCarbon := block.CreatedAtCarbon()
	require.NotNil(t, createdCarbon)
	require.Equal(t, block.CreatedAt(), createdCarbon.ToDateTimeString(carbon.UTC))

	updatedCarbon := block.UpdatedAtCarbon()
	require.NotNil(t, updatedCarbon)
	require.Equal(t, block.UpdatedAt(), updatedCarbon.ToDateTimeString(carbon.UTC))

	softDeletedCarbon := block.SoftDeletedAtCarbon()
	require.NotNil(t, softDeletedCarbon)
	require.True(t, softDeletedCarbon.Gte(carbon.Now(carbon.UTC)))
}

func TestBlockGetterMethods(t *testing.T) {
	block := NewBlock()

	// Test default values
	require.Equal(t, "", block.Content())
	require.Equal(t, "", block.Editor())
	require.Equal(t, "", block.Handle())
	require.Equal(t, "", block.Memo())
	require.Equal(t, "", block.Name())
	require.Equal(t, "", block.PageID())
	require.Equal(t, "", block.ParentID())
	require.Equal(t, "0", block.Sequence())
	require.Equal(t, 0, block.SequenceInt())
	require.Equal(t, "", block.SiteID())
	require.Equal(t, "", block.TemplateID())
	require.Equal(t, BLOCK_TYPE_HTML, block.Type())
}

func TestBlockStatusMethods(t *testing.T) {
	block := NewBlock()

	// Test default status (DRAFT)
	require.False(t, block.IsActive())
	require.False(t, block.IsInactive())

	// Test ACTIVE status
	block.SetStatus(BLOCK_STATUS_ACTIVE)
	require.True(t, block.IsActive())
	require.False(t, block.IsInactive())

	// Test INACTIVE status
	block.SetStatus(BLOCK_STATUS_INACTIVE)
	require.False(t, block.IsActive())
	require.True(t, block.IsInactive())

	// Test other status values
	block.SetStatus("unknown")
	require.False(t, block.IsActive())
	require.False(t, block.IsInactive())
}

func TestBlockSoftDeleteMethods(t *testing.T) {
	block := NewBlock()
	require.False(t, block.IsSoftDeleted())

	// Test with future date
	future := carbon.Now(carbon.UTC).AddHour()
	block.SetSoftDeletedAt(future.ToDateTimeString(carbon.UTC))
	require.False(t, block.IsSoftDeleted())
	require.Equal(t, future.ToDateTimeString(carbon.UTC), block.SoftDeletedAt())

	// Test with past date
	past := carbon.Now(carbon.UTC).SubHour()
	block.SetSoftDeletedAt(past.ToDateTimeString(carbon.UTC))
	require.True(t, block.IsSoftDeleted())
	require.Equal(t, past.ToDateTimeString(carbon.UTC), block.SoftDeletedAt())
}

func TestBlockMetasMethods(t *testing.T) {
	block := NewBlock()

	// Test empty metas
	metas, err := block.Metas()
	require.NoError(t, err)
	require.Empty(t, metas)

	// Test Meta lookup on empty metas
	require.Equal(t, "", block.Meta("nonexistent"))

	// Test SetMetas
	err = block.SetMetas(map[string]string{"layout": "main", "theme": "dark"})
	require.NoError(t, err)

	metas, err = block.Metas()
	require.NoError(t, err)
	require.Equal(t, "main", metas["layout"])
	require.Equal(t, "dark", metas["theme"])

	// Test Meta lookup
	require.Equal(t, "main", block.Meta("layout"))
	require.Equal(t, "dark", block.Meta("theme"))
	require.Equal(t, "", block.Meta("nonexistent"))

	// Test SetMeta
	err = block.SetMeta("newkey", "newvalue")
	require.NoError(t, err)
	require.Equal(t, "newvalue", block.Meta("newkey"))

	// Test UpsertMetas
	err = block.UpsertMetas(map[string]string{"layout": "sidebar", "color": "blue"})
	require.NoError(t, err)
	require.Equal(t, "sidebar", block.Meta("layout"))  // Updated
	require.Equal(t, "dark", block.Meta("theme"))      // Preserved
	require.Equal(t, "newvalue", block.Meta("newkey")) // Preserved
	require.Equal(t, "blue", block.Meta("color"))      // Added
}

func TestBlockSequenceMethods(t *testing.T) {
	block := NewBlock()

	// Test default sequence
	require.Equal(t, "0", block.Sequence())
	require.Equal(t, 0, block.SequenceInt())

	// Test SetSequenceInt
	block.SetSequenceInt(42)
	require.Equal(t, "42", block.Sequence())
	require.Equal(t, 42, block.SequenceInt())

	// Test SetSequence
	block.SetSequence("123")
	require.Equal(t, "123", block.Sequence())
	require.Equal(t, 123, block.SequenceInt())

	// Test invalid sequence
	block.SetSequence("invalid")
	require.Equal(t, "invalid", block.Sequence())
	require.Equal(t, 0, block.SequenceInt()) // Should default to 0 for invalid
}

func TestBlockCreatedAtMethods(t *testing.T) {
	block := NewBlock()

	// Test default CreatedAt
	createdAt := block.CreatedAt()
	require.NotEmpty(t, createdAt)

	createdAtCarbon := block.CreatedAtCarbon()
	require.NotNil(t, createdAtCarbon)
	require.Equal(t, createdAt, createdAtCarbon.ToDateTimeString(carbon.UTC))

	// Test SetCreatedAt
	testDate := "2023-12-25 10:30:00"
	block.SetCreatedAt(testDate)
	require.Equal(t, testDate, block.CreatedAt())

	createdAtCarbon = block.CreatedAtCarbon()
	require.NotNil(t, createdAtCarbon)
	require.Equal(t, testDate, createdAtCarbon.ToDateTimeString(carbon.UTC))
}

func TestBlockUpdatedAtMethods(t *testing.T) {
	block := NewBlock()

	// Test default UpdatedAt
	updatedAt := block.UpdatedAt()
	require.NotEmpty(t, updatedAt)

	updatedAtCarbon := block.UpdatedAtCarbon()
	require.NotNil(t, updatedAtCarbon)
	require.Equal(t, updatedAt, updatedAtCarbon.ToDateTimeString(carbon.UTC))

	// Test SetUpdatedAt
	testDate := "2023-12-25 15:45:00"
	block.SetUpdatedAt(testDate)
	require.Equal(t, testDate, block.UpdatedAt())

	updatedAtCarbon = block.UpdatedAtCarbon()
	require.NotNil(t, updatedAtCarbon)
	require.Equal(t, testDate, updatedAtCarbon.ToDateTimeString(carbon.UTC))
}

func TestBlockSoftDeletedAtMethods(t *testing.T) {
	block := NewBlock()

	// Test default SoftDeletedAt
	softDeletedAt := block.SoftDeletedAt()
	require.Equal(t, sb.MAX_DATETIME, softDeletedAt)

	softDeletedAtCarbon := block.SoftDeletedAtCarbon()
	require.NotNil(t, softDeletedAtCarbon)
	require.Equal(t, softDeletedAt, softDeletedAtCarbon.ToDateTimeString(carbon.UTC))

	// Test SetSoftDeletedAt
	testDate := "2023-12-25 20:00:00"
	block.SetSoftDeletedAt(testDate)
	require.Equal(t, testDate, block.SoftDeletedAt())

	softDeletedAtCarbon = block.SoftDeletedAtCarbon()
	require.NotNil(t, softDeletedAtCarbon)
	require.Equal(t, testDate, softDeletedAtCarbon.ToDateTimeString(carbon.UTC))
}

func TestBlockIDMethods(t *testing.T) {
	block := NewBlock()

	// Test default ID
	id := block.ID()
	require.NotEmpty(t, id)

	// Test SetID
	newID := "test-block-id-123"
	block.SetID(newID)
	require.Equal(t, newID, block.ID())
}

func TestBlockContentMethods(t *testing.T) {
	block := NewBlock()

	// Test default content
	require.Equal(t, "", block.Content())

	// Test SetContent
	content := "This is block content"
	block.SetContent(content)
	require.Equal(t, content, block.Content())
}

func TestBlockEditorMethods(t *testing.T) {
	block := NewBlock()

	// Test default editor
	require.Equal(t, "", block.Editor())

	// Test SetEditor
	editor := "test-editor"
	block.SetEditor(editor)
	require.Equal(t, editor, block.Editor())
}

func TestBlockHandleMethods(t *testing.T) {
	block := NewBlock()

	// Test default handle
	require.Equal(t, "", block.Handle())

	// Test SetHandle
	handle := "test-block-handle"
	block.SetHandle(handle)
	require.Equal(t, handle, block.Handle())
}

func TestBlockMemoMethods(t *testing.T) {
	block := NewBlock()

	// Test default memo
	require.Equal(t, "", block.Memo())

	// Test SetMemo
	memo := "This is a block memo"
	block.SetMemo(memo)
	require.Equal(t, memo, block.Memo())
}

func TestBlockNameMethods(t *testing.T) {
	block := NewBlock()

	// Test default name
	require.Equal(t, "", block.Name())

	// Test SetName
	name := "Test Block Name"
	block.SetName(name)
	require.Equal(t, name, block.Name())
}

func TestBlockPageIDMethods(t *testing.T) {
	block := NewBlock()

	// Test default page ID
	require.Equal(t, "", block.PageID())

	// Test SetPageID
	pageID := "test-page-id"
	block.SetPageID(pageID)
	require.Equal(t, pageID, block.PageID())
}

func TestBlockParentIDMethods(t *testing.T) {
	block := NewBlock()

	// Test default parent ID
	require.Equal(t, "", block.ParentID())

	// Test SetParentID
	parentID := "test-parent-id"
	block.SetParentID(parentID)
	require.Equal(t, parentID, block.ParentID())
}

func TestBlockSiteIDMethods(t *testing.T) {
	block := NewBlock()

	// Test default site ID
	require.Equal(t, "", block.SiteID())

	// Test SetSiteID
	siteID := "test-site-id"
	block.SetSiteID(siteID)
	require.Equal(t, siteID, block.SiteID())
}

func TestBlockTemplateIDMethods(t *testing.T) {
	block := NewBlock()

	// Test default template ID
	require.Equal(t, "", block.TemplateID())

	// Test SetTemplateID
	templateID := "test-template-id"
	block.SetTemplateID(templateID)
	require.Equal(t, templateID, block.TemplateID())
}

func TestBlockTypeMethods(t *testing.T) {
	block := NewBlock()

	// Test default type (HTML)
	require.Equal(t, BLOCK_TYPE_HTML, block.Type())

	// Test SetType
	blockType := "text"
	block.SetType(blockType)
	require.Equal(t, blockType, block.Type())
}

func TestBlockStatusSettersAndGetters(t *testing.T) {
	block := NewBlock()

	// Test default status
	require.Equal(t, BLOCK_STATUS_DRAFT, block.Status())

	// Test SetStatus
	block.SetStatus(BLOCK_STATUS_ACTIVE)
	require.Equal(t, BLOCK_STATUS_ACTIVE, block.Status())

	block.SetStatus(BLOCK_STATUS_INACTIVE)
	require.Equal(t, BLOCK_STATUS_INACTIVE, block.Status())

	block.SetStatus("custom-status")
	require.Equal(t, "custom-status", block.Status())
}
