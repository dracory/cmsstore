package cmsstore

import (
	"time"

	"github.com/dracory/neat/database/orm"
	"github.com/dracory/neat/database/soft_delete"
	neatuid "github.com/dracory/neat/support/uid"
	"github.com/dromara/carbon/v2"
)

// == CONSTRUCTOR =============================================================

// newVersioning creates a new versioning with a generated ID and current timestamp
func newVersioning() *versioning {
	o := &versioning{}
	o.SetID(neatuid.GenerateShortID())
	o.SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	o.SetSoftDeletedAt(VERSIONING_MAX_DATETIME)
	return o
}

// == CLASS ==================================================================

var _ VersioningInterface = (*versioning)(nil)

type versioning struct {
	orm.ShortID

	EntityTypeField string `db:"entity_type"`
	EntityIDField   string `db:"entity_id"`
	ContentField    string `db:"content"`

	CreatedAtField orm.CreatedAt
	soft_delete.SoftDeletesMaxDate
}

// == METHODS =================================================================

// IsSoftDeleted returns true if the versioning is soft deleted.
func (o *versioning) IsSoftDeleted() bool {
	return o.SoftDeletedAt.Before(time.Now().UTC())
}

// == SETTERS AND GETTERS =====================================================

// ID returns the id of the versioning.
func (o *versioning) ID() string {
	return o.ShortID.ID
}

// SetID sets the id of the versioning.
func (o *versioning) SetID(id string) VersioningInterface {
	o.ShortID.ID = id
	return o
}

// EntityType returns the entity type of the versioning.
func (o *versioning) EntityType() string {
	return o.EntityTypeField
}

// SetEntityType sets the entity type of the versioning.
func (o *versioning) SetEntityType(entityType string) VersioningInterface {
	o.EntityTypeField = entityType
	return o
}

// EntityID returns the entity id of the versioning.
func (o *versioning) EntityID() string {
	return o.EntityIDField
}

// SetEntityID sets the entity id of the versioning.
func (o *versioning) SetEntityID(entityID string) VersioningInterface {
	o.EntityIDField = entityID
	return o
}

// Content returns the content of the versioning.
func (o *versioning) Content() string {
	return o.ContentField
}

// SetContent sets the content of the versioning.
func (o *versioning) SetContent(content string) VersioningInterface {
	o.ContentField = content
	return o
}

// GetCreatedAt returns the created at time of the versioning.
func (o *versioning) GetCreatedAt() string {
	if o.CreatedAtField.CreatedAt.IsZero() {
		return ""
	}
	return carbon.CreateFromStdTime(o.CreatedAtField.CreatedAt).ToDateTimeString()
}

// GetCreatedAtCarbon returns the created at time of the versioning as a carbon object.
func (o *versioning) GetCreatedAtCarbon() *carbon.Carbon {
	return carbon.CreateFromStdTime(o.CreatedAtField.CreatedAt)
}

// SetCreatedAt sets the created at time of the versioning.
func (o *versioning) SetCreatedAt(createdAt string) VersioningInterface {
	if createdAt == "" {
		return o
	}
	o.CreatedAtField.CreatedAt = carbon.Parse(createdAt, carbon.UTC).StdTime()
	return o
}

// GetSoftDeletedAt returns the soft deleted at time of the versioning.
func (o *versioning) GetSoftDeletedAt() string {
	if o.SoftDeletedAt.IsZero() {
		return ""
	}
	return carbon.CreateFromStdTime(o.SoftDeletedAt).ToDateTimeString()
}

// GetSoftDeletedAtCarbon returns the soft deleted at time of the versioning as a carbon object.
func (o *versioning) GetSoftDeletedAtCarbon() *carbon.Carbon {
	return carbon.CreateFromStdTime(o.SoftDeletedAt)
}

// SetSoftDeletedAt sets the soft deleted at time of the versioning.
func (o *versioning) SetSoftDeletedAt(softDeletedAt string) VersioningInterface {
	if softDeletedAt == "" {
		return o
	}
	o.SoftDeletedAt = carbon.Parse(softDeletedAt, carbon.UTC).StdTime()
	return o
}
