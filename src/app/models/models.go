package models

import (
	"fmt"
	"time"
)

type Resp struct {
	Ret  int64       `form:"ret" json:"ret"`
	Msg  string      `form:"msg" json:"msg"`
	Data interface{} `form:"data" json:"data"`
}

type User struct {
	Id      int       `form:"id" ` //db:"id,primarykey, autoincrement"
	Name    string    `form:"name" binding:"required"  db:"name"`
	Avatar  string    `form:"avatar"  binding:"required"  db:"avatar"`
	Gender  int       `form:"gender" binding:"required"  db:"gender"`
	Updated time.Time `db:"update_time"`
	Created time.Time `db:"create_time"`
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

//TODO 预约活动是否已经过期
type AppointmentRecord struct {
	Id         int       `db:"id"`
	ActivityId int       `db:"activity_id"`
	UserId     int       `db:"user_id"`
	State      int       `db:"state"` //0 未开始，1 活动过期，3,取消
	Created    time.Time `db:"create_time"`
}

type PayRecord struct {
	Id             int       `db:"id"`
	ActivityId     int       `db:"activity_id"`
	ActivityUserId int       `db:"activity_user_id"`
	UserId         int       `db:"pay_user_id"`
	Amount         int       `db:"amount"`
	Type           int       `db:"type"` //0 支付观看，1奖赏
	Created        time.Time `db:"create_time"`
}

type PlayRecord struct {
	Id         int       `db:"id"`
	ActivityId int       `db:"activity_id"`
	UserId     int       `db:"play_user_id"`
	Type       int       `db:"type"` //0 直播，1点播
	Created    time.Time `db:"create_time"`
}

type Activity struct {
	Id               int       `form:"id"  db:"id"`
	Title            string    `form:"title"  binding:"required" db:"title"`
	Date             time.Time `db:"date"`
	ADate            int64     `form:"date"  db:"-"` /*binding:"required"*/
	Desc             string    `form:"desc" binding:"required" db:"desc"`
	FontCover        string    `form:"fontCover"  db:"front_cover"`       /*binding:"required"*/
	Type             int       `form:"type" binding:"required" db:"type"` //0直播，1点播
	Price            int       `form:"price"  db:"price"`
	Password         string    `form:"password"  db:"pwd"`
	BelongUserId     int       `form:"belongUserId"  db:"belong_user_id"`
	VideoId          string    `form:"videoId"  db:"video_id"`
	VideoType        int       `form:"videoType" binding:"required" db:"video_type"`
	VideoPullPath    string    `form:"videoPullPath"  db:"video_pull_path"`
	VideoPushPath    string    `form:"videoPushPath"  db:"video_push_path"`
	VideoStorePath   string    `form:"videoPushPath"  db:"video_store_path"`
	State            int       `db:"state"` //0.未开播，1.正在直播，2.可点播，3.已下线
	PlayCount        int       `db:"play_count"`
	AppointmentCount int       `db:"appointment_count"`
	Updated          time.Time `db:"update_time"`
	Created          time.Time `db:"create_time"`
}
