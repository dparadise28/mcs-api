package models

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
)

var UserCollectionName = "Users"

type User struct {
	ID               bson.ObjectId     `bson:"_id,omitempty" json:"user_id"`
	ConfirmationCode string            `bson:"confirmation_code" json:"confirmation_code"`
	IsStoreOwner     bool              `bson:"is_store_owner" json:"is_store_owner"`
	Confirmed        bool              `bson:"confirmed" json:"confirmed"`
	Password         string            `bson:"password" json:"password" validate:"required"`
	Phone            string            `bson:"phone" json:"phone"` // validate:"required"`
	Email            string            `bson:"email" json:"email" validate:"required,email"`
	Roles            UserRoles         `bson:"user_roles" json:"user_roles"`
	Stores           map[string]string `bson:"store_map" json:"store_map"`
	StripeCustomerID string            `bson:"stripe_customer_id" json:"stripe_customer_id"`
	// user roles struct is in /src/models/auth

	Login struct {
		Token string        `bson:"-" json:"authtoken"`
		UID   bson.ObjectId `bson:"-" json:"userID"`
	} `bson:"-" json:"login"`

	// helper fields
	DB          *mgo.Database `bson:"-" json:"-"`
	DBSession   *mgo.Session  `bson:"-" json:"-"`
	StripeToken string        `bson:"-" json:"-"`
}

type UserAPIResponse struct {
	ID           bson.ObjectId     `bson:"_id,omitempty" json:"user_id"`
	IsStoreOwner bool              `bson:"is_store_owner" json:"is_store_owner" validate:"required"`
	Confirmed    bool              `bson:"confirmed" json:"confirmed"`
	Email        string            `bson:"email" json:"email" validate:"required,email"`
	Roles        UserRolesResponse `bson:"user_roles" json:"user_roles"`
}

/*
func (u *User) DoesUserHaveAStripAccount() bool {
	c := u.DB.C(UserCollectionName).With(u.DBSession)
	if count == 1; count := c.Find(bson.M{
		"_id": u.ID,
		"stripe_customer_id": bson.M{
			"$exists": true,
		},
	}).Count() {}
	return err
}
*/
func (u *User) AddUserStripeCustomerAccount() error {
	c := u.DB.C(UserCollectionName).With(u.DBSession)
	log.Println(u.StripeCustomerID)
	change := mgo.Change{
		ReturnNew: true,
		Upsert:    false,
		Remove:    false,
		Update: bson.M{
			"$set": bson.M{
				"stripe_customer_id": u.StripeCustomerID,
			},
		},
	}
	log.Println(u.StripeCustomerID)
	_, err := c.Find(bson.M{
		"_id": u.ID,
	}).Apply(change, u)
	return err
}

func (u *User) ConfirmUserEmailLink(pw_reset bool) string {
	link := DOMAIN + "/api/user/confirm/email/" + u.ID.Hex() + "/" + u.ConfirmationCode
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

func (u *User) GetById() error {
	c := u.DB.C(UserCollectionName).With(u.DBSession)
	err := c.Find(bson.M{
		"_id": u.ID,
	}).One(u)
	return err
}
