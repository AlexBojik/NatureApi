package models

import (
	"fmt"
	"log"
	"os"
)

type File struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func CreateFile(name string) {
	db.Exec("INSERT INTO files (name) values (?)", name)
}

func GetFiles() []*File {
	res := make([]*File, 0)

	rows, err := db.Query("SELECT id, name FROM files")
	if err != nil {
		log.Print(err)
	}
	for rows.Next() {
		f := File{}
		err = rows.Scan(&f.Id, &f.Name)
		res = append(res, &f)
	}

	return res
}

func DeleteFile(id int) {
	rows, err := db.Query("SELECT name FROM files WHERE id = ?", id)
	if err != nil {
		log.Print(err)
	}
	name := ""
	for rows.Next() {
		err = rows.Scan(&name)
		if err != nil {
			fmt.Println(err)
		}
		os.Remove("files/" + name)
		db.Exec("DELETE FROM files WHERE id = ?", id)
	}
}
