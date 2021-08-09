package models

import (
	"bytes"
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

var db *sql.DB

type Dump struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func init() {
	e := godotenv.Load()
	if e != nil {
		fmt.Print(e)
	}

	username := os.Getenv("DATABASE_USER")
	password := os.Getenv("DATABASE_PASSWORD")
	dbName := os.Getenv("DATABASE")
	dbHost := os.Getenv("DATABASE_HOST")

	database, err := sql.Open("mysql", username+":"+password+"@tcp("+dbHost+")/"+dbName)
	if err != nil {
		log.Println(err)
	}

	db = database
}

func CloseDB() {
	db.Close()
}

func BackupDB() {
	username := os.Getenv("DATABASE_USER")
	password := os.Getenv("DATABASE_PASSWORD")
	dbName := os.Getenv("DATABASE")

	os.Mkdir("dumps", 0777)
	filename := "dumps/nature26_" + time.Now().Format("20060102T150405") + ".sql"

	cmd := exec.Command("mysqldump", "-u"+username, "-p"+password, dbName)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Println(err)
		return
	}
	if err := cmd.Start(); err != nil {
		log.Println(err)
		return
	}
	bytes, err := ioutil.ReadAll(stdout)
	if err != nil {
		log.Println(err)
		return
	}
	err = ioutil.WriteFile(filename, bytes, 0644)
	if err != nil {
		log.Println(err)
		return
	}

	filename = strings.Replace(filename, "dumps/", "", 1)
	db.Exec("INSERT INTO dumps (name) values (?)", filename)
}

func RestoreDB(name string) {
	username := os.Getenv("DATABASE_USER")
	password := os.Getenv("DATABASE_PASSWORD")
	dbName := os.Getenv("DATABASE")

	filename := "dumps/" + name

	cmd := exec.Command("mysql", "-u"+username, "-p"+password, dbName)

	dump, err := os.Open(filename)
	if err != nil {
		log.Println(err)
	}
	cmd.Stdin = dump
	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		log.Println(err)
	}
}

func GetDumpsNames() []*Dump {
	res := make([]*Dump, 0)

	rows, err := db.Query("SELECT id, name FROM  dumps")
	if err != nil {
		log.Print(err)
	}
	for rows.Next() {
		d := Dump{}
		err = rows.Scan(&d.Id, &d.Name)
		res = append(res, &d)
	}

	return res
}
