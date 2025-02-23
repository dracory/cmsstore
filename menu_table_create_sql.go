package cmsstore

import (
	"github.com/gouniverse/sb"
)

// SQLCreateTable returns a SQL string for creating the menu table
func (st *store) menuTableCreateSql() string {
	// Initialize a new SQL builder using the database driver name from the store
	sql := sb.NewBuilder(sb.DatabaseDriverName(st.db)).
		// Set the table name using the menuTableName from the store
		Table(st.menuTableName).
		// Define the ID column, which is the primary key with a length of 40 characters
		Column(sb.Column{
			Name:       COLUMN_ID,
			Type:       sb.COLUMN_TYPE_STRING,
			PrimaryKey: true,
			Length:     40,
		}).
		// Define the SITE_ID column with a length of 40 characters
		Column(sb.Column{
			Name:   COLUMN_SITE_ID,
			Type:   sb.COLUMN_TYPE_STRING,
			Length: 40,
		}).
		// Define the STATUS column with a length of 40 characters
		Column(sb.Column{
			Name:   COLUMN_STATUS,
			Type:   sb.COLUMN_TYPE_STRING,
			Length: 40,
		}).
		// Define the NAME column with a length of 255 characters
		Column(sb.Column{
			Name:   COLUMN_NAME,
			Type:   sb.COLUMN_TYPE_STRING,
			Length: 255,
		}).
		// Define the HANDLE column with a length of 40 characters
		Column(sb.Column{
			Name:   COLUMN_HANDLE,
			Type:   sb.COLUMN_TYPE_STRING,
			Length: 40,
		}).
		// Define the METAS column as text
		Column(sb.Column{
			Name: COLUMN_METAS,
			Type: sb.COLUMN_TYPE_TEXT,
		}).
		// Define the MEMO column as text
		Column(sb.Column{
			Name: COLUMN_MEMO,
			Type: sb.COLUMN_TYPE_TEXT,
		}).
		// Define the CREATED_AT column as a datetime
		Column(sb.Column{
			Name: COLUMN_CREATED_AT,
			Type: sb.COLUMN_TYPE_DATETIME,
		}).
		// Define the UPDATED_AT column as a datetime
		Column(sb.Column{
			Name: COLUMN_UPDATED_AT,
			Type: sb.COLUMN_TYPE_DATETIME,
		}).
		// Define the SOFT_DELETED_AT column as a datetime
		Column(sb.Column{
			Name: COLUMN_SOFT_DELETED_AT,
			Type: sb.COLUMN_TYPE_DATETIME,
		}).
		// Create the table if it does not exist
		CreateIfNotExists()

	// Return the generated SQL string
	return sql
}
