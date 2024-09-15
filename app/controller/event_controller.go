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

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"
)

type EventController interface {
	AddEvent(c *gin.Context)
	GetAllEvent(c *gin.Context)
	GetEventById(c *gin.Context)
	UpdateEventById(c *gin.Context)
	DeleteEventById(c *gin.Context)
	RegisterUserForEvent(c *gin.Context)
	UnregisterUserForEvent(c *gin.Context)
	GetAttendeesEmailById(c *gin.Context)
}

type EventControllerImpl struct {
	eventSvc    service.EventService
	registerSvc service.RegisterService
}

// AddEvent godoc
//
//	@Summary		Create a new event
//	@Description	Create a new event with the provided data. Requires JWT authentication.
//	@Tags			events
//	@Accept			json
//	@Produce		json
//	@Param			event	body		dao.Event							true	"Event data"
//	@Success		201		{object}	dto.ApiResponse[dao.EventResponse]	"Created"
//	@Failure		400		{object}	dto.ApiResponse[any]				"Bad request"
//	@Failure		401		{object}	dto.ApiResponse[any]				"Unauthorized"
//	@Failure		500		{object}	dto.ApiResponse[any]				"Internal server error"
//	@Router			/events [post]
//	@Security		BearerAuth
func (e EventControllerImpl) AddEvent(c *gin.Context) {
	defer pkg.PanicHandler(c)

	var request dao.Event
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Info("Error parsing request data: ", err)
		pkg.PanicException(constant.InvalidRequest)
	}

	validate := validator.New()
	if err := validate.StructExcept(request, "User"); err != nil {
		log.Info("Error validating request data: ", err)
		pkg.PanicException(constant.InvalidRequest)
	}

	request.UserID = c.GetInt("userId")

	event, err := e.eventSvc.AddEvent(request)
	if err != nil {
		pkg.PanicException(constant.UnknownError)
	}

	response := dao.EventResponse{
		ID:          event.ID,
		Name:        event.Name,
		Description: event.Description,
		Location:    event.Location,
		EventTime:   event.EventTime,
		UserID:      event.UserID,
	}

	c.JSON(http.StatusCreated, pkg.BuildResponse(constant.Success, response))
}

// GetAllEvent godoc
//
//	@Summary		Get all events
//	@Description	Retrieve a list of events
//	@Tags			events
//	@Produce		json
//	@Success		200	{object}	dto.ApiResponse[[]dao.EventResponse]	"Success"
//	@Failure		500	{object}	dto.ApiResponse[any]					"Internal server error"
//	@Router			/events [get]
func (e EventControllerImpl) GetAllEvent(c *gin.Context) {
	defer pkg.PanicHandler(c)

	events, err := e.eventSvc.GetAllEvent()
	if err != nil {
		pkg.PanicException(constant.UnknownError)
	}

	response := make([]dao.EventResponse, len(events))
	for i, event := range events {
		response[i] = dao.EventResponse{
			ID:          event.ID,
			Name:        event.Name,
			Description: event.Description,
			Location:    event.Location,
			EventTime:   event.EventTime,
			UserID:      event.UserID,
		}
	}

	c.JSON(http.StatusOK, pkg.BuildResponse(constant.Success, response))
}

// GetEventById godoc
//
//	@Summary		Get event by ID
//	@Description	Retrieve a specific event by its ID
//	@Tags			events
//	@Produce		json
//	@Param			id	path		int									true	"Event ID"
//	@Success		200	{object}	dto.ApiResponse[dao.EventResponse]	"Success"
//	@Failure		404	{object}	dto.ApiResponse[any]				"Not found"
//	@Failure		500	{object}	dto.ApiResponse[any]				"Internal server error"
//	@Router			/events/{id} [get]
func (e EventControllerImpl) GetEventById(c *gin.Context) {
	defer pkg.PanicHandler(c)

	eventId, _ := strconv.Atoi(c.Param("eventId"))

	event, err := e.eventSvc.GetEventById(eventId)
	if err != nil {
		var customErr *pkg.CustomError
		if errors.As(err, &customErr) {
			pkg.PanicException(customErr.Type)
		}

		pkg.PanicException(constant.UnknownError)
	}

	response := dao.EventResponse{
		ID:          event.ID,
		Name:        event.Name,
		Description: event.Description,
		Location:    event.Location,
		EventTime:   event.EventTime,
		UserID:      event.UserID,
	}

	c.JSON(http.StatusOK, pkg.BuildResponse(constant.Success, response))
}

// UpdateEventById godoc
//
//	@Summary		Update event by ID
//	@Description	Update an event with the provided data. Requires JWT authentication.
//	@Tags			events
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int									true	"Event ID"
//	@Param			event	body		dao.Event							true	"Updated event data"
//	@Success		200		{object}	dto.ApiResponse[dao.EventResponse]	"Success"
//	@Failure		400		{object}	dto.ApiResponse[any]				"Bad request"
//	@Failure		401		{object}	dto.ApiResponse[any]				"Unauthorized"
//	@Failure		404		{object}	dto.ApiResponse[any]				"Not found"
//	@Failure		500		{object}	dto.ApiResponse[any]				"Internal server error"
//	@Router			/events/{id} [put]
//	@Security		BearerAuth
func (e EventControllerImpl) UpdateEventById(c *gin.Context) {
	defer pkg.PanicHandler(c)

	eventId, _ := strconv.Atoi(c.Param("eventId"))
	userId := c.GetInt("userId")

	var request dao.Event
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Info("Error parsing request data: ", err)
		pkg.PanicException(constant.InvalidRequest)
	}

	event, err := e.eventSvc.UpdateEventById(request, eventId, userId)
	if err != nil {
		var customErr *pkg.CustomError
		if errors.As(err, &customErr) {
			pkg.PanicException(customErr.Type)
		}

		pkg.PanicException(constant.UnknownError)
	}

	response := dao.EventResponse{
		ID:          event.ID,
		Name:        event.Name,
		Description: event.Description,
		Location:    event.Location,
		EventTime:   event.EventTime,
		UserID:      event.UserID,
	}

	c.JSON(http.StatusOK, pkg.BuildResponse(constant.Success, response))
}

