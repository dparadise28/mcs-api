package db

import (
	mgo "gopkg.in/mgo.v2"
	"models"
)

var ReviewIndex = []mgo.Index{
	mgo.Index{
		Key: []string{
			"review_for",
		},
		Unique:     false,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	},
	mgo.Index{
		Key: []string{
			"review_for",
			"store_id",
		},
		Unique:     false,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	},
	mgo.Index{
		Key: []string{
			"review_for",
			"product_id",
		},
		Unique:     false,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	},
	mgo.Index{
		Key: []string{
			"review_for",
			"user_id",
		},
		Unique:     false,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	},
}

func EnsureReviewIndex() {
	session := Database.Session.Copy()
	defer session.Close()

	// grab the proper collection, create a new store id and attempt an insert
	c := Database.C(models.ReviewCollectionName).With(session)

	// Index
	for _, index := range ReviewIndex {
		err := c.EnsureIndex(index)
		if err != nil {
			panic(err)
		}
	}
}
