package cmsstore

import "github.com/golang-module/carbon/v2"

type BlockInterface interface {
	Data() map[string]string
	DataChanged() map[string]string
	MarkAsNotDirty()

	ID() string
	SetID(id string) BlockInterface

	CreatedAt() string
	SetCreatedAt(createdAt string) BlockInterface
	CreatedAtCarbon() carbon.Carbon

	Content() string
	SetContent(content string) BlockInterface

	Editor() string
	SetEditor(editor string) BlockInterface

	Handle() string
	SetHandle(handle string) BlockInterface

	Memo() string
	SetMemo(memo string) BlockInterface

	Meta(key string) string
	SetMeta(key, value string) error
	Metas() (map[string]string, error)
	SetMetas(metas map[string]string) error
	UpsertMetas(metas map[string]string) error

	Name() string
	SetName(name string) BlockInterface

	PageID() string
	SetPageID(pageID string) BlockInterface

	ParentID() string
	SetParentID(parentID string) BlockInterface

	Sequence() string
	SequenceInt() int
	SetSequenceInt(sequence int) BlockInterface
	SetSequence(sequence string) BlockInterface

	SiteID() string
	SetSiteID(siteID string) BlockInterface

	TemplateID() string
	SetTemplateID(templateID string) BlockInterface

	SoftDeletedAt() string
	SetSoftDeletedAt(softDeletedAt string) BlockInterface
	SoftDeletedAtCarbon() carbon.Carbon

	// Status returns the status of the block, i.e. BLOCK_STATUS_ACTIVE
	Status() string

	// SetStatus sets the status of the block, i.e. BLOCK_STATUS_ACTIVE
	SetStatus(status string) BlockInterface

	// Type returns the type of the block, i.e. "text"
	Type() string

	// SetType sets the type of the block, i.e. "text"
	SetType(blockType string) BlockInterface

	// UpdatedAt returns the last updated time of block
	UpdatedAt() string

	// SetUpdatedAt sets the last updated time of block
	SetUpdatedAt(updatedAt string) BlockInterface

	// UpdatedAtCarbon returns carbon.Carbon of the last updated time of block
	UpdatedAtCarbon() carbon.Carbon

	IsActive() bool
	IsInactive() bool
	IsSoftDeleted() bool
}
