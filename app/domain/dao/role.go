package dao

type Role struct {
	ID   int    `gorm:"column:id; primary_key; not null" json:"id"`
	Role string `gorm:"column:role; not null" json:"role"`
	BaseModel
}
