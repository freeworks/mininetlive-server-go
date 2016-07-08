package admin

import (
	. "app/common"
	//	config "app/config"
	//easemob "app/easemob"
	logger "app/logger"
	. "app/models"
	"app/sessionauth"
	"app/sessions"
	//	"fmt"
	//	"io"
	"net/http"
	//	"os"
	//	"strconv"
	//	"time"

	"github.com/coopernurse/gorp"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	//	cache "github.com/patrickmn/go-cache"
)

func Index(r render.Render) {
	logger.Debug("Index")
	r.HTML(200, "index", nil)
}

func PostLogin(args martini.Params, req *http.Request, session sessions.Session, r render.Render, dbmap *gorp.DbMap) {
	req.ParseForm()
	phone := req.PostFormValue("phone")
	password := req.PostFormValue("password")
	if ValidatePhone(phone) && ValidatePassword(password) {
		logger.Info("admin-login:" + phone + " " + password)
		var admin AdminModel
		err := dbmap.SelectOne(&admin, "SELECT * FROM t_admin WHERE phone = ? AND password = ?", phone, password)
		CheckErr(err, "Login select one")
		if err != nil {
			r.JSON(500, "用户名密码错误")
			return
		} else {
			err := sessionauth.AuthenticateSession(session, &admin)
			CheckErr(err, "Login AuthenticateSession")
			if err != nil {
				r.JSON(500, err)
			}
			logger.Info(req.URL)
			redirectParams := req.URL.Query()[sessionauth.RedirectParam]
			logger.Info(redirectParams)
			var redirectPath string
			if len(redirectParams) > 0 {
				redirectPath = redirectParams[0]
			} else {
				redirectPath = "/"
			}
			r.JSON(200, redirectPath)
			return
		}
	} else {
		r.JSON(500, "账号或密码错误！")
	}
}

func Logout(session sessions.Session, user sessionauth.User, r render.Render) {
	sessionauth.Logout(session, user)
	r.Redirect("/")
}

func GetLogin(r render.Render) {
	r.HTML(200, "login", nil)
}

func UpdatePassword() {

}

func GetOrderList() {

}

func QueryOrderList() {

}

func GetOrderChat() {

}

func GetIncomChart() {

}

func GetUserList(r render.Render, dbmap *gorp.DbMap) {
	var thiredUserModel []ThiredUserModel
	_, err := dbmap.Select(&thiredUserModel, "SELECT t_user.id,t_user.name,gender,avatar,balance,update_time,create_time,plat FROM t_user,t_oauth WHERE t_user.id = t_oauth.user_id")
	CheckErr(err, "GetUserlist select failed")
	var phoneUserModel []PhoneUserModel
	_, err = dbmap.Select(&phoneUserModel, "SELECT t_user.id , t_user.name,gender,avatar,balance,update_time,create_time,phone FROM t_user,t_local_auth WHERE t_user.id = t_local_auth.user_id")
	CheckErr(err, "GetUserlist select failed")
	if err == nil {
		newmap := map[string]interface{}{"thiredUserModel": thiredUserModel, "phoneUserModel": phoneUserModel}
		r.HTML(200, "userlist", newmap)
	} else {
		r.HTML(500, "userlist", nil)
	}
}

func GetActivityList() {

}

func GetActivity() {

}

func NewActivity() {

}

func UpdateActivity() {

}

func DeleteActivity() {

}
