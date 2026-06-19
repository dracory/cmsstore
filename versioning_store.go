package cmsstore

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/dracory/neat"
	contractsorm "github.com/dracory/neat/contracts/database/orm"
	contractsschema "github.com/dracory/neat/contracts/database/schema"
	"github.com/dromara/carbon/v2"
)

// == CONSTRUCTOR =============================================================

// newVersioningStore creates a new versioning store backed by the provided neat database.
func newVersioningStore(neatDB *neat.Database, tableName string) (*versioningStore, error) {
	if neatDB == nil {
		return nil, errors.New("versioning store: database is nil")
	}

	if tableName == "" {
		return nil, errors.New("versioning store: tableName is required")
	}

	return &versioningStore{
		tableName: tableName,
		db:        neatDB,
	}, nil
}

// == CLASS ===================================================================

type versioningStore struct {
	db        *neat.Database
	tableName string
}

// MigrateUp creates the versioning table
func (store *versioningStore) MigrateUp(ctx context.Context, tx ...*sql.Tx) error {
	if store.db.Schema().HasTable(store.tableName) {
		return nil
	}

	return store.db.Schema().Create(store.tableName, func(table contractsschema.Blueprint) {
		table.String(COLUMN_ID, 21)
		table.Primary(COLUMN_ID)
		table.String(COLUMN_ENTITY_TYPE, 40)
		table.String(COLUMN_ENTITY_ID, 40)
		table.Text(COLUMN_CONTENT)
		table.DateTime(COLUMN_CREATED_AT)
		table.DateTime(COLUMN_SOFT_DELETED_AT)
	})
}

// MigrateDown drops the versioning table
func (store *versioningStore) MigrateDown(ctx context.Context, tx ...*sql.Tx) error {
	if !store.db.Schema().HasTable(store.tableName) {
		return nil
	}

	return store.db.Schema().Drop(store.tableName)
}

// VersionCreate creates a new versioning
func (store *versioningStore) VersionCreate(ctx context.Context, version VersioningInterface) error {
	if ctx == nil {
		return errors.New("ctx is nil")
	}
	if version == nil {
		return errors.New("versioning store: version cannot be nil")
	}
	if version.ID() == "" {
		return errors.New("versioning store: version id should not be empty")
	}
	if version.EntityType() == "" {
		return errors.New("versioning store: version entity type should not be empty")
	}
	if version.EntityID() == "" {
		return errors.New("versioning store: version entity id should not be empty")
	}
	if version.GetCreatedAt() == "" {
		version.SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString())
	}
	if version.GetSoftDeletedAt() == "" {
		version.SetSoftDeletedAt(VERSIONING_MAX_DATETIME)
	}

	row := map[string]any{
		COLUMN_ID:              version.ID(),
		COLUMN_ENTITY_TYPE:     version.EntityType(),
		COLUMN_ENTITY_ID:       version.EntityID(),
		COLUMN_CONTENT:         version.Content(),
		COLUMN_CREATED_AT:      version.GetCreatedAtCarbon().StdTime(),
		COLUMN_SOFT_DELETED_AT: version.GetSoftDeletedAtCarbon().StdTime(),
	}

	return store.db.Query().Table(store.tableName).Create(row)
}

// VersionDelete deletes a versioning permanently
func (store *versioningStore) VersionDelete(ctx context.Context, version VersioningInterface) error {
	if ctx == nil {
		return errors.New("ctx is nil")
	}
	if version == nil {
		return errors.New("version is nil")
	}

	return store.VersionDeleteByID(ctx, version.ID())
}

// VersionDeleteByID deletes a versioning by ID permanently
func (store *versioningStore) VersionDeleteByID(ctx context.Context, id string) error {
	if ctx == nil {
		return errors.New("ctx is nil")
	}
	if id == "" {
		return errors.New("versioning id is empty")
	}

	_, err := store.db.Query().
		Table(store.tableName).
		Where(COLUMN_ID+" = ?", id).
		Delete()
	return err
}

// VersionFindByID finds a versioning by ID
func (store *versioningStore) VersionFindByID(ctx context.Context, id string) (VersioningInterface, error) {
	if id == "" {
		return nil, errors.New("versioning store: version id is required")
	}

	list, err := store.VersionList(ctx, newVersioningQuery().SetID(id).SetLimit(1))
	if err != nil {
		return nil, err
	}
	if len(list) > 0 {
		return list[0], nil
	}
	return nil, nil
}

