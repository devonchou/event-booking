package repository

import (
	"errors"
	"event-booking-api/app/domain/dao"
	"event-booking-api/app/pkg"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type RegisterRepository interface {
	Save(request *dao.Register) error
	Delete(eventId, userId int) error
	FindAttendeesEmailById(eventId int) ([]string, error)
}

type RegisterRepositoryImpl struct {
	db *gorm.DB
}

// Save stores a new user registration for an event to the database.
// It returns an error, if any.
func (r RegisterRepositoryImpl) Save(request *dao.Register) error {
	err := r.db.Save(request).Error
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			log.Info("Error saving register: ", err)
			return pkg.NewConflictError("register record already exist", err)
		}

		log.Error("Error saving register: ", err)
		return err
	}

	return nil
}

// Delete deletes the register entry by the given event and user ID from the database.
// It returns an error if the deletion fails.
func (r RegisterRepositoryImpl) Delete(eventId, userId int) error {
	err := r.db.Where("event_id = ? AND user_id = ?", eventId, userId).Delete(&dao.Register{}).Error
	if err != nil {
		log.Error("Error deleting register entry by event id and user id: ", err)
		return err
	}

	return nil
}

// FindAttendeesEmailByEventID retrieves the email addresses of all attendees for a given event ID.
// It returns a slice of emails and an error, if any.
func (r RegisterRepositoryImpl) FindAttendeesEmailById(eventId int) ([]string, error) {
	var emails []string

	err := r.db.Model(&dao.Register{}).
		Joins("JOIN users ON registers.user_id = users.id").
		Where("registers.event_id = ?", eventId).
		Pluck("users.email", &emails).Error

	if err != nil {
		log.Error("Error finding attendees email by event id: ", err)
		return nil, err
	}

	return emails, nil
}

func RegisterRepositoryInit(db *gorm.DB) *RegisterRepositoryImpl {
	if err := db.AutoMigrate(&dao.Register{}); err != nil {
		log.Fatal("Error AutoMigrating Register: ", err)
	}

	return &RegisterRepositoryImpl{
		db: db,
	}
}
