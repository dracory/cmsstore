package cmsstore

import "errors"

type blockQuery struct {
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

func NewBlockQuery() BlockQueryInterface {
	return &blockQuery{}
}

var _ BlockQueryInterface = (*blockQuery)(nil)

func (q *blockQuery) ID() string {
	return q.id
}

func (q *blockQuery) SetID(id string) (BlockQueryInterface, error) {
	if id == "" {
		return q, errors.New(ERROR_EMPTY_STRING)
	}

	q.id = id

	return q, nil
}

func (q *blockQuery) IDIn() []string {
	return q.idIn
}

func (q *blockQuery) SetIDIn(idIn []string) (BlockQueryInterface, error) {
	if len(idIn) < 1 {
		return q, errors.New(ERROR_EMPTY_ARRAY)
	}

	q.idIn = idIn

	return q, nil
}

func (q *blockQuery) NameLike() string {
	return q.nameLike
}

func (q *blockQuery) SetNameLike(nameLike string) (BlockQueryInterface, error) {
	if nameLike == "" {
		return q, errors.New(ERROR_EMPTY_STRING)
	}
	q.nameLike = nameLike
	return q, nil
}

func (q *blockQuery) Status() string {
	return q.status
}

func (q *blockQuery) SetStatus(status string) (BlockQueryInterface, error) {
	if status == "" {
		return q, errors.New(ERROR_EMPTY_STRING)
	}
	q.status = status
	return q, nil
}

func (q *blockQuery) StatusIn() []string {
	return q.statusIn
}

func (q *blockQuery) SetStatusIn(statusIn []string) (BlockQueryInterface, error) {
	if len(statusIn) < 1 {
		return q, errors.New(ERROR_EMPTY_ARRAY)
	}
	q.statusIn = statusIn
	return q, nil
}

func (q *blockQuery) Handle() string {
	return q.handle
}

func (q *blockQuery) SetHandle(handle string) (BlockQueryInterface, error) {
	if handle == "" {
		return q, errors.New(ERROR_EMPTY_STRING)
	}
	q.handle = handle
	return q, nil
}

func (q *blockQuery) CreatedAtGte() string {
	return q.createdAtGte
}

func (q *blockQuery) SetCreatedAtGte(createdAtGte string) (BlockQueryInterface, error) {
	if createdAtGte == "" {
		return q, errors.New(ERROR_EMPTY_STRING)
	}
	q.createdAtGte = createdAtGte
	return q, nil
}

func (q *blockQuery) CreatedAtLte() string {
	return q.createdAtLte
}

func (q *blockQuery) SetCreatedAtLte(createdAtLte string) (BlockQueryInterface, error) {
	if createdAtLte == "" {
		return q, errors.New(ERROR_EMPTY_STRING)
	}
	q.createdAtLte = createdAtLte
	return q, nil
}

func (q *blockQuery) Offset() int {
	return int(q.offset)
}

func (q *blockQuery) SetOffset(offset int) (BlockQueryInterface, error) {
	if offset < 0 {
		return q, errors.New(ERROR_NEGATIVE_NUMBER)
	}
	q.offset = int64(offset)
	return q, nil
}

func (q *blockQuery) Limit() int {
	return q.limit
}

func (q *blockQuery) SetLimit(limit int) (BlockQueryInterface, error) {
	if limit < 1 {
		return q, errors.New(ERROR_NEGATIVE_NUMBER)
	}
	q.limit = limit
	return q, nil
}

func (q *blockQuery) SortOrder() string {
	return q.sortOrder
}

func (q *blockQuery) SetSortOrder(sortOrder string) (BlockQueryInterface, error) {
	if sortOrder == "" {
		return q, errors.New(ERROR_EMPTY_STRING)
	}
	q.sortOrder = sortOrder
	return q, nil
}

func (q *blockQuery) OrderBy() string {
	return q.orderBy
}

func (q *blockQuery) SetOrderBy(orderBy string) (BlockQueryInterface, error) {
	if orderBy == "" {
		return q, errors.New(ERROR_EMPTY_STRING)
	}
	q.orderBy = orderBy
	return q, nil
}

func (q *blockQuery) CountOnly() bool {
	return q.countOnly
}

func (q *blockQuery) SetCountOnly(countOnly bool) BlockQueryInterface {
	q.countOnly = countOnly
	return q
}

func (q *blockQuery) WithSoftDeleted() bool {
	return q.withSoftDeleted
}

func (q *blockQuery) SetWithSoftDeleted(withSoftDeleted bool) BlockQueryInterface {
	q.withSoftDeleted = withSoftDeleted
	return q
}
