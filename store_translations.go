package cmsstore

import (
	"context"
	"errors"
	"log"
	"strings"

	contractsorm "github.com/dracory/neat/contracts/database/orm"
	"github.com/dromara/carbon/v2"
)

func (store *storeImplementation) TranslationCount(ctx context.Context, options TranslationQueryInterface) (int64, error) {
	if store.neatDB == nil {
		return -1, errors.New("cms store: database is nil")
	}

	if !store.translationsEnabled {
		return -1, errors.New("translations are disabled")
	}

	q, _, err := store.translationSelectQuery(options)

	if err != nil {
		return -1, err
	}

	var count int64
	err = q.Table(store.translationTableName).Count(&count)
	return count, err
}

func (store *storeImplementation) TranslationCreate(ctx context.Context, translation TranslationInterface) error {
	if store.neatDB == nil {
		return errors.New("translationstore: database is nil")
	}

	if translation == nil {
		return errors.New("translation is nil")
	}

	if translation.CreatedAt() == "" {
		translation.SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	}

	if translation.UpdatedAt() == "" {
		translation.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	}

	return store.withTransaction(ctx, func(txCtx context.Context) error {
		data := translation.Data()

		if store.debugEnabled {
			log.Println("TranslationCreate:", data)
		}

		err := store.neatDB.Query().Table(store.translationTableName).Create(data)

		if err != nil {
			return err
		}

		translation.MarkAsNotDirty()

		return store.versioningTrackEntity(txCtx, VERSIONING_TYPE_TRANSLATION, translation.ID(), translation)
	})
}

func (store *storeImplementation) TranslationDelete(ctx context.Context, translation TranslationInterface) error {
	if store.neatDB == nil {
		return errors.New("cmsstore: database is nil")
	}

	if translation == nil {
		return errors.New("translation is nil")
	}

	return store.TranslationDeleteByID(ctx, translation.ID())
}

func (store *storeImplementation) TranslationDeleteByID(ctx context.Context, id string) error {
	if store.neatDB == nil {
		return errors.New("cmsstore: database is nil")
	}

	if id == "" {
		return errors.New("translation id is empty")
	}

	if store.debugEnabled {
		log.Println("TranslationDeleteByID:", id)
	}

	_, err := store.neatDB.Query().Table(store.translationTableName).Where("id = ?", id).Delete()

	return err
}

