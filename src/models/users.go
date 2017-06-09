package models

import "gopkg.in/mgo.v2/bson"

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

func (u *User) GenerateToken() string {
	return "token"
}
