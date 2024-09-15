package dao

type User struct {
	ID       int    `gorm:"column:id; primary_key; not null" json:"-"`
	Email    string `gorm:"column:email; type:varchar(255); not null; uniqueIndex" json:"email" validate:"email"`
	Password string `gorm:"column:password; not null" json:"password,omitempty" validate:"required"`
	RoleID   int    `gorm:"column:role_id; not null" json:"role_id"`
	Role     Role   `gorm:"foreignKey:RoleID; references:ID" json:"-"`
	BaseModel
}

type UserResponse struct {
	ID     int    `json:"id"`
	Email  string `json:"email"`
	RoleID int    `json:"role_id"`
}
