package cmsstore

import (
	"context"
	"errors"
	"log"
	"strconv"
	"strings"

	contractsorm "github.com/dracory/neat/contracts/database/orm"
	"github.com/dromara/carbon/v2"
)

func (store *storeImplementation) MediaCount(ctx context.Context, options MediaQueryInterface) (int64, error) {
	if store.neatDB == nil {
		return -1, errors.New("cms store: database is nil")
	}

	if options != nil && !options.IsCountOnly() {
		options.SetCountOnly(true)
	}

	q, _, err := store.mediaSelectQuery(options)

	if err != nil {
		return -1, err
	}

	var count int64
	err = q.Table(store.mediaTableName).Count(&count)
	return count, err
}

func (store *storeImplementation) MediaCreate(ctx context.Context, media MediaInterface) error {
	if store.neatDB == nil {
		return errors.New("cms store: database is nil")
	}

	if media == nil {
		return errors.New("media is nil")
	}

	media.SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	media.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))

	return store.withTransaction(ctx, func(txCtx context.Context) error {
		data := media.Data()

		if store.debugEnabled {
			log.Println("MediaCreate:", data)
		}

		err := store.neatDB.Query().Table(store.mediaTableName).Create(data)

		if err != nil {
			return err
		}

		media.MarkAsNotDirty()

		return store.versioningTrackEntity(txCtx, VERSIONING_TYPE_MEDIA, media.ID(), media)
	})
}

func (store *storeImplementation) MediaDelete(ctx context.Context, media MediaInterface) error {
	if store.neatDB == nil {
		return errors.New("cms store: database is nil")
	}

	if media == nil {
		return errors.New("media is nil")
	}

	return store.MediaDeleteByID(ctx, media.ID())
}

func (store *storeImplementation) MediaDeleteByID(ctx context.Context, id string) error {
	if store.neatDB == nil {
		return errors.New("cms store: database is nil")
	}

	if id == "" {
		return errors.New("media id is empty")
	}

	if store.debugEnabled {
		log.Println("MediaDeleteByID:", id)
	}

	_, err := store.neatDB.Query().Table(store.mediaTableName).Where("id = ?", id).Delete()

	return err
}

func (store *storeImplementation) MediaFindByHandle(ctx context.Context, handle string) (MediaInterface, error) {
	if store.neatDB == nil {
		return nil, errors.New("cms store: database is nil")
	}

	if handle == "" {
		return nil, errors.New("media handle is empty")
	}

	list, err := store.MediaList(ctx, MediaQuery().
		SetHandle(handle).
		SetLimit(1))

	if err != nil {
		return nil, err
	}

	if len(list) > 0 {
		return list[0], nil
	}

	return nil, nil
}

func (store *storeImplementation) MediaFindByID(ctx context.Context, id string) (MediaInterface, error) {
	if id == "" {
		return nil, errors.New("media id is empty")
	}

	id = NormalizeID(id)

	list, err := store.MediaList(ctx, MediaQuery().SetID(id).SetLimit(1))

	if err != nil {
		return nil, err
	}

	if len(list) > 0 {
		return list[0], nil
	}

	if IsShortID(id) {
		unshortenedID := UnshortenID(id)
		if unshortenedID != id {
			list, err = store.MediaList(ctx, MediaQuery().SetID(unshortenedID).SetLimit(1))
			if err != nil {
				return nil, err
			}
			if len(list) > 0 {
				return list[0], nil
			}
		}
	}

	return nil, nil
}

func (store *storeImplementation) MediaList(ctx context.Context, query MediaQueryInterface) ([]MediaInterface, error) {
	if store.neatDB == nil {
		return []MediaInterface{}, errors.New("cms store: database is nil")
	}

	q, _, err := store.mediaSelectQuery(query)

	if err != nil {
		return []MediaInterface{}, err
	}

	type mediaRow struct {
		ID            string `db:"id"`
		EntityID      string `db:"entity_id"`
		EntityType    string `db:"entity_type"`
		SiteID        string `db:"site_id"`
		Title         string `db:"title"`
		Description   string `db:"description"`
		Memo          string `db:"memo"`
		MediaURL      string `db:"media_url"`
		MediaType     string `db:"media_type"`
		FileSize      string `db:"file_size"`
		FileExtension string `db:"file_extension"`
		Sequence      int    `db:"sequence"`
		Status        string `db:"status"`
		Handle        string `db:"handle"`
		Metas         string `db:"metas"`
		CreatedAt     string `db:"created_at"`
		UpdatedAt     string `db:"updated_at"`
		SoftDeletedAt string `db:"soft_deleted_at"`
	}

	var rows []mediaRow
	if err := q.Table(store.mediaTableName).Get(&rows); err != nil {
		return []MediaInterface{}, err
	}

	list := make([]MediaInterface, 0, len(rows))
	for _, r := range rows {
		modelMap := map[string]string{
			"id":              r.ID,
			"entity_id":       r.EntityID,
			"entity_type":     r.EntityType,
			"site_id":         r.SiteID,
			"title":           r.Title,
			"description":     r.Description,
			"memo":            r.Memo,
			"media_url":       r.MediaURL,
			"media_type":      r.MediaType,
			"file_size":       r.FileSize,
			"file_extension":  r.FileExtension,
			"sequence":        strconv.Itoa(r.Sequence),
			"status":          r.Status,
			"handle":          r.Handle,
			"metas":           r.Metas,
			"created_at":      r.CreatedAt,
			"updated_at":      r.UpdatedAt,
			"soft_deleted_at": r.SoftDeletedAt,
		}
		model := NewMediaFromExistingData(modelMap)
		list = append(list, model)
	}

	return list, nil
}

