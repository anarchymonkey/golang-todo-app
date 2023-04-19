package services

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
)

func GetContentsInItems(c *gin.Context, conn *pgxpool.Conn) {
	itemId, ok := c.Params.Get("id")

	if !ok {
		AbortWithMessage(c, "error reading item_id from query params")
		return
	}

	rows, err := conn.Query(context.Background(), `
	SELECT i.id item_id, c.id, c.content 
	FROM item_contents ic 
	JOIN items i on i.id = ic.item_id
	JOIN contents c on c.id = ic.content_id
	WHERE i.id=$1 AND ic.is_active=true AND i.is_active=true`, itemId)

	if err != nil {
		AbortWithMessage(c, fmt.Sprintf("error while running select query: ", err))
		return
	}

	var itemContents []ItemContentResponse
	for rows.Next() {
		var itemContent ItemContentResponse

		err := rows.Scan(&itemContent.ItemId, &itemContent.Id, &itemContent.Content)

		if err != nil {
			AbortWithMessage(c, fmt.Sprintf("error while scanning the rows %v", err))
			return
		}
		itemContents = append(itemContents, itemContent)
	}

	c.JSON(http.StatusOK, itemContents)
}

func AddContentInItem(c *gin.Context, conn *pgxpool.Conn) {
	const NoRowsAffected = 0
	// get the item id in params
	// this will be a transactional query, when we are adding content then we should add it in the items table too
	itemId, ok := c.Params.Get("id")

	if !ok {
		AbortWithMessage(c, "error while getting the itemId from params")
		return
	}
	var content Content

	if err := c.Bind(&content); err != nil {
		AbortWithMessage(c, fmt.Sprintf("error wile binding the request params with the content pointer %v", err))
		return
	}

	tran, err := conn.Begin(context.Background())

	if err != nil {
		// rollback the transcation if any error occours
		tran.Rollback(context.Background())
		AbortWithMessage(c, fmt.Sprintf("error while opening a transaction %v", err))
		return
	}

	// this will be the returning id from the query
	// use this id to insert into the item_contents table
	var contentId int
	err = conn.QueryRow(context.Background(), "INSERT INTO contents(content) VALUES($1) RETURNING id", content.Content).Scan(&contentId)

	if err != nil {
		tran.Rollback(context.Background())
		AbortWithMessage(c, fmt.Sprintf("error while running the query row %v", err))
		return
	}

	row, err := tran.Exec(context.Background(), "INSERT INTO item_contents(item_id, content_id) VALUES($1, $2)", itemId, contentId)

	if err != nil {
		tran.Rollback(context.Background())
		AbortWithMessage(c, fmt.Sprintf("error while inserting the item keys into the item_contents table %v", err))
		return
	}

	err = tran.Commit(context.Background())

	if err != nil {
		defer tran.Rollback(context.Background())
		AbortWithMessage(c, fmt.Sprintf("error while committing the transaction %v", err))
		return
	}

	// should I rollback if there are no rows inserted, what cases can be there?
	if row.RowsAffected() == NoRowsAffected {
		c.JSON(http.StatusOK, gin.H{
			"message": "No rows inserted",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "content added successfully!",
	})
}

func UpdateContentInItem(c *gin.Context, conn *pgxpool.Conn) {
	contentId, contentIdPresent := c.Params.Get("content_id")

	if !contentIdPresent {
		AbortWithMessage(c, "request params not present or nil")
		return
	}

	var content Content

	if err := c.Bind(&content); err != nil {
		AbortWithMessage(c, fmt.Sprintf("error while binding the content with reqeust body %v", err))
		return
	}

	row, err := conn.Exec(context.Background(), "UPDATE contents SET content=$2 WHERE id=$1", contentId, content.Content)

	if err != nil {
		AbortWithMessage(c, fmt.Sprintf("error while updating the content %v", err))
		return
	}

	if row.RowsAffected() == 0 {
		c.JSON(http.StatusOK, gin.H{
			"message": "no rows updated",
		})
		return
	}

	c.JSON(http.StatusOK, "successfully updated the content")
}

func DeleteContentInItem(c *gin.Context, conn *pgxpool.Conn) {
	itemId, itemIdPresent := c.Params.Get("id")
	contentId, contentIdPresent := c.Params.Get("content_id")

	if !itemIdPresent || !contentIdPresent {
		AbortWithMessage(c, "contentId not found in param, marked as required")
		return
	}

	tran, err := conn.Begin(context.Background())

	if err != nil {
		AbortWithMessage(c, fmt.Sprintf("error starting transaction %v", err))
		return
	}

	itemContentRows, err := tran.Exec(context.Background(), "DELETE from item_contents where item_id=$1 AND content_id=$2", itemId, contentId)

	if err != nil {
		tran.Rollback(context.Background())
		AbortWithMessage(c, fmt.Sprintf("error while deleting item rows %v", err))
		return
	}

	row, err := tran.Exec(context.Background(), "DELETE FROM contents WHERE id=$1", contentId)

	if err != nil {
		tran.Rollback(context.Background())
		AbortWithMessage(c, fmt.Sprintf("error while deleting item rows %v", err))
		return
	}

	err = tran.Commit(context.Background())

	if err != nil {
		AbortWithMessage(c, fmt.Sprintf("error while commiting updates %v", err))
		return
	}

	if itemContentRows.RowsAffected() == 0 || row.RowsAffected() == 0 {
		c.JSON(http.StatusOK, gin.H{
			"message": "rows are not affected",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "rows deleted successfully",
	})
}
