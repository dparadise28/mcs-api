package networking

import (
	"api"
	"models"
	"tools"
)

var Handles = map[string]interface{}{
	//"/h2":                                             Info,
	"/docs": Docs,
	//"/categories":                                     api.Categories,
	"/user/create":                                    api.UserCreate,
	"/user/login":                                     api.Login,
	"/user/retrieve":                                  api.GetUserById,
	"/user/confirm/email/:user_id/:confirmation_code": api.UserConfirmation,
	"/store/retrieve/full/:store_id":                  api.GetStoreById,
	"/store/create":                                   api.StoreCreate,
	"/transform":                                      tools.RemodelJ,
}

var APIRouteMap = map[string]map[string]interface{}{
	//"/h2":         {"control_method": "GET"},
	"/docs": {"control_method": "GET", "authenticate": []string{}},
	//"/categories": {"control_method": "GET"},

	// products api

	// users api
	"/user/create": {
		"control_method": "POST",
		"post_payload":   models.User{},
		"authenticate":   []string{},
	},
	"/user/retrieve": {
		"control_method": "GET",
		"authenticate": []string{
			models.ACCESSROLE_CONFIRMED_USER,
		},
	},
	"/user/login": {
		"control_method": "GET",
		"authenticate":   []string{},
	},
	"/user/confirm/email/:user_id/:confirmation_code": {
		"control_method": "GET",
		"authenticate":   []string{},
	},

	// stores api
	"/store/retrieve/full/:store_id": {
		"control_method": "GET",
		"authenticate": []string{
			models.ACCESSROLE_CONFIRMED_USER,
			models.ACCESSROLE_STOREOWNER,
		},
	},
	"/store/create": {
		"control_method": "POST",
		"post_payload":   models.Store{},
		"authenticate": []string{
			models.ACCESSROLE_CONFIRMED_USER,
		},
	},

	"/transform": {
		"control_method": "GET",
		"post_payload":   models.User{},
		"authenticate":   []string{},
	},
}
