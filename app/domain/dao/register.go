package dao

type Register struct {
	ID      int   `gorm:"column:id; primary_key; not null" json:"id"`
	EventID int   `gorm:"column:event_id; not null; uniqueIndex:idx_event_user" json:"event_id"`
	Event   Event `gorm:"foreignKey:EventID;references:ID" json:"-"`
	UserID  int   `gorm:"column:user_id; not null; uniqueIndex:idx_event_user" json:"user_id"`
	User    User  `gorm:"foreignKey:UserID;references:ID" json:"-"`
	BaseModel
}
