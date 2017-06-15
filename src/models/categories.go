package models

import (
	"gopkg.in/mgo.v2/bson"
)

var CategoryCollectionName = "Categories"

type Category struct {
	ID         bson.ObjectId   `bson:"_id,omitempty" json:"category_id"`
	Name       string          `bson:"name" json:"name" validate:"required"`
	StoreId    bson.ObjectId   `bson:"store_id" json:"store_id" validate:"required"`
	ProductIDS []bson.ObjectId `bson:"p_ids" json:"product_ids"`
}

// helper model for unpacking to avoid writing boilerplate
type StoreCategory struct {
	ID         bson.ObjectId   `bson:"_id,omitempty" json:"category_id"`
	Name       string          `bson:"name" json:"name" validate:"required"`
	StoreId    bson.ObjectId   `bson:"store_id" json:"store_id"`
	ProductIDS []bson.ObjectId `bson:"p_ids" json:"product_ids"`
	Products   []Product       `bson:"-" json:"products" validate:"required"`
}
