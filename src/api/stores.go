package api

import (
	"db"
	"encoding/json"
	//"errors"
	"github.com/julienschmidt/httprouter"
	"models"
	"net/http"
	"strconv"
	//"strings"
	"tools"
	//"gopkg.in/mgo.v2/bson"
	//"log"
)

func StoreSearch(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var store models.Store
	store.DB = db.Database
	store.DBSession = store.DB.Session.Copy()
	defer store.DBSession.Close()

	lon := r.URL.Query().Get("lon")
	lat := r.URL.Query().Get("lat")
	time := r.URL.Query().Get("time")
	if lon == "" || lat == "" || time == "" {
		models.WriteError(w, models.ErrBadRequest)
	}
	lon_float, _ := strconv.ParseFloat(lon, 1000000)
	lat_float, _ := strconv.ParseFloat(lat, 1000000)
	time_int, _ := strconv.Atoi(time)
	err, resp := store.FindStoresByLocation(lon_float, lat_float, models.MAX_DISTANCE, time_int)
	if err != nil {
		models.WriteNewError(w, err)
		return
	}
	json.NewEncoder(w).Encode(resp)
}

func GetStoreById(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var store models.Store
	store.DB = db.Database
	store.DBSession = store.DB.Session.Copy()
	defer store.DBSession.Close()

	store.RetrieveStoreByID(ps.ByName("store_id"))
	json.NewEncoder(w).Encode(store)
}

func GetFullStoreById(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var store models.Store
	store.DB = db.Database
	store.DBSession = store.DB.Session.Copy()
	defer store.DBSession.Close()

	err, resp := store.RetrieveFullStoreByID(ps.ByName("store_id"))
	if err != nil {
		models.WriteNewError(w, err)
		return
	}
	json.NewEncoder(w).Encode(resp)
}

func StoreCreate(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var store models.Store
	store.DB = db.Database
	store.DBSession = store.DB.Session.Copy()
	defer store.DBSession.Close()

	v := new(tools.DefaultValidator)
	if err := json.NewDecoder(r.Body).Decode(&store); err != nil {
		models.WriteNewError(w, err)
		return
	}
	if validationErr := v.ValidateIncomingJsonRequest(&store); validationErr.Status != 200 {
		models.WriteError(w, &validationErr)
		return
	}
	if insert_err := store.Insert(); insert_err != nil {
		models.WriteError(w, models.ErrResourceConflict)
		return
	}
	UserSetStoreOwnerPerms(w, r, store.ID.Hex(), store.Name)
	json.NewEncoder(w).Encode(store)
}

func StoreAddCCPaymentMethod(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var store models.Store
	store.DB = db.Database
	store.DBSession = store.DB.Session.Copy()
	defer store.DBSession.Close()

	v := new(tools.DefaultValidator)
	if err := json.NewDecoder(r.Body).Decode(&store); err != nil {
		models.WriteNewError(w, err)
		return
	}
	if validationErr := v.ValidateIncomingJsonRequest(&store.PaymentDetails); validationErr.Status != 200 {
		models.WriteError(w, &validationErr)
		return
	}
	act, err := CreateStoreStripeCustomAccountImpl(w, r, ps, &store)
	if err != nil {
		models.WriteNewError(w, err)
		return
	}
	store.PaymentDetails.StripeAccountID = act.ID
	store.AddStoreCCPaymentMethod()
	json.NewEncoder(w).Encode(store)
}

func StoreInfoUpdate(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var store models.Store
	store.DB = db.Database
	store.DBSession = store.DB.Session.Copy()
	defer store.DBSession.Close()

	v := new(tools.DefaultValidator)
	if err := json.NewDecoder(r.Body).Decode(&store); err != nil {
		models.WriteNewError(w, err)
		return
	}
	if validationErr := v.ValidateIncomingJsonRequest(&store); validationErr.Status != 200 {
		models.WriteError(w, &validationErr)
		return
	}
	store.UpdateStoreInfo()
	json.NewEncoder(w).Encode(store)
}
