package controller

import (
	"errors"
	"event-booking-api/app/constant"
	"event-booking-api/app/domain/dao"
	_ "event-booking-api/app/domain/dto"
	"event-booking-api/app/pkg"
	"event-booking-api/app/service"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
)

type UserController interface {
	AddUser(c *gin.Context)
	GetAllUser(c *gin.Context)
	GetUserById(c *gin.Context)
	UpdateUserById(c *gin.Context)
	DeleteUserById(c *gin.Context)
	LoginUser(c *gin.Context)
}

type UserControllerImpl struct {
	userSvc service.UserService
}

// AddUser godoc
//
//	@Summary		Create a new user
//	@Description	Create a new user with the provided data
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			user	body		dao.User							true	"User credentials"
//	@Success		201		{object}	dto.ApiResponse[dao.UserResponse]	"Created"
//	@Failure		400		{object}	dto.ApiResponse[any]				"Bad request"
//	@Failure		409		{object}	dto.ApiResponse[any]				"Conflict"
//	@Failure		500		{object}	dto.ApiResponse[any]				"Internal server error"
//	@Router			/users [post]
func (u UserControllerImpl) AddUser(c *gin.Context) {
	defer pkg.PanicHandler(c)

	var request dao.User
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Info("Error parsing request data: ", err)
		pkg.PanicException(constant.InvalidRequest)
	}

	if request.RoleID == 0 {
		request.RoleID = 2
	}

	validate := validator.New()
	if err := validate.Struct(request); err != nil {
		log.Info("Error validating request data: ", err)
		pkg.PanicException(constant.InvalidRequest)
	}

	user, err := u.userSvc.AddUser(request)
	if err != nil {
		var customErr *pkg.CustomError
		if errors.As(err, &customErr) {
			pkg.PanicException(customErr.Type)
		}

		pkg.PanicException(constant.UnknownError)
	}

	response := dao.UserResponse{
		ID:     user.ID,
		Email:  user.Email,
		RoleID: user.RoleID,
	}

	c.JSON(http.StatusCreated, pkg.BuildResponse(constant.Success, response))
}

// GetAllUser godoc
//
//	@Summary		Get all users
//	@Description	Retrieve a list of users. Admin only. Requires JWT authentication.
//	@Tags			users
//	@Produce		json
//	@Success		200	{object}	dto.ApiResponse[[]dao.UserResponse]	"Success"
//	@Failure		401	{object}	dto.ApiResponse[any]				"Unauthorized"
//	@Failure		500	{object}	dto.ApiResponse[any]				"Internal server error"
//	@Router			/users [get]
//	@Security		BearerAuth
func (u UserControllerImpl) GetAllUser(c *gin.Context) {
	defer pkg.PanicHandler(c)

	roleId := c.GetInt("roleId")
	if roleId != 1 {
		log.Info("Access denied. Not an Admin User")
		pkg.PanicException(constant.Unauthorized)
	}

	users, err := u.userSvc.GetAllUser()
	if err != nil {
		pkg.PanicException(constant.UnknownError)
	}

	response := make([]dao.UserResponse, len(users))
	for i, user := range users {
		response[i] = dao.UserResponse{
			ID:     user.ID,
			Email:  user.Email,
			RoleID: user.RoleID,
		}
	}

	c.JSON(http.StatusOK, pkg.BuildResponse(constant.Success, response))
}

// GetUserById godoc
//
//	@Summary		Get user by ID
//	@Description	Retrieve a specific user by its ID. Requires JWT authentication.
//	@Tags			users
//	@Produce		json
//	@Param			id	path		int									true	"User ID"
//	@Success		200	{object}	dto.ApiResponse[dao.UserResponse]	"Success"
//	@Failure		401	{object}	dto.ApiResponse[any]				"Unauthorized"
//	@Failure		404	{object}	dto.ApiResponse[any]				"Not found"
//	@Failure		500	{object}	dto.ApiResponse[any]				"Internal server error"
//	@Router			/users/{id} [get]
//	@Security		BearerAuth
func (u UserControllerImpl) GetUserById(c *gin.Context) {
	defer pkg.PanicHandler(c)

	pathUserId, _ := strconv.Atoi(c.Param("userId"))
	userId := c.GetInt("userId")
	if userId != pathUserId {
		log.Info("Access denied. Not a resource owner")
		pkg.PanicException(constant.Unauthorized)
	}

	user, err := u.userSvc.GetUserById(userId)
	if err != nil {
		var customErr *pkg.CustomError
		if errors.As(err, &customErr) {
			pkg.PanicException(customErr.Type)
		}

		pkg.PanicException(constant.UnknownError)
	}

	response := dao.UserResponse{
		ID:     user.ID,
		Email:  user.Email,
		RoleID: user.RoleID,
	}

	c.JSON(http.StatusOK, pkg.BuildResponse(constant.Success, response))
}