// VersionList returns a list of versionings matching the query options
func (store *versioningStore) VersionList(ctx context.Context, options VersioningQueryInterface) ([]VersioningInterface, error) {
	if ctx == nil {
		return nil, errors.New("ctx is nil")
	}

	type versionRow struct {
		ID            string    `db:"id"`
		EntityType    string    `db:"entity_type"`
		EntityID      string    `db:"entity_id"`
		Content       string    `db:"content"`
		CreatedAt     time.Time `db:"created_at"`
		SoftDeletedAt time.Time `db:"soft_deleted_at"`
	}

	q := store.buildQuery(options)
	q = q.Table(store.tableName)

	if len(options.Columns()) > 0 {
		q = q.Select(options.Columns())
	}

	var rows []versionRow
	if err := q.Get(&rows); err != nil {
		return []VersioningInterface{}, err
	}

	list := make([]VersioningInterface, 0, len(rows))
	for _, r := range rows {
		v := &versioning{}
		v.SetID(r.ID)
		v.SetEntityType(r.EntityType)
		v.SetEntityID(r.EntityID)
		v.SetContent(r.Content)
		v.CreatedAtField.CreatedAt = r.CreatedAt
		v.SoftDeletedAt = r.SoftDeletedAt
		list = append(list, v)
	}

	return list, nil
}

// VersionSoftDelete soft deletes a versioning
func (store *versioningStore) VersionSoftDelete(ctx context.Context, version VersioningInterface) error {
	if ctx == nil {
		return errors.New("ctx is nil")
	}
	if version == nil {
		return errors.New("version is nil")
	}

	version.SetSoftDeletedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))

	return store.VersionUpdate(ctx, version)
}

// VersionSoftDeleteByID soft deletes a versioning by ID
func (store *versioningStore) VersionSoftDeleteByID(ctx context.Context, id string) error {
	if ctx == nil {
		return errors.New("ctx is nil")
	}
	if id == "" {
		return errors.New("versioning id is empty")
	}

	version, err := store.VersionFindByID(ctx, id)
	if err != nil {
		return err
	}
	if version == nil {
		return errors.New("versioning not found")
	}

	return store.VersionSoftDelete(ctx, version)
}

// VersionUpdate updates a versioning.
//
// Note!! There is no reason to call this method other than marking
// the versioning as soft deleted
func (store *versioningStore) VersionUpdate(ctx context.Context, version VersioningInterface) error {
	if ctx == nil {
		return errors.New("ctx is nil")
	}
	if version == nil {
		return errors.New("version is nil")
	}

	row := map[string]any{
		COLUMN_SOFT_DELETED_AT: version.GetSoftDeletedAtCarbon().StdTime(),
	}

	_, err := store.db.Query().Table(store.tableName).Where(COLUMN_ID+" = ?", version.ID()).Update(row)
	return err
}

// == QUERY BUILDER ==========================================================

// buildQuery builds a neat query from the versioning query interface.
func (store *versioningStore) buildQuery(options VersioningQueryInterface) contractsorm.Query {
	q := store.db.Query()

	if options == nil {
		return q
	}

	if options.HasID() && options.ID() != "" {
		q = q.Where(COLUMN_ID+" = ?", options.ID())
	}

	if options.HasEntityType() && options.EntityType() != "" {
		q = q.Where(COLUMN_ENTITY_TYPE+" = ?", options.EntityType())
	}

	if options.HasEntityID() && options.EntityID() != "" {
		q = q.Where(COLUMN_ENTITY_ID+" = ?", options.EntityID())
	}

	if options.HasLimit() && options.Limit() > 0 {
		q = q.Limit(options.Limit())
	}

	if options.HasOffset() && options.Offset() > 0 {
		q = q.Offset(int(options.Offset()))
	}

	if options.HasOrderBy() && options.OrderBy() != "" {
		if options.HasSortOrder() && options.SortOrder() == "asc" {
			q = q.OrderBy(options.OrderBy())
		} else {
			q = q.OrderByDesc(options.OrderBy())
		}
	}

	if options.HasSoftDeletedIncluded() && options.SoftDeletedIncluded() {
		q = q.WithSoftDeleted()
	} else {
		q = q.Where(COLUMN_SOFT_DELETED_AT+" = ?", carbon.Parse(VERSIONING_MAX_DATETIME, carbon.UTC).StdTime())
	}

	return q
}
