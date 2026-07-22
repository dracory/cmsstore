# Media Storage Support

**Date**: 2026-07-22
**Status**: [Approved]

## Summary

Add media (file attachment) storage support to `cmsstore`, mirroring the proven implementation from `blogstore/store_media.go`. This enables attaching files (images, videos, documents) to any CMS entity (pages, blocks, sites, templates, etc.) with full CRUD, soft-delete, query filtering, and metadata support.

### Problem Being Solved

Currently, `cmsstore` has no way to manage media files attached to CMS entities. The `blogstore` package already has a working media implementation, but `cmsstore` cannot use it because:

1. `blogstore` media uses ORM struct-embedded traits (`orm.ShortID`, `orm.CreatedAt`, etc.) while `cmsstore` uses `dataobject.DataObject` (map-based storage)
2. `blogstore` uses a plain struct `MediaQueryOptions` while `cmsstore` uses interface-based query builders (`PageQueryInterface`, `BlockQueryInterface`, etc.) with `Has*()` / `Validate()` pattern
3. `cmsstore` has versioning integration (`withTransaction`, `versioningTrackEntity`, `MarshalToVersioning`) that `blogstore` media does not participate in
4. `cmsstore` uses `NormalizeID()` / `IsShortID()` / `UnshortenID()` for ID handling, while `blogstore` uses `GenerateShortID()` directly

### High-Level Solution

Port the media concept from `blogstore` and adapt it to `cmsstore` conventions:

- **`MediaInterface`** — follows the cmsstore interface style (e.g. `PageInterface`, `BlockInterface`) with `Data()`, `DataChanged()`, `MarkAsNotDirty()`, `MarshalToVersioning()`, and `dataobject.DataObject` base
- **`MediaQueryInterface`** — interface-based query builder with `Has*()` / `Validate()` pattern, matching `PageQueryInterface`
- **`store_media.go`** — store CRUD methods following `store_pages.go` patterns including `withTransaction` and versioning integration
- **Database migration** — media table creation in `MigrateUp` / `MigrateDown`
- **Store configuration** — `MediaEnabled` flag and `MediaTableName` in `NewStoreOptions`

## Background

### Current System State

`cmsstore` manages blocks, pages, sites, templates, menus, menu items, translations, and versioning. Each entity type follows a consistent pattern:

1. **Interface** (`interfaces.go`) — e.g. `PageInterface` with getters/setters, `Data()`, `DataChanged()`, `MarkAsNotDirty()`, `MarshalToVersioning()`, status predicates
2. **Model** (`page.go`) — `pageImplementation` embedding `dataobject.DataObject`, constructed via `NewPage()` / `NewPageFromExistingData()`
3. **Query Interface** (`page_query_interface.go`) — `PageQueryInterface` with `Has*()` / `Validate()` / `Set*()` methods
4. **Query Implementation** (`page_query.go`) — `pageQuery` struct with `parameters map[string]any`
5. **Store Methods** (`store_pages.go`) — `PageCreate`, `PageCount`, `PageDelete`, `PageDeleteByID`, `PageFindByID`, `PageFindByHandle`, `PageList`, `PageSoftDelete`, `PageSoftDeleteByID`, `PageUpdate` + private `pageSelectQuery`
6. **Store Config** (`store_new.go`) — table name in `NewStoreOptions`, field in `storeImplementation`, validation in `NewStore()`
7. **Migration** (`store.go`) — table creation in `MigrateUp`, drop in `MigrateDown`
8. **Constants** (`consts.go`) — column names, status constants

### Reference Implementation (blogstore)

The `blogstore` media implementation consists of:

- **`media.go`** — `MediaInterface` (131 lines interface) + `mediaImplementation` (struct with ORM traits) + `NewMedia()` constructor
- **`media_query.go`** — `MediaQueryOptions` plain struct (34 lines)
- **`store_media.go`** — 10 store methods + `buildMediaQuery` (394 lines total)
- **`constants.go`** — 4 media column constants + 3 status constants
- **`store.go`** — `mediaTableName` field, `MediaCreate`/`MediaCount`/etc. in `StoreInterface`, media table migration

