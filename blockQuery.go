package cmsstore

import "errors"

func BlockQuery() BlockQueryInterface {
	return &blockQuery{
		properties: make(map[string]interface{}),
	}
}

var _ BlockQueryInterface = (*blockQuery)(nil)

type blockQuery struct {
	properties map[string]interface{}
}

func (q *blockQuery) Validate() error {
	if q.HasCreatedAtGte() && q.CreatedAtGte() == "" {
		return errors.New("block query. created_at_gte cannot be empty")
	}

	if q.HasCreatedAtLte() && q.CreatedAtLte() == "" {
		return errors.New("block query. created_at_lte cannot be empty")
	}

	if q.HasID() && q.ID() == "" {
		return errors.New("block query. id cannot be empty")
	}

	if q.HasIDIn() && len(q.IDIn()) < 1 {
		return errors.New("block query. id_in cannot be empty array")
	}

	if q.HasLimit() && q.Limit() < 0 {
		return errors.New("block query. limit cannot be negative")
	}

	if q.HasHandle() && q.Handle() == "" {
		return errors.New("block query. handle cannot be empty")
	}

	if q.HasNameLike() && q.NameLike() == "" {
		return errors.New("block query. name_like cannot be empty")
	}

	if q.HasStatus() && q.Status() == "" {
		return errors.New("block query. status cannot be empty")
	}

	if q.HasOffset() && q.Offset() < 0 {
		return errors.New("block query. offset cannot be negative")
	}

	return nil
}

func (q *blockQuery) HasCreatedAtGte() bool {
	return q.hasProperty("created_at_gte")
}

func (q *blockQuery) HasCreatedAtLte() bool {
	return q.hasProperty("created_at_lte")
}

func (q *blockQuery) HasHandle() bool {
	return q.hasProperty("handle")
}

func (q *blockQuery) HasID() bool {
	return q.hasProperty("id")
}

func (q *blockQuery) HasIDIn() bool {
	return q.hasProperty("id_in")
}

func (q *blockQuery) HasLimit() bool {
	return q.hasProperty("limit")
}

func (q *blockQuery) HasNameLike() bool {
	return q.hasProperty("name_like")
}

func (q *blockQuery) HasOffset() bool {
	return q.hasProperty("offset")
}

func (q *blockQuery) HasOrderBy() bool {
	return q.hasProperty("order_by")
}

func (q *blockQuery) HasSoftDeleted() bool {
	return q.hasProperty("soft_deleted")
}

func (q *blockQuery) HasSortOrder() bool {
	return q.hasProperty("sort_order")
}

func (q *blockQuery) HasStatus() bool {
	return q.hasProperty("status")
}

func (q *blockQuery) HasStatusIn() bool {
	return q.hasProperty("status_in")
}

func (q *blockQuery) IsCountOnly() bool {
	if q.hasProperty("count_only") {
		return q.properties["count_only"].(bool)
	}

	return false
}

func (q *blockQuery) IncludeSoftDeleted() bool {
	if q.hasProperty("soft_deleted") {
		return q.properties["soft_deleted"].(bool)
	}

	return false
}

func (q *blockQuery) CreatedAtGte() string {
	return q.properties["created_at_gte"].(string)
}

func (q *blockQuery) CreatedAtLte() string {
	return q.properties["created_at_lte"].(string)
}

func (q *blockQuery) ID() string {
	return q.properties["id"].(string)
}

func (q *blockQuery) IDIn() []string {
	return q.properties["id_in"].([]string)
}

func (q *blockQuery) Limit() int {
	return q.properties["limit"].(int)
}

func (q *blockQuery) NameLike() string {
	return q.properties["name_like"].(string)
}

func (q *blockQuery) Offset() int {
	return q.properties["offset"].(int)
}

func (q *blockQuery) OrderBy() string {
	return q.properties["order_by"].(string)
}

func (q *blockQuery) SoftDeleteIncluded() bool {
	return q.properties["soft_delete_included"].(bool)
}

func (q *blockQuery) SortOrder() string {
	return q.properties["sort_order"].(string)
}

func (q *blockQuery) Status() string {
	return q.properties["status"].(string)
}

func (q *blockQuery) StatusIn() []string {
	return q.properties["status_in"].([]string)
}

func (q *blockQuery) Handle() string {
	return q.properties["handle"].(string)
}

func (q *blockQuery) SetCountOnly(countOnly bool) BlockQueryInterface {
	q.properties["count_only"] = countOnly
	return q
}

func (q *blockQuery) SetID(id string) BlockQueryInterface {
	q.properties["id"] = id
	return q
}

func (q *blockQuery) SetIDIn(idIn []string) BlockQueryInterface {
	q.properties["id_in"] = idIn
	return q
}

func (q *blockQuery) SetHandle(handle string) BlockQueryInterface {
	q.properties["handle"] = handle
	return q
}

func (q *blockQuery) SetLimit(limit int) BlockQueryInterface {
	q.properties["limit"] = limit
	return q
}

func (q *blockQuery) SetNameLike(nameLike string) BlockQueryInterface {
	q.properties["name_like"] = nameLike
	return q
}

func (q *blockQuery) SetOffset(offset int) BlockQueryInterface {
	q.properties["offset"] = offset
	return q
}

func (q *blockQuery) SetOrderBy(orderBy string) BlockQueryInterface {
	q.properties["order_by"] = orderBy
	return q
}

func (q *blockQuery) SetSoftDeleteIncluded(SoftDeleteIncluded bool) BlockQueryInterface {
	q.properties["soft_delete_included"] = SoftDeleteIncluded
	return q
}

func (q *blockQuery) SetSortOrder(sortOrder string) BlockQueryInterface {
	q.properties["sort_order"] = sortOrder
	return q
}

