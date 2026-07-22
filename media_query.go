package cmsstore

import "errors"

func MediaQuery() MediaQueryInterface {
	return &mediaQuery{
		parameters: make(map[string]any),
	}
}

type mediaQuery struct {
	parameters map[string]any
}

var _ MediaQueryInterface = (*mediaQuery)(nil)

func (p *mediaQuery) Validate() error {
	if p.parameters == nil {
		return errors.New("media query: parameters cannot be nil")
	}

	if p.HasID() && p.ID() == "" {
		return errors.New("media query: id cannot be empty")
	}

	if p.HasIDIn() && len(p.IDIn()) < 1 {
		return errors.New("media query: id_in cannot be empty array")
	}

	if p.HasHandle() && p.Handle() == "" {
		return errors.New("media query: handle cannot be empty")
	}

	if p.HasNameLike() && p.NameLike() == "" {
		return errors.New("media query: name_like cannot be empty")
	}

	if p.HasLimit() && p.Limit() < 0 {
		return errors.New("media query: limit cannot be negative")
	}

	if p.HasOffset() && p.Offset() < 0 {
		return errors.New("media query: offset cannot be negative")
	}

	if p.HasOrderBy() && p.OrderBy() == "" {
		return errors.New("media query: order_by cannot be empty")
	}

	if p.HasStatus() && p.Status() == "" {
		return errors.New("media query: status cannot be empty")
	}

	if p.HasStatusIn() && len(p.StatusIn()) < 1 {
		return errors.New("media query: status_in cannot be empty array")
	}

	return nil
}

func (p *mediaQuery) HasColumns() bool {
	return p.hasParameter(propertyKeyColumns)
}

func (p *mediaQuery) Columns() []string {
	if p.parameters[propertyKeyColumns] == nil {
		return []string{}
	}
	return p.parameters[propertyKeyColumns].([]string)
}

func (p *mediaQuery) SetColumns(columns []string) MediaQueryInterface {
	p.parameters[propertyKeyColumns] = columns
	return p
}

func (p *mediaQuery) HasID() bool {
	return p.hasParameter(propertyKeyId)
}

func (p *mediaQuery) ID() string {
	return p.parameters[propertyKeyId].(string)
}

func (p *mediaQuery) SetID(id string) MediaQueryInterface {
	p.parameters[propertyKeyId] = id
	return p
}

func (p *mediaQuery) HasIDIn() bool {
	return p.hasParameter(propertyKeyIdIn)
}

func (p *mediaQuery) IDIn() []string {
	return p.parameters[propertyKeyIdIn].([]string)
}

func (p *mediaQuery) SetIDIn(idIn []string) MediaQueryInterface {
	p.parameters[propertyKeyIdIn] = idIn
	return p
}

func (p *mediaQuery) HasEntityID() bool {
	return p.hasParameter(propertyKeyEntityID)
}

func (p *mediaQuery) EntityID() string {
	return p.parameters[propertyKeyEntityID].(string)
}

func (p *mediaQuery) SetEntityID(entityID string) MediaQueryInterface {
	p.parameters[propertyKeyEntityID] = entityID
	return p
}

func (p *mediaQuery) HasEntityType() bool {
	return p.hasParameter(propertyKeyEntityType)
}

func (p *mediaQuery) EntityType() string {
	return p.parameters[propertyKeyEntityType].(string)
}

func (p *mediaQuery) SetEntityType(entityType string) MediaQueryInterface {
	p.parameters[propertyKeyEntityType] = entityType
	return p
}

func (p *mediaQuery) HasSiteID() bool {
	return p.hasParameter(propertyKeySiteID)
}

func (p *mediaQuery) SiteID() string {
	return p.parameters[propertyKeySiteID].(string)
}

func (p *mediaQuery) SetSiteID(siteID string) MediaQueryInterface {
	p.parameters[propertyKeySiteID] = siteID
	return p
}

func (p *mediaQuery) HasHandle() bool {
	return p.hasParameter(propertyKeyHandle)
}

func (p *mediaQuery) Handle() string {
	return p.parameters[propertyKeyHandle].(string)
}

func (p *mediaQuery) SetHandle(handle string) MediaQueryInterface {
	p.parameters[propertyKeyHandle] = handle
	return p
}

func (p *mediaQuery) HasExtension() bool {
	return p.hasParameter(propertyKeyExtension)
}

func (p *mediaQuery) Extension() string {
	return p.parameters[propertyKeyExtension].(string)
}

