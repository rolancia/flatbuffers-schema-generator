package lib

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"strings"
)

func ExportAsJson(db *sql.DB) (map[string]interface{}, error) {
	tables, err := tables(db)
	if err != nil {
		return nil, err
	}

	ret := map[string]interface{}{}
	for _, table := range tables {
		ret[table], err = dump(db, table)
		if err != nil {
			return nil, err
		}
	}

	return ret, nil
}

func tables(db *sql.DB) ([]string, error) {
	rows, err := db.Query("SELECT name FROM sqlite_master WHERE type = 'table'")
	if err != nil {
		return nil, err
	}

	names := []string{}
	for rows.Next() {
		var name string
		err := rows.Scan(&name)
		if err != nil {
			return nil, err
		}
		if strings.HasPrefix(name, "_") {
			continue
		}
		names = append(names, name)
	}

	return names, nil
}

func dump(db *sql.DB, tableName string) ([]map[string]interface{}, error) {
	rows, err := db.Query(fmt.Sprintf("SELECT * FROM %s", tableName))
	if err != nil {
		return nil, err
	}

	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	ret := []map[string]interface{}{}
	for rows.Next() {

		columns := make([]interface{}, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i, _ := range columns {
			columnPointers[i] = &columns[i]
		}

		if err := rows.Scan(columnPointers...); err != nil {
			return nil, err
		}

		m := make(map[string]interface{})
		for i, colName := range cols {
			val := columnPointers[i].(*interface{})
			m[colName] = *val
		}

		ret = append(ret, m)
	}

	return ret, nil
}