// DeleteEventById godoc
//
//	@Summary		Delete event by ID
//	@Description	Delete a specific event by its ID. Requires JWT authentication.
//	@Tags			events
//	@Produce		json
//	@Param			id	path		int						true	"Event ID"
//	@Success		200	{object}	dto.ApiResponse[any]	"Success"
//	@Failure		401	{object}	dto.ApiResponse[any]	"Unauthorized"
//	@Failure		404	{object}	dto.ApiResponse[any]	"Not found"
//	@Failure		500	{object}	dto.ApiResponse[any]	"Internal server error"
//	@Router			/events/{id} [delete]
//	@Security		BearerAuth
func (e EventControllerImpl) DeleteEventById(c *gin.Context) {
	defer pkg.PanicHandler(c)

	eventId, _ := strconv.Atoi(c.Param("eventId"))
	userId := c.GetInt("userId")

	err := e.eventSvc.DeleteEventById(eventId, userId)
	if err != nil {
		var customErr *pkg.CustomError
		if errors.As(err, &customErr) {
			pkg.PanicException(customErr.Type)
		}

		pkg.PanicException(constant.UnknownError)
	}

	c.JSON(http.StatusOK, pkg.BuildResponse(constant.Success, pkg.Null()))
}

// RegisterUserForEvent godoc
//
//	@Summary		Register user for a specific event
//	@Description	Register user for a specific event by its ID. Requires JWT authentication.
//	@Tags			events
//	@Produce		json
//	@Param			id	path		int						true	"Event ID"
//	@Success		201	{object}	dto.ApiResponse[any]	"Created"
//	@Failure		401	{object}	dto.ApiResponse[any]	"Unauthorized"
//	@Failure		404	{object}	dto.ApiResponse[any]	"Not found"
//	@Failure		409	{object}	dto.ApiResponse[any]	"Conflict"
//	@Failure		500	{object}	dto.ApiResponse[any]	"Internal server error"
//	@Router			/events/{id}/register [post]
//	@Security		BearerAuth
func (e EventControllerImpl) RegisterUserForEvent(c *gin.Context) {
	defer pkg.PanicHandler(c)

	eventId, _ := strconv.Atoi(c.Param("eventId"))
	userId := c.GetInt("userId")

	err := e.registerSvc.RegisterUserForEvent(eventId, userId)
	if err != nil {
		var customErr *pkg.CustomError
		if errors.As(err, &customErr) {
			pkg.PanicException(customErr.Type)
		}

		pkg.PanicException(constant.UnknownError)
	}

	c.JSON(http.StatusCreated, pkg.BuildResponse(constant.Success, pkg.Null()))
}

// UnregisterUserForEvent godoc
//
//	@Summary		Unregister user for a specific event
//	@Description	Unregister user for a specific event by its ID. Requires JWT authentication.
//	@Tags			events
//	@Produce		json
//	@Param			id	path		int						true	"Event ID"
//	@Success		200	{object}	dto.ApiResponse[any]	"Success"
//	@Failure		401	{object}	dto.ApiResponse[any]	"Unauthorized"
//	@Failure		500	{object}	dto.ApiResponse[any]	"Internal server error"
//	@Router			/events/{id}/register [delete]
//	@Security		BearerAuth
func (e EventControllerImpl) UnregisterUserForEvent(c *gin.Context) {
	defer pkg.PanicHandler(c)

	eventId, _ := strconv.Atoi(c.Param("eventId"))
	userId := c.GetInt("userId")

	err := e.registerSvc.UnregisterUserForEvent(eventId, userId)
	if err != nil {
		pkg.PanicException(constant.UnknownError)
	}

	c.JSON(http.StatusOK, pkg.BuildResponse(constant.Success, pkg.Null()))
}

// GetAttendeesEmailById godoc
//
//	@Summary		Get all event attendees email
//	@Description	Retrieve a list of event attendees email. Requires JWT authentication.
//	@Tags			events
//	@Produce		json
//	@Param			id	path		int							true	"Event ID"
//	@Success		200	{object}	dto.ApiResponse[[]string]	"Success"
//	@Failure		401	{object}	dto.ApiResponse[any]		"Unauthorized"
//	@Failure		404	{object}	dto.ApiResponse[any]		"Not found"
//	@Failure		500	{object}	dto.ApiResponse[any]		"Internal server error"
//	@Router			/events/{id}/attendees [get]
//	@Security		BearerAuth
func (e EventControllerImpl) GetAttendeesEmailById(c *gin.Context) {
	defer pkg.PanicHandler(c)

	eventId, _ := strconv.Atoi(c.Param("eventId"))
	userId := c.GetInt("userId")

	emails, err := e.registerSvc.GetAttendeesEmailById(eventId, userId)
	if err != nil {
		var customErr *pkg.CustomError
		if errors.As(err, &customErr) {
			pkg.PanicException(customErr.Type)
		}

		pkg.PanicException(constant.UnknownError)
	}

	c.JSON(http.StatusOK, pkg.BuildResponse(constant.Success, emails))
}

func EventControllerInit(eventService service.EventService,
	registerService service.RegisterService) *EventControllerImpl {
	return &EventControllerImpl{
		eventSvc:    eventService,
		registerSvc: registerService,
	}
}