### Why This Change Is Needed

1. CMS pages and blocks often need associated images (hero images, thumbnails, gallery images)
2. Sites need logos and favicon media
3. Templates may reference media assets
4. A generic entity-attached media system is more flexible than adding image columns to each entity table
5. The `entity_id` + `entity_type` pattern allows media to be attached to any entity, not just one specific type

## Detailed Design

### 1. Constants (`consts.go`)

Add media-specific column and status constants:

```go
// Media columns
const (
    COLUMN_MEDIA_URL      = "media_url"
    COLUMN_MEDIA_TYPE     = "media_type"
    COLUMN_FILE_SIZE      = "file_size"
    COLUMN_FILE_EXTENSION = "file_extension"
)

// Media Statuses
const (
    MEDIA_STATUS_DRAFT    = "draft"
    MEDIA_STATUS_ACTIVE   = "active"
    MEDIA_STATUS_INACTIVE = "inactive"
)

// Versioning Types
const VERSIONING_TYPE_MEDIA = "media"
```

Note: `COLUMN_ENTITY_ID`, `COLUMN_ENTITY_TYPE`, `COLUMN_ID`, `COLUMN_TITLE`, `COLUMN_DESCRIPTION`, `COLUMN_MEMO`, `COLUMN_SEQUENCE`, `COLUMN_STATUS`, `COLUMN_METAS`, `COLUMN_MEMO`, `COLUMN_CREATED_AT`, `COLUMN_UPDATED_AT`, `COLUMN_SOFT_DELETED_AT` already exist in `consts.go`.

### 2. Media Interface (`interfaces.go`)

Add `MediaInterface` following the cmsstore interface conventions:

```go
type MediaInterface interface {
    Data() map[string]string
    DataChanged() map[string]string
    MarkAsNotDirty(...string)

    // Methods
    MarshalToVersioning() (string, error)

    // Setters and Getters
    ID() string
    SetID(id string) MediaInterface

    EntityID() string
    SetEntityID(entityID string) MediaInterface

    EntityType() string
    SetEntityType(entityType string) MediaInterface

    Title() string
    SetTitle(title string) MediaInterface

    Description() string
    SetDescription(description string) MediaInterface

    Memo() string
    SetMemo(memo string) MediaInterface

    URL() string
    SetURL(url string) MediaInterface

    Type() string
    SetType(mediaType string) MediaInterface

    Size() string
    SetSize(size string) MediaInterface

    Extension() string
    SetExtension(extension string) MediaInterface

    Sequence() string
    SequenceInt() int
    SetSequence(sequence string) MediaInterface
    SetSequenceInt(sequence int) MediaInterface

    Status() string
    SetStatus(status string) MediaInterface

    Handle() string
    SetHandle(handle string) MediaInterface

    SiteID() string
    SetSiteID(siteID string) MediaInterface

    // Metadata
    Meta(key string) string
    SetMeta(key, value string) error
    Metas() (map[string]string, error)
    SetMetas(metas map[string]string) error
    UpsertMetas(metas map[string]string) error

    // Timestamps
    CreatedAt() string
    SetCreatedAt(createdAt string) MediaInterface
    CreatedAtCarbon() *carbon.Carbon

    UpdatedAt() string
    SetUpdatedAt(updatedAt string) MediaInterface
    UpdatedAtCarbon() *carbon.Carbon

    SoftDeletedAt() string
    SetSoftDeletedAt(softDeletedAt string) MediaInterface
    SoftDeletedAtCarbon() *carbon.Carbon

    // Predicates
    IsActive() bool
    IsInactive() bool
    IsDraft() bool
    IsSoftDeleted() bool

    // Type Predicates
    IsImage() bool
    IsVideo() bool
}
```

