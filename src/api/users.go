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
	//user.Password = tools.HashPassword(user.Password)
	userPassword := user.Password
	user.Password = ""

	// copy db session for the stores collection and close on completion
	session := db.Database.Session.Copy()
	c := db.Database.C(models.UserCollectionName).With(session)

	if insert_err := c.Insert(&user); insert_err != nil {
		models.WriteError(w, models.ErrResourceConflict)
		return
	}
	defer session.Close()
	go func() {
		defer session.Close()
		hashedPassword := tools.HashPassword(userPassword)
		change := mgo.Change{
			ReturnNew: true,
			Upsert:    false,
			Remove:    false,
			Update: bson.M{
				"$set": bson.M{
					"password": hashedPassword,
				},
			},
		}
		u, _ := c.Find(bson.M{"_id": user.ID}).Apply(change, &user)
		log.Println(u)
	}()

	email := user.EmailConfirmation()
	tools.EmailQueue <- &email
	user.ScrubSensitiveInfo()
	json.NewEncoder(w).Encode(user)
}

func UserInfo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var user models.User
	session := db.Database.Session.Copy()
	defer session.Close()

	c := db.Database.C(models.UserCollectionName).With(session)
	//c.Find(bson.M{"email": }).One(&user)
	c.Find(bson.M{"email": r.URL.Query().Get("email")}).One(&user)
	log.Println(user)
	json.NewEncoder(w).Encode(user)
}

func GetUserById(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var user models.User
	session := db.Database.Session.Copy()
	defer session.Close()

	c := db.Database.C(models.UserCollectionName).With(session)
	c.Find(bson.M{"_id": bson.ObjectIdHex(ps.ByName("user_id"))}).One(&user)
	json.NewEncoder(w).Encode(user)
}

func Login(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var user models.User
	session := db.Database.Session.Copy()
	defer session.Close()

	c := db.Database.C(models.UserCollectionName).With(session)
	c.Find(bson.M{"email": ps.ByName("email")}).One(&user)
	token, err := user.Roles.GenetateLoginToken()
	if err != nil {
		models.WriteError(w, models.ErrInternalServer)
		return
	}
	json.NewEncoder(w).Encode(models.LoginResponse{
		Token: token,
	})
}
