package service

import (
	"event-booking-api/app/domain/dao"
	"event-booking-api/app/pkg"
	"event-booking-api/app/repository"

	log "github.com/sirupsen/logrus"
)

type RegisterService interface {
	RegisterUserForEvent(eventId, userId int) error
	UnregisterUserForEvent(eventId, userId int) error
	GetAttendeesEmailById(eventId, userId int) ([]string, error)
}

type RegisterServiceImpl struct {
	eventRepo    repository.EventRepository
	registerRepo repository.RegisterRepository
}

// RegisterUserForEvent registers a user for a specific event to the repository.
// It returns an error if the operation fails.
func (r RegisterServiceImpl) RegisterUserForEvent(eventId, userId int) error {
	log.Info("Start to execute register user for event")

	_, err := r.eventRepo.FindEventById(eventId)
	if err != nil {
		return err
	}

	register := dao.Register{
		EventID: eventId,
		UserID:  userId,
	}

	err = r.registerRepo.Save(&register)
	if err != nil {
		return err
	}

	return nil
}

// UnregisterUserForEvent unregisters a user for a specific event from the repository.
// It returns an error if the operation fails.
func (r RegisterServiceImpl) UnregisterUserForEvent(eventId, userId int) error {
	log.Info("Start to execute unregister user for event")

	err := r.registerRepo.Delete(eventId, userId)
	if err != nil {
		return err
	}

	return nil
}

// GetAttendeesEmailByEventID retrieves the email addresses of all the event attendees from the repository.
// Access is restricted to the resource owner.
// It returns a slice of emails and an error if the operation fails.
func (r RegisterServiceImpl) GetAttendeesEmailById(eventId, userId int) ([]string, error) {
	log.Info("Start to execute get attendees email by id")

	event, err := r.eventRepo.FindEventById(eventId)
	if err != nil {
		return nil, err
	}

	if event.UserID != userId {
		log.Info("Access denied. Not a resource owner")
		return nil, pkg.NewUnauthorizedError("Unauthorized", nil)
	}

	emails, err := r.registerRepo.FindAttendeesEmailById(eventId)
	if err != nil {
		return nil, err
	}

	return emails, nil
}

func RegisterServiceInit(eventRepository repository.EventRepository,
	registerRepository repository.RegisterRepository) *RegisterServiceImpl {
	return &RegisterServiceImpl{
		eventRepo:    eventRepository,
		registerRepo: registerRepository,
	}
}
