package main

func main() {
	a := App{}
	a.Initialize(
		"root",
		"1234",
		"todo_db")

	a.Run(":8010")
}
