package db

import (
	"database/sql"
	"github.com/coopernurse/gorp"
	."app/models"
	."app/common"
	)

func InitDb() *gorp.DbMap {
	db, err := sql.Open("mysql", "root:weiwanglive@tcp(106.75.19.205:3306)/minnetlive?parseTime=true")
	CheckErr(err, "sql.Open failed")
	// construct a gorp DbMap
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"}}

	dbmap.AddTableWithName(User{}, "t_user").SetKeys(true, "Id")
	dbmap.AddTableWithName(OAuth{}, "t_oauth").SetKeys(true, "Id")
	dbmap.AddTableWithName(LocalAuth{}, "t_local_auth").SetKeys(true, "Id")
	dbmap.AddTableWithName(Activity{}, "t_activity").SetKeys(true, "Id")

	//TODO cdreate table
	//	// add a table, setting the table name to 'posts' and
	//	// specifying that the Id property is an auto incrementing PK
	//	dbmap.AddTableWithName(Post{}, "posts").SetKeys(true, "Id")
	//	// create the table. in a production system you'd generally
	//	// use a migration tool, or create the tables via scripts
	//	err = dbmap.CreateTablesIfNotExists()
	//	checkErr(err, "Create tables failed")

	return dbmap
}