package api

import (
	"db"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	//"gopkg.in/go-playground/validator.v9"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"models"
	"net/http"
	"tools"
)

func UserConfirmation(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	//panic("uhoh")
	var user models.User
	session := db.Database.Session.Copy()
	defer session.Close()

	c := db.Database.C(models.UserCollectionName).With(session)
	change := mgo.Change{
		ReturnNew: true,
		Upsert:    false,
		Remove:    false,
		Update: bson.M{
			"$set": bson.M{
				"confirmed": true, "confirmation_code": "",
			},
		},
	}
	info, err := c.Find(bson.M{
		"_id":               bson.ObjectIdHex(ps.ByName("user_id")),
		"confirmation_code": ps.ByName("confirmation_code"),
	}).Apply(change, &user)

	user.ScrubSensitiveInfo()
	log.Println(err)
	log.Println(info)
	json.NewEncoder(w).Encode(user)
}

func UserCreate(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var user models.User
	v := new(tools.DefaultValidator)
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		models.WriteError(w, models.ErrBadRequest)
		return
	}
	if validationErr := v.ValidateIncomingJsonRequest(&user); validationErr.Status != 200 {
		models.WriteError(w, &validationErr)
		return
	}

	user.ID = bson.NewObjectId()
	user.Confirmed = false
	user.ConfirmationCode = tools.GenerateConfirmationCode()
	user.Password = tools.HashPassword(user.Password)

	// copy db session for the stores collection and close on completion
	session := db.Database.Session.Copy()
	defer session.Close()
	c := db.Database.C(models.UserCollectionName).With(session)

	// grab the proper collection, create a new store id and attempt an insert
	if insert_err := c.Insert(&user); insert_err != nil {
		models.WriteError(w, models.ErrResourceConflict)
		return
	}
	email_subject := "Thank You for signing up!"
	email_body := "Welcome! Please click on the following link to confirm your account \n" +
		"http://mycorner.store:8001/api/user/confirm/email/" +
		user.ID.Hex() + "/" + user.ConfirmationCode
	if email_sent := tools.SendEmailValidation(user.Email, email_subject, email_body); !email_sent {
		unreachable := *models.ErrRequestTimeout
		unreachable.Details["timeout"] = append(
			unreachable.Details["timeout"], "Unable to reach the email address provided.",
		)
		models.WriteError(w, &unreachable)
		return
	}

	user.ScrubSensitiveInfo()
	json.NewEncoder(w).Encode(user)
}
