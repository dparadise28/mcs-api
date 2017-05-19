package networking

import (
	"api"
	"tools"
)

var APIRouteMap = map[string]map[string]interface{}{
	"/h2": {
		"control_method": "GET",
		"handler_method": Info,
	},
	"/categories": {
		"control_method": "GET",
		"handler_method": api.Categories,
	},

	// products api

	// users api

	// stores api
	"/store/instance/:store_id": {
		"control_method": "GET",
		"handler_method": api.StoreSearch,
	},
	"/store/create": {
		"control_method": "POST",
		"handler_method": api.StoreCreate,
	},
	"/store/update/:store_id": {
		"control_method": "POST",
		"handler_method": api.StoreCreate,
	},
	"/store/categories/:store_id": {
		"control_method": "GET",
		"handler_method": api.StoreCategories,
	},

	"/transform": {
		"control_method": "GET",
		"handler_method": tools.RemodelJ,
	},
}
