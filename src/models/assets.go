package models

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var AssetCollectionName = "Assets"

type Asset struct {
	ID           bson.ObjectId `bson:"_id" json:"asset_id"`
	Link         string        `bson:"link" json:"link" validate:"required"`
	Title        string        `bson:"title" json:"display_title" validate:"required"`
	DisplayTitle string        `bson:"display_title" json:"title"`

	DB        *mgo.Database `bson:"-" json:"-"`
	DBSession *mgo.Session  `bson:"-" json:"-"`
}

func (a *Asset) SearchForAsset(queryTerm string) []Asset {
	var assets []Asset
	c := a.DB.C(AssetCollectionName).With(a.DBSession)
	c.Find(bson.M{
		"$text": bson.M{
			"$search": queryTerm,
		},
	}).Limit(100).All(&assets)
	if assets == nil {
		assets = []Asset{}
	}
	return assets
}
