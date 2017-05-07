package api

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func FindStores(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, string(`{
	"stores": [{
		"id": "asldk2l3kj4h5",
		"lat": 13.51,
		"lon": -13.51,
		"name": "sample store",
		"image": "https://...",
		"address": "",
		"distance": "x.x miles",
		"phone_number": "(xxx) xxx - xxxx",
		"long_description": "long store description...",
		"short_description": "short store description...",
		"working_hours": [
			{
				"day": "Sunday",
				"from": "09:00",
				"to": "17:00"
			},
			{
				"day": "Monday",
				"from": "09:00",
				"to": "17:00"
			},
			{
				"day": "Tuesday",
				"from": "09:00",
				"to": "17:00"
			},
			{
				"day": "Wednesday",
				"from": "09:00",
				"to": "17:00"
			},
			{
				"day": "Thursday",
				"from": "09:00",
				"to": "17:00"
			},
			{
				"day": "Friday",
				"from": "09:00",
				"to": "17:00"
			},
			{
				"day": "Saturday",
				"from": "09:00",
				"to": "17:00"
			}
		],
		"pickup": {
			"service_offered": false,
			"minimum_time_to_pickup": 30,
			"maximum_time_to_pickup": 90
		},
		"delivery": {
			"service_offered": false,
			"delivery_fee": "$0.00",
			"delivery_minimum": "$10.00",
			"delivery_distance": 2,
			"minimum_time_to_delivery": 1800
		}
	}]
}`))
}
