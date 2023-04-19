package services

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
)

// get items in group
func GetItemsInGroup(c *gin.Context, conn *pgxpool.Conn) {
	groupId, ok := c.Params.Get("id")

	if !ok {
		AbortWithMessage(c, fmt.Sprintf("error while getting the request params"))
		return
	}

	/*
		Analysis


					 Sort  (cost=17.48..17.48 rows=1 width=61) (actual time=0.063..0.064 rows=7 loops=1)
		   Sort Key: i.created_at
		   Sort Method: quicksort  Memory: 25kB
		   ->  Nested Loop  (cost=0.30..17.47 rows=1 width=61) (actual time=0.037..0.055 rows=7 loops=1)
		         ->  Nested Loop  (cost=0.15..9.27 rows=1 width=4) (actual time=0.028..0.038 rows=7 loops=1)
		               ->  Seq Scan on grouped_items gi  (cost=0.00..1.09 rows=1 width=8) (actual time=0.018..0.019 rows=7 loops=1)
		                     Filter: (is_active AND (group_id = 1))
		               ->  Index Scan using groups_pkey on groups g  (cost=0.15..8.17 rows=1 width=4) (actual time=0.002..0.002 rows=1 loops=7)
		                     Index Cond: (id = 1)
		                     Filter: is_active
		         ->  Index Scan using items_pkey on items i  (cost=0.15..8.17 rows=1 width=61) (actual time=0.002..0.002 rows=1 loops=7)
		               Index Cond: (id = gi.item_id)
		               Filter: is_active
		 Planning Time: 0.247 ms
		 Execution Time: 0.102 ms
		(15 rows)

	*/

	rows, err := conn.Query(context.Background(), `
	SELECT i.*
	FROM grouped_items gi 
	JOIN groups g on g.id = gi.group_id 
	JOIN items i on i.id = gi.item_id 
	WHERE gi.group_id=$1 AND gi.is_active = true AND g.is_active=true AND i.is_active=true 
	ORDER BY i.created_at ASC NULLS LAST;`, groupId)

	if err != nil {
		AbortWithMessage(c, fmt.Sprintf("error while querying the table: %v", err))
		return
	}
	defer rows.Close()

	var items []Item

	for rows.Next() {
		var item Item
		var remindAt *time.Time
		err := rows.Scan(&item.Id, &item.Content, &item.IsActive, &item.CreatedAt, &item.UpdatedAt, &remindAt)

		if err != nil {
			AbortWithMessage(c, fmt.Sprintf("error while scanning rows, %v", err))
			return
		}

		item.RemindAt = remindAt
		items = append(items, item)
	}

	c.JSON(http.StatusOK, items)
}

// Add the item listing to group
func AddItemToGroup(c *gin.Context, conn *pgxpool.Conn) {
	group_id, ok := c.Params.Get("id")

	if !ok {
		AbortWithMessage(c, "error while parsing the request body")
		return
	}

	var item Item

	if err := c.BindJSON(&item); err != nil {
		AbortWithMessage(c, "error while binding")
	}

	tran, err := conn.Begin(context.Background())

	if err != nil {
		AbortWithMessage(c, fmt.Sprintf("error while begining the transaction %v", err))
	}
	defer tran.Rollback(context.Background())

	var id int
	err = tran.QueryRow(context.Background(), "INSERT into items(content) values($1) RETURNING id", &item.Content).Scan(&id)

	if err != nil {
		// defer tran.Rollback(context.Background())
		AbortWithMessage(c, fmt.Sprintf("error while getting inserted item Id: %v", err))
		return
	}

	err = tran.QueryRow(context.Background(), "INSERT into grouped_items(group_id, item_id) VALUES($1, $2) RETURNING id", group_id, id).Scan(&id)

	fmt.Println("The id grouped_items is", id)

	if err != nil {
		// defer tran.Rollback(context.Background())
		AbortWithMessage(c, fmt.Sprintf("error while getting inserted item in grouped_items Id: %v", err))
		return
	}

	defer tran.Commit(context.Background())

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("items have been added, the id of the added item is %d", id),
	})
}

