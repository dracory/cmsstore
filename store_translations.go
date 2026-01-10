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

func (store *storeImplementation) TranslationCount(ctx context.Context, options TranslationQueryInterface) (int64, error) {
	if store.db == nil {
		return -1, errors.New("cms store: db is nil")
	}

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

	mapped, err := database.SelectToMapString(store.toQuerableContext(ctx), sqlStr, params...)

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

func (store *storeImplementation) TranslationCreate(ctx context.Context, translation TranslationInterface) error {
	if store.db == nil {
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

	_, err := database.Execute(store.toQuerableContext(ctx), sqlStr, params...)

	if err != nil {
		return err
	}

	translation.MarkAsNotDirty()

	if err := store.versioningTrackEntity(ctx, VERSIONING_TYPE_TRANSLATION, translation.ID(), translation); err != nil {
		return err
	}

	return nil
}

func (store *storeImplementation) TranslationDelete(ctx context.Context, translation TranslationInterface) error {
	if store.db == nil {
		return errors.New("cmsstore: database is nil")
	}

	if translation == nil {
		return errors.New("translation is nil")
	}

	return store.TranslationDeleteByID(ctx, translation.ID())
}

func (store *storeImplementation) TranslationDeleteByID(ctx context.Context, id string) error {
	if store.db == nil {
		return errors.New("cmsstore: database is nil")
	}

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

	_, err := database.Execute(store.toQuerableContext(ctx), sqlStr, params...)

	return err
}

func (store *storeImplementation) TranslationFindByHandle(ctx context.Context, handle string) (translation TranslationInterface, err error) {
	if store.db == nil {
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
	if store.db == nil {
		return nil, errors.New("cmsstore: database is nil")
	}

	if id == "" {
		return nil, errors.New("translation id is empty")
	}

	list, err := store.TranslationList(ctx, TranslationQuery().SetID(id).SetLimit(1))

	if err != nil {
		return nil, err
	}

	if len(list) > 0 {
		return list[0], nil
	}

	return nil, nil
}

func (store *storeImplementation) TranslationFindByHandleOrID(ctx context.Context, handleOrID string, language string) (translation TranslationInterface, err error) {
	if store.db == nil {
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
	if store.db == nil {
		return []TranslationInterface{}, errors.New("cmsstore: database is nil")
	}

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

	modelMaps, err := database.SelectToMapString(store.toQuerableContext(ctx), sqlStr)

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

func (store *storeImplementation) TranslationSoftDelete(ctx context.Context, translation TranslationInterface) error {
	if store.db == nil {
		return errors.New("cmsstore: database is nil")
	}

	if translation == nil {
		return errors.New("translation is nil")
	}

	translation.SetSoftDeletedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))

	return store.TranslationUpdate(ctx, translation)
}

func (store *storeImplementation) TranslationSoftDeleteByID(ctx context.Context, id string) error {
	if store.db == nil {
		return errors.New("cmsstore: database is nil")
	}

	if id == "" {
		return errors.New("translation id is empty")
	}

	translation, err := store.TranslationFindByID(ctx, id)

	if err != nil {
		return err
	}

	return store.TranslationSoftDelete(ctx, translation)
}

func (store *storeImplementation) TranslationUpdate(ctx context.Context, translation TranslationInterface) error {
	if store.db == nil {
		return errors.New("cmsstore: database is nil")
	}

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

	_, err := database.Execute(store.toQuerableContext(ctx), sqlStr, params...)

	if err != nil {
		return err
	}

	translation.MarkAsNotDirty()

	if err := store.versioningTrackEntity(ctx, VERSIONING_TYPE_TRANSLATION, translation.ID(), translation); err != nil {
		return err
	}

	return nil
}

func (store *storeImplementation) translationSelectQuery(options TranslationQueryInterface) (selectDataset *goqu.SelectDataset, columns []any, err error) {
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

	if options.HasHandleOrID() {
		q = q.Where(
			goqu.C(COLUMN_HANDLE).Eq(options.HandleOrID()),
			goqu.C(COLUMN_ID).Eq(options.HandleOrID()),
		)
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

	if !options.IsCountOnly() && options.HasOrderBy() {
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
