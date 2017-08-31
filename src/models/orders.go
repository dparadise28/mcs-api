package models

import (
	"errors"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/charge"
	"github.com/stripe/stripe-go/refund"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	//"log"
	"time"
)

var OrderCollectionName = "Orders"

type CashPickupOrderRequest struct {
	CartID           bson.ObjectId `bson:"cart_id" json:"cart_id" validate:"required"`
	StoreID          bson.ObjectId `bson:"store_id" json:"store_id" validate:"required"`
	UserInstructions string        `bson:"instructions" json:"instructions"`
}

type CashDeliveryOrderRequest struct {
	CartID           bson.ObjectId   `bson:"cart_id" json:"cart_id" validate:"required"`
	StoreID          bson.ObjectId   `bson:"store_id" json:"store_id" validate:"required"`
	Address          DeliveryAddress `bson:"address" json:"address" validate:"required,dive"`
	UserInstructions string          `bson:"instructions" json:"instructions"`
}

type CCPickupOrderRequest struct {
	Tip              uint          `bson:"tip" json:"tip"`
	CardID           string        `bson:"card_id" json:"card_id"`
	CartID           bson.ObjectId `bson:"cart_id" json:"cart_id" validate:"required"`
	StoreID          bson.ObjectId `bson:"store_id" json:"store_id" validate:"required"`
	UserInstructions string        `bson:"instructions" json:"instructions"`
}

type CCDeliveryOrderRequest struct {
	Tip              uint            `bson:"tip" json:"tip"`
	CardID           string          `bson:"card_id" json:"card_id"`
	CartID           bson.ObjectId   `bson:"cart_id" json:"cart_id" validate:"required"`
	StoreID          bson.ObjectId   `bson:"store_id" json:"store_id" validate:"required"`
	Address          DeliveryAddress `bson:"address" json:"address" validate:"required,dive"`
	UserInstructions string          `bson:"instructions" json:"instructions"`
}

type OrderStatusUpdate struct {
	OrderID   bson.ObjectId `json:"order_id"`
	NewStatus string        `json:"new_status"`
	Message   string        `json:"message"`
}

type OrderStatusLog struct {
	Status string    `bson:"charge_status" json:"data"`
	Date   time.Time `bson:"date" json:"date"`
	Msg    string    `bson:"msg" json:"message"`
}

type DeliveryAddress struct {
	City             string  `bson:"city" json:"city"`
	Phone            string  `bson:"phone" json:"phone" validate:"required"`
	Line1            string  `bson:"line1" json:"line1" validate:"required"`
	Route            string  `bson:"route" json:"route"`
	Country          string  `bson:"country" json:"country" validate:"required"`
	AptSuite         string  `bson:"apt_suite" json:"apt_suite"`
	Latitude         float64 `bson:"latitude" json:"latitude" validate:"required,min=-85.0511499,max=85.001"`
	Longitude        float64 `bson:"longitude" json:"longitude" validate:"required,min=-180.001,max=180.001"`
	PostalCode       string  `bson:"postal_code" json:"postal_code" validate:"required"`
	StreetNumber     string  `bson:"street_number" json:"street_number"`
	AdminAreaLvl1    string  `bson:"administrative_area_level_1" json:"administrative_area_level_1"`
	FormattedAddress string  `bson:"formatted_address" json:"formatted_address"`
}

type Order struct {
	ID                 bson.ObjectId    `bson:"_id" json:"id"`
	Tip                uint             `bson:"tip" json:"tip"`
	CardID             string           `bson:"card_id" json:"card_id"`
	UserID             bson.ObjectId    `bson:"user_id" json:"user_id"`
	CartID             bson.ObjectId    `bson:"cart_id" json:"cart_id"`
	Address            DeliveryAddress  `bson:"address,omitempty" json:"address,omitempty"`
	StoreID            bson.ObjectId    `bson:"store_id" json:"store_id"`
	ChargeID           string           `bson:"charge_id" json:"charge_id"`
	RefundID           string           `bson:"refund_id" json:"refurn_id"`
	StatusLog          []OrderStatusLog `bson:"status_log" json:"status_log"`
	CreatedAt          time.Time        `bson:"created_at" json:"created_at"`
	OrderType          string           `bson:"order_type" json:"order_type"`
	OrderStatus        string           `bson:"order_status" json:"order_status"`
	PaymentMethod      string           `bson:"payment_method" json:"payment_method"`
	UserInstructions   string           `bson:"instructions" json:"instructions"`
	StoreMessageToUser string           `bson:"store_msg_to_user" json:"store_message_to_user"`

	// helpers
	DestinationEmail string `bson:"-" json:"-"`
	Store            Store  `bson:"-" json:"-"`
	Cart             Cart   `bson:"-" json:"-"`
	User             User   `bson:"-" json:"-"`
	NewStatus        string `bson:"-" json:"-"`

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

	o.Store.DB = o.DB
	o.Store.DBSession = o.DBSession
	if err := o.Store.RetrieveStoreByOID(); err != nil {
		return err
	}

	o.Cart.DB = o.DB
	o.Cart.DBSession = o.DBSession
	if err := o.Cart.GetActiveCartsById(); err != nil {
		return err
	}
	return nil
}

