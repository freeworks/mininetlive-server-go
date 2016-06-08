package main

import (
	"github.com/go-martini/martini"
	_ "github.com/go-sql-driver/mysql"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
	."app/controller"
	."app/models"
	db "app/db"
	."app/upload"
	"net/http"
)


func main() {
	//setup db
	dbmap := db.InitDb()
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
		r.Post("", binding.Bind(Activity{}), NewActivity)
		r.Get("", GetActivityList)
		r.Get("/:id", GetActivity)
		r.Put("/:id", binding.Bind(Activity{}), UpdateActivity)
		r.Delete("/:id", DeleteActivity)
	})

	m.Group("/user", func(r martini.Router) {
		r.Post("/oauth/login", binding.Bind(OAuth{}), LoginOAuth)
		r.Post("/oauth/register", binding.Bind(OAuthUser{}), RegisterOAuth)
		r.Post("/login", binding.Bind(LocalAuth{}), Login)
		r.Post("/register", binding.Bind(LocalAuthUser{}), Register)
		r.Get("/:id", GetUser)
		r.Put("/:id", binding.Bind(User{}), UpdateUser)
		r.Delete("/:id", DeleteUser)
	})
	m.Post("/upload", Upload)
	m.Group("/admin", func(r martini.Router) {
		r.Get("", AdminMain)
	})

	m.Run()
}
