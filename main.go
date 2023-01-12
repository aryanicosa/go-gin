package main

import (
	"database/sql"
	"github.com/aryanicosa/go_gin_simple_bank/api"
	db "github.com/aryanicosa/go_gin_simple_bank/db/sqlc"
	_ "github.com/lib/pq" // _ means we use it without call any function from it directly
	"log"
)

const (
	dbDriver      = "postgres"
	dbSource      = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
	serverAddress = "0.0.0.0:8080"
)

func main() {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("can not connect to db:", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(serverAddress)
	if err != nil {
		log.Fatal("can not start server: ", err)
	}
}
