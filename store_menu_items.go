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

// MenuItemCount returns the count of menu items based on the provided query options.
func (store *storeImplementation) MenuItemCount(ctx context.Context, options MenuItemQueryInterface) (int64, error) {
	if store.neatDB == nil {
		return -1, errors.New("cms store: database is nil")
	}

	// Check if menus are enabled
	if !store.menusEnabled {
		return -1, errors.New("menus are disabled")
	}

	if options != nil && !options.IsCountOnly() {
		options.SetCountOnly(true)
	}

	// Generate the select query based on the options
	q, _, err := store.menuItemSelectQuery(options)

	if err != nil {
		return -1, err
	}

	var count int64
	err = q.Table(store.menuItemTableName).Count(&count)
	return count, err
}

// MenuItemCreate creates a new menu item in the database.
func (store *storeImplementation) MenuItemCreate(ctx context.Context, menuItem MenuItemInterface) error {
	if store.neatDB == nil {
		return errors.New("menuitemstore: database is nil")
	}

	// Check if menus are enabled
	if !store.menusEnabled {
		return errors.New("menus are disabled")
	}

	// Validate the menu item
	if menuItem == nil {
		return errors.New("menuItem is nil")
	}

	// Set the creation timestamp if not already set
	if menuItem.CreatedAt() == "" {
		menuItem.SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	}

	// Set the update timestamp if not already set
	if menuItem.UpdatedAt() == "" {
		menuItem.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	}

	return store.withTransaction(ctx, func(txCtx context.Context) error {
		// Get the data from the menu item
		data := menuItem.Data()

		// Log the SQL query if debug is enabled
		if store.debugEnabled {
			log.Println("MenuItemCreate:", data)
		}

		err := store.neatDB.Query().Table(store.menuItemTableName).Create(data)
		if err != nil {
			return err
		}

		// Mark the menu item as not dirty
		menuItem.MarkAsNotDirty()

		return store.versioningTrackEntity(txCtx, VERSIONING_TYPE_MENU_ITEM, menuItem.ID(), menuItem)
	})
}

// MenuItemDelete deletes a menu item from the database.
func (store *storeImplementation) MenuItemDelete(ctx context.Context, menuItem MenuItemInterface) error {
	// Check if menus are enabled
	if !store.menusEnabled {
		return errors.New("menus are disabled")
	}

	// Validate the menu item
	if menuItem == nil {
		return errors.New("menuItem is nil")
	}

	// Delete the menu item by its ID
	return store.MenuItemDeleteByID(ctx, menuItem.ID())
}

// MenuItemDeleteByID deletes a menu item from the database by its ID.
func (store *storeImplementation) MenuItemDeleteByID(ctx context.Context, id string) error {
	if store.neatDB == nil {
		return errors.New("menuItemstore: database is nil")
	}

	// Check if menus are enabled
	if !store.menusEnabled {
		return errors.New("menus are disabled")
	}

	// Validate the ID
	if id == "" {
		return errors.New("menuItem id is empty")
	}

	// Log the SQL query if debug is enabled
	if store.debugEnabled {
		log.Println("MenuItemDeleteByID:", id)
	}

	_, err := store.neatDB.Query().Table(store.menuItemTableName).Where("id = ?", id).Delete()

	return err
}

