package cmsstore

import (
	"context"
	"errors"
	"log"
	"strconv"
	"strings"

	contractsorm "github.com/dracory/neat/contracts/database/orm"
	"github.com/dromara/carbon/v2"
)

// BlockCount returns the count of blocks matching the provided query options.
func (store *storeImplementation) BlockCount(ctx context.Context, options BlockQueryInterface) (int64, error) {
	if store.neatDB == nil {
		return -1, errors.New("cms store: database is nil")
	}

	q, _, err := store.blockSelectQuery(options)

	if err != nil {
		return -1, err
	}

	var count int64
	err = q.Table(store.blockTableName).Count(&count)
	return count, err
}

// BlockCreate creates a new block in the database.
func (store *storeImplementation) BlockCreate(ctx context.Context, block BlockInterface) error {
	if store.neatDB == nil {
		return errors.New("blockstore: database is nil") // Return an error if the database connection is not established
	}

	if block == nil {
		return errors.New("block is nil") // Return an error if the block is not provided
	}

	block.SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC)) // Set the creation timestamp of the block
	block.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC)) // Set the update timestamp of the block

	return store.withTransaction(ctx, func(txCtx context.Context) error {
		data := block.Data() // Get the data from the block to be inserted

		if store.debugEnabled {
			log.Println("BlockCreate:", data)
		}

		err := store.neatDB.Query().Table(store.blockTableName).Create(data)

		if err != nil {
			return err // Return the error if the query execution failed
		}

		block.MarkAsNotDirty() // Mark the block as not dirty after successful insertion

		return store.versioningTrackEntity(txCtx, VERSIONING_TYPE_BLOCK, block.ID(), block)
	})
}

// BlockDelete deletes a block from the database by its ID.
func (store *storeImplementation) BlockDelete(ctx context.Context, block BlockInterface) error {
	if store.neatDB == nil {
		return errors.New("blockstore: database is nil")
	}

	if block == nil {
		return errors.New("block is nil")
	}

	return store.BlockDeleteByID(ctx, block.ID())
}

// BlockDeleteByID deletes a block from the database by its ID.
func (store *storeImplementation) BlockDeleteByID(ctx context.Context, id string) error {
	if store.neatDB == nil {
		return errors.New("blockstore: database is nil") // Return an error if the database connection is not established
	}

	if id == "" {
		return errors.New("block id is empty") // Return an error if the block ID is empty
	}

	if store.debugEnabled {
		log.Println("BlockDeleteByID:", id)
	}

	_, err := store.neatDB.Query().Table(store.blockTableName).Where("id = ?", id).Delete()

	return err // Return the error if the query execution failed
}

// BlockFindByHandle finds a block by its handle (unique identifier).
func (store *storeImplementation) BlockFindByHandle(ctx context.Context, handle string) (block BlockInterface, err error) {
	if store.neatDB == nil {
		return nil, errors.New("blockstore: database is nil")
	}

	if handle == "" {
		return nil, errors.New("block handle is empty")
	}

	list, err := store.BlockList(ctx, BlockQuery().SetHandle(handle).SetLimit(1))

	if err != nil {
		return nil, err
	}

	if len(list) > 0 {
		return list[0], nil
	}

	return nil, nil
}

// BlockFindByID finds a block by its ID.
func (store *storeImplementation) BlockFindByID(ctx context.Context, id string) (block BlockInterface, err error) {
	if store.neatDB == nil {
		return nil, errors.New("blockstore: database is nil")
	}

	if id == "" {
		return nil, errors.New("block id is empty") // Return an error if the block ID is empty
	}

	// Normalize ID to lowercase for consistent lookups
	id = NormalizeID(id)

	// Try direct lookup first (handles both 9-char and 32-char IDs)
	list, err := store.BlockList(ctx, BlockQuery().SetID(id).SetLimit(1)) // Get the list of blocks matching the ID

	if err != nil {
		return nil, err // Return the error if the query execution failed
	}

	if len(list) > 0 {
		return list[0], nil // Return the first block if found
	}

	// If not found and ID looks shortened, try unshortening
	if IsShortID(id) {
		unshortenedID := UnshortenID(id)
		if unshortenedID != id {
			list, err = store.BlockList(ctx, BlockQuery().SetID(unshortenedID).SetLimit(1))
			if err != nil {
				return nil, err
			}
			if len(list) > 0 {
				return list[0], nil
			}
		}
	}

	return nil, nil // Return nil if no block is found
}

