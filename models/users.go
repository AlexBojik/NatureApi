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
}

type UserGroups struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Admin    bool   `json:"admin"`
	Layers   bool   `json:"layers"`
	Dicts    bool   `json:"dicts"`
	Messages bool   `json:"messages"`
}

func GetUser(token string) *User {
	rows, err := db.Query(sql.User, token)
	if err != nil {
		log.Print(err)
	}
	var user = User{}

	if rows.Next() {
		err = rows.Scan(&user.Name, &user.Token, &user.Phone, &user.Email, &user.Snils, &user.RegAddr, &user.ProAddr, &user.Doc, &user.Admin, &user.Layers, &user.Dicts, &user.Messages)
		if err != nil {
			log.Print(err)
		}
	}
	return &user
}

func GetUserList() []*User {
	res := make([]*User, 0)

	rows, err := db.Query(sql.UserList)
	if err != nil {
		log.Print(err)
	}
	for rows.Next() {
		user := User{}
		err = rows.Scan(&user.Name, &user.Token, &user.Phone, &user.Email, &user.Snils, &user.RegAddr, &user.ProAddr, &user.Doc, &user.Admin, &user.Layers, &user.Dicts, &user.Messages)
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
	_, err := db.Exec(sql.UserUpdate, &user.Admin, &user.Layers, &user.Dicts, &user.Messages, &user.Token)
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
		err = rows.Scan(&ug.Id, &ug.Name, &ug.Admin, &ug.Layers, &ug.Dicts, &ug.Messages)
		res = append(res, &ug)
	}

	return res
}

func CreateUserGroup(ug *UserGroups) {
	_, err := db.Exec(sql.UserGroupCreate, &ug.Name, &ug.Admin, &ug.Layers, &ug.Dicts, &ug.Messages)
	if err != nil {
		log.Print(err)
	}
}

func UpdateUserGroup(ug *UserGroups) {
	_, err := db.Exec(sql.UserGroupUpdate, &ug.Name, &ug.Admin, &ug.Layers, &ug.Dicts, &ug.Messages, &ug.Id)
	if err != nil {
		log.Print(err)
	}
}