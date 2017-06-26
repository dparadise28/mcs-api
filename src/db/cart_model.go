package db

import (
	mgo "gopkg.in/mgo.v2"
	"models"
)

var CartIndex = []mgo.Index{
	mgo.Index{
		Key: []string{
			"user_id",
			"store_id",
			"active",
		},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	},
	mgo.Index{
		Key: []string{
			"store_id",
		},
		Unique:     false,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	},
}

func EnsureCartIndex() {
	session := Database.Session.Copy()
	defer session.Close()
	c := Database.C(models.CartCollectionName).With(session)

	for _, index := range CartIndex {
		err := c.EnsureIndex(index)
		if err != nil {
			panic(err)
		}
	}
}
