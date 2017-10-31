package api

import (
	"db"
	"encoding/json"
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
	var asset models.TemplateAsset
	asset.DB = db.Database
	asset.DBSession = asset.DB.Session.Copy()
	defer asset.DBSession.Close()

	pg := 1
	if p, err := strconv.Atoi(r.URL.Query().Get("p")); err == nil {
		pg = p
	}
	assets := asset.RetrieveTemplateCategoryAssets(
		bson.ObjectIdHex(r.URL.Query().Get("category_id")), pg,
	)
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
