package models

import (
	sql2 "database/sql"
	"fmt"
	"log"
	"water-api/sql"
)

type Filter struct {
	Type int    `json:"type"`
	Str  string `json:"str"`
}

type Check struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

type NameValue struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func CheckPosition(check *Check) []string {
	var res = make([]string, 0)
	q := "SELECT name FROM layers WHERE warning and id in (SELECT DISTINCT layerId FROM objects where id in (SELECT id FROM (SELECT id, g FROM coordinates c WHERE MBRWithin(Point(?, ?), g)) t WHERE ST_Distance(g, Point(?, ?)) < 100))"
	rows, err := db.Query(q, check.Lon, check.Lat, check.Lon, check.Lat)
	if err != nil {
		log.Print(err)
	}
	for rows.Next() {
		v := ""
		err = rows.Scan(&v)
		res = append(res, v)
	}
	return res
}

func GetFeaturesBy(filter *Filter, hasInfoRole bool) []*Object {
	var rows *sql2.Rows = nil
	var err error = nil
	switch filter.Type {
	case 1:
		rows, err = db.Query(sql.FeaturesByFilter, filter.Str, filter.Str, filter.Str)
		break
	case 2:
		sqlRaw := fmt.Sprintf(sql.FeaturesByIds, filter.Str)
		rows, err = db.Query(sqlRaw)
		break
	case 3:
		rows, err = db.Query(sql.FeaturesByRegion, filter.Str)
		break
	case 4:
		sqlRaw := fmt.Sprintf(sql.FeaturesByFields, filter.Str)
		rows, err = db.Query(sqlRaw)
		break
	}

	if err != nil {
		log.Print(err)
	}

	var objects = make([]*Object, 0)
	for rows.Next() {
		ob := Object{}
		err = rows.Scan(&ob.Id, &ob.LayerId, &ob.Name, &ob.GeoJson)

		description := ""

		rowsF, err := db.Query(sql.AdditionalFields, hasInfoRole, ob.Id)
		if err != nil {
			log.Println(err)
			continue
		}
		for rowsF.Next() {
			nv := NameValue{}
			err = rowsF.Scan(&nv.Name, &nv.Value)
			description += "<strong>" + nv.Name + ": </strong>" + nv.Value + "<br>"
		}

		ob.Description = description

		objects = append(objects, &ob)
	}
	return objects
}
