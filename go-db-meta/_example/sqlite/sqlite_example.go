package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"

	d "github.com/gsiems/go-db-meta/dbms"
)

func main() {

	File := "chinook.db"
	_, err = os.Stat(File)

	db, err := sql.Open("sqlite3", fmt.Sprintf("file:%s?mode=ro", File))
	FailOnErr(err)

	defer func() {
		if cerr := db.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()

	md, err := d.Init(db, d.SQLite)
	FailOnErr(err)

	fmt.Println(md.Name())

	catalog, err := md.CurrentCatalog()
	FailOnErr(err)
	fmt.Println(catalog)

	schemata, err := md.Schemata("", "")
	FailOnErr(err)
	fmt.Printf("%d schemas returned\n", len(schemata))

	tables, err := md.Tables("")
	FailOnErr(err)
	fmt.Printf("%d tables returned\n", len(tables))

	columns, err := md.Columns("", "")
	FailOnErr(err)
	fmt.Printf("%d columns returned\n", len(columns))

	indexes, err := md.Indexes("", "")
	FailOnErr(err)
	fmt.Printf("%d indexes returned\n", len(indexes))

	checkConstraints, err := md.CheckConstraints("", "")
	FailOnErr(err)
	fmt.Printf("%d checkConstraints returned\n", len(checkConstraints))

	domains, err := md.Domains("")
	FailOnErr(err)
	fmt.Printf("%d domains returned\n", len(domains))

	primaryKeys, err := md.PrimaryKeys("", "")
	FailOnErr(err)
	fmt.Printf("%d primaryKeys returned\n", len(primaryKeys))

	foreignKeys, err := md.ReferentialConstraints("", "")
	FailOnErr(err)
	fmt.Printf("%d foreignKeys returned\n", len(foreignKeys))

	uniqueConstraints, err := md.UniqueConstraints("", "")
	FailOnErr(err)
	fmt.Printf("%d uniqueConstraints returned\n", len(uniqueConstraints))

	dependencies, err := md.Dependencies("", "")
	FailOnErr(err)
	fmt.Printf("%d dependencies returned\n", len(dependencies))

	userTypes, err := md.Types("")
	FailOnErr(err)
	fmt.Printf("%d userTypes returned\n", len(userTypes))

}

func FailOnErr(err error) {
	os.Stderr.WriteString(fmt.Sprintf("%s\n", err))
	os.Exit(1)
}
