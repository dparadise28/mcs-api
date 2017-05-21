package models

import "gopkg.in/mgo.v2/bson"

var UserCollectionName = "Users"

type UserAddressBook struct {
	AddressName string  `bson:"address_name" json:"address_name"`
	Location    GeoJson `bson:"location" json:"location"`
}

type UserRoles struct {
	// list of store ids stored against their respective roles
	// (might want to include names and meta but ok for now)
	StoresOwned     []bson.ObjectId `bson:"stores_owned" json:"stores_owned"`
	StoresEmployeIn []bson.ObjectId `bson:"stores_employed_in" json:"stores_employed_in"`
}

type User struct {
	ID               bson.ObjectId   `bson:"_id,omitempty" json:"user_id"`
	ConfirmationCode string          `bson:"confirmation_code" json:"confirmation_code"`
	AddressBook      UserAddressBook `bson:"address_book" json:"address_book"`
	Confirmed        bool            `bson:"confirmed" json:"confirmed"`
	Password         string          `bson:"password" json:"password"`
	UserName         string          `bson:"username" json:"username"`
	Email            string          `bson:"email" json:"email"`
	Roles            UserRoles       `bson:"user_roles" json:"user_roles"`
}