func (o *Order) InsertOrder() error {
	o.OrderStatus = PENDING
	o.ID = bson.NewObjectId()
	c := o.DB.C(OrderCollectionName).With(o.DBSession)
	o.CreatedAt = time.Now()
	o.StatusLog = []OrderStatusLog{
		OrderStatusLog{
			Status: o.OrderStatus,
			Date:   o.CreatedAt,
			Msg:    "Your order is currently being reviwed by the store.",
		},
	}
	return c.Insert(&o)
}

func (o *Order) IsValidNewStatus(status *OrderStatusUpdate) bool {
	for _, allowedNewStatus := range ALLOWED_STATUS_PATH[o.OrderType][o.OrderStatus] {
		if allowedNewStatus == status.NewStatus {
			return true
		}
	}
	return false
}

func (o *Order) DoExternalUpdatesForNewStatus(status *OrderStatusUpdate) (error, Email) {
	stripe.Key = StripeSK
	if o.PaymentMethod == CC {
		if status.NewStatus == REJECTED || status.NewStatus == CANCELED {
			refundParams := &stripe.RefundParams{Charge: o.ChargeID}
			refundParams.Params.StripeAccount = o.Store.PaymentDetails.StripeAccountID
			refundResp, err := refund.New(refundParams)
			if err != nil {
				return err, Email{}
			}
			o.RefundID = refundResp.ID
		}
		if status.NewStatus == COMPLETED {
			capture := &stripe.CaptureParams{}
			capture.Params.StripeAccount = o.Store.PaymentDetails.StripeAccountID
			_, err := charge.Capture(o.ChargeID, capture)
			if err != nil {
				return err, Email{}
			}
		}
	}
	c := o.DB.C(OrderCollectionName).With(o.DBSession)
	_, err := c.Find(bson.M{
		"_id": o.ID,
	}).Apply(mgo.Change{
		ReturnNew: true,
		Upsert:    false,
		Remove:    false,
		Update: bson.M{
			"$set": bson.M{
				"order_status": status.NewStatus,
				"refund_id":    o.RefundID,
			},
			"$push": bson.M{
				"status_log": OrderStatusLog{
					Status: status.NewStatus,
					Date:   time.Now(),
					Msg:    status.Message,
				},
			},
		},
	}, o)
	return err, o.CompletedEmail()
}

func (o *Order) UpdateOrderStatus(status *OrderStatusUpdate) (error, Email) {
	if o.IsValidNewStatus(status) {
		return o.DoExternalUpdatesForNewStatus(status)
	}
	return errors.New("Invalid status selected."), Email{}
}

func (o *Order) GetActiveByID() error {
	c := o.DB.C(OrderCollectionName).With(o.DBSession)
	return c.Find(bson.M{
		"_id": o.ID,
		"order_status": bson.M{
			"$in": []string{
				IN_PROGRESS,
				EN_ROUTE,
				APPROVED,
				PENDING,
			},
		},
	}).One(o)
}

func (o *Order) GetByID() error {
	c := o.DB.C(OrderCollectionName).With(o.DBSession)
	return c.Find(bson.M{
		"_id": o.ID,
	}).One(o)
}

func (o *Order) RetrieveActiveStoreOrders() ([]Order, error) {
	c := o.DB.C(OrderCollectionName).With(o.DBSession)
	orders := []Order{}
	err := c.Find(bson.M{
		"store_id": o.StoreID,
		"order_status": bson.M{
			"$in": []string{
				IN_PROGRESS,
				EN_ROUTE,
				APPROVED,
				PENDING,
			},
		},
	}).All(&orders)
	return orders, err
}

func (o *Order) RetrieveAllUserOrders() ([]Order, error) {
	c := o.DB.C(OrderCollectionName).With(o.DBSession)
	orders := []Order{}
	err := c.Find(bson.M{
		"user_id": o.UserID,
	}).All(&orders)
	return orders, err
}
