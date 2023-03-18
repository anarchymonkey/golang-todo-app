package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type Todos struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	IsComplete bool   `json:"isComplete"`
	IsStriked  bool   `json:"isStriked"`
}

type TodoJson struct {
	Todos []Todos `json:"todos"`
}

func getTodos(c *gin.Context) {
	fileData, err := os.ReadFile("./todos.json")

	if err != nil {
		panic("An error has occoured, lets start panicing")
	}
	var todos TodoJson = TodoJson{}
	err = json.Unmarshal(fileData, &todos)
	fmt.Println("todos", todos)

	if err != nil {
		panic(fmt.Errorf("there was an error while deserializing the file %s", err))
	}

	c.IndentedJSON(http.StatusOK, todos.Todos)
}

func main() {
	fmt.Println("Welcome to the todo web application")

	router := gin.Default()

	router.GET("/todos", getTodos)
	router.Run(":8080")

}
