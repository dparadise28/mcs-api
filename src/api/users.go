package api

import (
	"crypto/rand"
	"db"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"gopkg.in/go-playground/validator.v9"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"models"
	"net/http"
	"tools"
)

func generateConfirmationCode() (string, bool) {
	n := 32
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", false
		log.Println(err)
	}
	s := fmt.Sprintf("%X", b)
	return s, true
}

func UserConfirmation(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var user models.User
	session := db.Database.Session.Copy()
	defer session.Close()

	// grab the proper collection, create a new store id and attempt an insert
	c := db.Database.C(models.UserCollectionName).With(session)
	user_oid := bson.ObjectIdHex(ps.ByName("user_id"))
	change := mgo.Change{
		Update:    bson.M{"$set": bson.M{"confirmed": true, "confirmation_code": ""}},
		Upsert:    false,
		Remove:    false,
		ReturnNew: true,
	}
	info, err := c.Find(bson.M{"_id": user_oid, "confirmation_code": ps.ByName("confirmation_code")}).Apply(change, &user)
	user.Password = ""
	log.Println(err)
	log.Println(info)
	json.NewEncoder(w).Encode(user)
}

func UserCreate(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// TODO: move this out of specific calls and into generic
	// deserialization to avoid code duplication
	var user models.User
	//var checked_user models.User
	validation := validator.New()
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	if err := validation.Struct(&user); err != nil {
		errors := []string{}
		for _, validationError := range err.(validator.ValidationErrors) {
			errors = append(errors, validationError.Namespace())
		}
		jsonErr, _ := json.Marshal(errors)
		http.Error(w, string(jsonErr), 400)
		return
	}
	// TODO handle errors
	confirmation_code, confimation_err := generateConfirmationCode()
	if !confimation_err {
		return
	}

	user.Confirmed = false
	user.ID = bson.NewObjectId()
	user.ConfirmationCode = confirmation_code

	// copy db session for the stores collection and close on completion
	session := db.Database.Session.Copy()
	defer session.Close()
	c := db.Database.C(models.UserCollectionName).With(session)

	/*if info, err := c.Find(bson.M{"username": user.UserName, "password": user.Password}).one(&checked_user); err == nil {
		http.Error(w, string(`{"success": false, "errors": ["User exists"]}`), 400)
		return
	}*/

	// grab the proper collection, create a new store id and attempt an insert
	if insert_err := c.Insert(&user); insert_err != nil {
		http.Error(w, insert_err.Error(), 400)
		return
	}
	email_subject := "Thank You for signing up!"
	email_body := "Welcome! Please click on the following link to confirm your account \n" +
		"http://mycorner.store:8001/api/user/confirm/email/" +
		user.ID.Hex() + "/" + confirmation_code
	if email_sent := tools.SendEmailValidation(user.Email, email_subject, email_body); !email_sent {
		http.Error(w, string(`{"success": false, "errors": ["Unable to reach the email address provided."]}`), 400)
		return
	}
	user.ConfirmationCode = ""
	json.NewEncoder(w).Encode(user)
}
