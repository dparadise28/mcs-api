package models

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var StoreCollectionName = "Stores"

type OpenHours struct {
	Hours  Hours `json:"hours"`
	IsOpen bool  `json:"open"`
}

type WeeklyWorkingHours struct {
	Sun OpenHours `bson:"sun" json:"sunday"`
	Mon OpenHours `bson:"mon" json:"monday"`
	Tue OpenHours `bson:"tue" json:"tuesday"`
	Wed OpenHours `bson:"wed" json:"wednesday"`
	Thu OpenHours `bson:"thu" json:"thursday"`
	Fri OpenHours `bson:"fri" json:"friday"`
	Sat OpenHours `bson:"sat" json:"saturday"`
}

type StoreDelivery struct {
	// we might want to offer a variable delivery
	// fee model at some point but this will do
	// for now. That can be split out form the base
	// store model
	Fee       uint   `bson:"fee,omitempty" json:"delivery_fee"`
	Offered   bool   `bson:"offered" json:"service_offered" validate:"required"`
	MaxDist   int    `bson:"max_distance,omitempty" json:"delivery_distance"`
	MinTime   uint8  `bson:"min_time,omitempty" json:"maximum_time_to_delivery"`
	MaxTime   uint8  `bson:"max_time,omitempty" json:"minimum_time_to_delivery"`
	MinAmount uint16 `bson:"min_amount,omitempty" json:"delivery_minimum"`
}

type StorePickup struct {
	Offered         bool  `json:"offered" validate:"required"`
	MinTime         uint8 `bson:"min_time" json:"minimum_time_to_pickup" validate:"max=90,min=1"`
	MaxTime         uint8 `bson:"max_time" json:"maximum_time_to_pickup" validate:"max=90,min=1"`
	PickupItemCount struct {
		Min uint8  `json:"min" validate:"max=255,min=1"`
		Max uint32 `json:"max" validate:"max=4294967295,min=1"`
	} `bson:"pickup_items" json:"pickup_items"`
}

type Store struct {
	ID              bson.ObjectId      `bson:"_id,omitempty" json:"store_id"`
	Name            string             `bson:"name" json:"name"`
	Image           string             `json:"image"`
	Phone           string             `json:"phone"`
	Pickup          StorePickup        `json:"pickup"`
	Address         Address            `json:"address"`
	TaxRate         float64            `json:"tax_rate" validate:"required"`
	Delivery        StoreDelivery      `json:"delivery"`
	Distance        float64            `bson:"distance,omitempty" json:"distance,omitempty"`
	WorkingHours    WeeklyWorkingHours `json:"working_hours"`
	LongDescription string             `json:"long_description"`
	// this field has has a fulltext index for
	// full text search so must so we must
	// ensure its length for now to avoid index
	// bloating untill switching to a more robust
	// search solution or building one
	ShortDescription string `json:"short_description" validate:"max=50"`

	// retain the list of platform and store
	// instance categories for order and filtering
	// needs
	CategoryList       []string `json:"store_categories"`
	PlatformCategories []string `json:"platform_categories"`
}

func (s *Store) Insert(db *mgo.Database) error {
	// copy db session for the stores collection and close on completion
	session := db.Session.Copy()
	defer session.Close()

	// grab the proper collection, create a new store id and attempt an insert
	c := db.C(StoreCollectionName).With(session)
	s.ID = bson.NewObjectId()
	s.Address.Location.Type = "Point"
	insert_err := c.Insert(&s)
	return insert_err
}
