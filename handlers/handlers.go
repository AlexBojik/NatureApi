package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	_ "go.mozilla.org/pkcs7"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
	"water-api/esia"
	"water-api/models"
	"water-api/utils"
)

var BaseLayerHandler = func(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		utils.Respond(w, models.GetBaseLayers())
		break
	case "POST":
		CreateBaseLayer(w, r)
		break
	case "PUT":
		UpdateBaseLayer(w, r)
		break
	}
}

var DictionariesHandler = func(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		utils.Respond(w, models.GetDictionaries())
		break
	case "POST":
		CreateDictionary(w, r)
		break
	case "PUT":
		UpdateDictionary(w, r)
		break
	}
}

var ValuesHandler = func(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		utils.Respond(w, models.GetDictionariesValuesList())
		break
	case "POST":
		CreateValue(w, r)
		break
	case "PUT":
		UpdateValue(w, r)
		break
	}
}

var DictionariesIdHandler = func(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	switch r.Method {
	case "GET":
		utils.Respond(w, models.GetValues(id))
		break
	case "DELETE":
		models.DeleteDictionary(id)
		break
	}
}

var ValuesIdHandler = func(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	models.DeleteValue(id)
}

var LayersHandler = func(w http.ResponseWriter, r *http.Request) {
	hasInfoRole := models.HasInfoRole(r.Header.Get("Token"))
	switch r.Method {
	case "GET":
		utils.Respond(w, models.GetLayers(hasInfoRole))
		break
	case "POST":
		CreateLayer(w, r)
		break
	case "PUT":
		UpdateLayer(w, r)
		break
	}
}

var ObjectsHandler = func(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		break
	case "POST":
		CreateObject(w, r)
		break
	case "PUT":
		//UpdateObject(w, r)
		break
	}
}

var DumpsHandler = func(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		utils.Respond(w, models.GetDumpsNames())
		break
	case "POST":
		models.BackupDB()
		utils.Respond(w, utils.Message(true, "Backup created"))
		break
	}
}

var RestoreHandler = func(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		models.RestoreDB(mux.Vars(r)["name"])
		utils.Respond(w, utils.Message(true, "Restored backup"))
		break
	}
}

var TablesHandler = func(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		utils.Respond(w, models.GetTables())
		break
	}
}

var ObjectsIdHandler = func(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	switch r.Method {
	case "DELETE":
		models.DeleteObjectById(id)
		utils.Respond(w, utils.Message(true, "Delete success"))
		break
	}
}

func CreateObject(w http.ResponseWriter, r *http.Request) {
	o := &models.Object{}
	err := json.NewDecoder(r.Body).Decode(o)
	if err != nil {
		utils.Respond(w, utils.Message(false, "Error while decoding request body"))
		return
	}

	id := models.CreateObject(o)
	utils.Created(w, id)
}

func CreateLayer(w http.ResponseWriter, r *http.Request) {
	l := &models.Layer{}
	err := json.NewDecoder(r.Body).Decode(l)
	if err != nil {
		utils.Respond(w, utils.Message(false, "Error while decoding request body"))
		return
	}

	id := models.CreateLayer(l)
	utils.Created(w, id)
}

func CreateDictionary(w http.ResponseWriter, r *http.Request) {
	d := &models.Dictionary{}
	err := json.NewDecoder(r.Body).Decode(d)
	if err != nil {
		utils.Respond(w, utils.Message(false, "Error while decoding request body"))
		return
	}

	id := models.CreateDictionary(d)
	utils.Created(w, id)
}

func UpdateDictionary(w http.ResponseWriter, r *http.Request) {
	d := &models.Dictionary{}
	err := json.NewDecoder(r.Body).Decode(d)
	if err != nil {
		fmt.Println(err)
		utils.Respond(w, utils.Message(false, "Error while decoding request body"))
		return
	}
	models.UpdateDictionary(d)
	utils.Respond(w, utils.Message(true, "Update success"))
}

func CreateValue(w http.ResponseWriter, r *http.Request) {
	v := &models.Value{}
	err := json.NewDecoder(r.Body).Decode(v)
	if err != nil {
		utils.Respond(w, utils.Message(false, "Error while decoding request body"))
		return
	}

	id := models.CreateValue(v)
	utils.Created(w, id)
}

func UpdateValue(w http.ResponseWriter, r *http.Request) {
	v := &models.Value{}
	err := json.NewDecoder(r.Body).Decode(v)
	if err != nil {
		fmt.Println(err)
		utils.Respond(w, utils.Message(false, "Error while decoding request body"))
		return
	}
	models.UpdateValue(v)
	utils.Respond(w, utils.Message(true, "Update success"))
}

