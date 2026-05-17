package models

type RegisterUser struct {
	Name     string `db:"name" json:"name" binding:"required,min=2"`
	Phone    string `db:"phone_no" json:"phone" binding:"required"`
	Password string `db:"password" json:"password" binding:"required,min=8"`
}

type LoginUserDetails struct {
	UserID       string `db:"user_id"`
	PlayerID     string `db:"player_id"`
	HashPassword string `db:"password"`
}

type SessionUserDetails struct {
	UserID   string `db:"user_id"`
	PlayerID string `db:"player_id"`
}

type LoginUser struct {
	Phone    string `db:"phone_no" json:"phone" binding:"required"`
	Password string `db:"password" json:"password" binding:"required"`
}
