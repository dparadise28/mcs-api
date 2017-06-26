package models

import (
	"log"
	"reflect"
)

type Email struct {
	To      string
	Body    string
	Subject string
}

type GeoJson struct {
	Type        string    `json:"-"`
	Coordinates []float64 `json:"coordinates"`
}

type Address struct {
	Route         string  `bson:"route" json:"route"`
	Country       string  `bson:"country" json:"country"`
	Location      GeoJson `bson:"location" json:"location"`
	Latitude      float64 `bson:"latitude" json:"latitude" validate:"required"`
	Longitude     float64 `bson:"longitude" json:"longitude" validate:"required"`
	PostalCode    string  `bson:"postal_code" json:"postal_code"`
	StreetNumber  string  `bson:"street_number" json:"street_number"`
	AdminAreaLvl1 string  `bson:"administrative_area_level_1" json:"administrative_area_level_1"`
}

type Hours struct {
	From uint16 `json:"from" validate:"ltefield=To"`
	To   uint16 `json:"to"`
}

func I(array interface{}) []interface{} {
	// any array type to array interface (useful for mongo multi insert/retrieve)
	v := reflect.ValueOf(array)
	t := v.Type()

	if t.Kind() != reflect.Slice {
		log.Panicf("`array` should be %s but got %s", reflect.Slice, t.Kind())
	}

	result := make([]interface{}, v.Len(), v.Len())

	for i := 0; i < v.Len(); i++ {
		result[i] = v.Index(i).Interface()
	}

	return result
}
