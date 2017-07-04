package models

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var AssetCollectionName = "Assets"

type Asset struct {
	ID           bson.ObjectId `bson:"_id" json:"asset_id"`
	Link         string        `bson:"link" json:"link" validate:"required"`
	Title        string        `bson:"title" json:"title" validate:"required"`
	DisplayTitle string        `bson:"display_title" json:"display_title"`

	DB        *mgo.Database `bson:"-" json:"-"`
	DBSession *mgo.Session  `bson:"-" json:"-"`
}

type AutocompleteAsset struct {
	ID           bson.ObjectId `bson:"_id" json:"asset_id"`
	Link         string        `bson:"link" json:"image" validate:"required"`
	Title        string        `bson:"title" json:"title" validate:"required"`
	DisplayTitle string        `bson:"display_title" json:"label"`

	DB        *mgo.Database `bson:"-" json:"-"`
	DBSession *mgo.Session  `bson:"-" json:"-"`
}

func (a *Asset) SearchForAsset(queryTerm string) []AutocompleteAsset {
	var assets []AutocompleteAsset
	c := a.DB.C(AssetCollectionName).With(a.DBSession)
	c.Find(bson.M{
		"$text": bson.M{
			"$search": queryTerm,
		},
	}).Limit(100).All(&assets)
	if assets == nil {
		assets = []AutocompleteAsset{}
	}
	return assets
}
