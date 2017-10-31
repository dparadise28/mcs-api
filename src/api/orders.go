// DRY the fuck out of this

package api

import (
	"db"
	"encoding/json"
	//"errors"
	"github.com/julienschmidt/httprouter"
	"github.com/stripe/stripe-go"
	//"github.com/stripe/stripe-go/account"
	"github.com/stripe/stripe-go/card"
	"github.com/stripe/stripe-go/charge"
	"github.com/stripe/stripe-go/token"
	//"github.com/stripe/stripe-go/customer"
	"gopkg.in/mgo.v2/bson"
	"log"
	"models"
	"net/http"
	//"strings"
	//"time"
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
	order.CartID = orderRequest.CartID
	order.StoreID = orderRequest.StoreID
	order.Address.Phone = orderRequest.Phone
	order.UserInstructions = orderRequest.UserInstructions
	order.PaymentMethod, order.OrderType = models.CASH, models.PICKUP
	order.DB = db.Database
	order.DBSession = order.DB.Session.Copy()
	defer order.DBSession.Close()

	if err := order.ExpandOrderInfo(); err != nil {
		models.WriteNewError(w, err)
		return
	}
	order.Cart.DB = order.DB
	order.Cart.DBSession = order.DBSession
	order.Cart.ApplyFee = false
	if err := order.Cart.CompleteCart(); err != nil {
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
	go tools.SendToSlack("notifications", "cash pickup.", order.ID.Hex())
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
	order.CartID = orderRequest.CartID
	order.StoreID = orderRequest.StoreID
	order.Address = orderRequest.Address
	order.UserInstructions = orderRequest.UserInstructions
	order.PaymentMethod, order.OrderType = models.CASH, models.DELIVERY
	order.DB = db.Database
	order.DBSession = order.DB.Session.Copy()
	defer order.DBSession.Close()

	b, _ := json.Marshal(order)
	log.Println(string(b))
	if err := order.ExpandOrderInfo(); err != nil {
		models.WriteNewError(w, err)
		return
	}
	order.Cart.DB = order.DB
	order.Cart.DBSession = order.DBSession
	if err := order.Cart.CompleteCart(); err != nil {
		models.WriteNewError(w, err)
		return
	}
	b, _ = json.Marshal(order)
	log.Println(string(b))
	if err := order.InsertOrder(); err != nil {
		models.WriteNewError(w, err)
		return
	}
	b, _ = json.Marshal(order)
	log.Println(string(b))
	userEmail := order.UserOrderConfirmationEmail(false)
	storeEmail := order.UserOrderConfirmationEmail(true)

	tools.EmailQueue <- &userEmail
	tools.EmailQueue <- &storeEmail
	go tools.SendToSlack("notifications", "cash delivery.", order.ID.Hex())
	json.NewEncoder(w).Encode(order)
}

func PayWithCCForPickup(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var orderRequest models.CCPickupOrderRequest
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
	order.Tip = orderRequest.Tip
	order.UserID = bson.ObjectIdHex(r.Header.Get(models.USERID_COOKIE_NAME))
	order.CartID = orderRequest.CartID
	order.CardID = orderRequest.CardID
	order.StoreID = orderRequest.StoreID
	order.Address.Phone = orderRequest.Phone
	order.UserInstructions = orderRequest.UserInstructions
	order.PaymentMethod, order.OrderType = models.CC, models.PICKUP
	order.DB = db.Database
	order.DBSession = order.DB.Session.Copy()
	defer order.DBSession.Close()

	if err := order.ExpandOrderInfo(); err != nil {
		models.WriteNewError(w, err)
		return
	}
	order.Cart.ApplyFee = false
	stripe.Key = models.StripeSK

	params := &stripe.ChargeParams{
		Fee:       uint64(order.Cart.Totals.Total * 0.048),
		Amount:    uint64(order.Cart.Totals.Total) + uint64(order.Tip),
		Currency:  "usd",
		NoCapture: true,
	}
	params.Params.StripeAccount = order.Store.PaymentDetails.StripeAccountID
	stripeSrc := &stripe.TokenParams{
		Customer: order.User.StripeCustomerID,
	}
	stripeSrc.Params.StripeAccount = order.Store.PaymentDetails.StripeAccountID
	if order.CardID != "" {
		c, err := card.Get(
			order.CardID,
			&stripe.CardParams{Customer: order.User.StripeCustomerID},
		)
		if err != nil {
			json.NewEncoder(w).Encode(err)
			return
		}
		stripeSrc.Card = &stripe.CardParams{
			Token: c.ID,
		}
		tok, tokErr := token.New(stripeSrc)
		if tokErr != nil {
			json.NewEncoder(w).Encode(tokErr)
			return
		}
		params.SetSource(tok.ID)
	} else {
		tok, tokErr := token.New(stripeSrc)
		if tokErr != nil {
			json.NewEncoder(w).Encode(tokErr)
			return
		}
		params.SetSource(tok.ID)
	}
	charge, err := charge.New(params)
	if err != nil {
		log.Println(charge, err)
		json.NewEncoder(w).Encode(err)
		return
	}
	order.ChargeID = charge.ID

	order.Cart.DB = order.DB
	order.Cart.ApplyFee = false
	order.Cart.DBSession = order.DBSession
	if err := order.Cart.CompleteCart(); err != nil {
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
	go tools.SendToSlack("notifications", "cc pickup! Fuck yeah!", order.ID.Hex())
	json.NewEncoder(w).Encode(order)
}

func PayWithCCForDelivery(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var orderRequest models.CCDeliveryOrderRequest
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
	order.Tip = orderRequest.Tip
	order.UserID = bson.ObjectIdHex(r.Header.Get(models.USERID_COOKIE_NAME))
	order.CartID = orderRequest.CartID
	order.CardID = orderRequest.CardID
	order.StoreID = orderRequest.StoreID
	order.Address = orderRequest.Address
	order.UserInstructions = orderRequest.UserInstructions
	order.PaymentMethod, order.OrderType = models.CC, models.DELIVERY
	order.DB = db.Database
	order.DBSession = order.DB.Session.Copy()
	defer order.DBSession.Close()

	if err := order.ExpandOrderInfo(); err != nil {
		models.WriteNewError(w, err)
		return
	}
	order.Cart.ApplyFee = true

	stripe.Key = models.StripeSK
	params := &stripe.ChargeParams{
		Fee:       uint64(order.Cart.Totals.Total * 0.048),
		Amount:    uint64(order.Cart.Totals.Total) + uint64(order.Tip),
		Currency:  "usd",
		NoCapture: true,
	}
	params.Params.StripeAccount = order.Store.PaymentDetails.StripeAccountID
	stripeSrc := &stripe.TokenParams{
		Customer: order.User.StripeCustomerID,
	}
	stripeSrc.Params.StripeAccount = order.Store.PaymentDetails.StripeAccountID
	if order.CardID != "" {
		c, err := card.Get(
			order.CardID,
			&stripe.CardParams{Customer: order.User.StripeCustomerID},
		)
		if err != nil {
			json.NewEncoder(w).Encode(err)
			return
		}
		stripeSrc.Card = &stripe.CardParams{
			Token: c.ID,
		}
		tok, tokErr := token.New(stripeSrc)
		if tokErr != nil {
			json.NewEncoder(w).Encode(tokErr)
			return
		}
		params.SetSource(tok.ID)
	} else {
		tok, tokErr := token.New(stripeSrc)
		if tokErr != nil {
			json.NewEncoder(w).Encode(tokErr)
			return
		}
		params.SetSource(tok.ID)
	}
	charge, err := charge.New(params)
	if err != nil {
		models.WriteNewError(w, err)
		return
	}
	order.ChargeID = charge.ID

	order.Cart.DB = order.DB
	order.Cart.DBSession = order.DBSession
	if err := order.Cart.CompleteCart(); err != nil {
		models.WriteNewError(w, err)
		return
	}
	if err := order.InsertOrder(); err != nil {
		models.WriteNewError(w, err)
		return
	}
	order.Cart.ApplyFee = true
	userEmail := order.UserOrderConfirmationEmail(false)
	storeEmail := order.UserOrderConfirmationEmail(true)

	tools.EmailQueue <- &userEmail
	tools.EmailQueue <- &storeEmail
	json.NewEncoder(w).Encode(order)
	go tools.SendToSlack("notifications", "cc delivery! Fuck yeah!", order.ID.Hex())
}

func GetActiveStoreOrders(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var order models.Order
	order.DB = db.Database
	order.DBSession = order.DB.Session.Copy()
	defer order.DBSession.Close()
	order.StoreID = bson.ObjectIdHex(ps.ByName("store_id"))
	orders, err := order.RetrieveActiveStoreOrders()
	if err != nil {
		models.WriteNewError(w, err)
		return
	}
	json.NewEncoder(w).Encode(orders)
}

func GetAllUserOrders(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var order models.Order
	order.DB = db.Database
	order.DBSession = order.DB.Session.Copy()
	defer order.DBSession.Close()
	order.UserID = bson.ObjectIdHex(r.Header.Get(models.USERID_COOKIE_NAME))
	orders, err := order.RetrieveAllUserOrders()
	if err != nil {
		models.WriteNewError(w, err)
		return
	}
	json.NewEncoder(w).Encode(orders)
}

func UpdateOrderStatus(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var status models.OrderStatusUpdate
	v := new(tools.DefaultValidator)
	if err := json.NewDecoder(r.Body).Decode(&status); err != nil {
		models.WriteNewError(w, err)
		return
	}
	if validationErr := v.ValidateIncomingJsonRequest(&status); validationErr.Status != 200 {
		models.WriteError(w, &validationErr)
		return
	}

	var order models.Order
	order.DB = db.Database
	order.DBSession = order.DB.Session.Copy()
	defer order.DBSession.Close()
	order.ID = status.OrderID
	if err := order.GetActiveByID(); err != nil {
		models.WriteNewError(w, err)
		return
	}
	b, _ := json.Marshal(order)
	log.Println(string(b))

	// abstract this shit; getting annoying as fuck to write every time
	order.DB = db.Database
	order.DBSession = order.DB.Session.Copy()
	defer order.DBSession.Close()
	order.Store.DB = order.DB
	order.Store.DBSession = order.DBSession
	order.Store.ID = order.StoreID
	if err := order.Store.RetrieveStoreByOID(); err != nil {
		models.WriteNewError(w, err)
		return
	}
	if order.UserID != bson.ObjectIdHex(r.Header.Get(models.USERID_COOKIE_NAME)) {
		order.DestinationEmail = order.Store.Email
	} else {
		order.User.DB = order.DB
		order.User.DBSession = order.DBSession
		order.User.ID = bson.ObjectIdHex(r.Header.Get(models.USERID_COOKIE_NAME))
		if err := order.User.GetById(); err != nil {
			models.WriteNewError(w, err)
			return
		}
		order.DestinationEmail = order.User.Email
	}

	err, email := order.UpdateOrderStatus(&status)
	if err != nil {
		models.WriteNewError(w, err)
		return
	}
	json.NewEncoder(w).Encode(order)
	if order.OrderStatus == models.CANCELED ||
		order.OrderStatus == models.REJECTED ||
		order.OrderStatus == models.COMPLETED {
		tools.EmailQueue <- &email
		go tools.SendToSlack(
			"notifications",
			"status updated to "+order.OrderStatus,
			order.ID.Hex(),
		)
	}
}
