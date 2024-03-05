package services

import (
	"context"
	"errors"
	"strconv"

	"go-media-stream/internal/domain"
	"go-media-stream/internal/utils"

	"github.com/sirupsen/logrus"
)

//go:generate mockery --name AuthProvider
type AuthProvider interface {
	CreateUser(ctx context.Context, user *domain.User) (int64, error)
	GetUserByLogin(ctx context.Context, login string, password string) (*domain.User, error)
	GetUserById(ctx context.Context, id int64) (*domain.User, error)
}

type AuthService struct {
	authProvider AuthProvider
}

func NewAuthServices(authProvider AuthProvider) *AuthService {
	return &AuthService{
		authProvider: authProvider,
	}
}

func (s *AuthService) SignUp(ctx context.Context, login string, password string) (string, error) {
	var (
		user *domain.User
		err  error
	)
	hashPassword := utils.GetMD5(password)
	user, err = s.authProvider.GetUserByLogin(ctx, login, hashPassword)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			lastIndex, err := s.authProvider.CreateUser(ctx, &domain.User{Login: login, Password: hashPassword})
			if err != nil {
				return "", err
			}
			user, err = s.authProvider.GetUserById(ctx, lastIndex)
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
	token, err := utils.GenerateJWT(strconv.FormatInt(user.ID, 10))
	if err != nil {
		return "", err
	}
	return token, nil
}