**Key differences from blogstore's `MediaInterface`:**
- Uses `Data()` / `DataChanged()` / `MarkAsNotDirty()` (cmsstore dirty-tracking pattern)
- Uses `MarshalToVersioning()` (cmsstore versioning pattern)
- Adds `EntityType()` / `SetEntityType()` — blogstore only has `EntityID`, but cmsstore has multiple entity types (page, block, site, template) so `EntityType` is needed to distinguish
- Adds `Handle()` / `SetHandle()` — consistent with all cmsstore entities
- Adds `SiteID()` / `SetSiteID()` — for site-scoped media queries
- Uses `Sequence()` returning `string` + `SequenceInt()` returning `int` (cmsstore pattern) instead of just `int`
- Method names without `Get` prefix (cmsstore convention: `ID()` not `GetID()`)

### 3. Media Model (`media.go`)

```go
package cmsstore

import (
    "strings"
    "github.com/dracory/dataobject"
    "github.com/dromara/carbon/v2"
)

type mediaImplementation struct {
    dataobject.DataObject
}

var _ MediaInterface = (*mediaImplementation)(nil)

func NewMedia() MediaInterface {
    o := &mediaImplementation{}
    o.SetID(GenerateShortID())
    o.SetEntityID("")
    o.SetEntityType("")
    o.SetTitle("")
    o.SetDescription("")
    o.SetMemo("")
    o.SetURL("")
    o.SetType("")
    o.SetSize("0")
    o.SetExtension("")
    o.SetSequence("0")
    o.SetStatus(MEDIA_STATUS_DRAFT)
    o.SetHandle("")
    o.SetSiteID("")
    o.SetMetas(map[string]string{})
    o.SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
    o.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
    o.SetSoftDeletedAt(MAX_DATETIME)
    return o
}

func NewMediaFromExistingData(data map[string]string) *mediaImplementation {
    o := &mediaImplementation{}
    o.Hydrate(data)
    return o
}

func (o *mediaImplementation) IsActive() bool {
    return o.Status() == MEDIA_STATUS_ACTIVE
}

func (o *mediaImplementation) IsInactive() bool {
    return o.Status() == MEDIA_STATUS_INACTIVE
}

func (o *mediaImplementation) IsDraft() bool {
    return o.Status() == MEDIA_STATUS_DRAFT
}

func (o *mediaImplementation) IsSoftDeleted() bool {
    return o.SoftDeletedAtCarbon().Compare("<", carbon.Now(carbon.UTC))
}

func (o *mediaImplementation) IsImage() bool {
    return strings.HasPrefix(o.Type(), "image/")
}

func (o *mediaImplementation) IsVideo() bool {
    return strings.HasPrefix(o.Type(), "video/")
}

func (o *mediaImplementation) MarshalToVersioning() (string, error) {
    versionedData := map[string]string{}
    for k, v := range o.Data() {
        if k == COLUMN_CREATED_AT ||
            k == COLUMN_UPDATED_AT ||
            k == COLUMN_SOFT_DELETED_AT {
            continue
        }
        versionedData[k] = v
    }
    // Marshal to JSON...
    // (Follow the same pattern as pageImplementation.MarshalToVersioning)
}
```

The `dataobject.DataObject` base provides `Data()`, `DataChanged()`, `MarkAsNotDirty()`, `Hydrate()`, and all the getter/setter methods via its internal map. This is the same pattern used by `pageImplementation`, `blockImplementation`, etc.

### 4. Media Query Interface (`media_query_interface.go`)

