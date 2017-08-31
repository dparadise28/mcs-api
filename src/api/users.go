package api

import (
	"db"
	"encoding/json"
	"errors"
	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"models"
	"net/http"
	"strings"
	"tools"
)

func UserConfirmation(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var user models.User
	session := db.Database.Session.Copy()
	defer session.Close()

	c := db.Database.C(models.UserCollectionName).With(session)
	query := bson.M{"$set": bson.M{"confirmed": true, "confirmation_code": ""}}
	if r.URL.Query().Get("password") != "" {
		query = bson.M{
			"$set": bson.M{
				"confirmed":         true,
				"confirmation_code": "",
				"password":          r.URL.Query().Get("password"),
			},
		}
	}
	change := mgo.Change{
		ReturnNew: true,
		Upsert:    false,
		Remove:    false,
		Update:    query,
	}
	_, err := c.Find(bson.M{
		"_id": bson.ObjectIdHex(ps.ByName("user_id")),
		"$or": []bson.M{
			bson.M{
				"confirmed":         false,
				"confirmation_code": ps.ByName("confirmation_code"),
			},
			bson.M{
				"confirmed":         true,
				"confirmation_code": "",
			},
		},
	}).Apply(change, &user)
	if err != nil {
		log.Println(err, ps.ByName("user_id"), ps.ByName("confirmation_code"))
		models.WriteNewError(w, err)
		return
	}

	user.UpdateTokenAndCookie(w)
	user.ScrubSensitiveInfo()

	// ::TODO:: dynamic external routing for things like this
	log.Println("redirecting")
	http.Redirect(w, r, "http://mycorner.store:8003/#/login", http.StatusTemporaryRedirect)
	return
}

func UserSetStoreOwnerPerms(w http.ResponseWriter, r *http.Request, storeId string, storeName string) {
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
				"user_roles.access." + storeId:    models.ACCESSROLE_STOREOWNER,
				"user_roles.store_map." + storeId: storeName,
			},
		},
	}
	uid := r.Header.Get(models.USERID_COOKIE_NAME)
	info, _ := c.Find(bson.M{
		"_id": bson.ObjectIdHex(uid),
	}).Apply(change, &user)
	user.UpdateTokenAndCookie(w)
	user.ScrubSensitiveInfo()
	log.Println(info)
}

func UserResendConfirmation(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var user models.User
	session := db.Database.Session.Copy()
	defer session.Close()

	c := db.Database.C(models.UserCollectionName).With(session)

	change := mgo.Change{
		ReturnNew: true,
		Upsert:    false,
		Remove:    false,
		Update:    bson.M{"$set": bson.M{"confirmation_code": tools.GenerateConfirmationCode()}},
	}
	_, err := c.Find(bson.M{"email": strings.ToLower(r.URL.Query().Get("email"))}).Apply(change, &user)
	if err != nil {
		models.WriteNewError(w, err)
		return
	}
	user.UpdateTokenAndCookie(w)

	reset_pw := false
	if r.URL.Query().Get("password") != "" {
		user.Password = tools.HashPassword(r.URL.Query().Get("password"))
		reset_pw = true
	}
	email := user.EmailConfirmation(reset_pw)

	user.Password = ""
	tools.EmailQueue <- &email
	user.ConfirmationCode = ""
	json.NewEncoder(w).Encode(user)
}

func UserCreate(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var user models.User
	v := new(tools.DefaultValidator)
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		models.WriteNewError(w, err)
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
	user.Email = strings.ToLower(user.Email)

	session := db.Database.Session.Copy()
	defer session.Close()
	c := db.Database.C(models.UserCollectionName).With(session)
	if insert_err := c.Insert(&user); insert_err != nil {
		models.WriteNewError(w, insert_err) //models.ErrResourceConflict)
		log.Println(user)
		return
	}

	user.UpdateTokenAndCookie(w)
	email := user.EmailConfirmation(false)
	user.ScrubSensitiveInfo()
	tools.EmailQueue <- &email
	json.NewEncoder(w).Encode(user)
}

func UserInfo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var user models.UserAPIResponse
	user.GetByEmail(db.Database, strings.ToLower(ps.ByName("email")))
	json.NewEncoder(w).Encode(user)
}

