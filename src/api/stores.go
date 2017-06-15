package api

import (
	"db"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2/bson"
	"models"
	"net/http"
	"tools"
)

func StoreSearch(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// TODO write queries and maybe split into different requests
	// search by location, search by category etc
	fmt.Println(ps.ByName("store_id"))
}

func GetStoreById(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var store models.Store
	session := db.Database.Session.Copy()
	defer session.Close()

	c := db.Database.C(models.StoreCollectionName).With(session)
	c.Find(bson.M{"_id": bson.ObjectIdHex(ps.ByName("store_id"))}).One(&store)
	json.NewEncoder(w).Encode(store)
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
