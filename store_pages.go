package cmsstore

import (
	"context"
	"errors"
	"log"
	"strings"

	contractsorm "github.com/dracory/neat/contracts/database/orm"
	"github.com/dromara/carbon/v2"
)

func (store *storeImplementation) PageCount(ctx context.Context, options PageQueryInterface) (int64, error) {
	if store.neatDB == nil {
		return -1, errors.New("cms store: database is nil")
	}

	if options != nil && !options.IsCountOnly() {
		options.SetCountOnly(true)
	}

	q, _, err := store.pageSelectQuery(options)

	if err != nil {
		return -1, err
	}

	var count int64
	err = q.Table(store.pageTableName).Count(&count)
	return count, err
}

func (store *storeImplementation) PageCreate(ctx context.Context, page PageInterface) error {
	if store.neatDB == nil {
		return errors.New("pagestore: database is nil")
	}

	if page == nil {
		return errors.New("page is nil")
	}

	page.SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	page.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))

	return store.withTransaction(ctx, func(txCtx context.Context) error {
		data := page.Data()

		if store.debugEnabled {
			log.Println("PageCreate:", data)
		}

		err := store.neatDB.Query().Table(store.pageTableName).Create(data)

		if err != nil {
			return err
		}

		page.MarkAsNotDirty()

		return store.versioningTrackEntity(txCtx, VERSIONING_TYPE_PAGE, page.ID(), page)
	})
}

func (store *storeImplementation) PageDelete(ctx context.Context, page PageInterface) error {
	if store.neatDB == nil {
		return errors.New("pagestore: database is nil")
	}

	if page == nil {
		return errors.New("page is nil")
	}

	return store.PageDeleteByID(ctx, page.ID())
}

func (store *storeImplementation) PageDeleteByID(ctx context.Context, id string) error {
	if store.neatDB == nil {
		return errors.New("pagestore: database is nil")
	}

	if id == "" {
		return errors.New("page id is empty")
	}

	if store.debugEnabled {
		log.Println("PageDeleteByID:", id)
	}

	_, err := store.neatDB.Query().Table(store.pageTableName).Where("id = ?", id).Delete()

	return err
}

