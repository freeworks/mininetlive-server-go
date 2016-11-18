package db

import (
	. "app/admin"
	. "app/common"
	//	config "app/config"
	models "app/models"
	wxpub "app/wxpub"
	"database/sql"
	"os"

	"github.com/coopernurse/gorp"
)

func InitDb() *gorp.DbMap {
	_, err := os.Open("martini-sessionauth.bin")
	if err == nil {
		os.Remove("martini-sessionauth.bin")
	}

	db, err := sql.Open("mysql", "root:weiwanglive@tcp(106.75.19.205:3306)/minnetlive?parseTime=true")
	CheckErr(err, "sql.Open failed")
	// construct a gorp DbMap
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"}}
	models.SetDbmap(dbmap)
	dbmap.AddTableWithName(models.User{}, "t_user").SetKeys(true, "Id")
	dbmap.AddTableWithName(models.OAuth{}, "t_oauth").SetKeys(true, "Id")
	dbmap.AddTableWithName(models.LocalAuth{}, "t_local_auth").SetKeys(true, "Id")
	dbmap.AddTableWithName(models.Activity{}, "t_activity").SetKeys(true, "Id")
	dbmap.AddTableWithName(models.Record{}, "t_record").SetKeys(true, "Id")
	dbmap.AddTableWithName(models.Order{}, "t_order").SetKeys(true, "Id")
	dbmap.AddTableWithName(models.Recomend{}, "t_recomend").SetKeys(true, "Id")
	dbmap.AddTableWithName(AdminModel{}, "t_admin").SetKeys(true, "Id")
	dbmap.AddTableWithName(models.NActivity{}, "t_activity").SetKeys(true, "Id")
	dbmap.AddTableWithName(models.QActivity{}, "t_activity").SetKeys(true, "Id")
	dbmap.AddTableWithName(models.InviteRelationship{}, "t_invite_relation").SetKeys(true, "Id")
	dbmap.AddTableWithName(wxpub.WXPub{}, "t_wxpub").SetKeys(true, "Id")
	dbmap.AddTableWithName(models.DividendRecord{}, "t_dividend_record").SetKeys(true, "Id")
	//	Uid         string    `form:"uid"  json:"uid" db:"uid"`
	//	NickName    string    `form:"nickname" json:"nickname" binding:"required"  db:"nickname"`
	//	Avatar      string    `form:"avatar" json:"avatar"  db:"avatar"`
	//	Gender      int       `form:"gender" json:"gender" db:"gender"` //binding:"required"  TODO 0 default not bindle
	//	Balance     int       `form:"balance" json:"balance" db:"balance"`
	//	InviteCode  string    `form:"-" json:"inviteCode" db:"invite_code"`
	//	Qrcode      string    `form:"qrcode" json:"qrcode" db:"qrcode"`
	//	Phone       string    `form:"phone" json:"phone" db:"phone"`
	//	user := user{
	//		Uid:      UID(),
	//		NickName:  config.NickName,
	//		Avatar:   config.Avatar,
	//		Gender:  1
	//		Phone:  config.Phone,
	//	}

	//	LocalAuth

	//	trans, err := dbmap.Begin()
	//	authUser.User.Uid = UID()
	//	authUser.User.InviteCode = GeneraVCode6()
	//	authUser.User.Qrcode = "http://h.hiphotos.baidu.com/image/pic/item/3bf33a87e950352a5936aa0a5543fbf2b2118b59.jpg"
	//	err = trans.Insert(&authUser.User)
	//	CheckErr(err, "Register insert user failed")
	//	authUser.LocalAuth.Token = Token()
	//	authUser.LocalAuth.Uid = authUser.User.Uid
	//	authUser.LocalAuth.Expires = time.Now().Add(time.Hour * 24 * 30)
	//	err = trans.Insert(&authUser.LocalAuth)
	//	CheckErr(err, "Register insert LocalAuth failed")
	//	err = trans.Commit()
	return dbmap
}
