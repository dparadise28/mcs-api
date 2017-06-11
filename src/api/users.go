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
	info, _ := c.Find(bson.M{
		"_id":               bson.ObjectIdHex(ps.ByName("user_id")),
		"confirmation_code": ps.ByName("confirmation_code"),
	}).Apply(change, &user)

	updatedClaims := models.GenerateTokenClaims(user.Roles.Access, user.Confirmed)
	updatedToken, _ := updatedClaims.CreateToken()
	authCookie := http.Cookie{
		Name:    "AUTH-TOKEN",
		Path:    "/api",
		Value:   updatedToken,
		Expires: models.COOKIE_TTL(),
	}
	http.SetCookie(w, &authCookie)
	user.ScrubSensitiveInfo()
	log.Println(info)
	json.NewEncoder(w).Encode(user)
}

func UserSetStoreOwnerPerms(w http.ResponseWriter, r *http.Request, storeId string) {
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
				"user_roles.access." + storeId: models.ACCESSROLE_STOREOWNER,
			},
		},
	}
	uid, _ := r.Cookie("UID")
	info, _ := c.Find(bson.M{
		"_id": bson.ObjectIdHex(uid.Value),
	}).Apply(change, &user)
	user.UpdateTokenAndCookie(w)

	user.ScrubSensitiveInfo()
	log.Println(info)
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
	userPassword := user.Password
	user.Password = ""

	// copy db session for the stores collection and close on completion
	session := db.Database.Session.Copy()
	c := db.Database.C(models.UserCollectionName).With(session)
	if insert_err := c.Insert(&user); insert_err != nil {
		models.WriteError(w, models.ErrResourceConflict)
		return
	}
	go func() {
		defer session.Close()
		hashedPassword := tools.HashPassword(userPassword)
		change := mgo.Change{
			ReturnNew: false,
			Upsert:    false,
			Remove:    false,
			Update: bson.M{
				"$set": bson.M{
					"password": hashedPassword,
				},
			},
		}
		if _, err := c.Find(bson.M{"_id": user.ID}).Apply(change, &user); err != nil {
			log.Printf(err.Error())
		}
	}()

	updatedClaims := models.GenerateTokenClaims(user.Roles.Access, user.Confirmed)
	updatedToken, _ := updatedClaims.CreateToken()
	authCookie := http.Cookie{
		Name:    "AUTH-TOKEN",
		Path:    "/api",
		Value:   updatedToken,
		Expires: models.COOKIE_TTL(),
	}
	http.SetCookie(w, &authCookie)
	uidCookie := http.Cookie{Name: "UID", Value: user.ID.Hex(), Expires: models.COOKIE_TTL(), Path: "/api"}
	http.SetCookie(w, &uidCookie)
	email := user.EmailConfirmation()
	tools.EmailQueue <- &email
	user.ScrubSensitiveInfo()
	json.NewEncoder(w).Encode(user)
}

func UserInfo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var user models.UserAPIResponse
	user.GetByEmail(db.Database, ps.ByName("email"))
	json.NewEncoder(w).Encode(user)
}

func GetUserById(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var user models.UserAPIResponse
	uid, _ := r.Cookie("UID")
	user.GetByIdStr(db.Database, uid.Value)
	json.NewEncoder(w).Encode(user)
}

func Login(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var user models.User
	user.GetByEmail(db.Database, r.URL.Query().Get("email"))
	token, err := user.GenetateLoginToken(r.URL.Query().Get("password"))
	if err != nil {
		if err.Error() == models.UNCONFIRMED_USER {
			models.WriteError(w, models.ErrUnconfirmedUser)
			return
		}
		models.WriteError(w, models.ErrUnauthorizedAccess)
		return
	}
	uidCookie := http.Cookie{Name: "UID", Value: user.ID.Hex(), Expires: models.COOKIE_TTL(), Path: "/api"}
	http.SetCookie(w, &uidCookie)
	authCookie := http.Cookie{Name: "AUTH-TOKEN", Value: token, Expires: models.COOKIE_TTL(), Path: "/api"}
	http.SetCookie(w, &authCookie)
	user.ScrubSensitiveInfo()
	json.NewEncoder(w).Encode(user)
}
