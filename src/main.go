package main

import "github.com/ZootHii/todo-golang-backend/src/controllers"

func main() {
	a := controllers.App{}
	a.Initialize(
		"root",
		"1234",
		"todo_db")

	a.Run(":8010")
}
