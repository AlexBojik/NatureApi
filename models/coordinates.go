package models

import (
	"fmt"
	"water-api/sql"
)

type Coordinates struct {
	Id  int    `json:"id"`
	Wkt string `json:"wkt"`
}

func UpdateCoordinates(c *Coordinates) {
	// TODO: обработка ошибок
	_, err := db.Exec(sql.CoordinatesUpdate, c.Wkt, c.Id)
	if err != nil {
		fmt.Println(err)
	}
}
