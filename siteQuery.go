package cmsstore

import "errors"

type siteQuery struct {
	id              string
	idIn            []string
	nameLike        string
	status          string
	statusIn        []string
	handle          string
	createdAtGte    string
	createdAtLte    string
	countOnly       bool
	offset          int64
	limit           int
	sortOrder       string
	orderBy         string
	withSoftDeleted bool
}

func NewSiteQuery() SiteQueryInterface {
	return &siteQuery{}
}

var _ SiteQueryInterface = (*siteQuery)(nil)

func (q *siteQuery) ID() string {
	return q.id
}

func (q *siteQuery) SetID(id string) (SiteQueryInterface, error) {
	if id == "" {
		return q, errors.New(ERROR_EMPTY_STRING)
	}

	q.id = id

	return q, nil
}

func (q *siteQuery) IDIn() []string {
	return q.idIn
}

func (q *siteQuery) SetIDIn(idIn []string) (SiteQueryInterface, error) {
	if len(idIn) < 1 {
		return q, errors.New(ERROR_EMPTY_ARRAY)
	}

	q.idIn = idIn

	return q, nil
}

func (q *siteQuery) NameLike() string {
	return q.nameLike
}

func (q *siteQuery) SetNameLike(nameLike string) (SiteQueryInterface, error) {
	if nameLike == "" {
		return q, errors.New(ERROR_EMPTY_STRING)
	}
	q.nameLike = nameLike
	return q, nil
}

func (q *siteQuery) Status() string {
	return q.status
}

func (q *siteQuery) SetStatus(status string) (SiteQueryInterface, error) {
	if status == "" {
		return q, errors.New(ERROR_EMPTY_STRING)
	}
	q.status = status
	return q, nil
}

func (q *siteQuery) StatusIn() []string {
	return q.statusIn
}

func (q *siteQuery) SetStatusIn(statusIn []string) (SiteQueryInterface, error) {
	if len(statusIn) < 1 {
		return q, errors.New(ERROR_EMPTY_ARRAY)
	}
	q.statusIn = statusIn
	return q, nil
}

func (q *siteQuery) Handle() string {
	return q.handle
}

func (q *siteQuery) SetHandle(handle string) (SiteQueryInterface, error) {
	if handle == "" {
		return q, errors.New(ERROR_EMPTY_STRING)
	}
	q.handle = handle
	return q, nil
}

func (q *siteQuery) CreatedAtGte() string {
	return q.createdAtGte
}

func (q *siteQuery) SetCreatedAtGte(createdAtGte string) (SiteQueryInterface, error) {
	if createdAtGte == "" {
		return q, errors.New(ERROR_EMPTY_STRING)
	}
	q.createdAtGte = createdAtGte
	return q, nil
}

func (q *siteQuery) CreatedAtLte() string {
	return q.createdAtLte
}

func (q *siteQuery) SetCreatedAtLte(createdAtLte string) (SiteQueryInterface, error) {
	if createdAtLte == "" {
		return q, errors.New(ERROR_EMPTY_STRING)
	}
	q.createdAtLte = createdAtLte
	return q, nil
}

func (q *siteQuery) Offset() int {
	return int(q.offset)
}

func (q *siteQuery) SetOffset(offset int) (SiteQueryInterface, error) {
	if offset < 0 {
		return q, errors.New(ERROR_NEGATIVE_NUMBER)
	}
	q.offset = int64(offset)
	return q, nil
}

func (q *siteQuery) Limit() int {
	return q.limit
}

func (q *siteQuery) SetLimit(limit int) (SiteQueryInterface, error) {
	if limit < 1 {
		return q, errors.New(ERROR_NEGATIVE_NUMBER)
	}
	q.limit = limit
	return q, nil
}

func (q *siteQuery) SortOrder() string {
	return q.sortOrder
}

func (q *siteQuery) SetSortOrder(sortOrder string) (SiteQueryInterface, error) {
	if sortOrder == "" {
		return q, errors.New(ERROR_EMPTY_STRING)
	}
	q.sortOrder = sortOrder
	return q, nil
}

func (q *siteQuery) OrderBy() string {
	return q.orderBy
}

func (q *siteQuery) SetOrderBy(orderBy string) (SiteQueryInterface, error) {
	if orderBy == "" {
		return q, errors.New(ERROR_EMPTY_STRING)
	}
	q.orderBy = orderBy
	return q, nil
}

func (q *siteQuery) CountOnly() bool {
	return q.countOnly
}

func (q *siteQuery) SetCountOnly(countOnly bool) SiteQueryInterface {
	q.countOnly = countOnly
	return q
}

func (q *siteQuery) WithSoftDeleted() bool {
	return q.withSoftDeleted
}

func (q *siteQuery) SetWithSoftDeleted(withSoftDeleted bool) SiteQueryInterface {
	q.withSoftDeleted = withSoftDeleted
	return q
}
