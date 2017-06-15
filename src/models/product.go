package models

import "gopkg.in/mgo.v2/bson"

var ProductCollectionName = "Products"

type Product struct {
	ID             bson.ObjectId `bson:"_id" json:"product_id"`
	StoreID        bson.ObjectId `bson:"store_id" json:"store_id"`
	CategoryID     bson.ObjectId `bson:"category_id" json:"store_id"`
	ProductTitle   string        `bson:"title" json:"title" validate:"required"`
	Description    string        `bson:"desc" json:"description"`
	DisplayPrice   string        `bson:"display_price" json:"display_price" validate:"required"`
	ProductRatings struct {
		ReviewCount           uint64  `bson:"review_count" json:"total_reviews"`
		ReviewPercentageScore float64 `bson:"pct_score" json:"review_percent"`
	}
}
