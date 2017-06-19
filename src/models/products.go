package models

import "gopkg.in/mgo.v2/bson"

var ProductCollectionName = "Products"

type Product struct {
	ID             bson.ObjectId `bson:"_id" json:"product_id"`
	StoreID        bson.ObjectId `bson:"store_id" json:"store_id"`
	SortOrder      uint16        `bson:"sort_order" json:"sort_order"`
	CategoryID     bson.ObjectId `bson:"category_id" json:"category_id"`
	Description    string        `bson:"desc" json:"description"`
	ProductTitle   string        `bson:"title" json:"title" validate:"required"`
	DisplayPrice   string        `bson:"-" json:"display_price"`
	PriceCents     uint32        `bson:"price_cents" json:"price_cents" validate:"required"`
	ProductRatings struct {
		ReviewCount           uint64  `bson:"review_count" json:"total_reviews"`
		ReviewPercentageScore float64 `bson:"pct_score" json:"review_percent"`
	}
}

type CartProduct struct {
	//StoreID      bson.ObjectId `bson:"-" json:"store_id"`
	ID           bson.ObjectId `bson:"_id" json:"product_id"`
	Quantity     uint16        `bson:"qty" json:"quantity"`
	PriceCents   uint32        `bson:"price_cents" json:"price_cents"`
	ProductTitle string        `bson:"title" json:"title"`
	Instructions string        `bson:"instructions" json:"instructions"`
}

type CartRequest struct {
	SID          bson.ObjectId `json:"store_id" validate:"required"`
	PID          bson.ObjectId `json:"product_id" validate:"required"`
	CID          bson.ObjectId `json:"cart_id"`
	QTY          uint16        `json:"quantity" validate:"required"`
	IsNewCart    bool          `json:"is_new_cart"`
	Instructions string        `json:"instructions"`
}
