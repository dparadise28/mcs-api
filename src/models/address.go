package models

import (
	"errors"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var AddressCollectionName = "addresses"

type Address struct {
	ID            bson.ObjectId `bson:"_id" json:"address_id"`
	Name          string        `bson:"name" json:"name"`
	City          string        `bson:"city" json:"city" validate:"required"`
	Line1         string        `bson:"line1" json:"line1" validate:"required"`
	Route         string        `bson:"route" json:"route"`
	UserID        bson.ObjectId `bson:"user_id" json:"user_id"`
	Default       bool          `bson:"default" json:"default"`
	Country       string        `bson:"country" json:"country" validate:"required"`
	Location      GeoJson       `bson:"location" json:"location"`
	Latitude      float64       `bson:"latitude" json:"latitude" validate:"required,min=-85.0511501,max=85.001"`
	Longitude     float64       `bson:"longitude" json:"longitude" validate:"required,min=-180.001,max=180.001"`
	PostalCode    string        `bson:"postal_code" json:"postal_code" validate:"required"`
	StreetNumber  string        `bson:"street_number" json:"street_number"`
	AdminAreaLvl1 string        `bson:"administrative_area_level_1" json:"administrative_area_level_1"`

	// helper fields
	NewDefaultID bson.ObjectId `bson:"-"`

	DB        *mgo.Database `bson:"-" json:"-"`
	DBSession *mgo.Session  `bson:"-" json:"-"`
}

func (a *Address) AddToUserAddressBook() error {
	c := a.DB.C(AddressCollectionName).With(a.DBSession)
	a.Location.Type = "Point"
	a.Location.Coordinates = []float64{a.Longitude, a.Latitude}
	a.ID = bson.NewObjectId()

	defaultCount, countErr := c.Find(bson.M{
		"user_id": a.UserID,
		"default": true,
	}).Count()
	if countErr != nil {
		return countErr
	}
	if a.Default && defaultCount > 0 {
		_, err := c.UpdateAll(bson.M{
			"user_id": a.UserID,
			"default": true,
		}, bson.M{
			"default": false,
		})
		if err != nil {
			return err
		}
	}
	if defaultCount == 0 {
		// if no addresses exist (default or otherwise) then set the
		// incoming addr as the default by default (mouthful)
		a.Default = true
	}
	return c.Insert(&a)
}

func (a *Address) RemoveFromUserAddressBook() error {
	c := a.DB.C(AddressCollectionName).With(a.DBSession)
	err := c.Remove(bson.M{
		"_id":     a.ID,
		"user_id": a.UserID,
		"default": false,
	})
	if err != nil {
		err = errors.New(
			"We could not remove the address selected at this time. " +
				"If you are trying to remove your current default " +
				"address please change your default first before " +
				"removing this address.",
		)
	}
	return err
}

func (a *Address) ChangeDefaultAddressInUserAddressBook() error {
	c := a.DB.C(AddressCollectionName).With(a.DBSession)
	_, err := c.UpdateAll(bson.M{
		"user_id": a.UserID,
		"default": true,
	}, bson.M{
		"$set": bson.M{
			"default": false,
		},
	})
	if err != nil {
		return err
	}
	change := mgo.Change{
		ReturnNew: true,
		Upsert:    false,
		Remove:    false,
		Update: bson.M{
			"$set": bson.M{
				"default": true,
			},
		},
	}
	_, err = c.Find(bson.M{
		"user_id": a.UserID,
		"_id":     a.NewDefaultID,
	}).Apply(change, a)
	return err
}

func (a *Address) RetrieveUserDefaultAddressInAddressBook() error {
	c := a.DB.C(AddressCollectionName).With(a.DBSession)
	err := c.Find(bson.M{
		"user_id": a.UserID,
		"default": true,
	}).One(a)
	return err
}

func (a *Address) RetrieveUserAddressBook() ([]Address, error) {
	c := a.DB.C(AddressCollectionName).With(a.DBSession)
	addrs := []Address{}
	err := c.Find(bson.M{
		"user_id": a.UserID,
	}).All(&addrs)
	return addrs, err
}
