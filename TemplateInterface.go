package cmsstore

import "github.com/golang-module/carbon/v2"

type TemplateInterface interface {
	Data() map[string]string
	DataChanged() map[string]string
	MarkAsNotDirty()

	ID() string
	SetID(id string) TemplateInterface

	CreatedAt() string
	SetCreatedAt(createdAt string) TemplateInterface
	CreatedAtCarbon() carbon.Carbon

	Content() string
	SetContent(content string) TemplateInterface

	Editor() string
	SetEditor(editor string) TemplateInterface

	Handle() string
	SetHandle(handle string) TemplateInterface

	Memo() string
	SetMemo(memo string) TemplateInterface

	Meta(key string) string
	SetMeta(key, value string) error
	Metas() (map[string]string, error)
	SetMetas(metas map[string]string) error
	UpsertMetas(metas map[string]string) error

	Name() string
	SetName(name string) TemplateInterface

	SoftDeletedAt() string
	SetSoftDeletedAt(softDeletedAt string) TemplateInterface
	SoftDeletedAtCarbon() carbon.Carbon

	Status() string
	SetStatus(status string) TemplateInterface

	UpdatedAt() string
	SetUpdatedAt(updatedAt string) TemplateInterface
	UpdatedAtCarbon() carbon.Carbon

	IsActive() bool
	IsInactive() bool
	IsSoftDeleted() bool
}
