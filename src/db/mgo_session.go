package db

import (
	"fmt"
	"log"
	//"time"
	mgo "gopkg.in/mgo.v2"
	//"gopkg.in/mgo.v2"
)

var (
	Database     *mgo.Database
	MongoDBHosts = ""
	AuthDatabase = ""
	AuthUserName = ""
	AuthPassword = ""
	MongoDBUri   = ""
)

func InitSession() {
	// We need this object to establish a session to our MongoDB.
	fmt.Println("setting mongo session")
	session, err := mgo.Dial(MongoDBUri)
	if err != nil {
		log.Fatalf("CreateSession: %s\n", err)
	}
	/*mongoDBDialInfo := mgo.DialInfo{
		Addrs:    []string{MongoDBHosts},
		Timeout:  60 * time.Second,
		Database: AuthDatabase,
		Username: AuthUserName,
		Password: AuthPassword,
	}

	// Create a session which maintains a pool of socket connections
	// to our MongoDB.
	mongoSession, err := mgo.DialWithInfo(mongoDBDialInfo)*/

	// Reads may not be entirely up-to-date, but they will always see the
	// history of changes moving forward, the data read will be consistent
	// across sequential queries in the same session, and modifications made
	// within the session will be observed in following queries (read-your-writes).
	// http://godoc.org/labix.org/v2/mgo#Session.SetMode
	//session.SetMode(mgo.Monotonic, true)
	Database = session.DB(AuthDatabase)
}
