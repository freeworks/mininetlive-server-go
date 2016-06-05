package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	//	"strconv"
	"time"

	"github.com/go-martini/martini"
	_ "github.com/go-sql-driver/mysql"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
	//	"gopkg.in/gorp.v1"
	"github.com/coopernurse/gorp"
)

type Resp struct {
	Ret  string      `form:"ret"`
	Msg  string      `form:"msg"`
	Data interface{} `form:"data"`
}

type User struct {
	ID     int    `form:"id"  db:"id,primarykey, autoincrement"`
	Name   string `form:"name" binding:"required"  db:"name"`
	Avatar string `form:"avatar" db:"avatar"`
	Gender int    `form:"gender" db:"gender"`
	Note   string `form:"note" db:"note" `
}

type OAuth struct {
	ID          int       `form:"id"  db:"id,primarykey, autoincrement"`
	UserId      int       `form:"userId"  db:"iuser_id"`
	Plat        string    `form:"plat" binding:"required" db:"plat"`
	OpenId      string    `form:"openid" binding:"required" db:"openid"`
	AccessToken string    `form:"access_token" binding:"required" db:"access_token"`
	ExpiresIn   int       `form:"expires_in" binding:"required" db:"-"` //- 忽略的意思
	Expires     time.Time `db:"expires"`
}

type RegisterModel struct {
	User  User
	OAuth OAuth
}

type Activity struct {
	ID    int
	Title string
}

func GetActivityList(r render.Render, dbmap *gorp.DbMap) {
	var activities []Activity
	_, err := dbmap.Select(&activities, "select * from t_activity")
	checkErr(err, "Select failed")
	newmap := map[string]interface{}{"metatitle": "this is my custom title", "activities": activities}
	r.HTML(200, "posts", newmap)
}

func GetActivity(args martini.Params, r render.Render, dbmap *gorp.DbMap) {
	var activities []Activity
	_, err := dbmap.Select(&activities, "select * from t_activity")
	checkErr(err, "Select failed")
	newmap := map[string]interface{}{"metatitle": "this is my custom title", "activities": activities}
	r.HTML(200, "posts", newmap)

	var activity Activity
	err = dbmap.SelectOne(&activity, "select * from t_activty where id =?", args["id"])
	//simple error check
	if err != nil {
		newmap := map[string]interface{}{"metatitle": "404 Error", "message": "This is not found"}
		r.HTML(404, "error", newmap)
	} else {
		newmap := map[string]interface{}{"metatitle": activity.Title + " more custom", "activity": activity}
		r.HTML(200, "activity", newmap)
	}
}

func NewActivity(activity Activity, r render.Render, dbmap *gorp.DbMap) {
	log.Println(activity)
	err := dbmap.Insert(&activity)
	checkErr(err, "Insert failed")
	newmap := map[string]interface{}{"metatitle": "created activity", "activity": activity}
	r.HTML(200, "post", newmap)
}

//登陆
func Login(req *http.Request, r render.Render, db *sql.DB, dbmap *gorp.DbMap) {
	plat := req.FormValue("plat")
	if plat == "SinaWeibo" || plat == "QQ" {
		openId, accessToken, expiresIn := req.FormValue("openid"), req.FormValue("access_token"), req.FormValue("expires_in")
		log.Println("login : openId->" + openId + "accessToken :" + accessToken + ",expiresIn:" + expiresIn)
		row := db.QueryRow("SELECT access_token FROM t_oauth WHERE plat = ? AND openid = ?", plat, openId)
		var token string
		err := row.Scan(&token)
		if err != nil {
			r.JSON(200, "unregister")
		} else {
			r.JSON(200, "ok")
		}
	} else {
		username, password := req.FormValue("username"), req.FormValue("password")
		row := db.QueryRow("SELECT id FROM t_local_oauth WHERE username = ? AND password = ?", username, password)
		var id int
		err := row.Scan(&id)
		fmt.Println(id)
		panic(err.Error()) // TODO
		if row == nil {
			r.JSON(200, "ok")
		} else {
			r.JSON(200, "user not found")
		}
	}
}

func Register(register RegisterModel, r render.Render, dbmap *gorp.DbMap) {
	var oauth OAuth
	err := dbmap.SelectOne(&oauth, "select * from t_oauth where openid=?", register.OAuth.OpenId)
	if err != nil {
		log.Fatal(err)
	}
	if &oauth == nil {
		trans, err := dbmap.Begin()
		if err != nil {
			log.Fatal(err)
		}
		trans.Insert(register.User)
		//		expiresInInt, err := strconv.Atoi(oauth.ExpiresIn)
		//		oauth.ExpiresIn = time.Now().Add(time.Second * time.Duration(expiresInInt))
		oauth.Expires = time.Now().Add(time.Second * time.Duration(oauth.ExpiresIn))
		oauth.UserId = register.User.ID
		trans.Insert(oauth)
		// if the commit is successful, a nil error is returned
		err = trans.Commit()
		if err == nil {
			r.JSON(200, "registered success ")
		}
	} else {
		r.JSON(200, "registered ")
	}

}

