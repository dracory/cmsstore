package cmsstore

import (
	"context"
	"errors"
	"log"
	"strconv"
	"strings"

	"github.com/doug-martin/goqu/v9"
	"github.com/dracory/database"
	"github.com/dracory/sb"
	"github.com/dromara/carbon/v2"
	"github.com/samber/lo"
)

// MenuItemCount returns the count of menu items based on the provided query options.
func (store *store) MenuItemCount(ctx context.Context, options MenuItemQueryInterface) (int64, error) {
	// Check if menus are enabled
	if !store.menusEnabled {
		return -1, errors.New("menus are disabled")
	}

	// Set the query to count only
	options.SetCountOnly(true)

	// Generate the select query based on the options
	q, _, err := store.menuItemSelectQuery(options)

	if err != nil {
		return -1, err
	}

	// Prepare the SQL query to count the number of menu items
	sqlStr, params, errSql := q.Prepared(true).
		Limit(1).
		Select(goqu.COUNT(goqu.Star()).As("count")).
		ToSQL()

	if errSql != nil {
		return -1, nil
	}

	// Log the SQL query if debug is enabled
	if store.debugEnabled {
		log.Println(sqlStr)
	}

	// Execute the query and get the result as a map
	mapped, err := database.SelectToMapString(store.toQuerableContext(ctx), sqlStr, params...)

	if err != nil {
		return -1, err
	}

	// Check if the result is empty
	if len(mapped) < 1 {
		return -1, nil
	}

	// Extract the count from the result
	countStr := mapped[0]["count"]

	// Convert the count string to an integer
	i, err := strconv.ParseInt(countStr, 10, 64)

	if err != nil {
		return -1, err
	}

	return i, nil
}

// MenuItemCreate creates a new menu item in the database.
func (store *store) MenuItemCreate(ctx context.Context, menuItem MenuItemInterface) error {
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

	// Get the data from the menu item
	data := menuItem.Data()

	// Prepare the SQL query to insert the menu item
	sqlStr, params, errSql := goqu.Dialect(store.dbDriverName).
		Insert(store.menuItemTableName).
		Prepared(true).
		Rows(data).
		ToSQL()

	if errSql != nil {
		return errSql
	}

	// Log the SQL query if debug is enabled
	if store.debugEnabled {
		log.Println(sqlStr)
	}

	// Check if the database connection is nil
	if store.db == nil {
		return errors.New("menuItemstore: database is nil")
	}

	// Execute the query to insert the menu item
	_, err := database.Execute(store.toQuerableContext(ctx), sqlStr, params...)

	if err != nil {
		return err
	}

	// Mark the menu item as not dirty
	menuItem.MarkAsNotDirty()

	return nil
}

