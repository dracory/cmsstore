package cmsstore

import (
	"github.com/dracory/sb"
)

// menuItemTableCreateSql returns a SQL string for creating the menu_item table
func (st *store) menuItemTableCreateSql() string {
	// Create a new SQL builder for the database driver used by the store
	sql := sb.NewBuilder(sb.DatabaseDriverName(st.db)).
		// Define the table name using the store's menuItemTableName
		Table(st.menuItemTableName).
		// Define the ID column as a primary key string with a length of 40
		Column(sb.Column{
			Name:       COLUMN_ID,
			Type:       sb.COLUMN_TYPE_STRING,
			PrimaryKey: true,
			Length:     40,
		}).
		// Define the MENU_ID column as a string with a length of 40
		Column(sb.Column{
			Name:   COLUMN_MENU_ID,
			Type:   sb.COLUMN_TYPE_STRING,
			Length: 40,
		}).
		// Define the STATUS column as a string with a length of 40
		Column(sb.Column{
			Name:   COLUMN_STATUS,
			Type:   sb.COLUMN_TYPE_STRING,
			Length: 40,
		}).
		// Define the NAME column as a string with a length of 255
		Column(sb.Column{
			Name:   COLUMN_NAME,
			Type:   sb.COLUMN_TYPE_STRING,
			Length: 255,
		}).
		// Define the PARENT_ID column as a string with a length of 40
		Column(sb.Column{
			Name:   COLUMN_PARENT_ID,
			Type:   sb.COLUMN_TYPE_STRING,
			Length: 40,
		}).
		// Define the SEQUENCE column as an integer with a length of 10
		Column(sb.Column{
			Name:   COLUMN_SEQUENCE,
			Type:   sb.COLUMN_TYPE_INTEGER,
			Length: 10,
		}).
		// Define the PAGE_ID column as a string with a length of 40
		Column(sb.Column{
			Name:   COLUMN_PAGE_ID,
			Type:   sb.COLUMN_TYPE_STRING,
			Length: 40,
		}).
		// Define the URL column as a string with a length of 255
		Column(sb.Column{
			Name:   COLUMN_URL,
			Type:   sb.COLUMN_TYPE_STRING,
			Length: 255,
		}).
		// Define the TARGET column as a string with a length of 40
		Column(sb.Column{
			Name:   COLUMN_TARGET,
			Type:   sb.COLUMN_TYPE_STRING,
			Length: 40,
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
		// Create the table if it does not already exist
		CreateIfNotExists()

	return sql
}
