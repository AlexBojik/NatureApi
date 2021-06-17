package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	_ "go.mozilla.org/pkcs7"
	"net/http"
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

var LayersHandler = func(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		utils.Respond(w, models.GetLayers())
		break
	case "POST":
		//CreateBaseLayer(w, r)
		break
	case "PUT":
		//UpdateBaseLayer(w, r)
		break
	}
}

var LayerHandler = func(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		id, _ := strconv.Atoi(mux.Vars(r)["id"])
		utils.Respond(w, models.GetLayer(id))
		break
	case "POST":
		//CreateBaseLayer(w, r)
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
	f := &models.Filter{}
	err := json.NewDecoder(r.Body).Decode(f)
	if err != nil {
		utils.Respond(w, utils.Message(false, "Error while decoding request body"))
		return
	}
	utils.Respond(w, models.GetFeaturesBy(f))
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
	utils.Respond(w, models.GetFieldsList())
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
	utils.Respond(w, models.GetUser(mux.Vars(r)["token"]))
}

var UserListHandler = func(w http.ResponseWriter, r *http.Request) {
	utils.Respond(w, models.GetUserList())
}

var MessageListHandler = func(w http.ResponseWriter, r *http.Request) {
	utils.Respond(w, models.GetMessageList())
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
		models.CreateUserGroup(ug)
	} else {
		models.UpdateUserGroup(ug)
	}
	utils.Created(w, int64(ug.Id))
}
