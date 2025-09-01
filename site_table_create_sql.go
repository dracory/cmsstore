package cmsstore

import (
	"github.com/dracory/sb"
)

// siteTableCreateSql returns a SQL string for creating the site table
func (st *store) siteTableCreateSql() string {
	// Start a new SQL builder for the database driver used by the store
	sql := sb.NewBuilder(sb.DatabaseDriverName(st.db)).
		// Set the table name to the store's site table name
		Table(st.siteTableName).
		// Define the ID column
		Column(sb.Column{
			Name:       COLUMN_ID,             // Name of the column
			Type:       sb.COLUMN_TYPE_STRING, // Data type of the column
			PrimaryKey: true,                  // Set as primary key
			Length:     40,                    // Length of the column
		}).
		// Define the STATUS column
		Column(sb.Column{
			Name:   COLUMN_STATUS,         // Name of the column
			Type:   sb.COLUMN_TYPE_STRING, // Data type of the column
			Length: 40,                    // Length of the column
		}).
		// Define the NAME column
		Column(sb.Column{
			Name:   COLUMN_NAME,           // Name of the column
			Type:   sb.COLUMN_TYPE_STRING, // Data type of the column
			Length: 255,                   // Length of the column
		}).
		// Define the DOMAIN_NAMES column
		Column(sb.Column{
			Name: COLUMN_DOMAIN_NAMES, // Name of the column
			Type: sb.COLUMN_TYPE_TEXT, // Data type of the column
		}).
		// Define the HANDLE column
		Column(sb.Column{
			Name:   COLUMN_HANDLE,         // Name of the column
			Type:   sb.COLUMN_TYPE_STRING, // Data type of the column
			Length: 40,                    // Length of the column
		}).
		// Define the METAS column
		Column(sb.Column{
			Name: COLUMN_METAS,        // Name of the column
			Type: sb.COLUMN_TYPE_TEXT, // Data type of the column
		}).
		// Define the MEMO column
		Column(sb.Column{
			Name: COLUMN_MEMO,         // Name of the column
			Type: sb.COLUMN_TYPE_TEXT, // Data type of the column
		}).
		// Define the CREATED_AT column
		Column(sb.Column{
			Name: COLUMN_CREATED_AT,       // Name of the column
			Type: sb.COLUMN_TYPE_DATETIME, // Data type of the column
		}).
		// Define the UPDATED_AT column
		Column(sb.Column{
			Name: COLUMN_UPDATED_AT,       // Name of the column
			Type: sb.COLUMN_TYPE_DATETIME, // Data type of the column
		}).
		// Define the SOFT_DELETED_AT column
		Column(sb.Column{
			Name: COLUMN_SOFT_DELETED_AT,  // Name of the column
			Type: sb.COLUMN_TYPE_DATETIME, // Data type of the column
		}).
		// Create the table if it does not exist
		CreateIfNotExists()

	return sql
}
