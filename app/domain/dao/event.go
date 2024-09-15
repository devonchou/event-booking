package dao

import "time"

type Event struct {
	ID          int       `gorm:"column:id; primary_key; not null" json:"-"`
	Name        string    `gorm:"column:name; not null" json:"name" validate:"required"`
	Description string    `gorm:"column:description; not null" json:"description" validate:"required"`
	Location    string    `gorm:"column:location; not null" json:"location" validate:"required"`
	EventTime   time.Time `gorm:"column:event_time; not null" json:"event_time" validate:"required"`
	UserID      int       `gorm:"column:user_id; not null" json:"-"`
	User        User      `gorm:"foreignKey:UserID; references:ID" json:"-"`
	BaseModel
}

type EventResponse struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Location    string    `json:"location"`
	EventTime   time.Time `json:"event_time"`
	UserID      int       `json:"user_id"`
}
