package api

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"models"
	"net/http"
	"tools"
)

func SendEmail(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var email models.EmailReq
	v := new(tools.DefaultValidator)
	if err := json.NewDecoder(r.Body).Decode(&email); err != nil {
		models.WriteError(w, models.ErrBadRequest)
		return
	}
	if validationErr := v.ValidateIncomingJsonRequest(&email); validationErr.Status != 200 {
		models.WriteError(w, &validationErr)
		return
	}
	tools.EmailQueue <- &models.Email{
		To:      email.To,
		Body:    email.Body,
		Subject: email.Subject,
	}
}
