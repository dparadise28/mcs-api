package api

import (
	"db"
	"encoding/json"
	"errors"
	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2/bson"
	"models"
	"net/http"
	"strconv"
	"tools"
)

func RetrieveStoreCategoryProducts(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var asset models.PaginatedProducts
	asset.DB = db.Database
	asset.DBSession = asset.DB.Session.Copy()
	defer asset.DBSession.Close()

	asset.PG = 1
	asset.Size = models.DefaultPageSize
	if p, err := strconv.Atoi(r.URL.Query().Get("p")); err == nil {
		asset.PG = p
	}
	if s, err := strconv.Atoi(r.URL.Query().Get("size")); err == nil {
		if _, ok := models.PageSizes[s]; ok {
			asset.Size = s
		}
	}
	asset.CID = bson.ObjectIdHex(r.URL.Query().Get("category_id"))
	if r.Header.Get(models.STOREID_HEADER_NAME) != "" {
		asset.SID = bson.ObjectIdHex(r.Header.Get(models.STOREID_HEADER_NAME))
	} else {
		if r.URL.Query().Get("store_id") == "" {
			models.WriteNewError(w, errors.New("Please provide a valid store id."))
			return
		}
		asset.SID = bson.ObjectIdHex(r.URL.Query().Get("store_id"))
	}
	asset.RetrieveStoreProductsByCategory()
	json.NewEncoder(w).Encode(asset)
}

func AddStoreProduct(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var product models.Product

	v := new(tools.DefaultValidator)
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		models.WriteError(w, models.ErrBadRequest)
		return
	}
	if validationErr := v.ValidateIncomingJsonRequest(&product); validationErr.Status != 200 {
		models.WriteError(w, &validationErr)
		return
	}
	product.DB = db.Database
	product.DBSession = product.DB.Session.Copy()
	defer product.DBSession.Close()
	if err := product.AddProductToStoreCategory(); err != nil {
		models.WriteError(w, models.ErrBadRequest)
	}
	json.NewEncoder(w).Encode(product)
}

func AddStoreProducts(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var product models.Product
	var newProducts []models.NewProduct

	v := new(tools.DefaultValidator)
	if err := json.NewDecoder(r.Body).Decode(&newProducts); err != nil {
		models.WriteError(w, models.ErrBadRequest)
		return
	}
	if validationErr := v.ValidateIncomingJsonRequest(&newProducts); validationErr.Status != 200 {
		models.WriteError(w, &validationErr)
		return
	}
	product.DB = db.Database
	product.DBSession = product.DB.Session.Copy()
	defer product.DBSession.Close()
	err, resp := product.AddProducts(
		newProducts,
		bson.ObjectIdHex(r.Header.Get(models.STOREID_HEADER_NAME)),
	)
	if err != nil {
		models.WriteError(w, models.ErrBadRequest)
	}
	json.NewEncoder(w).Encode(resp)
}

func UpdateStoreProduct(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var product models.Product

	v := new(tools.DefaultValidator)
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		models.WriteError(w, models.ErrBadRequest)
		return
	}
	if validationErr := v.ValidateIncomingJsonRequest(&product); validationErr.Status != 200 {
		models.WriteError(w, &validationErr)
		return
	}

	product.DB = db.Database
	product.DBSession = product.DB.Session.Copy()
	defer product.DBSession.Close()
	product.UpdateStoreProduct()
	json.NewEncoder(w).Encode(product)
}

func EnableStoreProduct(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var product models.Product

	v := new(tools.DefaultValidator)
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		models.WriteError(w, models.ErrBadRequest)
		return
	}
	// set temporary name locally to stop validator from shitting itself
	product.ProductTitle = "temp"
	product.PriceCents = 0
	if validationErr := v.ValidateIncomingJsonRequest(&product); validationErr.Status != 200 {
		models.WriteError(w, &validationErr)
		return
	}
	product.DB = db.Database
	product.DBSession = product.DB.Session.Copy()
	defer product.DBSession.Close()
	if err := product.ActivateStoreProduct(); err != nil {
		models.WriteError(w, models.ErrBadRequest)
	}
	json.NewEncoder(w).Encode(product)
}
