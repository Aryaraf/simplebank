package db

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5"

	_"github.com/lib/pq"
)

const (
	dbSource = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
)

var testQueries *Queries
var testDB *pgx.Conn
var testStore *Store

func TestMain(m *testing.M) {
	var err error
 
	testDB, err = pgx.Connect(context.Background(), dbSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	testQueries = New(testDB)
	testStore = &Store{
		db: testDB,
		realDB: testDB,
		Queries: New(testDB),
	}

	os.Exit(m.Run())
}