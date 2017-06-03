package app

import (
	"log"

	"gopkg.in/mgo.v2"
)

// Database is our database struct, we didn't use a global variable
// to improve swapping database
type Database struct {
	// the database name
	name string
	// The database session
	session *mgo.Session
	// the database context
	ctx *mgo.Database
}

// Setup connects the database
func (db *Database) Setup() (err error) {
	// create a new session for the database
	db.session, err = mgo.Dial("localhost:27017")
	if err != nil {
		log.Println("Error connecting to database")
		return err
	}
	db.ctx = db.session.DB(db.name)

	log.Printf("Connected to database name=%s", db.name)
	return err
}

// Closes the database connection
func (db *Database) Close() {
	db.Close()
}
func (db *Database) CopySession() *mgo.Session {
	// Copy creates a new session that listens to different socket
	return db.session.Copy()
}

// Returns the databse context
func (db *Database) Collection(collection string, session *mgo.Session) *mgo.Collection {
	return db.ctx.C(collection).With(session)
}

// NewDatabase creates a new database
func NewDatabase(name string) *Database {
	db := Database{name: name}
	err := db.Setup()
	if err != nil {
		panic(err)
	}
	return &db
}