func (store *storeImplementation) PageFindByHandle(ctx context.Context, handle string) (page PageInterface, err error) {
	if store.neatDB == nil {
		return nil, errors.New("pagestore: database is nil")
	}

	if handle == "" {
		return nil, errors.New("page handle is empty")
	}

	list, err := store.PageList(ctx, PageQuery().
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

func (store *storeImplementation) PageFindByID(ctx context.Context, id string) (page PageInterface, err error) {
	if id == "" {
		return nil, errors.New("page id is empty")
	}

	// Normalize ID to lowercase for consistent lookups
	id = NormalizeID(id)

	// Try direct lookup first (handles both 9-char and 32-char IDs)
	list, err := store.PageList(ctx, PageQuery().SetID(id).SetLimit(1))

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
			list, err = store.PageList(ctx, PageQuery().SetID(unshortenedID).SetLimit(1))
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

func (store *storeImplementation) PageList(ctx context.Context, query PageQueryInterface) ([]PageInterface, error) {
	if store.neatDB == nil {
		return []PageInterface{}, errors.New("pagestore: database is nil")
	}

	q, _, err := store.pageSelectQuery(query)

	if err != nil {
		return []PageInterface{}, err
	}

	type pageRow struct {
		ID                string `db:"id"`
		SiteID            string `db:"site_id"`
		TemplateID        string `db:"template_id"`
		Name              string `db:"name"`
		Handle            string `db:"handle"`
		Alias             string `db:"alias"`
		Status            string `db:"status"`
		Title             string `db:"title"`
		Content           string `db:"content"`
		Editor            string `db:"editor"`
		CanonicalURL      string `db:"canonical_url"`
		MetaKeywords      string `db:"meta_keywords"`
		MetaDescription   string `db:"meta_description"`
		MetaRobots        string `db:"meta_robots"`
		MiddlewaresAfter  string `db:"middlewares_after"`
		MiddlewaresBefore string `db:"middlewares_before"`
		Metas             string `db:"metas"`
		Memo              string `db:"memo"`
		CreatedAt         string `db:"created_at"`
		UpdatedAt         string `db:"updated_at"`
		SoftDeletedAt     string `db:"soft_deleted_at"`
	}

	var rows []pageRow
	if err := q.Table(store.pageTableName).Get(&rows); err != nil {
		return []PageInterface{}, err
	}

	list := make([]PageInterface, 0, len(rows))
	for _, r := range rows {
		modelMap := map[string]string{
			"id":                 r.ID,
			"site_id":            r.SiteID,
			"template_id":        r.TemplateID,
			"name":               r.Name,
			"handle":             r.Handle,
			"alias":              r.Alias,
			"status":             r.Status,
			"title":              r.Title,
			"content":            r.Content,
			"editor":             r.Editor,
			"canonical_url":      r.CanonicalURL,
			"meta_keywords":      r.MetaKeywords,
			"meta_description":   r.MetaDescription,
			"meta_robots":        r.MetaRobots,
			"middlewares_after":  r.MiddlewaresAfter,
			"middlewares_before": r.MiddlewaresBefore,
			"metas":              r.Metas,
			"memo":               r.Memo,
			"created_at":         r.CreatedAt,
			"updated_at":         r.UpdatedAt,
			"soft_deleted_at":    r.SoftDeletedAt,
		}
		model := NewPageFromExistingData(modelMap)
		list = append(list, model)
	}

	return list, nil
}

func (store *storeImplementation) PageSoftDelete(ctx context.Context, page PageInterface) error {
	if page == nil {
		return errors.New("page is nil")
	}

	page.SetSoftDeletedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))

	return store.PageUpdate(ctx, page)
}

func (store *storeImplementation) PageSoftDeleteByID(ctx context.Context, id string) error {
	if store.neatDB == nil {
		return errors.New("pagestore: database is nil")
	}

	page, err := store.PageFindByID(ctx, id)

	if err != nil {
		return err
	}

	if page == nil {
		return errors.New("page not found")
	}

	return store.PageSoftDelete(ctx, page)
}

func (store *storeImplementation) PageUpdate(ctx context.Context, page PageInterface) error {
	if store.neatDB == nil {
		return errors.New("pagestore: database is nil")
	}

	if page == nil {
		return errors.New("page is nil")
	}

	page.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString())

	return store.withTransaction(ctx, func(txCtx context.Context) error {
		dataChanged := page.DataChanged()

		delete(dataChanged, COLUMN_ID) // ID is not updateable

		if len(dataChanged) < 1 {
			return nil
		}

		if store.debugEnabled {
			log.Println("PageUpdate:", dataChanged)
		}

		_, err := store.neatDB.Query().Table(store.pageTableName).Where("id = ?", page.ID()).Update(dataChanged)
		if err != nil {
			return err
		}

		page.MarkAsNotDirty()

		return store.versioningTrackEntity(txCtx, VERSIONING_TYPE_PAGE, page.ID(), page)
	})
}

func (store *storeImplementation) pageSelectQuery(options PageQueryInterface) (query contractsorm.Query, selectColumns []any, err error) {
	if options == nil {
		return nil, []any{}, errors.New("page options cannot be nil")
	}

	if err := options.Validate(); err != nil {
		return nil, []any{}, err
	}

	q := store.neatDB.Query().Table(store.pageTableName)

	if options.HasAlias() {
		q = q.Where(COLUMN_ALIAS+" = ?", options.Alias())
	}

	if options.HasAliasLike() {
		q = q.Where(COLUMN_ALIAS+" LIKE ?", options.AliasLike())
	}

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

	if options.HasNameLike() {
		q = q.Where(COLUMN_NAME+" LIKE ?", options.NameLike())
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

	if options.HasTemplateID() {
		q = q.Where(COLUMN_TEMPLATE_ID+" = ?", options.TemplateID())
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
		return q, columns, nil // soft deleted pages requested specifically
	}

	q = q.Where(COLUMN_SOFT_DELETED_AT+" > ?", carbon.Now(carbon.UTC).ToDateTimeString())

	return q, columns, nil
}
