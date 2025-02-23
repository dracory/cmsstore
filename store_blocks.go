package cmsstore

import (
	"context"
	"errors"
	"log"
	"strconv"
	"strings"

	"github.com/doug-martin/goqu/v9"
	"github.com/dromara/carbon/v2"
	"github.com/gouniverse/base/database"
	"github.com/gouniverse/sb"
	"github.com/samber/lo"
)

// BlockCount returns the count of blocks matching the provided query options.
func (store *store) BlockCount(ctx context.Context, options BlockQueryInterface) (int64, error) {
	if store.db == nil {
		return -1, errors.New("cms store: db is nil") // Return an error if the database connection is not established
	}

	options.SetCountOnly(true) // Set the query to count only

	q, _, err := store.blockSelectQuery(options) // Generate the select query for the block count

	if err != nil {
		return -1, err // Return the error if the query generation failed
	}

	sqlStr, params, errSql := q.Prepared(true).
		Limit(1).                                    // Limit the query to 1 result
		Select(goqu.COUNT(goqu.Star()).As("count")). // Select the count of all matching blocks
		ToSQL()                                      // Convert the query to SQL string and parameters

	if errSql != nil {
		return -1, nil // Return an error if the SQL conversion failed
	}

	if store.debugEnabled {
		log.Println(sqlStr) // Log the SQL query if debug mode is enabled
	}

	mapped, err := database.SelectToMapString(store.toQuerableContext(ctx), sqlStr, params...) // Execute the SQL query and get the result

	if err != nil {
		return -1, err // Return the error if the query execution failed
	}

	if len(mapped) < 1 {
		return -1, nil // Return an error if no results were found
	}

	countStr := mapped[0]["count"] // Extract the count string from the result

	i, err := strconv.ParseInt(countStr, 10, 64) // Convert the count string to an integer

	if err != nil {
		return -1, err // Return the error if the conversion failed
	}

	return i, nil // Return the count of blocks
}

// BlockCreate creates a new block in the database.
func (store *store) BlockCreate(ctx context.Context, block BlockInterface) error {
	if store.db == nil {
		return errors.New("blockstore: database is nil") // Return an error if the database connection is not established
	}

	if block == nil {
		return errors.New("block is nil") // Return an error if the block is not provided
	}

	block.SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC)) // Set the creation timestamp of the block
	block.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC)) // Set the update timestamp of the block

	data := block.Data() // Get the data from the block to be inserted

	sqlStr, params, errSql := goqu.Dialect(store.dbDriverName).
		Insert(store.blockTableName). // Insert into the block table
		Prepared(true).
		Rows(data). // Insert the block data
		ToSQL()     // Convert the query to SQL string and parameters

	if errSql != nil {
		return errSql // Return the error if the SQL conversion failed
	}

	if store.debugEnabled {
		log.Println(sqlStr) // Log the SQL query if debug mode is enabled
	}

	_, err := database.Execute(store.toQuerableContext(ctx), sqlStr, params...) // Execute the SQL query

	if err != nil {
		return err // Return the error if the query execution failed
	}

	block.MarkAsNotDirty() // Mark the block as not dirty after successful insertion

	return nil // Return success
}

// BlockDelete deletes a block from the database by its ID.
func (store *store) BlockDelete(ctx context.Context, block BlockInterface) error {
	if store.db == nil {
		return errors.New("blockstore: database is nil") // Return an error if the database connection is not established
	}

	if block == nil {
		return errors.New("block is nil") // Return an error if the block is not provided
	}

	return store.BlockDeleteByID(ctx, block.ID()) // Delete the block by its ID
}

// BlockDeleteByID deletes a block from the database by its ID.
func (store *store) BlockDeleteByID(ctx context.Context, id string) error {
	if store.db == nil {
		return errors.New("blockstore: database is nil") // Return an error if the database connection is not established
	}

	if id == "" {
		return errors.New("block id is empty") // Return an error if the block ID is empty
	}

	sqlStr, params, errSql := goqu.Dialect(store.dbDriverName).
		Delete(store.blockTableName). // Delete from the block table
		Prepared(true).
		Where(goqu.C("id").Eq(id)). // Where the block ID matches the provided ID
		ToSQL()                     // Convert the query to SQL string and parameters

	if errSql != nil {
		return errSql // Return the error if the SQL conversion failed
	}

	if store.debugEnabled {
		log.Println(sqlStr) // Log the SQL query if debug mode is enabled
	}

	_, err := database.Execute(store.toQuerableContext(ctx), sqlStr, params...) // Execute the SQL query

	return err // Return the error if the query execution failed
}

