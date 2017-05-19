package api

import (
	"db"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2/bson"
	"models"
	"net/http"
)

func StoreSearch(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Println(ps.ByName("store_id"))
}

func StoreCreate(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// TODO: move this out of specific calls and into generic
	// deserialization to avoid code duplication
	var store models.Store
	err := json.NewDecoder(r.Body).Decode(&store)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	// copy db session for the stores collection and close on completion
	session := db.Database.Session.Copy()
	defer session.Close()

	// grab the proper collection, create a new store id and attempt an insert
	c := db.Database.C(models.StoreCollectionName).With(session)
	store.ID = bson.NewObjectId()
	if insert_err := c.Insert(&store); insert_err != nil {
		http.Error(w, insert_err.Error(), 400)
		return
	}
	json.NewEncoder(w).Encode(store)
}
