package api

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func StoreInfo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, string(`{
	"id": "o8ifu23o8finlk",
	"lat": 13.51,
	"lon": -13.51,
	"image": "https://...",
	"phone": "(xxx) xxx - xxxx",
	"distance": "1.5 miles",
	"store_name": "Sample Store",
	"short_description": "short desc about store instance."
}`))
}

func StoreCategories(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, string(`{
	"categories": [{
		"category_name": "c1",
		"id": "c_1",
		"products": [{
			"name": "p1",
			"id": "fa87asoiudf9",
			"description": "short prod description",
			"display_price": "$xx.xx",
			"image": "https://..."
		}]
	}]
}`))
}
