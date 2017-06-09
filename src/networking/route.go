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
	"/user/login/:email":                              api.Login,
	"/user/retrieve/:user_id":                         api.GetUserById,
	"/user/confirm/email/:user_id/:confirmation_code": api.UserConfirmation,
	"/store/retrieve/:store_id":                       api.GetStoreById,
	"/store/create":                                   api.StoreCreate,
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
	"/user/retrieve/:user_id": {
		"control_method": "GET",
	},
	"/user/login/:email": {
		"control_method": "GET",
	},
	"/user/confirm/email/:user_id/:confirmation_code": {
		"control_method": "GET",
	},

	// stores api
	"/store/retrieve/:store_id": {
		"control_method": "GET",
	},
	"/store/create": {
		"control_method": "POST",
		"post_payload":   models.Store{},
	},

	"/transform": {
		"control_method": "GET",
		"post_payload":   models.User{},
	},
}