```go
type MediaQueryInterface interface {
    Validate() error

    Columns() []string
    HasColumns() bool
    SetColumns(columns []string) MediaQueryInterface

    HasID() bool
    ID() string
    SetID(id string) MediaQueryInterface

    HasIDIn() bool
    IDIn() []string
    SetIDIn(idIn []string) MediaQueryInterface

    HasEntityID() bool
    EntityID() string
    SetEntityID(entityID string) MediaQueryInterface

    HasEntityType() bool
    EntityType() string
    SetEntityType(entityType string) MediaQueryInterface

    HasSiteID() bool
    SiteID() string
    SetSiteID(siteID string) MediaQueryInterface

    HasHandle() bool
    Handle() string
    SetHandle(handle string) MediaQueryInterface

    HasExtension() bool
    Extension() string
    SetExtension(extension string) MediaQueryInterface

    HasType() bool
    Type() string
    SetType(mediaType string) MediaQueryInterface

    HasStatus() bool
    Status() string
    SetStatus(status string) MediaQueryInterface

    HasStatusIn() bool
    StatusIn() []string
    SetStatusIn(statusIn []string) MediaQueryInterface

    HasNameLike() bool
    NameLike() string
    SetNameLike(nameLike string) MediaQueryInterface

    HasCountOnly() bool
    IsCountOnly() bool
    SetCountOnly(countOnly bool) MediaQueryInterface

    HasLimit() bool
    Limit() int
    SetLimit(limit int) MediaQueryInterface

    HasOffset() bool
    Offset() int
    SetOffset(offset int) MediaQueryInterface

    HasSortOrder() bool
    SortOrder() string
    SetSortOrder(sortOrder string) MediaQueryInterface

    HasOrderBy() bool
    OrderBy() string
    SetOrderBy(orderBy string) MediaQueryInterface

    HasSoftDeletedIncluded() bool
    SoftDeletedIncluded() bool
    SetSoftDeletedIncluded(softDeleteIncluded bool) MediaQueryInterface
}

func MediaQuery() MediaQueryInterface {
    return &mediaQuery{
        parameters: make(map[string]any),
    }
}
```

**Key additions over blogstore's `MediaQueryOptions`:**
- `EntityType` filter — allows querying media for a specific entity type (e.g. all media for pages vs blocks)
- `SiteID` filter — allows querying media within a specific site
- `Handle` filter — consistent with other cmsstore query interfaces
- `StatusIn` filter — batch status filtering (consistent with `PageQueryInterface`)
- `NameLike` filter — replaces blogstore's `Search` field with the cmsstore naming convention
- Interface-based with `Has*()` methods — enables nil/empty distinction for query building

### 5. Media Query Implementation (`media_query.go`)

Follows the `page_query.go` pattern exactly — a `mediaQuery` struct with `parameters map[string]any`, implementing all `MediaQueryInterface` methods with `Has*()` checks and `Validate()`.

### 6. Store Methods (`store_media.go`)

Follows the `store_pages.go` pattern:

```go
func (store *storeImplementation) MediaCreate(ctx context.Context, media MediaInterface) error
func (store *storeImplementation) MediaCount(ctx context.Context, options MediaQueryInterface) (int64, error)
func (store *storeImplementation) MediaDelete(ctx context.Context, media MediaInterface) error
func (store *storeImplementation) MediaDeleteByID(ctx context.Context, id string) error
func (store *storeImplementation) MediaFindByID(ctx context.Context, id string) (MediaInterface, error)
func (store *storeImplementation) MediaFindByHandle(ctx context.Context, handle string) (MediaInterface, error)
func (store *storeImplementation) MediaList(ctx context.Context, query MediaQueryInterface) ([]MediaInterface, error)
func (store *mediaImplementation) MediaListByEntityID(ctx context.Context, entityID string, entityType string) ([]MediaInterface, error)
func (store *storeImplementation) MediaSoftDelete(ctx context.Context, media MediaInterface) error
func (store *storeImplementation) MediaSoftDeleteByID(ctx context.Context, id string) error
func (store *storeImplementation) MediaUpdate(ctx context.Context, media MediaInterface) error
```

**Key implementation details:**

- `MediaCreate` uses `store.withTransaction()` + `store.versioningTrackEntity()` (same as `PageCreate`)
- `MediaUpdate` uses `DataChanged()` to only update modified fields + `MarkAsNotDirty()` after update (same as `PageUpdate`)
- `MediaFindByID` uses `NormalizeID()` / `IsShortID()` / `UnshortenID()` (same as `PageFindByID`)
- `MediaList` uses a private `mediaSelectQuery()` method that builds the query from `MediaQueryInterface` (same as `pageSelectQuery()`)
- `MediaListByEntityID` is a convenience method that calls `MediaList` with `EntityID` + `EntityType` filters and `OrderBy: COLUMN_SEQUENCE`, `SortOrder: SORT_ORDER_ASC`

