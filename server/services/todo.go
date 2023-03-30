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
	Id        int        `json:"id"`
	Content   string     `json:"content"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	RemindAt  *time.Time `json:"remind_at"`
	IsActive  bool       `json:"is_active"`
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
	defer rows.Close()

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

func DeleteGroupById(c *gin.Context, conn *pgxpool.Conn) {

	// does not need group id, item id would suffice

	idToDelete, ok := c.Params.Get("id")

	if !ok {
		abortWithMessage(c, "item id is required but not present!")
		return
	}

	row, err := conn.Exec(context.Background(), "DELETE from groups where id=$1", idToDelete)

	if err != nil {
		abortWithMessage(c, fmt.Sprintf("error while running the delete query %v", err))
	}

	if row.RowsAffected() == 0 {
		c.JSON(http.StatusOK, gin.H{
			"message": "Cannot be deleted at the moment!",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "The group has been deleted",
	})
}

// get items in group
func GetItemsInGroup(c *gin.Context, conn *pgxpool.Conn) {
	groupId, ok := c.Params.Get("id")

	if !ok {
		abortWithMessage(c, fmt.Sprintf("error while getting the request params"))
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
		abortWithMessage(c, fmt.Sprintf("error while querying the table: %v", err))
		return
	}
	defer rows.Close()

	var items []Item

	for rows.Next() {
		var item Item
		var remindAt *time.Time
		err := rows.Scan(&item.Id, &item.Content, &item.IsActive, &item.CreatedAt, &item.UpdatedAt, &remindAt)

		if err != nil {
			abortWithMessage(c, fmt.Sprintf("error while scanning rows, %v", err))
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
	// does not need group id, item id would suffice

	idToUpdate, ok := c.Params.Get("id")

	if !ok {
		abortWithMessage(c, "item id is required but not present!")
		return
	}

	var item Item

	if err := c.BindJSON(&item); err != nil {
		abortWithMessage(c, fmt.Sprintf("error while reading the request body %v", err))
		return
	}

	row, err := conn.Exec(context.Background(), "UPDATE items SET (content, is_active, remind_at) = ($2, $3, $4) where id=$1", idToUpdate, item.Content, item.IsActive, item.RemindAt)

	if err != nil {
		abortWithMessage(c, fmt.Sprintf("error while running the update query %v", err))
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
		abortWithMessage(c, "error getting params, or params not present")
		return
	}

	// This should be a transaction as grouped_items also need to be set is_active = false

	tran, err := conn.Begin(context.Background())

	if err != nil {
		abortWithMessage(c, fmt.Sprintf("error while starting a transaction, %v", err))
		return
	}

	groupedItemRowDetails, err := tran.Exec(context.Background(), "DELETE from grouped_items where group_id=$1 AND item_id=$2", groupId, itemId)

	if err != nil {
		tran.Rollback(context.Background())
		abortWithMessage(c, fmt.Sprintf("error while running a delete query on grouped_items: id %s, with err %v", itemId, err))
		return
	}

	itemRowDetails, err := conn.Exec(context.Background(), "DELETE FROM items where id=$1", itemId)

	if err != nil {
		tran.Rollback(context.Background())
		abortWithMessage(c, fmt.Sprintf("error while running a delete query on item id %s, with err %v", itemId, err))
		return
	}

	fmt.Println("rows affected", itemRowDetails.RowsAffected(), groupedItemRowDetails.RowsAffected())

	err = tran.Commit(context.Background())

	if err != nil {
		abortWithMessage(c, fmt.Sprintf("error while commiting the transaction %v", err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "data deleted successfully",
	})

}

func GetContentsInItems(c *gin.Context, conn *pgxpool.Conn) {

}

func AddContentInItem(c *gin.Context, conn *pgxpool.Conn) {

}

func UpdateContentInItem(c *gin.Context, conn *pgxpool.Conn) {

}

func DeleteContentInItem(c *gin.Context, conn *pgxpool.Conn) {

}
