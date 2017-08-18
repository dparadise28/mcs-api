package db

import (
	mgo "gopkg.in/mgo.v2"
	"models"
)

var OrderIndex = []mgo.Index{
	mgo.Index{
		Key: []string{
			"user_id",
			"order_status",
		},
		Unique:     false,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	},
	mgo.Index{
		Key: []string{
			"store_id",
			"order_status",
		},
		Unique:     false,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	},
}

func EnsureOrderIndex() {
	session := Database.Session.Copy()
	defer session.Close()
	c := Database.C(models.OrderCollectionName).With(session)

	for _, index := range OrderIndex {
		err := c.EnsureIndex(index)
		if err != nil {
			panic(err)
		}
	}
}
