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

func createRandomTodo(t *testing.T) (models.SignleDataResponseModel, int, string) {
	rand.Seed(time.Now().UnixNano())
	randomWhatTodo := ("Todo " + strconv.Itoa(rand.Int()))
	testPostTodo := &models.Todo{WhatTodo: randomWhatTodo}

	jsonStr, errMarshal := json.Marshal(testPostTodo)
	req, errRequest := http.NewRequest("POST", "/api/todo", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response := ExecuteRequest(req)
	sdrm := &models.SignleDataResponseModel{}
	errUnmarshal := json.Unmarshal([]byte(response.Body.Bytes()), sdrm)

	require.NoError(t, errMarshal)
	require.NoError(t, errRequest)
	require.NoError(t, errUnmarshal)

	return *sdrm, response.Code, randomWhatTodo
}

func TestGetTodos(t *testing.T) {
	DeleteAndRestartTable()

	for i := 0; i < 10; i++ {
		createRandomTodo(t)
	}

	req, errRequest := http.NewRequest("GET", "/api/todos", nil)
	response := ExecuteRequest(req)

	require.Equal(t, http.StatusOK, response.Code)

	drm := &models.DataResponseModel{}
	errUnmarshal := json.Unmarshal([]byte(response.Body.Bytes()), &drm)

	require.NoError(t, errRequest)
	require.NoError(t, errUnmarshal)
	require.NotEmpty(t, drm.Data)
	require.Len(t, drm.Data, 10)
	for _, todo := range drm.Data {
		require.NotEmpty(t, todo)
	}
}

func TestGetNonExistentTodo(t *testing.T) {
	DeleteAndRestartTable()

	req, errRequest := http.NewRequest("GET", "/api/todos/11", nil)
	response := ExecuteRequest(req)

	require.Equal(t, http.StatusNotFound, response.Code)

	rm := &models.ResponseModel{}
	errUnmarshal := json.Unmarshal([]byte(response.Body.Bytes()), rm)

	require.NoError(t, errRequest)
	require.NoError(t, errUnmarshal)
	require.Equal(t, "Todo not found", rm.Message)

}

func TestCreateTodo(t *testing.T) {
	DeleteAndRestartTable()
	randomTodo, responseCode, randomWhatTodo := createRandomTodo(t)

	require.Equal(t, http.StatusCreated, responseCode)
	require.Equal(t, randomWhatTodo, randomTodo.Data.WhatTodo)
	require.Equal(t, 1, randomTodo.Data.ID)
	require.NotZero(t, randomTodo.Data.CreatedAt)

}

func TestGetTodo(t *testing.T) {
	DeleteAndRestartTable()
	randomTodo, _, _ := createRandomTodo(t)

	req, errRequest := http.NewRequest("GET", "/api/todos/"+strconv.Itoa(randomTodo.Data.ID), nil)
	response := ExecuteRequest(req)
	sdrm := &models.SignleDataResponseModel{}
	errUnmarshal := json.Unmarshal([]byte(response.Body.Bytes()), sdrm)

	require.Equal(t, http.StatusOK, response.Code)
	require.NotEmpty(t, sdrm.Data)
	require.NoError(t, errRequest)
	require.NoError(t, errUnmarshal)
	require.Equal(t, randomTodo.Data.ID, sdrm.Data.ID)
	require.Equal(t, randomTodo.Data.WhatTodo, sdrm.Data.WhatTodo)
	require.Equal(t, randomTodo.Data.CreatedAt, sdrm.Data.CreatedAt)

}

func TestUpdateTodo(t *testing.T) {
	DeleteAndRestartTable()
	randomTodo, _, _ := createRandomTodo(t)

	var jsonStr = []byte(`{"what_todo":"test todo - updated what_todo"}`)
	req, errRequest := http.NewRequest("PUT", "/api/todos/"+strconv.Itoa(randomTodo.Data.ID), bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	response := ExecuteRequest(req)
	sdrm := &models.SignleDataResponseModel{}
	errUnmarshal := json.Unmarshal([]byte(response.Body.Bytes()), sdrm)

	require.Equal(t, http.StatusOK, response.Code)
	require.NoError(t, errRequest)
	require.NoError(t, errUnmarshal)
	require.NotEmpty(t, sdrm.Data)
	require.Equal(t, randomTodo.Data.ID, sdrm.Data.ID)
	require.NotEqual(t, randomTodo.Data.WhatTodo, sdrm.Data.WhatTodo)
	require.Equal(t, randomTodo.Data.CreatedAt, sdrm.Data.CreatedAt)

}

func TestDeleteTodo(t *testing.T) {
	DeleteAndRestartTable()
	randomTodo, responseCode, _ := createRandomTodo(t)
	require.Equal(t, http.StatusCreated, responseCode)

	req, errRequest := http.NewRequest("DELETE", "/api/todos/"+strconv.Itoa(randomTodo.Data.ID), nil)
	response := ExecuteRequest(req)

	require.NoError(t, errRequest)
	require.Equal(t, http.StatusOK, response.Code)

	req, errRequest = http.NewRequest("GET", "/api/todos/"+strconv.Itoa(randomTodo.Data.ID), nil)
	response = ExecuteRequest(req)

	require.NoError(t, errRequest)
	require.Equal(t, http.StatusNotFound, response.Code)
}
