package repository

import (
	"errors"
	"event-booking-api/app/domain/dao"
	"event-booking-api/app/pkg"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type EventRepository interface {
	Save(request *dao.Event) (dao.Event, error)
	FindAllEvent() ([]dao.Event, error)
	FindEventById(id int) (dao.Event, error)
	DeleteEventById(id int) error
}

type EventRepositoryImpl struct {
	db *gorm.DB
}

// Save stores the event to the database.
// It returns the saved dao.Event and an error, if any.
func (e EventRepositoryImpl) Save(request *dao.Event) (dao.Event, error) {
	err := e.db.Save(request).Error
	if err != nil {
		log.Error("Error saving event: ", err)
		return dao.Event{}, err
	}

	return *request, nil
}

// FindAllEvent retrieves all events from the database.
// It returns a slice of dao.Event and an error, if any.
func (e EventRepositoryImpl) FindAllEvent() ([]dao.Event, error) {
	var events []dao.Event

	err := e.db.Select("id, name, description, location, event_time, user_id").
		Find(&events).Error
	if err != nil {
		log.Error("Error finding all events: ", err)
		return nil, err
	}

	return events, nil
}

// FindEventById retrieves a event by the given ID from the database.
// It returns the dao.Event and an error, if any.
func (e EventRepositoryImpl) FindEventById(id int) (dao.Event, error) {
	event := dao.Event{ID: id}

	err := e.db.First(&event).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Info("Error finding event by id: ", err)
			return dao.Event{}, pkg.NewNotFoundError("Event not found", err)
		}

		log.Error("Error finding event by id: ", err)
		return dao.Event{}, err
	}

	return event, nil
}

// DeleteEventById deletes the event by the given ID from the database.
// It returns an error if the deletion fails.
func (e EventRepositoryImpl) DeleteEventById(id int) error {
	err := e.db.Delete(&dao.Event{}, id).Error
	if err != nil {
		log.Error("Error deleting event: ", err)
		return err
	}

	return nil
}

func EventRepositoryInit(db *gorm.DB) *EventRepositoryImpl {
	if err := db.AutoMigrate(&dao.Event{}); err != nil {
		log.Fatal("Error AutoMigrating Event: ", err)
	}

	return &EventRepositoryImpl{
		db: db,
	}
}
