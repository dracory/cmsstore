package cmsstore

import "errors"

// page_query.go provides the implementation for the PageQuery interface, which
// is used to construct and validate page queries.

// == CONSTRUCTOR ============================================================

// PageQuery returns a new instance of PageQueryInterface with an initialized
// map to hold query parameters.
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

// Validate checks the validity of the page query parameters.
func (p *pageQuery) Validate() error {
	if p.parameters == nil {
		return errors.New("page query: parameters cannot be nil")
	}

	if p.HasAliasLike() && p.AliasLike() == "" {
		return errors.New("page query: alias_like cannot be empty")
	}

	if p.HasCreatedAtGte() && p.CreatedAtGte() == "" {
		return errors.New("page query: created_at_gte cannot be empty")
	}

	if p.HasCreatedAtLte() && p.CreatedAtLte() == "" {
		return errors.New("page query: created_at_lte cannot be empty")
	}

	if p.HasID() && p.ID() == "" {
		return errors.New("page query: id cannot be empty")
	}

	if p.HasIDIn() && len(p.IDIn()) < 1 {
		return errors.New("page query: id_in cannot be empty array")
	}

	if p.HasLimit() && p.Limit() < 0 {
		return errors.New("page query: limit cannot be negative")
	}

	if p.HasHandle() && p.Handle() == "" {
		return errors.New("page query: handle cannot be empty")
	}

	if p.HasNameLike() && p.NameLike() == "" {
		return errors.New("page query: name_like cannot be empty")
	}

	if p.HasOffset() && p.Offset() < 0 {
		return errors.New("page query: offset cannot be negative")
	}

	if p.HasOrderBy() && p.OrderBy() == "" {
		return errors.New("page query: order_by cannot be empty")
	}

	if p.HasStatus() && p.Status() == "" {
		return errors.New("page query: status cannot be empty")
	}

	if p.HasStatusIn() && len(p.StatusIn()) < 1 {
		return errors.New("page query: status_in cannot be empty array")
	}

	if p.HasTemplateID() && p.TemplateID() == "" {
		return errors.New("page query: template_id cannot be empty")
	}

	return nil
}

// HasColumns checks if the Columns parameter is set.
func (p *pageQuery) HasColumns() bool {
	return p.hasParameter(propertyKeyColumns)
}

// HasAlias checks if the Alias parameter is set.
func (p *pageQuery) HasAlias() bool {
	return p.hasParameter(propertyKeyAlias)
}

// Alias returns the value of the Alias parameter.
func (p *pageQuery) Alias() string {
	return p.parameters[propertyKeyAlias].(string)
}

// SetAlias sets the value of the Alias parameter.
func (p *pageQuery) SetAlias(alias string) PageQueryInterface {
	p.parameters[propertyKeyAlias] = alias
	return p
}

// HasAliasLike checks if the AliasLike parameter is set.
func (p *pageQuery) HasAliasLike() bool {
	return p.hasParameter(propertyKeyAliasLike)
}

// AliasLike returns the value of the AliasLike parameter.
func (p *pageQuery) AliasLike() string {
	return p.parameters[propertyKeyAliasLike].(string)
}

// SetAliasLike sets the value of the AliasLike parameter.
func (p *pageQuery) SetAliasLike(nameLike string) PageQueryInterface {
	p.parameters[propertyKeyAliasLike] = nameLike
	return p
}

// Columns returns the value of the Columns parameter.
func (p *pageQuery) Columns() []string {
	if p.parameters[propertyKeyColumns] == nil {
		return []string{}
	}
	return p.parameters[propertyKeyColumns].([]string)
}

// SetColumns sets the value of the Columns parameter.
func (p *pageQuery) SetColumns(columns []string) PageQueryInterface {
	p.parameters[propertyKeyColumns] = columns
	return p
}

// HasCreatedAtGte checks if the CreatedAtGte parameter is set.
func (p *pageQuery) HasCreatedAtGte() bool {
	return p.hasParameter(propertyKeyCreatedAtGte)
}

// CreatedAtGte returns the value of the CreatedAtGte parameter.
func (p *pageQuery) CreatedAtGte() string {
	return p.parameters[propertyKeyCreatedAtGte].(string)
}

// SetCreatedAtGte sets the value of the CreatedAtGte parameter.
func (p *pageQuery) SetCreatedAtGte(createdAtGte string) PageQueryInterface {
	p.parameters[propertyKeyCreatedAtGte] = createdAtGte
	return p
}

