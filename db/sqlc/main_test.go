package db

import (
	"context"
	"log"
	"os"
	"testing"

	
	"github.com/jackc/pgx/v5/pgxpool"

	_ "github.com/lib/pq"
)

const (
	dbSource = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
)

var testQueries *Queries
var testDB *pgxpool.Pool
var testStore *Store

func TestMain(m *testing.M) {
	var err error
 
	config, err := pgxpool.ParseConfig(dbSource)
	if err != nil {
		log.Fatal("cannot parse config:", err)
	}

	testDB, err = pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatal("cannot connect to db:", err)

	
	}
	
	testQueries = New(testDB)
	testStore = NewStore(testDB)

	os.Exit(m.Run())
}