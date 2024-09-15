package service

import (
	"event-booking-api/app/domain/dao"
	"event-booking-api/app/repository"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	AddUser(request dao.User) (dao.User, error)
	GetAllUser() ([]dao.User, error)
	GetUserById(userId int) (dao.User, error)
	UpdateUserById(request dao.User, userId int) (dao.User, error)
	DeleteUserById(userId int) error
	LoginUser(request dao.User) (string, error)
}

type UserServiceImpl struct {
	userRepo repository.UserRepository
}

// AddUser adds a new user to the repository by hashing the provided password.
// It returns the added dao.User and an error if the operation fails.
func (u UserServiceImpl) AddUser(request dao.User) (dao.User, error) {
	log.Info("Start to execute add user")

	hash, _ := bcrypt.GenerateFromPassword([]byte(request.Password), 14)
	request.Password = string(hash)

	user, err := u.userRepo.Save(&request)
	if err != nil {
		return dao.User{}, err
	}

	return user, nil
}

// GetAllUser retrieves all users from the repository.
// It returns a slice of dao.User and an error if the operation fails.
func (u UserServiceImpl) GetAllUser() ([]dao.User, error) {
	log.Info("Start to execute get all user")

	users, err := u.userRepo.FindAllUser()
	if err != nil {
		return nil, err
	}

	return users, nil
}

// GetUserById retrieves a user from the repository by their ID.
// It returns the dao.User with the specified ID and an error if the operation fails.
func (u UserServiceImpl) GetUserById(userId int) (dao.User, error) {
	log.Info("Start to execute get user by id")

	user, err := u.userRepo.FindUserById(userId)
	if err != nil {
		return dao.User{}, err
	}

	return user, nil
}

// UpdateUserById updates a user's details by their ID.
// It modifies the user's email, password if provided in the request.
// It returns the updated dao.User and an error if the operation fails.
func (u UserServiceImpl) UpdateUserById(request dao.User, userId int) (dao.User, error) {
	log.Info("Start to execute update user by id")

	user, err := u.userRepo.FindUserById(userId)
	if err != nil {
		return dao.User{}, err
	}

	if request.Email != "" {
		user.Email = request.Email
	}
	if request.Password != "" {
		hash, _ := bcrypt.GenerateFromPassword([]byte(request.Password), 14)
		user.Password = string(hash)
	}

	user, err = u.userRepo.Save(&user)
	if err != nil {
		return dao.User{}, err
	}

	return user, nil
}

// DeleteUserById removes a user from the repository by their ID.
// It returns an error if the operation fails.
func (u UserServiceImpl) DeleteUserById(userId int) error {
	log.Info("Start to execute delete user by id")

	err := u.userRepo.DeleteUserById(userId)
	if err != nil {
		return err
	}

	return nil
}

// LoginUser verifies user credentials and generates a JWT token if the credentials are valid.
// It returns the JWT token and an error if the operation fails.
func (u UserServiceImpl) LoginUser(request dao.User) (string, error) {
	log.Info("Start to verify user login credentials")

	foundUser, err := u.userRepo.VerifyUser(request)
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": foundUser.ID,
		"email":   foundUser.Email,
		"role_id": foundUser.RoleID,
		"exp":     time.Now().Add(time.Hour * 2).Unix(),
	})

	return token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
}

func UserServiceInit(userRepository repository.UserRepository) *UserServiceImpl {
	return &UserServiceImpl{
		userRepo: userRepository,
	}
}
