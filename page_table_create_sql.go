package cmsstore

import (
	"github.com/gouniverse/sb"
)

// pageTableCreateSql returns a SQL string for creating the page table
func (st *store) pageTableCreateSql() string {
	// Create a new SQL builder for the database driver used by the store
	sql := sb.NewBuilder(sb.DatabaseDriverName(st.db)).
		// Define the table name using the store's page table name
		Table(st.pageTableName).
		// Define the ID column as a primary key string with a length of 40 characters
		Column(sb.Column{
			Name:       COLUMN_ID,
			Type:       sb.COLUMN_TYPE_STRING,
			PrimaryKey: true,
			Length:     40,
		}).
		// Define the SITE_ID column as a string with a length of 40 characters
		Column(sb.Column{
			Name:   COLUMN_SITE_ID,
			Type:   sb.COLUMN_TYPE_STRING,
			Length: 40,
		}).
		// Define the STATUS column as a string with a length of 40 characters
		Column(sb.Column{
			Name:   COLUMN_STATUS,
			Type:   sb.COLUMN_TYPE_STRING,
			Length: 40,
		}).
		// Define the ALIAS column as a string with a length of 255 characters
		Column(sb.Column{
			Name:   COLUMN_ALIAS,
			Type:   sb.COLUMN_TYPE_STRING,
			Length: 255,
		}).
		// Define the NAME column as a string with a length of 255 characters
		Column(sb.Column{
			Name:   COLUMN_NAME,
			Type:   sb.COLUMN_TYPE_STRING,
			Length: 255,
		}).
		// Define the TITLE column as a string with a length of 255 characters
		Column(sb.Column{
			Name:   COLUMN_TITLE,
			Type:   sb.COLUMN_TYPE_STRING,
			Length: 255,
		}).
		// Define the CONTENT column as a long text field
		Column(sb.Column{
			Name: COLUMN_CONTENT,
			Type: sb.COLUMN_TYPE_LONGTEXT,
		}).
		// Define the EDITOR column as a string with a length of 40 characters
		Column(sb.Column{
			Name:   COLUMN_EDITOR,
			Type:   sb.COLUMN_TYPE_STRING,
			Length: 40,
		}).
		// Define the TEMPLATE_ID column as a string with a length of 40 characters
		Column(sb.Column{
			Name:   COLUMN_TEMPLATE_ID,
			Type:   sb.COLUMN_TYPE_STRING,
			Length: 40,
		}).
		// Define the CANONICAL_URL column as a string with a length of 255 characters
		Column(sb.Column{
			Name:   COLUMN_CANONICAL_URL,
			Type:   sb.COLUMN_TYPE_STRING,
			Length: 255,
		}).
		// Define the META_KEYWORDS column as a string with a length of 255 characters
		Column(sb.Column{
			Name:   COLUMN_META_KEYWORDS,
			Type:   sb.COLUMN_TYPE_STRING,
			Length: 255,
		}).
		// Define the META_DESCRIPTION column as a string with a length of 255 characters
		Column(sb.Column{
			Name:   COLUMN_META_DESCRIPTION,
			Type:   sb.COLUMN_TYPE_STRING,
			Length: 255,
		}).
		// Define the META_ROBOTS column as a string with a length of 255 characters
		Column(sb.Column{
			Name:   COLUMN_META_ROBOTS,
			Type:   sb.COLUMN_TYPE_STRING,
			Length: 255,
		}).
		// Define the HANDLE column as a string with a length of 40 characters
		Column(sb.Column{
			Name:   COLUMN_HANDLE,
			Type:   sb.COLUMN_TYPE_STRING,
			Length: 40,
		}).
		// Define the MIDDLEWARES_AFTER column as a text field
		Column(sb.Column{
			Name: COLUMN_MIDDLEWARES_AFTER,
			Type: sb.COLUMN_TYPE_TEXT,
		}).
		// Define the MIDDLEWARES_BEFORE column as a text field
		Column(sb.Column{
			Name: COLUMN_MIDDLEWARES_BEFORE,
			Type: sb.COLUMN_TYPE_TEXT,
		}).
		// Define the METAS column as a text field
		Column(sb.Column{
			Name: COLUMN_METAS,
			Type: sb.COLUMN_TYPE_TEXT,
		}).
		// Define the MEMO column as a text field
		Column(sb.Column{
			Name: COLUMN_MEMO,
			Type: sb.COLUMN_TYPE_TEXT,
		}).
		// Define the CREATED_AT column as a datetime field
		Column(sb.Column{
			Name: COLUMN_CREATED_AT,
			Type: sb.COLUMN_TYPE_DATETIME,
		}).
		// Define the UPDATED_AT column as a datetime field
		Column(sb.Column{
			Name: COLUMN_UPDATED_AT,
			Type: sb.COLUMN_TYPE_DATETIME,
		}).
		// Define the SOFT_DELETED_AT column as a datetime field
		Column(sb.Column{
			Name: COLUMN_SOFT_DELETED_AT,
			Type: sb.COLUMN_TYPE_DATETIME,
		}).
		// Generate the SQL statement to create the table if it does not already exist
		CreateIfNotExists()

	return sql
}
