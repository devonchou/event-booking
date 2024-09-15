package config

import (
	"event-booking-api/app/controller"
	"event-booking-api/app/repository"
	"event-booking-api/app/service"
)

type Initialization struct {
	roleRepo     repository.RoleRepository
	userRepo     repository.UserRepository
	eventRepo    repository.EventRepository
	registerRepo repository.RegisterRepository
	userSvc      service.UserService
	eventSvc     service.EventService
	registerSvc  service.RegisterService
	UserCtrl     controller.UserController
	EventCtrl    controller.EventController
}

func NewInitialization(roleRepo repository.RoleRepository,
	userRepo repository.UserRepository,
	eventRepo repository.EventRepository,
	registerRepo repository.RegisterRepository,
	userSvc service.UserService,
	eventSvc service.EventService,
	registerSvc service.RegisterService,
	userCtrl controller.UserController,
	eventCtrl controller.EventController,
) *Initialization {
	return &Initialization{
		roleRepo:     roleRepo,
		userRepo:     userRepo,
		eventRepo:    eventRepo,
		registerRepo: registerRepo,
		userSvc:      userSvc,
		eventSvc:     eventSvc,
		registerSvc:  registerSvc,
		UserCtrl:     userCtrl,
		EventCtrl:    eventCtrl,
	}
}
