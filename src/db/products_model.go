package db

import (
	mgo "gopkg.in/mgo.v2"
	"models"
)

var ProductIndex = []mgo.Index{
	mgo.Index{
		Key: []string{
			"category_id",
			"store_id",
			"title",
		},
		Unique:     false,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	},
}

func EnsureProductIndex() {
	session := Database.Session.Copy()
	defer session.Close()

	// grab the proper collection, create a new store id and attempt an insert
	c := Database.C(models.ProductCollectionName).With(session)

	// Index
	for _, index := range ProductIndex {
		err := c.EnsureIndex(index)
		if err != nil {
			panic(err)
		}
	}
}
