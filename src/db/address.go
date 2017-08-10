package db

import (
	mgo "gopkg.in/mgo.v2"
	"models"
)

var AddrIndex = []mgo.Index{
	mgo.Index{
		Key: []string{
			"user_id",
			"default",
		},
		Unique:     false,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	},
	mgo.Index{
		Key: []string{
			"user_id",
			"location.coordinates",
		},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	},
}

func EnsureAddressIndex() {
	session := Database.Session.Copy()
	defer session.Close()

	c := Database.C(models.AddressCollectionName).With(session)
	for _, index := range AddrIndex {
		err := c.EnsureIndex(index)
		if err != nil {
			panic(err)
		}
	}
}
