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
	// get connection from connection pool
	conn, err := db.AcquireConnectionFromPool(env.dbPool)

	if err != nil {
		log.Fatal(err)
	}

	// release the connection when it's done so that it can be aquired by someone else
	defer conn.Release()

	return func(ctx *gin.Context) {
		fn(ctx, conn)
	}
}

func main() {

	fmt.Println("Welcome to the todo web application")

	var dbConfig db.DbConfig = db.DbConfig{
		Username: os.Getenv("PG_USERNAME_V1"),
		Password: os.Getenv("PG_PASS_V1"),
		DbName:   "postgres",
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

	router.GET("/groups", getHandlerWithBindedEnvironment(services.GetGroups, env))
	router.POST("/group/add", getHandlerWithBindedEnvironment(services.AddGroup, env))
	router.PUT("/group/:id/update", getHandlerWithBindedEnvironment(services.UpdateGroupById, env))
	router.GET("/group/:id/items", getHandlerWithBindedEnvironment(services.GetItemsInGroup, env))
	router.POST("item/add", getHandlerWithBindedEnvironment(services.AddItemToGroup, env))
	router.PUT("/item/:id/update", getHandlerWithBindedEnvironment(services.UpdateItemInGroup, env))

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
