package db

import (
	"gopkg.in/mgo.v2"
	"log"
	"os"
		)

var Session *mgo.Session

///CHANGE THESE
var mongohosts = "daas_database_1"//"mongodb://localhost"//os.Getenv("MONGO_URL")


func init() {
	log.Printf("Datastore Initialized");
	var err error
	Session, err = mgo.Dial(mongohosts)
	if err != nil {
		log.Fatalf("Error connecting to the database: %s\n", err.Error())
	}
	Session.SetMode(mgo.Monotonic, false)
    Session.DB(os.Getenv("DB_NAME"))
}