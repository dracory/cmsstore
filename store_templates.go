package cmsstore

import (
	"context"
	"errors"
	"log"
	"strings"

	contractsorm "github.com/dracory/neat/contracts/database/orm"
	"github.com/dromara/carbon/v2"
)

func (store *storeImplementation) TemplateCount(ctx context.Context, options TemplateQueryInterface) (int64, error) {
	if store.neatDB == nil {
		return -1, errors.New("cms store: database is nil")
	}

	q, _, err := store.templateSelectQuery(options)

	if err != nil {
		return -1, err
	}

	var count int64
	err = q.Table(store.templateTableName).Count(&count)
	return count, err
}

func (store *storeImplementation) TemplateCreate(ctx context.Context, template TemplateInterface) error {
	if store.neatDB == nil {
		return errors.New("templatestore: database is nil")
	}

	if template == nil {
		return errors.New("template is nil")
	}
	if template.CreatedAt() == "" {
		template.SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	}

	if template.UpdatedAt() == "" {
		template.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	}

	return store.withTransaction(ctx, func(txCtx context.Context) error {
		data := template.Data()

		if store.debugEnabled {
			log.Println("TemplateCreate:", data)
		}

		err := store.neatDB.Query().Table(store.templateTableName).Create(data)

		if err != nil {
			return err
		}

		template.MarkAsNotDirty()

		return store.versioningTrackEntity(txCtx, VERSIONING_TYPE_TEMPLATE, template.ID(), template)
	})
}

func (store *storeImplementation) TemplateDelete(ctx context.Context, template TemplateInterface) error {
	if store.neatDB == nil {
		return errors.New("templatestore: database is nil")
	}

	if template == nil {
		return errors.New("template is nil")
	}

	return store.TemplateDeleteByID(ctx, template.ID())
}

func (store *storeImplementation) TemplateDeleteByID(ctx context.Context, id string) error {
	if store.neatDB == nil {
		return errors.New("templatestore: database is nil")
	}

	if id == "" {
		return errors.New("template id is empty")
	}

	if store.debugEnabled {
		log.Println("TemplateDeleteByID:", id)
	}

	_, err := store.neatDB.Query().Table(store.templateTableName).Where("id = ?", id).Delete()

	return err
}

func (store *storeImplementation) TemplateFindByHandle(ctx context.Context, handle string) (template TemplateInterface, err error) {
	if store.neatDB == nil {
		return nil, errors.New("templatestore: database is nil")
	}

	if handle == "" {
		return nil, errors.New("template handle is empty")
	}

	list, err := store.TemplateList(ctx, TemplateQuery().
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

func (store *storeImplementation) TemplateFindByID(ctx context.Context, id string) (template TemplateInterface, err error) {
	if id == "" {
		return nil, errors.New("template id is empty")
	}

	// Normalize ID to lowercase for consistent lookups
	id = NormalizeID(id)

	// Try direct lookup first (handles both 9-char and 32-char IDs)
	list, err := store.TemplateList(ctx, TemplateQuery().SetID(id).SetLimit(1))

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
			list, err = store.TemplateList(ctx, TemplateQuery().SetID(unshortenedID).SetLimit(1))
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

func (store *storeImplementation) TemplateList(ctx context.Context, query TemplateQueryInterface) ([]TemplateInterface, error) {
	if store.neatDB == nil {
		return []TemplateInterface{}, errors.New("templatestore: database is nil")
	}

	q, _, err := store.templateSelectQuery(query)

	if err != nil {
		return []TemplateInterface{}, err
	}

	type templateRow struct {
		ID            string `db:"id"`
		SiteID        string `db:"site_id"`
		Name          string `db:"name"`
		Handle        string `db:"handle"`
		Status        string `db:"status"`
		Content       string `db:"content"`
		Editor        string `db:"editor"`
		Metas         string `db:"metas"`
		Memo          string `db:"memo"`
		CreatedAt     string `db:"created_at"`
		UpdatedAt     string `db:"updated_at"`
		SoftDeletedAt string `db:"soft_deleted_at"`
	}

	var rows []templateRow
	if err := q.Table(store.templateTableName).Get(&rows); err != nil {
		return []TemplateInterface{}, err
	}

	list := make([]TemplateInterface, 0, len(rows))
	for _, r := range rows {
		modelMap := map[string]string{
			"id":              r.ID,
			"site_id":         r.SiteID,
			"name":            r.Name,
			"handle":          r.Handle,
			"status":          r.Status,
			"content":         r.Content,
			"editor":          r.Editor,
			"metas":           r.Metas,
			"memo":            r.Memo,
			"created_at":      r.CreatedAt,
			"updated_at":      r.UpdatedAt,
			"soft_deleted_at": r.SoftDeletedAt,
		}
		model := NewTemplateFromExistingData(modelMap)
		list = append(list, model)
	}

	return list, nil
}

func (store *storeImplementation) TemplateSoftDelete(ctx context.Context, template TemplateInterface) error {
	if template == nil {
		return errors.New("template is nil")
	}

	template.SetSoftDeletedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))

	return store.TemplateUpdate(ctx, template)
}

func (store *storeImplementation) TemplateSoftDeleteByID(ctx context.Context, id string) error {
	if store.neatDB == nil {
		return errors.New("templatestore: database is nil")
	}

	template, err := store.TemplateFindByID(ctx, id)

	if err != nil {
		return err
	}

	if template == nil {
		return errors.New("template not found")
	}

	return store.TemplateSoftDelete(ctx, template)
}

func (store *storeImplementation) TemplateUpdate(ctx context.Context, template TemplateInterface) error {
	if store.neatDB == nil {
		return errors.New("templatestore: database is nil")
	}

	if template == nil {
		return errors.New("template is nil")
	}

	template.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString())

	return store.withTransaction(ctx, func(txCtx context.Context) error {
		dataChanged := template.DataChanged()

		delete(dataChanged, COLUMN_ID) // ID is not updateable

		if len(dataChanged) < 1 {
			return nil
		}

		if store.debugEnabled {
			log.Println("TemplateUpdate:", dataChanged)
		}

		_, err := store.neatDB.Query().Table(store.templateTableName).Where("id = ?", template.ID()).Update(dataChanged)

		if err != nil {
			return err
		}

		template.MarkAsNotDirty()

		return store.versioningTrackEntity(txCtx, VERSIONING_TYPE_TEMPLATE, template.ID(), template)
	})
}

func (store *storeImplementation) templateSelectQuery(options TemplateQueryInterface) (query contractsorm.Query, columns []any, err error) {
	if options == nil {
		return nil, nil, errors.New("template query cannot be nil")
	}

	if err := options.Validate(); err != nil {
		return nil, nil, err
	}

	q := store.neatDB.Query().Table(store.templateTableName)

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
		return q, columns, nil // soft deleted templates requested specifically
	}

	q = q.Where(COLUMN_SOFT_DELETED_AT+" > ?", carbon.Now(carbon.UTC).ToDateTimeString())

	return q, columns, nil
}
