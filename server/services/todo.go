package services

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
)

// List of groups present todo_app.groups
type Groups struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	isActive    bool   `json:"is_active"`
	isPublic    bool   `json:"is_public"`
}

// resembles todo_app.items
type Item struct {
	Id        int    `json:"id"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	RemindAt  string `json:"remind_at"`
	isActive  bool   `json:"is_active"`
}

type GroupedItems struct {
	Id        int    `json:"id"`
	GroupId   int    `json:"group_id"`
	ItemId    int    `json:"item_id"`
	CreatedAt string `json:"created_at"`
	IsActive  bool   `json:"is_active"`
}

func abortWithMessage[T any](c *gin.Context, message T) {
	c.AbortWithStatusJSON(http.StatusBadRequest, message)
}

func GetGroups(c *gin.Context, conn *pgxpool.Conn) {

}

func AddGroup(c *gin.Context, conn *pgxpool.Conn) {

}

func UpdateGroupById(c *gin.Context, conn *pgxpool.Conn) {

}

// get items in group
func GetItemsInGroup(c *gin.Context, conn *pgxpool.Conn) {
	requestBody, ok := c.Params.Get("id")

	fmt.Println("Request body", requestBody)

	if !ok {
		abortWithMessage(c, "error while parsing the request body")
	}
}

// Add the item listing to group
func AddItemToGroup(c *gin.Context, conn *pgxpool.Conn) {

}

// Update the item listing to group
func UpdateItemInGroup(c *gin.Context, conn *pgxpool.Conn) {

}
