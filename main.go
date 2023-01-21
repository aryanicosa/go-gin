package main

import (
	"database/sql"
	"github.com/aryanicosa/go_gin_simple_bank/api"
	db "github.com/aryanicosa/go_gin_simple_bank/db/sqlc"
	"github.com/aryanicosa/go_gin_simple_bank/util"
	_ "github.com/lib/pq" // _ means we use it without call any function from it directly
	"log"
)

func main() {
	config, err := util.LoadConfig(".")
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("can not connect to db:", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("can not start server: ", err)
	}
}
