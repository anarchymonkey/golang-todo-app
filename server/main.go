package main

import (
	"fmt"
	"net/http"

	"sync"

	"github.com/gin-gonic/gin"
	s "golang.org/x/exp/slices"
)

type Todo struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	IsComplete bool   `json:"isComplete"`
	IsStriked  bool   `json:"isStriked"`
}

type SyncLocks struct {
	mu        sync.Mutex
	currentId int
}

func (m *SyncLocks) generateRandomId() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.currentId++
	return m.currentId
}

var todos []Todo = []Todo{}

func getTodos(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, todos)
}

var syncLock SyncLocks = SyncLocks{
	currentId: 0,
}

func addTodo(c *gin.Context) {

	var todo Todo = Todo{
		Id:         syncLock.generateRandomId(),
		IsStriked:  false,
		IsComplete: false,
	}
	if err := c.BindJSON(&todo); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "something went wrong"})
		return
	}
	todos = append(todos, todo)

	c.IndentedJSON(http.StatusOK, gin.H{"message": todo})
}

func updateTodo(c *gin.Context) {
	var todo Todo = Todo{}

	if err := c.BindJSON(&todo); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Wrong id provided or the request id sucks"})
		return
	}

	found := s.IndexFunc(todos, func(t Todo) bool {
		return t.Id == todo.Id
	})

	fmt.Println("The id found is", found)
	if found != -1 {

		for idx, _ := range todos {
			if idx == found {
				todos = s.Replace(todos, idx, idx+1, todo)
				break
			}
		}
		c.IndentedJSON(http.StatusOK, gin.H{
			"message": "Found this shit",
		})
		return
	}

	c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
		"message": todos,
	})

}

func main() {

	fmt.Println("Welcome to the todo web application")

	router := gin.Default()

	router.Use(corsMiddleware())

	router.GET("/todos", getTodos)
	router.POST("/todo/add", addTodo)
	router.PUT("/todo/update", updateTodo)

	// similar to http.ListenAndServe()
	router.Run(":8080")

}

func corsMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, mode")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusOK)
			return
		}
		c.Next()
	}

}
