package api

import (
	"db"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"gopkg.in/go-playground/validator.v9"
	"gopkg.in/mgo.v2/bson"
	"models"
	"net/http"
	"tools"
)

var validate *validator.Validate

func StoreSearch(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Println(ps.ByName("store_id"))
}

func StoreCreate(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var store models.Store
	v := new(tools.DefaultValidator)
	if err := json.NewDecoder(r.Body).Decode(&store); err != nil {
		models.WriteError(w, models.ErrBadRequest)
		return
	}
	if validationErr := v.ValidateIncomingJsonRequest(&store); validationErr.Status != 200 {
		models.WriteError(w, &validationErr)
		return
	}

	// copy db session for the stores collection and close on completion
	session := db.Database.Session.Copy()
	defer session.Close()

	// grab the proper collection, create a new store id and attempt an insert
	c := db.Database.C(models.StoreCollectionName).With(session)
	store.ID = bson.NewObjectId()
	store.Location.Type = "Point"
	if insert_err := c.Insert(&store); insert_err != nil {
		models.WriteError(w, models.ErrResourceConflict)
		return
	}
	json.NewEncoder(w).Encode(store)
}
