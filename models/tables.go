package models

import (
	"log"
)

type Table struct {
	Name string `json:"name"`
}

func GetTables() []*Table {
	res := make([]*Table, 0)

	rows, err := db.Query("SHOW TABLES")
	if err != nil {
		log.Print(err)
	}
	for rows.Next() {
		bl := Table{}
		err = rows.Scan(&bl.Name)
		res = append(res, &bl)
	}

	return res
}
