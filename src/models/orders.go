package models

import (
	"errors"
	"gopkg.in/mgo.v2"
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

	PENDING     = "PENDING"
	APPROVED    = "APPROVED"
	REJECTED    = "REJECTED"
	IN_PROGRESS = "IN-PROGRESS"
	CANCELED    = "CANCELED"
	EN_ROUTE    = "EN-ROUTE"
	COMPLETED   = "COMPLETED"

	FAILED_PAYMENT_HOLD    = "FAILED-PAYMENT-HOLD"
	PAYMENT_ON_HOLD        = "PAYMENT-ON-HOLD"
	FAILED_PAYMENT_CAPTURE = "FAILED-PAYMENT-CAPTURE"
	PAYMENT_CAPTURED       = "PAYMENT-CAPTURED"
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

type CashPickupOrderRequest struct {
	CartID           bson.ObjectId `bson:"cart_id" json:"id" validate:"required"`
	StoreID          bson.ObjectId `bson:"store_id" json:"store_id" validate:"required"`
	UserInstructions string        `bson:"instructions" json:"instructions"`
}

type CashDeliveryOrderRequest struct {
	CartID           bson.ObjectId `bson:"cart_id" json:"id" validate:"required"`
	StoreID          bson.ObjectId `bson:"store_id" json:"store_id" validate:"required"`
	AddressID        bson.ObjectId `bson:"address_id" json:"address_id" validate:"required"`
	UserInstructions string        `bson:"instructions" json:"instructions"`
}

type CCPickupOrderRequest struct {
	Tip              uint          `bson:"tip" json:"tip"`
	CardID           bson.ObjectId `bson:"card_id" json:"card_id"`
	CartID           bson.ObjectId `bson:"cart_id" json:"id" validate:"required"`
	StoreID          bson.ObjectId `bson:"store_id" json:"store_id" validate:"required"`
	UserInstructions string        `bson:"instructions" json:"instructions"`
}

type CCDeliveryOrderRequest struct {
	Tip              uint          `bson:"tip" json:"tip"`
	CardID           bson.ObjectId `bson:"card_id" json:"card_id"`
	CartID           bson.ObjectId `bson:"cart_id" json:"id" validate:"required"`
	StoreID          bson.ObjectId `bson:"store_id" json:"store_id" validate:"required"`
	AddressID        bson.ObjectId `bson:"address_id" json:"address_id" validate:"required"`
	UserInstructions string        `bson:"instructions" json:"instructions"`
}

type OrderPayment struct {
	ChargeStatus string `bson:"charge_status" json:"charge_status"`
	ErrorMessage string `bson:"error_message" json:"error_message"`
}

type Order struct {
	ID                 bson.ObjectId `bson:"_id" json:"id"`
	Tip                uint          `bson:"tip" json:"tip"`
	CardID             bson.ObjectId `bson:"card_id" json:"card_id"`
	UserID             bson.ObjectId `bson:"user_id" json:"user_id"`
	CartID             bson.ObjectId `bson:"cart_id" json:"id"`
	StoreID            bson.ObjectId `bson:"store_id" json:"store_id"`
	ChargeID           bson.ObjectId `bson:"charge_id" json:"charge_id"`
	AddressID          bson.ObjectId `bson:"address_id" json:"address_id"`
	OrderType          string        `bson:"order_type" json:"order_type"`
	OrderStatus        string        `bson:"order_status" json:"order_status"`
	PaymentMethod      string        `bson:"payment_method" json:"payment_method"`
	UserInstructions   string        `bson:"instructions" json:"instructions"`
	StoreMessageToUser string        `bson:"store_msg_to_user" json:"store_message_to_user"`
	Times              struct {
		LastUpdatedAt time.Time `bson:"last_updated_at" json:"last_updated_at"`
		CompletedAt   time.Time `bson:"completed_at" json:"completed_at"`
		CreatedAt     time.Time `bson:"fulfilled_at" json:"fulfilled_at"`
	} `bson:"times" json:"times"`

	// helpers
	Address Address `bson:"-" json:"-"`
	Store   Store   `bson:"-" json:"-"`
	Cart    Cart    `bson:"-" json:"-"`
	User    User    `bson:"-" json:"-"`

	DB        *mgo.Database `bson:"-" json:"-"`
	DBSession *mgo.Session  `bson:"-" json:"-"`
}

func (o *Order) ExpandOrderInfo() error {
	o.User.ID, o.Store.ID, o.Cart.ID = o.UserID, o.StoreID, o.CartID
	o.User.DB = o.DB
	o.User.DBSession = o.DBSession
	if err := o.User.GetById(); err != nil {
		return err
	}
	if o.User.StripeCustomerID == "" && o.CardID == "" && o.PaymentMethod != CC {
		errors.New("Please add a cc to your wallet")
	}

	o.store.DB = o.db
	o.store.DBSession = o.DBSession
	if err := o.Store.RetrieveStoreByOID(); err != nil {
		return err
	}

	o.Cart.DB = o.DB
	o.Cart.DBSession = o.DBSession
	if err := o.Cart.GetCartsById(); err != nil {
		return err
	}

	if o.OrderType == DELIVERY {
		o.Address.DB = o.DB
		o.Address.DBSession = o.DBSession
		if err := o.User.GetAddressById(); err != nil {
			return err
		}
	}
	o.OrderStatus = PENDING
	return nil
}

func (o *Order) InsertOrder() error {
	c := o.DB.C(OrderCollectionName).With(o.DBSession)
	return c.Insert(&o)
}
