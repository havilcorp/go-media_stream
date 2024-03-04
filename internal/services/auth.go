package services

import (
	"context"
	"errors"
	"strconv"

	"go-media-stream/internal/domain"
	"go-media-stream/internal/utils"

	"github.com/sirupsen/logrus"
)

type Auther interface {
	CreateUser(ctx context.Context, user domain.User) (int64, error)
	GetUserByLogin(ctx context.Context, login string, password string) (domain.User, error)
	GetUserById(ctx context.Context, id int64) (domain.User, error)
}

type AuthService struct {
	auth Auther
}

func NewAuthServices(auth Auther) *AuthService {
	return &AuthService{
		auth: auth,
	}
}

func (s *AuthService) SignUp(ctx context.Context, login string, password string) (string, error) {
	var (
		user domain.User
		err  error
	)
	hashPassword := utils.GetMD5(password)
	user, err = s.auth.GetUserByLogin(ctx, login, hashPassword)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			lastIndex, err := s.auth.CreateUser(ctx, domain.User{Login: login, Password: hashPassword})
			if err != nil {
				return "", err
			}
			user, err = s.auth.GetUserById(ctx, lastIndex)
			if err != nil {
				return "", err
			}
		} else if errors.Is(err, domain.ErrWrongPassword) {
			return "", domain.ErrWrongPassword
		} else {
			logrus.Error(err)
			return "", err
		}
	}
	token, err := utils.GenerateJWT(strconv.Itoa(user.ID))
	if err != nil {
		return "", err
	}
	return token, nil
}
