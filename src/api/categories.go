package api

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func Categories(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, string(`{
	"categories": [
		{"name": "c1", "id": "tier_1_c_1"},
		{"name": "c2", "id": "tier_1_c_2"},
		{"name": "c3", "id": "tier_1_c_3"},
		{"name": "c4", "id": "tier_1_c_4"},
		{"name": "c5", "id": "tier_1_c_5"},
		{"name": "c6", "id": "tier_1_c_6"},
		{"name": "c7", "id": "tier_1_c_7"},
		{"name": "c8", "id": "tier_1_c_8"},
		{"name": "c9", "id": "tier_1_c_9"}
	]
}`))
}
