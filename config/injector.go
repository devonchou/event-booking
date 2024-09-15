// go:build wireinject
//go:build wireinject
// +build wireinject

package config

import (
	"event-booking-api/app/controller"
	"event-booking-api/app/repository"
	"event-booking-api/app/service"

	"github.com/google/wire"
)

var db = wire.NewSet(ConnectToDB)

var roleRepoSet = wire.NewSet(repository.RoleRepositoryInit,
	wire.Bind(new(repository.RoleRepository), new(*repository.RoleRepositoryImpl)),
)

var userRepoSet = wire.NewSet(repository.UserRepositoryInit,
	wire.Bind(new(repository.UserRepository), new(*repository.UserRepositoryImpl)),
)

var eventRepoSet = wire.NewSet(repository.EventRepositoryInit,
	wire.Bind(new(repository.EventRepository), new(*repository.EventRepositoryImpl)),
)

var registerRepoSet = wire.NewSet(repository.RegisterRepositoryInit,
	wire.Bind(new(repository.RegisterRepository), new(*repository.RegisterRepositoryImpl)),
)

var userSvcSet = wire.NewSet(service.UserServiceInit,
	wire.Bind(new(service.UserService), new(*service.UserServiceImpl)),
)

var eventSvcSet = wire.NewSet(service.EventServiceInit,
	wire.Bind(new(service.EventService), new(*service.EventServiceImpl)),
)

var registerSvcSet = wire.NewSet(service.RegisterServiceInit,
	wire.Bind(new(service.RegisterService), new(*service.RegisterServiceImpl)),
)

var userCtrlSet = wire.NewSet(controller.UserControllerInit,
	wire.Bind(new(controller.UserController), new(*controller.UserControllerImpl)),
)

var eventCtrlSet = wire.NewSet(controller.EventControllerInit,
	wire.Bind(new(controller.EventController), new(*controller.EventControllerImpl)),
)

func Init() *Initialization {
	wire.Build(
		NewInitialization,
		db,
		roleRepoSet,
		userRepoSet,
		eventRepoSet,
		registerRepoSet,
		userSvcSet,
		eventSvcSet,
		registerSvcSet,
		userCtrlSet,
		eventCtrlSet,
	)
	return nil
}
