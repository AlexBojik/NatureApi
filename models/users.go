package models

import (
	"log"
	"water-api/sql"
)

type User struct {
	Name     string `json:"name"`
	Token    string `json:"token"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Snils    string `json:"snils"`
	RegAddr  string `json:"regAddr"`
	ProAddr  string `json:"proAddr"`
	Doc      string `json:"doc"`
	Admin    bool   `json:"admin"`
	Layers   bool   `json:"layers"`
	Dicts    bool   `json:"dicts"`
	Messages bool   `json:"messages"`
	Info     bool   `json:"info"`
	GroupId  int    `json:"group"`
}

type UserGroups struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Admin    bool   `json:"admin"`
	Layers   bool   `json:"layers"`
	Dicts    bool   `json:"dicts"`
	Messages bool   `json:"messages"`
	Info     bool   `json:"info"`
}

func GetUser(token string) *User {
	rows, err := db.Query(sql.User, token)
	if err != nil {
		log.Print(err)
	}
	var user = User{}

	if rows.Next() {
		err = rows.Scan(&user.Name, &user.Token, &user.Phone, &user.Email, &user.Snils, &user.RegAddr, &user.ProAddr, &user.Doc, &user.Admin, &user.Layers, &user.Dicts, &user.Messages, &user.Info, &user.GroupId)
		if err != nil {
			log.Print(err)
		}
	}
	return &user
}

func GetUserList(id int) []*User {
	res := make([]*User, 0)

	rows, err := db.Query(sql.UserList, id, id)
	if err != nil {
		log.Print(err)
	}
	for rows.Next() {
		user := User{}
		err = rows.Scan(&user.Name, &user.Token, &user.Phone, &user.Email, &user.Snils, &user.RegAddr, &user.ProAddr, &user.Doc, &user.Admin, &user.Layers, &user.Dicts, &user.Messages, &user.Info, &user.GroupId)
		res = append(res, &user)
	}

	return res
}

func CreateUser(user *User) {
	_, err := db.Exec(sql.UserCreate, &user.Name, &user.Token, &user.Phone, &user.Email, &user.Snils, &user.RegAddr, &user.ProAddr, &user.Doc)
	if err != nil {
		log.Print(err)
	}
}

func UpdateUser(user *User) {
	_, err := db.Exec(sql.UserUpdate, &user.Admin, &user.Layers, &user.Dicts, &user.Messages, &user.Info, &user.GroupId, &user.Token)
	if err != nil {
		log.Print(err)
	}
}

func GetUserGroups() []*UserGroups {
	res := make([]*UserGroups, 0)

	rows, err := db.Query(sql.UserGroupList)
	if err != nil {
		log.Print(err)
	}
	for rows.Next() {
		ug := UserGroups{}
		err = rows.Scan(&ug.Id, &ug.Name, &ug.Admin, &ug.Layers, &ug.Dicts, &ug.Messages, &ug.Info)
		res = append(res, &ug)
	}

	return res
}

func CreateUserGroup(ug *UserGroups) int64 {
	row, err := db.Exec(sql.UserGroupCreate, &ug.Name, &ug.Admin, &ug.Layers, &ug.Dicts, &ug.Messages, &ug.Info)
	if err != nil {
		log.Print(err)
	}
	id, err := row.LastInsertId()
	if err == nil {
		return id
	}
	return 0
}

func UpdateUserGroup(ug *UserGroups) {
	_, err := db.Exec(sql.UserGroupUpdate, &ug.Name, &ug.Admin, &ug.Layers, &ug.Dicts, &ug.Messages, &ug.Info, &ug.Id)
	if err != nil {
		log.Print(err)
	}
}

func HasInfoRole(token string) bool {
	return GetUser(token).Info
}
