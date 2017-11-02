package api

import (
	"db"
	"encoding/json"
	"errors"
	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2/bson"
	"log"
	"models"
	"net/http"
	"strconv"
	"strings"
	"tools"
)

func SearchAssets(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var asset models.Asset
	asset.DB = db.Database
	asset.DBSession = asset.DB.Session.Copy()
	defer asset.DBSession.Close()

	assets := asset.SearchForAsset(strings.ToLower(ps.ByName("query_term")))
	json.NewEncoder(w).Encode(assets)
}

func RetrieveTemplateAsset(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var asset models.Asset
	asset.DB = db.Database
	asset.DBSession = asset.DB.Session.Copy()
	defer asset.DBSession.Close()

	log.Println(r.URL.Query().Get("id"))
	asset.RetrieveTemplateAssetById(bson.ObjectIdHex(r.URL.Query().Get("id")))
	json.NewEncoder(w).Encode(asset)
}

func RetrieveTemplateAssets(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var assets models.PaginatedTemplateAssets
	assets.DB = db.Database
	assets.DBSession = assets.DB.Session.Copy()
	defer assets.DBSession.Close()

	assets.PG = 1
	if p, err := strconv.Atoi(r.URL.Query().Get("p")); err == nil {
		assets.PG = p
	}
	assets.Size = models.DefaultPageSize
	if s, err := strconv.Atoi(r.URL.Query().Get("size")); err == nil {
		if _, ok := models.PageSizes[s]; ok {
			assets.Size = s
		}
	}
	if r.Header.Get(models.STOREID_HEADER_NAME) != "" {
		assets.SID = bson.ObjectIdHex(r.Header.Get(models.STOREID_HEADER_NAME))
	} else {
		if r.URL.Query().Get("store_id") == "" {
			models.WriteNewError(w, errors.New("Please provide a valid store id."))
			return
		}
		assets.SID = bson.ObjectIdHex(r.URL.Query().Get("store_id"))
	}

	assets.CID = bson.ObjectIdHex(r.URL.Query().Get("category_id"))
	assets.RetrieveTemplateCategoryAssets()
	json.NewEncoder(w).Encode(assets)
}

func CreateAsset(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var asset models.NewAsset
	v := new(tools.DefaultValidator)
	if err := json.NewDecoder(r.Body).Decode(&asset); err != nil {
		models.WriteNewError(w, err)
		return
	}
	if validationErr := v.ValidateIncomingJsonRequest(&asset); validationErr.Status != 200 {
		models.WriteError(w, &validationErr)
		return
	}
	asset.DB = db.Database
	asset.DBSession = asset.DB.Session.Copy()
	defer asset.DBSession.Close()

	err, new_asset := asset.UploadAsset()
	if err != nil {
		models.WriteNewError(w, err)
		return
	}
	json.NewEncoder(w).Encode(new_asset)
}