**`mediaSelectQuery` implementation:**

```go
func (store *storeImplementation) mediaSelectQuery(options MediaQueryInterface) (query contractsorm.Query, selectColumns []any, err error) {
    if options == nil {
        return nil, []any{}, errors.New("media options cannot be nil")
    }

    if err := options.Validate(); err != nil {
        return nil, []any{}, err
    }

    q := store.neatDB.Query().Table(store.mediaTableName)

    if options.HasID() {
        q = q.Where(COLUMN_ID+" = ?", options.ID())
    }

    if options.HasIDIn() {
        // Build IN clause (same pattern as pageSelectQuery)
    }

    if options.HasEntityID() {
        q = q.Where(COLUMN_ENTITY_ID+" = ?", options.EntityID())
    }

    if options.HasEntityType() {
        q = q.Where(COLUMN_ENTITY_TYPE+" = ?", options.EntityType())
    }

    if options.HasSiteID() {
        q = q.Where(COLUMN_SITE_ID+" = ?", options.SiteID())
    }

    if options.HasHandle() {
        q = q.Where(COLUMN_HANDLE+" = ?", options.Handle())
    }

    if options.HasExtension() {
        q = q.Where(COLUMN_FILE_EXTENSION+" = ?", options.Extension())
    }

    if options.HasType() {
        q = q.Where(COLUMN_MEDIA_TYPE+" = ?", options.Type())
    }

    if options.HasStatus() {
        q = q.Where(COLUMN_STATUS+" = ?", options.Status())
    }

    if options.HasStatusIn() {
        // Build IN clause (same pattern as pageSelectQuery)
    }

    if options.HasNameLike() {
        q = q.Where(COLUMN_NAME+" LIKE ?", options.NameLike())
    }

    // Pagination, sorting, soft-delete filtering
    // (Same pattern as pageSelectQuery)

    if !options.IsCountOnly() {
        if options.HasLimit() {
            q = q.Limit(options.Limit())
        }
        if options.HasOffset() {
            q = q.Offset(options.Offset())
        }
    }

    sortOrder := SORT_ORDER_DESC
    if options.HasSortOrder() {
        sortOrder = options.SortOrder()
    }

    if !options.IsCountOnly() && options.HasOrderBy() {
        if strings.EqualFold(sortOrder, SORT_ORDER_ASC) {
            q = q.OrderBy(options.OrderBy(), "ASC")
        } else {
            q = q.OrderBy(options.OrderBy(), "DESC")
        }
    }

    if options.SoftDeletedIncluded() {
        return q, []any{}, nil
    }

    q = q.Where(COLUMN_SOFT_DELETED_AT+" > ?", carbon.Now(carbon.UTC).ToDateTimeString())

    return q, []any{}, nil
}
```

### 7. Store Configuration (`store_new.go`)

Add to `NewStoreOptions`:

```go
type NewStoreOptions struct {
    // ... existing fields ...

    // MediaEnabled enables media support
    MediaEnabled bool

    // MediaTableName is the name of the media database table to be created/used
    MediaTableName string
}
```

Add validation in `NewStore()`:

```go
if opts.MediaEnabled && opts.MediaTableName == "" {
    return nil, errors.New("cms store: MediaTableName is required")
}
```

Add to `storeImplementation` in `store.go`:

```go
type storeImplementation struct {
    // ... existing fields ...

    // Media
    mediaEnabled  bool
    mediaTableName string
}
```

Add to `NewStore()` store initialization:

```go
store := &storeImplementation{
    // ... existing fields ...

    mediaEnabled:   opts.MediaEnabled,
    mediaTableName: opts.MediaTableName,
}
```

Add `MediaEnabled()` method to `StoreInterface`:

```go
// In StoreInterface
MediaEnabled() bool
```

### 8. Store Interface (`interfaces.go`)

Add media methods to `StoreInterface`:

```go
type StoreInterface interface {
    // ... existing methods ...

    // Media
    MediaEnabled() bool
    MediaCreate(ctx context.Context, media MediaInterface) error
    MediaCount(ctx context.Context, options MediaQueryInterface) (int64, error)
    MediaDelete(ctx context.Context, media MediaInterface) error
    MediaDeleteByID(ctx context.Context, id string) error
    MediaFindByID(ctx context.Context, id string) (MediaInterface, error)
    MediaFindByHandle(ctx context.Context, handle string) (MediaInterface, error)
    MediaList(ctx context.Context, query MediaQueryInterface) ([]MediaInterface, error)
    MediaListByEntityID(ctx context.Context, entityID string, entityType string) ([]MediaInterface, error)
    MediaSoftDelete(ctx context.Context, media MediaInterface) error
    MediaSoftDeleteByID(ctx context.Context, id string) error
    MediaUpdate(ctx context.Context, media MediaInterface) error
}
```

### 9. Database Migration (`store.go`)

Add media table creation in `MigrateUp`:

```go
// Create media table if enabled
if store.mediaEnabled {
    if !store.neatDB.Schema().HasTable(store.mediaTableName) {
        err := store.neatDB.Schema().Create(store.mediaTableName, func(table contractsschema.Blueprint) {
            table.String(COLUMN_ID, 40)
            table.Primary(COLUMN_ID)
            table.String(COLUMN_ENTITY_ID, 40)
            table.String(COLUMN_ENTITY_TYPE, 40)
            table.String(COLUMN_SITE_ID, 40)
            table.String(COLUMN_TITLE, 255)
            table.Text(COLUMN_DESCRIPTION)
            table.String(COLUMN_MEMO, 255)
            table.LongText(COLUMN_MEDIA_URL)
            table.String(COLUMN_MEDIA_TYPE, 100)
            table.String(COLUMN_FILE_SIZE, 50).Default("0")
            table.String(COLUMN_FILE_EXTENSION, 20)
            table.Integer(COLUMN_SEQUENCE)
            table.String(COLUMN_STATUS, 40)
            table.String(COLUMN_HANDLE, 40)
            table.Text(COLUMN_METAS)
            table.DateTime(COLUMN_CREATED_AT)
            table.DateTime(COLUMN_UPDATED_AT)
            table.DateTime(COLUMN_SOFT_DELETED_AT)
        })
        if err != nil {
            return err
        }
    }
}
```

Add media table drop in `MigrateDown`:

```go
if store.mediaEnabled {
    if store.neatDB.Schema().HasTable(store.mediaTableName) {
        err := store.neatDB.Schema().Drop(store.mediaTableName)
        if err != nil {
            return err
        }
    }
}
```

### 10. Versioning Integration

Add `VERSIONING_TYPE_MEDIA = "media"` to `consts.go`.

The `MarshalToVersioning()` method on `mediaImplementation` follows the same pattern as `pageImplementation.MarshalToVersioning()` — excludes `created_at`, `updated_at`, `soft_deleted_at` from the versioned snapshot.

The `versioningTrackEntity` call in `MediaCreate` and `MediaUpdate` uses `VERSIONING_TYPE_MEDIA` as the entity type.

### 11. Tests

Create test files following the cmsstore test patterns:

**`media_test.go`** — unit tests for the media model:
- `TestNewMedia` — verifies all default values
- `TestMediaSettersAndGetters` — all setter/getter pairs
- `TestMediaTimestamps` — timestamp handling
- `TestMediaGetData` — `Data()` returns all fields
- `TestMediaIsSoftDeleted` — soft delete predicate
- `TestMediaSequence` — sequence edge cases
- `TestMediaStatusPredicates` — `IsActive()`, `IsInactive()`, `IsDraft()`
- `TestMediaTypePredicates` — `IsImage()`, `IsVideo()`
- `TestMediaMetas` — metadata operations (`SetMetas`, `GetMetas`, `SetMeta`, `UpsertMetas`, `MetaRemove`)

**`media_query_test.go`** — unit tests for the query builder:
- `TestMediaQueryValidate` — validation logic
- `TestMediaQuerySettersAndGetters` — all query setter/getter pairs

