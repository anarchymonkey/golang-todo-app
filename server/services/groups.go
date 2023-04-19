package services

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
)

func GetGroups(c *gin.Context, conn *pgxpool.Conn) {
	var ctx context.Context = context.Background()

	rows, err := conn.Query(ctx, "SELECT * from groups WHERE is_active=true")

	if err != nil {
		AbortWithMessage(c, "error fetching rows from query")
		return
	}
	defer rows.Close()

	var groups []Group = []Group{}

	for rows.Next() {
		var group Group = Group{}
		err := rows.Scan(&group.Id, &group.Title, &group.Description, &group.CreatedAt, &group.UpdatedAt, &group.IsActive, &group.IsPublic)

		if err != nil {
			AbortWithMessage(c, fmt.Sprintf("error while scanning rows %v", err))
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
		AbortWithMessage(c, fmt.Sprintf("error while fetching the request body %v", err))
		return
	}

	row, err := conn.Exec(context.Background(), "INSERT INTO groups(title, description) VALUES($1, $2);", group.Title, group.Description)

	if err != nil {
		AbortWithMessage(c, fmt.Sprintf("error while executing the insert statement %v", err))
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
		AbortWithMessage(c, fmt.Sprintf("error while binding json %v", err))
		return
	}

	updatedRow, err := conn.Exec(context.Background(), "UPDATE groups set (title, description, is_active, is_public) = ($1, $2, $3, $4)", &group.Title, &group.Description, &group.IsActive, &group.IsPublic)

	if err != nil {
		AbortWithMessage(c, fmt.Sprintf("error while updating row json %v", err))
		return
	}

	if updatedRow.RowsAffected() == NoRowsAffected {
		c.JSON(http.StatusOK, "Could not insert data!")
		return
	}

	c.JSON(http.StatusOK, fmt.Sprintf("Data has been updated! affected rows = %d", updatedRow.RowsAffected()))
}

func DeleteGroupById(c *gin.Context, conn *pgxpool.Conn) {

	idToDelete, ok := c.Params.Get("id")

	if !ok {
		AbortWithMessage(c, "id not present")
		return
	}

	tran, err := conn.Begin(context.Background())

	if err != nil {
		AbortWithMessage(c, fmt.Sprintf("Error while begining a transaction %v", err))
		return
	}

	// delete group in grouped items and return the item ids

	ids, err := tran.Query(context.Background(), "DELETE from grouped_items where group_id=$1 RETURNING item_id", idToDelete)

	if err != nil {
		tran.Rollback(context.Background())
		AbortWithMessage(c, fmt.Sprintf("error while deleting grouped_items: %v", err))
		return
	}

	var itemIds []string

	for ids.Next() {
		var item int

		err := ids.Scan(&item)

		if err != nil {
			tran.Rollback(context.Background())
			AbortWithMessage(c, fmt.Sprintf("error while scanning rows: %v", err))
			defer ids.Close()
			return
		}

		itemIds = append(itemIds, strconv.Itoa(item))
	}

	// delete group id

	_, err = tran.Exec(context.Background(), "DELETE from groups where id=$1", idToDelete)

	if err != nil {
		tran.Rollback(context.Background())
		AbortWithMessage(c, fmt.Sprintf("error while deleting group: %v", err))
		return

	}

	// delete items and contents from item_contents where item_id in item_ids

	fmt.Println("The item ids", itemIds)

	ids, err = tran.Query(context.Background(), fmt.Sprintf("DELETE from item_contents where item_id IN (%s) RETURNING content_id", strings.Join(itemIds, ",")))

	if err != nil {
		tran.Rollback(context.Background())
		AbortWithMessage(c, fmt.Sprintf("error while deleting item contents: %v", err))
		return
	}

	var contentIds []string
	for ids.Next() {
		var contentId int

		err := ids.Scan(&contentId)

		if err != nil {
			defer ids.Close()
			tran.Rollback(context.Background())
			return
		}
		contentIds = append(contentIds, strconv.Itoa(contentId))
	}

	// delete itemIds from items

	_, err = tran.Exec(context.Background(), fmt.Sprintf("DELETE from items where id IN (%s)", strings.Join(itemIds, ",")))

	if err != nil {
		tran.Rollback(context.Background())
		AbortWithMessage(c, fmt.Sprintf("error while deleting items: %v", err))
		return
	}

	// delete contents

	_, err = tran.Exec(context.Background(), fmt.Sprintf("DELETE from contents where id IN (%s)", strings.Join(contentIds, ",")))

	if err != nil {
		tran.Rollback(context.Background())
		AbortWithMessage(c, fmt.Sprintf("error while deleting items: %v", err))
		return
	}

	err = tran.Commit(context.Background())

	if err != nil {
		tran.Rollback(context.Background())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "group and it's items deleted successfully",
	})
}
