package dao

import (
	"time"

	"gorm.io/gorm"
)

type BaseModel struct {
	CreatedAt time.Time      `gorm:"->:false;<-:create; column:created_at" json:"-"`
	UpdatedAt time.Time      `gorm:"->:false;<-:update; column:updated_at" json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"->:false; column:deleted_at" json:"-"`
}