// UpdateUserById godoc
//
//	@Summary		Update user by ID
//	@Description	Update an user with the provided data. Requires JWT authentication.
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int									true	"User ID"
//	@Param			user	body		dao.User							true	"Updated user credentials"
//	@Success		200		{object}	dto.ApiResponse[dao.UserResponse]	"Success"
//	@Failure		400		{object}	dto.ApiResponse[any]				"Bad request"
//	@Failure		401		{object}	dto.ApiResponse[any]				"Unauthorized"
//	@Failure		404		{object}	dto.ApiResponse[any]				"Not found"
//	@Failure		409		{object}	dto.ApiResponse[any]				"Conflict"
//	@Failure		500		{object}	dto.ApiResponse[any]				"Internal server error"
//	@Router			/users/{id} [put]
//	@Security		BearerAuth
func (u UserControllerImpl) UpdateUserById(c *gin.Context) {
	defer pkg.PanicHandler(c)

	pathUserId, _ := strconv.Atoi(c.Param("userId"))
	userId := c.GetInt("userId")
	if userId != pathUserId {
		log.Info("Access denied. Not a resource owner")
		pkg.PanicException(constant.Unauthorized)
	}

	var request dao.User
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Info("Error parsing request data: ", err)
		pkg.PanicException(constant.InvalidRequest)
	}

	if request.Email != "" {
		validate := validator.New()
		if err := validate.Var(request.Email, "email"); err != nil {
			log.Info("Error validating request data: ", err)
			pkg.PanicException(constant.InvalidRequest)
		}
	}

	user, err := u.userSvc.UpdateUserById(request, userId)
	if err != nil {
		var customErr *pkg.CustomError
		if errors.As(err, &customErr) {
			pkg.PanicException(customErr.Type)
		}

		pkg.PanicException(constant.UnknownError)
	}

	response := dao.UserResponse{
		ID:     user.ID,
		Email:  user.Email,
		RoleID: user.RoleID,
	}

	c.JSON(http.StatusOK, pkg.BuildResponse(constant.Success, response))
}

// DeleteUserById godoc
//
//	@Summary		Delete user by ID
//	@Description	Delete a specific user by its ID. Requires JWT authentication.
//	@Tags			users
//	@Produce		json
//	@Param			id	path		int						true	"User ID"
//	@Success		200	{object}	dto.ApiResponse[any]	"Success"
//	@Failure		401	{object}	dto.ApiResponse[any]	"Unauthorized"
//	@Failure		500	{object}	dto.ApiResponse[any]	"Internal server error"
//	@Router			/users/{id} [delete]
//	@Security		BearerAuth
func (u UserControllerImpl) DeleteUserById(c *gin.Context) {
	defer pkg.PanicHandler(c)

	pathUserId, _ := strconv.Atoi(c.Param("userId"))
	userId := c.GetInt("userId")
	if userId != pathUserId {
		log.Info("Access denied. Not a resource owner")
		pkg.PanicException(constant.Unauthorized)
	}

	err := u.userSvc.DeleteUserById(userId)
	if err != nil {
		pkg.PanicException(constant.UnknownError)
	}

	c.JSON(http.StatusOK, pkg.BuildResponse(constant.Success, pkg.Null()))
}

// LoginUser godoc
//
//	@Summary		Authenticate a user
//	@Description	Authenticate a user with the provided credentials and return a JWT token
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			user	body		dao.User				true	"User credentials"
//	@Success		200		{object}	dto.ApiResponse[string]	"Success"
//	@Failure		400		{object}	dto.ApiResponse[any]	"Bad request"
//	@Failure		401		{object}	dto.ApiResponse[any]	"Unauthorized"
//	@Failure		500		{object}	dto.ApiResponse[any]	"Internal server error"
//	@Router			/users/login [post]
func (u UserControllerImpl) LoginUser(c *gin.Context) {
	defer pkg.PanicHandler(c)

	var request dao.User
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Info("Error parsing request data: ", err)
		pkg.PanicException(constant.InvalidRequest)
	}

	validate := validator.New()
	if err := validate.Struct(request); err != nil {
		log.Info("Error validating request data: ", err)
		pkg.PanicException(constant.InvalidRequest)
	}

	token, err := u.userSvc.LoginUser(request)
	if err != nil {
		var customErr *pkg.CustomError
		if errors.As(err, &customErr) {
			pkg.PanicException(customErr.Type)
		}

		pkg.PanicException(constant.UnknownError)
	}

	c.JSON(http.StatusOK, pkg.BuildResponse(constant.Success, token))
}

func UserControllerInit(userService service.UserService) *UserControllerImpl {
	return &UserControllerImpl{
		userSvc: userService,
	}
}
