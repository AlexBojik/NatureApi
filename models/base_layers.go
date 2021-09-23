package models

import (
	"log"
	"water-api/sql"
)

type BaseLayer struct {
	Id          int     `json:"id"`
	Name        string  `json:"name"`
	Url         string  `json:"url"`
	Description string  `json:"description"`
	MinZoom     float64 `json:"minZoom"`
	MaxZoom     float64 `json:"maxZoom"`
}

func GetBaseLayers() []*BaseLayer {
	res := make([]*BaseLayer, 0)

	rows, err := db.Query(sql.BaseLayerList)
	if err != nil {
		log.Print(err)
		return res
	}
	defer rows.Close()

	for rows.Next() {
		bl := BaseLayer{}
		err = rows.Scan(&bl.Id, &bl.Name, &bl.Url, &bl.Description, &bl.MinZoom, &bl.MaxZoom)
		res = append(res, &bl)
	}

	return res
}

func CreateBaseLayer(layer *BaseLayer) int64 {
	// TODO: обработка ошибок
	res, _ := db.Exec(sql.BaseLayersCreate, layer.Name, layer.Url, layer.Description)
	id, _ := res.LastInsertId()
	return id
}

func UpdateBaseLayer(layer *BaseLayer) {
	// TODO: обработка ошибок
	db.Exec(sql.BaseLayerUpdate, layer.Name, layer.Url, layer.Description, layer.Id)
}

func DeleteBaseLayer(id int) {
	// TODO: обработка ошибок
	db.Exec(sql.BaseLayerDelete, id)
}
