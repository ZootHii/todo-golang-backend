package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/ZootHii/todo-golang-backend/src/models"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"strconv"
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
	log.Fatal(http.ListenAndServe(":8010", a.Router))
}

func (a *App) GetTodo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		rm := models.ResponseModel{Success: false, Message: "Invalid todo ID"}
		RespondWithJSON(w, http.StatusBadRequest, rm)
		//respondWithError(w, http.StatusBadRequest, "Invalid todo ID")
		return
	}

	t := models.Todo{ID: id}
	if err := t.GetTodo(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			rm := models.ResponseModel{Success: false, Message: "Todo not found"}
			RespondWithJSON(w, http.StatusNotFound, rm)
			//respondWithError(w, http.StatusNotFound, "Todo not found")
		default:
			rm := models.ResponseModel{Success: false, Message: err.Error()}
			RespondWithJSON(w, http.StatusInternalServerError, rm)
			//respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	sdrm := models.SignleDataResponseModel{ResponseModel: models.ResponseModel{Success: true, Message: "data returned successfully"}, Data: t}
	RespondWithJSON(w, http.StatusOK, sdrm)
}

func RespondWithError(w http.ResponseWriter, code int, message string) {
	RespondWithJSON(w, code, map[string]string{"error": message})
}

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func (a *App) GetTodos(w http.ResponseWriter, r *http.Request) {

	//count, _ := strconv.Atoi(r.FormValue("count"))
	//start, _ := strconv.Atoi(r.FormValue("start"))

	/*if count > 10 || count < 1 {
		count = 10
	}
	if start < 0 {
		start = 0
	}*/

	todos, err := models.GetTodos(a.DB /*, start, count*/)
	if err != nil {
		rm := models.ResponseModel{Success: false, Message: err.Error()}
		RespondWithJSON(w, http.StatusInternalServerError, rm)
		return
	}

	drm := models.DataResponseModel{ResponseModel: models.ResponseModel{Success: true, Message: "datas returned successfully"}, Data: todos}
	RespondWithJSON(w, http.StatusOK, drm)
}

func (a *App) CreateTodo(w http.ResponseWriter, r *http.Request) {
	var t models.Todo
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&t); err != nil {
		rm := models.ResponseModel{Success: false, Message: "Invalid request payload"}
		RespondWithJSON(w, http.StatusBadRequest, rm)
		return
	}
	defer r.Body.Close()

	if err := t.CreateTodo(a.DB); err != nil {
		rm := models.ResponseModel{Success: false, Message: err.Error()}
		RespondWithJSON(w, http.StatusBadRequest, rm)
		return
	}

	if err := t.GetTodo(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			rm := models.ResponseModel{Success: false, Message: "Todo not found"}
			RespondWithJSON(w, http.StatusNotFound, rm)
			//respondWithError(w, http.StatusNotFound, "Todo not found")
		default:
			rm := models.ResponseModel{Success: false, Message: err.Error()}
			RespondWithJSON(w, http.StatusInternalServerError, rm)
			//respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	sdrm := models.SignleDataResponseModel{ResponseModel: models.ResponseModel{Success: true, Message: "data created successfully"}, Data: t}
	RespondWithJSON(w, http.StatusCreated, sdrm)
}

func (a *App) UpdateTodo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		rm := models.ResponseModel{Success: false, Message: "Invalid todo ID"}
		RespondWithJSON(w, http.StatusBadRequest, rm)
		return
	}

	t := models.Todo{ID: id}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&t); err != nil {
		rm := models.ResponseModel{Success: false, Message: "Invalid resquest payload"}
		RespondWithJSON(w, http.StatusBadRequest, rm)
		return
	}

	defer r.Body.Close()

	if err := t.UpdateTodo(a.DB); err != nil {
		rm := models.ResponseModel{Success: false, Message: err.Error()}
		RespondWithJSON(w, http.StatusBadRequest, rm)
		return
	}

	if err := t.GetTodo(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			rm := models.ResponseModel{Success: false, Message: "Todo not found"}
			RespondWithJSON(w, http.StatusNotFound, rm)
			//respondWithError(w, http.StatusNotFound, "Todo not found")
		default:
			rm := models.ResponseModel{Success: false, Message: err.Error()}
			RespondWithJSON(w, http.StatusInternalServerError, rm)
			//respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	sdrm := models.SignleDataResponseModel{ResponseModel: models.ResponseModel{Success: true, Message: "data updated successfully"}, Data: t}
	RespondWithJSON(w, http.StatusOK, sdrm)
}

func (a *App) DeleteTodo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		rm := models.ResponseModel{Success: false, Message: "Invalid todo ID"}
		RespondWithJSON(w, http.StatusBadRequest, rm)
		return
	}

	t := models.Todo{ID: id}
	if err := t.DeleteTodo(a.DB); err != nil {
		rm := models.ResponseModel{Success: false, Message: err.Error()}
		RespondWithJSON(w, http.StatusInternalServerError, rm)
		return
	}

	rm := models.ResponseModel{Success: true, Message: "data deleted successfully"}
	RespondWithJSON(w, http.StatusOK, rm)
}

func (a *App) InitializeRoutes() {
	a.Router.HandleFunc("/todos", a.GetTodos).Methods("GET")
	a.Router.HandleFunc("/todo", a.CreateTodo).Methods("POST")
	a.Router.HandleFunc("/todos/{id:[0-9]+}", a.GetTodo).Methods("GET")
	a.Router.HandleFunc("/todos/{id:[0-9]+}", a.UpdateTodo).Methods("PUT")
	a.Router.HandleFunc("/todos/{id:[0-9]+}", a.DeleteTodo).Methods("DELETE")
}