func (q *blockQuery) SetStatus(status string) BlockQueryInterface {
	q.properties["status"] = status
	return q
}

func (q *blockQuery) SetStatusIn(statusIn []string) BlockQueryInterface {
	q.properties["status_in"] = statusIn
	return q
}

func (q *blockQuery) hasProperty(key string) bool {
	_, ok := q.properties[key]
	return ok
}

// type blockQuery struct {
// 	id              string
// 	idIn            []string
// 	nameLike        string
// 	status          string
// 	statusIn        []string
// 	handle          string
// 	createdAtGte    string
// 	createdAtLte    string
// 	countOnly       bool
// 	offset          int64
// 	limit           int
// 	sortOrder       string
// 	orderBy         string
// 	withSoftDeleted bool
// }

// func NewBlockQuery() BlockQueryInterface {
// 	return &blockQuery{}
// }

// var _ BlockQueryInterface = (*blockQuery)(nil)

// func (q *blockQuery) ID() string {
// 	return q.id
// }

// func (q *blockQuery) SetID(id string) (BlockQueryInterface, error) {
// 	if id == "" {
// 		return q, errors.New(ERROR_EMPTY_STRING)
// 	}

// 	q.id = id

// 	return q, nil
// }

// func (q *blockQuery) IDIn() []string {
// 	return q.idIn
// }

// func (q *blockQuery) SetIDIn(idIn []string) (BlockQueryInterface, error) {
// 	if len(idIn) < 1 {
// 		return q, errors.New(ERROR_EMPTY_ARRAY)
// 	}

// 	q.idIn = idIn

// 	return q, nil
// }

// func (q *blockQuery) NameLike() string {
// 	return q.nameLike
// }

// func (q *blockQuery) SetNameLike(nameLike string) (BlockQueryInterface, error) {
// 	if nameLike == "" {
// 		return q, errors.New(ERROR_EMPTY_STRING)
// 	}
// 	q.nameLike = nameLike
// 	return q, nil
// }

// func (q *blockQuery) Status() string {
// 	return q.status
// }

// func (q *blockQuery) SetStatus(status string) (BlockQueryInterface, error) {
// 	if status == "" {
// 		return q, errors.New(ERROR_EMPTY_STRING)
// 	}
// 	q.status = status
// 	return q, nil
// }

// func (q *blockQuery) StatusIn() []string {
// 	return q.statusIn
// }

// func (q *blockQuery) SetStatusIn(statusIn []string) (BlockQueryInterface, error) {
// 	if len(statusIn) < 1 {
// 		return q, errors.New(ERROR_EMPTY_ARRAY)
// 	}
// 	q.statusIn = statusIn
// 	return q, nil
// }

// func (q *blockQuery) Handle() string {
// 	return q.handle
// }

// func (q *blockQuery) SetHandle(handle string) (BlockQueryInterface, error) {
// 	if handle == "" {
// 		return q, errors.New(ERROR_EMPTY_STRING)
// 	}
// 	q.handle = handle
// 	return q, nil
// }

// func (q *blockQuery) CreatedAtGte() string {
// 	return q.createdAtGte
// }

// func (q *blockQuery) SetCreatedAtGte(createdAtGte string) (BlockQueryInterface, error) {
// 	if createdAtGte == "" {
// 		return q, errors.New(ERROR_EMPTY_STRING)
// 	}
// 	q.createdAtGte = createdAtGte
// 	return q, nil
// }

// func (q *blockQuery) CreatedAtLte() string {
// 	return q.createdAtLte
// }

// func (q *blockQuery) SetCreatedAtLte(createdAtLte string) (BlockQueryInterface, error) {
// 	if createdAtLte == "" {
// 		return q, errors.New(ERROR_EMPTY_STRING)
// 	}
// 	q.createdAtLte = createdAtLte
// 	return q, nil
// }

// func (q *blockQuery) Offset() int {
// 	return int(q.offset)
// }

// func (q *blockQuery) SetOffset(offset int) (BlockQueryInterface, error) {
// 	if offset < 0 {
// 		return q, errors.New(ERROR_NEGATIVE_NUMBER)
// 	}
// 	q.offset = int64(offset)
// 	return q, nil
// }

// func (q *blockQuery) Limit() int {
// 	return q.limit
// }

// func (q *blockQuery) SetLimit(limit int) (BlockQueryInterface, error) {
// 	if limit < 1 {
// 		return q, errors.New(ERROR_NEGATIVE_NUMBER)
// 	}
// 	q.limit = limit
// 	return q, nil
// }

// func (q *blockQuery) SortOrder() string {
// 	return q.sortOrder
// }

// func (q *blockQuery) SetSortOrder(sortOrder string) (BlockQueryInterface, error) {
// 	if sortOrder == "" {
// 		return q, errors.New(ERROR_EMPTY_STRING)
// 	}
// 	q.sortOrder = sortOrder
// 	return q, nil
// }

// func (q *blockQuery) OrderBy() string {
// 	return q.orderBy
// }

// func (q *blockQuery) SetOrderBy(orderBy string) (BlockQueryInterface, error) {
// 	if orderBy == "" {
// 		return q, errors.New(ERROR_EMPTY_STRING)
// 	}
// 	q.orderBy = orderBy
// 	return q, nil
// }

// func (q *blockQuery) CountOnly() bool {
// 	return q.countOnly
// }

// func (q *blockQuery) SetCountOnly(countOnly bool) BlockQueryInterface {
// 	q.countOnly = countOnly
// 	return q
// }

// func (q *blockQuery) WithSoftDeleted() bool {
// 	return q.withSoftDeleted
// }

// func (q *blockQuery) SetWithSoftDeleted(withSoftDeleted bool) BlockQueryInterface {
// 	q.withSoftDeleted = withSoftDeleted
// 	return q
// }
