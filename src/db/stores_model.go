package db

import (
	mgo "gopkg.in/mgo.v2"
	"models"
)

var StoreIndex = []mgo.Index{
	// for mgo index struct type
	// http://bazaar.launchpad.net/+branch/mgo/v2/view/head:/session.go#L889
	mgo.Index{
		Key: []string{
			"name",
			"address.location.coordinates",
		},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	},
	mgo.Index{
		Key:        []string{"$2dsphere:address.location"},
		Unique:     false,
		DropDups:   false,
		Background: true,
		Sparse:     true,
	},
}

func EnsureStoreIndex() {
	session := Database.Session.Copy()
	defer session.Close()

	// grab the proper collection, create a new store id and attempt an insert
	c := Database.C(models.StoreCollectionName).With(session)

	// Index
	for _, index := range StoreIndex {
		err := c.EnsureIndex(index)
		if err != nil {
			panic(err)
		}
	}
}
