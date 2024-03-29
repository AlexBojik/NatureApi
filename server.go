package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	http "net/http"
	"os"
	h "water-api/handlers"
	"water-api/models"
)

func main() {
	router := mux.NewRouter()

	// BaseLayers
	router.HandleFunc("/base_layers", h.BaseLayerHandler).Methods("GET", "POST", "PUT")
	router.HandleFunc("/base_layers/{id:[0-9]+}", h.BaseLayerHandlerDelete).Methods("DELETE")

	// Dictionaries
	router.HandleFunc("/dictionaries", h.DictionariesHandler).Methods("GET", "POST", "PUT")
	router.HandleFunc("/values", h.ValuesHandler).Methods("GET", "POST", "PUT")
	router.HandleFunc("/dictionaries/{id:[0-9]+}", h.DictionariesIdHandler).Methods("GET", "DELETE")
	router.HandleFunc("/values/{id:[0-9]+}", h.ValuesIdHandler).Methods("DELETE")

	// Layers
	router.HandleFunc("/layers", h.LayersHandler).Methods("GET", "POST", "PUT")
	router.HandleFunc("/layers/{id:[0-9]+}", h.LayerHandler).Methods("GET", "DELETE")
	router.HandleFunc("/cluster/{id:[0-9]+}", h.ClusterHandler).Methods("GET")

	// filter
	router.HandleFunc("/filter", h.FilterHandler).Methods("POST")

	// additional fields
	router.HandleFunc("/fields", h.FieldsHandler).Methods("GET", "POST", "PUT")
	router.HandleFunc("/fields/{id:[0-9]+}", h.FieldsIdHandler).Methods("GET", "DELETE")

	//coordinates
	router.HandleFunc("/coordinates", h.CoordinatesHandler).Methods("PUT")

	// check
	router.HandleFunc("/check", h.CheckHandler).Methods("POST")

	// user
	router.HandleFunc("/user/{token}", h.UserHandler).Methods("GET", "DELETE")
	router.HandleFunc("/user", h.UserCreateHandler).Methods("POST")
	router.HandleFunc("/user_put", h.UserPutHandler).Methods("POST")
	router.HandleFunc("/user_list/{id:[0-9]+}", h.UserListHandler).Methods("GET")
	router.HandleFunc("/user_groups", h.UserGroupsHandler).Methods("GET", "POST", "PUT")

	//messages
	router.HandleFunc("/send", h.MessageHandler).Methods("POST")
	router.HandleFunc("/new_messages", h.NewMessageHandler).Methods("GET")
	router.HandleFunc("/send_messages", h.SendMessageHandler).Methods("GET")
	router.HandleFunc("/messages", h.MessageListHandler).Methods("GET", "PUT")
	router.HandleFunc("/messages/{id}", h.MessageGetHandler).Methods("GET")

	// esia
	router.HandleFunc("/auth", h.AuthHandler).Methods("GET")
	router.HandleFunc("/webauth", h.WebAuthHandler).Methods("GET")
	router.HandleFunc("/esia", h.EsiaCodeHandler).Methods("GET")
	router.HandleFunc("/webesia", h.WebEsiaCodeHandler).Methods("GET")

	// news
	router.HandleFunc("/news", h.NewsHandler).Methods("GET", "POST", "PUT")
	router.HandleFunc("/news_list", h.NewsFilteredListHandler).Methods("GET")
	router.HandleFunc("/news/{id:[0-9]+}", h.NewsHandlerDelete).Methods("DELETE")

	// objects
	router.HandleFunc("/objects", h.ObjectsHandler).Methods("POST", "PUT")
	router.HandleFunc("/objects/{id:[0-9]+}", h.ObjectsIdHandler).Methods("DELETE")

	// dumps
	router.HandleFunc("/dumps", h.DumpsHandler).Methods("GET", "POST")
	router.HandleFunc("/dumps/{id:[0-9]+}", h.DumpsHandler).Methods("DELETE")

	// restore
	router.HandleFunc("/restore/{name}", h.RestoreHandler).Methods("GET")
	router.HandleFunc("/tables", h.TablesHandler).Methods("GET")

	//files
	router.HandleFunc("/files", h.FilesHandler).Methods("GET", "POST")
	router.HandleFunc("/files/{id:[0-9]+}", h.FilesHandler).Methods("DELETE")

	router.PathPrefix("/image/").Handler(http.StripPrefix("/image/", http.FileServer(http.Dir("images"))))
	router.PathPrefix("/dump/").Handler(http.StripPrefix("/dump/", http.FileServer(http.Dir("dumps"))))
	router.PathPrefix("/file/").Handler(http.StripPrefix("/file/", http.FileServer(http.Dir("files"))))

	os.Mkdir("images", 0777)
	port := os.Getenv("PORT")
	fmt.Println("Server is listening at port: ", port)


	cors := handlers.CORS(
		handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization", "Token"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS", "DELETE"}),
		handlers.AllowedOrigins([]string{"*"}))(router)

	http.ListenAndServe(":"+port, cors)
	defer models.CloseDB()
}
