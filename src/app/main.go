package main

import (
	admin "app/admin"
	. "app/controller"
	db "app/db"
	logger "app/logger"
	. "app/models"
	pay "app/pay"
	sessionauth "app/sessionauth"
	sessions "app/sessions"
	. "app/upload"

	"github.com/go-martini/martini"
	_ "github.com/go-sql-driver/mysql"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
)

func main() {
	initLog("api.log", logger.ALL, true)
	dbmap := db.InitDb()
	defer dbmap.Db.Close()
	m := martini.Classic()
	m.Use(logger.Logger())
	m.Use(render.Renderer())
	m.Map(dbmap)
	m.Post("/auth/login", binding.Bind(LocalAuth{}), Login)
	m.Post("/auth/register", binding.Bind(LocalAuthUser{}), Register)
	m.Post("/auth/logout", Logout)
	m.Post("/oauth/login", binding.Bind(OAuth{}), LoginOAuth)
	m.Post("/oauth/register", binding.Bind(OAuthUser{}), RegisterOAuth)
	m.Group("/common", func(r martini.Router) {
		r.Post("/inviteCode", PostInviteCode)
	})
	m.Post("/upload", Upload)
	m.Group("/account", func(r martini.Router) {
		r.Post("/info", GetAccountInfo)
		r.Get("/playRecordList", GetPlayRecordList)
		r.Get("/payRecordList", GetPayRecordList)
		r.Get("/appointmentRecordList", GetAppointmentRecordList)
		r.Put("/name", UpdateAccountNickName)
		r.Put("/phone", UpdateAccountPhone)
	})
	m.Group("/user", func(r martini.Router) {
		r.Get("/info/:uid", GetUser)
	})
	m.Group("/activity", func(r martini.Router) {
		r.Get("/list", GetActivityList)
		r.Get("/list/more", GetMoreActivityList)
		r.Get("/detail/:id", GetActivityDetail)
		r.Post("/appointment", AppointmentActivity)
		r.Post("/play", PlayActivity)
	})
	m.Group("/pay", func(r martini.Router) {
		r.Get("/charge", pay.GetCharge)
		r.Post("/webhook", pay.Webhook)
		r.Post("/withdraw", pay.Withdraw)
	})
	m.NotFound(func(r render.Render) {
		r.JSON(404, "接口不存在/请求方法错误")
	})

	go func() {
		initLog("admin.log", logger.ALL, true)
		admin.SetDBMap(dbmap)
		m := martini.Classic()
		m.Map(dbmap)
		m.Use(logger.Logger())
		m.Use(render.Renderer())
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
		m.Put("/activity/update/:id", binding.Bind(Activity{}), admin.UpdateActivity)
		m.Delete("/activity/:id", sessionauth.LoginRequired, admin.DeleteActivity)
		m.Get("/addactivity", sessionauth.LoginRequired, admin.AddActivity)
		m.Get("/admin", sessionauth.LoginRequired, admin.GetAdminList)
		m.Get("/user", sessionauth.LoginRequired, admin.GetUserList)
		m.Get("/income", sessionauth.LoginRequired, admin.GetIncome)
		m.Post("/upload", Upload)
		m.RunOnAddr(":8081")
	}()

	m.RunOnAddr(":8080")
}

func initLog(filename string, level logger.LEVEL, console bool) {
	logger.SetConsole(console)
	logger.SetRollingDaily(".", filename)
	logger.SetLevel(level)
}
