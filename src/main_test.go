package main_test

import (
	"bytes"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"github.com/ZootHii/todo-golang-backend/src/controllers"
	"github.com/ZootHii/todo-golang-backend/src/models"
)

var a controllers.App

func TestMain(m *testing.M) {
	a = controllers.App{}
	a.Initialize(
		"root",
		"1234",
		"todo_db")

	EnsureTableExists()
	DeleteAndRestartTable()
	code := m.Run()

	DeleteAndRestartTable()

	os.Exit(code)
}

const tableCreationQuery = `CREATE TABLE IF NOT EXISTS todos 
(
	id bigserial PRIMARY KEY,
	what_todo varchar NOT NULL,
	created_at timestamptz NOT NULL DEFAULT (now())

);`

func ExecuteRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)

	return rr
}

func CheckResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected : %d\nActual : %d\n", expected, actual)
	}
}

func AddTodos(count int) {
	if count < 1 {
		count = 1
	}

	for i := 0; i < count; i++ {
		a.DB.Exec("INSERT INTO todos(what_todo) VALUES($1)", "Todo "+strconv.Itoa(rand.Intn(1000)*(i+1)))
	}
}

func EnsureTableExists() {
	if _, err := a.DB.Exec(tableCreationQuery); err != nil {
		log.Fatal(err)
	}
}

func DeleteAndRestartTable() {
	a.DB.Exec("TRUNCATE TABLE todos RESTART IDENTITY;")
}

func TestEmptyTable(t *testing.T) {
	DeleteAndRestartTable()

	req, _ := http.NewRequest("GET", "/todos", nil)
	response := ExecuteRequest(req)

	CheckResponseCode(t, http.StatusOK, response.Code)

	body := response.Body.String()
	//var data map[string]interface{}

	drm := &models.DataResponseModel{}
	err := json.Unmarshal([]byte(body), &drm)
	if err != nil {
		panic(err)
	}

	if len(drm.Data) != 0 {
		t.Errorf("Expected : 0\nGot : %d", len(drm.Data))
	}
}

func TestGetNonExistentTodo(t *testing.T) {
	DeleteAndRestartTable()

	req, _ := http.NewRequest("GET", "/todos/11", nil)
	response := ExecuteRequest(req)

	CheckResponseCode(t, http.StatusNotFound, response.Code)

	rm := &models.ResponseModel{}
	err := json.Unmarshal([]byte(response.Body.Bytes()), rm)
	if err != nil {
		log.Fatal(err)
	}

	if rm.Message != "Todo not found" {
		t.Errorf("Expected : message to be Todo not found\nGot : '%s'", rm.Message)
	}

	/*var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["message"] != "Todo not found" {
		t.Errorf("Expected the 'error' key of the response to be set to 'Todo not found'. Got '%s'", m["message"])
	}*/
}

func TestCreateTodo(t *testing.T) {

	DeleteAndRestartTable()

	var jsonStr = []byte(`{"what_todo":"test todo"}`)
	req, _ := http.NewRequest("POST", "/todo", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response := ExecuteRequest(req)
	CheckResponseCode(t, http.StatusCreated, response.Code)

	sdrm := &models.SignleDataResponseModel{}
	err := json.Unmarshal([]byte(response.Body.Bytes()), sdrm)
	if err != nil {
		log.Fatal(err)
	}

	if sdrm.Data.WhatTodo != "test todo" {
		t.Errorf("Expected : what_todo to be 'test todo'\nGot : '%v'", sdrm.Data.WhatTodo)
	}

	if sdrm.Data.ID != 1.0 {
		t.Errorf("Expected : id to be '1'\nGot : '%v'", sdrm.Data.ID)
	}
}

func TestGetTodo(t *testing.T) {
	DeleteAndRestartTable()
	AddTodos(1)

	req, _ := http.NewRequest("GET", "/todos/1", nil)
	response := ExecuteRequest(req)

	sdrm := &models.SignleDataResponseModel{}
	err := json.Unmarshal([]byte(response.Body.Bytes()), sdrm)
	if err != nil {
		log.Fatal(err)
	}

	if sdrm.Data.ID != 1.0 {
		t.Errorf("Expected : id to be '1'\nGot : '%v'", sdrm.Data.ID)
	}

	CheckResponseCode(t, http.StatusOK, response.Code)

}

func TestUpdateTodo(t *testing.T) {

	DeleteAndRestartTable()
	AddTodos(1)

	req, _ := http.NewRequest("GET", "/todos/1", nil)
	response := ExecuteRequest(req)
	//var originalTodo map[string]interface{}

	originalSdrm := &models.SignleDataResponseModel{}
	json.Unmarshal(response.Body.Bytes(), &originalSdrm)

	var jsonStr = []byte(`{"what_todo":"test todo - updated what_todo"}`)
	req, _ = http.NewRequest("PUT", "/todos/1", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response = ExecuteRequest(req)

	CheckResponseCode(t, http.StatusOK, response.Code)

	//var m map[string]interface{}
	sdrm := &models.SignleDataResponseModel{}
	json.Unmarshal(response.Body.Bytes(), &sdrm)

	if sdrm.Data.ID != originalSdrm.Data.ID {
		t.Errorf("Expected : id to remain the same (%v)\nGot : %v", originalSdrm.Data.ID, sdrm.Data.ID)
	}

	if sdrm.Data.WhatTodo == originalSdrm.Data.WhatTodo {
		t.Errorf("Expected : what_todo to change from '%v' to '%v'\nGot : '%v'", originalSdrm.Data.WhatTodo, sdrm.Data.WhatTodo, sdrm.Data.WhatTodo)
	}
}

func TestDeleteTodo(t *testing.T) {
	DeleteAndRestartTable()
	AddTodos(1)

	req, _ := http.NewRequest("GET", "/todos/1", nil)
	response := ExecuteRequest(req)
	CheckResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("DELETE", "/todos/1", nil)
	response = ExecuteRequest(req)

	CheckResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("GET", "/todos/1", nil)
	response = ExecuteRequest(req)
	CheckResponseCode(t, http.StatusNotFound, response.Code)
}
