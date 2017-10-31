package db

import (
	mgo "gopkg.in/mgo.v2"
	"models"
)

var AssetIndex = []mgo.Index{
	mgo.Index{
		Key: []string{
			"template_category_id",
		},
		Unique:     false,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	},
	mgo.Index{
		Key: []string{
			"$text:title",
		},
		Unique:     false,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	},
}

func EnsureAssetIndex() {
	session := Database.Session.Copy()
	defer session.Close()

	c := Database.C(models.AssetCollectionName).With(session)
	for _, index := range AssetIndex {
		err := c.EnsureIndex(index)
		if err != nil {
			panic(err)
		}
	}
}
