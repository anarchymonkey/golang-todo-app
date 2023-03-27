package services

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
)

// List of groups present todo_app.groups
type Group struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	IsActive    bool   `json:"is_active"`
	IsPublic    bool   `json:"is_public"`
}

// resembles todo_app.items
type Item struct {
	Id        int    `json:"id"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	RemindAt  string `json:"remind_at"`
	IsActive  bool   `json:"is_active"`
}

type GroupedItem struct {
	Id        int    `json:"id"`
	GroupId   int    `json:"group_id"`
	ItemId    int    `json:"item_id"`
	CreatedAt string `json:"created_at"`
	IsActive  bool   `json:"is_active"`
}

func abortWithMessage[T any](c *gin.Context, message T) {
	c.AbortWithStatusJSON(http.StatusBadRequest, message)
}

// this is done
func GetGroups(c *gin.Context, conn *pgxpool.Conn) {
	var ctx context.Context = context.Background()

	rows, err := conn.Query(ctx, "SELECT * from groups WHERE is_active=true")

	defer conn.Release()

	fmt.Println("rows", rows)

	if err != nil {
		abortWithMessage(c, "error fetching rows from query")
	}

	var groups []Group = []Group{}

	for rows.Next() {
		var group Group = Group{}

		err := rows.Scan(&group.Id, &group.Title, &group.IsActive, &group.IsPublic, &group.CreatedAt, &group.UpdatedAt, &group.Description)

		fmt.Println("group", group)

		if err != nil {
			abortWithMessage(c, fmt.Sprintf("error while scanning rows %v", err))
		}

		groups = append(groups, group)
	}
	c.JSON(http.StatusOK, groups)
}

func AddGroup(c *gin.Context, conn *pgxpool.Conn) {
	const NoRowsAffected = 0

	var group Group

	if err := c.BindJSON(&group); err != nil {
		abortWithMessage(c, fmt.Sprintf("error while fetching the request body %v", err))
	}

	row, err := conn.Exec(context.Background(), "INSERT INTO groups(title, description) VALUES($1, $2);", group.Title, group.Description)
	defer conn.Release()

	if err != nil {
		abortWithMessage(c, fmt.Sprintf("error while executing the insert statement %v", err))
		return
	}

	if row.RowsAffected() == NoRowsAffected {
		c.JSON(http.StatusOK, "Could not insert data!")
		return
	}

	c.JSON(http.StatusOK, fmt.Sprintf("Data has been inserted, rows affected = %d", row.RowsAffected()))

}

func UpdateGroupById(c *gin.Context, conn *pgxpool.Conn) {

	const NoRowsAffected = 0

	var group Group

	if err := c.BindJSON(&group); err != nil {
		abortWithMessage(c, fmt.Sprintf("error while binding json %v", err))
		return
	}

	updatedRow, err := conn.Exec(context.Background(), "UPDATE groups set (title, description, is_active, is_public) = ($1, $2, $3, $4)", &group.Title, &group.Description, &group.IsActive, &group.IsPublic)
	defer conn.Release()
	if err != nil {
		abortWithMessage(c, fmt.Sprintf("error while updating row json %v", err))
		return
	}

	if updatedRow.RowsAffected() == NoRowsAffected {
		c.JSON(http.StatusOK, "Could not insert data!")
		return
	}

	c.JSON(http.StatusOK, fmt.Sprintf("Data has been updated! affected rows = %d", updatedRow.RowsAffected()))
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
