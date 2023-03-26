package services

import (
	"net/http"

	"sync"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
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

func GetTodos(c *gin.Context, db *pgxpool.Pool) {
	c.IndentedJSON(http.StatusOK, todos)
}

var syncLock SyncLocks = SyncLocks{
	currentId: 0,
}

func abortWithMessage[T any](c *gin.Context, message T) {
	c.AbortWithStatusJSON(http.StatusBadRequest, message)
}

func AddTodo(c *gin.Context, db *pgxpool.Pool) {

	var todo Todo = Todo{
		Id:         syncLock.generateRandomId(),
		IsStriked:  false,
		IsComplete: false,
	}

	if err := c.BindJSON(&todo); err != nil {
		abortWithMessage(c, "Error while fetching the request body")
		return
	}

	todos = append(todos, todo)
	c.IndentedJSON(http.StatusOK, gin.H{"message": todo})
}

func UpdateTodo(c *gin.Context, db *pgxpool.Pool) {
	var todo Todo = Todo{}

	if err := c.BindJSON(&todo); err != nil {
		abortWithMessage(c, "Error while updating the todo or Id is wrong")
		return
	}

	found := s.IndexFunc(todos, func(t Todo) bool {
		return t.Id == todo.Id
	})

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

	abortWithMessage(c, todo)

}
