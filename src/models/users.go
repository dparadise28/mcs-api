package models

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var UserCollectionName = "Users"

type User struct {
	ID               bson.ObjectId `bson:"_id,omitempty" json:"user_id"`
	ConfirmationCode string        `bson:"confirmation_code" json:"confirmation_code"`
	AddressBook      []Address     `bson:"address_book" json:"address_book"`
	DefaultAddr      bson.ObjectId `bson:"default_address" json:"default_address"`
	Confirmed        bool          `bson:"confirmed" json:"confirmed"`
	Password         string        `bson:"password" json:"password" validate:"required"`
	Email            string        `bson:"email" json:"email" validate:"required,email"`
	Roles            UserRoles     `bson:"user_roles" json:"user_roles"`
	// user roles struct is in /src/models/auth

	Login struct {
		Token string        `bson:"-" json:"authtoken"`
		UID   bson.ObjectId `bson:"-" json:"userID"`
	} `bson:"-" json:"login"`
}

type UserAPIResponse struct {
	ID bson.ObjectId `bson:"_id,omitempty" json:"user_id"`
	//AddressBook []UserAddress `bson:"address_book" json:"address_book"`
	Confirmed bool      `bson:"confirmed" json:"confirmed"`
	Email     string    `bson:"email" json:"email" validate:"required,email"`
	Roles     UserRoles `bson:"user_roles" json:"user_roles"`
}

func (a *Address) AddAddressToUserAddressBook(u_id bson.ObjectId) error {
	c := a.DB.C(UserCollectionName).With(a.DBSession)
	a.Location.Coordinates = []float64{a.Latitude, a.Latitude}
	pushQuery := bson.M{
		"$push": bson.M{
			"address_book": a,
		},
	}
	change := mgo.Change{
		ReturnNew: true,
		Upsert:    false,
		Remove:    false,
		Update:    pushQuery,
	}
	_, err := c.Find(bson.M{
		"_id": u_id,
	}).Apply(change, a)
	return err
}

func (u *User) ConfirmUserEmailLink(pw_reset bool) string {
	link := "http://mycorner.store:8080/api/user/confirm/email/" + u.ID.Hex() + "/" + u.ConfirmationCode
	if pw_reset {
		link += "?password=" + u.Password
	}
	return link
}

func (u *User) EmailConfirmation(pw_reset bool) Email {
	// Email struct model can be found in the general.go file
	emailSubject := "Thank You for signing up!"
	if pw_reset {
		emailSubject = "Forgot your password ehy?"
	}
	emailBody := ConirmationEmail(
		u.ID.Hex(),
		u.ConfirmationCode,
		u.Email,
		u.ConfirmUserEmailLink(pw_reset),
	)
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
