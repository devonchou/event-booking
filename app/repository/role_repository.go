package repository

import (
	"event-booking-api/app/domain/dao"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type RoleRepository interface {
	FindAllRole()
}

type RoleRepositoryImpl struct {
	db *gorm.DB
}

func (r RoleRepositoryImpl) FindAllRole() {
	panic("implement me")
}

func RoleRepositoryInit(db *gorm.DB) *RoleRepositoryImpl {
	if err := db.AutoMigrate(&dao.Role{}); err != nil {
		log.Fatal("Error AutoMigrating Role: ", err)
	}

	return &RoleRepositoryImpl{
		db: db,
	}
}