func GetUserById(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var user models.UserAPIResponse
	uid := r.Header.Get(models.USERID_COOKIE_NAME)
	user.GetByIdStr(db.Database, uid)
	json.NewEncoder(w).Encode(user)
}

func Login(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var user models.User
	user.GetByEmail(db.Database, strings.ToLower(r.URL.Query().Get("email")))
	_, err := user.GenetateLoginTokenAndSetHeaders(r.URL.Query().Get("password"), w)
	if err != nil {
		if err.Error() == models.UNCONFIRMED_USER {
			models.WriteError(w, models.ErrUnconfirmedUser)
			return
		}
		models.WriteError(w, models.ErrUnauthorizedAccess)
		return
	}
	user.ScrubSensitiveInfo()
	json.NewEncoder(w).Encode(user)
}

func AddUserAddrToAddrBook(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var addr models.Address
	addr.DB = db.Database
	if err := json.NewDecoder(r.Body).Decode(&addr); err != nil {
		models.WriteNewError(w, err)
		return
	}
	v := new(tools.DefaultValidator)
	if validationErr := v.ValidateIncomingJsonRequest(&addr); validationErr.Status != 200 {
		models.WriteError(w, &validationErr)
		return
	}
	if addr.Name == "" {
		models.WriteNewError(w, errors.New(
			"Please provide the name associated with this address.",
		))
		return
	}
	addr.UserID = bson.ObjectIdHex(r.Header.Get(models.USERID_COOKIE_NAME))
	addr.DBSession = addr.DB.Session.Copy()
	defer addr.DBSession.Close()

	err := addr.AddToUserAddressBook()
	log.Println(err)
	if err != nil {
		models.WriteNewError(w, err)
		return
	}
	json.NewEncoder(w).Encode(addr)
}

func RemoveUserAddrFromAddrBook(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var addr models.Address
	addr.DB = db.Database
	addr.DBSession = addr.DB.Session.Copy()
	defer addr.DBSession.Close()

	if r.URL.Query().Get("address_id") == "" {
		models.WriteNewError(w, errors.New(
			"You must provide an address_id in the query params associated "+
				"with the address you'd like to remove.",
		))
	}
	addr.UserID = bson.ObjectIdHex(r.Header.Get(models.USERID_COOKIE_NAME))
	addr.ID = bson.ObjectIdHex(r.URL.Query().Get("address_id"))
	err := addr.RemoveFromUserAddressBook()
	if err != nil {
		models.WriteNewError(w, err)
		return
	}
	json.NewEncoder(w).Encode(addr)
}

func UpdateUserDefaultInAddrBook(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var addr models.Address
	addr.DB = db.Database
	addr.DBSession = addr.DB.Session.Copy()
	defer addr.DBSession.Close()

	if r.URL.Query().Get("address_id") == "" {
		models.WriteNewError(w, errors.New(
			"You must provide an address_id in the query params associated "+
				"with the address you'd like to set as your new default.",
		))
	}
	addr.UserID = bson.ObjectIdHex(r.Header.Get(models.USERID_COOKIE_NAME))
	addr.NewDefaultID = bson.ObjectIdHex(r.URL.Query().Get("address_id"))
	err := addr.ChangeDefaultAddressInUserAddressBook()
	if err != nil {
		models.WriteNewError(w, err)
		return
	}
	json.NewEncoder(w).Encode(addr)
}

func GetUserAddrBook(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var addr models.Address
	addr.DB = db.Database
	addr.DBSession = addr.DB.Session.Copy()
	defer addr.DBSession.Close()

	addr.UserID = bson.ObjectIdHex(r.Header.Get(models.USERID_COOKIE_NAME))
	addrs, err := addr.RetrieveUserAddressBook()
	if err != nil {
		models.WriteNewError(w, err)
		return
	}
	json.NewEncoder(w).Encode(addrs)
}

func GetUserDefaultAddr(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var addr models.Address
	addr.DB = db.Database
	addr.DBSession = addr.DB.Session.Copy()
	defer addr.DBSession.Close()

	addr.UserID = bson.ObjectIdHex(r.Header.Get(models.USERID_COOKIE_NAME))
	err := addr.RetrieveUserDefaultAddressInAddressBook()
	if err != nil {
		models.WriteNewError(w, err)
		return
	}
	json.NewEncoder(w).Encode(addr)
}