// Update the item listing to group
func UpdateItemInGroup(c *gin.Context, conn *pgxpool.Conn) {
	// does not need group id, item id would suffice

	idToUpdate, ok := c.Params.Get("id")

	if !ok {
		AbortWithMessage(c, "item id is required but not present!")
		return
	}

	var item Item

	if err := c.BindJSON(&item); err != nil {
		AbortWithMessage(c, fmt.Sprintf("error while reading the request body %v", err))
		return
	}

	row, err := conn.Exec(context.Background(), "UPDATE items SET (content, is_active, remind_at) = ($2, $3, $4) where id=$1", idToUpdate, item.Content, item.IsActive, item.RemindAt)

	if err != nil {
		AbortWithMessage(c, fmt.Sprintf("error while running the update query %v", err))
	}

	if row.RowsAffected() == 0 {
		c.JSON(http.StatusOK, gin.H{
			"message": "Cannot be updated at the moment!",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "updated successfully!",
	})
}

func DeleteItemInGroup(c *gin.Context, conn *pgxpool.Conn) {
	groupId, groupIdPresent := c.Params.Get("id")
	itemId, itemIdPresent := c.Params.Get("item_id")

	if !groupIdPresent || !itemIdPresent {
		AbortWithMessage(c, "error getting params, or params not present")
		return
	}

	// This should be a transaction as grouped_items also need to be set is_active = false

	tran, err := conn.Begin(context.Background())

	if err != nil {
		AbortWithMessage(c, fmt.Sprintf("error while starting a transaction, %v", err))
		return
	}

	// delete item_contents

	ids, err := tran.Query(context.Background(), "DELETE from item_contents WHERE item_id=$1 RETURNING content_id as id;", itemId)

	if err != nil {
		tran.Rollback(context.Background())
		AbortWithMessage(c, fmt.Sprintf("error while deleting item from item_contents: err = %v", err))
		return
	}

	var contentIds []string

	for ids.Next() {
		var content int

		err := ids.Scan(&content)

		if err != nil {
			tran.Rollback(context.Background())
			AbortWithMessage(c, fmt.Sprintf("error while getting content ids: err = %v", err))
			return
		}

		// if there is some content id then only push the id
		if content != 0 {
			contentIds = append(contentIds, strconv.Itoa(content))
		}

	}

	// if content ids are available then only delete otherwise move on to the next step
	if len(contentIds) != 0 {
		fmt.Println("content ids", contentIds, fmt.Sprintf("DELETE from contents where id IN (%s)", strings.Join(contentIds, ",")))
		_, err := tran.Exec(context.Background(), fmt.Sprintf("DELETE from contents where id IN (%s)", strings.Join(contentIds, ",")))

		if err != nil {
			tran.Rollback(context.Background())
			AbortWithMessage(c, fmt.Sprintf("error while deleting content with ids = %s: err = %v", contentIds, err))
			return
		}
	}

	// delete grouped_items

	_, err = tran.Exec(context.Background(), "DELETE from grouped_items WHERE group_id=$1 AND item_id=$2", groupId, itemId)

	if err != nil {
		tran.Rollback(context.Background())
		AbortWithMessage(c, fmt.Sprintf("error while deleting grouped_items with group_id=%s and item_id=%s: err=%v", groupId, itemId, err))
		return
	}

	// delete items

	_, err = tran.Exec(context.Background(), "DELETE from items WHERE id=$1", itemId)

	if err != nil {
		tran.Rollback(context.Background())
		AbortWithMessage(c, fmt.Sprintf("error while deleting items with item_id=%s: err=%v", itemId, err))
		return
	}

	err = tran.Commit(context.Background())

	if err != nil {
		tran.Rollback(context.Background())
		AbortWithMessage(c, fmt.Sprintf("error while commiting the transaction: err=%v", err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "successfully deleted items",
	})

}
