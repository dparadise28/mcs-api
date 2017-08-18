package api

import (
	"db"
	"encoding/json"
	"errors"
	"github.com/julienschmidt/httprouter"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/account"
	"github.com/stripe/stripe-go/card"
	"github.com/stripe/stripe-go/customer"
	"gopkg.in/mgo.v2/bson"
	//"log"
	"models"
	"net/http"
	"strings"
	"time"
	"tools"
)

func PayWithCashForPickup(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var orderRequest models.CashPickupOrderRequest
	v := new(tools.DefaultValidator)
	if err := json.NewDecoder(r.Body).Decode(&orderRequest); err != nil {
		models.WriteNewError(w, err)
		return
	}
	if validationErr := v.ValidateIncomingJsonRequest(&orderRequest); validationErr.Status != 200 {
		models.WriteError(w, &validationErr)
		return
	}
	var order models.Order
	order.UserID = bson.ObjectIdHex(r.Header.Get(models.USERID_COOKIE_NAME))
	order.CartID = orderRequest.CardID
	order.StoreID = orderRequest.StoreID
	order.UserInstructions = orderRequest.UserInstructions
	order.PaymentMethod, order.OrderType = models.CASH, models.PICKUP
	order.DB, order.DBSession = db.Database, order.DB.Session.Copy()
	defer order.DBSession.Close()

	if err := order.ExpandOrderInfo(); err != nil {
		models.WriteNewError(w, err)
		return
	}
	if err := cart.CompleteCart(); err != nil {
		models.WriteNewError(w, err)
		return
	}
	if err := order.InsertOrder(); err != nil {
		models.WriteNewError(w, err)
		return
	}
	userEmail := order.UserOrderConfirmationEmail(false)
	storeEmail := order.UserOrderConfirmationEmail(true)

	tools.EmailQueue <- &userEmail
	tools.EmailQueue <- &storeEmail
	json.NewEncoder(w).Encode(order)
}

func PayWithCashForDelivery(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var orderRequest models.CashDeliveryOrderRequest
	v := new(tools.DefaultValidator)
	if err := json.NewDecoder(r.Body).Decode(&orderRequest); err != nil {
		models.WriteNewError(w, err)
		return
	}
	if validationErr := v.ValidateIncomingJsonRequest(&orderRequest); validationErr.Status != 200 {
		models.WriteError(w, &validationErr)
		return
	}
	var order models.Order
	order.UserID = bson.ObjectIdHex(r.Header.Get(models.USERID_COOKIE_NAME))
	order.CartID = orderRequest.CardID
	order.StoreID = orderRequest.StoreID
	order.AddressID = orderRequest.AddressID
	order.UserInstructions = orderRequest.UserInstructions
	order.PaymentMethod, order.OrderType = models.CASH, models.DELIVERY
	order.DB, order.DBSession = db.Database, order.DB.Session.Copy()
	defer order.DBSession.Close()

	if err := order.ExpandOrderInfo(); err != nil {
		models.WriteNewError(w, err)
		return
	}
	if err := cart.CompleteCart(); err != nil {
		models.WriteNewError(w, err)
		return
	}
	if err := order.InsertOrder(); err != nil {
		models.WriteNewError(w, err)
		return
	}
	userEmail := order.UserOrderConfirmationEmail(false)
	storeEmail := order.UserOrderConfirmationEmail(true)

	tools.EmailQueue <- &userEmail
	tools.EmailQueue <- &storeEmail
	json.NewEncoder(w).Encode(order)
}

func PayWithCCForPickup(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var orderRequest models.CashPickupOrderRequest
	v := new(tools.DefaultValidator)
	if err := json.NewDecoder(r.Body).Decode(&orderRequest); err != nil {
		models.WriteNewError(w, err)
		return
	}
	if validationErr := v.ValidateIncomingJsonRequest(&orderRequest); validationErr.Status != 200 {
		models.WriteError(w, &validationErr)
		return
	}
	var order models.Order
	order.UserID = bson.ObjectIdHex(r.Header.Get(models.USERID_COOKIE_NAME))
	order.CartID = orderRequest.CartID
	order.CardID = orderRequest.CardID
	order.StoreID = orderRequest.StoreID
	order.UserInstructions = orderRequest.UserInstructions
	order.PaymentMethod, order.OrderType = models.CASH, models.PICKUP
	order.DB, order.DBSession = db.Database, order.DB.Session.Copy()
	defer order.DBSession.Close()

	if err := order.ExpandOrderInfo(); err != nil {
		models.WriteNewError(w, err)
		return
	}
	stripe.Key = models.StripeSK

	params := &stripe.ChargeParams{
		Amount:        order.Cart.Totals.Total,
		Currency:      "usd",
		Fee:           orders.Cart.Totals.Total * .035,
		StripeAccount: order.Store.StripeAccountID,
		Capture:       false,
	}
	if order.CardID != "" {
		c, err := card.Get(
			order.CardID,
			&stripe.CardParams{Customer: order.User.StripeCustomerID},
		)
		if err != nil {
			json.NewEncoder(w).Encode(err)
			return
		}
		params.SetSource(order.CardID)
	} else {
		params.Customer = order.User.StripeCustomerID
	}
	charge, err := charge.New(params)
	if err != nil {
		models.WriteNewError(w, err)
		return
	}
	order.ChargeID = charge.ID

	if err := cart.CompleteCart(); err != nil {
		models.WriteNewError(w, err)
		return
	}
	if err := order.InsertOrder(); err != nil {
		models.WriteNewError(w, err)
		return
	}
	userEmail := order.UserOrderConfirmationEmail(false)
	storeEmail := order.UserOrderConfirmationEmail(true)

	tools.EmailQueue <- &userEmail
	tools.EmailQueue <- &storeEmail
	json.NewEncoder(w).Encode(order)
}
