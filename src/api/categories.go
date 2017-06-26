package api

import (
	"db"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"models"
	"net/http"
	"tools"
)

func GetCategoriesByStoreId(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var category models.Category
	category.DB = db.Database
	category.DBSession = category.DB.Session.Copy()
	defer category.DBSession.Close()

	enabled_only_categories := true
	enabled_only_products := true
	if r.URL.Query().Get("include_disabled_categories") == "true" {
		enabled_only_categories = false
	}
	if r.URL.Query().Get("include_disabled_products") == "true" {
		enabled_only_products = false
	}
	_, resp := category.RetrieveFullCategoriesByStoreID(
		ps.ByName("store_id"),
		enabled_only_categories,
		enabled_only_products,
	)
	json.NewEncoder(w).Encode(resp)
}

func AddStoreCategory(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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
	if err := category.AddStoreCategory(); err != nil {
		models.WriteError(w, models.ErrBadRequest)
	}
	json.NewEncoder(w).Encode(category)
}

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
