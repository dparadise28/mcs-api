package models

import (
	"gopkg.in/mgo.v2/bson"
)

var CategoryCollectionName = "Categories"

type Category struct {
	ID        bson.ObjectId `bson:"_id,omitempty" json:"category_id"`
	Name      string        `bson:"name" json:"name" validate:"required"`
	StoreId   bson.ObjectId `bson:"store_id" json:"store_id" validate:"required"`
	SortOrder uint16        `bson:"sort_order" json:"sort_order"`
}

// helper model for unpacking to avoid writing boilerplate
type StoreCategory struct {
	ID        bson.ObjectId `bson:"_id,omitempty" json:"category_id"`
	Name      string        `bson:"name" json:"name" validate:"required"`
	StoreId   bson.ObjectId `bson:"store_id" json:"store_id"`
	SortOrder uint16        `bson:"sort_order" json:"sort_order"`
	Products  []Product     `bson:"-" json:"products" validate:"required"`
}
