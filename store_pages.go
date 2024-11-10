package cmsstore

import (
	"errors"
	"log"
	"strconv"
	"strings"

	"github.com/doug-martin/goqu/v9"
	"github.com/golang-module/carbon/v2"
	"github.com/gouniverse/sb"
	"github.com/samber/lo"
)

func (store *store) PageCount(options PageQueryInterface) (int64, error) {
	options.SetCountOnly(true)

	q := store.pageSelectQuery(options)

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

func (store *store) PageCreate(page PageInterface) error {
	page.SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	page.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))

	data := page.Data()

	sqlStr, params, errSql := goqu.Dialect(store.dbDriverName).
		Insert(store.pageTableName).
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
		return errors.New("pagestore: database is nil")
	}

	_, err := store.db.Exec(sqlStr, params...)

	if err != nil {
		return err
	}

	page.MarkAsNotDirty()

	return nil
}

func (store *store) PageDelete(page PageInterface) error {
	if page == nil {
		return errors.New("page is nil")
	}

	return store.PageDeleteByID(page.ID())
}

func (store *store) PageDeleteByID(id string) error {
	if id == "" {
		return errors.New("page id is empty")
	}

	sqlStr, params, errSql := goqu.Dialect(store.dbDriverName).
		Delete(store.pageTableName).
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

func (store *store) PageFindByHandle(hadle string) (page PageInterface, err error) {
	if hadle == "" {
		return nil, errors.New("page handle is empty")
	}

	query := NewPageQuery()

	query, err = query.SetHandle(hadle)

	if err != nil {
		return nil, err
	}

	query, err = query.SetLimit(1)

	if err != nil {
		return nil, err
	}

	list, err := store.PageList(query)

	if err != nil {
		return nil, err
	}

	if len(list) > 0 {
		return list[0], nil
	}

	return nil, nil
}

func (store *store) PageFindByID(id string) (page PageInterface, err error) {
	if id == "" {
		return nil, errors.New("page id is empty")
	}

	query := NewPageQuery()

	query, err = query.SetID(id)

	if err != nil {
		return nil, err
	}

	query, err = query.SetLimit(1)

	if err != nil {
		return nil, err
	}

	list, err := store.PageList(query)

	if err != nil {
		return nil, err
	}

	if len(list) > 0 {
		return list[0], nil
	}

	return nil, nil
}

func (store *store) PageList(query PageQueryInterface) ([]PageInterface, error) {
	q := store.pageSelectQuery(query)

	sqlStr, _, errSql := q.Select().ToSQL()

	if errSql != nil {
		return []PageInterface{}, nil
	}

	if store.debugEnabled {
		log.Println(sqlStr)
	}

	if store.db == nil {
		return []PageInterface{}, errors.New("pagestore: database is nil")
	}

	db := sb.NewDatabase(store.db, store.dbDriverName)

	if db == nil {
		return []PageInterface{}, errors.New("pagestore: database is nil")
	}

	modelMaps, err := db.SelectToMapString(sqlStr)

	if err != nil {
		return []PageInterface{}, err
	}

	list := []PageInterface{}

	lo.ForEach(modelMaps, func(modelMap map[string]string, index int) {
		model := NewPageFromExistingData(modelMap)
		list = append(list, model)
	})

	return list, nil
}

func (store *store) PageSoftDelete(page PageInterface) error {
	if page == nil {
		return errors.New("page is nil")
	}

	page.SetSoftDeletedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))

	return store.PageUpdate(page)
}

func (store *store) PageSoftDeleteByID(id string) error {
	page, err := store.PageFindByID(id)

	if err != nil {
		return err
	}

	return store.PageSoftDelete(page)
}

func (store *store) PageUpdate(page PageInterface) error {
	if page == nil {
		return errors.New("page is nil")
	}

	page.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString())

	dataChanged := page.DataChanged()

	delete(dataChanged, COLUMN_ID) // ID is not updateable

	if len(dataChanged) < 1 {
		return nil
	}

	sqlStr, params, errSql := goqu.Dialect(store.dbDriverName).
		Update(store.pageTableName).
		Prepared(true).
		Set(dataChanged).
		Where(goqu.C(COLUMN_ID).Eq(page.ID())).
		ToSQL()

	if errSql != nil {
		return errSql
	}

	if store.debugEnabled {
		log.Println(sqlStr)
	}

	if store.db == nil {
		return errors.New("pagestore: database is nil")
	}

	_, err := store.db.Exec(sqlStr, params...)

	page.MarkAsNotDirty()

	return err
}

func (store *store) pageSelectQuery(options PageQueryInterface) *goqu.SelectDataset {
	q := goqu.Dialect(store.dbDriverName).From(store.pageTableName)

	if options.ID() != "" {
		q = q.Where(goqu.C(COLUMN_ID).Eq(options.ID()))
	}

	if len(options.IDIn()) > 0 {
		q = q.Where(goqu.C(COLUMN_ID).In(options.IDIn()))
	}

	if options.Handle() != "" {
		q = q.Where(goqu.C(COLUMN_HANDLE).Eq(options.Handle()))
	}

	if options.NameLike() != "" {
		q = q.Where(goqu.C(COLUMN_NAME).ILike(options.NameLike()))
	}

	if options.Status() != "" {
		q = q.Where(goqu.C(COLUMN_STATUS).Eq(options.Status()))
	}

	if len(options.StatusIn()) > 0 {
		q = q.Where(goqu.C(COLUMN_STATUS).In(options.StatusIn()))
	}

	if options.CreatedAtGte() != "" && options.CreatedAtLte() != "" {
		q = q.Where(
			goqu.C(COLUMN_CREATED_AT).Gte(options.CreatedAtGte()),
			goqu.C(COLUMN_CREATED_AT).Lte(options.CreatedAtLte()),
		)
	} else if options.CreatedAtGte() != "" {
		q = q.Where(goqu.C(COLUMN_CREATED_AT).Gte(options.CreatedAtGte()))
	} else if options.CreatedAtLte() != "" {
		q = q.Where(goqu.C(COLUMN_CREATED_AT).Lte(options.CreatedAtLte()))
	}

	if !options.CountOnly() {
		if options.Limit() > 0 {
			q = q.Limit(uint(options.Limit()))
		}

		if options.Offset() > 0 {
			q = q.Offset(uint(options.Offset()))
		}
	}

	sortOrder := sb.DESC
	if options.SortOrder() != "" {
		sortOrder = options.SortOrder()
	}

	if options.OrderBy() != "" {
		if strings.EqualFold(sortOrder, sb.ASC) {
			q = q.Order(goqu.I(options.OrderBy()).Asc())
		} else {
			q = q.Order(goqu.I(options.OrderBy()).Desc())
		}
	}

	if options.WithSoftDeleted() {
		return q // soft deleted pages requested specifically
	}

	softDeleted := goqu.C(COLUMN_SOFT_DELETED_AT).
		Gt(carbon.Now(carbon.UTC).ToDateTimeString())

	return q.Where(softDeleted)
}
