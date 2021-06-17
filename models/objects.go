package models

import geojson "github.com/paulmach/go.geojson"

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
}
