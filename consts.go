package cmsstore

// Block Statuses
const (
	BLOCK_STATUS_DRAFT    = "draft"
	BLOCK_STATUS_ACTIVE   = "active"
	BLOCK_STATUS_INACTIVE = "inactive"
)

// Error Messages for Validation
const (
	ERROR_EMPTY_ARRAY     = "array cannot be empty"
	ERROR_EMPTY_STRING    = "string cannot be empty"
	ERROR_NEGATIVE_NUMBER = "number cannot be negative"
)

// Column Names for Database Queries
const (
	COLUMN_ALIAS              = "alias"
	COLUMN_CANONICAL_URL      = "canonical_url"
	COLUMN_CONTENT            = "content"
	COLUMN_CREATED_AT         = "created_at"
	COLUMN_DOMAIN_NAMES       = "domain_names"
	COLUMN_EDITOR             = "editor"
	COLUMN_ID                 = "id"
	COLUMN_HANDLE             = "handle"
	COLUMN_MEMO               = "memo"
	COLUMN_MENU_ID            = "menu_id"
	COLUMN_META_DESCRIPTION   = "meta_description"
	COLUMN_META_KEYWORDS      = "meta_keywords"
	COLUMN_META_ROBOTS        = "meta_robots"
	COLUMN_METAS              = "metas"
	COLUMN_NAME               = "name"
	COLUMN_MIDDLEWARES_BEFORE = "middlewares_before"
	COLUMN_MIDDLEWARES_AFTER  = "middlewares_after"
	COLUMN_PAGE_ID            = "page_id"
	COLUMN_PARENT_ID          = "parent_id"
	COLUMN_SEQUENCE           = "sequence"
	COLUMN_SITE_ID            = "site_id"
	COLUMN_SOFT_DELETED_AT    = "soft_deleted_at"
	COLUMN_STATUS             = "status"
	COLUMN_TARGET             = "target"
	COLUMN_TYPE               = "type"
	COLUMN_TEMPLATE_ID        = "template_id"
	COLUMN_TITLE              = "title"
	COLUMN_UPDATED_AT         = "updated_at"
	COLUMN_URL                = "url"
)

// Menu Statuses
const (
	MENU_STATUS_DRAFT    = "draft"
	MENU_STATUS_ACTIVE   = "active"
	MENU_STATUS_INACTIVE = "inactive"
)

// Menu Item Statuses
const (
	MENU_ITEM_STATUS_DRAFT    = "draft"
	MENU_ITEM_STATUS_ACTIVE   = "active"
	MENU_ITEM_STATUS_INACTIVE = "inactive"
)

// Middleware Types
const (
	MIDDLEWARE_TYPE_BEFORE = "before"
	MIDDLEWARE_TYPE_AFTER  = "after"
)

// Page Statuses
const (
	PAGE_STATUS_DRAFT    = "draft"
	PAGE_STATUS_ACTIVE   = "active"
	PAGE_STATUS_INACTIVE = "inactive"
)

// Page Editor Types
const (
	PAGE_EDITOR_BLOCKAREA   = "blockarea"
	PAGE_EDITOR_BLOCKEDITOR = "blockeditor"
	PAGE_EDITOR_CODEMIRROR  = "codemirror"
	PAGE_EDITOR_HTMLAREA    = "htmlarea"
	PAGE_EDITOR_MARKDOWN    = "markdown"
	PAGE_EDITOR_TEXTAREA    = "textarea"
)

// Site Statuses
const (
	SITE_STATUS_DRAFT    = "draft"
	SITE_STATUS_ACTIVE   = "active"
	SITE_STATUS_INACTIVE = "inactive"
)

// Template Statuses
const (
	TEMPLATE_STATUS_DRAFT    = "draft"
	TEMPLATE_STATUS_ACTIVE   = "active"
	TEMPLATE_STATUS_INACTIVE = "inactive"
)

// Translation Statuses
const (
	TRANSLATION_STATUS_DRAFT    = "draft"
	TRANSLATION_STATUS_ACTIVE   = "active"
	TRANSLATION_STATUS_INACTIVE = "inactive"
)

// Versioning Types
const (
	VERSIONING_TYPE_BLOCK       = "block"
	VERSIONING_TYPE_PAGE        = "page"
	VERSIONING_TYPE_TEMPLATE    = "template"
	VERSIONING_TYPE_TRANSLATION = "translation"
	VERSIONING_TYPE_SITE        = "site"
)

// Query Parameter Keys
const (
	propertyKeyColumns            = "columns"
	propertyKeyCreatedAtGte       = "created_at_gte"
	propertyKeyCreatedAtLte       = "created_at_lte"
	propertyKeyHandle             = "handle"
	propertyKeyId                 = "id"
	propertyKeyIdIn               = "id_in"
	propertyKeyLimit              = "limit"
	propertyKeyNameLike           = "name_like"
	propertyKeyOffset             = "offset"
	propertyKeyOrderBy            = "order_by"
	propertyKeyPageID             = "page_id"
	propertyKeyParentID           = "parent_id"
	propertyKeySequence           = "sequence"
	propertyKeySiteID             = "site_id"
	propertyKeySoftDeleteIncluded = "soft_delete_included"
	propertyKeySortOrder          = "sort_order"
	propertyKeyStatus             = "status"
	propertyKeyStatusIn           = "status_in"
	propertyKeyTemplateID         = "template_id"
	propertyKeyCountOnly          = "count_only"
	propertyKeyDomainName         = "domain_name"
	propertyKeyAlias              = "alias"
	propertyKeyAliasLike          = "alias_like"
	propertyKeyHandleOrID         = "handle_or_id"
)
