package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var jwtKey = []byte("bebasapasaja:134")
var tokenName = "token"

type Claims struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	UserType int    `json:"user_type"`
	jwt.StandardClaims
}

func GenerateToken(w http.ResponseWriter, id int, name string, userType int) {
	tokenExpireTime := time.Now().Add(5 * time.Minute)

	// create claims with user data
	claims := &Claims{
		ID:       id,
		Name:     name,
		UserType: userType,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: tokenExpireTime.Unix(),
		},
	}

	// encrypt claims to jwt token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(jwtKey)
	if err != nil {
		return
	}

	// set token to cookie
	http.SetCookie(w, &http.Cookie{
		Name:     tokenName,
		Value:    signedToken,
		Expires:  tokenExpireTime,
		Secure:   false,
		HttpOnly: true,
	})
}

func ResetUserToken(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     tokenName,
		Value:    "",
		Expires:  time.Now(),
		Secure:   false,
		HttpOnly: true,
	})
}

func SendUnAuthorizedResponse(w http.ResponseWriter) {
	// success response
	var response MessageResponse
	response.Status = 401
	response.Message = "UnAuthorized Attempt"

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func Authenticate(next http.HandlerFunc, accessType int) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		isValidToken := ValidateUserToken(r, accessType)
		if !isValidToken {
			SendUnAuthorizedResponse(w)
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

func ValidateUserToken(r *http.Request, accessType int) bool {
	isAccessTokenValid, id, email, userType := ValidateTokenFormCookies(r)
	fmt.Print(id, email, userType, accessType, isAccessTokenValid)

	if isAccessTokenValid {
		isUserValid := userType == accessType
		if isUserValid {
			return true
		}
	}
	return false
}

func ValidateTokenFormCookies(r *http.Request) (bool, int, string, int) {

	if cookie, err := r.Cookie(tokenName); err != nil {
		accessToken := cookie.Value
		accessClaims := &Claims{}
		parsedToken, err := jwt.ParseWithClaims(accessToken, accessClaims, func(accessToken *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil && parsedToken.Valid {
			return true, accessClaims.ID, accessClaims.Name, accessClaims.UserType
		}
	}
	return false, -1, "", -1
}
