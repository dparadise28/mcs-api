package models

import "gopkg.in/mgo.v2/bson"

var StoreCollectionName = "Stores"

type Hours struct {
	From uint16 `json:"from"`
	To   uint16 `json:"to"`
}

type OpenHours struct {
	// split into am and pm to allow for efficient
	// index and search in the case of non inclusive
	// ranges cases:
	// 		9am-5pm
	//      6pm-2am
	// TODO:time zone considerations
	AM     Hours `json:"am"`
	PM     Hours `json:"pm"`
	IsOpen bool  `json:"open"`
}

type WeeklyWorkingHours struct {
	Sun OpenHours `json:"sun"`
	Mon OpenHours `json:"mon"`
	Tue OpenHours `json:"tue"`
	Wed OpenHours `json:"wed"`
	Thu OpenHours `json:"thu"`
	Fri OpenHours `json:"fri"`
	Sat OpenHours `json:"sat"`
}

type GeoJson struct {
	Type        string    `json:"-"`
	Coordinates []float64 `json:"coordinates"`
}

type StoreDelivery struct {
	// we might want to offer a variable delivery
	// fee model at some point but this will do
	// for now. That can be split out form the base
	// store model
	Fee       uint   `bson:"delivery_fee,omitempty" json:"delivery_fee"`
	Offered   bool   `json:"offered" validate:"required"`
	MaxDist   int    `bson:"max_distance,omitempty" json:"max_distance"`
	MinTime   uint8  `bson:"min_time,omitempty" json:"min_time" validate:"max=180,min=0"`
	MaxTime   uint8  `bson:"max_time,omitempty" json:"max_time" validate:"max=180,min=0"`
	MinAmount uint16 `bson:"min_amount,omitempty" json:"min_amount"`
}

type StorePickup struct {
	Offered bool  `json:"offered" validate:"required"`
	MinTime uint8 `json:"min_time" validate:"max=90,min=0"`
	MaxTime uint8 `json:"max_time" validate:"max=90,min=0"`
}

type Store struct {
	ID              bson.ObjectId      `bson:"_id,omitempty" json:"store_id"`
	Name            string             `bson:"name" json:"name"`
	Image           string             `json:"image"`
	Phone           string             `json:"phone"`
	Pickup          StorePickup        `json:"pickup"`
	Delivery        StoreDelivery      `json:"delivery"`
	Distance        float64            `bson:"distance,omitempty" json:"distance,omitempty"`
	Address1        string             `json:"address_line_1"`
	Address2        string             `json:"address_line_2"`
	Location        GeoJson            `bson:"location" json:"location"`
	WorkingHours    WeeklyWorkingHours `json:"working_hours"`
	LongDescription string             `json:"long_description"`
	// this field has has a fulltext index for
	// full text search so must so we must
	// ensure its length for now to avoid index
	// bloating untill switching to a more robust
	// search solution or building one
	ShortDescription string `json:"short_description"`

	// retain the list of platform and store
	// instance categories for order and filtering
	// needs
	CategoryList       []string `json:"store_categories"`
	PlatformCategories []string `json:"platform_categories"`
}
