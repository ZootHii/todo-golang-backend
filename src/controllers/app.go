package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

func (a *App) Initialize(user, password, dbname string) {
	connectionString :=
		fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", user, password, dbname)

	var err error
	a.DB, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	a.Router = mux.NewRouter()

	a.InitializeRoutes()
}

func (a *App) Run(addr string) {

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:4200"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"},
	})

	handler := c.Handler(a.Router)

	log.Fatal(http.ListenAndServe(":8010", handler))
}

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func (a *App) InitializeRoutes() {

	a.Router.HandleFunc("/api/todos", a.GetTodos).Methods("GET")
	a.Router.HandleFunc("/api/todo", a.CreateTodo).Methods("POST")
	a.Router.HandleFunc("/api/todos/{id:[0-9]+}", a.GetTodo).Methods("GET")
	a.Router.HandleFunc("/api/todos/{id:[0-9]+}", a.UpdateTodo).Methods("PUT")
	a.Router.HandleFunc("/api/todos/{id:[0-9]+}", a.DeleteTodo).Methods("DELETE")
}
