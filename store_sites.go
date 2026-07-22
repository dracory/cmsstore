package cmsstore

import (
	"context"
	"errors"
	"log"
	"strings"

	contractsorm "github.com/dracory/neat/contracts/database/orm"
	"github.com/dromara/carbon/v2"
)

func (store *storeImplementation) SiteCount(ctx context.Context, options SiteQueryInterface) (int64, error) {
	if store.neatDB == nil {
		return -1, errors.New("cms store: database is nil")
	}

	if options == nil {
		return -1, errors.New("site options cannot be nil")
	}

	if options != nil && !options.IsCountOnly() {
		options.SetCountOnly(true)
	}

	q, _, err := store.siteSelectQuery(options)

	if err != nil {
		return -1, err
	}

	var count int64
	err = q.Table(store.siteTableName).Count(&count)
	return count, err
}

func (store *storeImplementation) SiteCreate(ctx context.Context, site SiteInterface) error {
	if store.neatDB == nil {
		return errors.New("sitestore: database is nil")
	}

	if site == nil {
		return errors.New("site is nil")
	}

	site.SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	site.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))

	return store.withTransaction(ctx, func(txCtx context.Context) error {
		data := site.Data()

		if store.debugEnabled {
			log.Println("SiteCreate:", data)
		}

		err := store.neatDB.Query().Table(store.siteTableName).Create(data)

		if err != nil {
			return err
		}

		site.MarkAsNotDirty()

		return store.versioningTrackEntity(txCtx, VERSIONING_TYPE_SITE, site.ID(), site)
	})
}

func (store *storeImplementation) SiteDelete(ctx context.Context, site SiteInterface) error {
	if store.neatDB == nil {
		return errors.New("sitestore: database is nil")
	}

	if site == nil {
		return errors.New("site is nil")
	}

	return store.SiteDeleteByID(ctx, site.ID())
}

func (store *storeImplementation) SiteDeleteByID(ctx context.Context, id string) error {
	if store.neatDB == nil {
		return errors.New("sitestore: database is nil")
	}

	if id == "" {
		return errors.New("site id is empty")
	}

	if store.debugEnabled {
		log.Println("SiteDeleteByID:", id)
	}

	_, err := store.neatDB.Query().Table(store.siteTableName).Where("id = ?", id).Delete()

	return err
}

func (store *storeImplementation) SiteFindByDomainName(ctx context.Context, domainName string) (site SiteInterface, err error) {
	if store.neatDB == nil {
		return nil, errors.New("sitestore: database is nil")
	}

	if domainName == "" {
		return nil, errors.New("site domain is empty")
	}

	list, err := store.SiteList(ctx, SiteQuery().
		SetDomainName(domainName).
		SetLimit(1))

	if err != nil {
		return nil, err
	}

	if len(list) > 0 {
		return list[0], nil
	}

	return nil, nil
}

