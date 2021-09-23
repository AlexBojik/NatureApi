package models

import (
	"log"
	"water-api/sql"
)

type DictionaryValue struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type Field struct {
	Id         int      `json:"id"`
	Name       string   `json:"name"`
	Type       int      `json:"type"`
	Limitation bool     `json:"limitation"`
	Options    []*Value `json:"options"`
}

func GetDictionariesValuesList() []*DictionaryValue {
	res := make([]*DictionaryValue, 0)

	rows, err := db.Query(sql.FieldsList)
	if err != nil {
		log.Print(err)
		return res
	}
	defer rows.Close()

	for rows.Next() {
		f := DictionaryValue{}
		err = rows.Scan(&f.Id, &f.Name)
		res = append(res, &f)
	}

	return res
}

func GetFieldsList() []*Field {
	res := make([]*Field, 0)

	rows, err := db.Query("select id, name, type, limitation from fields")
	if err != nil {
		log.Print(err)
		return res
	}
	defer rows.Close()

	for rows.Next() {
		f := Field{}
		err = rows.Scan(&f.Id, &f.Name, &f.Type, &f.Limitation)

		var opts = make([]*Value, 0)
		rOpt, _ := db.Query("select id, name from dictionary_values where dictId = ?", f.Type)
		for rOpt.Next() {
			v := Value{}
			_ = rOpt.Scan(&v.Id, &v.Name)
			opts = append(opts, &v)
		}
		f.Options = opts

		res = append(res, &f)
	}

	return res
}

func GetFieldsByLayerId(id int) []*Field {
	res := make([]*Field, 0)

	rows, err := db.Query("select id, name, type from fields where layerId = ?", id)
	if err != nil {
		log.Print(err)
		return res
	}

	defer rows.Close()

	for rows.Next() {
		f := Field{}
		err = rows.Scan(&f.Id, &f.Name, f.Type)
		res = append(res, &f)
	}

	return res
}

func DeleteFieldById(id int) {
	db.Exec("DELETE FROM fields where id = ?", id)
}

func CreateField(f *Field) int64 {
	res, _ := db.Exec("INSERT into fields (name, layerId, type, limitation) values (?, 0, ?, ?)", f.Name, f.Type, f.Limitation)
	id, _ := res.LastInsertId()
	return id
}

func UpdateField(f *Field) {
	db.Exec("UPDATE fields SET name = ?, type = ?, limitation = ? where id = ?", f.Name, f.Type, f.Limitation, f.Id)
}
