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
	Id         int       `form:"id" json:"-" ` //db:"id,primarykey, autoincrement"
	Uid        string    `form:"uid"  json:"uid" db:"uid"`
	NickName   string    `form:"nickname" json:"nickname" binding:"required"  db:"nickname"`
	Avatar     string    `form:"avatar" json:"nickname"  db:"avatar"`
	Gender     int       `form:"gender" json:"gender" db:"gender"` //binding:"required"  TODO 0 default not bindle
	Balance    int       `form:"balance" json:"balance" db:"balance"`
	InviteCode string    `form:"inviteCode" json:"inviteCode" db:"invite_code"`
	Qrcode     string    `form:"qrcode" json:"qrcode" db:"qrcode"`
	Phone      string    `form:"phone" json:"phone" db:"phone"`
	Updated    time.Time `json:"-" db:"update_time"`
	Created    time.Time `json:"-" db:"create_time"`
}

func (u User) String() string {
	return fmt.Sprintf("[%s,%s, %s, %d]", u.Id, u.Uid, u.NickName, u.Gender)
}

type OAuth struct {
	Id          int       `form:"id" json:"-"` //  `form:"id"  db:"id,primarykey, autoincrement"`
	Uid         string    `form:"uid" json:"uid" db:"uid"`
	Plat        string    `form:"plat" json:"plat" binding:"required" db:"plat"`
	OpenId      string    `form:"openid" json:"openid" binding:"required" db:"openid"`
	AccessToken string    `form:"access_token" json:"access_token" binding:"required" db:"access_token"`
	ExpiresIn   int       `form:"expires_in" json:"-" binding:"required" db:"-"` //- 忽略的意思
	Expires     time.Time `db:"expires" json:"-" `
}

type LocalAuth struct {
	Id       int       `form:"id" json:"-"` //  `form:"id"  db:"id,primarykey, autoincrement"`
	Uid      string    `form:"uid" json:"uid" db:"uid"`
	Phone    string    `form:"phone" json:"phone" binding:"required"  db:"phone"`
	Password string    `form:"password" json:"password" binding:"required" db:"password"`
	Token    string    `db:"token" json:"token"`
	Expires  time.Time `db:"expires" json:"expires"`
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
	Id      int       `db:"id" json:"-"`
	Aid     string    `db:"aid" json:"aid"`
	Uid     string    `db:"uid" json:"uid"`
	State   int       `db:"state" json:"state"` //0 未开始，1 活动过期，3,取消
	Created time.Time `db:"create_time" json:"createTime"`
}

type PayRecord struct {
	Id      int       `db:"id" json:"-"`
	Aid     string    `db:"aid" json:"aid"`
	Uid     int       `db:"uid"  json:"uid"`
	Amount  int       `db:"amount" json:"amount"`
	Type    int       `db:"type" json:"type"` //0 支付观看，1奖赏
	Created time.Time `db:"create_time" json:"createTime"`
}

type PlayRecord struct {
	Id      int       `db:"id" json:"-"`
	Aid     string    `db:"aid" json:"aid"`
	Uid     string    `db:"uid" json:"uid"`
	Type    int       `db:"type" json:"type"` //0 直播，1点播
	Created time.Time `db:"create_time" json:"create_time"`
}

type Activity struct {
	Id               int       `form:"id" json:"-" db:"id"`
	Aid              string    `form:"aid" json:"aid" db:"aid"`
	Title            string    `form:"title" json:"title"  binding:"required" db:"title"`
	Date             time.Time `db:"date" json:"date"`
	ADate            int64     `form:"date" json:"-"  db:"-"` /*binding:"required"*/
	Desc             string    `form:"desc" json:"desc" binding:"required" db:"desc"`
	FontCover        string    `form:"fontCover" json:"fontCover" binding:"required" db:"front_cover"`
	Type             int       `form:"type" json:"type" binding:"required" db:"type"` //0直播，1点播
	Price            int       `form:"price" json:"price"  db:"price"`
	Password         string    `form:"password" json:"-" db:"pwd"`
	Uid              string    `form:"uid" json:"uid" db:"uid"`
	VideoId          string    `form:"videoId" json:"videoId" db:"video_id"`
	VideoType        int       `form:"videoType" json:"videoType" binding:"required" db:"video_type"` //0 免费， 1收费
	VideoPullPath    string    `form:"videoPullPath" json:"videoPullPath" db:"video_pull_path"`
	VideoPushPath    string    `form:"videoPushPath" json:"videoPushPath" db:"video_push_path"`
	VideoStorePath   string    `form:"videoPushPath" json:"-" db:"video_store_path"`
	State            int       `json:"state" db:"state"` //0.未开播，1.正在直播，2.可点播，3.已下线
	PlayCount        int       `json:"playCount" db:"play_count"`
	AppointmentCount int       `json:"appointmentCount" db:"appointment_count"`
	Updated          time.Time `json:"-" db:"update_time"`
	Created          time.Time `json:"-" db:"create_time"`
}

func (a Activity) String() string {
	return fmt.Sprintf("[%s, %s, %s]", a.Id, a.Title, a.FontCover)
}
