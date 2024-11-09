package cmsstore

import "errors"

type templateQuery struct {
	id              string
	idIn            []string
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

func NewTemplateQuery() TemplateQueryInterface {
	return &templateQuery{}
}

var _ TemplateQueryInterface = (*templateQuery)(nil)

func (q *templateQuery) ID() string {
	return q.id
}

func (q *templateQuery) SetID(id string) (TemplateQueryInterface, error) {
	if id == "" {
		return q, errors.New(ERROR_EMPTY_STRING)
	}

	q.id = id

	return q, nil
}

func (q *templateQuery) IDIn() []string {
	return q.idIn
}

func (q *templateQuery) SetIDIn(idIn []string) (TemplateQueryInterface, error) {
	if len(idIn) < 1 {
		return q, errors.New(ERROR_EMPTY_ARRAY)
	}

	q.idIn = idIn

	return q, nil
}

func (q *templateQuery) Status() string {
	return q.status
}

func (q *templateQuery) SetStatus(status string) (TemplateQueryInterface, error) {
	if status == "" {
		return q, errors.New(ERROR_EMPTY_STRING)
	}
	q.status = status
	return q, nil
}

func (q *templateQuery) StatusIn() []string {
	return q.statusIn
}

func (q *templateQuery) SetStatusIn(statusIn []string) (TemplateQueryInterface, error) {
	if len(statusIn) < 1 {
		return q, errors.New(ERROR_EMPTY_ARRAY)
	}
	q.statusIn = statusIn
	return q, nil
}

func (q *templateQuery) Handle() string {
	return q.handle
}

func (q *templateQuery) SetHandle(handle string) (TemplateQueryInterface, error) {
	if handle == "" {
		return q, errors.New(ERROR_EMPTY_STRING)
	}
	q.handle = handle
	return q, nil
}

func (q *templateQuery) CreatedAtGte() string {
	return q.createdAtGte
}

func (q *templateQuery) SetCreatedAtGte(createdAtGte string) (TemplateQueryInterface, error) {
	if createdAtGte == "" {
		return q, errors.New(ERROR_EMPTY_STRING)
	}
	q.createdAtGte = createdAtGte
	return q, nil
}

func (q *templateQuery) CreatedAtLte() string {
	return q.createdAtLte
}

func (q *templateQuery) SetCreatedAtLte(createdAtLte string) (TemplateQueryInterface, error) {
	if createdAtLte == "" {
		return q, errors.New(ERROR_EMPTY_STRING)
	}
	q.createdAtLte = createdAtLte
	return q, nil
}

func (q *templateQuery) Offset() int {
	return int(q.offset)
}

func (q *templateQuery) SetOffset(offset int) (TemplateQueryInterface, error) {
	if offset < 0 {
		return q, errors.New(ERROR_NEGATIVE_NUMBER)
	}
	q.offset = int64(offset)
	return q, nil
}

func (q *templateQuery) Limit() int {
	return q.limit
}

func (q *templateQuery) SetLimit(limit int) (TemplateQueryInterface, error) {
	if limit < 1 {
		return q, errors.New(ERROR_NEGATIVE_NUMBER)
	}
	q.limit = limit
	return q, nil
}

func (q *templateQuery) SortOrder() string {
	return q.sortOrder
}

func (q *templateQuery) SetSortOrder(sortOrder string) (TemplateQueryInterface, error) {
	if sortOrder == "" {
		return q, errors.New(ERROR_EMPTY_STRING)
	}
	q.sortOrder = sortOrder
	return q, nil
}

func (q *templateQuery) OrderBy() string {
	return q.orderBy
}

func (q *templateQuery) SetOrderBy(orderBy string) (TemplateQueryInterface, error) {
	if orderBy == "" {
		return q, errors.New(ERROR_EMPTY_STRING)
	}
	q.orderBy = orderBy
	return q, nil
}

func (q *templateQuery) CountOnly() bool {
	return q.countOnly
}

func (q *templateQuery) SetCountOnly(countOnly bool) TemplateQueryInterface {
	q.countOnly = countOnly
	return q
}

func (q *templateQuery) WithSoftDeleted() bool {
	return q.withSoftDeleted
}

func (q *templateQuery) SetWithSoftDeleted(withSoftDeleted bool) TemplateQueryInterface {
	q.withSoftDeleted = withSoftDeleted
	return q
}
