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

	"/stores": {
		"control_method": "GET",
		"handler_method": api.FindStores,
	},
	"/store/info": {
		"control_method": "GET",
		"handler_method": api.StoreInfo,
	},
	"/store/categories": {
		"control_method": "GET",
		"handler_method": api.StoreCategories,
	},

	"/transform": {
		"control_method": "GET",
		"handler_method": tools.RemodelJ,
	},
}
