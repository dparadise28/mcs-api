package models

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

var OrderCollectionName = "Orders"

var OrderType = map[string]int{
	"delivery": 0,
	"pickup":   1,
}

var OrderStatus = map[string]int{
	"PENDING":     0,
	"REJECTED":    1,
	"APPROVED":    2,
	"IN-PROGRESS": 3,
	"EN-ROUTE":    4,
	"COMPLETED":   5,
	"CANCELED":    6,
}

type Order struct {
	// ::TODO:: update model with payment info when integrating stripe
	ID                 bson.ObjectId `bson:"_id" json:"id"`
	UserID             bson.ObjectId `bson:"user_id" json:"user_id"`
	CartID             bson.ObjectId `bson:"cart_id" json:"id"`
	StoreID            bson.ObjectId `bson:"store_id" json:"store_id"`
	StoreName          string        `bson:"store_name" json:"store_name"`
	OrderType          uint8         `bson:"order_type" json:"order_type"`
	OrderStatus        uint8         `bson:"order_status" json:"order_status"`
	DateCreated        time.Time     `bson:"created_at" json:"created_at"`
	ProductNames       []string      `bson:"product_names" json:"product_names"`
	DateFulfilled      time.Time     `bson:"fulfilled_at" json:"fulfilled_at"`
	UserInstructions   string        `bson:"instructions" json:"instructions"`
	StoreMessageToUser string        `bson:"store_msg2user" json:"store_message_to_user"`

	// index of address in users address book (in the user model now)
	DeliveryAddressIndex uint8 `bson:"address_index" json:"delivery_address_index"`
}
