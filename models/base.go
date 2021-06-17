package models

import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
)

var db *sql.DB

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
