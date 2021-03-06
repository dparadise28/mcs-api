package models

import (
	"encoding/json"
	"net/http"
	"strings"
)

type Error struct {
	Id        string              `json:"id"`
	Status    int                 `json:"status"`
	Title     string              `json:"title"`
	Details   map[string][]string `json:"details"`
	JsonError error               `json:"json_error,omitempty"`
}

func WriteError(w http.ResponseWriter, err *Error) {
	w.WriteHeader(err.Status)
	json.NewEncoder(w).Encode(err)
}

func WriteNewError(w http.ResponseWriter, err error) {
	r := strings.NewReplacer(
		"json: ", "",
		"unmarshal", "deserialize",
		"Go struct", "json",
	)
	ErrBadRequest = &Error{
		"FAILED", 400, "Bad request", map[string][]string{
			"malformed": []string{
				strings.ToLower(r.Replace(err.Error())),
			},
		}, err,
	}

	w.WriteHeader(ErrBadRequest.Status)
	json.NewEncoder(w).Encode(ErrBadRequest)
}

var (
	ErrSuccess = &Error{
		"SUCCESS", 200, "Successful request", map[string][]string{}, nil,
	}
	ErrRequestTimeout = &Error{
		"gateway_timeout", 504, "Gateway Timeout", map[string][]string{
			"timeout": []string{
				"External resource unavailable.",
			},
		}, nil,
	}
	ErrResourceConflict = &Error{
		"resource_conflict", 409, "resource_conflict", map[string][]string{
			"conflict": []string{
				"The request could not be completed due to a conflict with the current state of the resource.",
			},
		}, nil,
	}
	ErrUnauthorizedAccess = &Error{
		"bad_request", 403, "Unauthorized Access", map[string][]string{
			"unauthorized": []string{
				"Unauthorized Access",
			},
		}, nil,
	}
	ErrExpiredJWToken = &Error{
		"bad_request", 403, "Expired Token", map[string][]string{
			"expired": []string{
				"Please log back",
			},
		}, nil,
	}
	ErrUnconfirmedUser = &Error{
		"bad_request", 403, "Unauthorized Access", map[string][]string{
			"unauthorize": []string{
				"Please Confirm Your Email Address",
			},
		}, nil,
	}
	ErrBadRequest = &Error{
		"bad_request", 400, "Bad request", map[string][]string{
			"malformed": []string{
				"Malformed request body.",
			},
		}, nil,
	}
	ErrMissingPayload = &Error{
		"bad_request", 400, "Bad request", map[string][]string{
			"json_decoding": []string{
				"Malformed request body. Must Provide JSON object.",
			},
		}, nil,
	}
	ErrInternalServer = &Error{
		"internal_server_error", 500, "Internal Server Error", map[string][]string{
			"internal_server_error": []string{
				"Something went wrong.",
			},
		}, nil,
	}
)
