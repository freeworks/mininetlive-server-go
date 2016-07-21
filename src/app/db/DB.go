package db

import (
	. "app/admin"
	. "app/common"
	models "app/models"
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
	models.Dbmap = dbmap
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
	return dbmap
}
