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

func (p *pageQuery) Validate() error {
	if p.parameters == nil {
		return errors.New("page query. parameters cannot be nil")
	}

	if p.HasAliasLike() && p.AliasLike() == "" {
		return errors.New("page query. alias_like cannot be empty")
	}

	if p.HasCreatedAtGte() && p.CreatedAtGte() == "" {
		return errors.New("page query. created_at_gte cannot be empty")
	}

	if p.HasCreatedAtLte() && p.CreatedAtLte() == "" {
		return errors.New("page query. created_at_lte cannot be empty")
	}

	if p.HasID() && p.ID() == "" {
		return errors.New("page query. id cannot be empty")
	}

	if p.HasIDIn() && len(p.IDIn()) < 1 {
		return errors.New("page query. id_in cannot be empty array")
	}

	if p.HasLimit() && p.Limit() < 0 {
		return errors.New("page query. limit cannot be negative")
	}

	if p.HasHandle() && p.Handle() == "" {
		return errors.New("page query. handle cannot be empty")
	}

	if p.HasNameLike() && p.NameLike() == "" {
		return errors.New("page query. name_like cannot be empty")
	}

	if p.HasOffset() && p.Offset() < 0 {
		return errors.New("page query. offset cannot be negative")
	}

	if p.HasOrderBy() && p.OrderBy() == "" {
		return errors.New("page query. order_by cannot be empty")
	}

	if p.HasStatus() && p.Status() == "" {
		return errors.New("page query. status cannot be empty")
	}

	if p.HasStatusIn() && len(p.StatusIn()) < 1 {
		return errors.New("page query. status_in cannot be empty array")
	}

	if p.HasTemplateID() && p.TemplateID() == "" {
		return errors.New("page query. template_id cannot be empty")
	}

	return nil
}

func (p *pageQuery) HasAlias() bool {
	return p.hasParameter(propertyKeyAlias)
}

func (p *pageQuery) Alias() string {
	return p.parameters[propertyKeyAlias].(string)
}

func (p *pageQuery) SetAlias(alias string) PageQueryInterface {
	p.parameters[propertyKeyAlias] = alias
	return p
}

func (p *pageQuery) HasAliasLike() bool {
	return p.hasParameter(propertyKeyAliasLike)
}

func (p *pageQuery) AliasLike() string {
	return p.parameters[propertyKeyAliasLike].(string)
}

func (p *pageQuery) SetAliasLike(nameLike string) PageQueryInterface {
	p.parameters[propertyKeyAliasLike] = nameLike
	return p
}

func (p *pageQuery) Columns() []string {
	if p.parameters[propertyKeyColumns] == nil {
		return []string{}
	}
	return p.parameters[propertyKeyColumns].([]string)
}

func (p *pageQuery) SetColumns(columns []string) PageQueryInterface {
	p.parameters[propertyKeyColumns] = columns
	return p
}

func (p *pageQuery) HasCreatedAtGte() bool {
	return p.hasParameter(propertyKeyCreatedAtGte)
}

func (p *pageQuery) CreatedAtGte() string {
	return p.parameters[propertyKeyCreatedAtGte].(string)
}

func (p *pageQuery) SetCreatedAtGte(createdAtGte string) PageQueryInterface {
	p.parameters[propertyKeyCreatedAtGte] = createdAtGte
	return p
}

func (p *pageQuery) HasCreatedAtLte() bool {
	return p.hasParameter(propertyKeyCreatedAtLte)
}

func (p *pageQuery) CreatedAtLte() string {
	return p.parameters[propertyKeyCreatedAtLte].(string)
}

func (p *pageQuery) SetCreatedAtLte(createdAtLte string) PageQueryInterface {
	p.parameters[propertyKeyCreatedAtLte] = createdAtLte
	return p
}

func (p *pageQuery) HasCountOnly() bool {
	return p.hasParameter(propertyKeyCountOnly)
}

func (p *pageQuery) IsCountOnly() bool {
	if !p.HasCountOnly() {
		return false
	}
	return p.parameters[propertyKeyCountOnly].(bool)
}

func (p *pageQuery) SetCountOnly(isCountOnly bool) PageQueryInterface {
	p.parameters[propertyKeyCountOnly] = isCountOnly
	return p
}

func (p *pageQuery) HasHandle() bool {
	return p.hasParameter(propertyKeyHandle)
}

func (p *pageQuery) Handle() string {
	return p.parameters[propertyKeyHandle].(string)
}

func (p *pageQuery) SetHandle(handle string) PageQueryInterface {
	p.parameters[propertyKeyHandle] = handle
	return p
}

func (p *pageQuery) HasID() bool {
	return p.hasParameter(propertyKeyId)
}

