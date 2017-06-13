package models

import (
	"gopkg.in/mgo.v2/bson"
)

var ReviewCollectionName = "Reviews"

var ReviewTypes = []string{
	"Product",
	"Store",
}

var ReviewTypesMap = map[uint8]string{
	0: ReviewTypes[0],
	1: ReviewTypes[1],
}

type Review struct {
	ReviewId   bson.ObjectId `bson:"_id" json:"id"`
	ProductId  bson.ObjectId `bson:"product_id" json:"product_id"`
	StoreId    bson.ObjectId `bson:"store_id" json:"store_id"`
	UserId     bson.ObjectId `bson:"user_id" json:"user_id"`
	Score      uint8         `bson:"score" json:"score" validate:"required"`
	Comment    string        `bson:"comment" json:"comment"`
	ReviewType uint8         `bson:"review_type" json:"review_type"`
}
