package app

import (
	"github.com/hashicorp/go-memdb"
)

var db *memdb.MemDB

func setupDatabase() (*memdb.MemDB, error) {
	// Create the database schema
	schema := &memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{
			"user": {
				Name: "user",
				Indexes: map[string]*memdb.IndexSchema{
					"id": {
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "Email"},
					},
				},
			},
		},
	}

	// Create a new database
	db, err := memdb.NewMemDB(schema)
	if err != nil {
		return db, err
	}
	return db, nil
}

// Database setups
func Database() *memdb.MemDB {
	if db == nil {
		var err error
		db, err = setupDatabase()

		if err != nil {
			panic(err)
		}

		return db
	}
	return db
}
