package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"log"
	"os"
	"water-api/utils"
)

var db *sql.DB

func main() {
	dbinit()
	test()
}

func dbinit() {
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
	rows, _ := db.Query("select id, data from images")
	for rows.Next() {
		id := 0
		data := ""
		rows.Scan(&id, &data)
		img := utils.Base64toPng(data)
		db.Exec("UPDATE images set data = ? WHERE id = ?", img, id)
	}
}

func test() {

}