// MenuItemFindByID finds a menu item by its ID.
func (store *storeImplementation) MenuItemFindByID(ctx context.Context, id string) (menuItem MenuItemInterface, err error) {
	// Check if menus are enabled
	if !store.menusEnabled {
		return nil, errors.New("menus are disabled")
	}

	// Validate the ID
	if id == "" {
		return nil, errors.New("menuItem id is empty")
	}

	// Normalize ID to lowercase for consistent lookups
	id = NormalizeID(id)

	// Try direct lookup first (handles both 9-char and 32-char IDs)
	list, err := store.MenuItemList(ctx, MenuItemQuery().SetID(id).SetLimit(1))

	if err != nil {
		return nil, err
	}

	// Return the first item if found
	if len(list) > 0 {
		return list[0], nil
	}

	// If not found and ID looks shortened, try unshortening
	if IsShortID(id) {
		unshortenedID := UnshortenID(id)
		if unshortenedID != id {
			list, err = store.MenuItemList(ctx, MenuItemQuery().SetID(unshortenedID).SetLimit(1))
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

// MenuItemList returns a list of menu items based on the provided query options.
func (store *storeImplementation) MenuItemList(ctx context.Context, query MenuItemQueryInterface) ([]MenuItemInterface, error) {
	// Check if menus are enabled
	if !store.menusEnabled {
		return []MenuItemInterface{}, errors.New("menus are disabled")
	}

	if store.neatDB == nil {
		return []MenuItemInterface{}, errors.New("menuItemstore: database is nil")
	}

	// Generate the select query based on the options
	q, _, err := store.menuItemSelectQuery(query)

	if err != nil {
		return []MenuItemInterface{}, err
	}

	type menuItemRow struct {
		ID            string `db:"id"`
		SiteID        string `db:"site_id"`
		MenuID        string `db:"menu_id"`
		PageID        string `db:"page_id"`
		ParentID      string `db:"parent_id"`
		Name          string `db:"name"`
		Handle        string `db:"handle"`
		URL           string `db:"url"`
		Target        string `db:"target"`
		Sequence      int    `db:"sequence"`
		Status        string `db:"status"`
		Metas         string `db:"metas"`
		Memo          string `db:"memo"`
		CreatedAt     string `db:"created_at"`
		UpdatedAt     string `db:"updated_at"`
		SoftDeletedAt string `db:"soft_deleted_at"`
	}

	var rows []menuItemRow
	if err := q.Table(store.menuItemTableName).Get(&rows); err != nil {
		return []MenuItemInterface{}, err
	}

	list := make([]MenuItemInterface, 0, len(rows))
	for _, r := range rows {
		modelMap := map[string]string{
			"id":              r.ID,
			"site_id":         r.SiteID,
			"menu_id":         r.MenuID,
			"page_id":         r.PageID,
			"parent_id":       r.ParentID,
			"name":            r.Name,
			"handle":          r.Handle,
			"url":             r.URL,
			"target":          r.Target,
			"sequence":        strconv.Itoa(r.Sequence),
			"status":          r.Status,
			"metas":           r.Metas,
			"memo":            r.Memo,
			"created_at":      r.CreatedAt,
			"updated_at":      r.UpdatedAt,
			"soft_deleted_at": r.SoftDeletedAt,
		}
		model := NewMenuItemFromExistingData(modelMap)
		list = append(list, model)
	}

	return list, nil
}

// MenuItemSoftDelete soft deletes a menu item by setting the soft_deleted_at timestamp.
func (store *storeImplementation) MenuItemSoftDelete(ctx context.Context, menuItem MenuItemInterface) error {
	// Check if menus are enabled
	if !store.menusEnabled {
		return errors.New("menus are disabled")
	}

	// Validate the menu item
	if menuItem == nil {
		return errors.New("menuItem is nil")
	}

	// Set the soft deleted timestamp
	menuItem.SetSoftDeletedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))

	// Update the menu item
	return store.MenuItemUpdate(ctx, menuItem)
}

// MenuItemSoftDeleteByID soft deletes a menu item by its ID.
func (store *storeImplementation) MenuItemSoftDeleteByID(ctx context.Context, id string) error {
	// Check if menus are enabled
	if !store.menusEnabled {
		return errors.New("menus are disabled")
	}

	// Find the menu item by ID
	menuItem, err := store.MenuItemFindByID(ctx, id)

	if err != nil {
		return err
	}

	// Soft delete the menu item
	return store.MenuItemSoftDelete(ctx, menuItem)
}

// MenuItemUpdate updates an existing menu item in the database.
func (store *storeImplementation) MenuItemUpdate(ctx context.Context, menuItem MenuItemInterface) error {
	if store.neatDB == nil {
		return errors.New("menuitemstore: database is nil")
	}

	// Check if menus are enabled
	if !store.menusEnabled {
		return errors.New("menus are disabled")
	}

	// Validate the menu item
	if menuItem == nil {
		return errors.New("menuItem is nil")
	}

	// Set the update timestamp
	menuItem.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString())

	return store.withTransaction(ctx, func(txCtx context.Context) error {
		// Get the changed data from the menu item
		dataChanged := menuItem.DataChanged()

		// Remove the ID from the changed data as it is not updateable
		delete(dataChanged, COLUMN_ID)

		// Check if there are any changes to update
		if len(dataChanged) < 1 {
			return nil
		}

		// Log the SQL query if debug is enabled
		if store.debugEnabled {
			log.Println("MenuItemUpdate:", dataChanged)
		}

		_, err := store.neatDB.Query().Table(store.menuItemTableName).Where("id = ?", menuItem.ID()).Update(dataChanged)
		if err != nil {
			return err
		}

		// Mark the menu item as not dirty
		menuItem.MarkAsNotDirty()

		return store.versioningTrackEntity(txCtx, VERSIONING_TYPE_MENU_ITEM, menuItem.ID(), menuItem)
	})
}

// menuItemSelectQuery generates a select query based on the provided query options.
func (store *storeImplementation) menuItemSelectQuery(options MenuItemQueryInterface) (query contractsorm.Query, columns []any, err error) {
	// Validate the query options
	if options == nil {
		return nil, nil, errors.New("menuItem query cannot be nil")
	}

	if err := options.Validate(); err != nil {
		return nil, nil, err
	}

	// Start building the select query
	q := store.neatDB.Query().Table(store.menuItemTableName)

	// Apply filters based on the query options
	if options.HasCreatedAtGte() && options.HasCreatedAtLte() {
		q = q.Where(COLUMN_CREATED_AT+" >= ? AND "+COLUMN_CREATED_AT+" <= ?", options.CreatedAtGte(), options.CreatedAtLte())
	} else if options.HasCreatedAtGte() {
		q = q.Where(COLUMN_CREATED_AT+" >= ?", options.CreatedAtGte())
	} else if options.HasCreatedAtLte() {
		q = q.Where(COLUMN_CREATED_AT+" <= ?", options.CreatedAtLte())
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

	if options.HasMenuID() {
		q = q.Where(COLUMN_MENU_ID+" = ?", options.MenuID())
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

	// Apply pagination options if not counting only
	if !options.IsCountOnly() {
		if options.HasLimit() {
			q = q.Limit(options.Limit())
		}

		if options.HasOffset() {
			q = q.Offset(options.Offset())
		}
	}

	// Apply sorting options
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

	// Collect the columns to select
	columns = []any{}

	for _, column := range options.Columns() {
		columns = append(columns, column)
	}

	// Include soft deleted items if requested
	if options.SoftDeletedIncluded() {
		return q, columns, nil
	}

	// Exclude soft deleted items by default
	q = q.Where(COLUMN_SOFT_DELETED_AT+" > ?", carbon.Now(carbon.UTC).ToDateTimeString())

	return q, columns, nil
}