func (store *storeImplementation) BlockList(ctx context.Context, query BlockQueryInterface) ([]BlockInterface, error) {
	if store.neatDB == nil {
		return []BlockInterface{}, errors.New("blockstore: database is nil")
	}

	if query == nil {
		return []BlockInterface{}, nil
	}

	q, _, err := store.blockSelectQuery(query)

	if err != nil {
		return []BlockInterface{}, err
	}

	type blockRow struct {
		ID            string `db:"id"`
		SiteID        string `db:"site_id"`
		PageID        string `db:"page_id"`
		TemplateID    string `db:"template_id"`
		ParentID      string `db:"parent_id"`
		Name          string `db:"name"`
		Handle        string `db:"handle"`
		Type          string `db:"type"`
		Content       string `db:"content"`
		Sequence      int    `db:"sequence"`
		Status        string `db:"status"`
		CreatedAt     string `db:"created_at"`
		UpdatedAt     string `db:"updated_at"`
		SoftDeletedAt string `db:"soft_deleted_at"`
	}

	var rows []blockRow
	if err := q.Table(store.blockTableName).Get(&rows); err != nil {
		return []BlockInterface{}, err
	}

	list := make([]BlockInterface, 0, len(rows))
	for _, r := range rows {
		modelMap := map[string]string{
			"id":              r.ID,
			"site_id":         r.SiteID,
			"page_id":         r.PageID,
			"template_id":     r.TemplateID,
			"parent_id":       r.ParentID,
			"name":            r.Name,
			"handle":          r.Handle,
			"type":            r.Type,
			"content":         r.Content,
			"sequence":        strconv.Itoa(r.Sequence),
			"status":          r.Status,
			"created_at":      r.CreatedAt,
			"updated_at":      r.UpdatedAt,
			"soft_deleted_at": r.SoftDeletedAt,
		}
		model := NewBlockFromExistingData(modelMap)
		list = append(list, model)
	}

	return list, nil
}

func (store *storeImplementation) BlockSoftDelete(ctx context.Context, block BlockInterface) error {
	if store.neatDB == nil {
		return errors.New("blockstore: database is nil")
	}

	if block == nil {
		return errors.New("block is nil")
	}

	block.SetSoftDeletedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))

	return store.BlockUpdate(ctx, block)
}

func (store *storeImplementation) BlockSoftDeleteByID(ctx context.Context, id string) error {
	if store.neatDB == nil {
		return errors.New("blockstore: database is nil")
	}

	block, err := store.BlockFindByID(ctx, id)

	if err != nil {
		return err
	}

	if block == nil {
		return errors.New("block not found")
	}

	return store.BlockSoftDelete(ctx, block)
}

func (store *storeImplementation) BlockUpdate(ctx context.Context, block BlockInterface) error {
	if store.neatDB == nil {
		return errors.New("blockstore: database is nil")
	}

	if block == nil {
		return errors.New("block is nil")
	}

	block.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString())

	return store.withTransaction(ctx, func(txCtx context.Context) error {
		dataChanged := block.DataChanged()

		delete(dataChanged, COLUMN_ID) // ID is not updateable

		if len(dataChanged) < 1 {
			return nil
		}

		if store.debugEnabled {
			log.Println("BlockUpdate:", dataChanged)
		}

		_, err := store.neatDB.Query().Table(store.blockTableName).Where("id = ?", block.ID()).Update(dataChanged)
		if err != nil {
			return err
		}

		block.MarkAsNotDirty()

		return store.versioningTrackEntity(txCtx, VERSIONING_TYPE_BLOCK, block.ID(), block)
	})
}

