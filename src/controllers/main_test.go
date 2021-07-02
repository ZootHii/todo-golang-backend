package controllers_test

import (
	/*"bytes"
	"encoding/json"*/
	"log"
	/*"math/rand"
	"net/http"*/
	//"net/http/httptest"
	"os"
	//"strconv"
	"testing"

	"github.com/ZootHii/todo-golang-backend/src/controllers"
	/*"github.com/ZootHii/todo-golang-backend/src/models"
	"github.com/stretchr/testify/require"*/)

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

/*func ExecuteRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)

	return rr
}

func AddTodos(count int) {
	if count < 1 {
		count = 1
	}

	for i := 0; i < count; i++ {
		a.DB.Exec("INSERT INTO todos(what_todo) VALUES($1)", "Todo "+strconv.Itoa(rand.Intn(1000)*(i+1)))
	}
}
*/
func EnsureTableExists() {
	if _, err := a.DB.Exec(tableCreationQuery); err != nil {
		log.Fatal(err)
	}
}

func DeleteAndRestartTable() {
	a.DB.Exec("TRUNCATE TABLE todos RESTART IDENTITY;")
}

/*func TestEmptyTable(t *testing.T) {
	DeleteAndRestartTable()

	req, _ := http.NewRequest("GET", "/api/todos", nil)
	response := ExecuteRequest(req)

	require.Equal(t, http.StatusOK, response.Code)

	drm := &models.DataResponseModel{}
	err := json.Unmarshal([]byte(response.Body.Bytes()), &drm)

	require.NoError(t, err)
	require.Empty(t, drm.Data)
}

func TestGetNonExistentTodo(t *testing.T) {
	DeleteAndRestartTable()

	req, _ := http.NewRequest("GET", "/api/todos/11", nil)
	response := ExecuteRequest(req)

	require.Equal(t, http.StatusNotFound, response.Code)

	rm := &models.ResponseModel{}
	err := json.Unmarshal([]byte(response.Body.Bytes()), rm)

	require.NoError(t, err)
	require.Equal(t, "Todo not found", rm.Message)

}

func TestCreateTodo(t *testing.T) {
	DeleteAndRestartTable()

	var jsonStr = []byte(`{"what_todo":"test todo"}`)
	req, _ := http.NewRequest("POST", "/api/todo", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response := ExecuteRequest(req)

	require.Equal(t, http.StatusCreated, response.Code)

	sdrm := &models.SignleDataResponseModel{}
	err := json.Unmarshal([]byte(response.Body.Bytes()), sdrm)

	require.NoError(t, err)
	require.Equal(t, "test todo", sdrm.Data.WhatTodo)
	require.Equal(t, 1, sdrm.Data.ID)
	require.NotZero(t, sdrm.Data.CreatedAt)

}

func TestGetTodo(t *testing.T) {
	DeleteAndRestartTable()
	AddTodos(1)

	req, _ := http.NewRequest("GET", "/api/todos/1", nil)
	response := ExecuteRequest(req)

	require.Equal(t, http.StatusOK, response.Code)

	sdrm := &models.SignleDataResponseModel{}
	err := json.Unmarshal([]byte(response.Body.Bytes()), sdrm)

	require.NoError(t, err)
	require.Equal(t, 1, sdrm.Data.ID)

}

func TestUpdateTodo(t *testing.T) {
	DeleteAndRestartTable()
	AddTodos(1)

	req, _ := http.NewRequest("GET", "/api/todos/1", nil)
	response := ExecuteRequest(req)

	originalSdrm := &models.SignleDataResponseModel{}
	json.Unmarshal(response.Body.Bytes(), &originalSdrm)

	var jsonStr = []byte(`{"what_todo":"test todo - updated what_todo"}`)
	req, _ = http.NewRequest("PUT", "/api/todos/1", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response = ExecuteRequest(req)

	require.Equal(t, http.StatusOK, response.Code)

	sdrm := &models.SignleDataResponseModel{}
	err := json.Unmarshal([]byte(response.Body.Bytes()), sdrm)

	require.NoError(t, err)
	require.Equal(t, originalSdrm.Data.ID, sdrm.Data.ID)
	require.NotEqual(t, originalSdrm.Data.WhatTodo, sdrm.Data.WhatTodo)

}

func TestDeleteTodo(t *testing.T) {
	DeleteAndRestartTable()
	AddTodos(1)

	req, _ := http.NewRequest("GET", "/api/todos/1", nil)
	response := ExecuteRequest(req)

	require.Equal(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("DELETE", "/api/todos/1", nil)
	response = ExecuteRequest(req)

	require.Equal(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("GET", "/api/todos/1", nil)
	response = ExecuteRequest(req)

	require.Equal(t, http.StatusNotFound, response.Code)
}
*/
