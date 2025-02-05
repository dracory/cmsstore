package cmsstore

import "errors"

// == CONSTRUCTOR ============================================================

func PageQuery() PageQueryInterface {
	return &pageQuery{
		parameters: make(map[string]any),
	}
}

// == TYPE ===================================================================

type pageQuery struct {
	parameters map[string]any
}

// == INTERFACE VERIFICATION =================================================

var _ PageQueryInterface = (*pageQuery)(nil)

// == INTERFACE IMPLEMENTATION ===============================================

func (q *pageQuery) Validate() error {
	if q.parameters == nil {
		return errors.New("page query. parameters cannot be nil")
	}

	if q.HasAliasLike() && q.AliasLike() == "" {
		return errors.New("page query. alias_like cannot be empty")
	}

	if q.HasCreatedAtGte() && q.CreatedAtGte() == "" {
		return errors.New("page query. created_at_gte cannot be empty")
	}

	if q.HasCreatedAtLte() && q.CreatedAtLte() == "" {
		return errors.New("page query. created_at_lte cannot be empty")
	}

	if q.HasID() && q.ID() == "" {
		return errors.New("page query. id cannot be empty")
	}

	if q.HasIDIn() && len(q.IDIn()) < 1 {
		return errors.New("page query. id_in cannot be empty array")
	}

	if q.HasLimit() && q.Limit() < 0 {
		return errors.New("page query. limit cannot be negative")
	}

	if q.HasHandle() && q.Handle() == "" {
		return errors.New("page query. handle cannot be empty")
	}

	if q.HasNameLike() && q.NameLike() == "" {
		return errors.New("page query. name_like cannot be empty")
	}

	if q.HasOffset() && q.Offset() < 0 {
		return errors.New("page query. offset cannot be negative")
	}

	if q.HasOrderBy() && q.OrderBy() == "" {
		return errors.New("page query. order_by cannot be empty")
	}

	if q.HasStatus() && q.Status() == "" {
		return errors.New("page query. status cannot be empty")
	}

	if q.HasStatusIn() && len(q.StatusIn()) < 1 {
		return errors.New("page query. status_in cannot be empty array")
	}

	if q.HasTemplateID() && q.TemplateID() == "" {
		return errors.New("page query. template_id cannot be empty")
	}

	return nil
}

func (p *pageQuery) HasAlias() bool {
	return p.hasParameter("alias")
}

func (p *pageQuery) Alias() string {
	return p.parameters["alias"].(string)
}

func (p *pageQuery) SetAlias(alias string) PageQueryInterface {
	p.parameters["alias"] = alias
	return p
}

func (p *pageQuery) HasAliasLike() bool {
	return p.hasParameter("alias_like")
}

func (p *pageQuery) AliasLike() string {
	return p.parameters["alias_like"].(string)
}

func (p *pageQuery) SetAliasLike(nameLike string) PageQueryInterface {
	p.parameters["alias_like"] = nameLike
	return p
}

func (p *pageQuery) Columns() []string {
	if p.parameters["columns"] == nil {
		return []string{}
	}
	return p.parameters["columns"].([]string)
}

func (p *pageQuery) SetColumns(columns []string) PageQueryInterface {
	p.parameters["columns"] = columns
	return p
}

func (p *pageQuery) HasCreatedAtGte() bool {
	return p.hasParameter("created_at_gte")
}

func (p *pageQuery) CreatedAtGte() string {
	return p.parameters["created_at_gte"].(string)
}

func (p *pageQuery) SetCreatedAtGte(createdAtGte string) PageQueryInterface {
	p.parameters["created_at_gte"] = createdAtGte
	return p
}

func (p *pageQuery) HasCreatedAtLte() bool {
	return p.hasParameter("created_at_lte")
}

func (p *pageQuery) CreatedAtLte() string {
	return p.parameters["created_at_lte"].(string)
}

func (p *pageQuery) SetCreatedAtLte(createdAtLte string) PageQueryInterface {
	p.parameters["created_at_lte"] = createdAtLte
	return p
}

func (p *pageQuery) HasCountOnly() bool {
	return p.hasParameter("count_only")
}

func (p *pageQuery) IsCountOnly() bool {
	if !p.HasCountOnly() {
		return false
	}
	return p.parameters["count_only"].(bool)
}

func (p *pageQuery) SetCountOnly(isCountOnly bool) PageQueryInterface {
	p.parameters["count_only"] = isCountOnly
	return p
}

