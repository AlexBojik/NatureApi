package models

import (
	"log"
	"water-api/sql"
)

type Field struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func GetFieldsList() []*Field {
	res := make([]*Field, 0)

	rows, err := db.Query(sql.FieldsList)
	if err != nil {
		log.Print(err)
	}
	for rows.Next() {
		f := Field{}
		err = rows.Scan(&f.Id, &f.Name)
		res = append(res, &f)
	}

	return res
}