func UpdateLayer(w http.ResponseWriter, r *http.Request) {
	l := &models.Layer{}
	err := json.NewDecoder(r.Body).Decode(l)
	if err != nil {
		fmt.Println(err)
		utils.Respond(w, utils.Message(false, "Error while decoding request body"))
		return
	}
	models.UpdateLayer(l)
	utils.Respond(w, utils.Message(true, "Update success"))
}

var LayerHandler = func(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	switch r.Method {
	case "GET":
		utils.Respond(w, models.GetLayer(id))
		break
	case "DELETE":
		models.DeleteLayer(id)
		utils.Respond(w, utils.Message(true, "Delete success"))
		break
	case "PUT":
		//UpdateBaseLayer(w, r)
		break
	}
}

var ClusterHandler = func(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	utils.Respond(w, models.GetCluster(id))
}

var FilterHandler = func(w http.ResponseWriter, r *http.Request) {
	hasInfoRole := models.HasInfoRole(r.Header.Get("Token"))
	f := &models.Filter{}
	err := json.NewDecoder(r.Body).Decode(f)
	if err != nil {
		utils.Respond(w, utils.Message(false, "Error while decoding request body"))
		return
	}
	utils.Respond(w, models.GetFeaturesBy(f, hasInfoRole))
}

var CheckHandler = func(w http.ResponseWriter, r *http.Request) {
	ch := &models.Check{}
	err := json.NewDecoder(r.Body).Decode(ch)
	if err != nil {
		utils.Respond(w, utils.Message(false, "Error while decoding request body"))
		return
	}
	utils.Respond(w, models.CheckPosition(ch))
}

var FieldsHandler = func(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		utils.Respond(w, models.GetFieldsList())
		break
	case "POST":
		CreateField(w, r)
		break
	case "PUT":
		UpdateField(w, r)
		break
	}
}

var FieldsIdHandler = func(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	switch r.Method {
	case "GET":
		utils.Respond(w, models.GetFieldsByLayerId(id))
		break
	case "DELETE":
		models.DeleteFieldById(id)
		utils.Respond(w, utils.Message(true, "Delete success"))
		break
	}
}

var CoordinatesHandler = func(w http.ResponseWriter, r *http.Request) {
	c := &models.Coordinates{}
	err := json.NewDecoder(r.Body).Decode(c)
	if err != nil {
		utils.Respond(w, utils.Message(false, "Error while decoding request body"))
		return
	}
	models.UpdateCoordinates(c)
	utils.Respond(w, utils.Message(true, "Update success"))
}

var AuthHandler = func(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, esia.GetAuthUrl(true), 302)
}

var WebAuthHandler = func(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, esia.GetAuthUrl(false), 302)
}

var EsiaCodeHandler = func(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, esia.GetRedirectUrl(r, true), 302)
}

var WebEsiaCodeHandler = func(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, esia.GetRedirectUrl(r, false), 302)
}

var UserHandler = func(w http.ResponseWriter, r *http.Request) {
	token := mux.Vars(r)["token"]
	switch r.Method {
	case "GET":
		utils.Respond(w, models.GetUser(token))
	case "DELETE":
		models.DeleteUser(token)
		utils.Respond(w, utils.Message(true, "Delete success"))
	}
}

var UserListHandler = func(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	utils.Respond(w, models.GetUserList(id))
}

var MessageListHandler = func(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		utils.Respond(w, models.GetMessageList())
		break
	case "POST":
		//CreateBaseLayer(w, r)
		break
	case "PUT":
		UpdateMessage(w, r)
		break
	}
}

var UpdateMessage = func(w http.ResponseWriter, r *http.Request) {
	msg := &models.Message{}
	err := json.NewDecoder(r.Body).Decode(msg)
	if err != nil {
		fmt.Println(err)
		utils.Respond(w, utils.Message(false, "Error while decoding request body"))
		return
	}

	models.UpdateMessage(msg)
	utils.Created(w, msg.Id)
}

var MessageHandler = func(w http.ResponseWriter, r *http.Request) {
	utils.Respond(w, models.PostMessage(r))
}

var MessageGetHandler = func(w http.ResponseWriter, r *http.Request) {
	utils.Respond(w, models.GetMessage(mux.Vars(r)["id"]))
}

var NewMessageHandler = func(w http.ResponseWriter, r *http.Request) {
	utils.Respond(w, models.GetMessagesCount())
}

var SendMessageHandler = func(w http.ResponseWriter, r *http.Request) {
	utils.Respond(w, models.SendMessages())
}

var UserPutHandler = func(w http.ResponseWriter, r *http.Request) {
	usr := &models.User{}
	err := json.NewDecoder(r.Body).Decode(usr)
	if err != nil {
		utils.Respond(w, utils.Message(false, "Error while decoding request body"))
		return
	}

	models.UpdateUser(usr)
	utils.Created(w, 0)
}

