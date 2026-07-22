package cmsstore

import (
	"context"
	"errors"
	"log"
	"strings"

	contractsorm "github.com/dracory/neat/contracts/database/orm"
	"github.com/dromara/carbon/v2"
)

// MenuCount returns the count of menus that match the provided query options.
func (store *storeImplementation) MenuCount(ctx context.Context, options MenuQueryInterface) (int64, error) {
	if store.neatDB == nil {
		return -1, errors.New("cms store: database is nil")
	}

	if !store.menusEnabled {
		return -1, errors.New("menus are disabled")
	}

	if options != nil && !options.IsCountOnly() {
		options.SetCountOnly(true)
	}

	q, _, err := store.menuSelectQuery(options)
	if err != nil {
		return -1, err
	}

	var count int64
	err = q.Table(store.menuTableName).Count(&count)
	return count, err
}

// MenuCreate creates a new menu in the database.
func (store *storeImplementation) MenuCreate(ctx context.Context, menu MenuInterface) error {
	if store.neatDB == nil {
		return errors.New("menustore: database is nil")
	}

	if !store.menusEnabled {
		return errors.New("menus are disabled")
	}

	if menu == nil {
		return errors.New("menu is nil")
	}
	if menu.CreatedAt() == "" {
		menu.SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	}
	if menu.UpdatedAt() == "" {
		menu.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	}

	return store.withTransaction(ctx, func(txCtx context.Context) error {
		data := menu.Data()

		if store.debugEnabled {
			log.Println("MenuCreate:", data)
		}

		err := store.neatDB.Query().Table(store.menuTableName).Create(data)
		if err != nil {
			return err
		}

		menu.MarkAsNotDirty()

		return store.versioningTrackEntity(txCtx, VERSIONING_TYPE_MENU, menu.ID(), menu)
	})
}

// MenuDelete deletes a menu from the database by its ID.
func (store *storeImplementation) MenuDelete(ctx context.Context, menu MenuInterface) error {
	if !store.menusEnabled {
		return errors.New("menus are disabled")
	}

	if menu == nil {
		return errors.New("menu is nil")
	}

	return store.MenuDeleteByID(ctx, menu.ID())
}

// MenuDeleteByID deletes a menu from the database by its ID.
func (store *storeImplementation) MenuDeleteByID(ctx context.Context, id string) error {
	if store.neatDB == nil {
		return errors.New("menustore: database is nil")
	}

	if !store.menusEnabled {
		return errors.New("menus are disabled")
	}

	if id == "" {
		return errors.New("menu id is empty")
	}

	if store.debugEnabled {
		log.Println("MenuDeleteByID:", id)
	}

	_, err := store.neatDB.Query().Table(store.menuTableName).Where("id = ?", id).Delete()

	return err
}

