package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
)

type DbConfig struct {
	// The username of the psql service
	Username string

	// The password required, black if no password is given
	Password string

	// The database to connect with
	DbName string

	// The port to establish a connection with, default: 5432
	PORT int
}

// global error definitions
const (
	CONFIG_PARSE_ERROR = "error while parsing the database config \n%v"
	DB_CONNECT_ERROR   = "error while connecting with the database \n%v"
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
