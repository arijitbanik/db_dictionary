package ora

import (
	"database/sql"

	m "github.com/gsiems/go-db-meta/model"
)

// ReferentialConstraints defines the query for obtaining the
// referential constraints for the (schemaName, tableName) parameters
// (as either the parent or child) and returns the results of executing
// the query
func ReferentialConstraints(db *sql.DB, schemaName, tableName string) ([]m.ReferentialConstraint, error) {

	q := `
WITH args AS (
    SELECT :1 AS schema_name,
            :2 AS table_name
        FROM dual
),
cols AS (
    SELECT col.owner,
            col.table_name,
            col.constraint_name,
            listagg ( col.column_name, ', ' ) WITHIN GROUP (
                ORDER BY col.position ) AS table_columns
        FROM dba_cons_columns col
        WHERE col.owner NOT IN (` + systemTables + ` )
        GROUP BY col.owner,
            col.table_name,
            col.constraint_name
)
SELECT sys_context ( 'userenv', 'DB_NAME' ) AS table_catalog,
        con.owner AS table_schema,
        con.table_name,
        col.table_columns,
        --con.constraint_type,
        con.constraint_name,
        sys_context ( 'userenv', 'DB_NAME' ) AS ref_table_catalog,
        rcon.owner AS ref_table_schema,
        rcon.table_name AS ref_table_name,
        rcol.table_columns AS ref_table_columns,
        --rcon.constraint_type AS ref_constraint_type,
        rcon.constraint_name AS ref_constraint_name,
        NULL AS match_option, -- TODO
        con.delete_rule,
        'RESTRICT' AS update_rule,
        CASE con.status
            WHEN 'ENABLED' THEN 'YES'
            WHEN 'DISABLED' THEN 'NO'
            ELSE con.status
            END AS is_enforced,
        --con.deferrable AS is_deferrable,
        --con.deferred AS initially_deferred,
        NULL AS comments
    FROM dba_constraints con
    JOIN dba_constraints rcon
        ON ( con.r_owner = rcon.owner
            AND con.r_constraint_name = rcon.constraint_name )
    CROSS JOIN args
    JOIN cols col
        ON ( con.owner = col.owner
            AND con.constraint_name = col.constraint_name
            AND con.table_name = col.table_name )
    JOIN cols rcol
        ON ( rcon.owner = rcol.owner
            AND rcon.constraint_name = rcol.constraint_name
            AND rcon.table_name = rcol.table_name )
    WHERE con.constraint_type = 'R'
        AND con.owner NOT IN (` + systemTables + ` )
        AND rcon.owner NOT IN ( ` + systemTables + ` )
        AND ( ( ( con.owner = args.schema_name OR ( args.schema_name IS NULL AND args.table_name IS NULL ) )
                AND ( con.table_name = args.table_name OR args.table_name IS NULL ) )
            OR ( ( rcon.owner = args.schema_name OR ( args.schema_name IS NULL AND args.table_name IS NULL ) )
                AND ( rcon.table_name = args.table_name OR args.table_name IS NULL ) ) )
`
	return m.ReferentialConstraints(db, q, schemaName, tableName)
}
