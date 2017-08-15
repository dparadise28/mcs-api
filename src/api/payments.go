package api

import (
	"db"
	"encoding/json"
	"errors"
	"github.com/julienschmidt/httprouter"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/account"
	"github.com/stripe/stripe-go/card"
	"github.com/stripe/stripe-go/customer"
	//"gopkg.in/mgo.v2/bson"
	//"log"
	"models"
	"net/http"
	"strings"
	"time"
	"tools"
)

func CreateStoreStripeCustomAccount(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var store models.Store
	v := new(tools.DefaultValidator)
	if err := json.NewDecoder(r.Body).Decode(&store.PaymentDetails); err != nil {
		models.WriteNewError(w, err)
		return
	}
	if validationErr := v.ValidateIncomingJsonRequest(&store.PaymentDetails); validationErr.Status != 200 {
		models.WriteError(w, &validationErr)
		return
	}

	resp, err := CreateStoreStripeCustomAccountImpl(w, r, ps, &store)
	if err != nil {
		models.WriteNewError(w, err)
	}
	json.NewEncoder(w).Encode(resp)
}

// ::TODO:: this should prob be part of the store.PaymentDetails struct
func CreateStoreStripeCustomAccountImpl(w http.ResponseWriter, r *http.Request, ps httprouter.Params, store *models.Store) (*stripe.Account, error) {
	stripe.Key = models.StripeSK
	act := stripe.Account{}
	BusinessType := stripe.Company
	if store.PaymentDetails.BusinessType == "individual" {
		BusinessType = stripe.Individual
	} else if store.PaymentDetails.BusinessType != "company" {
		return &act, errors.New("Invalid option provided for business_type. Field must either individual or company")
	}
	params := &stripe.AccountParams{
		Type:    stripe.AccountTypeCustom,
		Country: models.CountryNameToCountryCode[store.PaymentDetails.LegalEntity.BillingAddress.Country],
		TOSAcceptance: &stripe.TOSAcceptanceParams{
			IP:        strings.Split(r.RemoteAddr, ":")[0],
			Date:      time.Now().Unix(),
			UserAgent: strings.Join(r.Header["User-Agent"], " "),
		},
		LegalEntity: &stripe.LegalEntity{
			PersonalIDProvided: true,
			BusinessTaxID:      store.PaymentDetails.LegalEntity.BusinessTaxID,
			BusinessName:       store.PaymentDetails.LegalEntity.BusinessName,
			PersonalID:         store.PaymentDetails.LegalEntity.PersonalID,
			First:              store.PaymentDetails.LegalEntity.Owner.First,
			Last:               store.PaymentDetails.LegalEntity.Owner.Last,
			Type:               BusinessType,
			SSN:                store.PaymentDetails.LegalEntity.SSNLast4,
			DOB: stripe.DOB{
				Day:   int(store.PaymentDetails.LegalEntity.Owner.DOB.Day),
				Month: int(store.PaymentDetails.LegalEntity.Owner.DOB.Month),
				Year:  int(store.PaymentDetails.LegalEntity.Owner.DOB.Year),
			},
			Address: stripe.Address{
				Country: models.CountryNameToCountryCode[store.PaymentDetails.LegalEntity.BillingAddress.Country],
				City:    store.PaymentDetails.LegalEntity.BillingAddress.City,
				Zip:     store.PaymentDetails.LegalEntity.BillingAddress.PostalCode,
				Line1:   store.PaymentDetails.LegalEntity.BillingAddress.Line1,
				State:   store.PaymentDetails.LegalEntity.BillingAddress.AdminAreaLvl1,
			},
		},
		ExternalAccount: &stripe.AccountExternalAccountParams{
			Token:    r.URL.Query().Get("stripe_src"),
			Country:  models.CountryNameToCountryCode[store.PaymentDetails.LegalEntity.BillingAddress.Country],
			Currency: "usd",
		},
	}
	new_act, err := account.New(params)
	if err != nil {
		return &act, err
	} else {
		return new_act, nil
	}
}

