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

func (store *store) BlockCount(options BlockQueryInterface) (int64, error) {
	options.SetCountOnly(true)

	q, err := store.blockSelectQuery(options)

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

func (store *store) BlockCreate(block BlockInterface) error {
	block.SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	block.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))

	data := block.Data()

	sqlStr, params, errSql := goqu.Dialect(store.dbDriverName).
		Insert(store.blockTableName).
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
		return errors.New("blockstore: database is nil")
	}

	_, err := store.db.Exec(sqlStr, params...)

	if err != nil {
		return err
	}

	block.MarkAsNotDirty()

	return nil
}

func (store *store) BlockDelete(block BlockInterface) error {
	if block == nil {
		return errors.New("block is nil")
	}

	return store.BlockDeleteByID(block.ID())
}

func (store *store) BlockDeleteByID(id string) error {
	if id == "" {
		return errors.New("block id is empty")
	}

	sqlStr, params, errSql := goqu.Dialect(store.dbDriverName).
		Delete(store.blockTableName).
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

func (store *store) BlockFindByHandle(hadle string) (block BlockInterface, err error) {
	if hadle == "" {
		return nil, errors.New("block handle is empty")
	}

	list, err := store.BlockList(BlockQuery().SetHandle(hadle).SetLimit(1))

	if err != nil {
		return nil, err
	}

	if len(list) > 0 {
		return list[0], nil
	}

	return nil, nil
}

func (store *store) BlockFindByID(id string) (block BlockInterface, err error) {
	if id == "" {
		return nil, errors.New("block id is empty")
	}

	list, err := store.BlockList(BlockQuery().SetID(id).SetLimit(1))

	if err != nil {
		return nil, err
	}

	if len(list) > 0 {
		return list[0], nil
	}

	return nil, nil
}

func (store *store) BlockList(query BlockQueryInterface) ([]BlockInterface, error) {
	q, err := store.blockSelectQuery(query)

	if err != nil {
		return []BlockInterface{}, err
	}

	sqlStr, _, errSql := q.Select().ToSQL()

	if errSql != nil {
		return []BlockInterface{}, nil
	}

	if store.debugEnabled {
		log.Println(sqlStr)
	}

	if store.db == nil {
		return []BlockInterface{}, errors.New("blockstore: database is nil")
	}

	db := sb.NewDatabase(store.db, store.dbDriverName)

	if db == nil {
		return []BlockInterface{}, errors.New("blockstore: database is nil")
	}

	modelMaps, err := db.SelectToMapString(sqlStr)

	if err != nil {
		return []BlockInterface{}, err
	}

	list := []BlockInterface{}

	lo.ForEach(modelMaps, func(modelMap map[string]string, index int) {
		model := NewBlockFromExistingData(modelMap)
		list = append(list, model)
	})

	return list, nil
}

func (store *store) BlockSoftDelete(block BlockInterface) error {
	if block == nil {
		return errors.New("block is nil")
	}

	block.SetSoftDeletedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))

	return store.BlockUpdate(block)
}

func (store *store) BlockSoftDeleteByID(id string) error {
	block, err := store.BlockFindByID(id)

	if err != nil {
		return err
	}

	return store.BlockSoftDelete(block)
}

func (store *store) BlockUpdate(block BlockInterface) error {
	if block == nil {
		return errors.New("block is nil")
	}

	block.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString())

	dataChanged := block.DataChanged()

	delete(dataChanged, COLUMN_ID) // ID is not updateable

	if len(dataChanged) < 1 {
		return nil
	}

	sqlStr, params, errSql := goqu.Dialect(store.dbDriverName).
		Update(store.blockTableName).
		Prepared(true).
		Set(dataChanged).
		Where(goqu.C(COLUMN_ID).Eq(block.ID())).
		ToSQL()

	if errSql != nil {
		return errSql
	}

	if store.debugEnabled {
		log.Println(sqlStr)
	}

	if store.db == nil {
		return errors.New("blockstore: database is nil")
	}

	_, err := store.db.Exec(sqlStr, params...)

	block.MarkAsNotDirty()

	return err
}

func (store *store) blockSelectQuery(options BlockQueryInterface) (*goqu.SelectDataset, error) {
	if options == nil {
		return nil, errors.New("block query: cannot be nil")
	}

	if err := options.Validate(); err != nil {
		return nil, err
	}

	q := goqu.Dialect(store.dbDriverName).From(store.blockTableName)

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

	if options.HasHandle() {
		q = q.Where(goqu.C(COLUMN_HANDLE).Eq(options.Handle()))
	}

	if options.HasID() {
		q = q.Where(goqu.C(COLUMN_ID).Eq(options.ID()))
	}

	if options.HasIDIn() {
		q = q.Where(goqu.C(COLUMN_ID).In(options.IDIn()))
	}

	if options.HasNameLike() {
		q = q.Where(goqu.C(COLUMN_NAME).Like(`%` + options.NameLike() + `%`))
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
	if options.HasSortOrder() && options.SortOrder() != "" {
		sortOrder = options.SortOrder()
	}

	if options.HasOrderBy() && options.OrderBy() != "" {
		if strings.EqualFold(sortOrder, sb.ASC) {
			q = q.Order(goqu.I(options.OrderBy()).Asc())
		} else {
			q = q.Order(goqu.I(options.OrderBy()).Desc())
		}
	}

	if options.SoftDeleteIncluded() {
		return q, nil // soft deleted blocks requested specifically
	}

	softDeleted := goqu.C(COLUMN_SOFT_DELETED_AT).
		Gt(carbon.Now(carbon.UTC).ToDateTimeString())

	return q.Where(softDeleted), nil
}
