package auth

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4/middleware"
)

const TTL = time.Hour * 24 * 365 * 2

var SIGNED_KEY = []byte("verysecret")

var IsLoggedIn = middleware.JWTWithConfig(middleware.JWTConfig{
	SigningKey: SIGNED_KEY,
})

type JWTToken struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func GenerateJWT(fields map[string]interface{}) (*JWTToken, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["sub"] = 1
	claims["exp"] = time.Now().Add(TTL).Unix()

	for k, v := range fields {
		claims[k] = v
	}

	t, err := token.SignedString(SIGNED_KEY)
	if err != nil {
		return nil, err
	}

	refreshToken := jwt.New(jwt.SigningMethodHS256)
	rtClaims := refreshToken.Claims.(jwt.MapClaims)
	rtClaims["sub"] = 1
	rtClaims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	rt, err := refreshToken.SignedString(SIGNED_KEY)
	if err != nil {
		return nil, err
	}

	return &JWTToken{
		AccessToken:  t,
		RefreshToken: rt,
	}, nil
}
