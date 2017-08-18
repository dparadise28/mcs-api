package networking

import (
	"api"
	"fmt"
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

var APIRouteMap map[string]map[string]interface{}

// TODO clean this shit up
var APIRouteList = []map[string]map[string]interface{}{
	map[string]map[string]interface{}{
		"/h2": {
			"control_method": "GET",
			"authenticate":   []string{},
			"max_rps":        100,
			"api_method":     Info,
			"description": []string{
				"Heartbeat",
			},
		},
	},
	map[string]map[string]interface{}{
		"/docs": {
			"control_method": "GET",
			"authenticate":   []string{},
			"max_rps":        100,
			"api_method":     "api.Docs",
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
				"",
				"",
				"Helper/Convenience apis:",
				"	In general these are supper lightweight apis and can be called",
				"	liberally (as always, within reason). These convenience apis",
				"	will become more useful as the platform grows and the number",
				"	of conventions we need to maintain grows with it. You can use",
				"	the results to help form valid store creation requests and",
				"	the like.",
			},
		},
	},

	// users api
	map[string]map[string]interface{}{
		"/user/create": {
			"control_method": "POST",
			"post_payload":   models.User{},
			"authenticate":   []string{},
			"max_rps":        nil,
			"api_method":     api.UserCreate,
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
	},
	map[string]map[string]interface{}{
		"/user/confirm/email/:user_id/:confirmation_code": {
			"control_method": "GET",
			"authenticate":   []string{},
			"max_rps":        nil,
			"api_method":     api.UserConfirmation,
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
	},
	map[string]map[string]interface{}{
		"/user/confirmation/resend": {
			"control_method": "GET",
			"authenticate":   []string{},
			"max_rps":        nil,
			"api_method":     api.UserResendConfirmation,
			"description": []string{
				"call structure: protocol://base_url:port/api/user/confirmation/resend?email={user_email}&password={password}",
				"the port will likely not be required but left in as an example if the env used to test",
				"is binded to that address/port",
				"",
				"When a user forgets his/her password this is the call that should be envoked to",
				"resend a confirmation code. The call WILL NOT rest the user password. What it does",
				"do is reset the confirmation code associated with the user to allow them to reset",
				"the password (if one is provided). If one is not provided then the link generated",
				"from this api only attempts to confirm the user. If one is provided and then the",
				"password is reset to the one in the confirmation link. The confirmation code is",
				"always reset after the call to confim the user is made and another can be generated",
				"at any time using this call. This ensures that acounts are safe given that the attacker",
				"would need to be able to access the users email to reset the pw. If a loose confirmation",
				"code is unaccounted for (in the case of one being generated, the pw is remembered and",
				"and it is never reset) the account is not vlunerable given that the attacker would",
				"need to be able to generate a pw that is useable by the platform to allow him access.",
				"This is assuming that the attacker is able to crack a randomly generated 32 char",
				"confirmation code, guess the salt and cost of the pw hash and the random 24 char",
				"user id hex.",
			},
		},
	},
	map[string]map[string]interface{}{
		"/user/login": {
			"control_method": "GET",
			"authenticate":   []string{},
			"max_rps":        nil,
			"api_method":     api.Login,
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
	},
	map[string]map[string]interface{}{
		"/user/retrieve": {
			"control_method": "GET",
			"authenticate": []string{
				models.ACCESSROLE_CONFIRMED_USER,
			},
			"max_rps":    nil,
			"api_method": api.GetUserById,
			"description": []string{
				"Simply a call to retrieve the user info. Must be confirmed to be allowed",
				"access to this call.",
			},
		},
	},
	map[string]map[string]interface{}{
		"/user/address_book/add": {
			"control_method": "POST",
			"post_payload":   models.Address{},
			"authenticate": []string{
				models.ACCESSROLE_CONFIRMED_USER,
			},
			"max_rps":    nil,
			"api_method": api.AddUserAddrToAddrBook,
			"description": []string{
				"Uniqueness of address records are quarenteed for (user_id, long, lat) pairs.",
				"This is not necessarily an atomic transaction. A count of the number of defaults",
				"in a users address book is retrieved. If a request is made to add the first",
				"address then that address is set as the default (irrespective of the state",
				"of the default field in the request payload). If the number of defaults for",
				"for the users address book associated with the request is greater then 0 and",
				"the request payload has the default flag set to true then the current defaults",
				"will be reset to false and the new address will become the users default address.",
				"In general it should be either 0 or 1; any other state is an edge case that is",
				"unintended and should not happen but is accounted for in the case of unforseen",
				"ways of acquiring unexpected states. In the case of the default flag being set",
				"to false and the number of defaults in the addr book is > 0 then a regular write",
				"is attempted.",
			},
		},
	},
	map[string]map[string]interface{}{
		"/user/address_book/retrieve": {
			"control_method": "GET",
			"authenticate": []string{
				models.ACCESSROLE_CONFIRMED_USER,
			},
			"max_rps":    nil,
			"api_method": api.GetUserAddrBook,
			"description": []string{
				"Nothing really special about this guy. It simply responds with the users full",
				"address book in no guarenteed order.",
			},
		},
	},
	map[string]map[string]interface{}{
		"/user/address_book/retrieve/default": {
			"control_method": "GET",
			"authenticate": []string{
				models.ACCESSROLE_CONFIRMED_USER,
			},
			"max_rps":    nil,
			"api_method": api.GetUserDefaultAddr,
			"description": []string{
				"Just a helper api to return the users default address.",
			},
		},
	},
	map[string]map[string]interface{}{
		"/user/address_book/default/change": {
			"control_method": "GET",
			"authenticate": []string{
				models.ACCESSROLE_CONFIRMED_USER,
			},
			"max_rps":    nil,
			"api_method": api.UpdateUserDefaultInAddrBook,
			"description": []string{
				"Update users default address by setting any current defaults in the users address",
				"book to regular addresses and updating the record associated with the required query",
				"param argument (address_id) to be the default new default.",
			},
		},
	},
	map[string]map[string]interface{}{
		"/user/address_book/remove_by_id": {
			"control_method": "GET",
			"authenticate": []string{
				models.ACCESSROLE_CONFIRMED_USER,
			},
			"max_rps":    nil,
			"api_method": api.RemoveUserAddrFromAddrBook,
			"description": []string{
				"Exactly what is sounds like. The only (intentional) gotchya here is that a default",
				"address may not be removed. You must change the default address first and then remove",
				"the address no longer needed. This is done so that we always have a default address",
				"for a user (if at least one address is stored). That way the client is more likely",
				"to perform searches while making the users standard experience one step simpler.",
			},
		},
	},
	map[string]map[string]interface{}{
		"/user/wallet/add": {
			"control_method": "GET",
			"authenticate": []string{
				models.ACCESSROLE_CONFIRMED_USER,
			},
			"max_rps":    nil,
			"api_method": api.CreateCustomerStripeReuseableAccount,
			"description": []string{
				"This is an interface to add a cc associated with the token in url query params",
				"(stripe_src=tok_str). If there is no stripe customer account associated with the",
				"user id in the request headers then one will automatically be created with the",
				"card provided in the url. If the user has already activated their waller then",
				"the card will be added to it. The stripe customer account id is then set in the",
				"jwt token for later user. You must either grab the token from the headers, response",
				"body, or re-login in order to proceed with actions to your wallet.",
			},
		},
	},
	map[string]map[string]interface{}{
		"/user/wallet/retrieve": {
			"control_method": "GET",
			"authenticate": []string{
				models.ACCESSROLE_CONFIRMED_USER,
			},
			"max_rps":    nil,
			"api_method": api.GetUserStipeCustomerAccount,
			"description": []string{
				"Simply retrieves a users wallet if one exists for the user making the request",
			},
		},
	},
	map[string]map[string]interface{}{
		"/user/wallet/change/default": {
			"control_method": "GET",
			"authenticate": []string{
				models.ACCESSROLE_CONFIRMED_USER,
			},
			"max_rps":    nil,
			"api_method": api.SetUserDefaultStipeCC,
			"description": []string{
				"Change the default cc in your wallet by adding new_default_cc=card_id in the",
				"query params.",
			},
		},
	},
	map[string]map[string]interface{}{
		"/user/wallet/remove": {
			"control_method": "GET",
			"authenticate": []string{
				models.ACCESSROLE_CONFIRMED_USER,
			},
			"max_rps":    nil,
			"api_method": api.DeleteUserStipeCC,
			"description": []string{
				"Remove a cc from your wallet by adding card_id=card_id in the query params.",
			},
		},
	},

	// stores api
	map[string]map[string]interface{}{
		"/store/create": {
			"control_method": "POST",
			"post_payload":   models.Store{},
			"authenticate": []string{
				models.ACCESSROLE_CONFIRMED_USER,
			},
			"max_rps":    nil,
			"api_method": api.StoreCreate,
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
	},
	map[string]map[string]interface{}{
		"/store/search": {
			"control_method": "GET",
			"authenticate":   []string{},
			"max_rps":        nil,
			"api_method":     api.StoreSearch,
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
	},
	map[string]map[string]interface{}{
		"/store/info/retrieve/:store_id": {
			"control_method": "GET",
			"authenticate":   []string{},
			"max_rps":        nil,
			"api_method":     api.GetStoreById,
			"description": []string{
				"This call retuns a full store without any products",
			},
		},
	},
	map[string]map[string]interface{}{
		"/store/info/update/:store_id": {
			"control_method": "POST",
			"post_payload":   models.Store{},
			"authenticate": []string{
				models.ACCESSROLE_CONFIRMED_USER,
			},
			"max_rps":    nil,
			"api_method": api.StoreInfoUpdate,
			"description": []string{
				"Store info update call.",
			},
		},
	},
	map[string]map[string]interface{}{
		"/store/retrieve/full/:store_id": {
			"control_method": "GET",
			"authenticate": []string{
				models.ACCESSROLE_CONFIRMED_USER,
				models.ACCESSROLE_STOREOWNER,
			},
			"max_rps":    nil,
			"api_method": api.GetFullStoreById,
			"description": []string{
				"Use this endpoint to retrieve a store with all of its",
				"enabled and disabled categories/products.",
			},
		},
	},

	// store categories
	map[string]map[string]interface{}{
		"/store/categories/retrieve/:store_id": {
			"control_method": "GET",
			"authenticate":   []string{},
			"max_rps":        nil,
			"api_method":     api.GetCategoriesByStoreId,
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
	},
	map[string]map[string]interface{}{
		"/store/category/create/:store_id": {
			"control_method": "POST",
			"authenticate": []string{
				models.ACCESSROLE_CONFIRMED_USER,
				models.ACCESSROLE_STOREOWNER,
			},
			"post_payload": models.Category{},
			"max_rps":      nil,
			"api_method":   api.AddStoreCategory,
			"description": []string{
				"Creates a new category and appends it as the",
				"last in the seq. It will not automatically be set",
				"to enabled however a category is not required to",
				"contain products as of the code in its current state",
				"For that reason, you can send your payload to include",
				"an enabled field and the category will show up in user",
				"subsequent user facing searches. We can change the",
				"search/retrieve apis to not include empty categories",
				"but it may be a nice feature for something like",
				"'CATEGORY_x COMING SOON!'",
			},
		},
	},
	map[string]map[string]interface{}{
		"/store/category/update/:store_id": {
			"control_method": "POST",
			"authenticate": []string{
				models.ACCESSROLE_CONFIRMED_USER,
				models.ACCESSROLE_STOREOWNER,
			},
			"post_payload": models.Category{},
			"max_rps":      nil,
			"api_method":   api.UpdateStoreCategory,
			"desctiption": []string{
				"This is an interface exposed to update the category name.",
				"",
				"required fields in post payload:",
				"	store_id, category_id, name",
			},
		},
	},
	map[string]map[string]interface{}{
		"/store/category/activate/:store_id": {
			"control_method": "POST",
			"authenticate": []string{
				models.ACCESSROLE_CONFIRMED_USER,
				models.ACCESSROLE_STOREOWNER,
			},
			"post_payload": models.Category{},
			"max_rps":      nil,
			"api_method":   api.EnableStoreCategory,
			"description": []string{
				"Toggles the store category enabled field based on the value",
				"it is set to in the request. This will affect user searches.",
				"",
				"This api is not fully fleshed out yet. I say this because",
				"reinstating a category means bringing it back in the full set",
				"of results. To do this means reordering the full set of store",
				"categories, or, simply assigning newly disbaled categories a",
				"sort index equal to the cap set on categories (which has not",
				"been introduced yet) and have the store reorder the categories",
				"through the ui with the exposed api method (if they choose to",
				"do so). I think the later is a better approach in terms of",
				"complexity and product design but its something to talk through",
				"since it is inherently a product design and affects user",
				"user interactions and ui/ux flow. If the former is chosen and",
				"reinstating a category in a specific position is desired then",
				"a sort index will be required along with the fields listed below.",
				"In either case more work will need to be done in order to update",
				"the api to the desired specs. The extent of the work depends on",
				"the path chosen.",
				"",
				"required fields in post payload:",
				"	store_id, category_id, enabled",
			},
		},
	},
	map[string]map[string]interface{}{
		"/store/categories/reorder/:store_id": {
			"control_method": "POST",
			"authenticate": []string{
				models.ACCESSROLE_CONFIRMED_USER,
				models.ACCESSROLE_STOREOWNER,
			},
			"post_payload": models.CategoryOrder{},
			"max_rps":      nil,
			"api_method":   api.ReorderStoreCategories,
			"description": []string{
				"This will group a transactional operation to update a batch of",
				"records based on the ids provided. As it currently stands, ALL",
				"active category ids must be specified in the payload in the order",
				"in which they are desired to appear in subsequent publicly facing",
				"category apis. The transaction can not be rolled back in its current",
				"state. That also means that if a batch fails then the previous",
				"order will remain unaltered.",
			},
		},
	},

	// platform assets
	map[string]map[string]interface{}{
		"/assets/image/search/:query_term": {
			"control_method": "GET",
			"authenticate":   []string{},
			"max_rps":        nil,
			"api_method":     api.SearchAssets,
			"description": []string{
				"Api for finding images loaded in the platform given a query",
				"term specified in the url. The results are capped at 100 items",
				"and are limitted to what is currently loaded in the db serving",
				"the env the api is currently on",
			},
		},
	},
	map[string]map[string]interface{}{
		"/assets/image/upload": {
			"control_method": "POST",
			"authenticate":   []string{},
			"post_payload":   models.NewAsset{},
			"max_rps":        nil,
			"api_method":     api.CreateAsset,
			"description": []string{
				"Api for uploading assets",
			},
		},
	},

	// store products
	map[string]map[string]interface{}{
		"/store/category/product/create/:store_id": {
			"control_method": "POST",
			"authenticate": []string{
				models.ACCESSROLE_CONFIRMED_USER,
				models.ACCESSROLE_STOREOWNER,
			},
			"post_payload": models.Product{},
			"max_rps":      nil,
			"api_method":   api.AddStoreProduct,
			"description": []string{
				"Creates a new product and attaches it to the end of the",
				"category specified in the payload. This will automatically",
				"set the order index as the last item in the category.",
				"It will not automatically enable the product that is being",
				"loaded. If the enabled field in the payload is set to true",
				"it will show up in user searches (assuming the category it",
				"is being added to is also enabled)",
			},
		},
	},
	map[string]map[string]interface{}{
		"/store/category/product/update/:store_id": {
			"control_method": "POST",
			"authenticate": []string{
				models.ACCESSROLE_CONFIRMED_USER,
				models.ACCESSROLE_STOREOWNER,
			},
			"post_payload": models.Product{},
			"max_rps":      nil,
			"api_method":   api.UpdateStoreProduct,
			"desctiption": []string{
				"This is an interface exposed to update a full product record.",
				"",
				"required fields in post payload:",
				"	description, title, price_cents, category_id, new_category_id",
				"",
				"If the category_id provided is different from that of the",
				"new_category_id the category will be thrown into the new",
				"category with the same order index as it had prior to the",
				"update. Ordering can not be mutated in this call. If this",
				"becomes a desired feature then we can bake that in but",
				"the overhead/complexity of doing so now did not make sense.",
			},
		},
	},
	map[string]map[string]interface{}{
		"/store/category/product/activate/:store_id": {
			"control_method": "POST",
			"authenticate": []string{
				models.ACCESSROLE_CONFIRMED_USER,
				models.ACCESSROLE_STOREOWNER,
			},
			"post_payload": models.Product{},
			"max_rps":      nil,
			"api_method":   api.EnableStoreProduct,
			"description": []string{
				"Toggles the store product enabled field based on the value",
				"it's set to in the request. This will affect user searches.",
				"",
				"The same disclaimers stated in /store/category/activate/:store_id",
				"can analogously be said about this call. For more information",
				"please refer to the description provided there.",
				"",
				"required fields in post payload:",
				"	store_id, category_id, product_id, enabled",
			},
		},
	},
	map[string]map[string]interface{}{
		"/store/category/products/reorder/:store_id": {
			"control_method": "POST",
			"authenticate": []string{
				models.ACCESSROLE_CONFIRMED_USER,
				models.ACCESSROLE_STOREOWNER,
			},
			"post_payload": models.ProductOrder{},
			"max_rps":      nil,
			"api_method":   api.ReorderStoreProducts,
			"description": []string{
				"please refer to the analogous categories call: ",
				"	/store/categories/reorder/:store_id",
			},
		},
	},

	// cart
	map[string]map[string]interface{}{
		"/cart/update/product/quantity": {
			"control_method": "POST",
			"post_payload":   models.CartRequest{},
			"authenticate": []string{
				models.ACCESSROLE_CONFIRMED_USER,
			},
			"max_rps":    nil,
			"api_method": api.UpdateCartProductQuantity,
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
				"Another thing we may want to do is update this api to track",
				"previous quantities and not only final quantity states prior",
				"to cart state changes (active, abandoned, completed) to retain",
				"a full picture of realstive user interactions. Howevere without",
				"any data modeling to determine the necessity of such tracking",
				"its hard to justify the cost with respect to the additional",
				"effort required to implement such a feature given no data",
				"equates to a lack of evidence of those actions being meaningful",
				"to analytical inferences. We will probably do this when there",
				"is sufficient volume flowing through the system that will allow",
				"us to intelligently make such decisions.",
				"",
				"There is still much to do with carts. One big one is guest carts",
			},
		},
	},
	map[string]map[string]interface{}{
		"/carts/retrieve/active": {
			"control_method": "GET",
			"post_payload":   nil,
			"authenticate": []string{
				models.ACCESSROLE_CONFIRMED_USER,
			},
			"max_rps":    nil,
			"api_method": api.RetrieveUserActiveCarts,
			"description": []string{
				"Retrieve all active carts for the user identified by the id in",
				"the headers.",
			},
		},
	},
	map[string]map[string]interface{}{
		"/carts/retrieve/completed": {
			"control_method": "GET",
			"post_payload":   nil,
			"authenticate": []string{
				models.ACCESSROLE_CONFIRMED_USER,
			},
			"max_rps":    nil,
			"api_method": api.RetrieveUserCompletedCarts,
			"description": []string{
				"Retrieve all completed carts for the user identified by the id in",
				"the headers.",
			},
		},
	},
	map[string]map[string]interface{}{
		"/cart/abandon/:cart_id": {
			"control_method": "GET",
			"post_payload":   nil,
			"authenticate": []string{
				models.ACCESSROLE_CONFIRMED_USER,
			},
			"max_rps":    nil,
			"api_method": api.AbandonUserActiveCart,
			"description": []string{
				"This is analogous to a clear/drop cart button you would",
				"normally see on a ui. We do not delete carts. We simply",
				"update their state to reflect the user is no longer interested",
				"in their cart in its current state. We do this for future",
				"analytics inferences we may want to make from user interations.",
			},
		},
	},
	map[string]map[string]interface{}{
		"/cart/re-activate/:cart_id": {
			"control_method": "GET",
			"post_payload":   nil,
			"authenticate": []string{
				models.ACCESSROLE_CONFIRMED_USER,
			},
			"max_rps":    nil,
			"api_method": api.ReActiveUserCart,
			"description": []string{
				"This will basically duplicate a previously checked out cart.",
				"If you currently have an active cart associated with the same",
				"store as the cart about to be re-activated this call will error",
				"out with the approriate error message telling you you must",
				"drop that active cart first. If try to re activate a cart with",
				"any other state then completed then this call will not work for",
				"that either. The carts safe to reactivate (with the exception of",
				"of potentially already having an active cart) can be retrieved",
				"from the /carts/retrieve/completed api.",
			},
		},
	},

	// payments
	map[string]map[string]interface{}{
		"/payment/store/create/account": {
			"control_method": "POST",
			"post_payload":   models.Store{},
			"authenticate": []string{
				models.ACCESSROLE_CONFIRMED_USER,
			},
			"max_rps":    nil,
			"api_method": api.CreateStoreStripeCustomAccount,
			"description": []string{
				"create a new custom stripe account for onboarding stores",
			},
		},
	},

	// constants
	map[string]map[string]interface{}{
		"/helper/payment/methods": {
			"control_method": "GET",
			"post_payload":   nil,
			"authenticate":   []string{},
			"max_rps":        nil,
			"api_method":     api.PaymentMethods,
			"description": []string{
				"Just a little helper method to expose the currently availible",
				"payment methods.",
				"This is a supper lightweight method and can be",
				"called liberally (as always, within reason). These convenience",
				"methods will become more useful as the platform grows and the",
				"number of conventions we need to maintain grows with it. You can",
				"use the results to help form valid store creation requests and",
				"the like.",
			},
		},
	},
	map[string]map[string]interface{}{
		"/helper/order/methods": {
			"control_method": "GET",
			"post_payload":   nil,
			"authenticate":   []string{},
			"max_rps":        nil,
			"api_method":     api.OrderMethods,
			"description": []string{
				"Another convenince method for exposing currently availible order",
				"methods.",
			},
		},
	},
	map[string]map[string]interface{}{
		"/helper/order/statuses/all": {
			"control_method": "GET",
			"post_payload":   nil,
			"authenticate":   []string{},
			"max_rps":        nil,
			"api_method":     api.AllOrderStatuses,
			"description": []string{
				"Convenince method for exposing all currently supported order",
				"statuses for the various order methods.",
			},
		},
	},
	map[string]map[string]interface{}{
		"/helper/order/statuses/delivery": {
			"control_method": "GET",
			"post_payload":   nil,
			"authenticate":   []string{},
			"max_rps":        nil,
			"api_method":     api.DeliveryOrderStatuses,
			"description": []string{
				"Convenince method for exposing currently supported delivery order",
				"statuses.",
			},
		},
	},
	map[string]map[string]interface{}{
		"/helper/order/statuses/pickup": {
			"control_method": "GET",
			"post_payload":   nil,
			"authenticate":   []string{},
			"max_rps":        nil,
			"api_method":     api.PickupOrderStatuses,
			"description": []string{
				"Convenince method for exposing currently supported pickup order",
				"statuses.",
			},
		},
	},

	// reviews
	map[string]map[string]interface{}{
		"/review/order/add": {
			"control_method": "POST",
			"post_payload":   models.Review{},
			"authenticate":   []string{},
			"max_rps":        nil,
			"api_method":     api.ReviewOrder,
			"description": []string{
				"This call follows the same convention as the store review but also",
				"requires an order_id.",
			},
		},
	},
	map[string]map[string]interface{}{
		"/review/store/add": {
			"control_method": "POST",
			"post_payload":   models.Review{},
			"authenticate":   []string{},
			"max_rps":        nil,
			"api_method":     api.ReviewStore,
			"description": []string{
				"This call is exposed to review a store. It is pretty interactive",
				"and should display valid interactive errors. The required fields",
				"for this call are score, and store_id. The other fields will be",
				"filled out automatically for this call, with the exception of the",
				"optional field, comment, which you can specify any time. The max",
				"field length allowed for the comment is 200 characters.",
			},
		},
	},
	map[string]map[string]interface{}{
		"/review/platform/add": {
			"control_method": "POST",
			"post_payload":   models.Review{},
			"authenticate":   []string{},
			"max_rps":        nil,
			"api_method":     api.ReviewPlatform,
			"description": []string{
				"This follows the same convention as adding a store review but",
				"without the restriction of requiring a store id and there is no",
				"api exposed to retrieve these reviews. We can expose this when",
				"we have a platform admin ui to view these.",
			},
		},
	},
	map[string]map[string]interface{}{
		"/review/store/retrieve": {
			"control_method": "GET",
			"post_payload":   nil,
			"authenticate":   []string{},
			"max_rps":        nil,
			"api_method":     api.GetStoreReviews,
			"description": []string{
				"Api to expose all reviews currently made by customers of the store",
				"specified in the query params (url?store_id=store_id)",
			},
		},
	},

	map[string]map[string]interface{}{
		"/transform": {
			"control_method": "GET",
			"post_payload":   models.User{},
			"authenticate": []string{
				models.ACCESSROLE_ADMIN,
			},
			"max_rps":    nil,
			"api_method": tools.RemodelJ,
			"desctiption": []string{
				"An api to translate jsons. Just a product of a late night.",
				"Its nothing to worry about and might be useful at some point",
				"but for now its just part of the code base and may discapear",
				"I'm not gonna write instructions for this since its not",
				"currently in use and subject to many changes. Its probably best",
				"Best not to expose it which is why it is restricted to admin use.",
			},
		},
	},
}

func init() {
	routeMap := map[string]map[string]interface{}{}

	for _, route := range APIRouteList {
		for routeEndPoint, routeSpecs := range route {
			routeMap[routeEndPoint] = routeSpecs

			switch v := routeSpecs["api_method"].(type) {
			case string:
				routeSpecs["api_method"] = Docs
				fmt.Println(v)
			}
		}
	}
	APIRouteMap = routeMap
}
