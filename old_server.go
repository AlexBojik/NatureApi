//package main
//
//import (
//	"database/sql"
//	"encoding/json"
//	"fmt"
//	"github.com/go-martini/martini"
//	_ "github.com/go-sql-driver/mysql"
//	"github.com/martini-contrib/cors"
//	"github.com/martini-contrib/render"
//	"github.com/paulmach/go.geojson"
//	"log"
//	"net/http"
//	"strconv"
//	"strings"
//)
//
//var db *sql.DB
//
//type Geojson struct {
//	Type string                     `json:"type"`
//	Data *geojson.FeatureCollection `json:"data"`
//}
//
//type Filter struct {
//	WKT string `json:"wkt""`
//}
//
//type FilterStr struct {
//	Str string `json:"str""`
//}
//type Object struct {
//	Id          int               `json:"id"`
//	Name        string            `json:"name"`
//	LayerId     int               `json:"layerId"`
//	GeoJson     *geojson.Geometry `json:"geoJson"`
//	Description string            `json:"description"`
//	Fields      []*FieldValue     `json:"fields"`
//}
//
//type NameValue struct {
//	Name  string `json:"name"`
//	Value string `json:"value"`
//}
//
//type Dict struct {
//	Name string `json:"name"`
//	Id   int    `json:"id"`
//}
//
//type Value struct {
//	Name   string `json:"name"`
//	Id     int    `json:"id"`
//	DictId int    `json:"dictId"`
//}
//
//type Field struct {
//	Name    string   `json:"name"`
//	Id      int      `json:"id"`
//	Type    int      `json:"type"`
//	Options []*Value `json:"options"`
//}
//
//type FieldValue struct {
//	Id       int    `json:"id"`
//	Value    string `json:"value"`
//	ValueNum int    `json:"valueNum"`
//	FieldId  int    `json:"fieldId"`
//}
//
//type Layer struct {
//	Id                int       `json:"id"`
//	Name              string    `json:"name"`
//	Layers            []*Layer  `json:"children"`
//	Table             []*Object `json:"table"`
//	Type              string    `json:"type"`
//	Lock              bool      `json:"lock"`
//	Params            string    `json:"params"`
//	Visible           bool      `json:"visible"`
//	Expandable        bool      `json:"expandable"`
//	Icon              string    `json:"icon"`
//	Color             string    `json:"color"`
//	CommonName        string    `json:"commonName"`
//	CommonDescription string    `json:"commonDescription"`
//	Symbol            string    `json:"symbol"`
//	parentId          int
//	Fields            []*Field `json:"fields"`
//}
//
//func layersHandler(r render.Render) {
//	set := make(map[int]*Layer)
//	var layers = make([]*Layer, 0)
//
//	rows, err := db.Query("select id, name, type, is_lock, parent_id, params, visible, icon, color, commonName, commonDescription, symbol from layers order by level")
//	if err != nil {
//		log.Print(err)
//	}
//	for rows.Next() {
//		l := Layer{}
//		// l.Layers = make([]*Layer, 0)
//		err = rows.Scan(&l.Id, &l.Name, &l.Type, &l.Lock, &l.parentId, &l.Params, &l.Visible, &l.Icon, &l.Color, &l.CommonName, &l.CommonDescription, &l.Symbol)
//		l.Expandable = true
//		set[l.Id] = &l
//		if l.parentId > 0 {
//			lp := set[l.parentId]
//			if lp.Layers == nil {
//				lp.Layers = make([]*Layer, 0)
//			}
//			lp.Layers = append(lp.Layers, &l)
//		} else {
//			layers = append(layers, set[l.Id])
//		}
//	}
//	r.JSON(200, layers)
//}
//
//func layersHandler1(r render.Render) {
//	set := make(map[int]*Layer)
//	var layers = make([]*Layer, 0)
//
//	rows, err := db.Query("select id, name, type, is_lock, parent_id, params, visible, icon, color, commonName, commonDescription, symbol from layers where NOT (id = 9 OR id = 10) order by level")
//	if err != nil {
//		log.Print(err)
//	}
//	for rows.Next() {
//		l := Layer{}
//		// l.Layers = make([]*Layer, 0)
//		err = rows.Scan(&l.Id, &l.Name, &l.Type, &l.Lock, &l.parentId, &l.Params, &l.Visible, &l.Icon, &l.Color, &l.CommonName, &l.CommonDescription, &l.Symbol)
//		l.Expandable = true
//		set[l.Id] = &l
//		if l.parentId > 0 {
//			lp := set[l.parentId]
//			if lp.Layers == nil {
//				lp.Layers = make([]*Layer, 0)
//			}
//			lp.Layers = append(lp.Layers, &l)
//		} else {
//			layers = append(layers, set[l.Id])
//		}
//	}
//	r.JSON(200, layers)
//}
//
//func fieldsHandler(params martini.Params, r render.Render) {
//	var fields = make([]Field, 0)
//
//	rows, err := db.Query("select id, name, type from fields where layerId = ?", params["id"])
//	if err != nil {
//		log.Print(err)
//	}
//	for rows.Next() {
//		f := Field{}
//		err = rows.Scan(&f.Id, &f.Name, &f.Type)
//
//		var opts = make([]*Value, 0)
//		rOpt, err := db.Query("select id, name from dictionary_values where dictId = ?", f.Type)
//		if err != nil {
//			log.Print(err)
//		}
//		for rOpt.Next() {
//			v := Value{}
//			err = rOpt.Scan(&v.Id, &v.Name)
//			opts = append(opts, &v)
//		}
//
//		f.Options = opts
//		fields = append(fields, f)
//	}
//	r.JSON(200, fields)
//}
//
//func layerHandler(params martini.Params, r render.Render) {
//	var objects = make([]Object, 0)
//	const q = "select id, layerId, name, ST_AsGeoJSON(geometry) as geojson, description from objects where layerId = ?"
//	rows, err := db.Query(q, params["id"])
//	if err != nil {
//		log.Print(err)
//	}
//	for rows.Next() {
//		ob := Object{}
//		err = rows.Scan(&ob.Id, &ob.LayerId, &ob.Name, &ob.GeoJson, &ob.Description)
//		objects = append(objects, ob)
//	}
//	r.JSON(200, objects)
//}
//
//func layerHandler1(params martini.Params, r render.Render) {
//	const q = "select id, layerId, name, ST_AsGeoJSON(geometry), description as geojson from objects where layerId = ?"
//	rows, err := db.Query(q, params["id"])
//	if err != nil {
//		log.Print(err)
//	}
//
//	var objects = make([]*geojson.Feature, 0)
//	for rows.Next() {
//		ob := Object{}
//		err = rows.Scan(&ob.Id, &ob.LayerId, &ob.Name, &ob.GeoJson, &ob.Description)
//		f := geojson.Feature{}
//		f.Geometry = ob.GeoJson
//		f.ID = ob.Id
//		prop := map[string]interface{}{"name": ob.Name, "layerId": ob.LayerId}
//		f.Properties = prop
//		objects = append(objects, &f)
//	}
//	res := geojson.FeatureCollection{}
//	res.Features = objects[:]
//
//	geores := Geojson{}
//	geores.Type = "geojson"
//	geores.Data = &res
//
//	r.JSON(200, res)
//}
//
//func dictValuesHandler(params martini.Params, r render.Render) {
//	var values = make([]Value, 0)
//
//	rows, err := db.Query("select id, name from dictionary_values where dictId=?", params["id"])
//	if err != nil {
//		log.Print(err)
//	}
//	for rows.Next() {
//		v := Value{}
//		err = rows.Scan(&v.Id, &v.Name)
//		values = append(values, v)
//	}
//	r.JSON(200, values)
//}
//
//func fieldValueHandler(params martini.Params, r render.Render) {
//	var values = make([]NameValue, 0)
//
//	query := "select f.name, IFNULL(dv.name, fv.value) as value from fields_values fv INNER JOIN fields f ON f.id = fv.fieldId LEFT JOIN dictionary_values dv  ON f.`type` = dv.dictId AND fv.value_num = dv.id where objectId = ? order by fv.fieldId"
//	rows, err := db.Query(query, params["id"])
//	if err != nil {
//		log.Print(err)
//	}
//	for rows.Next() {
//		v := NameValue{}
//		err = rows.Scan(&v.Name, &v.Value)
//		values = append(values, v)
//	}
//	r.JSON(http.StatusOK, values)
//}
//
//type ObjectFields struct {
//	Id     int    `json:"id"`
//	Fields []*NameValue `json:"fields"`
//}
//
//func fieldValuesHandler(params martini.Params, r render.Render) {
//	var res = make([]ObjectFields, 0)
//	var objs = strings.Split(params["ids"], ",")
//
//	for _, id := range objs {
//		var of = ObjectFields{}
//		query := "select f.name, IFNULL(dv.name, fv.value) as value from fields_values fv INNER JOIN fields f ON f.id = fv.fieldId LEFT JOIN dictionary_values dv  ON f.`type` = dv.dictId AND fv.value_num = dv.id where objectId = ? order by fv.fieldId"
//		rows, err := db.Query(query, id)
//		if err != nil {
//			log.Print(err)
//		}
//		fields := make([]*NameValue, 0)
//		for rows.Next() {
//			nv := NameValue{}
//			err = rows.Scan(&nv.Name, &nv.Value)
//			fields = append(fields, &nv)
//		}
//		of.Id, _ = strconv.Atoi(id)
//		of.Fields = fields
//		res = append(res, of)
//	}
//	r.JSON(http.StatusOK, res)
//}
//
//func checkHandler(params martini.Params, r render.Render) {
//	var res = make([]string, 0)
//	q := "SELECT name FROM layers WHERE id in (SELECT DISTINCT layerId FROM objects where id in (SELECT id FROM geom WHERE MBRWithin(ST_SRID(Point(?,?), 4326), g)))"
//	rows, err := db.Query(q, params["lat"], params["lon"])
//	if err != nil {
//		log.Print(err)
//	}
//	for rows.Next() {
//		v := Value{}
//		err = rows.Scan(&v.Name)
//		res = append(res, v.Name)
//	}
//	r.JSON(200, strings.Join(res, ", "))
//}
//
//func dictHandler(r render.Render) {
//	var dicts = make([]Dict, 0)
//
//	rows, err := db.Query("select id, name from dictionaries")
//	if err != nil {
//		log.Print(err)
//	}
//	for rows.Next() {
//		d := Dict{}
//		err = rows.Scan(&d.Id, &d.Name)
//		dicts = append(dicts, d)
//	}
//	r.JSON(200, dicts)
//}
//
//func loadLayersHandler(req *http.Request, r render.Render) {
//	var l []*Layer
//	_ = json.NewDecoder(req.Body).Decode(&l)
//	recursiveLoad(l, 0)
//	r.Status(http.StatusCreated)
//}
//
//func loadDictionariesHandler(req *http.Request, r render.Render) {
//	var d *Dict
//	err := json.NewDecoder(req.Body).Decode(&d)
//	if err != nil {
//		fmt.Println(err)
//		r.Status(http.StatusInternalServerError)
//		return
//	}
//
//	_, err = db.Exec("insert into dictionaries (name) values (?)", d.Name)
//
//	if err != nil {
//		fmt.Println(err)
//		r.Status(http.StatusInternalServerError)
//		return
//	}
//
//	r.Status(http.StatusCreated)
//}
//
//func loadObjectsHandler(req *http.Request, r render.Render) {
//	var ob *Object
//	_ = json.NewDecoder(req.Body).Decode(&ob)
//
//	query := "insert into objects (layerId, name, tempGeo, description) values (?, ?, ?, ?)"
//	out, err := ob.GeoJson.MarshalJSON()
//	if err != nil {
//		fmt.Println(err)
//		r.Text(http.StatusInternalServerError, err.Error())
//		return
//	}
//	str := string(out)
//	res, err := db.Exec(query, ob.LayerId, ob.Name, str, ob.Description)
//	if err != nil {
//		fmt.Println(err)
//		r.Text(http.StatusInternalServerError, err.Error())
//		return
//	}
//
//	id, err := res.LastInsertId()
//	if err != nil {
//		fmt.Println(err)
//		r.Text(http.StatusInternalServerError, err.Error())
//		return
//	}
//	query = "insert into fields_values (objectId, fieldId, value, value_num) values (?, ?, ?, ?)"
//	for _, v := range ob.Fields {
//		_, err = db.Exec(query, id, v.FieldId, v.Value, v.ValueNum)
//		if err != nil {
//			fmt.Println(err)
//			r.Text(http.StatusInternalServerError, err.Error())
//			return
//		}
//	}
//
//	//query = "UPDATE objects SET geometry = ST_GeomFromGeoJSON(tempGeo), tempGeo = '' WHERE tempGeo <> ''"
//	//_, err = db.Exec(query)
//	//if err != nil {
//	//	fmt.Println(err)
//	//	r.Text(http.StatusInternalServerError, err.Error())
//	//	return
//	//}
//
//	r.Status(http.StatusCreated)
//}
//
//func filterHandler(req *http.Request, r render.Render) {
//	var f *Filter
//	_ = json.NewDecoder(req.Body).Decode(&f)
//
//	var objects = make([]Object, 0)
//	const q = "select id, layerId, name, ST_AsGeoJSON(geometry) as geojson, description FROM objects WHERE MBRContains(ST_SRID(ST_GeomFromText(?), 4326), geometry)"
//	rows, err := db.Query(q, f.WKT)
//	if err != nil {
//		log.Print(err)
//	}
//
//	for rows.Next() {
//		ob := Object{}
//		err = rows.Scan(&ob.Id, &ob.LayerId, &ob.Name, &ob.GeoJson, &ob.Description)
//		objects = append(objects, ob)
//	}
//	r.JSON(http.StatusOK, objects)
//}
//
//func filterStrHandler(req *http.Request, r render.Render) {
//	var f *FilterStr
//	_ = json.NewDecoder(req.Body).Decode(&f)
//
//	var objects = make([]Object, 0)
//	const q = "select id, layerId, name, ST_AsGeoJSON(geometry) as geojson, description FROM objects WHERE name like ? limit 100"
//	rows, err := db.Query(q, "%"+f.Str+"%")
//	if err != nil {
//		log.Print(err)
//	}
//
//	for rows.Next() {
//		ob := Object{}
//		err = rows.Scan(&ob.Id, &ob.LayerId, &ob.Name, &ob.GeoJson, &ob.Description)
//		objects = append(objects, ob)
//	}
//	r.JSON(http.StatusOK, objects)
//}
//
//func recursiveLoad(layers []*Layer, parent int64) {
//	//for _, l := range layers {
//	//	if l.Type =="group" {
//	//		query := "insert into layers (name, type, is_lock, parent_id, params, visible) values (?,?,?,?,?,?)"
//	//		res, _ := db.Exec(query, l.Name, l.Type, false, parent, l.Params, false)
//	//		id, _ := res.LastInsertId()
//	//		recursiveLoad(l.Layers, id)
//	//	} else {
//	//		query := "insert into layers (name, type, is_lock, parent_id, params, visible) values (?,?,?,?,?,?)"
//	//		res, _ := db.Exec(query, l.Name, l.Type, false, parent, l.Params, false)
//	//		id, _ := res.LastInsertId()
//	//		for _, row := range l.Table {
//	//			query := "insert into objects (name, description, layerId, coordinates) values (?,?,?,?)"
//	//			_, _ = db.Exec(query, row.Name, row.Description, id, row.Coordinates)
//	//		}
//	//	}
//	//}
//}
//
//func addLayersHandler(req *http.Request, r render.Render) {
//	var l Layer
//	err := json.NewDecoder(req.Body).Decode(&l)
//
//	query := "insert into layers (name, type, is_lock, parent_id, params, visible) values (?,?,?,?,?,?,?)"
//	_, err = db.Exec(query, l.Id, l.Name, l.Type, l.Lock, l.parentId, l.Params, l.Visible)
//
//	if err != nil {
//		fmt.Println(err)
//	}
//
//	fmt.Println(l)
//	r.Status(http.StatusCreated)
//}
//
//func deleteObjectHandler(params martini.Params, r render.Render) {
//	const q = "delete from objects where id = ?"
//	_, err := db.Exec(q, params["id"])
//	if err != nil {
//		fmt.Println(err)
//		r.Status(http.StatusInternalServerError)
//		return
//	}
//	r.Status(http.StatusOK)
//}
//
//func addDictValuesHandler(req *http.Request, r render.Render) {
//	var v Value
//	_ = json.NewDecoder(req.Body).Decode(&v)
//
//	row, err := db.Exec("insert into dictionary_values (name, dictId) values(?, ?)", v.Name, v.DictId)
//	if err != nil {
//		fmt.Println(err)
//		r.Status(http.StatusInternalServerError)
//		return
//	}
//	id, _ := row.LastInsertId()
//	r.JSON(http.StatusCreated, id)
//}
//
//type Image struct {
//	Jpeg string `json:"jpeg"`
//}
//
//type UserMessage struct {
//	Images []Image `json:"images"`
//	Text   string  `json:"text"`
//	Lat    float32 `json:"lat"`
//	Lon    float32 `json:"lon"`
//	UserId int     `json:"userId"`
//}
//
//func sendHandler(req *http.Request, r render.Render) {
//	var um UserMessage
//	_ = json.NewDecoder(req.Body).Decode(&um)
//	point := "POINT(" + fmt.Sprintf("%f", um.Lat) + " " + fmt.Sprintf("%f", um.Lon) + ")"
//	row, err := db.Exec("insert into user_messages (user_id, text, point) values(?, ?, ST_GeomFromText(?))", um.UserId, um.Text, point)
//	if err != nil {
//		fmt.Println(err)
//		r.Status(http.StatusInternalServerError)
//		return
//	}
//	id, err := row.LastInsertId()
//	if err != nil {
//		fmt.Println(err)
//		r.Status(http.StatusInternalServerError)
//		return
//	}
//	query := "insert into images (message_id, data) values(?, ?)"
//	for _, img := range um.Images {
//		_, err = db.Exec(query, id, img.Jpeg)
//		if err != nil {
//			fmt.Println(err)
//			r.Error(http.StatusInternalServerError)
//			return
//		}
//	}
//	r.Status(http.StatusCreated)//	var um UserMessage
//	_ = json.NewDecoder(req.Body).Decode(&um)
//	point := "POINT(" + fmt.Sprintf("%f", um.Lat) + " " + fmt.Sprintf("%f", um.Lon) + ")"
//	row, err := db.Exec("insert into user_messages (user_id, text, point) values(?, ?, ST_GeomFromText(?))", um.UserId, um.Text, point)
//	if err != nil {
//		fmt.Println(err)
//		r.Status(http.StatusInternalServerError)
//		return
//	}
//	id, err := row.LastInsertId()
//	if err != nil {
//		fmt.Println(err)
//		r.Status(http.StatusInternalServerError)
//		return
//	}
//	query := "insert into images (message_id, data) values(?, ?)"
//	for _, img := range um.Images {
//		_, err = db.Exec(query, id, img.Jpeg)
//		if err != nil {
//			fmt.Println(err)
//			r.Error(http.StatusInternalServerError)
//			return
//		}
//	}
//	r.Status(http.StatusCreated)
//}
//
//func updateLayerHandler(req *http.Request, r render.Render) {
//	var l Layer
//	err := json.NewDecoder(req.Body).Decode(&l)
//
//	if err != nil {
//		fmt.Println(err)
//		r.Status(http.StatusInternalServerError)
//		return
//	}
//
//	query := "update layers set name = ?, color = ? where id = ?"
//	_, err = db.Exec(query, l.Name, l.Color, l.Id)
//	if err != nil {
//		fmt.Println(err)
//		r.Status(http.StatusInternalServerError)
//		return
//	}
//	query = "delete from fields where layerId = ?"
//	_, err = db.Exec(query, l.Id)
//	if err != nil {
//		fmt.Println(err)
//		r.Status(http.StatusInternalServerError)
//		return
//	}
//	for _, field := range l.Fields {
//		query := "insert into fields (name, type, layerId) values (?,?,?)"
//		_, err = db.Exec(query, field.Name, field.Type, l.Id)
//		if err != nil {
//			r.Status(http.StatusInternalServerError)
//			return
//		}
//	}
//	r.Status(http.StatusCreated)
//}
//
//// object handlers
//
//type Obj1 struct {
//	Name  string `json:"name"`
//	Layer int    `json:"layer"`
//	Uid   string `json:"uid"`
//}
//
//func objectPostHandler(r render.Render, req *http.Request) {
//	var obj Obj1
//
//	err := json.NewDecoder(req.Body).Decode(&obj)
//	if err != nil {
//		r.Error(http.StatusInternalServerError)
//		return
//	}
//
//	query := "insert into objects_1(uid, name, layer_id) values (?, ?, ?)"
//	res, err := db.Exec(query, obj.Uid, obj.Name, obj.Layer)
//	if err != nil {
//		r.Error(http.StatusInternalServerError)
//		return
//	}
//
//	id, err := res.LastInsertId()
//	if err != nil {
//		r.Error(http.StatusInternalServerError)
//		return
//	}
//
//	r.JSON(http.StatusCreated, id)
//}
//
//type Message struct {
//	Id       int    `json:"id"`
//	UserName string `json:"userName"`
//	Text     string `json:"text"`
//	Status   int    `json:"status"`
//}
//
//func messagesHandler(r render.Render) {
//	var ms = make([]Message, 0)
//
//	rows, err := db.Query("SELECT um.id, um.text, u.name, um.status FROM user_messages um INNER JOIN users u ON um.user_id = u.id")
//	if err != nil {
//		log.Print(err)
//	}
//	for rows.Next() {
//		m := Message{}
//		err = rows.Scan(&m.Id, &m.Text, &m.UserName, &m.Status)
//		ms = append(ms, m)
//	}
//	r.JSON(200, ms)
//}
//
//func newMessagesHandler(r render.Render) {
//	var count = 0
//
//	row := db.QueryRow("SELECT COUNT(id) FROM user_messages WHERE status = 0")
//	row.Scan(&count)
//	r.JSON(200, count)
//}
//
//func messageHandler(r render.Render, p martini.Params) {
//	images := make([]Image, 0)
//	rows, err := db.Query("SELECT data FROM images where message_id = ?", p["id"])
//	if err != nil {
//		log.Print(err)
//	}
//	for rows.Next() {
//		image := Image{}
//		err = rows.Scan(&image.Jpeg)
//		images = append(images, image)
//	}
//	r.JSON(200, images)
//}
//
//func main() {
//	database, err := sql.Open("mysql", "root:5550123Aa@tcp(localhost:3306)/waterobjects")
//	if err != nil {
//		log.Println(err)
//	}
//	defer func() {
//		_ = database.Close()
//	}()
//
//	db = database
//
//	m := martini.Classic()
//	m.Use(render.Renderer())
//	m.Use(cors.Allow(&cors.Options{
//		AllowMethods:     []string{"GET, PATCH, PUT, POST, DELETE, OPTIONS"},
//		AllowHeaders:     []string{"Content-Type"},
//		AllowAllOrigins:  true,
//		AllowCredentials: true,
//	}))
//	m.Get("/layers", layersHandler)
//	m.Get("/layers1", layersHandler1)
//	m.Get("/dictionaries", dictHandler)
//	m.Get("/dictionaries/:id", dictValuesHandler)
//	m.Get("/layers/:id", layerHandler)
//	m.Get("/layer/:id", layerHandler1)
//	m.Get("/fields/:id", fieldsHandler)
//	m.Get("/fields_values/:id", fieldValueHandler)
//	m.Get("/field_values/:ids", fieldValuesHandler)
//	m.Get("/check/:lat/:lon", checkHandler)
//	m.Get("/messages", messagesHandler)
//	m.Get("/new_messages", newMessagesHandler)
//	m.Get("/messages/:id", messageHandler)
//
//	m.Post("/", addLayersHandler)
//	m.Post("/layers", updateLayerHandler)
//	m.Post("/load", loadLayersHandler)
//	m.Post("/objects", loadObjectsHandler)
//	m.Post("/dictionaries", loadDictionariesHandler)
//	m.Post("/value", addDictValuesHandler)
//	m.Post("/object", objectPostHandler)
//	m.Post("/filter", filterHandler)
//	m.Post("/filter_str", filterStrHandler)
//	m.Post("/send", sendHandler)
//
//	m.Delete("/objects/:id", deleteObjectHandler)
//	m.Run()
//}
