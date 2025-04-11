package repository

import (
	"todolist/internal/adapters"
	"todolist/internal/models"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type User struct {
	ID       uuid.UUID `gorm:"primaryKey;column:id;type:uuid;default:gen_random_uuid()"`
	Name     string    `gorm:"unique;column:name"`
	Password string    `gorm:"column:password"`
}

func ToDaUser(user models.UserAuth) User {
	return User{
		Name:     user.Name,
		Password: user.Password,
	}
}

func FromDaUser(user User) models.User {
	return models.User{
		ID:       user.ID,
		Name:     user.Name,
		Password: user.Password,
	}
}

type UserRepositoryAdapter struct {
	db *gorm.DB
}

func NewUserRepositoryAdapter(srcDB *gorm.DB) adapters.IUserRepository {
	return &UserRepositoryAdapter{
		db: srcDB,
	}
}

func (repo *UserRepositoryAdapter) GetUserByID(id uuid.UUID) (*models.User, error) {
	var userDA User
	userDA.ID = id
	tx := repo.db.First(&userDA)

	if tx.Error == gorm.ErrRecordNotFound {
		return nil, models.ErrUserNotFound
	}

	if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "Error getting user by ID")
	}
	user := FromDaUser(userDA)
	return &user, nil
}

func (repo *UserRepositoryAdapter) GetUserByName(name string) (*models.User, error) {
	var userDA User
	tx := repo.db.Where("name = ?", name).First(&userDA)

	if tx.Error == gorm.ErrRecordNotFound {
		return nil, models.ErrUserNotFound
	}

	if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "Error getting user by ID")
	}
	user := FromDaUser(userDA)
	return &user, nil
}

func (repo *UserRepositoryAdapter) CreateUser(user *models.UserAuth) error {
	userDa := ToDaUser(*user)
	tx := repo.db.Create(&userDa)

	if tx.Error != nil {
		return errors.Wrap(tx.Error, "error in creating user")
	}
	return nil
}

func (repo *UserRepositoryAdapter) CheckTaskOwnership(userID uuid.UUID, taskID uuid.UUID) (bool, error) {
	var isOwned bool

	err := repo.db.
		Raw(
			"SELECT EXISTS(SELECT 1 FROM task WHERE id = ? AND user_id = ?)",
			taskID,
			userID,
		).
		Scan(&isOwned).
		Error

	if err != nil {
		return false, errors.Wrap(err, "failed to verify task ownership")
	}

	return isOwned, nil
}

func (repo *UserRepositoryAdapter) CheckCategoriesOwnership(userID uuid.UUID, categories []uuid.UUID) (bool, error) {
	if len(categories) == 0 {
		return true, nil
	}

	var allOwned bool
	err := repo.db.Raw(`
        SELECT NOT EXISTS (
            SELECT 1 FROM categories 
            WHERE id IN (?) 
            AND user_id != ?
        )`,
		pq.Array(categories),
		userID,
	).Scan(&allOwned).Error

	return allOwned, errors.Wrap(err, "failed raw ownership check")
}