func CreateCustomerStripeReuseableAccount(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var user models.User
	token, _ := models.GetJWTContent(r.Header.Get(models.JWT_COOKIE_NAME))
	claims, ok := token.Claims.(*models.CustomClaims)
	user.GetByIdStr(db.Database, r.Header.Get(models.USERID_COOKIE_NAME))
	user.StripeToken = r.URL.Query().Get("stripe_src")
	if user.StripeToken == "" {
		models.WriteNewError(w, errors.New(
			"Please provide a card token associated with the"+
				" card you would like to add to your wallet in"+
				" the query params (stripe_src=tok_str).",
		))
	}
	user.DB = db.Database
	user.DBSession = user.DB.Session.Copy()
	stripe.Key = models.StripeSK

	if user.StripeCustomerID == claims.StripeCustomerID {
		if ok && len(claims.StripeCustomerID) == 0 {
			customerParams := &stripe.CustomerParams{
				Email: user.Email,
				Desc:  user.ID.Hex(),
			}
			customerParams.SetSource(user.StripeToken)
			customer, customerCreateErr := customer.New(customerParams)
			if customerCreateErr != nil {
				json.NewEncoder(w).Encode(customerCreateErr)
				return
			}
			user.StripeCustomerID = customer.ID
			if err := user.AddUserStripeCustomerAccount(); err != nil {
				models.WriteNewError(w, err)
				return
			}
			user.UpdateTokenAndCookie(w)
			json.NewEncoder(w).Encode(user)
			return
		}
		target, err := card.New(&stripe.CardParams{
			Customer: user.StripeCustomerID,
			Token:    user.StripeToken,
		})
		if err != nil {
			json.NewEncoder(w).Encode(err)
			return
		}
		user.StripeCustomerID = claims.StripeCustomerID
		user.UpdateTokenAndCookie(w)
		json.NewEncoder(w).Encode(target)
		return
	}
	models.WriteNewError(w, errors.New("It is unacceptable to impersonate other users here."))
}

func GetUserStipeCustomerAccount(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var user models.User
	token, _ := models.GetJWTContent(r.Header.Get(models.JWT_COOKIE_NAME))
	claims, _ := token.Claims.(*models.CustomClaims)
	user.GetByIdStr(db.Database, r.Header.Get(models.USERID_COOKIE_NAME))
	user.DB = db.Database
	user.DBSession = user.DB.Session.Copy()
	stripe.Key = models.StripeSK

	if user.StripeCustomerID == claims.StripeCustomerID {
		if len(user.StripeCustomerID) == 0 {
			json.NewEncoder(w).Encode(stripe.Customer{})
			return
		}
		c, err := customer.Get(user.StripeCustomerID, nil)
		if err != nil {
			json.NewEncoder(w).Encode(err)
			return
		}
		json.NewEncoder(w).Encode(c)
		return
	}
	models.WriteNewError(w, errors.New("It is unacceptable to impersonate other users here."))
}

func SetUserDefaultStipeCC(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var user models.User
	token, _ := models.GetJWTContent(r.Header.Get(models.JWT_COOKIE_NAME))
	claims, _ := token.Claims.(*models.CustomClaims)
	user.GetByIdStr(db.Database, r.Header.Get(models.USERID_COOKIE_NAME))
	user.DB = db.Database
	user.DBSession = user.DB.Session.Copy()
	stripe.Key = models.StripeSK

	if user.StripeCustomerID == claims.StripeCustomerID {
		if len(user.StripeCustomerID) == 0 {
			json.NewEncoder(w).Encode(stripe.Customer{})
			return
		}
		c, err := customer.Update(user.StripeCustomerID, &stripe.CustomerParams{
			DefaultSource: r.URL.Query().Get("new_default_cc"),
		})
		if err != nil {
			json.NewEncoder(w).Encode(err)
			return
		}
		json.NewEncoder(w).Encode(c)
		return
	}
	models.WriteNewError(w, errors.New("It is unacceptable to impersonate other users here."))
}

func DeleteUserStipeCC(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var user models.User
	token, _ := models.GetJWTContent(r.Header.Get(models.JWT_COOKIE_NAME))
	claims, _ := token.Claims.(*models.CustomClaims)
	user.GetByIdStr(db.Database, r.Header.Get(models.USERID_COOKIE_NAME))
	user.DB = db.Database
	user.DBSession = user.DB.Session.Copy()
	stripe.Key = models.StripeSK

	if user.StripeCustomerID == claims.StripeCustomerID {
		if len(user.StripeCustomerID) == 0 {
			json.NewEncoder(w).Encode(stripe.Customer{})
			return
		}
		c, err := card.Del(
			r.URL.Query().Get("card_id"),
			&stripe.CardParams{Customer: user.StripeCustomerID},
		)
		if err != nil {
			json.NewEncoder(w).Encode(err)
			return
		}
		json.NewEncoder(w).Encode(c)
		return
	}
	models.WriteNewError(w, errors.New("It is unacceptable to impersonate other users here."))
}
