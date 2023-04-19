package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
)

// gets the db url wrt the dbConfig
func getDBUrlFromConfig(dbConfig *DbConfig) string {
	return fmt.Sprintf("postgres://%s:%s@localhost:%d/%s", dbConfig.Username, dbConfig.Password, dbConfig.PORT, dbConfig.DbName)
}

// params: a dbConfig which is required to establish a connection
// returns: a connection pool and error if any otherwise nil
func (dbConfig *DbConfig) GetDbConnectionPool() (*pgxpool.Pool, error) {

	poolConfig, err := pgxpool.ParseConfig(getDBUrlFromConfig(dbConfig))

	if err != nil {
		return nil, fmt.Errorf(CONFIG_PARSE_ERROR, err)
	}

	// connect to db with the parsed config
	conn, err := pgxpool.ConnectConfig(context.Background(), poolConfig)

	if err != nil {
		return nil, fmt.Errorf(DB_CONNECT_ERROR, err)
	}

	return conn, nil
}

// this is used to aquire a connection from the pool
func AcquireConnectionFromPool(pool *pgxpool.Pool) (*pgxpool.Conn, error) {
	conn, err := pool.Acquire(context.Background())

	if err != nil {
		return nil, fmt.Errorf("error while creating a connection from the pool")
	}
	return conn, nil
}