// MenuFindByHandle finds a menu by its handle.
func (store *storeImplementation) MenuFindByHandle(ctx context.Context, handle string) (menu MenuInterface, err error) {
	if !store.menusEnabled {
		return nil, errors.New("menus are disabled")
	}

	if handle == "" {
		return nil, errors.New("menu handle is empty")
	}

	list, err := store.MenuList(ctx, MenuQuery().
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

// MenuFindByID finds a menu by its ID.
func (store *storeImplementation) MenuFindByID(ctx context.Context, id string) (menu MenuInterface, err error) {
	if !store.menusEnabled {
		return nil, errors.New("menus are disabled")
	}

	if id == "" {
		return nil, errors.New("menu id is empty")
	}

	// Normalize ID to lowercase for consistent lookups
	id = NormalizeID(id)

	// Try direct lookup first (handles both 9-char and 32-char IDs)
	list, err := store.MenuList(ctx, MenuQuery().SetID(id).SetLimit(1))
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
			list, err = store.MenuList(ctx, MenuQuery().SetID(unshortenedID).SetLimit(1))
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

// MenuList returns a list of menus that match the provided query options.
func (store *storeImplementation) MenuList(ctx context.Context, query MenuQueryInterface) ([]MenuInterface, error) {
	if !store.menusEnabled {
		return []MenuInterface{}, errors.New("menus are disabled")
	}

	if store.neatDB == nil {
		return []MenuInterface{}, errors.New("menustore: database is nil")
	}

	q, _, err := store.menuSelectQuery(query)
	if err != nil {
		return []MenuInterface{}, err
	}

	type menuRow struct {
		ID            string `db:"id"`
		SiteID        string `db:"site_id"`
		Name          string `db:"name"`
		Handle        string `db:"handle"`
		Status        string `db:"status"`
		Metas         string `db:"metas"`
		Memo          string `db:"memo"`
		CreatedAt     string `db:"created_at"`
		UpdatedAt     string `db:"updated_at"`
		SoftDeletedAt string `db:"soft_deleted_at"`
	}

	var rows []menuRow
	if err := q.Table(store.menuTableName).Get(&rows); err != nil {
		return []MenuInterface{}, err
	}

	list := make([]MenuInterface, 0, len(rows))
	for _, r := range rows {
		modelMap := map[string]string{
			"id":              r.ID,
			"site_id":         r.SiteID,
			"name":            r.Name,
			"handle":          r.Handle,
			"status":          r.Status,
			"metas":           r.Metas,
			"memo":            r.Memo,
			"created_at":      r.CreatedAt,
			"updated_at":      r.UpdatedAt,
			"soft_deleted_at": r.SoftDeletedAt,
		}
		model := NewMenuFromExistingData(modelMap)
		list = append(list, model)
	}

	return list, nil
}

// MenuSoftDelete marks a menu as soft-deleted by setting the soft_deleted_at timestamp.
func (store *storeImplementation) MenuSoftDelete(ctx context.Context, menu MenuInterface) error {
	if !store.menusEnabled {
		return errors.New("menus are disabled")
	}

	if menu == nil {
		return errors.New("menu is nil")
	}

	menu.SetSoftDeletedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))

	return store.MenuUpdate(ctx, menu)
}

// MenuSoftDeleteByID marks a menu as soft-deleted by its ID.
func (store *storeImplementation) MenuSoftDeleteByID(ctx context.Context, id string) error {
	if !store.menusEnabled {
		return errors.New("menus are disabled")
	}

	menu, err := store.MenuFindByID(ctx, id)
	if err != nil {
		return err
	}

	return store.MenuSoftDelete(ctx, menu)
}

// MenuUpdate updates an existing menu in the database.
func (store *storeImplementation) MenuUpdate(ctx context.Context, menu MenuInterface) error {
	if store.neatDB == nil {
		return errors.New("menustore: database is nil")
	}

	if !store.menusEnabled {
		return errors.New("menus are disabled")
	}

	if menu == nil {
		return errors.New("menu is nil")
	}

	menu.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString())

	return store.withTransaction(ctx, func(txCtx context.Context) error {
		dataChanged := menu.DataChanged()
		delete(dataChanged, COLUMN_ID) // ID is not updateable

		if len(dataChanged) < 1 {
			return nil
		}

		if store.debugEnabled {
			log.Println("MenuUpdate:", dataChanged)
		}

		_, err := store.neatDB.Query().Table(store.menuTableName).Where("id = ?", menu.ID()).Update(dataChanged)
		if err != nil {
			return err
		}

		menu.MarkAsNotDirty()

		return store.versioningTrackEntity(txCtx, VERSIONING_TYPE_MENU, menu.ID(), menu)
	})
}

// menuSelectQuery constructs a SQL query for selecting menus based on the provided query options.
func (store *storeImplementation) menuSelectQuery(options MenuQueryInterface) (query contractsorm.Query, columns []any, err error) {
	if options == nil {
		return nil, nil, errors.New("menu query cannot be nil")
	}

	if err := options.Validate(); err != nil {
		return nil, nil, err
	}

	q := store.neatDB.Query().Table(store.menuTableName)

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
		return q, columns, nil // soft deleted menus requested specifically
	}

	q = q.Where(COLUMN_SOFT_DELETED_AT+" > ?", carbon.Now(carbon.UTC).ToDateTimeString())

	return q, columns, nil
}
