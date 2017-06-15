package db

import (
	mgo "gopkg.in/mgo.v2"
	"models"
)

var CategoryIndex = []mgo.Index{
	// for mgo index struct type
	// http://bazaar.launchpad.net/+branch/mgo/v2/view/head:/session.go#L889
	mgo.Index{
		Key: []string{
			"name",
			"store_id",
		},
		Unique:     false,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	},
}

func EnsureCategoryIndex() {
	session := Database.Session.Copy()
	defer session.Close()

	// grab the proper collection, create a new store id and attempt an insert
	c := Database.C(models.CategoryCollectionName).With(session)

	// Index
	for _, index := range CategoryIndex {
		err := c.EnsureIndex(index)
		if err != nil {
			panic(err)
		}
	}
}