var UserCreateHandler = func(w http.ResponseWriter, r *http.Request) {
	usr := &models.User{}
	err := json.NewDecoder(r.Body).Decode(usr)
	if err != nil {
		utils.Respond(w, utils.Message(false, "Error while decoding request body"))
		return
	}

	models.CreateUser(usr)
	utils.Created(w, 0)
}

var CreateBaseLayer = func(w http.ResponseWriter, r *http.Request) {
	bl := &models.BaseLayer{}
	err := json.NewDecoder(r.Body).Decode(bl)
	if err != nil {
		utils.Respond(w, utils.Message(false, "Error while decoding request body"))
		return
	}

	id := models.CreateBaseLayer(bl)
	utils.Created(w, id)
}

var BaseLayerHandlerDelete = func(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	models.DeleteBaseLayer(id)
	utils.Respond(w, utils.Message(true, "Delete success"))
}

var UpdateBaseLayer = func(w http.ResponseWriter, r *http.Request) {
	bl := &models.BaseLayer{}
	err := json.NewDecoder(r.Body).Decode(bl)
	if err != nil {
		fmt.Println(err)
		utils.Respond(w, utils.Message(false, "Error while decoding request body"))
		return
	}
	models.UpdateBaseLayer(bl)
	utils.Respond(w, utils.Message(true, "Update success"))
}

var NewsHandler = func(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		utils.Respond(w, models.GetNews(false))
		break
	case "POST":
		CreateNews(w, r)
		break
	case "PUT":
		UpdateNews(w, r)
		break
	}
}

var FilesHandler = func(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		utils.Respond(w, models.GetFiles())
		break
	case "POST":
		uploadHandler(w, r)
		break
	case "DELETE":
		id, _ := strconv.Atoi(mux.Vars(r)["id"])
		models.DeleteFile(id)
		break
	}
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	file, handler, err := r.FormFile("fileKey")
	if err != nil {
		fmt.Println(err)
	}

	defer file.Close()
	os.Mkdir("files", 0777)

	f, err := os.OpenFile("files/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	//_, _ = io.WriteString(w, "File "+fileName+" Uploaded successfully")
	_, _ = io.Copy(f, file)

	models.CreateFile(handler.Filename)
	utils.Respond(w, utils.Message(true, "Upload success"))
}

var NewsFilteredListHandler = func(w http.ResponseWriter, r *http.Request) {
	utils.Respond(w, models.GetNews(true))
}

var CreateNews = func(w http.ResponseWriter, r *http.Request) {
	nw := &models.News{}
	err := json.NewDecoder(r.Body).Decode(nw)
	if err != nil {
		utils.Respond(w, utils.Message(false, "Error while decoding request body"))
		return
	}
	nw.Created = time.Now()
	id := models.CreateNews(nw)
	utils.Created(w, id)
}

var NewsHandlerDelete = func(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	models.DeleteNews(id)
	utils.Respond(w, utils.Message(true, "Delete success"))
}

var UpdateNews = func(w http.ResponseWriter, r *http.Request) {
	nw := &models.News{}
	err := json.NewDecoder(r.Body).Decode(nw)
	if err != nil {
		fmt.Println(err)
		utils.Respond(w, utils.Message(false, "Error while decoding request body"))
		return
	}
	models.UpdateNews(nw)
	utils.Respond(w, utils.Message(true, "Update success"))
}

var UserGroupsHandler = func(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		utils.Respond(w, models.GetUserGroups())
		break
	case "POST":
		PostUserGroup(w, r)
		break
	case "PUT":
		//UpdateNews(w, r)
		break
	}
}

var PostUserGroup = func(w http.ResponseWriter, r *http.Request) {
	ug := &models.UserGroups{}
	err := json.NewDecoder(r.Body).Decode(ug)
	if err != nil {
		utils.Respond(w, utils.Message(false, "Error while decoding request body"))
		return
	}
	if ug.Id == 0 {
		ug.Id = int(models.CreateUserGroup(ug))
	} else {
		models.UpdateUserGroup(ug)
	}
	utils.Created(w, int64(ug.Id))
}

func CreateField(w http.ResponseWriter, r *http.Request) {
	f := &models.Field{}
	err := json.NewDecoder(r.Body).Decode(f)
	if err != nil {
		utils.Respond(w, utils.Message(false, "Error while decoding request body"))
		return
	}

	id := models.CreateField(f)
	utils.Created(w, id)
}

func UpdateField(w http.ResponseWriter, r *http.Request) {
	f := &models.Field{}
	err := json.NewDecoder(r.Body).Decode(f)
	if err != nil {
		fmt.Println(err)
		utils.Respond(w, utils.Message(false, "Error while decoding request body"))
		return
	}
	models.UpdateField(f)
	utils.Respond(w, utils.Message(true, "Update success"))
}