func (p *pageQuery) HasHandle() bool {
	return p.hasParameter("handle")
}

func (p *pageQuery) Handle() string {
	return p.parameters["handle"].(string)
}

func (p *pageQuery) SetHandle(handle string) PageQueryInterface {
	p.parameters["handle"] = handle
	return p
}

func (p *pageQuery) HasID() bool {
	return p.hasParameter("id")
}

func (p *pageQuery) ID() string {
	return p.parameters["id"].(string)
}

func (p *pageQuery) SetID(id string) PageQueryInterface {
	p.parameters["id"] = id
	return p
}

func (p *pageQuery) HasIDIn() bool {
	return p.hasParameter("id_in")
}

func (p *pageQuery) IDIn() []string {
	return p.parameters["id_in"].([]string)
}

func (p *pageQuery) SetIDIn(idIn []string) PageQueryInterface {
	p.parameters["id_in"] = idIn
	return p
}

func (p *pageQuery) HasLimit() bool {
	return p.hasParameter("limit")
}

func (p *pageQuery) Limit() int {
	return p.parameters["limit"].(int)
}

func (p *pageQuery) SetLimit(limit int) PageQueryInterface {
	p.parameters["limit"] = limit
	return p
}

func (p *pageQuery) HasNameLike() bool {
	return p.hasParameter("name_like")
}

func (p *pageQuery) NameLike() string {
	return p.parameters["name_like"].(string)
}

func (p *pageQuery) SetNameLike(nameLike string) PageQueryInterface {
	p.parameters["name_like"] = nameLike
	return p
}

func (p *pageQuery) HasOffset() bool {
	return p.hasParameter("offset")
}

func (p *pageQuery) Offset() int {
	return p.parameters["offset"].(int)
}

func (p *pageQuery) SetOffset(offset int) PageQueryInterface {
	p.parameters["offset"] = offset
	return p
}

func (p *pageQuery) HasOrderBy() bool {
	return p.hasParameter("order_by")
}

func (p *pageQuery) OrderBy() string {
	return p.parameters["order_by"].(string)
}

func (p *pageQuery) SetOrderBy(orderBy string) PageQueryInterface {
	p.parameters["order_by"] = orderBy
	return p
}

func (p *pageQuery) HasSiteID() bool {
	return p.hasParameter("site_id")
}

func (p *pageQuery) SiteID() string {
	return p.parameters["site_id"].(string)
}

func (p *pageQuery) SetSiteID(siteID string) PageQueryInterface {
	p.parameters["site_id"] = siteID
	return p
}

func (p *pageQuery) HasSoftDeletedIncluded() bool {
	return p.hasParameter("soft_deleted_included")
}

func (p *pageQuery) SoftDeletedIncluded() bool {
	if !p.HasSoftDeletedIncluded() {
		return false
	}
	return p.parameters["soft_deleted_included"].(bool)
}

func (p *pageQuery) SetSoftDeletedIncluded(softDeletedIncluded bool) PageQueryInterface {
	p.parameters["soft_deleted_included"] = softDeletedIncluded
	return p
}

func (p *pageQuery) HasSortOrder() bool {
	return p.hasParameter("sort_order")
}

func (p *pageQuery) SortOrder() string {
	return p.parameters["sort_order"].(string)
}

func (p *pageQuery) SetSortOrder(sortOrder string) PageQueryInterface {
	p.parameters["sort_order"] = sortOrder
	return p
}

func (p *pageQuery) HasStatus() bool {
	return p.hasParameter("status")
}

func (p *pageQuery) Status() string {
	return p.parameters["status"].(string)
}

func (p *pageQuery) SetStatus(status string) PageQueryInterface {
	p.parameters["status"] = status
	return p
}

func (p *pageQuery) HasStatusIn() bool {
	return p.hasParameter("status_in")
}

func (p *pageQuery) StatusIn() []string {
	return p.parameters["status_in"].([]string)
}

func (p *pageQuery) SetStatusIn(statusIn []string) PageQueryInterface {
	p.parameters["status_in"] = statusIn
	return p
}

func (p *pageQuery) HasTemplateID() bool {
	return p.hasParameter("template_id")
}

func (p *pageQuery) TemplateID() string {
	return p.parameters["template_id"].(string)
}

func (p *pageQuery) SetTemplateID(templateID string) PageQueryInterface {
	p.parameters["template_id"] = templateID
	return p
}

func (p *pageQuery) hasParameter(name string) bool {
	_, ok := p.parameters[name]
	return ok
}