// HasCreatedAtLte checks if the CreatedAtLte parameter is set.
func (p *pageQuery) HasCreatedAtLte() bool {
	return p.hasParameter(propertyKeyCreatedAtLte)
}

// CreatedAtLte returns the value of the CreatedAtLte parameter.
func (p *pageQuery) CreatedAtLte() string {
	return p.parameters[propertyKeyCreatedAtLte].(string)
}

// SetCreatedAtLte sets the value of the CreatedAtLte parameter.
func (p *pageQuery) SetCreatedAtLte(createdAtLte string) PageQueryInterface {
	p.parameters[propertyKeyCreatedAtLte] = createdAtLte
	return p
}

// HasCountOnly checks if the CountOnly parameter is set.
func (p *pageQuery) HasCountOnly() bool {
	return p.hasParameter(propertyKeyCountOnly)
}

// IsCountOnly returns the value of the CountOnly parameter.
func (p *pageQuery) IsCountOnly() bool {
	if !p.HasCountOnly() {
		return false
	}
	return p.parameters[propertyKeyCountOnly].(bool)
}

// SetCountOnly sets the value of the CountOnly parameter.
func (p *pageQuery) SetCountOnly(isCountOnly bool) PageQueryInterface {
	p.parameters[propertyKeyCountOnly] = isCountOnly
	return p
}

// HasHandle checks if the Handle parameter is set.
func (p *pageQuery) HasHandle() bool {
	return p.hasParameter(propertyKeyHandle)
}

// Handle returns the value of the Handle parameter.
func (p *pageQuery) Handle() string {
	return p.parameters[propertyKeyHandle].(string)
}

// SetHandle sets the value of the Handle parameter.
func (p *pageQuery) SetHandle(handle string) PageQueryInterface {
	p.parameters[propertyKeyHandle] = handle
	return p
}

// HasID checks if the ID parameter is set.
func (p *pageQuery) HasID() bool {
	return p.hasParameter(propertyKeyId)
}

// ID returns the value of the ID parameter.
func (p *pageQuery) ID() string {
	return p.parameters[propertyKeyId].(string)
}

// SetID sets the value of the ID parameter.
func (p *pageQuery) SetID(id string) PageQueryInterface {
	p.parameters[propertyKeyId] = id
	return p
}

// HasIDIn checks if the IDIn parameter is set.
func (p *pageQuery) HasIDIn() bool {
	return p.hasParameter(propertyKeyIdIn)
}

// IDIn returns the value of the IDIn parameter.
func (p *pageQuery) IDIn() []string {
	return p.parameters[propertyKeyIdIn].([]string)
}

// SetIDIn sets the value of the IDIn parameter.
func (p *pageQuery) SetIDIn(idIn []string) PageQueryInterface {
	p.parameters[propertyKeyIdIn] = idIn
	return p
}

// HasLimit checks if the Limit parameter is set.
func (p *pageQuery) HasLimit() bool {
	return p.hasParameter(propertyKeyLimit)
}

// Limit returns the value of the Limit parameter.
func (p *pageQuery) Limit() int {
	return p.parameters[propertyKeyLimit].(int)
}

// SetLimit sets the value of the Limit parameter.
func (p *pageQuery) SetLimit(limit int) PageQueryInterface {
	p.parameters[propertyKeyLimit] = limit
	return p
}

// HasNameLike checks if the NameLike parameter is set.
func (p *pageQuery) HasNameLike() bool {
	return p.hasParameter(propertyKeyNameLike)
}

// NameLike returns the value of the NameLike parameter.
func (p *pageQuery) NameLike() string {
	return p.parameters[propertyKeyNameLike].(string)
}

// SetNameLike sets the value of the NameLike parameter.
func (p *pageQuery) SetNameLike(nameLike string) PageQueryInterface {
	p.parameters[propertyKeyNameLike] = nameLike
	return p
}

// HasOffset checks if the Offset parameter is set.
func (p *pageQuery) HasOffset() bool {
	return p.hasParameter(propertyKeyOffset)
}

// Offset returns the value of the Offset parameter.
func (p *pageQuery) Offset() int {
	return p.parameters[propertyKeyOffset].(int)
}

// SetOffset sets the value of the Offset parameter.
func (p *pageQuery) SetOffset(offset int) PageQueryInterface {
	p.parameters[propertyKeyOffset] = offset
	return p
}

