package controllers_test

import (
	"log"

	"os"
	"testing"

	"github.com/ZootHii/todo-golang-backend/src/controllers"
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

func EnsureTableExists() {
	if _, err := a.DB.Exec(tableCreationQuery); err != nil {
		log.Fatal(err)
	}
}

func DeleteAndRestartTable() {
	a.DB.Exec("TRUNCATE TABLE todos RESTART IDENTITY;")
}
