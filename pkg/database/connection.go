package connection;

import (
    "database/sql"
    "fmt"
    _ "github.com/lib/pq"
)

type DatabaseConnection interface {
    Connect() sql.DB
}

type PostgreSqlConnector struct {
    HostName string
    Port int
    UserName string
    Password string // yikes
    DatabaseName string
}

func (psqlConnection PostgreSqlConnector) Connect() *sql.DB {
    connectionStr := fmt.Sprintf("host=%s port=%d user=%s password=%s " +
    "dbname=%s sslmode=disable",
    psqlConnection.HostName,
    psqlConnection.Port,
    psqlConnection.UserName,
    psqlConnection.Password,
    psqlConnection.DatabaseName)

    db, err := sql.Open("postgres", connectionStr)
    if err != nil {
        panic(err)
    }

    return db
}

