package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-martini/martini"
	_ "github.com/go-sql-driver/mysql"
	//	"github.com/martini-contrib/binding"
	"github.com/jinzhu/gorm"
	"github.com/martini-contrib/render"
)

type Resp struct {
	Ret  string      `form:"ret"`
	Msg  string      `form:"msg"`
	Data interface{} `form:"data"`
}

type User struct {
	gorm.Model
	Name   string `form:"name" binding:"required"`
	Avatar string `form:"avatar"`
	Gender bool   `form:"gender"`
	Note   string `form:"note" `
}

type Activity struct {
	gorm.Model
	ID    int
	Title string
}

func GetActivity(params martini.Params) {
	//	id := params["name"]
}

func NewActivity() {

}

func UpdateActivity() {

}

func DeleteActivity() {

}

//登陆
func Login(req *http.Request, r render.Render, db *sql.DB) {
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

//注册
func Register(req *http.Request, r render.Render, db *sql.DB) {
	plat := req.FormValue("plat")
	if plat == "SinaWeibo" || plat == "QQ" {
		openId, accessToken, expiresIn := req.FormValue("openid"), req.FormValue("access_token"), req.FormValue("expires_in")
		username, gender, avatar := req.FormValue("name"), req.FormValue("gender"), req.FormValue("avatar")
		row := db.QueryRow("SELECT access_token FROM t_oauth WHERE plat = ? AND openid = ?", plat, openId)
		var token string
		row.Scan(&token)
		if token == "" {
			tx, err := db.Begin()
			if err != nil {
				log.Fatal(err)
			}
			stmt, err := tx.Prepare("INSERT INTO t_user(name,gender,avatar) VALUES (?,?,?)")
			defer stmt.Close()
			if err != nil {
				log.Fatal(err)
			}
			stmt2, err2 := db.Prepare("INSERT INTO t_oauth(user_id,plat,openid,access_token,expires) VALUES (?,?,?,?,?)")
			defer stmt2.Close()
			if err2 != nil {
				log.Fatal(err2)
			}

			defer tx.Rollback()
			res, err := stmt.Exec(username, gender, avatar)
			id, err := res.LastInsertId()
			expiresInInt, err := strconv.Atoi(expiresIn)
			_, err = stmt2.Exec(id, plat, openId, accessToken, time.Now().Add(time.Second*time.Duration(expiresInInt)))
			tx.Commit()
			if err != nil || err2 != nil {
				log.Println(err)
				log.Println(err2)
			} else {
				r.JSON(200, "register OK")
			}
		} else {
			r.JSON(200, "fail registered ")
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

//更新修改资料
func UpdateUser() {

}

func GetUser(params martini.Params) {
	//	id := params["name"]
}

//删除user
func DeleteUser() {

}

func main() {
	m := martini.Classic()
	m.NotFound(func() {
		// 处理 404
	})
	m.Use(func(res http.ResponseWriter, req *http.Request) {
		//		if req.Header.Get("X-API-KEY") != "secret123" {
		//			res.WriteHeader(http.StatusUnauthorized)
		//		}
	})
	m.Use(render.Renderer())

	//db
	db, err := sql.Open("mysql", "root:weiwanglive@tcp(106.75.19.205:3306)/minnetlive")
	if err != nil {
		panic(err.Error()) // TODO
	}
	defer db.Close()
	m.Map(db)

	m.Get("/", func(r render.Render) string {
		return "Hello world!"
	})

	m.Group("/activity", func(r martini.Router) {
		r.Get("/:id", GetActivity)
		r.Post("/new", NewActivity)
		r.Put("/update/:id", UpdateActivity)
		r.Delete("/delete/:id", DeleteActivity)
	})

	m.Group("/user", func(r martini.Router) {
		r.Get("/:id", GetUser)
		r.Post("/login", Login)
		//		r.Post("/register", binding.Bind(User{}), Register)
		r.Post("/register", Register)
		r.Put("/update/:id", UpdateUser)
		r.Delete("/delete/:id", DeleteUser)
	})

	m.Run()
}