// BlockFindByHandle finds a block by its handle (unique identifier).
func (store *store) BlockFindByHandle(ctx context.Context, handle string) (block BlockInterface, err error) {
	if store.db == nil {
		return nil, errors.New("blockstore: database is nil") // Return an error if the database connection is not established
	}

	if handle == "" {
		return nil, errors.New("block handle is empty") // Return an error if the block handle is empty
	}

	list, err := store.BlockList(ctx, BlockQuery().SetHandle(handle).SetLimit(1)) // Get the list of blocks matching the handle

	if err != nil {
		return nil, err // Return the error if the query execution failed
	}

	if len(list) > 0 {
		return list[0], nil // Return the first block if found
	}

	return nil, nil // Return nil if no block is found
}

// BlockFindByID finds a block by its ID.
func (store *store) BlockFindByID(ctx context.Context, id string) (block BlockInterface, err error) {
	if store.db == nil {
		return nil, errors.New("blockstore: database is nil") // Return an error if the database connection is not established
	}

	if id == "" {
		return nil, errors.New("block id is empty") // Return an error if the block ID is empty
	}

	list, err := store.BlockList(ctx, BlockQuery().SetID(id).SetLimit(1)) // Get the list of blocks matching the ID

	if err != nil {
		return nil, err // Return the error if the query execution failed
	}

	if len(list) > 0 {
		return list[0], nil // Return the first block if found
	}

	return nil, nil // Return nil if no block is found
}

func (store *store) BlockList(ctx context.Context, query BlockQueryInterface) ([]BlockInterface, error) {
	if store.db == nil {
		return []BlockInterface{}, errors.New("blockstore: database is nil")
	}

	if query == nil {
		return []BlockInterface{}, nil
	}

	q, columns, err := store.blockSelectQuery(query)

	if err != nil {
		return []BlockInterface{}, err
	}

	sqlStr, _, errSql := q.Select(columns...).ToSQL()

	if errSql != nil {
		return []BlockInterface{}, nil
	}

	if store.debugEnabled {
		log.Println(sqlStr)
	}

	modelMaps, err := database.SelectToMapString(store.toQuerableContext(ctx), sqlStr)

	if err != nil {
		return []BlockInterface{}, err
	}

	list := []BlockInterface{}

	lo.ForEach(modelMaps, func(modelMap map[string]string, index int) {
		model := NewBlockFromExistingData(modelMap)
		list = append(list, model)
	})

	return list, nil
}

func (store *store) BlockSoftDelete(ctx context.Context, block BlockInterface) error {
	if store.db == nil {
		return errors.New("blockstore: database is nil")
	}

	if block == nil {
		return errors.New("block is nil")
	}

	block.SetSoftDeletedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))

	return store.BlockUpdate(ctx, block)
}

func (store *store) BlockSoftDeleteByID(ctx context.Context, id string) error {
	if store.db == nil {
		return errors.New("blockstore: database is nil")
	}

	block, err := store.BlockFindByID(ctx, id)

	if err != nil {
		return err
	}

	return store.BlockSoftDelete(ctx, block)
}

func (store *store) BlockUpdate(ctx context.Context, block BlockInterface) error {
	if store.db == nil {
		return errors.New("blockstore: database is nil")
	}

	if block == nil {
		return errors.New("block is nil")
	}

	block.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString())

	dataChanged := block.DataChanged()

	delete(dataChanged, COLUMN_ID) // ID is not updateable

	if len(dataChanged) < 1 {
		return nil
	}

	sqlStr, params, errSql := goqu.Dialect(store.dbDriverName).
		Update(store.blockTableName).
		Prepared(true).
		Set(dataChanged).
		Where(goqu.C(COLUMN_ID).Eq(block.ID())).
		ToSQL()

	if errSql != nil {
		return errSql
	}

	if store.debugEnabled {
		log.Println(sqlStr)
	}

	if store.db == nil {
		return errors.New("blockstore: database is nil")
	}

	_, err := database.Execute(store.toQuerableContext(ctx), sqlStr, params...)

	block.MarkAsNotDirty()

	return err
}

