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

func (store *store) TranslationCount(options TranslationQueryInterface) (int64, error) {
	options.SetCountOnly(true)

	q, _, err := store.translationSelectQuery(options)

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

func (store *store) TranslationCreate(translation TranslationInterface) error {
	if translation == nil {
		return errors.New("translation is nil")
	}
	if translation.CreatedAt() == "" {
		translation.SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	}

	if translation.UpdatedAt() == "" {
		translation.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	}

	data := translation.Data()

	sqlStr, params, errSql := goqu.Dialect(store.dbDriverName).
		Insert(store.translationTableName).
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
		return errors.New("translationstore: database is nil")
	}

	_, err := store.db.Exec(sqlStr, params...)

	if err != nil {
		return err
	}

	translation.MarkAsNotDirty()

	return nil
}

func (store *store) TranslationDelete(translation TranslationInterface) error {
	if translation == nil {
		return errors.New("translation is nil")
	}

	return store.TranslationDeleteByID(translation.ID())
}

func (store *store) TranslationDeleteByID(id string) error {
	if id == "" {
		return errors.New("translation id is empty")
	}

	sqlStr, params, errSql := goqu.Dialect(store.dbDriverName).
		Delete(store.translationTableName).
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

func (store *store) TranslationFindByHandle(hadle string) (translation TranslationInterface, err error) {
	if hadle == "" {
		return nil, errors.New("translation handle is empty")
	}

	list, err := store.TranslationList(TranslationQuery().
		SetHandle(hadle).
		SetLimit(1))

	if err != nil {
		return nil, err
	}

	if len(list) > 0 {
		return list[0], nil
	}

	return nil, nil
}

func (store *store) TranslationFindByID(id string) (translation TranslationInterface, err error) {
	if id == "" {
		return nil, errors.New("translation id is empty")
	}

	list, err := store.TranslationList(TranslationQuery().SetID(id).SetLimit(1))

	if err != nil {
		return nil, err
	}

	if len(list) > 0 {
		return list[0], nil
	}

	return nil, nil
}

func (store *store) TranslationLanguageDefault() string {
	return store.translationLanguageDefault
}

func (store *store) TranslationLanguages() map[string]string {
	return store.translationLanguages
}

func (store *store) TranslationList(query TranslationQueryInterface) ([]TranslationInterface, error) {
	q, columns, err := store.translationSelectQuery(query)

	if err != nil {
		return []TranslationInterface{}, err
	}

	sqlStr, _, errSql := q.Select(columns...).ToSQL()

	if errSql != nil {
		return []TranslationInterface{}, nil
	}

	if store.debugEnabled {
		log.Println(sqlStr)
	}

	if store.db == nil {
		return []TranslationInterface{}, errors.New("translationstore: database is nil")
	}

	db := sb.NewDatabase(store.db, store.dbDriverName)

	if db == nil {
		return []TranslationInterface{}, errors.New("translationstore: database is nil")
	}

	modelMaps, err := db.SelectToMapString(sqlStr)

	if err != nil {
		return []TranslationInterface{}, err
	}

	list := []TranslationInterface{}

	lo.ForEach(modelMaps, func(modelMap map[string]string, index int) {
		model := NewTranslationFromExistingData(modelMap)
		list = append(list, model)
	})

	return list, nil
}

func (store *store) TranslationSoftDelete(translation TranslationInterface) error {
	if translation == nil {
		return errors.New("translation is nil")
	}

	translation.SetSoftDeletedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))

	return store.TranslationUpdate(translation)
}

func (store *store) TranslationSoftDeleteByID(id string) error {
	translation, err := store.TranslationFindByID(id)

	if err != nil {
		return err
	}

	return store.TranslationSoftDelete(translation)
}

func (store *store) TranslationUpdate(translation TranslationInterface) error {
	if translation == nil {
		return errors.New("translation is nil")
	}

	translation.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString())

	dataChanged := translation.DataChanged()

	delete(dataChanged, COLUMN_ID) // ID is not updateable

	if len(dataChanged) < 1 {
		return nil
	}

	sqlStr, params, errSql := goqu.Dialect(store.dbDriverName).
		Update(store.translationTableName).
		Prepared(true).
		Set(dataChanged).
		Where(goqu.C(COLUMN_ID).Eq(translation.ID())).
		ToSQL()

	if errSql != nil {
		return errSql
	}

	if store.debugEnabled {
		log.Println(sqlStr)
	}

	if store.db == nil {
		return errors.New("translationstore: database is nil")
	}

	_, err := store.db.Exec(sqlStr, params...)

	translation.MarkAsNotDirty()

	return err
}

func (store *store) translationSelectQuery(options TranslationQueryInterface) (selectDataset *goqu.SelectDataset, columns []any, err error) {
	if options == nil {
		return nil, nil, errors.New("translation query cannot be nil")
	}

	if err := options.Validate(); err != nil {
		return nil, nil, err
	}

	q := goqu.Dialect(store.dbDriverName).From(store.translationTableName)

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
		return q, columns, nil // soft deleted translations requested specifically
	}

	softDeleted := goqu.C(COLUMN_SOFT_DELETED_AT).
		Gt(carbon.Now(carbon.UTC).ToDateTimeString())

	return q.Where(softDeleted), columns, nil
}
