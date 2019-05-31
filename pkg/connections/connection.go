package connections

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"

	_ "github.com/lib/pq"
)

type DatabaseConnection interface {
	Connect() sql.DB
}

type postgreSqlConnection struct {
	HostName     string
	Port         int64
	UserName     string
	password     string
	DatabaseName string
	Connection   *sql.DB
}

type PostgreSqlConnector struct{}

var psqlConnection *postgreSqlConnection

func (conn PostgreSqlConnector) Connect() *sql.DB {
	if psqlConnection != nil {
		return psqlConnection.Connection
	}

	hostName := os.Getenv("DB_HOST")
	port, err := strconv.ParseInt(os.Getenv("DB_PORT"), 10, 64)
	if err != nil {
		panic(err)
	}
	user := os.Getenv("DB_USER")
	pwd := os.Getenv("DB_PWD")
	dbName := os.Getenv("DB_NAME")
	connectionStr := fmt.Sprintf("host=%s port=%d user=%s password=%s "+
		"dbname=%s sslmode=disable",
		hostName,
		port,
		user,
		pwd,
		dbName)

	db, err := sql.Open("postgres", connectionStr)
	if err != nil {
		panic(err)
	}

	psqlConnection := postgreSqlConnection{hostName, port, user, pwd, dbName, db}

	return psqlConnection.Connection
}
