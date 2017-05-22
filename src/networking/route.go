package networking

import (
	"api"
	"models"
	"tools"
)

var Handles = map[string]interface{}{
	"/h2":                                             Info,
	"/docs":                                           Info,
	"/categories":                                     api.Categories,
	"/user/create":                                    api.UserCreate,
	"/user/confirm/email/:user_id/:confirmation_code": api.UserConfirmation,
	"/store/instance/:store_id":                       api.StoreSearch,
	"/store/create":                                   api.StoreCreate,
	"/store/update/:store_id":                         api.StoreCreate,
	"/store/categories/:store_id":                     api.StoreCategories,
	"/transform":                                      tools.RemodelJ,
}

var APIRouteMap = map[string]map[string]interface{}{
	"/h2":         {"control_method": "GET"},
	"/docs":       {"control_method": "GET"},
	"/categories": {"control_method": "GET"},

	// products api

	// users api
	"/user/create": {
		"control_method": "POST",
		"post_payload":   models.User{},
	},
	"/user/confirm/email/:user_id/:confirmation_code": {
		"control_method": "GET",
	},

	// stores api
	"/store/instance/:store_id": {
		"control_method": "GET",
	},
	"/store/create": {
		"control_method": "POST",
		"post_payload":   models.Store{},
	},
	"/store/update/:store_id": {
		"control_method": "POST",
		"post_payload":   models.Store{},
	},
	"/store/categories/:store_id": {
		"control_method": "GET",
	},

	"/transform": {
		"control_method": "GET",
		"post_payload":   models.User{},
	},
}