func (p *pageQuery) ID() string {
	return p.parameters[propertyKeyId].(string)
}

func (p *pageQuery) SetID(id string) PageQueryInterface {
	p.parameters[propertyKeyId] = id
	return p
}

func (p *pageQuery) HasIDIn() bool {
	return p.hasParameter(propertyKeyIdIn)
}

func (p *pageQuery) IDIn() []string {
	return p.parameters[propertyKeyIdIn].([]string)
}

func (p *pageQuery) SetIDIn(idIn []string) PageQueryInterface {
	p.parameters[propertyKeyIdIn] = idIn
	return p
}

func (p *pageQuery) HasLimit() bool {
	return p.hasParameter(propertyKeyLimit)
}

func (p *pageQuery) Limit() int {
	return p.parameters[propertyKeyLimit].(int)
}

func (p *pageQuery) SetLimit(limit int) PageQueryInterface {
	p.parameters[propertyKeyLimit] = limit
	return p
}

func (p *pageQuery) HasNameLike() bool {
	return p.hasParameter(propertyKeyNameLike)
}

func (p *pageQuery) NameLike() string {
	return p.parameters[propertyKeyNameLike].(string)
}

func (p *pageQuery) SetNameLike(nameLike string) PageQueryInterface {
	p.parameters[propertyKeyNameLike] = nameLike
	return p
}

func (p *pageQuery) HasOffset() bool {
	return p.hasParameter(propertyKeyOffset)
}

func (p *pageQuery) Offset() int {
	return p.parameters[propertyKeyOffset].(int)
}

func (p *pageQuery) SetOffset(offset int) PageQueryInterface {
	p.parameters[propertyKeyOffset] = offset
	return p
}

func (p *pageQuery) HasOrderBy() bool {
	return p.hasParameter(propertyKeyOrderBy)
}

func (p *pageQuery) OrderBy() string {
	return p.parameters[propertyKeyOrderBy].(string)
}

func (p *pageQuery) SetOrderBy(orderBy string) PageQueryInterface {
	p.parameters[propertyKeyOrderBy] = orderBy
	return p
}

func (p *pageQuery) HasSiteID() bool {
	return p.hasParameter(propertyKeySiteID)
}

func (p *pageQuery) SiteID() string {
	return p.parameters[propertyKeySiteID].(string)
}

func (p *pageQuery) SetSiteID(siteID string) PageQueryInterface {
	p.parameters[propertyKeySiteID] = siteID
	return p
}

func (p *pageQuery) HasSoftDeletedIncluded() bool {
	return p.hasParameter(propertyKeySoftDeleteIncluded)
}

func (p *pageQuery) SoftDeletedIncluded() bool {
	if !p.HasSoftDeletedIncluded() {
		return false
	}
	return p.parameters[propertyKeySoftDeleteIncluded].(bool)
}

func (p *pageQuery) SetSoftDeletedIncluded(softDeletedIncluded bool) PageQueryInterface {
	p.parameters[propertyKeySoftDeleteIncluded] = softDeletedIncluded
	return p
}

func (p *pageQuery) HasSortOrder() bool {
	return p.hasParameter(propertyKeySortOrder)
}

func (p *pageQuery) SortOrder() string {
	return p.parameters[propertyKeySortOrder].(string)
}

func (p *pageQuery) SetSortOrder(sortOrder string) PageQueryInterface {
	p.parameters[propertyKeySortOrder] = sortOrder
	return p
}

func (p *pageQuery) HasStatus() bool {
	return p.hasParameter(propertyKeyStatus)
}

func (p *pageQuery) Status() string {
	return p.parameters[propertyKeyStatus].(string)
}

func (p *pageQuery) SetStatus(status string) PageQueryInterface {
	p.parameters[propertyKeyStatus] = status
	return p
}

func (p *pageQuery) HasStatusIn() bool {
	return p.hasParameter(propertyKeyStatusIn)
}

func (p *pageQuery) StatusIn() []string {
	return p.parameters[propertyKeyStatusIn].([]string)
}

func (p *pageQuery) SetStatusIn(statusIn []string) PageQueryInterface {
	p.parameters[propertyKeyStatusIn] = statusIn
	return p
}

func (p *pageQuery) HasTemplateID() bool {
	return p.hasParameter(propertyKeyTemplateID)
}

func (p *pageQuery) TemplateID() string {
	return p.parameters[propertyKeyTemplateID].(string)
}

func (p *pageQuery) SetTemplateID(templateID string) PageQueryInterface {
	p.parameters[propertyKeyTemplateID] = templateID
	return p
}

func (p *pageQuery) hasParameter(name string) bool {
	_, ok := p.parameters[name]
	return ok
}
