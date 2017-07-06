package app

import (
	"github.com/gorilla/mux"
	"net/http"
	"github.com/Zenika/MARCEL/backend/medias"
	"github.com/rs/cors"
	"log"
	"os"
	//"github.com/Zenika/MARCEL/backend/plugins"
	"github.com/Zenika/MARCEL/backend/apidoc"
)

//current version of the API
const MARCEL_API_VERSION = "1"
var logFile string = os.Getenv("MARCEL_LOG_FILE")

type App struct {
	Router http.Handler
}

func (a *App) Initialize() {

	a.initializeLog()

	a.initializeData()

	a.initializeRoutes()
}

func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.Router))
	log.Printf("Server is started and listening on port %v", addr)
}

func (a *App) initializeRoutes() {

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		// AllowedOrigins:   []string{"http://localhost:*"},
		AllowedMethods:   []string{"GET", "POST", "DELETE", "OPTION", "PUT"},
		AllowCredentials: true,
	})

	r := mux.NewRouter()
	s := r.PathPrefix("/api/v" + MARCEL_API_VERSION).Subrouter()
	s.HandleFunc("/medias", medias.GetAllHandler).Methods("GET")
	s.HandleFunc("/medias/{idMedia}", medias.GetHandler).Methods("GET")
	s.HandleFunc("/medias/{idMedia}", medias.PostHandler).Methods("POST")
	s.HandleFunc("/medias/create", medias.CreateHandler).Methods("GET")
	r.HandleFunc("/swagger.json", apidoc.GetConfigHandler).Methods("GET")

	a.Router = c.Handler(r)
}

func (a* App) initializeLog() {
	if len(logFile) == 0 {
		logFile = "marcel.log"
	}
	f, err := os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	log.SetOutput(f)
}

func (a* App) initializeData() {

	//Load plugins list from DB
	//plugins.LoadPluginsCatalog()

	//Load Medias configuration from DB
	medias.LoadMedias()
}

