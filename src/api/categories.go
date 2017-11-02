package api

import (
	"db"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2/bson"
	// "log"
	"errors"
	"models"
	"net/http"
	"tools"
)

func AddStoreTemplateCategories(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var reqCats []models.StoreCategory
	v := new(tools.DefaultValidator)
	if err := json.NewDecoder(r.Body).Decode(&reqCats); err != nil {
		models.WriteNewError(w, err)
		return
	}
	if validationErr := v.ValidateIncomingJsonRequest(&reqCats); validationErr.Status != 200 {
		models.WriteError(w, &validationErr)
		return
	}

	var category models.Category
	category.DB = db.Database
	category.DBSession = category.DB.Session.Copy()
	defer category.DBSession.Close()
	root := bson.NewObjectId()
	err, categories := category.AddStoreCategories(
		reqCats,
		root, // bson.ObjectIdHex(r.Header.Get(models.STOREID_HEADER_NAME)),
		true,
	)
	if err != nil {
		models.WriteNewError(w, err)
		return
	}
	json.NewEncoder(w).Encode(categories)
}

func RetrieveTier1StoreTemplateCategories(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var category models.Category
	category.DB = db.Database
	category.DBSession = category.DB.Session.Copy()
	defer category.DBSession.Close()
	categories := category.FindStoreTemplateTier1Categories()
	// bson.ObjectIdHex(r.Header.Get(models.STOREID_HEADER_NAME)),
	json.NewEncoder(w).Encode(categories)
}

func RetrieveTier2StoreTemplateCategories(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var category models.Category
	category.DB = db.Database
	category.DBSession = category.DB.Session.Copy()
	defer category.DBSession.Close()
	categories := category.FindStoreTemplateTier2Categories(
		bson.ObjectIdHex(r.URL.Query().Get("category_id")),
		// bson.ObjectIdHex(r.Header.Get(models.STOREID_HEADER_NAME)),
	)
	json.NewEncoder(w).Encode(categories)
}

func RetrieveTier1StoreCategories(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var category models.Category
	if r.Header.Get(models.STOREID_HEADER_NAME) != "" {
		category.SID = bson.ObjectIdHex(r.Header.Get(models.STOREID_HEADER_NAME))
	} else {
		if r.URL.Query().Get("store_id") == "" {
			models.WriteNewError(w, errors.New("Please provide a valid store id."))
			return
		}
		category.SID = bson.ObjectIdHex(r.URL.Query().Get("store_id"))
	}
	category.DB = db.Database
	category.DBSession = category.DB.Session.Copy()
	defer category.DBSession.Close()
	categories := category.FindStoreTier1Categories()
	json.NewEncoder(w).Encode(categories)
}

func RetrieveTier2StoreCategories(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var category models.Category
	if r.Header.Get(models.STOREID_HEADER_NAME) != "" {
		category.SID = bson.ObjectIdHex(r.Header.Get(models.STOREID_HEADER_NAME))
	} else {
		if r.URL.Query().Get("store_id") == "" {
			models.WriteNewError(w, errors.New("Please provide a valid store id."))
			return
		}
		category.SID = bson.ObjectIdHex(r.URL.Query().Get("store_id"))
	}
	if r.URL.Query().Get("category_id") == "" {
		models.WriteNewError(w, errors.New("Please provide a valid category id."))
	}
	category.CID = bson.ObjectIdHex(r.URL.Query().Get("category_id"))
	category.DB = db.Database
	category.DBSession = category.DB.Session.Copy()
	defer category.DBSession.Close()
	categories := category.FindStoreTier2Categories()
	json.NewEncoder(w).Encode(categories)
}

func AddStoreCategories(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var categories models.StoreCategoryIds

	v := new(tools.DefaultValidator)
	if err := json.NewDecoder(r.Body).Decode(&categories); err != nil {
		models.WriteError(w, models.ErrBadRequest)
		return
	}
	if validationErr := v.ValidateIncomingJsonRequest(&categories); validationErr.Status != 200 {
		models.WriteError(w, &validationErr)
		return
	}
	categories.DB = db.Database
	categories.DBSession = categories.DB.Session.Copy()
	defer categories.DBSession.Close()
	if err := categories.AddStoreCategories(bson.ObjectIdHex(r.Header.Get(models.STOREID_HEADER_NAME))); err != nil {
		models.WriteError(w, models.ErrBadRequest)
	}
	json.NewEncoder(w).Encode(categories)
}
