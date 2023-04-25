package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/anarchymonkey/golang-todo-server/db"
	"github.com/anarchymonkey/golang-todo-server/services"
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/gin-gonic/gin"
)

type Env struct {
	dbPool *pgxpool.Pool
}

func getHandlerWithBindedEnvironment(fn func(*gin.Context, *pgxpool.Conn), env *Env) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		conn, err := db.AcquireConnectionFromPool(env.dbPool)

		if err != nil {
			log.Fatal(err)
		}

		defer conn.Release()

		fn(ctx, conn)
	}
}

func main() {
	const PORT = 8080

	fmt.Println("Welcome to the todo web application")

	var dbConfig db.DbConfig = db.DbConfig{
		Username: os.Getenv("PG_USERNAME_V1"),
		Password: os.Getenv("PG_PASS_V1"),
		DbName:   "todo_app",
		PORT:     5432,
	}

	// use the connection pool to pass it through a handler
	dbConnPool, err := dbConfig.GetDbConnectionPool()

	if err != nil {
		log.Fatal(err)
	}
	defer dbConnPool.Close()

	var env *Env = &Env{
		dbPool: dbConnPool,
	}

	router := gin.Default()

	router.Use(corsMiddleware())

	// can refactor with router.group later on

	router.GET("/groups", getHandlerWithBindedEnvironment(services.GetGroups, env))
	router.GET("/group/:id/items", getHandlerWithBindedEnvironment(services.GetItemsInGroup, env))
	router.GET("/item/:id/contents", getHandlerWithBindedEnvironment(services.GetContentsInItems, env))

	// POST
	router.POST("/group/add", getHandlerWithBindedEnvironment(services.AddGroup, env))
	router.POST("/group/:id/item/add", getHandlerWithBindedEnvironment(services.AddItemToGroup, env))
	router.POST("/item/:id/content/add", getHandlerWithBindedEnvironment(services.AddContentInItem, env))

	// PUT
	router.PUT("/group/:id/update", getHandlerWithBindedEnvironment(services.UpdateGroupById, env))
	router.PUT("/item/:id/update", getHandlerWithBindedEnvironment(services.UpdateItemInGroup, env))
	router.PUT("/content/:content_id/update", getHandlerWithBindedEnvironment(services.UpdateContentInItem, env))

	// DELETE

	// This needs the group id and item_id respectively cause we want to delete the row from the mapping table also

	router.DELETE("/group/:id/delete", getHandlerWithBindedEnvironment(services.DeleteGroupById, env))
	router.DELETE("/group/:id/item/:item_id/delete", getHandlerWithBindedEnvironment(services.DeleteItemInGroup, env))
	router.DELETE("/item/:id/content/:content_id/delete", getHandlerWithBindedEnvironment(services.DeleteContentInItem, env))

	// run the server
	router.Run(fmt.Sprintf(":%d", PORT))

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
