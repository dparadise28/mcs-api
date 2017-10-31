package api

import (
	"db"
	"encoding/json"
	"errors"
	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2/bson"
	"log"
	"models"
	"net/http"
	"time"
	"tools"
)

func UpdateCartProductQuantity(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var cart models.Cart
	cart.DB = db.Database
	cart.DBSession = cart.DB.Session.Copy()
	defer cart.DBSession.Close()

	var cartReq models.CartRequest
	v := new(tools.DefaultValidator)
	log.Println(r.Body)
	if err := json.NewDecoder(r.Body).Decode(&cartReq); err != nil {
		models.WriteNewError(w, err)
		return
	}
	log.Println(cart)
	if validationErr := v.ValidateIncomingJsonRequest(&cartReq); validationErr.Status != 200 {
		models.WriteError(w, &validationErr)
		return
	}

	cart.StoreID = cartReq.SID
	cart.UserID = bson.ObjectIdHex(r.Header.Get(models.USERID_COOKIE_NAME))
	cart.LastUpdated = time.Now()
	cart.CartState = models.CartStates["ACTIVE"]
	if count, count_err := cart.ActiveUserCartCountForStore(); count == 0 {
		var s models.Store
		cart.ID = bson.NewObjectId()
		cart.DateCreated = time.Now()
		s_collection := cart.DB.C(models.StoreCollectionName).With(cart.DBSession)
		if err := s_collection.Find(bson.M{
			"_id": cart.StoreID,
		}).One(&s); err != nil {
			models.WriteNewError(w, err)
			return
		}
		cart.StoreTaxRate = s.TaxRate
		cart.DeliveryFee = s.Delivery.Fee
		cart.StoreName = s.Name
		cart.IsNew = true
	} else if count_err != nil {
		models.WriteNewError(w, count_err)
		return
	} else if count == 1 && cartReq.CID.Hex() == "" {
		models.WriteNewError(w, errors.New("Must provide a valid cart id"))
		return
	} else {
		cart.ID = cartReq.CID
	}
	cart.UpdateProductQuantity(cartReq.PID, cartReq.Instructions, cartReq.QTY)
	if cart.Products == nil {
		models.WriteError(w, models.ErrBadRequest)
		return
	}
	RetrieveUserActiveCarts(w, r, ps)
	//json.NewEncoder(w).Encode(cart)
}

func RetrieveUserActiveCarts(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	RetrieveUserCartsByStatus(w, r, ps, models.CartStates["ACTIVE"])
}

func RetrieveUserCompletedCarts(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	RetrieveUserCartsByStatus(w, r, ps, models.CartStates["COMPLETED"])
}

func RetrieveUserCartsByStatus(w http.ResponseWriter, r *http.Request, ps httprouter.Params, status int) {
	var cart models.Cart
	cart.DB = db.Database
	cart.DBSession = cart.DB.Session.Copy()
	cart.UserID = bson.ObjectIdHex(r.Header.Get(models.USERID_COOKIE_NAME))
	cart.CartState = status
	defer cart.DBSession.Close()

	carts, err := cart.RetrieveUserCartsByStatus()
	if err != nil {
		models.WriteNewError(w, err)
		return
	}
	json.NewEncoder(w).Encode(carts)
}

func AbandonUserActiveCart(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var cart models.Cart
	cart.DB = db.Database
	cart.DBSession = cart.DB.Session.Copy()
	defer cart.DBSession.Close()

	cart.ID = bson.ObjectIdHex(ps.ByName("cart_id"))
	cart.UserID = bson.ObjectIdHex(r.Header.Get(models.USERID_COOKIE_NAME))
	if err := cart.AbandonCart(); err != nil {
		models.WriteNewError(w, err)
		return
	}
	json.NewEncoder(w).Encode(cart)
}

func ReActiveUserCart(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var cart models.Cart
	cart.DB = db.Database
	cart.DBSession = cart.DB.Session.Copy()
	defer cart.DBSession.Close()

	cart.ID = bson.ObjectIdHex(ps.ByName("cart_id"))
	cart.UserID = bson.ObjectIdHex(r.Header.Get(models.USERID_COOKIE_NAME))
	if err := cart.ReActivateCart(); err != nil {
		models.WriteNewError(w, err)
		return
	}
	json.NewEncoder(w).Encode(cart)
}

func RetrieveStoreCartByID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var cart models.Cart
	cart.DB = db.Database
	cart.DBSession = cart.DB.Session.Copy()
	defer cart.DBSession.Close()

	cart.ID = bson.ObjectIdHex(ps.ByName("cart_id"))
	cart.StoreID = bson.ObjectIdHex(r.Header.Get(models.STOREID_HEADER_NAME))
	if err := cart.RetrieveStoreCartByID(); err != nil {
		models.WriteNewError(w, err)
		return
	}
	json.NewEncoder(w).Encode(cart)
}
