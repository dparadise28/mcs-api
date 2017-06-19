package api

import (
	"db"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2/bson"
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
	if err := json.NewDecoder(r.Body).Decode(&cartReq); err != nil {
		models.WriteError(w, models.ErrBadRequest)
		return
	}
	if validationErr := v.ValidateIncomingJsonRequest(&cartReq); validationErr.Status != 200 {
		models.WriteError(w, &validationErr)
		return
	}

	cart.StoreID = cartReq.SID
	cart.UserID = bson.ObjectIdHex(r.Header.Get(models.USERID_COOKIE_NAME))
	cart.LastUpdated = time.Now()
	cart.Abandoned = false
	if cartReq.IsNewCart {
		var s models.Store
		cart.ID = bson.NewObjectId()
		cart.DateCreated = time.Now()
		s_collection := cart.DB.C(models.StoreCollectionName).With(cart.DBSession)
		if err := s_collection.Find(bson.M{"_id": cart.StoreID}).One(&s); err != nil {
			models.WriteError(w, models.ErrBadRequest)
			return
		}
		cart.StoreTaxRate = s.TaxRate
		cart.DeliveryFee = s.Delivery.Fee
		cart.IsNew = cartReq.IsNewCart
	} else if cartReq.CID.Hex() != "" && !cartReq.IsNewCart {
		cart.ID = cartReq.CID
	} else {
		models.WriteError(w, models.ErrBadRequest)
		return
	}
	cart.UpdateProductQuantity(cartReq.PID, cartReq.Instructions, cartReq.QTY)
	if cart.Products == nil {
		models.WriteError(w, models.ErrBadRequest)
		return
	}
	json.NewEncoder(w).Encode(cart)
}
