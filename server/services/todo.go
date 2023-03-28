package services

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
)

// List of groups present todo_app.groups
type Group struct {
	Id          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	IsActive    bool      `json:"is_active"`
	IsPublic    bool      `json:"is_public"`
}

// resembles todo_app.items
type Item struct {
	Id        int       `json:"id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	RemindAt  string    `json:"remind_at"`
	IsActive  bool      `json:"is_active"`
}

type GroupedItem struct {
	Id        int       `json:"id"`
	GroupId   int       `json:"group_id"`
	ItemId    int       `json:"item_id"`
	CreatedAt time.Time `json:"created_at"`
	IsActive  bool      `json:"is_active"`
}

func abortWithMessage[T any](c *gin.Context, message T) {
	c.AbortWithStatusJSON(http.StatusBadRequest, message)
}

func GetGroups(c *gin.Context, conn *pgxpool.Conn) {
	var ctx context.Context = context.Background()

	rows, err := conn.Query(ctx, "SELECT * from groups WHERE is_active=true")

	if err != nil {
		abortWithMessage(c, "error fetching rows from query")
		return
	}

	var groups []Group = []Group{}

	for rows.Next() {
		var group Group = Group{}
		err := rows.Scan(&group.Id, &group.Title, &group.Description, &group.CreatedAt, &group.UpdatedAt, &group.IsActive, &group.IsPublic)

		if err != nil {
			abortWithMessage(c, fmt.Sprintf("error while scanning rows %v", err))
			return
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
		return
	}

	row, err := conn.Exec(context.Background(), "INSERT INTO groups(title, description) VALUES($1, $2);", group.Title, group.Description)

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
	// group_id, ok := c.Params.Get("id")

	// fmt.Println("Request body", group_id)

	// if !ok {
	// 	abortWithMessage(c, "error while parsing the request body")
	// 	return
	// }

	// rows, err := conn.Query(context.Background(), "SELECT g.*, gi.* from grouped_items gi join groups g on g.id = gi.group_id where gi.group_id=$1", group_id)

	// if err != nil {
	// 	abortWithMessage(c, fmt.Sprintf("error while running select query %v", err))
	// 	return
	// }

	// for rows.Next() {
	// 	err := rows.Values()
	// }

}

// Add the item listing to group
func AddItemToGroup(c *gin.Context, conn *pgxpool.Conn) {
	group_id, ok := c.Params.Get("id")

	if !ok {
		abortWithMessage(c, "error while parsing the request body")
		return
	}

	var item Item

	if err := c.BindJSON(&item); err != nil {
		abortWithMessage(c, "error while binding")
	}

	tran, err := conn.Begin(context.Background())

	if err != nil {
		abortWithMessage(c, fmt.Sprintf("error while begining the transaction %v", err))
	}
	defer tran.Rollback(context.Background())

	var id int
	err = tran.QueryRow(context.Background(), "INSERT into items(content) values($1) RETURNING id", &item.Content).Scan(&id)

	if err != nil {
		// defer tran.Rollback(context.Background())
		abortWithMessage(c, fmt.Sprintf("error while getting inserted item Id: %v", err))
		return
	}

	err = tran.QueryRow(context.Background(), "INSERT into grouped_items(group_id, item_id) VALUES($1, $2) RETURNING id", group_id, id).Scan(&id)

	fmt.Println("The id grouped_items is", id)

	if err != nil {
		// defer tran.Rollback(context.Background())
		abortWithMessage(c, fmt.Sprintf("error while getting inserted item in grouped_items Id: %v", err))
		return
	}

	defer tran.Commit(context.Background())

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("items have been added, the id of the added item is %d", id),
	})
}

// Update the item listing to group
func UpdateItemInGroup(c *gin.Context, conn *pgxpool.Conn) {

}
