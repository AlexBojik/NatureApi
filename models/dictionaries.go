package models

import (
	"log"
	"water-api/sql"
)

type Dictionary struct {
	Name string `json:"name"`
	Id   int    `json:"id"`
}

type Value struct {
	Name   string `json:"name"`
	Id     int    `json:"id"`
	DictId int    `json:"dictId"`
}

func GetDictionaries() []*Dictionary {
	res := make([]*Dictionary, 0)

	rows, err := db.Query(sql.DictionaryList)
	if err != nil {
		log.Print(err)
	}
	for rows.Next() {
		dct := Dictionary{}
		err = rows.Scan(&dct.Id, &dct.Name)
		res = append(res, &dct)
	}

	return res
}

func CreateDictionary(dict *Dictionary) int64 {
	res, _ := db.Exec("INSERT INTO dictionaries (name) values (?)", dict.Name)
	id, _ := res.LastInsertId()
	return id
}

func UpdateDictionary(dict *Dictionary) {
	db.Exec("UPDATE dictionaries SET name = ? WHERE id = ?", dict.Name, dict.Id)
}

func DeleteDictionary(id int) {
	db.Exec("DELETE from dictionaries where id = ?", id)
}

func GetValues(id int) []*Value {
	res := make([]*Value, 0)

	rows, err := db.Query("select id, name, dictId from dictionary_values where dictId = ?", id)
	if err != nil {
		log.Print(err)
	}
	for rows.Next() {
		val := Value{}
		err = rows.Scan(&val.Id, &val.Name, &val.DictId)
		res = append(res, &val)
	}

	return res
}

func CreateValue(val *Value) int64 {
	res, _ := db.Exec("INSERT INTO dictionary_values (name, dictId) values (?, ?)", val.Name, val.DictId)
	id, _ := res.LastInsertId()
	return id
}

func UpdateValue(val *Value) {
	db.Exec("UPDATE dictionary_values SET name = ? WHERE id = ?", val.Name, val.Id)
}

func DeleteValue(id int) {
	db.Exec("DELETE from dictionary_values where id = ?", id)
}
