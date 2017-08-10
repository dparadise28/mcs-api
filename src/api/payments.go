package api

import (
	//"db"
	"encoding/json"
	"errors"
	"github.com/julienschmidt/httprouter"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/account"
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
	if err := json.NewDecoder(r.Body).Decode(&store); err != nil {
		models.WriteNewError(w, err)
		return
	}
	if validationErr := v.ValidateIncomingJsonRequest(&store); validationErr.Status != 200 {
		models.WriteError(w, &validationErr)
		return
	}

	resp, err := CreateStoreStripeCustomAccountImpl(w, r, ps, &store)
	if err != nil {
		models.WriteNewError(w, err)
	}
	json.NewEncoder(w).Encode(resp)
}

// ::TODO:: this should prob be part of the store struct
func CreateStoreStripeCustomAccountImpl(w http.ResponseWriter, r *http.Request, ps httprouter.Params, store *models.Store) (*stripe.Account, error) {
	stripe.Key = models.StripeSK
	act := stripe.Account{}
	BusinessType := stripe.Company
	if store.BusinessType == "individual" {
		BusinessType = stripe.Individual
	} else if store.BusinessType != "company" {
		return &act, errors.New("Invalid option provided for business_type. Field must either individual or company")
	}
	params := &stripe.AccountParams{
		Type:    stripe.AccountTypeCustom,
		Country: models.CountryNameToCountryCode[store.LegalEntity.BillingAddress.Country],
		TOSAcceptance: &stripe.TOSAcceptanceParams{
			IP:        strings.Split(r.RemoteAddr, ":")[0],
			Date:      time.Now().Unix(),
			UserAgent: strings.Join(r.Header["User-Agent"], " "),
		},
		LegalEntity: &stripe.LegalEntity{
			PersonalIDProvided: true,
			BusinessTaxID:      store.LegalEntity.BusinessTaxID,
			BusinessName:       store.LegalEntity.BusinessName,
			PersonalID:         store.LegalEntity.PersonalID,
			First:              store.LegalEntity.Owner.First,
			Last:               store.LegalEntity.Owner.Last,
			Type:               BusinessType,
			SSN:                store.LegalEntity.SSNLast4,
			DOB: stripe.DOB{
				Day:   int(store.LegalEntity.Owner.DOB.Day),
				Month: int(store.LegalEntity.Owner.DOB.Month),
				Year:  int(store.LegalEntity.Owner.DOB.Year),
			},
			Address: stripe.Address{
				Country: models.CountryNameToCountryCode[store.LegalEntity.BillingAddress.Country],
				City:    store.LegalEntity.BillingAddress.City,
				Zip:     store.LegalEntity.BillingAddress.PostalCode,
				Line1:   store.LegalEntity.BillingAddress.Line1,
				State:   store.LegalEntity.BillingAddress.AdminAreaLvl1,
			},
		},
		ExternalAccount: &stripe.AccountExternalAccountParams{
			Token:    r.URL.Query().Get("stripe_src"),
			Country:  models.CountryNameToCountryCode[store.LegalEntity.BillingAddress.Country],
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

/*func CreateCustomerStripeReuseableAccount(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	//var store models.Stripe
	stripe_email = r.URL.Query().Get("store_email")
	stripe_src = r.URL.Query().Get("stripe_src")
	stripe.Key = models.StripeSK

	customerParams := &stripe.CustomerParams{
		Email: stripe_email,
	}
	customerParams.SetSource(stripe_src)
	if c, err := customer.New(customerParams); err != nil {
		models.WriteNewError(w, err)
	} else {
		json.NewEncoder(w).Encode(c)
	}
}*/

/*
	v := new(tools.DefaultValidator)
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		models.WriteError(w, models.ErrBadRequest)
		return
	}
	if validationErr := v.ValidateIncomingJsonRequest(&product); validationErr.Status != 200 {
		models.WriteError(w, &validationErr)
		return
	}
	product.DB = db.Database
	product.DBSession = product.DB.Session.Copy()
	defer product.DBSession.Close()
	if err := product.AddProductToStoreCategory(); err != nil {
		models.WriteError(w, models.ErrBadRequest)
	}
	json.NewEncoder(w).Encode(product)
}
*/