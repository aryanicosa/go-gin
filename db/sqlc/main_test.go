package db

import (
	"database/sql"
	"github.com/aryanicosa/go_gin_simple_bank/util"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq" // _ means we use it without call any function from it directly
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	var err error

	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("can not load config:", err)
	}
	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("can not connect to db:", err)
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}