func (store *storeImplementation) SiteFindByHandle(ctx context.Context, handle string) (site SiteInterface, err error) {
	if store.neatDB == nil {
		return nil, errors.New("sitestore: database is nil")
	}

	if handle == "" {
		return nil, errors.New("site handle is empty")
	}

	list, err := store.SiteList(ctx, SiteQuery().
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

func (store *storeImplementation) SiteFindByID(ctx context.Context, id string) (site SiteInterface, err error) {
	if id == "" {
		return nil, errors.New("site id is empty")
	}

	// Normalize ID to lowercase for consistent lookups
	id = NormalizeID(id)

	// Try direct lookup first (handles both 9-char and 32-char IDs)
	list, err := store.SiteList(ctx, SiteQuery().SetID(id).SetLimit(1))

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
			list, err = store.SiteList(ctx, SiteQuery().SetID(unshortenedID).SetLimit(1))
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

func (store *storeImplementation) SiteList(ctx context.Context, query SiteQueryInterface) ([]SiteInterface, error) {
	if store.neatDB == nil {
		return []SiteInterface{}, errors.New("sitestore: database is nil")
	}

	q, _, err := store.siteSelectQuery(query)

	if err != nil {
		return []SiteInterface{}, err
	}

	type siteRow struct {
		ID            string `db:"id"`
		Name          string `db:"name"`
		Handle        string `db:"handle"`
		Status        string `db:"status"`
		DomainNames   string `db:"domain_names"`
		Metas         string `db:"metas"`
		Memo          string `db:"memo"`
		CreatedAt     string `db:"created_at"`
		UpdatedAt     string `db:"updated_at"`
		SoftDeletedAt string `db:"soft_deleted_at"`
	}

	var rows []siteRow
	if err := q.Table(store.siteTableName).Get(&rows); err != nil {
		return []SiteInterface{}, err
	}

	list := make([]SiteInterface, 0, len(rows))
	for _, r := range rows {
		modelMap := map[string]string{
			"id":              r.ID,
			"name":            r.Name,
			"handle":          r.Handle,
			"status":          r.Status,
			"domain_names":    r.DomainNames,
			"metas":           r.Metas,
			"memo":            r.Memo,
			"created_at":      r.CreatedAt,
			"updated_at":      r.UpdatedAt,
			"soft_deleted_at": r.SoftDeletedAt,
		}
		model := NewSiteFromExistingData(modelMap)
		list = append(list, model)
	}

	return list, nil
}

func (store *storeImplementation) SiteSoftDelete(ctx context.Context, site SiteInterface) error {
	if site == nil {
		return errors.New("site is nil")
	}

	site.SetSoftDeletedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))

	return store.SiteUpdate(ctx, site)
}

func (store *storeImplementation) SiteSoftDeleteByID(ctx context.Context, id string) error {
	if store.neatDB == nil {
		return errors.New("sitestore: database is nil")
	}

	site, err := store.SiteFindByID(ctx, id)

	if err != nil {
		return err
	}

	if site == nil {
		return errors.New("site not found")
	}

	return store.SiteSoftDelete(ctx, site)
}

func (store *storeImplementation) SiteUpdate(ctx context.Context, site SiteInterface) error {
	if store.neatDB == nil {
		return errors.New("sitestore: database is nil")
	}

	if site == nil {
		return errors.New("site is nil")
	}

	site.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString())

	return store.withTransaction(ctx, func(txCtx context.Context) error {
		dataChanged := site.DataChanged()

		delete(dataChanged, COLUMN_ID) // ID is not updateable

		if len(dataChanged) < 1 {
			return nil
		}

		if store.debugEnabled {
			log.Println("SiteUpdate:", dataChanged)
		}

		_, err := store.neatDB.Query().Table(store.siteTableName).Where("id = ?", site.ID()).Update(dataChanged)
		if err != nil {
			return err
		}

		site.MarkAsNotDirty()

		return store.versioningTrackEntity(txCtx, VERSIONING_TYPE_SITE, site.ID(), site)
	})
}

func (store *storeImplementation) siteSelectQuery(options SiteQueryInterface) (query contractsorm.Query, columns []any, err error) {
	if options == nil {
		return nil, []any{}, errors.New("site options cannot be nil")
	}

	if err := options.Validate(); err != nil {
		return nil, []any{}, err
	}

	q := store.neatDB.Query().Table(store.siteTableName)

	if options.HasCreatedAtGte() && options.HasCreatedAtLte() {
		q = q.Where(COLUMN_CREATED_AT+" >= ? AND "+COLUMN_CREATED_AT+" <= ?", options.CreatedAtGte(), options.CreatedAtLte())
	} else if options.HasCreatedAtGte() {
		q = q.Where(COLUMN_CREATED_AT+" >= ?", options.CreatedAtGte())
	} else if options.HasCreatedAtLte() {
		q = q.Where(COLUMN_CREATED_AT+" <= ?", options.CreatedAtLte())
	}

	if options.HasDomainName() {
		q = q.Where(COLUMN_DOMAIN_NAMES+" LIKE ?", `%"`+options.DomainName()+`"%`)
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
		q = q.Where(COLUMN_NAME+" LIKE ?", `%`+options.NameLike()+`%`)
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

	if !options.IsCountOnly() {
		for _, column := range options.Columns() {
			columns = append(columns, column)
		}
	}

	if options.SoftDeletedIncluded() {
		return q, columns, nil // soft deleted sites requested specifically
	}

	q = q.Where(COLUMN_SOFT_DELETED_AT+" > ?", carbon.Now(carbon.UTC).ToDateTimeString())

	return q, columns, nil
}
