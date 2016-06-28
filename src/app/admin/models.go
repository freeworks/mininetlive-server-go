package admin

import (
	"time"
)

type AdminModel struct {
	Id            int64     `form:"id" db:"id"`
	Uid          string    ` db:"uid"`
	Phone      string    	`form:"phone" db:"phone"`
	NickName      string    `form:"nickName" db:"nickname"`
	Password      string    `form:"password" db:"password"`
	Avatar        string    `form:"avatar"  db:"avatar"`
	EasemobUUID    string    `form:"-"  db:"easemob_uuid"`
	Updated       time.Time `db:"update_time"`
	Created       time.Time `db:"create_time"`
	Authenticated bool      `form:"-" db:"-"`
}
