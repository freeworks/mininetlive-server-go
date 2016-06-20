package admin

import (
	. "app/models"
	"time"
)

// AdminModel can be any struct that represents a user in my system
type AdminModel struct {
	Id            int64     `form:"id" db:"id"`
	Uuid          string    ` db:"uuid"`
	Username      string    `form:"name" db:"username"`
	Password      string    `form:"password" db:"password"`
	Avatar        string    `form:"avatar"  db:"avatar"`
	Updated       time.Time `db:"update_time"`
	Created       time.Time `db:"create_time"`
	authenticated bool      `form:"-" db:"-"`
}

type ThiredUserModel struct {
	User
	Plat string
}

type PhoneUserModel struct {
	User
	Phone string
}
