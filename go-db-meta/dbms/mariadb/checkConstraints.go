package mariadb

import (
	"database/sql"

	m "github.com/gsiems/go-db-meta/model"
)

// CheckConstraints defines the query for obtaining the check
// constraints for the tables specified by the (schemaName, tableName)
// parameters and returns the results of executing the query
func CheckConstraints(db *sql.DB, schemaName, tableName string) ([]m.CheckConstraint, error) {

	// Supported since MariaDB 10.2.1

	q := `
SELECT con.constraint_catalog AS table_catalog,
        tab.table_schema,
        tab.table_name,
        con.constraint_name,
        con.check_clause,
        'Enabled' AS status,
        NULL AS comments
    FROM information_schema.check_constraints con
    INNER JOIN information_schema.table_constraints tab
        ON ( con.constraint_schema = tab.constraint_schema
            AND con.constraint_name = tab.constraint_name )
    CROSS JOIN (
        SELECT coalesce ( ?, '' ) AS schema_name,
                coalesce ( ?, '' ) AS table_name
        ) AS args
    WHERE tab.table_schema NOT IN ( 'information_schema', 'mysql', 'performance_schema', 'sys' )
        AND ( tab.table_schema = args.schema_name OR ( args.schema_name = '' AND args.table_name = '' ) )
        AND ( tab.table_name = args.table_name OR args.table_name = '' )
`
	return m.CheckConstraints(db, q, schemaName, tableName)
}
