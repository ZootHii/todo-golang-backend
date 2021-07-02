package controllers_test

import (
	"bytes"
	"encoding/json"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/ZootHii/todo-golang-backend/src/models"
	"github.com/stretchr/testify/require"
)

func ExecuteRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()

	a.Router.ServeHTTP(rr, req)

	return rr
}

func createRandomTodo(t *testing.T) (models.SignleDataResponseModel, int, string, error) {
	rand.Seed(time.Now().UnixNano())
	randomWhatTodo := ("Todo " + strconv.Itoa(rand.Int()))
	testPostTodo := &models.Todo{WhatTodo: randomWhatTodo}

	jsonStr, _ := json.Marshal(testPostTodo)
	req, _ := http.NewRequest("POST", "/api/todo", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response := ExecuteRequest(req)
	sdrm := &models.SignleDataResponseModel{}
	err := json.Unmarshal([]byte(response.Body.Bytes()), sdrm)

	return *sdrm, response.Code, randomWhatTodo, err

	//a.DB.Exec("INSERT INTO todos(what_todo) VALUES($1)", "Todo "+strconv.Itoa(rand.Intn(1000)*( /*i+*/ 1)))
}

func TestGetTodos(t *testing.T) {
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
	randomTodo, responseCode, randomWhatTodo, errR := createRandomTodo(t)

	/*var jsonStr = []byte(`{"what_todo":"test todo"}`)
	req, _ := http.NewRequest("POST", "/api/todo", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response := ExecuteRequest(req)*/

	//require.Equal(t, http.StatusCreated, response.Code)

	require.Equal(t, http.StatusCreated, responseCode)

	//sdrm := &models.SignleDataResponseModel{}
	//err := json.Unmarshal([]byte(response.Body.Bytes()), sdrm)

	require.NoError(t, errR)
	require.Equal(t, randomWhatTodo, randomTodo.Data.WhatTodo)
	require.Equal(t, 1, randomTodo.Data.ID)
	require.NotZero(t, randomTodo.Data.CreatedAt)

}

func TestGetTodo(t *testing.T) {
	DeleteAndRestartTable()
	randomTodo, _, _, _ := createRandomTodo(t)

	//require.Equal(t, http.StatusOK, returnTodo.Data)
	req, _ := http.NewRequest("GET", "/api/todos/"+strconv.Itoa(randomTodo.Data.ID), nil)
	response := ExecuteRequest(req)

	require.Equal(t, http.StatusOK, response.Code)

	sdrm := &models.SignleDataResponseModel{}
	err := json.Unmarshal([]byte(response.Body.Bytes()), sdrm)

	require.NoError(t, err)
	require.Equal(t, randomTodo.Data.ID, sdrm.Data.ID)
	require.Equal(t, randomTodo.Data, sdrm.Data)

}

func TestUpdateTodo(t *testing.T) {
	DeleteAndRestartTable()
	randomTodo, _, _, _ := createRandomTodo(t)

	req, _ := http.NewRequest("GET", "/api/todos/"+strconv.Itoa(randomTodo.Data.ID), nil)
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
	randomTodo, _, _, _ := createRandomTodo(t)

	req, _ := http.NewRequest("GET", "/api/todos/"+strconv.Itoa(randomTodo.Data.ID), nil)
	response := ExecuteRequest(req)

	require.Equal(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("DELETE", "/api/todos/"+strconv.Itoa(randomTodo.Data.ID), nil)
	response = ExecuteRequest(req)

	require.Equal(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("GET", "/api/todos/"+strconv.Itoa(randomTodo.Data.ID), nil)
	response = ExecuteRequest(req)

	require.Equal(t, http.StatusNotFound, response.Code)
}
