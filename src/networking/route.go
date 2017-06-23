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

	// user
	"/user/create":                                    api.UserCreate,
	"/user/login":                                     api.Login,
	"/user/retrieve":                                  api.GetUserById,
	"/user/confirmation/resend":                       api.UserResendConfirmation,
	"/user/confirm/email/:user_id/:confirmation_code": api.UserConfirmation,

	// store
	"/store/create":                  api.StoreCreate,
	"/store/search":                  api.StoreSearch,
	"/store/info/:store_id":          api.GetStoreById,
	"/store/retrieve/full/:store_id": api.GetFullStoreById,

	// store categories
	"/store/categories/:store_id":        api.GetCategoriesByStoreId,
	"/store/category/create/:store_id":   api.AddStoreCategory,
	"/store/category/udpate/:store_id":   api.UpdateStoreCategory,
	"/store/category/activate/:store_id": api.EnableStoreCategory,

	// store products
	//"/store/product/update/:store_id": api.EnableStoreCategory,
	//"/store/product/activate/:store_id": api.EnableStoreCategory,

	// cart
	"/cart/update/product/quantity": api.UpdateCartProductQuantity,
	"/transform":                    tools.RemodelJ,
}

var APIRouteMap = map[string]map[string]interface{}{
	"/h2": {
		"control_method": "GET",
		"authenticate":   []string{},
		"max_rps":        100,
		"description": []string{
			"Heartbeat",
		},
	},
	"/docs": {
		"control_method": "GET",
		"authenticate":   []string{},
		"max_rps":        100,
		"description": []string{
			"Your lookin at it",
			"",
			"Things to know or questions that may arise",
			"",
			"When the store id is required in both the url schame",
			"as well as the post payload that is a consequence of the",
			"way authentication works. This needs to be ironed out but",
			"is a solid enough approach for now and will be easy to",
			"update/enhance later. To elaborate, when checking whether",
			"a user is authorized to use the service being requested the",
			"token provided is decrypted; if a token is not provided and",
			"the authenticate array found in these docs is populated with",
			"permissions then the user will be refused access imidiately.",
			"The contents of the token include a map of the permissions",
			"where the permissions are keyed on by the store id",
			"(<store_id>: 'perm1,perm2,perm3,...,perm4').",
			"",
			"To avoid additional overhead, complexity and latency the",
			"the contents of the request body are not deserialized until",
			"until the user is verrified. Since the size of the request",
			"(generally speaking) is smaller then that of the payload it",
			"it is faster/lighter to retrieve the value from the url then",
			"from the payload. So the general authentication flow is as",
			"follows:",
			"",
			"Check whether the authenticate key in this structure for a",
			"a given call is restricted to certain permissions.",
			"",
			"if that is the case (always check the validity of token and desirialize)",
			"	for the special case of user confirmed, simply check that",
			"	the value set for confirmed in the tokens content is set to true",
			"",
			"	if more permissions are specified in the auth list",
			"		1) grab the store id in the url",
			"		2) grab the value keyed by the store_id in the user perms",
			"		   map found when deserializing the token",
			"		3) itterate through the permissions in the auth array here",
			"		   and check if they are contained in the value retrieved in",
			"		   step 2.",
		},
	},

	// users api
	"/user/create": {
		"control_method": "POST",
		"post_payload":   models.User{},
		"authenticate":   []string{},
		"max_rps":        nil,
		"description": []string{
			"This is a call for new users to create an account on our platform.",
			"The required fileds in the post payload are username, email, and",
			"password. Username may be removed depending on the need/desire for",
			"it. It is not actually used anywhere and there is no index set on it.",
			"It may be useful for future admin use but for now its just another",
			"field passed when creating a user",
			"",
			"",
			"-When creating a user uniqueness is determined by the email only.",
			"",
			"-Passwords are not stored in a human readable format since it will.",
			"",
			"-When creating a user you are sent a verrification email that is",
			"that is required to confirm/activate a user for anything requiring",
			"that permission",
			"",
			"-The token sent on successful completion of this request will always",
			"contain a user id and token. The token is for an unconfirmed user. To",
			"confirm the user hit the url specified in the email, and login to your",
			"account",
			"",
			"",
			"More on passwords:",
			"Passwords as said before are not stored in a human readable format.",
			"What is stored is a hash of the original password. It is irreversible.",
			"That means that it is a one way transformation. As in all hashing algorithms,",
			"an input can and should produce a variable output. Which means that when",
			"validating an input you will only be able to check that the input is valid",
			"given a prior hash of the input. For more info you can look into one way hash",
			"functions. We use bcrypt.",
			"",
			"If a user forgets his/her passwords there is a /user/confirmation/resend call.",
			"More information on how this mechanism works can be found in the description",
			"of the call itself.",
		},
	},
	"/user/confirm/email/:user_id/:confirmation_code": {
		"control_method": "GET",
		"authenticate":   []string{},
		"max_rps":        nil,
		"description": []string{
			"This is the url that is formed in the confirmation email and sent to the user",
			"on user creation (/user/create). Assuming the correct user id and confirmation",
			"code (these valuse are set during user creation and checked during this call)",
			"the confirmation flag in the user record is flipped to true and the confirmation",
			"code is reset to an empty string to avoid potentially unlikely attempts to reset",
			"passwords due to unsafely kept confirmation codes or brute force attacks.",
			"",
			"If a user is confirmed and the confirmation code exists (assuming the one provided",
			"matches what is associated with the user in the db record) you may pass an optional",
			"query param to reset the passord. The structure would look like this:",
			"",
			"protocol://base_url:port/user/confirm/email/:user_id/:confirmation_code?password={new_pw}",
		},
	},
	"/user/confirmation/resend": {
		"control_method": "GET",
		"authenticate":   []string{},
		"max_rps":        nil,
		"description": []string{
			"call structure: protocol://base_url:port/api/user/confirmation/resend&email={user_email}",
			"the port will likely not be required but left in as an example if the env used to test",
			"is binded to that address/port",
			"",
			"When a user forgets his/her password this is the call that should be envoked to",
			"resend a confirmation code. The call WILL NOT rest the user password. What it does",
			"do is reset the confirmation code associated with the user to allow them to reset",
			"the password. In the email you should find a confirmation code and a user id which",
			"will be required on the ui when resetting the password. The reset happens when",
			"confirming the user. For more information on this please refer to:",
			"/user/confirm/email/:user_id/:confirmation_code",
		},
	},
	"/user/login": {
		"control_method": "GET",
		"authenticate":   []string{},
		"max_rps":        nil,
		"description": []string{
			"If your a user and the email/pw combo is correct after all the proper validation then",
			"you will recieve a response that contains an authtoken and a userID which you can throw",
			"in your headers and go about your business.",
			"",
			"call structure:",
			"",
			"protocol://base_url/user/login?email={email}&password={password}",
		},
	},
	"/user/retrieve": {
		"control_method": "GET",
		"authenticate": []string{
			models.ACCESSROLE_CONFIRMED_USER,
		},
		"max_rps": nil,
		"description": []string{
			"Simply a call to retrieve the user info. Must be confirmed to be allowed",
			"access to this call.",
		},
	},

	// store categories
	"/store/categories/:store_id": {
		"control_method": "GET",
		"authenticate":   []string{},
		"max_rps":        nil,
		"desctiption": []string{
			"Returns the full set of hydrated store categories containing",
			"products in the order set by a store. This will filter out",
			"any disabled categories.",
			"",
			"To include all disabled store categories and or productsuse the query",
			"params include_disabled_categories=true&include_disabled_products=true",
			"respectively.",
			"",
			"If there is a product that belongs to a diabled category and",
			"the include_disabled_categories is set to true then the product",
			"not show up in the response. To find that product or all products",
			"use the /store/products/:store_id call with the include_disabled_products",
			"query param",
		},
	},
	"/store/category/create/:store_id": {
		"control_method": "POST",
		"authenticate": []string{
			models.ACCESSROLE_CONFIRMED_USER,
			models.ACCESSROLE_STOREOWNER,
		},
		"post_payload": models.Category{},
		"max_rps":      nil,
		"description": []string{
			"Creates a new category and appends it as the",
			"last in the seq. It is automatically set to be",
			"enabled and will show up for users.",
		},
	},
	"/store/category/udpate/:store_id": {
		"control_method": "POST",
		"authenticate": []string{
			models.ACCESSROLE_CONFIRMED_USER,
			models.ACCESSROLE_STOREOWNER,
		},
		"post_payload": models.Category{},
		"max_rps":      nil,
		"desctiption": []string{
			"This is an interface exposed to update the category name.",
		},
	},
	"/store/category/activate/:store_id": {
		"control_method": "POST",
		"authenticate": []string{
			models.ACCESSROLE_CONFIRMED_USER,
			models.ACCESSROLE_STOREOWNER,
		},
		"post_payload": models.Category{},
		"max_rps":      nil,
		"description": []string{
			"Toggles the store category enabled field based on the value",
			"it is set to in the request. This will affect user searches.",
		},
	},

	// stores api
	"/store/create": {
		"control_method": "POST",
		"post_payload":   models.Store{},
		"authenticate":   []string{
		//models.ACCESSROLE_CONFIRMED_USER,
		},
		"max_rps": nil,
		"description": []string{
			"Not too much to say about this one. Use the post payload as a template",
			"to create a store on the platform. If created successully the store id",
			"is added to your user record with the access role of owner. Once created",
			"a new token is sent in the headers. It may be useful, but not currently",
			"avalible, to throw a new token with updated user permissions (given the",
			"newly created store) to allow certain platforms access to it without",
			"requiring that the user re-login for access to a token.",
		},
	},
	"/store/search": {
		"control_method": "GET",
		"authenticate":   []string{},
		"max_rps":        nil,
		"description": []string{
			"This is a geo based search call to find stores near a user. For",
			"the time being the only required arguments are lon, lat and time.",
			"",
			"Sample request:",
			"	protocol://base_url:port/api/store/search?lon=-73.123&lat=40.123&time=400",
			"",
			"Some explination on how search works now. You are not able to specify",
			"a maximum distance for the search call. Nor is the stores maximum",
			"delivery distance taken into account. This is simply to provide all",
			"all stores within a staticaly defined disntance from the lon lat specified",
			"in the query params to show stores that may not deliver but will allow",
			"for pickup. This call is subject to change over time and the lack of",
			"complexity and buisness logic is for this reason. As the platform aquires",
			"more stores and the desire for facets/filters/fulltext search/etc become",
			"a reality we will certainly implent these. As of right now, we do not have",
			"the volume to provide such search features and exposing them will simply",
			"inhibt user experience.",
		},
	},
	"/store/info/:store_id": {
		"control_method": "GET",
		"authenticate":   []string{},
		"max_rps":        nil,
		"description": []string{
			"This call retuns a full store without any products",
		},
	},
	"/store/retrieve/full/:store_id": {
		"control_method": "GET",
		"authenticate": []string{
			models.ACCESSROLE_CONFIRMED_USER,
			models.ACCESSROLE_STOREOWNER,
		},
		"max_rps": nil,
		"description": []string{
			"Use this endpoint to retrieve a store with all of its",
			"enabled and disabled categories/products.",
		},
	},

	// cart
	"/cart/update/product/quantity": {
		"control_method": "POST",
		"post_payload":   models.CartRequest{},
		"authenticate":   []string{},
		"max_rps":        nil,
		"description": []string{
			"This api is meant to mutate the state of your cart.",
			"If you have never made this call before the payload will",
			"not be expecting a cart id. Store id is always required",
			"as carts are tied to users and stores. There may be",
			"**AT MOST** one --active-- cart per user per store. User to",
			"carts is one to many while active cart per store to users is",
			"1-1. You may re-activate your cart at any time by retreiving",
			"previous orders and reinstating your cart. This is actually",
			"duplicating your cart in the to avoid warping the data which",
			"leads to invalid analytics. If there is a desire to revive",
			"abandoned carts (an abandoned cart is simply an inactive cart",
			"that was never checked out) we can work on this feature.",
			"",
			"There is still much to do with carts. One big one is guest carts",
		},
	},

	"/transform": {
		"control_method": "GET",
		"post_payload":   models.User{},
		"authenticate": []string{
			models.ACCESSROLE_ADMIN,
		},
		"max_rps": nil,
		"desctiption": []string{
			"An api to translate jsons. Just a product of a late night.",
			"Its nothing to worry about and might be useful at some point",
			"but for now its just part of the code base and may discapear",
			"I'm not gonna write instructions for this since its not",
			"currently in use and subject to many changes. Its probably best",
			"Best not to expose it which is why it is restricted to admin use.",
		},
	},
}
