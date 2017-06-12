package models

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var UserCollectionName = "Users"

type UserAddress struct {
	AddressName string  `bson:"address_name" json:"address_name"`
	Location    GeoJson `bson:"location" json:"location"`
}

type User struct {
	ID               bson.ObjectId `bson:"_id,omitempty" json:"user_id"`
	ConfirmationCode string        `bson:"confirmation_code" json:"confirmation_code"`
	AddressBook      []UserAddress `bson:"address_book" json:"address_book"`
	Confirmed        bool          `bson:"confirmed" json:"confirmed"`
	Password         string        `bson:"password" json:"password" validate:"required"`
	UserName         string        `bson:"username" json:"username" validate:"required"`
	Email            string        `bson:"email" json:"email" validate:"required,email"`
	Roles            UserRoles     `bson:"user_roles" json:"user_roles"`
	// user roles struct is in /src/models/auth

	Login struct {
		Token string        `bson:"-" json:"AUTH-TOKEN"`
		UID   bson.ObjectId `bson:"-" json:"UID"`
	} `bson:"-" json:"login"`
}

type UserAPIResponse struct {
	ID          bson.ObjectId `bson:"_id,omitempty" json:"user_id"`
	AddressBook []UserAddress `bson:"address_book" json:"address_book"`
	Confirmed   bool          `bson:"confirmed" json:"confirmed"`
	UserName    string        `bson:"username" json:"username" validate:"required"`
	Email       string        `bson:"email" json:"email" validate:"required,email"`
	Roles       UserRoles     `bson:"user_roles" json:"user_roles"`
}

func (u *User) EmailConfirmation() Email {
	// Email struct model can be found in the general.go file
	emailSubject := "Thank You for signing up!"
	emailBody := "Welcome! Please click on the following link to confirm your account \n" +
		"http://mycorner.store:8001/api/user/confirm/email/" + u.ID.Hex() + "/" + u.ConfirmationCode
	return Email{u.Email, emailBody, emailSubject}
}

func (u *User) ScrubSensitiveInfo() {
	u.Password, u.ConfirmationCode = "", ""
}

func (u *User) GetByEmail(db *mgo.Database, email string) {
	session := db.Session.Copy()
	defer session.Close()

	c := db.C(UserCollectionName).With(session)
	c.Find(bson.M{"email": email}).One(u)
}

func (u *User) GetByIdStr(db *mgo.Database, id string) {
	session := db.Session.Copy()
	defer session.Close()

	c := db.C(UserCollectionName).With(session)
	c.Find(bson.M{"_id": bson.ObjectIdHex(id)}).One(u)
}

func (u *UserAPIResponse) GetByEmail(db *mgo.Database, email string) {
	session := db.Session.Copy()
	defer session.Close()

	c := db.C(UserCollectionName).With(session)
	c.Find(bson.M{"email": email}).One(u)
}

func (u *UserAPIResponse) GetByIdStr(db *mgo.Database, id string) {
	session := db.Session.Copy()
	defer session.Close()

	c := db.C(UserCollectionName).With(session)
	c.Find(bson.M{"_id": bson.ObjectIdHex(id)}).One(u)
}