////注册
//func Register(req *http.Request, oauth Oauth, user User, r render.Render) {
//	plat := req.FormValue("plat")
//	if plat == "SinaWeibo" || plat == "QQ" {
//		openId, accessToken, expiresIn := req.FormValue("openid"), req.FormValue("access_token"), req.FormValue("expires_in")
//		username, gender, avatar := req.FormValue("name"), req.FormValue("gender"), req.FormValue("avatar")
//		row := db.QueryRow("SELECT access_token FROM t_oauth WHERE plat = ? AND openid = ?", plat, openId)
//		var token string
//		row.Scan(&token)
//		if token == "" {
//			trans, err := dbmap.Begin()
//			if err != nil {
//				return err
//			}
//			trans.Insert(user)
//			oauth.userId = user.Id
//			trans.Insert(oauth)
//			// if the commit is successful, a nil error is returned
//			return trans.Commit()

//			tx, err := db.Begin()
//			if err != nil {
//				log.Fatal(err)
//			}
//			db.Create(&user)
//			stmt2, err2 := db.Prepare("INSERT INTO t_oauth(user_id,plat,openid,access_token,expires) VALUES (?,?,?,?,?)")
//			defer stmt2.Close()
//			if err2 != nil {
//				log.Fatal(err2)
//			}

//			defer tx.Rollback()
//			res, err := stmt.Exec(username, gender, avatar)
//			id, err := res.LastInsertId()
//			expiresInInt, err := strconv.Atoi(expiresIn)
//			_, err = stmt2.Exec(id, plat, openId, accessToken, time.Now().Add(time.Second*time.Duration(expiresInInt)))
//			tx.Commit()
//			if err != nil || err2 != nil {
//				log.Println(err)
//				log.Println(err2)
//			} else {
//				r.JSON(200, "register OK")
//			}
//		} else {
//			r.JSON(200, "fail registered ")
//		}
//	} else {
//		username, password := req.FormValue("username"), req.FormValue("password")
//		row := db.QueryRow("SELECT id FROM t_local_oauth WHERE username = ? AND password = ?", username, password)
//		var id int
//		err := row.Scan(&id)
//		fmt.Println(id)
//		panic(err.Error()) // TODO
//		if row == nil {
//			r.JSON(200, "ok")
//		} else {
//			r.JSON(200, "user not found")
//		}
//	}
//}

func UpdateUser(user User, dbmap *gorp.DbMap) {
	dbmap.Update(&user)
}

func DeleteUser(args martini.Params, dbmap *gorp.DbMap) {
	_, err := dbmap.Exec("delete from t_user where id=?", args["id"])
	checkErr(err, "Exec failed")
}

func GetUser(args martini.Params, r render.Render, dbmap *gorp.DbMap) {
	var user User
	err := dbmap.SelectOne(&user, "select * from t_user where id=?", args["id"])
	//simple error check
	if err != nil {
		r.JSON(200, "not found")
	} else {
		newmap := map[string]interface{}{"metatitle": user.Name + " more custom", "user": user}
		r.JSON(200, newmap)
	}
}

func main() {
	//setup db
	dbmap := initDb()
	defer dbmap.Db.Close()

	m := martini.Classic()

	m.Map(dbmap)

	m.NotFound(func() {
		// 处理 404
	})
	m.Use(func(res http.ResponseWriter, req *http.Request) {
		//		if req.Header.Get("X-API-KEY") != "secret123" {
		//			res.WriteHeader(http.StatusUnauthorized)
		//		}
	})
	m.Use(render.Renderer())

	m.Group("/activity", func(r martini.Router) {
		r.Get("/", GetActivityList)
		r.Get("/:id", GetActivity)
		r.Post("/", NewActivity)
		//		r.Put("/update/:id", UpdateActivity)
		//		r.Delete("/delete/:id", DeleteActivity)
	})

	m.Group("/user", func(r martini.Router) {
		r.Post("/login", Login)
		r.Post("/register", binding.Bind(RegisterModel{}), Register)
		r.Get("/:id", GetUser)
		r.Put("/:id", binding.Bind(User{}), UpdateUser)
		r.Delete("/:id", DeleteUser)
	})

	m.Run()
}

func initDb() *gorp.DbMap {
	db, err := sql.Open("mysql", "root:weiwanglive@tcp(106.75.19.205:3306)/minnetlive")
	checkErr(err, "sql.Open failed")
	// construct a gorp DbMap
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"}}

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

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}
