package networking

import (
	"github.com/dgrijalva/jwt-go"
	//"github.com/dgrijalva/jwt-go/request"
	hr "github.com/julienschmidt/httprouter"
	"log"
	"models"
	"net/http"
	"strings"
)

var AccessControlAllowMethods = "GET, OPTION, HEAD, PATCH, PUT, POST, DELETE"

func SetBaseHeaders(respWrtr http.ResponseWriter) {
	respWrtr.Header().Set("Sec-Websocket-Version", "13")
	respWrtr.Header().Set("Access-Control-Allow-Origin", "*")
	respWrtr.Header().Set("Content-Type", "application/json")
	respWrtr.Header().Set("Access-Control-Allow-Methods", AccessControlAllowMethods)
}

func ValidatedToken(w http.ResponseWriter, r *http.Request, ps hr.Params, ep string) (bool, *models.Error) {
	// ps (http router params); ep (endpoint found in src/networking/route.go::APIRouteMap)
	tokenStr, e := r.Cookie("AUTH-TOKEN")
	if e != nil {
		log.Println("Missing Token")
		return false, models.ErrUnauthorizedAccess
	}
	log.Println(tokenStr, e)
	token, err := jwt.ParseWithClaims(tokenStr.Value, &models.CustomClaims{},
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
			authCookie := http.Cookie{
				Name:    "AUTH-TOKEN",
				Path:    "/api",
				Value:   updatedToken,
				Expires: models.COOKIE_TTL(),
			}
			http.SetCookie(w, &authCookie)
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
