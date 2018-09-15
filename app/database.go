package app

import (
	"log"

	"gopkg.in/mgo.v2"
)

// Database helps initialize a new storage for the application and stores the database related fields
type Database struct {
	// The database name
	name string
	// The data source name
	dsn string
	// The database session
	session *mgo.Session
	// the database context
	ctx *mgo.Database
}

// Setup connects the database to the data source and returns an error if a connection fails to establish.
func (db *Database) Setup() (err error) {
	// create a new session for the database
	db.session, err = mgo.Dial(db.dsn)
	if err != nil {
		log.Println("Error connecting to database")
		return err
	}
	db.ctx = db.session.DB(db.name)

	log.Printf("Connected to database name=%s", db.name)
	return err
}

// Close terminates the database connection
func (db *Database) Close() {
	db.Close()
}

// NewSession creates a new session for the database
func (db *Database) NewSession() *mgo.Session {
	return db.session.Copy()
}

// Collection returns a new collection that is tied to a particular session
func (db *Database) Collection(collection string, session *mgo.Session) *mgo.Collection {
	return db.ctx.C(collection).With(session)
}

// NewDatabase creates a new database
func NewDatabase(name, dsn string) *Database {
	db := Database{name: name, dsn: dsn}
	err := db.Setup()
	if err != nil {
		panic(err)
	}
	return &db
}
