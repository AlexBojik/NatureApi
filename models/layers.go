package models

import (
	"fmt"
	geojson "github.com/paulmach/go.geojson"
	"log"
	"water-api/sql"
)

type GroupLayer struct {
	Id      int     `json:"id"`
	Name    string  `json:"name"`
	Icon    string  `json:"icon"`
	Layers  []Layer `json:"layers"`
	IsGroup bool    `json:"isGroup"`
}

type Layer struct {
	Group             int    `json:"group"`
	Id                int    `json:"id"`
	Name              string `json:"name"`
	Type              string `json:"type"`
	Url               string `json:"url"`
	Color             string `json:"color"`
	CommonName        string `json:"commonName"`
	CommonDescription string `json:"commonDescription"`
	Symbol            string `json:"symbol"`
	IsGroup           bool   `json:"isGroup"`
	Cluster           bool   `json:"cluster"`
	LineWidth         int    `json:"lineWidth"`
	LineColor         string `json:"lineColor"`
	Warning           bool   `json:"warning"`
	Order             int    `json:"order"`
	Limitation        bool   `json:"limitation"`
}

func GetLayers(hasInfoRole bool) []*GroupLayer {
	res := make([]*GroupLayer, 0)

	rows, err := db.Query(sql.GroupLayerList)
	if err != nil {
		log.Print(err)
	}
	for rows.Next() {
		bl := GroupLayer{}
		bl.Layers = make([]Layer, 0)
		bl.IsGroup = true
		err = rows.Scan(&bl.Id, &bl.Name, &bl.Icon)

		rowsL, err := db.Query(sql.LayerList, bl.Id, hasInfoRole)
		if err != nil {
			log.Print(err)
		}
		for rowsL.Next() {
			l := Layer{}
			l.IsGroup = false
			rowsL.Scan(&l.Group,
				&l.Id,
				&l.Name,
				&l.Type,
				&l.Url,
				&l.Color,
				&l.CommonName,
				&l.CommonDescription,
				&l.Symbol,
				&l.Cluster,
				&l.LineWidth,
				&l.LineColor,
				&l.Warning,
				&l.Limitation,
				&l.Order)
			bl.Layers = append(bl.Layers, l)
		}

		res = append(res, &bl)
	}

	return res
}

func GetLayer(id int) *geojson.FeatureCollection {
	rows, err := db.Query(sql.Layer, id)
	if err != nil {
		log.Print(err)
	}

	var objects = make([]*geojson.Feature, 0)
	for rows.Next() {
		ob := Object{}
		err = rows.Scan(&ob.Id, &ob.LayerId, &ob.Name, &ob.GeoJson)

		f := geojson.Feature{}
		f.Geometry = ob.GeoJson
		f.ID = ob.Id
		f.Properties = map[string]interface{}{"name": ob.Name, "layerId": ob.LayerId}

		objects = append(objects, &f)
	}
	res := geojson.FeatureCollection{}
	res.Features = objects[:]

	return &res
}

func GetCluster(id int) *geojson.FeatureCollection {
	rows, err := db.Query(sql.ClusterCoordinate, id)
	if err != nil {
		log.Print(err)
	}

	var objects = make([]*geojson.Feature, 0)
	for rows.Next() {
		ob := Object{}
		err = rows.Scan(&ob.Id, &ob.LayerId, &ob.Name, &ob.GeoJson)

		f := geojson.Feature{}
		f.Geometry = ob.GeoJson
		f.ID = ob.Id
		f.Properties = map[string]interface{}{"name": ob.Name, "layerId": ob.LayerId}

		objects = append(objects, &f)
	}
	res := geojson.FeatureCollection{}
	res.Features = objects[:]

	return &res
}

func CreateLayer(l *Layer) int64 {
	res, _ := db.Exec(sql.LayerCreate,
		l.Name,
		l.Type,
		l.Group,
		l.Url,
		l.Color,
		l.CommonName,
		l.CommonDescription,
		l.Warning,
		l.Symbol,
		l.Cluster,
		l.Order,
		l.LineWidth,
		l.LineColor,
		l.Limitation)
	id, _ := res.LastInsertId()
	return id
}

func UpdateLayer(l *Layer) {
	_, err := db.Exec(sql.LayerUpdate,
		l.Name,
		l.Type,
		l.Group,
		l.Url,
		l.Color,
		l.CommonName,
		l.CommonDescription,
		l.Warning,
		l.Symbol,
		l.Cluster,
		l.Order,
		l.LineWidth,
		l.LineColor,
		l.Limitation,
		l.Id)

	if err != nil {
		fmt.Print(err)
	}
}

func DeleteLayer(id int) {
	db.Exec("DELETE FROM layers where id = ?", id)
}
