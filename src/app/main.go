package main

import (
	. "app/controller"
	db "app/db"
	. "app/models"
	. "app/upload"
	"net/http"

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
	m.Group("/activity", func(r martini.Router) {
		r.Post("", binding.Bind(Activity{}), NewActivity)
		r.Get("", GetActivityList)
		r.Get("/:id", GetActivity)
		r.Put("/:id", binding.Bind(Activity{}), UpdateActivity)
		r.Delete("/:id", DeleteActivity)
		r.Post("/appointment/:id", AppointmentActivity)
		r.Post("/pay/:id", PayActivity)
		r.Post("/play/:id", PlayActivity)
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

	go func() {
		m := martini.Classic()
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
			//		IndentXML:  true,                        // Output human readable XML
			//		HTMLContentType: "application/xhtml+xml",     // Output XHTML content type instead of default "text/html"
		}))
		m.Use(func(res http.ResponseWriter, req *http.Request) {
			//		if req.Header.Get("X-API-KEY") != "secret123" {
			//			res.WriteHeader(http.StatusUnauthorized)
			//		}
		})
		m.Group("/admin", func(r martini.Router) {
			r.Get("", Index)
			r.Get("/login", LoginAdmin)
		})
		m.RunOnAddr("127.0.0.1:8081")
	}()
	m.RunOnAddr("127.0.0.1:8080")
}
