package service

import (
	"event-booking-api/app/domain/dao"
	"event-booking-api/app/pkg"
	"event-booking-api/app/repository"

	log "github.com/sirupsen/logrus"
)

type EventService interface {
	AddEvent(request dao.Event) (dao.Event, error)
	GetAllEvent() ([]dao.Event, error)
	GetEventById(eventId int) (dao.Event, error)
	UpdateEventById(request dao.Event, eventId, userId int) (dao.Event, error)
	DeleteEventById(eventId, userId int) error
}

type EventServiceImpl struct {
	eventRepo repository.EventRepository
}

// AddEvent adds a new event to the repository.
// It returns the added dao.Event and an error if the operation fails.
func (e EventServiceImpl) AddEvent(request dao.Event) (dao.Event, error) {
	log.Info("Start to execute add event")

	event, err := e.eventRepo.Save(&request)
	if err != nil {
		return dao.Event{}, err
	}

	return event, nil
}

// GetAllEvent retrieves all events from the repository.
// It returns a slice of dao.Event and an error if the operation fails.
func (e EventServiceImpl) GetAllEvent() ([]dao.Event, error) {
	log.Info("Start to execute get all event")

	events, err := e.eventRepo.FindAllEvent()
	if err != nil {
		return nil, err
	}

	return events, nil
}

// GetEventById retrieves a event from the repository by their ID.
// It returns the dao.Event with the specified ID and an error if the operation fails.
func (e EventServiceImpl) GetEventById(eventId int) (dao.Event, error) {
	log.Info("Start to execute get event by id")

	event, err := e.eventRepo.FindEventById(eventId)
	if err != nil {
		return dao.Event{}, err
	}

	return event, nil
}

// UpdateEventById updates a event's details by their ID.
// Access is restricted to the resource owner.
// It modifies the event's name, description, location, event time if provided in the request.
// It returns the updated dao.Event and an error if the operation fails.
func (e EventServiceImpl) UpdateEventById(request dao.Event, eventId, userId int) (dao.Event, error) {
	log.Info("Start to execute update event by id")

	event, err := e.eventRepo.FindEventById(eventId)
	if err != nil {
		return dao.Event{}, err
	}

	if event.UserID != userId {
		log.Info("Access denied. Not a resource owner")
		return dao.Event{}, pkg.NewUnauthorizedError("Unauthorized", nil)
	}

	if request.Name != "" {
		event.Name = request.Name
	}
	if request.Description != "" {
		event.Description = request.Description
	}
	if request.Location != "" {
		event.Location = request.Location
	}
	if !request.EventTime.IsZero() {
		event.EventTime = request.EventTime
	}

	event, err = e.eventRepo.Save(&event)
	if err != nil {
		return dao.Event{}, err
	}

	return event, nil
}

// DeleteEventById removes a event from the repository by their ID.
// Access is restricted to the resource owner.
// It returns an error if the operation fails.
func (e EventServiceImpl) DeleteEventById(eventId, userId int) error {
	log.Info("Start to execute delete event by id")

	event, err := e.eventRepo.FindEventById(eventId)
	if err != nil {
		return err
	}

	if event.UserID != userId {
		log.Info("Access denied. Not a resource owner")
		return pkg.NewUnauthorizedError("Unauthorized", nil)
	}

	err = e.eventRepo.DeleteEventById(eventId)
	if err != nil {
		return err
	}

	return nil
}

func EventServiceInit(eventRepository repository.EventRepository) *EventServiceImpl {
	return &EventServiceImpl{
		eventRepo: eventRepository,
	}
}
