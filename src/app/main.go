package main

import (
	admin "app/admin"
	. "app/controller"
	db "app/db"
	. "app/models"
	sessionauth "app/sessionauth"
	sessions "app/sessions"
	. "app/upload"

	"github.com/go-martini/martini"
	_ "github.com/go-sql-driver/mysql"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
)

func main() {
	//setup db
	dbmap := db.InitDb()
	defer dbmap.Db.Close()
	m := martini.Classic()
	m.Map(dbmap)
	m.Group("/user", func(r martini.Router) {
		r.Post("/oauth/login", binding.Bind(OAuth{}), LoginOAuth)
		r.Post("/oauth/register", binding.Bind(OAuthUser{}), RegisterOAuth)
		r.Post("/login", binding.Bind(LocalAuth{}), Login)
		r.Post("/register", binding.Bind(LocalAuthUser{}), Register)
		r.Get("/:id", GetUser)
		r.Put("/:id", binding.Bind(User{}), UpdateUser)
		r.Delete("/:id", DeleteUser)
	})
	m.Group("/activity", func(r martini.Router) {
		r.Post("", binding.Bind(Activity{}), NewActivity)
		r.Get("", GetAllActivity)
		r.Get("/:id", GetActivity)
		r.Put("/:id", binding.Bind(Activity{}), UpdateActivity)
		r.Delete("/:id", DeleteActivity)
		r.Post("/appointment/:id", AppointmentActivity)
		r.Delete("/appointment", CancelAppointmentActivity)
		r.Post("/pay/:id", PayActivity)
		r.Post("/play/:id", PlayActivity)
	})
	m.Get("appointmentrecords", GetAppointmentRecords)
	m.Get("playrecords", GetPlayRecords)
	m.Get("payrecords", GetPayRecords)
	m.Post("/upload", Upload)

	go func() {
		admin.SetDBMap(dbmap)
		m := martini.Classic()
		m.Map(dbmap)
		m.NotFound(func() {
			// 处理 404
		})
		m.Use(render.Renderer(render.Options{
			Directory: "templates", // Specify what path to load the templates from.
			//Layout:     "layout",                   // Specify a layout template. Layouts can call {{ yield }} to render the current template.
			Extensions: []string{".tmpl", ".html"}, // Specify extensions to load for templates.
			//		Funcs:           []template.FuncMap{AppHelpers}, // Specify helper function maps for templates to access.
			Delims:     render.Delims{"{[{", "}]}"}, // Sets delimiters to the specified strings.
			Charset:    "UTF-8",                     // Sets encoding for json and html content-types. Default is "UTF-8".
			IndentJSON: true,                        // Output human readable JSON
		}))

		m.Use(sessions.Sessions("my_session", []byte("secret123")))
		m.Use(sessionauth.SessionUser(admin.GenerateAnonymousUser))
		sessionauth.RedirectUrl = "/login"
		sessionauth.RedirectParam = "next"

		m.Get("/", sessionauth.LoginRequired, admin.Index)
		m.Post("/login", admin.Login)
		m.Get("/login", admin.RedirectLogin)
		m.Get("/logout", sessionauth.LoginRequired, admin.Logout)
		m.Get("/activity", sessionauth.LoginRequired, admin.GetActivityList)
		m.Post("/activity", sessionauth.LoginRequired, binding.Bind(Activity{}), admin.NewActivity)
		m.Get("/addactivity", sessionauth.LoginRequired, admin.AddActivity)
		m.Get("/admin", sessionauth.LoginRequired, admin.GetAdminList)
		m.Get("/user", sessionauth.LoginRequired, admin.GetUserList)
		m.Get("/income", sessionauth.LoginRequired, admin.GetIncome)
		m.RunOnAddr(":8081")
	}()

	m.RunOnAddr(":8080")
}
