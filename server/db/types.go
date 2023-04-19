package db

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
