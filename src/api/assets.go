package api

import (
	"db"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"models"
	"net/http"
	"strings"
)

func SearchAssets(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var asset models.Asset
	asset.DB = db.Database
	asset.DBSession = asset.DB.Session.Copy()
	defer asset.DBSession.Close()

	assets := asset.SearchForAsset(strings.ToLower(ps.ByName("query_term")))
	json.NewEncoder(w).Encode(assets)
}
