package models

import (
	"gopkg.in/mgo.v2/bson"
)

var CategoryCollectionName = "Categories"

type Category struct {
	ID       bson.ObjectId `bson:"_id,omitempty" json:"user_id"`
	Name     string        `bson:"name" json:"name" validate:"required"`
	StoreId  bson.ObjectId `bson:"store_id" json:"store_id" validate:"required"`
	Products []struct {
		ProductName string `bson:"name" json:"name"`
		Price       int
		ProductId   bson.ObjectId `bson:"p_id" json:"product_id"`
	} `bson:"products" json:"products"`
}
