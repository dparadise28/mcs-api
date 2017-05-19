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
)

var validate *validator.Validate

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

	validation := validator.New()
	if err = validation.Struct(&store); err != nil {
		errors := []string{}
		for _, validationError := range err.(validator.ValidationErrors) {
			errors = append(errors, validationError.Namespace())
		}
		jsonErr, _ := json.Marshal(errors)
		http.Error(w, string(jsonErr), 400)
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
		http.Error(w, insert_err.Error(), 400)
		return
	}
	json.NewEncoder(w).Encode(store)
}
