package networking

import (
	"github.com/dgrijalva/jwt-go"
	hr "github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
	"log"
	"models"
	"net/http"
	"strings"
)

var AccessControlAllowMethods = "GET, OPTION, HEAD, PATCH, PUT, POST, DELETE"

func SetBaseHeaders(respWrtr http.ResponseWriter, req *http.Request) {
	respWrtr.Header().Set("Sec-Websocket-Version", "13")
	respWrtr.Header().Set("Access-Control-Allow-Origin", req.Header.Get("Origin"))
	respWrtr.Header().Set("Content-Type", "application/json")
	respWrtr.Header().Set("Access-Control-Allow-Methods", AccessControlAllowMethods)
	respWrtr.Header().Set("Access-Control-Allow-Credentials", "true")
	respWrtr.Header().Set(
		"Access-Control-Allow-Headers",
		"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Requested-With, userID, authtoken",
	)
	c := cors.New(cors.Options{
		AllowedOrigins: []string{req.Header.Get("Origin")},
	})
	c.HandlerFunc(respWrtr, req)
}

func ValidatedToken(w http.ResponseWriter, r *http.Request, ps hr.Params, ep string) (bool, *models.Error) {
	tokenStr := r.Header.Get(models.JWT_COOKIE_NAME)
	if tokenStr == "" {
		log.Println("Missing Token")
		return false, models.ErrUnauthorizedAccess
	}
	log.Println(tokenStr)
	token, err := jwt.ParseWithClaims(tokenStr, &models.CustomClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(models.JWT_SIGNATURE), nil
		},
	)

	if err == nil {
		if claims, ok := token.Claims.(*models.CustomClaims); ok && token.Valid {
			log.Println(ep)
			for _, accessRole := range APIRouteMap[ep]["authenticate"].([]string) {
				if accessRole == models.ACCESSROLE_CONFIRMED_USER {
					if !claims.Confirmed {
						log.Println("Unconfirmed")
						return false, models.ErrUnconfirmedUser
					}
				} else {
					if ps.ByName("store_id") != "" {
						perms, found := claims.Perms[ps.ByName("store_id")]
						if !found || !strings.Contains(perms, accessRole) {
							log.Println("unauthorize request!")
							return false, models.ErrUnauthorizedAccess
						}
						if !strings.Contains(perms, accessRole) {
							return false, models.ErrUnauthorizedAccess
						}
					} else {
						log.Println("unauthorize: ", accessRole)
						return false, models.ErrUnauthorizedAccess
					}
				}
			}
			// checks pass so lets update the token with the new expiration time
			updatedClaims := models.GenerateTokenClaims(claims.Perms, claims.Confirmed)
			updatedToken, _ := updatedClaims.CreateToken()
			w.Header().Set(models.JWT_COOKIE_NAME, updatedToken)
			return true, models.ErrSuccess
		} else {
			log.Println("InvalidToken: ", token, tokenStr)
			return false, models.ErrUnauthorizedAccess
		}
	} else {
		if strings.Contains(err.Error(), "token is expired by") {
			return false, models.ErrExpiredJWToken
		}
		return false, models.ErrUnauthorizedAccess
	}
}
