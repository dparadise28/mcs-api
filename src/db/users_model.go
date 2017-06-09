package db

import (
	mgo "gopkg.in/mgo.v2"
	"models"
)

var UserIndex = []mgo.Index{
	// for mgo index struct type
	// http://bazaar.launchpad.net/+branch/mgo/v2/view/head:/session.go#L889
	mgo.Index{
		Key:        []string{"email"},
		Unique:     false,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	},
}

func EnsureUserIndex() {
	session := Database.Session.Copy()
	defer session.Close()

	// grab the proper collection, create a new store id and attempt an insert
	c := Database.C(models.UserCollectionName).With(session)

	// Index
	for _, index := range UserIndex {
		err := c.EnsureIndex(index)
		if err != nil {
			panic(err)
		}
	}
}
