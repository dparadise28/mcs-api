package api

import (
	"db"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"models"
	"net/http"
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

func CreateAsset(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var asset models.NewAsset
	v := new(tools.DefaultValidator)
	if err := json.NewDecoder(r.Body).Decode(&asset); err != nil {
		models.WriteError(w, models.ErrBadRequest)
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