func (store *storeImplementation) TranslationFindByHandle(ctx context.Context, handle string) (translation TranslationInterface, err error) {
	if store.neatDB == nil {
		return nil, errors.New("cmsstore: database is nil")
	}

	if handle == "" {
		return nil, errors.New("translation handle is empty")
	}

	list, err := store.TranslationList(ctx, TranslationQuery().
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

func (store *storeImplementation) TranslationFindByID(ctx context.Context, id string) (translation TranslationInterface, err error) {
	if store.neatDB == nil {
		return nil, errors.New("cmsstore: database is nil")
	}

	if id == "" {
		return nil, errors.New("translation id is empty")
	}

	// Normalize ID to lowercase for consistent lookups
	id = NormalizeID(id)

	// Try direct lookup first (handles both 9-char and 32-char IDs)
	list, err := store.TranslationList(ctx, TranslationQuery().SetID(id).SetLimit(1))

	if err != nil {
		return nil, err
	}

	if len(list) > 0 {
		return list[0], nil
	}

	// If not found and ID looks shortened, try unshortening
	if IsShortID(id) {
		unshortenedID := UnshortenID(id)
		if unshortenedID != id {
			list, err = store.TranslationList(ctx, TranslationQuery().SetID(unshortenedID).SetLimit(1))
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

func (store *storeImplementation) TranslationFindByHandleOrID(ctx context.Context, handleOrID string, language string) (translation TranslationInterface, err error) {
	if store.neatDB == nil {
		return nil, errors.New("cmsstore: database is nil")
	}

	if handleOrID == "" {
		return nil, errors.New("translation id is empty")
	}

	list, err := store.TranslationList(ctx, TranslationQuery().
		SetHandleOrID(handleOrID).
		SetLimit(1))

	if err != nil {
		return nil, err
	}

	if len(list) > 0 {
		return list[0], nil
	}

	return nil, nil
}

func (store *storeImplementation) TranslationLanguageDefault() string {
	return store.translationLanguageDefault
}

func (store *storeImplementation) TranslationLanguages() map[string]string {
	return store.translationLanguages
}

func (store *storeImplementation) TranslationList(ctx context.Context, query TranslationQueryInterface) ([]TranslationInterface, error) {
	if store.neatDB == nil {
		return []TranslationInterface{}, errors.New("cmsstore: database is nil")
	}

	if !store.translationsEnabled {
		return []TranslationInterface{}, errors.New("translations are disabled")
	}

	q, _, err := store.translationSelectQuery(query)

	if err != nil {
		return []TranslationInterface{}, err
	}

	type translationRow struct {
		ID            string `db:"id"`
		SiteID        string `db:"site_id"`
		Name          string `db:"name"`
		Handle        string `db:"handle"`
		Status        string `db:"status"`
		Language      string `db:"language"`
		Content       string `db:"content"`
		CreatedAt     string `db:"created_at"`
		UpdatedAt     string `db:"updated_at"`
		SoftDeletedAt string `db:"soft_deleted_at"`
	}

	var rows []translationRow
	if err := q.Table(store.translationTableName).Get(&rows); err != nil {
		return []TranslationInterface{}, err
	}

	list := make([]TranslationInterface, 0, len(rows))
	for _, r := range rows {
		modelMap := map[string]string{
			"id":              r.ID,
			"site_id":         r.SiteID,
			"name":            r.Name,
			"handle":          r.Handle,
			"status":          r.Status,
			"language":        r.Language,
			"content":         r.Content,
			"created_at":      r.CreatedAt,
			"updated_at":      r.UpdatedAt,
			"soft_deleted_at": r.SoftDeletedAt,
		}
		model := NewTranslationFromExistingData(modelMap)
		list = append(list, model)
	}

	return list, nil
}

func (store *storeImplementation) TranslationSoftDelete(ctx context.Context, translation TranslationInterface) error {
	if store.neatDB == nil {
		return errors.New("cmsstore: database is nil")
	}

	if translation == nil {
		return errors.New("translation is nil")
	}

	translation.SetSoftDeletedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))

	return store.TranslationUpdate(ctx, translation)
}

func (store *storeImplementation) TranslationSoftDeleteByID(ctx context.Context, id string) error {
	if store.neatDB == nil {
		return errors.New("cmsstore: database is nil")
	}

	if id == "" {
		return errors.New("translation id is empty")
	}

	translation, err := store.TranslationFindByID(ctx, id)

	if err != nil {
		return err
	}

	if translation == nil {
		return errors.New("translation not found")
	}

	return store.TranslationSoftDelete(ctx, translation)
}

func (store *storeImplementation) TranslationUpdate(ctx context.Context, translation TranslationInterface) error {
	if store.neatDB == nil {
		return errors.New("cmsstore: database is nil")
	}

	if translation == nil {
		return errors.New("translation is nil")
	}

	translation.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString())

	return store.withTransaction(ctx, func(txCtx context.Context) error {
		dataChanged := translation.DataChanged()

		delete(dataChanged, COLUMN_ID) // ID is not updateable

		if len(dataChanged) < 1 {
			return nil
		}

		if store.debugEnabled {
			log.Println("TranslationUpdate:", dataChanged)
		}

		_, err := store.neatDB.Query().Table(store.translationTableName).Where("id = ?", translation.ID()).Update(dataChanged)

		if err != nil {
			return err
		}

		translation.MarkAsNotDirty()

		return store.versioningTrackEntity(txCtx, VERSIONING_TYPE_TRANSLATION, translation.ID(), translation)
	})
}

func (store *storeImplementation) translationSelectQuery(options TranslationQueryInterface) (query contractsorm.Query, columns []any, err error) {
	if options == nil {
		return nil, []any{}, errors.New("translation query cannot be nil")
	}

	if err := options.Validate(); err != nil {
		return nil, []any{}, err
	}

	q := store.neatDB.Query().Table(store.translationTableName)

	if options.HasCreatedAtGte() && options.HasCreatedAtLte() {
		q = q.Where(COLUMN_CREATED_AT+" >= ? AND "+COLUMN_CREATED_AT+" <= ?", options.CreatedAtGte(), options.CreatedAtLte())
	} else if options.HasCreatedAtGte() {
		q = q.Where(COLUMN_CREATED_AT+" >= ?", options.CreatedAtGte())
	} else if options.HasCreatedAtLte() {
		q = q.Where(COLUMN_CREATED_AT+" <= ?", options.CreatedAtLte())
	}

	if options.HasHandle() {
		q = q.Where(COLUMN_HANDLE+" = ?", options.Handle())
	}

	if options.HasHandleOrID() {
		q = q.Where("("+COLUMN_HANDLE+" = ? OR "+COLUMN_ID+" = ?)", options.HandleOrID(), options.HandleOrID())
	}

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

	if options.HasNameLike() {
		q = q.Where(COLUMN_NAME+" LIKE ?", options.NameLike())
	}

	if options.HasSiteID() {
		q = q.Where(COLUMN_SITE_ID+" = ?", options.SiteID())
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

	columns = []any{}

	for _, column := range options.Columns() {
		columns = append(columns, column)
	}

	if options.SoftDeletedIncluded() {
		return q, columns, nil // soft deleted translations requested specifically
	}

	q = q.Where(COLUMN_SOFT_DELETED_AT+" > ?", carbon.Now(carbon.UTC).ToDateTimeString())

	return q, columns, nil
}
