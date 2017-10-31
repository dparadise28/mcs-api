package api

import (
	"db"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2/bson"
	// "log"
	"models"
	"net/http"
	"tools"
)

/*------------------template manipulation retrieval--------------------*/

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

// ----------------------------------------------------------------------*/
/*
func AddStoreCategoriesFromTemplate(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var reqCats []models.Category
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
	err, categories := category.AddStoreCategoriesFromTemplate(
		reqCats,
		bson.ObjectIdHex(r.Header.Get(models.STOREID_HEADER_NAME)),
	)
	if err != nil {
		models.WriteNewError(w, err)
		return
	}
	json.NewEncoder(w).Encode(categories)
}

func RetrieveEnabledTier1StoreCategories(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	RetrieveTier1StoreCategories(w, r, ps, true)
}

func RetrieveDisabledTier1StoreCategories(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	RetrieveTier1StoreCategories(w, r, ps, false)
}

func RetrieveTier1StoreCategories(w http.ResponseWriter, r *http.Request, ps httprouter.Params, enabled bool) {
	var category models.Category
	category.DB = db.Database
	category.DBSession = category.DB.Session.Copy()
	defer category.DBSession.Close()
	categories := category.FindStoreTier1Categories(
		bson.ObjectIdHex(r.Header.Get(models.STOREID_HEADER_NAME)), enabled,
	)
	json.NewEncoder(w).Encode(categories)
}

func RetrieveEnabledTier2StoreCategories(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	RetrieveTier2StoreCategories(w, r, ps, true)
}

func RetrieveDisabledTier2StoreCategories(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	RetrieveTier2StoreCategories(w, r, ps, false)
}

func RetrieveTier2StoreCategories(w http.ResponseWriter, r *http.Request, ps httprouter.Params, enabled bool) {
	var category models.Category
	category.DB = db.Database
	category.DBSession = category.DB.Session.Copy()
	defer category.DBSession.Close()
	categories := category.FindStoreTemplateTier2Categories(
		enabled,
		bson.ObjectIdHex(r.URL.Query().Get("category_id")),
		bson.ObjectIdHex(r.Header.Get(models.STOREID_HEADER_NAME)),
	)
	json.NewEncoder(w).Encode(categories)
}
*/

/*
func UpdateStoreCategory(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var category models.Category

	v := new(tools.DefaultValidator)
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		models.WriteError(w, models.ErrBadRequest)
		return
	}
	if validationErr := v.ValidateIncomingJsonRequest(&category); validationErr.Status != 200 {
		models.WriteError(w, &validationErr)
		return
	}
	category.DB = db.Database
	category.DBSession = category.DB.Session.Copy()
	defer category.DBSession.Close()
	category.UpdateStoreCategoryName()
	json.NewEncoder(w).Encode(category)
}

func EnableStoreCategory(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var category models.Category

	v := new(tools.DefaultValidator)
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		models.WriteError(w, models.ErrBadRequest)
		return
	}
	// set temporary name locally to stop validator from shitting itself
	category.Name = "temp"
	if validationErr := v.ValidateIncomingJsonRequest(&category); validationErr.Status != 200 {
		models.WriteError(w, &validationErr)
		return
	}
	category.DB = db.Database
	category.DBSession = category.DB.Session.Copy()
	defer category.DBSession.Close()
	if err := category.ActivateStoreCategory(); err != nil {
		models.WriteError(w, models.ErrBadRequest)
	}
	json.NewEncoder(w).Encode(category)
}

func ReorderStoreCategories(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var order models.CategoryOrder

	v := new(tools.DefaultValidator)
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		models.WriteError(w, models.ErrBadRequest)
		return
	}
	if validationErr := v.ValidateIncomingJsonRequest(&order); validationErr.Status != 200 {
		models.WriteError(w, &validationErr)
		return
	}
	order.DB = db.Database
	order.DBSession = order.DB.Session.Copy()
	defer order.DBSession.Close()
	if err := order.ReorderStoreCategories(); err != nil {
		models.WriteError(w, models.ErrBadRequest)
	}
	json.NewEncoder(w).Encode(order)
}
*/
