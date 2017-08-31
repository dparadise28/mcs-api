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
	"tools"
)

func AddReview(w http.ResponseWriter, r *http.Request, ps httprouter.Params, reviewType string) {
	var review models.Review
	review.DB = db.Database
	review.DBSession = review.DB.Session.Copy()
	review.ReviewFor = reviewType
	review.UserId = bson.ObjectIdHex(r.Header.Get(models.USERID_COOKIE_NAME))
	defer review.DBSession.Close()

	v := new(tools.DefaultValidator)
	if err := json.NewDecoder(r.Body).Decode(&review); err != nil {
		models.WriteNewError(w, err)
		return
	}
	if validationErr := v.ValidateIncomingJsonRequest(&review); validationErr.Status != 200 {
		models.WriteError(w, &validationErr)
		return
	}
	if err := review.AddReview(); err != nil {
		log.Println(err)
		models.WriteNewError(w, err)
		return
	}
	json.NewEncoder(w).Encode(review)
}

func ReviewPlatform(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	AddReview(w, r, ps, models.PlatformReview)
}

func ReviewStore(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	AddReview(w, r, ps, models.StoreReview)
}

func ReviewOrder(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	AddReview(w, r, ps, models.OrderReview)
}

func GetStoreReviews(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	//review.CurrentPage = r.URL.Query().Get("p")
	if r.URL.Query().Get("store_id") == "" {
		models.WriteNewError(w, errors.New(
			"Please specify the store you would like to view the ratings of"+
				" by specifying te store id in the query params (store_id=id)",
		))
		return
	}
	var review models.Review
	review.StoreId = bson.ObjectIdHex(r.URL.Query().Get("store_id"))
	//if review.CurrentPage == "" {
	review.CurrentPage = 1
	//}
	review.DB = db.Database
	review.DBSession = review.DB.Session.Copy()
	defer review.DBSession.Close()

	err, resp := review.StoreReviews()
	if err != nil {
		models.WriteNewError(w, err)
		return
	}
	json.NewEncoder(w).Encode(resp)
}
