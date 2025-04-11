package adapters

import (
	"todolist/internal/api/handlers"
	"todolist/internal/models"
	auth_utils "todolist/internal/pkg/authUtils"

	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"

	"github.com/sirupsen/logrus"
)

type IUserRepository interface {
	GetUserByLogin(login string) (*models.User, error)
	GetUserByID(id uint64) (*models.User, error)
	CreateUser(user *models.UserAuth) error
}

type AuthAdapter struct {
	logger       *logrus.Logger
	userRepo     IUserRepository
	key          string
	tokenHandler auth_utils.ITokenHandler
}

func NewAuthService(loggerSrc *logrus.Logger, repo IUserRepository, token auth_utils.ITokenHandler, k string) handlers.AuthProvider {
	return &AuthAdapter{
		logger:       loggerSrc,
		userRepo:     repo,
		tokenHandler: token,
		key:          k,
	}
}

func (serv *AuthAdapter) SignUp(candidate *models.UserAuth) error {
	var err error
	if candidate.Name == "" {
		err = errors.New("Failed to login with empty login")
		serv.logger.Info(err)
		return err
	}

	if candidate.Password == "" {
		err = errors.Errorf("Empty password for user with login %s", candidate.Name)
		serv.logger.Info(err)
		return err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(candidate.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.Wrapf(err, "Error in generating hash for password %s", candidate.Password)
	}

	candidateHashedPasswd := *candidate
	candidateHashedPasswd.Password = string(hash)

	err = serv.userRepo.CreateUser(&candidateHashedPasswd)
	if err != nil {
		err = errors.Wrapf(err, "Failed to create user: %s", candidate.Name)
		serv.logger.Warn(err)
		return err
	}
	serv.logger.Infof("auth svc - successfully signed up as user with login %v", candidate.Name)
	return nil
}

func (serv *AuthAdapter) SignIn(candidate *models.UserAuth) (string, error) {
	var user *models.User
	var err error
	var tokenStr string
	if candidate.Name == "" {
		err = errors.New("Failed to login with empty login")
		serv.logger.Info(err)
		return "", err
	}

	if candidate.Password == "" {
		err = errors.Errorf("Empty password for user with login %s", candidate.Name)
		serv.logger.Info(err)
		return "", err
	}
	user, err = serv.userRepo.GetUserByLogin(candidate.Name)

	if err != nil {
		err = errors.Wrapf(err, "Failed to get user %s", candidate.Name)
		serv.logger.Error(err)
		return "", err
	}
	err = bcrypt.CompareHashAndPassword([]byte(candidate.Password), []byte(user.Password))
	if err != nil {
		err = errors.Wrapf(err, "Invalid password for user %s", candidate.Name)
		serv.logger.Warn(err)
		return "", err
	}
	tokenStr, err = serv.tokenHandler.GenerateToken(*user, serv.key)
	if err != nil {
		err = errors.Wrapf(err, "Failed to generate token for user: %s", candidate.Name)
		serv.logger.Warn(err)
		return "", err
	}
	serv.logger.Infof("auth svc - successfully signed in as user with login %v", candidate.Name)
	return tokenStr, nil
}
