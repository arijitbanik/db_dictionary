package ora

import (
	"database/sql"

	m "github.com/gsiems/go-db-meta/model"
)

// CheckConstraints defines the query for obtaining the check
// constraints for the tables specified by the (schemaName, tableName)
// parameters and returns the results of executing the query
func CheckConstraints(db *sql.DB, schemaName, tableName string) ([]m.CheckConstraint, error) {

	q := `
WITH args AS (
    SELECT :1 AS schema_name,
            :2 AS table_name
        FROM dual
)
SELECT sys_context ( 'userenv', 'DB_NAME' ) AS table_catalog,
        con.owner AS table_schema,
        con.table_name,
        con.constraint_name,
        con.search_condition AS check_clause,
        initcap ( con.status) AS status,
        NULL AS comments
    FROM dba_constraints con
    CROSS JOIN args
    WHERE con.constraint_type = 'C'
        AND con.owner NOT IN ( ` + systemTables + ` )
        AND ( con.owner = args.schema_name OR ( args.schema_name IS NULL AND args.table_name IS NULL ) )
        AND ( con.table_name = args.table_name OR args.table_name IS NULL )
`
	return m.CheckConstraints(db, q, schemaName, tableName)
}
