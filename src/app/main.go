package main

import (
	admin "app/admin"
	config "app/config"
	. "app/controller"
	db "app/db"
	intervaler "app/intervaler"
	logger "app/logger"
	. "app/models"
	pay "app/pay"
	push "app/push"
	sessionauth "app/sessionauth"
	sessions "app/sessions"
	"time"

	"github.com/go-martini/martini"
	_ "github.com/go-sql-driver/mysql"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
	cache "github.com/patrickmn/go-cache"
)

func main() {
	logger.SetConsole(true)
	logger.SetRollingDaily(config.LogDir, "mininetlive.log")
	logger.SetLevel(logger.ALL)
	dbmap := db.InitDb()
	defer dbmap.Db.Close()
	c := cache.New(cache.NoExpiration, 30*time.Second)
	m := martini.Classic()
	m.Map(dbmap)
	m.Map(c)
	m.Use(logger.Logger())
	m.Use(render.Renderer()) // 解析模板，默认路径为templates
	m.Post("/auth/login", binding.Bind(LocalAuth{}), Login)
	m.Post("/auth/register", binding.Bind(LocalAuthUser{}), Register)
	m.Post("/auth/logout", Logout)
	m.Post("/auth/vcode", GetVCode)
	m.Post("/auth/verify/phone", VerifyPhone)
	m.Post("/oauth/login", binding.Bind(OAuth{}), LoginOAuth)
	m.Post("/oauth/register", binding.Bind(OAuthUser{}), RegisterOAuth)
	m.Group("/common", func(r martini.Router) {
		r.Post("/inviteCode", PostInviteCode)
	})
	m.Group("/account", func(r martini.Router) {
		r.Post("/info", GetAccountInfo)
		r.Get("/record/play/list", GetPlayRecordList)
		r.Get("/record/pay/list", GetPayRecordList)
		r.Get("/record/appointment/list", GetAppointmentRecordList)

		r.Post("/nickname", UpdateAccountNickName)
		r.Post("/vcode", GetVCodeForUpdatePhone)
		r.Post("/phone", UpdateAccountPhone)
		r.Post("/avatar", UploadAccountAvatar)
	})
	m.Group("/user", func(r martini.Router) {
		r.Get("/info/:uid", GetUser)
	})
	m.Group("/activity", func(r martini.Router) {
		r.Get("/list", GetHomeList)
		r.Get("/list/more/:lastAid", GetMoreActivityList)
		r.Get("/live/list", GetLiveActivityList)
		r.Get("/detail/:id", GetActivityDetail)
		r.Post("/appointment", AppointmentActivity)
		r.Post("/play", PlayActivity)
		r.Post("/group/join", JoinGroup)
		r.Post("/group/leave", LeaveGroup)
		r.Post("/group/member/list", GetLiveActivityMemberList)
		r.Post("/group/member/count", GetLiveActivityMemberCount)

	})
	m.Group("/pay", func(r martini.Router) {
		r.Post("/charge", pay.GetCharge)
		r.Post("/webhook", pay.Webhook)
		r.Post("/withdraw", pay.Transfer)
	})
	m.Get("/live/CallbackRecordFinish", CallbackRecordFinish)
	m.Get("/live/CallbackLiveBegin", CallbackLiveBegin)
	m.Get("/live/CallbackLiveEnd", CallbackLiveEnd)
	m.Get("/test/push/:type", push.TestPush)

	m.NotFound(func(r render.Render) {
		r.JSON(404, "接口不存在/请求方法错误")
	})
	go intervaler.PollGroupOnlineUser(c, dbmap)

	go func() {
		m := martini.Classic() // 默认配置静态目录public
		c := cache.New(cache.NoExpiration, 30*time.Second)
		m.Map(c)
		m.Map(dbmap)
		m.Use(logger.Logger())
		m.Use(render.Renderer())
		m.Use(sessions.Sessions("my_session", []byte("secret123")))
		m.Use(sessionauth.SessionUser(admin.GenerateAnonymousUser))
		sessionauth.RedirectUrl = "/login"
		sessionauth.RedirectParam = "next"

		m.Post("/", sessionauth.LoginRequired, admin.Index)

		m.Post("/login", admin.PostLogin)
		m.Get("/login", admin.GetLogin)
		m.Post("/getVCode", admin.GetVCode)
		m.Post("/password/update", admin.UpdatePassword)
		m.Get("/logout", sessionauth.LoginRequired, admin.Logout)

		m.Post("/upload", admin.Upload)

		m.Group("/activity", func(r martini.Router) {
			r.Get("/list", admin.GetActivityList)
			r.Get("/detail/:id", admin.GetActivity)
			r.Post("/new", binding.Bind(NActivity{}), admin.NewActivity)
			r.Put("/update/:id", binding.Bind(NActivity{}), admin.UpdateActivity)
			r.Delete("/delete/:id", admin.DeleteActivity)
		}, sessionauth.LoginRequired)

		m.Get("/order/list", sessionauth.LoginRequired, admin.GetOrderList)
		m.Get("/order/chart/:graph", admin.GetOrderChat)
		m.Get("/income/chart/:graph", admin.GetIncomChart)

		m.Get("/user/list", sessionauth.LoginRequired, admin.GetUserList)

		m.RunOnAddr(":8081")
	}()

	m.RunOnAddr(":8080")

}
