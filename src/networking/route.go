package networking

import (
	"api"
	"models"
	"tools"
)

// max requests per sec as defined globally for calls not needing fine grained
// control (improve throughput by avoiding cpu overclocking and other shit)
//
// 		it can probably be higher given current benchmarks but
//		will eventually be dynmic for all calls and no static
//		values will need to be set. but until that fun lets
//		just keep it simple and controlled with a safe number
const MAX_RPS = 10000

var Handles = map[string]interface{}{
	"/h2":   Info,
	"/docs": Docs,
	//"/categories":                                     api.Categories,
	"/user/create":                                    api.UserCreate,
	"/user/login":                                     api.Login,
	"/user/retrieve":                                  api.GetUserById,
	"/user/confirm/email/:user_id/:confirmation_code": api.UserConfirmation,
	"/store/retrieve/full/:store_id":                  api.GetStoreById,
	"/store/create":                                   api.StoreCreate,
	"/store/categories/create":                        api.StoreCategories,
	"/transform":                                      tools.RemodelJ,
}

var APIRouteMap = map[string]map[string]interface{}{
	"/h2":   {"control_method": "GET", "authenticate": []string{}, "max_rps": 100},
	"/docs": {"control_method": "GET", "authenticate": []string{}, "max_rps": 100},
	//"/categories": {"control_method": "GET"},

	// products api
	"/store/categories/create": {
		"control_method": "POST",
		"post_payload":   []models.StoreCategory{},
		"authenticate":   []string{},
		"max_rps":        nil,
	},

	// users api
	"/user/create": {
		"control_method": "POST",
		"post_payload":   models.User{},
		"authenticate":   []string{},
		"max_rps":        nil,
	},
	"/user/retrieve": {
		"control_method": "GET",
		"authenticate": []string{
			models.ACCESSROLE_CONFIRMED_USER,
		},
		"max_rps": nil,
	},
	"/user/login": {
		"control_method": "GET",
		"authenticate":   []string{},
		"max_rps":        nil,
	},
	"/user/confirm/email/:user_id/:confirmation_code": {
		"control_method": "GET",
		"authenticate":   []string{},
		"max_rps":        nil,
	},

	// stores api
	"/store/retrieve/full/:store_id": {
		"control_method": "GET",
		"authenticate": []string{
			models.ACCESSROLE_CONFIRMED_USER,
			models.ACCESSROLE_STOREOWNER,
		},
		"max_rps": nil,
	},
	"/store/create": {
		"control_method": "POST",
		"post_payload":   models.Store{},
		"authenticate": []string{
			models.ACCESSROLE_CONFIRMED_USER,
		},
		"max_rps": nil,
	},

	"/transform": {
		"control_method": "GET",
		"post_payload":   models.User{},
		"authenticate":   []string{},
		"max_rps":        nil,
	},
}
