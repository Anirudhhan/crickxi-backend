package models

type RegisterUser struct {
	Name     string `db:"name" json:"name" binding:"required,min=2"`
	Phone    string `db:"phone" json:"phone" binding:"required"`
	Password string `db:"password" json:"password" binding:"required,min=8"`
}