**`store_media_test.go`** — integration tests for store operations:
- `TestStoreMediaCreate` — create media, verify timestamps set
- `TestStoreMediaCreateErrors` — nil media, empty entity_id
- `TestStoreMediaFindByID` — create then find by ID
- `TestStoreMediaFindByHandle` — create then find by handle
- `TestStoreMediaList` — create multiple, list all, list by entity
- `TestStoreMediaListByEntityID` — create multiple for same entity, verify ordering by sequence
- `TestStoreMediaUpdate` — create then update, verify changes persisted
- `TestStoreMediaSoftDelete` — create then soft delete, verify excluded from default list
- `TestStoreMediaSoftDeleteByID` — create then soft delete by ID
- `TestStoreMediaDelete` — hard delete verification
- `TestStoreMediaCount` — count with various filters

## Implementation Plan

### Step 1: Constants
- Add media column constants, status constants, and `VERSIONING_TYPE_MEDIA` to `consts.go`

### Step 2: Interface
- Add `MediaInterface` to `interfaces.go`
- Add media methods to `StoreInterface` in `interfaces.go`

### Step 3: Model
- Create `media.go` with `mediaImplementation`, `NewMedia()`, `NewMediaFromExistingData()`, predicates, and `MarshalToVersioning()`

### Step 4: Query Interface + Implementation
- Create `media_query_interface.go` with `MediaQueryInterface`
- Create `media_query.go` with `mediaQuery` struct implementation

### Step 5: Store Configuration
- Add `MediaEnabled` and `MediaTableName` to `NewStoreOptions` in `store_new.go`
- Add `mediaEnabled` and `mediaTableName` fields to `storeImplementation` in `store.go`
- Add validation in `NewStore()`
- Add `MediaEnabled()` method

### Step 6: Store Methods
- Create `store_media.go` with all CRUD methods + `mediaSelectQuery`
- Include versioning integration via `withTransaction` and `versioningTrackEntity`

### Step 7: Migration
- Add media table creation to `MigrateUp` in `store.go`
- Add media table drop to `MigrateDown` in `store.go`

### Step 8: Tests
- Create `media_test.go` (model unit tests)
- Create `media_query_test.go` (query unit tests)
- Create `store_media_test.go` (store integration tests)

### Step 9: Test Utils Update
- Update `testutils/utils.go` to include `MediaEnabled: true` and `MediaTableName: "cms_media"` in test store initialization

### Dependencies
- No new external dependencies required — all needed packages (`dataobject`, `neat`, `carbon`) are already in `go.mod`

### Migration Strategy
- Media is opt-in via `MediaEnabled: true` in `NewStoreOptions`
- Existing users are unaffected — no schema changes unless they enable media
- `MigrateUp` only creates the media table if `mediaEnabled` is true

## Risks and Mitigations

### Risk 1: Versioning Overhead
**Risk:** Every media create/update writes a versioning record, which may be excessive for frequently updated media (e.g. re-uploads).
**Mitigation:** Versioning only runs if `store.versioningEnabled` is true. The `versioningTrackEntity` method already checks this. Users can disable versioning if not needed.

### Risk 2: Entity ID Collision
**Risk:** Different entity types (page, block) could have the same short ID, causing confusion when querying media by `EntityID` alone.
**Mitigation:** The `EntityType` field and `MediaListByEntityID(entityID, entityType)` signature enforce type-scoped queries. The `mediaSelectQuery` always filters by both `EntityID` and `EntityType` when both are provided.

### Risk 3: Large Media Metadata
**Risk:** The `metas` column is `TEXT` which may be insufficient for very large metadata payloads.
**Mitigation:** `TEXT` in SQLite can hold up to 1 billion bytes. For MySQL/PostgreSQL, `TEXT` is also sufficient for metadata. If needed, this can be changed to `LONGTEXT` in the migration.

### Risk 4: Performance with Many Media per Entity
**Risk:** Entities with hundreds of media files may cause slow queries.
**Mitigation:** The query interface supports `Limit` / `Offset` pagination. The `sequence` column enables ordered retrieval. Indexes on `(entity_id, entity_type)` and `(site_id)` can be added if performance becomes an issue.
