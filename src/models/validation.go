package models

import (
	"gopkg.in/go-playground/validator.v9"
)

// declare global package scope var (cant be used outside of api
var validate *validator.Validate
