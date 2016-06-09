package db

import (
	. "app/admin"
	. "app/common"
	. "app/models"
	"database/sql"
	"os"

	"github.com/coopernurse/gorp"
)

var dbmap *gorp.DbMap

func InitDb() *gorp.DbMap {
	_, err := os.Open("martini-sessionauth.bin")
	if err == nil {
		os.Remove("martini-sessionauth.bin")
	}

	db, err := sql.Open("mysql", "root:weiwanglive@tcp(106.75.19.205:3306)/minnetlive?parseTime=true")
	CheckErr(err, "sql.Open failed")
	// construct a gorp DbMap
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"}}

	dbmap.AddTableWithName(User{}, "t_user").SetKeys(true, "Id")
	dbmap.AddTableWithName(OAuth{}, "t_oauth").SetKeys(true, "Id")
	dbmap.AddTableWithName(LocalAuth{}, "t_local_auth").SetKeys(true, "Id")
	dbmap.AddTableWithName(Activity{}, "t_activity").SetKeys(true, "Id")
	dbmap.AddTableWithName(PlayRecord{}, "t_play_record").SetKeys(true, "Id")
	dbmap.AddTableWithName(PayRecord{}, "t_pay_record").SetKeys(true, "Id")
	dbmap.AddTableWithName(AppointmentRecord{}, "t_appointment_record").SetKeys(true, "Id")
	dbmap.AddTableWithName(AdminModel{}, "t_admin").SetKeys(true, "Id")

	return dbmap
}
