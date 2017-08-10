package models

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

const (
	OrderCollectionName = "Orders"

	PAYMENT_METHODS_KEY = "payment_methods"
	CASH                = "cash"
	CC                  = "stripe_cc"

	ORDER_METHODS_KEY = "order_methods"
	DELIVERY          = "delivery"
	PICKUP            = "pickup"

	ALL_STATUSES_KEY      = "status"
	DELIVERY_STATUSES_KEY = "delivery_status"
	PICKUP_STATUSES_KEY   = "pickup_status"
	CANCELED              = "CANCELED"
	IN_PROGRESS           = "IN-PROGRESS"
	COMPLETED             = "COMPLETED"
	REJECTED              = "REJECTED"
	EN_ROUTE              = "EN-ROUTE"
	APPROVED              = "APPROVED"
	PENDING               = "PENDING"
)

var (
	StripeSK = ""

	OrderMethod    = []string{DELIVERY, PICKUP}
	PaymentMethods = []string{CASH, CC}
	OrderStatuses  = []string{
		CANCELED,
		IN_PROGRESS,
		COMPLETED,
		REJECTED,
		EN_ROUTE,
		APPROVED,
		PENDING,
	}
	DeliveryOrderStatuses = OrderStatuses
	PickupOrderStatuses   = []string{
		CANCELED,
		IN_PROGRESS,
		COMPLETED,
		REJECTED,
		APPROVED,
		PENDING,
	}

	// helper map for all constants to allow for kw lookup
	CONST_MAP = map[string]map[string]string{
		PAYMENT_METHODS_KEY: map[string]string{
			CASH: CASH,
			CC:   CC,
		},
		ORDER_METHODS_KEY: map[string]string{
			DELIVERY: DELIVERY,
			PICKUP:   PICKUP,
		},
		ALL_STATUSES_KEY: map[string]string{
			CANCELED:    CANCELED,
			IN_PROGRESS: IN_PROGRESS,
			COMPLETED:   COMPLETED,
			REJECTED:    REJECTED,
			EN_ROUTE:    EN_ROUTE,
			APPROVED:    APPROVED,
			PENDING:     PENDING,
		},
		DELIVERY_STATUSES_KEY: map[string]string{
			CANCELED:    CANCELED,
			IN_PROGRESS: IN_PROGRESS,
			COMPLETED:   COMPLETED,
			REJECTED:    REJECTED,
			EN_ROUTE:    EN_ROUTE,
			APPROVED:    APPROVED,
			PENDING:     PENDING,
		},
		PICKUP_STATUSES_KEY: map[string]string{
			CANCELED:    CANCELED,
			IN_PROGRESS: IN_PROGRESS,
			COMPLETED:   COMPLETED,
			REJECTED:    REJECTED,
			APPROVED:    APPROVED,
			PENDING:     PENDING,
		},
	}
)

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
	PaymentMethod      uint8         `bson:"payment_method" json:"payment_method"`
	UserInstructions   string        `bson:"instructions" json:"instructions"`
	StoreMessageToUser string        `bson:"store_msg2user" json:"store_message_to_user"`

	// index of address in users address book (in the user model now)
	DeliveryAddressIndex uint8 `bson:"address_index" json:"delivery_address_index"`
}
