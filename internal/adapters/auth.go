package adapters

import (
	"todolist/internal/api/handlers"
	"todolist/internal/models"
	auth_utils "todolist/internal/pkg/authUtils"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"

	"github.com/sirupsen/logrus"
)

type IUserRepository interface {
	GetUserByName(name string) (*models.User, error)
	GetUserByID(id uuid.UUID) (*models.User, error)
	CreateUser(user *models.UserAuth) error
	CheckTaskOwnership(userID uuid.UUID, taskID uuid.UUID) (bool, error)
	CheckCategoriesOwnership(userID uuid.UUID, categories []uuid.UUID) (bool, error)
}

type UserAdapter struct {
	logger       *logrus.Logger
	userRepo     IUserRepository
	key          string
	tokenHandler auth_utils.ITokenHandler
}

func NewAuthService(loggerSrc *logrus.Logger, repo IUserRepository, token auth_utils.ITokenHandler, k string) handlers.AuthProvider {
	return &UserAdapter{
		logger:       loggerSrc,
		userRepo:     repo,
		tokenHandler: token,
		key:          k,
	}
}

func (serv *UserAdapter) SignUp(candidate *models.UserAuth) error {
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

func (serv *UserAdapter) SignIn(candidate *models.UserAuth) (string, error) {
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
	user, err = serv.userRepo.GetUserByName(candidate.Name)

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

func (serv *UserAdapter) CheckTaskOwnership(userID uuid.UUID, taskID uuid.UUID) (bool, error) {
	isTaskOwned, err := serv.CheckTaskOwnership(userID, taskID)
	if err != nil {
		serv.logger.Infof("Error in checking task ownership for task for user %v: %v", userID, err)
		return false, errors.Wrap(err, "Error in checking task ownership")
	}
	return isTaskOwned, nil
}

func (serv *UserAdapter) CheckCategoriesOwnership(userID uuid.UUID, categories []uuid.UUID) (bool, error) {
	areCategoriesOwned, err := serv.CheckCategoriesOwnership(userID, categories)
	if err != nil {
		serv.logger.Infof("Error in checking categories ownership for task for user %v: %v", userID, err)
		return false, errors.Wrap(err, "Error in checking task ownership")
	}
	return areCategoriesOwned, nil
}
