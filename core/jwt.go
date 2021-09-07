package core

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

var (
	mySigningKey = []byte("zozoo.net")
	expireAt     = time.Hour * 30 * 24
)

// MyCustomClaims 定制
type MyCustomClaims struct {
	UserID string `json:"user_id"`
	jwt.StandardClaims
}

//生成jwtToken
func GenerateToken(userID string) (string, error) {

	// Create the Claims
	claims := &MyCustomClaims{
		userID,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(expireAt).Unix(),
			Issuer:    "ganganlee",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(mySigningKey)
}

//验证jwtToken
func ValidateToken(tokenString string) (string, error) {

	token, err := jwt.ParseWithClaims(tokenString, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return mySigningKey, nil
	})

	if claims, ok := token.Claims.(*MyCustomClaims); ok && token.Valid {
		return claims.UserID, nil
	}
	return "", err
}
