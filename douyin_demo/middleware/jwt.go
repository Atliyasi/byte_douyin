package middleware

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

var jwtKey = []byte("chx")

type MyClaims struct {
	UserId int `json:"user_id"`
	jwt.StandardClaims
}

func GenToken(id int) (string, error) {
	claim := MyClaims{
		UserId: id,
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix(),
			ExpiresAt: time.Now().Unix() + 7*24*60*60,
			Issuer:    "chx",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	//加密
	ss, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	return ss, nil

}

func ParseToken(ss string) (*jwt.Token, *MyClaims, error) {
	claims := &MyClaims{}
	token, err := jwt.ParseWithClaims(ss, claims, func(t *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	return token, claims, err
}