func (store *store) blockSelectQuery(options BlockQueryInterface) (selectDataset *goqu.SelectDataset, columns []any, err error) {
	if options == nil {
		return nil, []any{}, errors.New("block query: cannot be nil")
	}

	if err := options.Validate(); err != nil {
		return nil, []any{}, err
	}

	q := goqu.Dialect(store.dbDriverName).From(store.blockTableName)

	if options.HasCreatedAtGte() && options.HasCreatedAtLte() {
		q = q.Where(
			goqu.C(COLUMN_CREATED_AT).Gte(options.CreatedAtGte()),
			goqu.C(COLUMN_CREATED_AT).Lte(options.CreatedAtLte()),
		)
	} else if options.HasCreatedAtGte() {
		q = q.Where(goqu.C(COLUMN_CREATED_AT).Gte(options.CreatedAtGte()))
	} else if options.HasCreatedAtLte() {
		q = q.Where(goqu.C(COLUMN_CREATED_AT).Lte(options.CreatedAtLte()))
	}

	if options.HasHandle() {
		q = q.Where(goqu.C(COLUMN_HANDLE).Eq(options.Handle()))
	}

	if options.HasID() {
		q = q.Where(goqu.C(COLUMN_ID).Eq(options.ID()))
	}

	if options.HasIDIn() {
		q = q.Where(goqu.C(COLUMN_ID).In(options.IDIn()))
	}

	if options.HasNameLike() {
		q = q.Where(goqu.C(COLUMN_NAME).Like(`%` + options.NameLike() + `%`))
	}

	if options.HasPageID() {
		q = q.Where(goqu.C(COLUMN_PAGE_ID).Eq(options.PageID()))
	}

	if options.HasParentID() {
		q = q.Where(goqu.C(COLUMN_PARENT_ID).Eq(options.ParentID()))
	}

	if options.HasSequence() {
		q = q.Where(goqu.C(COLUMN_SEQUENCE).Eq(options.Sequence()))
	}

	if options.HasSiteID() {
		q = q.Where(goqu.C(COLUMN_SITE_ID).Eq(options.SiteID()))
	}

	if options.HasStatus() {
		q = q.Where(goqu.C(COLUMN_STATUS).Eq(options.Status()))
	}

	if options.HasStatusIn() {
		q = q.Where(goqu.C(COLUMN_STATUS).In(options.StatusIn()))
	}

	if options.HasTemplateID() {
		q = q.Where(goqu.C(COLUMN_TEMPLATE_ID).Eq(options.TemplateID()))
	}

	if !options.IsCountOnly() {
		if options.HasLimit() {
			q = q.Limit(uint(options.Limit()))
		}

		if options.HasOffset() {
			q = q.Offset(uint(options.Offset()))
		}
	}

	sortOrder := sb.DESC
	if options.HasSortOrder() && options.SortOrder() != "" {
		sortOrder = options.SortOrder()
	}

	if options.HasOrderBy() && options.OrderBy() != "" {
		if strings.EqualFold(sortOrder, sb.ASC) {
			q = q.Order(goqu.I(options.OrderBy()).Asc())
		} else {
			q = q.Order(goqu.I(options.OrderBy()).Desc())
		}
	}

	columns = []any{}

	for _, column := range options.Columns() {
		columns = append(columns, column)
	}

	if options.SoftDeleteIncluded() {
		return q, columns, nil // soft deleted blocks requested specifically
	}

	softDeleted := goqu.C(COLUMN_SOFT_DELETED_AT).
		Gt(carbon.Now(carbon.UTC).ToDateTimeString())

	return q.Where(softDeleted), columns, nil
}