// HasOrderBy checks if the OrderBy parameter is set.
func (p *pageQuery) HasOrderBy() bool {
	return p.hasParameter(propertyKeyOrderBy)
}

// OrderBy returns the value of the OrderBy parameter.
func (p *pageQuery) OrderBy() string {
	return p.parameters[propertyKeyOrderBy].(string)
}

// SetOrderBy sets the value of the OrderBy parameter.
func (p *pageQuery) SetOrderBy(orderBy string) PageQueryInterface {
	p.parameters[propertyKeyOrderBy] = orderBy
	return p
}

// HasSiteID checks if the SiteID parameter is set.
func (p *pageQuery) HasSiteID() bool {
	return p.hasParameter(propertyKeySiteID)
}

// SiteID returns the value of the SiteID parameter.
func (p *pageQuery) SiteID() string {
	return p.parameters[propertyKeySiteID].(string)
}

// SetSiteID sets the value of the SiteID parameter.
func (p *pageQuery) SetSiteID(siteID string) PageQueryInterface {
	p.parameters[propertyKeySiteID] = siteID
	return p
}

// HasSoftDeletedIncluded checks if the SoftDeletedIncluded parameter is set.
func (p *pageQuery) HasSoftDeletedIncluded() bool {
	return p.hasParameter(propertyKeySoftDeleteIncluded)
}

// SoftDeletedIncluded returns the value of the SoftDeletedIncluded parameter.
func (p *pageQuery) SoftDeletedIncluded() bool {
	if !p.HasSoftDeletedIncluded() {
		return false
	}
	return p.parameters[propertyKeySoftDeleteIncluded].(bool)
}

// SetSoftDeletedIncluded sets the value of the SoftDeletedIncluded parameter.
func (p *pageQuery) SetSoftDeletedIncluded(softDeletedIncluded bool) PageQueryInterface {
	p.parameters[propertyKeySoftDeleteIncluded] = softDeletedIncluded
	return p
}

// HasSortOrder checks if the SortOrder parameter is set.
func (p *pageQuery) HasSortOrder() bool {
	return p.hasParameter(propertyKeySortOrder)
}

// SortOrder returns the value of the SortOrder parameter.
func (p *pageQuery) SortOrder() string {
	return p.parameters[propertyKeySortOrder].(string)
}

// SetSortOrder sets the value of the SortOrder parameter.
func (p *pageQuery) SetSortOrder(sortOrder string) PageQueryInterface {
	p.parameters[propertyKeySortOrder] = sortOrder
	return p
}

// HasStatus checks if the Status parameter is set.
func (p *pageQuery) HasStatus() bool {
	return p.hasParameter(propertyKeyStatus)
}

// Status returns the value of the Status parameter.
func (p *pageQuery) Status() string {
	return p.parameters[propertyKeyStatus].(string)
}

// SetStatus sets the value of the Status parameter.
func (p *pageQuery) SetStatus(status string) PageQueryInterface {
	p.parameters[propertyKeyStatus] = status
	return p
}

// HasStatusIn checks if the StatusIn parameter is set.
func (p *pageQuery) HasStatusIn() bool {
	return p.hasParameter(propertyKeyStatusIn)
}

// StatusIn returns the value of the StatusIn parameter.
func (p *pageQuery) StatusIn() []string {
	return p.parameters[propertyKeyStatusIn].([]string)
}

// SetStatusIn sets the value of the StatusIn parameter.
func (p *pageQuery) SetStatusIn(statusIn []string) PageQueryInterface {
	p.parameters[propertyKeyStatusIn] = statusIn
	return p
}

// HasTemplateID checks if the TemplateID parameter is set.
func (p *pageQuery) HasTemplateID() bool {
	return p.hasParameter(propertyKeyTemplateID)
}

// TemplateID returns the value of the TemplateID parameter.
func (p *pageQuery) TemplateID() string {
	return p.parameters[propertyKeyTemplateID].(string)
}

// SetTemplateID sets the value of the TemplateID parameter.
func (p *pageQuery) SetTemplateID(templateID string) PageQueryInterface {
	p.parameters[propertyKeyTemplateID] = templateID
	return p
}

// hasParameter checks if a parameter is set.
func (p *pageQuery) hasParameter(name string) bool {
	_, ok := p.parameters[name]
	return ok
}
