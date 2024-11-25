package cmsstore

import (
	"errors"
	"log"
	"strconv"
	"strings"

	"github.com/doug-martin/goqu/v9"
	"github.com/dromara/carbon/v2"
	"github.com/gouniverse/sb"
	"github.com/samber/lo"
)

func (store *store) MenuItemCount(options MenuItemQueryInterface) (int64, error) {
	if !store.menusEnabled {
		return -1, errors.New("menus are disabled")
	}

	options.SetCountOnly(true)

	q, _, err := store.menuItemSelectQuery(options)

	if err != nil {
		return -1, err
	}

	sqlStr, params, errSql := q.Prepared(true).
		Limit(1).
		Select(goqu.COUNT(goqu.Star()).As("count")).
		ToSQL()

	if errSql != nil {
		return -1, nil
	}

	if store.debugEnabled {
		log.Println(sqlStr)
	}

	db := sb.NewDatabase(store.db, store.dbDriverName)
	mapped, err := db.SelectToMapString(sqlStr, params...)
	if err != nil {
		return -1, err
	}

	if len(mapped) < 1 {
		return -1, nil
	}

	countStr := mapped[0]["count"]

	i, err := strconv.ParseInt(countStr, 10, 64)

	if err != nil {
		return -1, err

	}

	return i, nil
}

func (store *store) MenuItemCreate(menuItem MenuItemInterface) error {
	if !store.menusEnabled {
		return errors.New("menus are disabled")
	}

	if menuItem == nil {
		return errors.New("menuItem is nil")
	}
	if menuItem.CreatedAt() == "" {
		menuItem.SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	}

	if menuItem.UpdatedAt() == "" {
		menuItem.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	}

	data := menuItem.Data()

	sqlStr, params, errSql := goqu.Dialect(store.dbDriverName).
		Insert(store.menuItemTableName).
		Prepared(true).
		Rows(data).
		ToSQL()

	if errSql != nil {
		return errSql
	}

	if store.debugEnabled {
		log.Println(sqlStr)
	}

	if store.db == nil {
		return errors.New("menuItemstore: database is nil")
	}

	_, err := store.db.Exec(sqlStr, params...)

	if err != nil {
		return err
	}

	menuItem.MarkAsNotDirty()

	return nil
}

func (store *store) MenuItemDelete(menuItem MenuItemInterface) error {
	if !store.menusEnabled {
		return errors.New("menus are disabled")
	}

	if menuItem == nil {
		return errors.New("menuItem is nil")
	}

	return store.MenuItemDeleteByID(menuItem.ID())
}

func (store *store) MenuItemDeleteByID(id string) error {
	if !store.menusEnabled {
		return errors.New("menus are disabled")
	}

	if id == "" {
		return errors.New("menuItem id is empty")
	}

	sqlStr, params, errSql := goqu.Dialect(store.dbDriverName).
		Delete(store.menuItemTableName).
		Prepared(true).
		Where(goqu.C("id").Eq(id)).
		ToSQL()

	if errSql != nil {
		return errSql
	}

	if store.debugEnabled {
		log.Println(sqlStr)
	}

	_, err := store.db.Exec(sqlStr, params...)

	return err
}

func (store *store) MenuItemFindByID(id string) (menuItem MenuItemInterface, err error) {
	if !store.menusEnabled {
		return nil, errors.New("menus are disabled")
	}

	if id == "" {
		return nil, errors.New("menuItem id is empty")
	}

	list, err := store.MenuItemList(MenuItemQuery().SetID(id).SetLimit(1))

	if err != nil {
		return nil, err
	}

	if len(list) > 0 {
		return list[0], nil
	}

	return nil, nil
}

func (store *store) MenuItemList(query MenuItemQueryInterface) ([]MenuItemInterface, error) {
	if !store.menusEnabled {
		return []MenuItemInterface{}, errors.New("menus are disabled")
	}

	q, columns, err := store.menuItemSelectQuery(query)

	if err != nil {
		return []MenuItemInterface{}, err
	}

	sqlStr, _, errSql := q.Select(columns...).ToSQL()

	if errSql != nil {
		return []MenuItemInterface{}, nil
	}

	if store.debugEnabled {
		log.Println(sqlStr)
	}

	if store.db == nil {
		return []MenuItemInterface{}, errors.New("menuItemstore: database is nil")
	}

	db := sb.NewDatabase(store.db, store.dbDriverName)

	if db == nil {
		return []MenuItemInterface{}, errors.New("menuItemstore: database is nil")
	}

	modelMaps, err := db.SelectToMapString(sqlStr)

	if err != nil {
		return []MenuItemInterface{}, err
	}

	list := []MenuItemInterface{}

	lo.ForEach(modelMaps, func(modelMap map[string]string, index int) {
		model := NewMenuItemFromExistingData(modelMap)
		list = append(list, model)
	})

	return list, nil
}

func (store *store) MenuItemSoftDelete(menuItem MenuItemInterface) error {
	if !store.menusEnabled {
		return errors.New("menus are disabled")
	}

	if menuItem == nil {
		return errors.New("menuItem is nil")
	}

	menuItem.SetSoftDeletedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))

	return store.MenuItemUpdate(menuItem)
}

func (store *store) MenuItemSoftDeleteByID(id string) error {
	if !store.menusEnabled {
		return errors.New("menus are disabled")
	}

	menuItem, err := store.MenuItemFindByID(id)

	if err != nil {
		return err
	}

	return store.MenuItemSoftDelete(menuItem)
}

func (store *store) MenuItemUpdate(menuItem MenuItemInterface) error {
	if !store.menusEnabled {
		return errors.New("menus are disabled")
	}

	if menuItem == nil {
		return errors.New("menuItem is nil")
	}

	menuItem.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString())

	dataChanged := menuItem.DataChanged()

	delete(dataChanged, COLUMN_ID) // ID is not updateable

	if len(dataChanged) < 1 {
		return nil
	}

	sqlStr, params, errSql := goqu.Dialect(store.dbDriverName).
		Update(store.menuItemTableName).
		Prepared(true).
		Set(dataChanged).
		Where(goqu.C(COLUMN_ID).Eq(menuItem.ID())).
		ToSQL()

	if errSql != nil {
		return errSql
	}

	if store.debugEnabled {
		log.Println(sqlStr)
	}

	if store.db == nil {
		return errors.New("menuItemstore: database is nil")
	}

	_, err := store.db.Exec(sqlStr, params...)

	menuItem.MarkAsNotDirty()

	return err
}

func (store *store) menuItemSelectQuery(options MenuItemQueryInterface) (selectDataset *goqu.SelectDataset, columns []any, err error) {
	if options == nil {
		return nil, nil, errors.New("menuItem query cannot be nil")
	}

	if err := options.Validate(); err != nil {
		return nil, nil, err
	}

	q := goqu.Dialect(store.dbDriverName).From(store.menuItemTableName)

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

	if !options.IsCountOnly() {
		if options.HasLimit() {
			q = q.Limit(uint(options.Limit()))
		}

		if options.HasOffset() {
			q = q.Offset(uint(options.Offset()))
		}
	}

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

	columns = []any{}

	for _, column := range options.Columns() {
		columns = append(columns, column)
	}

	if options.SoftDeletedIncluded() {
		return q, columns, nil // soft deleted menuItems requested specifically
	}

	softDeleted := goqu.C(COLUMN_SOFT_DELETED_AT).
		Gt(carbon.Now(carbon.UTC).ToDateTimeString())

	return q.Where(softDeleted), columns, nil
}
