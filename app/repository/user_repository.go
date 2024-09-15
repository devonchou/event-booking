package repository

import (
	"errors"
	"event-booking-api/app/domain/dao"
	"event-booking-api/app/pkg"

	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserRepository interface {
	Save(request *dao.User) (dao.User, error)
	FindAllUser() ([]dao.User, error)
	FindUserById(id int) (dao.User, error)
	DeleteUserById(id int) error
	VerifyUser(request dao.User) (dao.User, error)
}

type UserRepositoryImpl struct {
	db *gorm.DB
}

// Save stores the user to the database.
// It returns the saved dao.User and an error, if any.
func (u UserRepositoryImpl) Save(request *dao.User) (dao.User, error) {
	err := u.db.Save(request).Error
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			log.Info("Error saving user: ", err)
			return dao.User{}, pkg.NewConflictError("Email already used", err)
		}

		log.Error("Error saving user: ", err)
		return dao.User{}, err
	}

	return *request, nil
}

// FindAllUser retrieves all users from the database.
// It returns a slice of dao.User and an error, if any.
func (u UserRepositoryImpl) FindAllUser() ([]dao.User, error) {
	var users []dao.User

	err := u.db.Select("id, email, role_id").Find(&users).Error
	if err != nil {
		log.Error("Error finding all users: ", err)
		return nil, err
	}

	return users, nil
}

// FindUserById retrieves a user by the given ID from the database.
// It returns the dao.User and an error, if any.
func (u UserRepositoryImpl) FindUserById(id int) (dao.User, error) {
	user := dao.User{ID: id}

	err := u.db.First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Info("Error finding user by ID: ", err)
			return dao.User{}, pkg.NewNotFoundError("User not found", err)
		}

		log.Error("Error finding user by ID: ", err)
		return dao.User{}, err
	}

	return user, nil
}

// DeleteUserById deletes the user by the given ID from the database.
// It returns an error if the deletion fails.
func (u UserRepositoryImpl) DeleteUserById(id int) error {
	err := u.db.Delete(&dao.User{}, id).Error
	if err != nil {
		log.Error("Error deleting user: ", err)
		return err
	}

	return nil
}

// VerifyUser verifies the user's credentials by checking the provided email and password.
// It returns the found dao.User if credentials are valid, or an error otherwise.
func (u UserRepositoryImpl) VerifyUser(request dao.User) (dao.User, error) {
	var foundUser dao.User

	err := u.db.Where("email = ?", request.Email).First(&foundUser).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Info("Error verifying user: ", err)
			return dao.User{}, pkg.NewUnauthorizedError("Invalid credentials", err)
		}

		log.Error("Error verifying user: ", err)
		return dao.User{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(request.Password))
	if err != nil {
		log.Info("Error verifying user: ", err)
		return dao.User{}, pkg.NewUnauthorizedError("Invalid credentials", err)
	}

	return foundUser, nil
}

func UserRepositoryInit(db *gorm.DB) *UserRepositoryImpl {
	if err := db.AutoMigrate(&dao.User{}); err != nil {
		log.Fatal("Error AutoMigrating User: ", err)
	}

	return &UserRepositoryImpl{
		db: db,
	}
}
