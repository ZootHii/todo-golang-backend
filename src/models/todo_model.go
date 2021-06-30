package models

import (
	"database/sql"
	"sort"
)

type Todo struct {
	ID        int    `json:"id"`
	WhatTodo  string `json:"what_todo"`
	CreatedAt string `json:"created_at"`
}

func (t *Todo) GetTodo(db *sql.DB) error {
	return db.QueryRow("SELECT what_todo, created_at FROM todos WHERE id=$1",
		t.ID).Scan(&t.WhatTodo, &t.CreatedAt)
}

func (t *Todo) UpdateTodo(db *sql.DB) error {
	_, err :=
		db.Exec("UPDATE todos SET what_todo=$1 WHERE id=$2",
			t.WhatTodo, t.ID)

	return err
}

func (t *Todo) DeleteTodo(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM todos WHERE id=$1", t.ID)

	return err
}

func (t *Todo) CreateTodo(db *sql.DB) error {
	err := db.QueryRow(
		"INSERT INTO todos(what_todo) VALUES($1) RETURNING id",
		t.WhatTodo).Scan(&t.ID)

	if err != nil {
		return err
	}

	return nil
}

func GetTodos(db *sql.DB /*, start, count int*/) ([]Todo, error) {
	/*rows, err := db.Query(
	"SELECT id, what_todo, created_at FROM todos LIMIT $1 OFFSET $2",
	count, start)*/

	rows, err := db.Query("SELECT id, what_todo, created_at FROM todos")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	todos := []Todo{}

	for rows.Next() {
		var t Todo
		if err := rows.Scan(&t.ID, &t.WhatTodo, &t.CreatedAt); err != nil {
			return nil, err
		}
		todos = append(todos, t)

		// sort starting from latest created
		sort.Slice(todos, func(i, j int) bool {
			return todos[i].CreatedAt > todos[j].CreatedAt
		})
	}

	return todos, nil
}