func (store *storeImplementation) blockSelectQuery(options BlockQueryInterface) (query contractsorm.Query, columns []any, err error) {
	if options == nil {
		return nil, []any{}, errors.New("block query: cannot be nil")
	}

	if err := options.Validate(); err != nil {
		return nil, []any{}, err
	}

	q := store.neatDB.Query().Table(store.blockTableName)

	if options.HasCreatedAtGte() && options.HasCreatedAtLte() {
		q = q.Where(COLUMN_CREATED_AT+" >= ? AND "+COLUMN_CREATED_AT+" <= ?", options.CreatedAtGte(), options.CreatedAtLte())
	} else if options.HasCreatedAtGte() {
		q = q.Where(COLUMN_CREATED_AT+" >= ?", options.CreatedAtGte())
	} else if options.HasCreatedAtLte() {
		q = q.Where(COLUMN_CREATED_AT+" <= ?", options.CreatedAtLte())
	}

	if options.HasHandle() {
		q = q.Where(COLUMN_HANDLE+" = ?", options.Handle())
	}

	if options.HasID() {
		q = q.Where(COLUMN_ID+" = ?", options.ID())
	}

	if options.HasIDIn() {
		idIn := options.IDIn()
		if len(idIn) > 0 {
			placeholders := make([]string, len(idIn))
			args := make([]any, len(idIn))
			for i, v := range idIn {
				placeholders[i] = "?"
				args[i] = v
			}
			q = q.Where(COLUMN_ID+" IN ("+strings.Join(placeholders, ", ")+")", args...)
		}
	}

	if options.HasNameLike() {
		q = q.Where(COLUMN_NAME+" LIKE ?", "%"+options.NameLike()+"%")
	}

	if options.HasPageID() {
		q = q.Where(COLUMN_PAGE_ID+" = ?", options.PageID())
	}

	if options.HasParentID() {
		q = q.Where(COLUMN_PARENT_ID+" = ?", options.ParentID())
	}

	if options.HasSequence() {
		q = q.Where(COLUMN_SEQUENCE+" = ?", options.Sequence())
	}

	if options.HasSiteID() {
		q = q.Where(COLUMN_SITE_ID+" = ?", options.SiteID())
	}

	if options.HasStatus() {
		q = q.Where(COLUMN_STATUS+" = ?", options.Status())
	}

	if options.HasStatusIn() {
		statusIn := options.StatusIn()
		if len(statusIn) > 0 {
			placeholders := make([]string, len(statusIn))
			args := make([]any, len(statusIn))
			for i, v := range statusIn {
				placeholders[i] = "?"
				args[i] = v
			}
			q = q.Where(COLUMN_STATUS+" IN ("+strings.Join(placeholders, ", ")+")", args...)
		}
	}

	if options.HasTemplateID() {
		q = q.Where(COLUMN_TEMPLATE_ID+" = ?", options.TemplateID())
	}

	if options.HasType() {
		q = q.Where(COLUMN_TYPE+" = ?", options.Type())
	}

	if !options.IsCountOnly() {
		if options.HasLimit() {
			q = q.Limit(options.Limit())
		}

		if options.HasOffset() {
			q = q.Offset(options.Offset())
		}
	}

	sortOrder := SORT_ORDER_DESC
	if options.HasSortOrder() && options.SortOrder() != "" {
		sortOrder = options.SortOrder()
	}

	if !options.IsCountOnly() && options.HasOrderBy() && options.OrderBy() != "" {
		if strings.EqualFold(sortOrder, SORT_ORDER_ASC) {
			q = q.OrderBy(options.OrderBy(), "ASC")
		} else {
			q = q.OrderBy(options.OrderBy(), "DESC")
		}
	}

	columns = []any{}

	for _, column := range options.Columns() {
		columns = append(columns, column)
	}

	if options.SoftDeleteIncluded() {
		return q, columns, nil // soft deleted blocks requested specifically
	}

	q = q.Where(COLUMN_SOFT_DELETED_AT+" > ?", carbon.Now(carbon.UTC).ToDateTimeString())

	return q, columns, nil
}
