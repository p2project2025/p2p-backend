package midleware

import (
	"errors"
	"p2p/config"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	Email  string
	UserID string
	Role   string
	jwt.RegisteredClaims
}

func GenerateJWT(email, userID, role string) (string, error) {
	expireTime := time.Now().Add(60 * time.Minute)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &Claims{
		Email:  email,
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireTime),
		},
	})

	tokenString, err := token.SignedString([]byte(config.Cfg.JWTSecret))

	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ValidateToken(tokenString string) (Claims, error) {
	claims := Claims{}
	token, err := jwt.ParseWithClaims(tokenString, &claims,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(config.Cfg.JWTSecret), nil
		},
	)
	if err != nil || !token.Valid {
		return claims, errors.New("not valid token")
	}
	if time.Now().Unix() > claims.ExpiresAt.Unix() {
		return claims, errors.New("token expired re-login")
	}
	return claims, nil
}