func (p *mediaQuery) SetExtension(extension string) MediaQueryInterface {
	p.parameters[propertyKeyExtension] = extension
	return p
}

func (p *mediaQuery) HasType() bool {
	return p.hasParameter(propertyKeyType)
}

func (p *mediaQuery) Type() string {
	return p.parameters[propertyKeyType].(string)
}

func (p *mediaQuery) SetType(mediaType string) MediaQueryInterface {
	p.parameters[propertyKeyType] = mediaType
	return p
}

func (p *mediaQuery) HasStatus() bool {
	return p.hasParameter(propertyKeyStatus)
}

func (p *mediaQuery) Status() string {
	return p.parameters[propertyKeyStatus].(string)
}

func (p *mediaQuery) SetStatus(status string) MediaQueryInterface {
	p.parameters[propertyKeyStatus] = status
	return p
}

func (p *mediaQuery) HasStatusIn() bool {
	return p.hasParameter(propertyKeyStatusIn)
}

func (p *mediaQuery) StatusIn() []string {
	return p.parameters[propertyKeyStatusIn].([]string)
}

func (p *mediaQuery) SetStatusIn(statusIn []string) MediaQueryInterface {
	p.parameters[propertyKeyStatusIn] = statusIn
	return p
}

func (p *mediaQuery) HasNameLike() bool {
	return p.hasParameter(propertyKeyNameLike)
}

func (p *mediaQuery) NameLike() string {
	return p.parameters[propertyKeyNameLike].(string)
}

func (p *mediaQuery) SetNameLike(nameLike string) MediaQueryInterface {
	p.parameters[propertyKeyNameLike] = nameLike
	return p
}

func (p *mediaQuery) HasCountOnly() bool {
	return p.hasParameter(propertyKeyCountOnly)
}

func (p *mediaQuery) IsCountOnly() bool {
	if !p.HasCountOnly() {
		return false
	}
	return p.parameters[propertyKeyCountOnly].(bool)
}

func (p *mediaQuery) SetCountOnly(isCountOnly bool) MediaQueryInterface {
	p.parameters[propertyKeyCountOnly] = isCountOnly
	return p
}

func (p *mediaQuery) HasLimit() bool {
	return p.hasParameter(propertyKeyLimit)
}

func (p *mediaQuery) Limit() int {
	return p.parameters[propertyKeyLimit].(int)
}

func (p *mediaQuery) SetLimit(limit int) MediaQueryInterface {
	p.parameters[propertyKeyLimit] = limit
	return p
}

func (p *mediaQuery) HasOffset() bool {
	return p.hasParameter(propertyKeyOffset)
}

func (p *mediaQuery) Offset() int {
	return p.parameters[propertyKeyOffset].(int)
}

func (p *mediaQuery) SetOffset(offset int) MediaQueryInterface {
	p.parameters[propertyKeyOffset] = offset
	return p
}

func (p *mediaQuery) HasSortOrder() bool {
	return p.hasParameter(propertyKeySortOrder)
}

func (p *mediaQuery) SortOrder() string {
	return p.parameters[propertyKeySortOrder].(string)
}

func (p *mediaQuery) SetSortOrder(sortOrder string) MediaQueryInterface {
	p.parameters[propertyKeySortOrder] = sortOrder
	return p
}

func (p *mediaQuery) HasOrderBy() bool {
	return p.hasParameter(propertyKeyOrderBy)
}

func (p *mediaQuery) OrderBy() string {
	return p.parameters[propertyKeyOrderBy].(string)
}

func (p *mediaQuery) SetOrderBy(orderBy string) MediaQueryInterface {
	p.parameters[propertyKeyOrderBy] = orderBy
	return p
}

func (p *mediaQuery) HasSoftDeletedIncluded() bool {
	return p.hasParameter(propertyKeySoftDeleteIncluded)
}

func (p *mediaQuery) SoftDeletedIncluded() bool {
	if !p.HasSoftDeletedIncluded() {
		return false
	}
	return p.parameters[propertyKeySoftDeleteIncluded].(bool)
}

func (p *mediaQuery) SetSoftDeletedIncluded(softDeletedIncluded bool) MediaQueryInterface {
	p.parameters[propertyKeySoftDeleteIncluded] = softDeletedIncluded
	return p
}

func (p *mediaQuery) hasParameter(name string) bool {
	_, ok := p.parameters[name]
	return ok
}
