package models

import (
		"time"
		"fmt"
		)

type Resp struct {
	Ret  int64       `form:"ret" json:"ret"`
	Msg  string      `form:"msg" json:"msg"`
	Data interface{} `form:"data" json:"data"`
}

type User struct {
	Id     int    `form:"id" ` //db:"id,primarykey, autoincrement"
	Name   string `form:"name" binding:"required"  db:"name"`
	Avatar string `form:"avatar"  binding:"required"  db:"avatar"`
	Gender int    `form:"gender" binding:"required"  db:"gender"`
}

func (u User) String() string {
    return fmt.Sprintf("[%s, %s, %d]", u.Id, u.Name, u.Gender)
}

type OAuth struct {
	Id          int       `form:"id"` //  `form:"id"  db:"id,primarykey, autoincrement"`
	UserId      int       `form:"userId"  db:"user_id"`
	Plat        string    `form:"plat" binding:"required" db:"plat"`
	OpenId      string    `form:"openid" binding:"required" db:"openid"`
	AccessToken string    `form:"access_token" binding:"required" db:"access_token"`
	ExpiresIn   int       `form:"expires_in" binding:"required" db:"-"` //- 忽略的意思
	Expires     time.Time `db:"expires"`
}

type LocalAuth struct {
	Id       int       `form:"id"` //  `form:"id"  db:"id,primarykey, autoincrement"`
	UserId   int       `form:"userId"  db:"user_id"`
	Phone    string    `form:"phone" binding:"required"  db:"phone"`
	Password string    `form:"password" binding:"required" db:"password"`
	Token    string    `db:"token"`
	Expires  time.Time `db:"expires"`
}

type OAuthUser struct {
	User  User
	OAuth OAuth
}

type LocalAuthUser struct {
	User      User
	LocalAuth LocalAuth
}

type Activity struct {
	Id             int       `form:"id"  db:"id"`
	Title          string    `form:"title"  binding:"required" db:"title"`
	Date           time.Time `db:"date"`
	ADate          int64     `form:"date" binding:"required" db:"-"`
	Desc           string    `form:"desc" binding:"required" db:"desc"`
	FontCover      string    `form:"fontCover" binding:"required" db:"front_cover"`
	Type           int       `form:"type" binding:"required" db:"type"`
	Price          int       `form:"price"  db:"price"`
	Password       string    `form:"password"  db:"pwd"`
	BelongUserId   int       `form:"belongUserId"  db:"belong_user_id"`
	VideoId        string    `form:"videoId"  db:"video_id"`
	VideoType      int       `form:"videoType" binding:"required" db:"video_type"`
	VideoPullPath  string    `form:"videoPullPath"  db:"video_pull_path"`
	VideoPushPath  string    `form:"videoPushPath"  db:"video_push_path"`
	VideoStorePath string    `form:"videoPushPath"  db:"video_store_path"`
}
