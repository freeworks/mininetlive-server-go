package models

import (
	"fmt"
	"time"
	"reflect"
	"strconv"
	"database/sql/driver"
	"github.com/coopernurse/gorp"
)

var Dbmap *gorp.DbMap

type Resp struct {
	Ret  int64       `form:"ret" json:"ret"`
	Msg  string      `form:"msg" json:"msg"`
	Data interface{} `form:"data" json:"data"`
}

type User struct {
	Id          int       `form:"id" json:"-" ` //db:"id,primarykey, autoincrement"
	Uid         string    `form:"uid"  json:"uid" db:"uid"`
	EasemobUuid string    `json:"easemobUuid" db:"easemob_uuid"`
	NickName    string    `form:"nickname" json:"nickname" binding:"required"  db:"nickname"`
	Avatar      string    `form:"avatar" json:"avatar"  db:"avatar"`
	Gender      int       `form:"gender" json:"gender" binding:"required" db:"gender"` //binding:"required"  TODO 0 default not bindle
	Balance     int       `form:"balance" json:"balance" db:"balance"`
	InviteCode  string    `form:"inviteCode" json:"inviteCode" db:"invite_code"`
	Qrcode      string    `form:"qrcode" json:"qrcode" db:"qrcode"`
	Phone       string    `form:"phone" json:"phone" db:"phone"`
	BeInvitedUid string  `json:"-" db:"be_invited_uid"`
	Updated     time.Time `json:"-" db:"update_time"`
	Created     time.Time `json:"-" db:"create_time"`
}

func (u User) Value() (driver.Value, error) {
	return u.Uid, nil
}

func (u *User) Scan(value interface{}) (err error) {
	switch src := value.(type) {
	case []byte:
		u.Uid = string(src)
		Dbmap.SelectOne(u,"SELECT * FROM t_user WHERE uid  = ?",u.Uid)
	case int:
		u.Uid = strconv.Itoa(src)
		Dbmap.SelectOne(u,"SELECT * FROM t_user WHERE uid  = ?",u.Uid)
	default:
		typ := reflect.TypeOf(value)
		return fmt.Errorf("Expected person value to be convertible to int64, got %v (type %s)", value, typ)
	}
	return
}

func (u *User) PreInsert(s gorp.SqlExecutor) error {
    u.Created = time.Now()
    u.Updated = u.Created
    return nil
}

func (u *User) PreUpdate(s gorp.SqlExecutor) error {
    u.Updated = time.Now()
    return nil
}

func (u *User) String() string {
	return fmt.Sprintf("[%d,%s, %s, %d]", u.Id, u.Uid, u.NickName, u.Gender)
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

func (a *AppointmentRecord) PreInsert(s gorp.SqlExecutor) error {
    a.Created = time.Now()
    return nil
}

type PayRecord struct {
	Id      int       `db:"id" json:"-"`
	Aid     string    `db:"aid" json:"aid"`
	Uid     int       `db:"uid"  json:"uid"`
	Amount  int       `db:"amount" json:"amount"`
	Type    int       `db:"type" json:"type"` //0 支付观看，1奖赏
	Created time.Time `db:"create_time" json:"createTime"`
}

func (p *PayRecord) PreInsert(s gorp.SqlExecutor) error {
    p.Created = time.Now()
    return nil
}

type PlayRecord struct {
	Id      int       `db:"id" json:"-"`
	Aid     string    `db:"aid" json:"aid"`
	Uid     string    `db:"uid" json:"uid"`
	Type    int       `db:"type" json:"type"` //0 直播，1点播
	Created time.Time `db:"create_time" json:"create_time"`
}

func (pl *PlayRecord) PreInsert(s gorp.SqlExecutor) error {
    pl.Created = time.Now()
    return nil
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
	Owner  			 User 	   `db:"uid",json:"owner"`
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



func (a *Activity) PreInsert(s gorp.SqlExecutor) error {
    a.Created = time.Now()
    a.Updated = a.Created
    return nil
}

func (a *Activity) PreUpdate(s gorp.SqlExecutor) error {
    a.Updated = time.Now()
    return nil
}

func (a Activity) String() string {
	return fmt.Sprintf("[%s, %s, %s]", a.Id, a.Title, a.FontCover)
}
