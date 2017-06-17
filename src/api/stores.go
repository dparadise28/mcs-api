package api

import (
	"db"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"strconv"
	//"gopkg.in/mgo.v2/bson"
	"models"
	"net/http"
	"tools"
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
	_, resp := store.FindStoresByLocation(lon_float, lat_float, models.MAX_DISTANCE, time_int)
	json.NewEncoder(w).Encode(resp)
}

func GetStoreById(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var store models.Store
	store.DB = db.Database
	store.DBSession = store.DB.Session.Copy()
	defer store.DBSession.Close()

	_, resp := store.RetrieveFullStoreByID(ps.ByName("store_id"))
	json.NewEncoder(w).Encode(resp)
}

func StoreCreate(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var store models.Store
	store.DB = db.Database
	store.DBSession = store.DB.Session.Copy()
	defer store.DBSession.Close()

	v := new(tools.DefaultValidator)
	if err := json.NewDecoder(r.Body).Decode(&store); err != nil {
		models.WriteError(w, models.ErrBadRequest)
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
	UserSetStoreOwnerPerms(w, r, store.ID.Hex())
	json.NewEncoder(w).Encode(store)
}
