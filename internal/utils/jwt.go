package utils

import (
	"go-media-stream/internal/domain"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
)

var jwtkey = []byte("very-secret-key") // TODO: JWT secret

func GenerateJWT(userId string) (string, error) {
	payload := jwt.MapClaims{
		"sub": userId,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	t, err := token.SignedString(jwtkey)
	if err != nil {
		return "", err
	}
	return t, nil
}

func VerifyJWT(token string) (string, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, domain.ErrUserNotFound
		}
		return []byte(jwtkey), nil
	}
	payload := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(token, &payload, keyFunc)
	if err != nil {
		logrus.Error(err)
		return "", err
	}
	userId, err := payload.GetSubject()
	if err != nil {
		logrus.Error(err)
		return "", err
	}
	return userId, nil
}
