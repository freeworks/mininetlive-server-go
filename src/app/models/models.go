package models

import (
	. "app/common"
	"database/sql/driver"
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/coopernurse/gorp"
)

var Dbmap *gorp.DbMap

type Resp struct {
	Ret  int64       `form:"ret" json:"ret"`
	Msg  string      `form:"msg" json:"msg"`
	Data interface{} `form:"data" json:"data"`
}

type User struct {
	Id           int       `form:"id" json:"-" ` //db:"id,primarykey, autoincrement"
	Uid          string    `form:"uid"  json:"uid" db:"uid"`
	EasemobUuid  string    `json:"easemobUuid" db:"easemob_uuid"`
	NickName     string    `form:"nickname" json:"nickname" binding:"required"  db:"nickname"`
	Avatar       string    `form:"avatar" json:"avatar"  db:"avatar"`
	Gender       int       `form:"gender" json:"gender" binding:"required" db:"gender"` //binding:"required"  TODO 0 default not bindle
	Balance      int       `form:"balance" json:"balance" db:"balance"`
	InviteCode   string    `form:"inviteCode" json:"inviteCode" db:"invite_code"`
	Qrcode       string    `form:"qrcode" json:"qrcode" db:"qrcode"`
	Phone        string    `form:"phone" json:"phone" db:"phone"`
	BeInvitedUid string    `json:"-" db:"be_invited_uid"`
	Updated      time.Time `json:"-" db:"update_time"`
	Created      time.Time `json:"-" db:"create_time"`
}

func (u User) Value() (driver.Value, error) {
	return u.Uid, nil
}

func (u *User) Scan(value interface{}) (err error) {
	switch src := value.(type) {
	case []byte:
		u.Uid = string(src)
		Dbmap.SelectOne(u, "SELECT * FROM t_user WHERE uid  = ?", u.Uid)
	case int:
		u.Uid = strconv.Itoa(src)
		Dbmap.SelectOne(u, "SELECT * FROM t_user WHERE uid  = ?", u.Uid)
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

type Record struct {
	Id      int      `db:"id" json:"-"`
	Aid     string   `db:"aid" json:"aid"`
	Uid     string   `db:"uid" json:"-"`
	Type    int      `db:"type" json:"-"` //0 预约，1，观看，2 支付，购买
	Created JsonTime `db:"create_time" json:"createTime"`
}

func (pl *Record) PreInsert(s gorp.SqlExecutor) error {
	pl.Created = JsonTime{time.Now(), true}
	return nil
}

type Activity struct {
	Id               int       `json:"id" db:"id"`
	Aid              string    `json:"aid" db:"aid"`
	Title            string    `form:"title" json:"title"  binding:"required" db:"title"`
	Date             JsonTime  `json:"date" db:"date"`
	ADate            int64     `form:"date" json:"-"  db:"-"` /*binding:"required"*/
	Desc             string    `form:"desc" json:"desc" binding:"required" db:"desc"`
	FrontCover       string    `form:"frontCover" json:"frontCover" binding:"required" db:"front_cover"`
	Price            int       `form:"price" json:"price"  db:"price"`
	Password         string    `form:"password" json:"-" db:"pwd"`
	StreamId         string    `json:"streamId" json:"streamId" db:"stream_id"`
	StreamType       int       `form:"streamType" json:"streamType" db:"stream_type"` //0 直播，1 视频
	LivePullPath     string    `json:"livePullPath" db:"live_pull_path"`
	VideoPath        string    `json:"videoPath" db:"video_path"`
	ActivityState    int       `json:"activityState" db:"activity_state"`                                      //0 未开播， 1 直播中 2 直播结束
	ActivityType     int       `form:"activityType" json:"activityType" binding:"required" db:"activity_type"` //0免费，1收费
	PlayCount        int       `json:"playCount" db:"play_count"`
	AppointmentCount int       `json:"appointmentCount" db:"appointment_count"`
	PayState         int       `json:"payState" db:"-"`
	AppointState     int       `json:"appoinState" db:"-"`
	GroupId          string    `json:"groupId" db:"group_id"`
	OnlineCount      int       `json:"onlineCount" db:"online_count"`
	IsRecommend      int       `json:"-" db:"is_recommend"`
	Updated          time.Time `json:"-" db:"update_time"`
	Created          JsonTime  `json:"createTime" db:"create_time"`
}

func (a *Activity) PreInsert(s gorp.SqlExecutor) error {
	a.Created = JsonTime{time.Now(), true}
	a.Updated = a.Created.Time
	return nil
}

func (a *Activity) PreUpdate(s gorp.SqlExecutor) error {
	a.Updated = time.Now()
	return nil
}

func (a Activity) String() string {
	return fmt.Sprintf("[%d, %s, %s]", a.Id, a.Aid, a.Title)
}

type QActivity struct {
	Activity
	Owner        User   `json:"owner" db:"uid"`
	LivePushPath string `json:"-" db:"live_push_path"`
}

type NActivity struct {
	Activity
	Uid          string `db:"uid"`
	IsRecord     bool   `from:"isRecord" json:"-" db:"-"`
	LivePushPath string `json:"livePushPath" db:"live_push_path"`
}

type Recomend struct {
	Id      string    `db:"id"`
	aid     string    `db:"aid"`
	Updated time.Time `db:"update_time"`
	Created time.Time `db:"create_time"`
}

type Order struct {
	Id       int      `json:"-" db:"id"`
	OrderNo  string   `json:"no" db:"no"`
	Amount   uint64   `json:"amount" db:"amount"`
	Channel  string   `json:"channel" db:"channel"`
	ClientIP string   `json:"ip" db:"client_ip"`
	Subject  string   `json:"subject" db:"subject"`
	Aid      string   `json:"aid" db:"aid"`
	PayType  int      `json:"type" db:"type"`
	Created  JsonTime `json:"createTime" db:"create_time"`
}
