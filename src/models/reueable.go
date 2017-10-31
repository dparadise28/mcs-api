package models

import (
	"log"
	"reflect"
)

var (
	DOMAIN      = ""
	UI_DIR_PATH = ""
)

type Email struct {
	To      string
	Body    string
	Subject string
}

type EmailReq struct {
	To      string `json:"to" validate:"required"`
	Body    string `json:"body" validate:"required"`
	Subject string `json:"subject" validate:"required"`
}

type GeoJson struct {
	Type        string    `json:"-"`
	Coordinates []float64 `json:"coordinates"`
}

type Hours struct {
	From uint16 `json:"from" validate:"ltecsfield=To"`
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

func Keys(m map[string]interface{}) (keys []string) {
	for key := range m {
		keys = append(keys, key)
	}
	return
}
