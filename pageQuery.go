package cmsstore

import "errors"

type pageQuery struct {
	id              string
	idIn            []string
	handle          string
	nameLike        string
	status          string
	statusIn        []string
	createdAtGte    string
	createdAtLte    string
	countOnly       bool
	offset          int64
	limit           int
	sortOrder       string
	orderBy         string
	withSoftDeleted bool
}

func NewPageQuery() PageQueryInterface {
	return &pageQuery{}
}

var _ PageQueryInterface = (*pageQuery)(nil)

func (q *pageQuery) ID() string {
	return q.id
}

func (q *pageQuery) SetID(id string) (PageQueryInterface, error) {
	if id == "" {
		return q, errors.New(ERROR_EMPTY_STRING)
	}

	q.id = id

	return q, nil
}

func (q *pageQuery) IDIn() []string {
	return q.idIn
}

func (q *pageQuery) SetIDIn(idIn []string) (PageQueryInterface, error) {
	if len(idIn) < 1 {
		return q, errors.New(ERROR_EMPTY_ARRAY)
	}

	q.idIn = idIn

	return q, nil
}

func (q *pageQuery) NameLike() string {
	return q.nameLike
}

func (q *pageQuery) SetNameLike(nameLike string) (PageQueryInterface, error) {
	if nameLike == "" {
		return q, errors.New(ERROR_EMPTY_STRING)
	}
	q.nameLike = nameLike
	return q, nil
}

func (q *pageQuery) Status() string {
	return q.status
}

func (q *pageQuery) SetStatus(status string) (PageQueryInterface, error) {
	if status == "" {
		return q, errors.New(ERROR_EMPTY_STRING)
	}
	q.status = status
	return q, nil
}

func (q *pageQuery) StatusIn() []string {
	return q.statusIn
}

func (q *pageQuery) SetStatusIn(statusIn []string) (PageQueryInterface, error) {
	if len(statusIn) < 1 {
		return q, errors.New(ERROR_EMPTY_ARRAY)
	}
	q.statusIn = statusIn
	return q, nil
}

func (q *pageQuery) Handle() string {
	return q.handle
}

func (q *pageQuery) SetHandle(handle string) (PageQueryInterface, error) {
	if handle == "" {
		return q, errors.New(ERROR_EMPTY_STRING)
	}
	q.handle = handle
	return q, nil
}

func (q *pageQuery) CreatedAtGte() string {
	return q.createdAtGte
}

func (q *pageQuery) SetCreatedAtGte(createdAtGte string) (PageQueryInterface, error) {
	if createdAtGte == "" {
		return q, errors.New(ERROR_EMPTY_STRING)
	}
	q.createdAtGte = createdAtGte
	return q, nil
}

func (q *pageQuery) CreatedAtLte() string {
	return q.createdAtLte
}

func (q *pageQuery) SetCreatedAtLte(createdAtLte string) (PageQueryInterface, error) {
	if createdAtLte == "" {
		return q, errors.New(ERROR_EMPTY_STRING)
	}
	q.createdAtLte = createdAtLte
	return q, nil
}

func (q *pageQuery) Offset() int {
	return int(q.offset)
}

func (q *pageQuery) SetOffset(offset int) (PageQueryInterface, error) {
	if offset < 0 {
		return q, errors.New(ERROR_NEGATIVE_NUMBER)
	}
	q.offset = int64(offset)
	return q, nil
}

func (q *pageQuery) Limit() int {
	return q.limit
}

func (q *pageQuery) SetLimit(limit int) (PageQueryInterface, error) {
	if limit < 1 {
		return q, errors.New(ERROR_NEGATIVE_NUMBER)
	}
	q.limit = limit
	return q, nil
}

func (q *pageQuery) SortOrder() string {
	return q.sortOrder
}

func (q *pageQuery) SetSortOrder(sortOrder string) (PageQueryInterface, error) {
	if sortOrder == "" {
		return q, errors.New(ERROR_EMPTY_STRING)
	}
	q.sortOrder = sortOrder
	return q, nil
}

func (q *pageQuery) OrderBy() string {
	return q.orderBy
}

func (q *pageQuery) SetOrderBy(orderBy string) (PageQueryInterface, error) {
	if orderBy == "" {
		return q, errors.New(ERROR_EMPTY_STRING)
	}
	q.orderBy = orderBy
	return q, nil
}

func (q *pageQuery) CountOnly() bool {
	return q.countOnly
}

func (q *pageQuery) SetCountOnly(countOnly bool) PageQueryInterface {
	q.countOnly = countOnly
	return q
}

func (q *pageQuery) WithSoftDeleted() bool {
	return q.withSoftDeleted
}

func (q *pageQuery) SetWithSoftDeleted(withSoftDeleted bool) PageQueryInterface {
	q.withSoftDeleted = withSoftDeleted
	return q
}
