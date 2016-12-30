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
	. "app/push"
	sessionauth "app/sessionauth"
	sessions "app/sessions"
	. "app/wxpub"
	"time"

	"github.com/go-martini/martini"
	_ "github.com/go-sql-driver/mysql"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
	cache "github.com/patrickmn/go-cache"
)

func main() {
	logger.Info(time.Now())
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
	m.Post("/auth/password/reset", RestPassword)
	m.Post("/auth/phone/bind", BindPhone)
	m.Post("/auth/verify/phone", VerifyPhone)
	m.Post("/auth/bindPush", BindPush)
	m.Post("/oauth/login", binding.Bind(OAuth{}), LoginOAuth)
	m.Post("/oauth/register", binding.Bind(OAuthUser{}), RegisterOAuth)

	m.Group("/common", func(r martini.Router) {
		r.Post("/inviteCode", PostInviteCode)
	})
	m.Group("/account", func(r martini.Router) {
		r.Post("/info", GetAccountInfo)
		r.Get("/balance", GetBalance)
		r.Get("/record/play/list", GetPlayRecordList)
		r.Get("/record/pay/list", GetPayRecordList)
		r.Get("/record/appointment/list", GetAppointmentRecordList)
		r.Get("/record/transfer/list", GetWithdrawRecordList)
		r.Get("/record/dividend/list", GetDividendRecordList)
		r.Post("/nickname", UpdateAccountNickName)
		r.Post("/vcode", GetVCodeForUpdatePhone)
		r.Post("/phone", UpdateAccountPhone)
		r.Post("/avatar", UpdateAccountAvatar)
	})
	m.Group("/user", func(r martini.Router) {
		r.Get("/info/:uid", GetUser)
	})
	m.Group("/activity", func(r martini.Router) {
		r.Get("/list", GetHomeList)
		r.Get("/list/more/:lastAid", GetMoreActivityList)
		r.Get("/live/list", GetLiveActivityList)
		r.Get("/detail/:aid", GetActivityDetail)
		r.Post("/appointment", AppointmentActivity)
		r.Post("/play", PlayActivity)
		r.Post("/join", JoinGroup)
		r.Post("/leave", LeaveGroup)
		r.Get("/member/list", GetLiveActivityMemberList)
		r.Get("/member/count", GetLiveActivityMemberCount)

	})
	m.Group("/pay", func(r martini.Router) {
		r.Post("/charge", pay.GetCharge)
		r.Post("/webhook", pay.Webhook)
		r.Post("/transfer", pay.Transfer)
	})

	m.Get("/callback/record/finish", CallbackRecordFinish)
	m.Get("/callback/live/begin", CallbackLiveBegin)
	m.Get("/callback/live/end", CallbackLiveEnd)
	m.Get("/share/:platform/activity/:id", GetSharePage)
	m.Get("/wxpub/recv", RecvWXPubMsg)
	m.Post("/wxpub/recv", RecvWXPubMsg)
	m.Post("/wxpub/vcode", GetVCodeForWxPub)
	m.Post("/wxpub/config", GetVCodeForWxPub)
	m.Post("/wxpub/bindphone", BindWxPubPhone)
	m.Post("/wxpub/jsconfig", GetConfig)
	m.NotFound(func(r render.Render) {
		r.JSON(404, "接口不存在/请求方法错误")
	})

	m.Group("/debug", func(r martini.Router) {
		r.Post("/push", TestPush)
		r.Post("/testJSConfig", GetConfig)
	})

	go intervaler.PollSyncPingxx(dbmap)

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
		m.Post("/logout", sessionauth.LoginRequired, admin.Logout)

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

	m.RunOnAddr(":80")
}
