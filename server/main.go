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

func (m *SyncLocks) generateRandomId() {
	m.mu.Lock()
	m.currentId = m.currentId + 1
	m.mu.Unlock()
}

var todos []Todo = []Todo{
	Todo{
		Id:         1,
		Name:       "This is a test todo",
		IsComplete: false,
		IsStriked:  false,
	},
}

func getTodos(c *gin.Context) {

	fmt.Println("todos", todos)

	c.IndentedJSON(http.StatusOK, todos)
}

// this is outside so that it keeps the context of random Id, I will replace it by postgres but for now I dont have the headache of generating ids anymore
var newSyncLock SyncLocks = SyncLocks{
	currentId: 2,
}

func addTodo(c *gin.Context) {
	// get the request and read it, add it to the todos array

	go newSyncLock.generateRandomId()

	var todo Todo = Todo{
		Id:         newSyncLock.currentId,
		IsStriked:  false,
		IsComplete: false,
	}
	if err := c.BindJSON(&todo); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "something went wrong"})
		return
	}

	todos = append(todos, todo)

	c.IndentedJSON(http.StatusOK, gin.H{"message": "Successfully added a todo"})
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
