package models

import (
	"fmt"
	geojson "github.com/paulmach/go.geojson"
	"water-api/sql"
)

type GeoJson struct {
	Type string                     `json:"type"`
	Data *geojson.FeatureCollection `json:"data"`
}

type Object struct {
	Id          int               `json:"id"`
	Name        string            `json:"name"`
	LayerId     int               `json:"layerId"`
	GeoJson     *geojson.Geometry `json:"geoJson"`
	Description string            `json:"description"`
	Fields      []*FieldValue     `json:"fields"`
}

type FieldValue struct {
	Id       int    `json:"id"`
	Value    string `json:"value"`
	ValueNum int    `json:"valueNum"`
	FieldId  int    `json:"fieldId"`
}

func CreateObject(o *Object) int64 {
	out, err := o.GeoJson.MarshalJSON()
	if err != nil {
		fmt.Println(err)
		return 0
	}
	str := string(out)
	res, err := db.Exec(sql.ObjectCreate, o.LayerId, o.Name)
	if err != nil {
		fmt.Println(err)
		return 0
	}

	id, err := res.LastInsertId()
	if err != nil {
		fmt.Println(err)
		return 0
	}

	query := "insert into coordinates (id, g) values (?, ST_GeomFromGeoJSON('" + str + "'))"
	_, err = db.Exec(query, id)
	if err != nil {
		fmt.Println(err)
		return 0
	}

	for _, v := range o.Fields {
		_, err = db.Exec(sql.FieldsValuesCreate, id, v.FieldId, v.Value, v.ValueNum)
		if err != nil {
			fmt.Println(err)
			return 0
		}
	}

	return id
}