// MenuItemDelete deletes a menu item from the database.
func (store *store) MenuItemDelete(ctx context.Context, menuItem MenuItemInterface) error {
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
func (store *store) MenuItemDeleteByID(ctx context.Context, id string) error {
	// Check if menus are enabled
	if !store.menusEnabled {
		return errors.New("menus are disabled")
	}

	// Validate the ID
	if id == "" {
		return errors.New("menuItem id is empty")
	}

	// Prepare the SQL query to delete the menu item
	sqlStr, params, errSql := goqu.Dialect(store.dbDriverName).
		Delete(store.menuItemTableName).
		Prepared(true).
		Where(goqu.C("id").Eq(id)).
		ToSQL()

	if errSql != nil {
		return errSql
	}

	// Log the SQL query if debug is enabled
	if store.debugEnabled {
		log.Println(sqlStr)
	}

	// Execute the query to delete the menu item
	_, err := database.Execute(store.toQuerableContext(ctx), sqlStr, params...)

	return err
}

// MenuItemFindByID finds a menu item by its ID.
func (store *store) MenuItemFindByID(ctx context.Context, id string) (menuItem MenuItemInterface, err error) {
	// Check if menus are enabled
	if !store.menusEnabled {
		return nil, errors.New("menus are disabled")
	}

	// Validate the ID
	if id == "" {
		return nil, errors.New("menuItem id is empty")
	}

	// List menu items with the specified ID and limit to 1
	list, err := store.MenuItemList(ctx, MenuItemQuery().SetID(id).SetLimit(1))

	if err != nil {
		return nil, err
	}

	// Return the first item if found
	if len(list) > 0 {
		return list[0], nil
	}

	return nil, nil
}

// MenuItemList returns a list of menu items based on the provided query options.
func (store *store) MenuItemList(ctx context.Context, query MenuItemQueryInterface) ([]MenuItemInterface, error) {
	// Check if menus are enabled
	if !store.menusEnabled {
		return []MenuItemInterface{}, errors.New("menus are disabled")
	}

	// Generate the select query based on the options
	q, columns, err := store.menuItemSelectQuery(query)

	if err != nil {
		return []MenuItemInterface{}, err
	}

	// Prepare the SQL query to select the menu items
	sqlStr, _, errSql := q.Select(columns...).ToSQL()

	if errSql != nil {
		return []MenuItemInterface{}, nil
	}

	// Log the SQL query if debug is enabled
	if store.debugEnabled {
		log.Println(sqlStr)
	}

	// Check if the database connection is nil
	if store.db == nil {
		return []MenuItemInterface{}, errors.New("menuItemstore: database is nil")
	}

	// Execute the query and get the result as a map
	modelMaps, err := database.SelectToMapString(store.toQuerableContext(ctx), sqlStr)

	if err != nil {
		return []MenuItemInterface{}, err
	}

	// Convert the map to a list of menu items
	list := []MenuItemInterface{}

	lo.ForEach(modelMaps, func(modelMap map[string]string, index int) {
		model := NewMenuItemFromExistingData(modelMap)
		list = append(list, model)
	})

	return list, nil
}

// MenuItemSoftDelete soft deletes a menu item by setting the soft_deleted_at timestamp.
func (store *store) MenuItemSoftDelete(ctx context.Context, menuItem MenuItemInterface) error {
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
func (store *store) MenuItemSoftDeleteByID(ctx context.Context, id string) error {
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
func (store *store) MenuItemUpdate(ctx context.Context, menuItem MenuItemInterface) error {
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

	// Get the changed data from the menu item
	dataChanged := menuItem.DataChanged()

	// Remove the ID from the changed data as it is not updateable
	delete(dataChanged, COLUMN_ID)

	// Check if there are any changes to update
	if len(dataChanged) < 1 {
		return nil
	}

	// Prepare the SQL query to update the menu item
	sqlStr, params, errSql := goqu.Dialect(store.dbDriverName).
		Update(store.menuItemTableName).
		Prepared(true).
		Set(dataChanged).
		Where(goqu.C(COLUMN_ID).Eq(menuItem.ID())).
		ToSQL()

	if errSql != nil {
		return errSql
	}

	// Log the SQL query if debug is enabled
	if store.debugEnabled {
		log.Println(sqlStr)
	}

	// Check if the database connection is nil
	if store.db == nil {
		return errors.New("menuItemstore: database is nil")
	}

	// Execute the query to update the menu item
	_, err := database.Execute(store.toQuerableContext(ctx), sqlStr, params...)

	// Mark the menu item as not dirty
	menuItem.MarkAsNotDirty()

	return err
}

// menuItemSelectQuery generates a select query based on the provided query options.
func (store *store) menuItemSelectQuery(options MenuItemQueryInterface) (selectDataset *goqu.SelectDataset, columns []any, err error) {
	// Validate the query options
	if options == nil {
		return nil, nil, errors.New("menuItem query cannot be nil")
	}

	if err := options.Validate(); err != nil {
		return nil, nil, err
	}

	// Start building the select query
	q := goqu.Dialect(store.dbDriverName).From(store.menuItemTableName)

	// Apply filters based on the query options
	if options.HasCreatedAtGte() && options.HasCreatedAtLte() {
		q = q.Where(
			goqu.C(COLUMN_CREATED_AT).Gte(options.CreatedAtGte()),
			goqu.C(COLUMN_CREATED_AT).Lte(options.CreatedAtLte()),
		)
	} else if options.HasCreatedAtGte() {
		q = q.Where(goqu.C(COLUMN_CREATED_AT).Gte(options.CreatedAtGte()))
	} else if options.HasCreatedAtLte() {
		q = q.Where(goqu.C(COLUMN_CREATED_AT).Lte(options.CreatedAtLte()))
	}

	if options.HasID() {
		q = q.Where(goqu.C(COLUMN_ID).Eq(options.ID()))
	}

	if options.HasIDIn() {
		q = q.Where(goqu.C(COLUMN_ID).In(options.IDIn()))
	}

	if options.HasMenuID() {
		q = q.Where(goqu.C(COLUMN_MENU_ID).Eq(options.MenuID()))
	}

	if options.HasNameLike() {
		q = q.Where(goqu.C(COLUMN_NAME).Like(options.NameLike()))
	}

	if options.HasSiteID() {
		q = q.Where(goqu.C(COLUMN_SITE_ID).Eq(options.SiteID()))
	}

	if options.HasStatus() {
		q = q.Where(goqu.C(COLUMN_STATUS).Eq(options.Status()))
	}

	if options.HasStatusIn() {
		q = q.Where(goqu.C(COLUMN_STATUS).In(options.StatusIn()))
	}

	// Apply pagination options if not counting only
	if !options.IsCountOnly() {
		if options.HasLimit() {
			q = q.Limit(uint(options.Limit()))
		}

		if options.HasOffset() {
			q = q.Offset(uint(options.Offset()))
		}
	}

	// Apply sorting options
	sortOrder := sb.DESC
	if options.HasSortOrder() {
		sortOrder = options.SortOrder()
	}

	if options.HasOrderBy() {
		if strings.EqualFold(sortOrder, sb.ASC) {
			q = q.Order(goqu.I(options.OrderBy()).Asc())
		} else {
			q = q.Order(goqu.I(options.OrderBy()).Desc())
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
	softDeleted := goqu.C(COLUMN_SOFT_DELETED_AT).
		Gt(carbon.Now(carbon.UTC).ToDateTimeString())

	return q.Where(softDeleted), columns, nil
}