func (store *storeImplementation) MediaListByEntityID(ctx context.Context, entityID string, entityType string) ([]MediaInterface, error) {
	return store.MediaList(ctx, MediaQuery().
		SetEntityID(entityID).
		SetEntityType(entityType).
		SetOrderBy(COLUMN_SEQUENCE).
		SetSortOrder(SORT_ORDER_ASC))
}

func (store *storeImplementation) MediaSoftDelete(ctx context.Context, media MediaInterface) error {
	if media == nil {
		return errors.New("media is nil")
	}

	media.SetSoftDeletedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))

	return store.MediaUpdate(ctx, media)
}

func (store *storeImplementation) MediaSoftDeleteByID(ctx context.Context, id string) error {
	if store.neatDB == nil {
		return errors.New("cms store: database is nil")
	}

	media, err := store.MediaFindByID(ctx, id)

	if err != nil {
		return err
	}

	if media == nil {
		return errors.New("media not found")
	}

	return store.MediaSoftDelete(ctx, media)
}

func (store *storeImplementation) MediaUpdate(ctx context.Context, media MediaInterface) error {
	if store.neatDB == nil {
		return errors.New("cms store: database is nil")
	}

	if media == nil {
		return errors.New("media is nil")
	}

	media.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString())

	return store.withTransaction(ctx, func(txCtx context.Context) error {
		dataChanged := media.DataChanged()

		delete(dataChanged, COLUMN_ID)

		if len(dataChanged) < 1 {
			return nil
		}

		if store.debugEnabled {
			log.Println("MediaUpdate:", dataChanged)
		}

		_, err := store.neatDB.Query().Table(store.mediaTableName).Where("id = ?", media.ID()).Update(dataChanged)
		if err != nil {
			return err
		}

		media.MarkAsNotDirty()

		return store.versioningTrackEntity(txCtx, VERSIONING_TYPE_MEDIA, media.ID(), media)
	})
}

func (store *storeImplementation) mediaSelectQuery(options MediaQueryInterface) (query contractsorm.Query, selectColumns []any, err error) {
	if options == nil {
		return nil, []any{}, errors.New("media options cannot be nil")
	}

	if err := options.Validate(); err != nil {
		return nil, []any{}, err
	}

	q := store.neatDB.Query().Table(store.mediaTableName)

	if options.HasID() {
		q = q.Where(COLUMN_ID+" = ?", options.ID())
	}

	if options.HasIDIn() {
		idIn := options.IDIn()
		if len(idIn) > 0 {
			placeholders := make([]string, len(idIn))
			args := make([]any, len(idIn))
			for i, v := range idIn {
				placeholders[i] = "?"
				args[i] = v
			}
			q = q.Where(COLUMN_ID+" IN ("+strings.Join(placeholders, ", ")+")", args...)
		}
	}

	if options.HasEntityID() {
		q = q.Where(COLUMN_ENTITY_ID+" = ?", options.EntityID())
	}

	if options.HasEntityType() {
		q = q.Where(COLUMN_ENTITY_TYPE+" = ?", options.EntityType())
	}

	if options.HasSiteID() {
		q = q.Where(COLUMN_SITE_ID+" = ?", options.SiteID())
	}

	if options.HasHandle() {
		q = q.Where(COLUMN_HANDLE+" = ?", options.Handle())
	}

	if options.HasExtension() {
		q = q.Where(COLUMN_FILE_EXTENSION+" = ?", options.Extension())
	}

	if options.HasType() {
		q = q.Where(COLUMN_MEDIA_TYPE+" = ?", options.Type())
	}

	if options.HasStatus() {
		q = q.Where(COLUMN_STATUS+" = ?", options.Status())
	}

	if options.HasStatusIn() {
		statusIn := options.StatusIn()
		if len(statusIn) > 0 {
			placeholders := make([]string, len(statusIn))
			args := make([]any, len(statusIn))
			for i, v := range statusIn {
				placeholders[i] = "?"
				args[i] = v
			}
			q = q.Where(COLUMN_STATUS+" IN ("+strings.Join(placeholders, ", ")+")", args...)
		}
	}

	if options.HasNameLike() {
		q = q.Where(COLUMN_TITLE+" LIKE ?", options.NameLike())
	}

	if !options.IsCountOnly() {
		if options.HasLimit() {
			q = q.Limit(options.Limit())
		}

		if options.HasOffset() {
			q = q.Offset(options.Offset())
		}
	}

	sortOrder := SORT_ORDER_DESC
	if options.HasSortOrder() {
		sortOrder = options.SortOrder()
	}

	if !options.IsCountOnly() && options.HasOrderBy() {
		if strings.EqualFold(sortOrder, SORT_ORDER_ASC) {
			q = q.OrderBy(options.OrderBy(), "ASC")
		} else {
			q = q.OrderBy(options.OrderBy(), "DESC")
		}
	}

	columns := []any{}

	for _, column := range options.Columns() {
		columns = append(columns, column)
	}

	if options.SoftDeletedIncluded() {
		return q, columns, nil
	}

	q = q.Where(COLUMN_SOFT_DELETED_AT+" > ?", carbon.Now(carbon.UTC).ToDateTimeString())

	return q, columns, nil
}
