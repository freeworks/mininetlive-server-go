package controller

import (
	. "app/common"
	logger "app/logger"
	. "app/models"
	"net/http"
	"time"

	"github.com/coopernurse/gorp"
	"github.com/martini-contrib/render"
	cache "github.com/patrickmn/go-cache"
)

//登陆
func LoginOAuth(oauth OAuth, r render.Render, dbmap *gorp.DbMap) {
	logger.Info("LoginOAuth....")
	err := dbmap.SelectOne(&oauth, "select * from t_oauth where openid=?", oauth.OpenId)
	CheckErr(err, "LoginOAuth select oauth")
	if err != nil {
		r.JSON(200, Resp{1000, "未注册", nil})
	} else {
		dbmap.Update(&oauth)
		var user User
		err := dbmap.SelectOne(&user, "select * from t_user where uid=?", oauth.Uid)
		CheckErr(err, "LoginOAuth select user ")
		if err != nil {
			r.JSON(200, Resp{1002, "用户资料信息不存在", nil})
		} else {
			r.JSON(200, Resp{0, "登陆成功", map[string]interface{}{"token": oauth.AccessToken, "user": user}})
		}
	}
}

func RegisterOAuth(register OAuthUser, r render.Render, dbmap *gorp.DbMap) {
	var oauth OAuth
	err := dbmap.SelectOne(&oauth, "select * from t_oauth where openid=?", register.OAuth.OpenId)
	CheckErr(err, "RegisterOAuth selectOne failed")
	if err != nil && oauth.OpenId != register.OAuth.OpenId {
		// err = dbmap.Insert(&register.User)
		// err = dbmap.Insert(&register.OAuth)
		trans, err := dbmap.Begin()
		CheckErr(err, "RegisterOAuth begin trans"+register.User.String())
		register.User.Created = time.Now()
		register.User.Uid = GeneraToken16()
		trans.Insert(&register.User)
		register.OAuth.Expires = time.Now().Add(time.Second * time.Duration(register.OAuth.ExpiresIn))
		register.OAuth.Uid = register.User.Uid
		trans.Insert(&register.OAuth)
		err = trans.Commit()
		CheckErr(err, "RegisterOAuth trans commit ")
		if err == nil {
			r.JSON(200, Resp{0, "注册成功", map[string]interface{}{"user": register.User}})
		} else {
			r.JSON(200, Resp{1003, "注册失败，服务器异常", nil})
		}
	} else {
		r.JSON(200, Resp{1004, "该账号已经注册", nil})
	}
}

func Login(localAuth LocalAuth, r render.Render, dbmap *gorp.DbMap) {
	var auth LocalAuth
	err := dbmap.SelectOne(&auth, "select * from t_local_auth where phone=? and password = ?", localAuth.Phone, localAuth.Password)
	CheckErr(err, "Login selectOne failed")
	if err != nil {
		r.JSON(200, Resp{1005, "账户或密码错误", nil})
	} else {
		auth.Token = GeneraToken16()
		auth.Expires = time.Now().Add(time.Hour * 24 * 30)
		_, err := dbmap.Update(&auth)
		CheckErr(err, "Login update auth")
		obj, err := dbmap.Get(User{}, auth.Uid)
		CheckErr(err, "Login get user")
		if obj == nil {
			r.JSON(200, Resp{1002, "用户资料信息不存在", nil})
		}
		user := obj.(*User)
		r.JSON(200, Resp{0, "登陆成功", map[string]interface{}{"token": auth.Token, "user": user}})
	}
}

func Register(authUser LocalAuthUser, r render.Render, dbmap *gorp.DbMap) {
	var auth LocalAuth
	err := dbmap.SelectOne(&auth, "select * from t_local_auth where phone=?", authUser.LocalAuth.Phone)
	CheckErr(err, "Register selectOne failed")
	if err != nil && auth.Phone == "" {
		// err = dbmap.Insert(&authUser.LocalAuth)
		trans, err := dbmap.Begin()
		CheckErr(err, "Register begin trans failed")
		authUser.User.Updated = time.Now()
		authUser.User.Created = time.Now()
		err = trans.Insert(&authUser.User)
		CheckErr(err, "Register insert user failed")
		authUser.LocalAuth.Uid = authUser.User.Uid
		authUser.LocalAuth.Expires = time.Now().Add(time.Hour * 24 * 30)
		err = trans.Insert(&authUser.LocalAuth)
		CheckErr(err, "Register insert LocalAuth failed")
		err = trans.Commit()
		CheckErr(err, "Register commit failed")
		if err == nil {
			r.JSON(200, Resp{0, "注册成功", map[string]interface{}{"token": auth.Token, "user": authUser.User}})
		} else {
			r.JSON(200, Resp{1003, "注册失败，服务器异常", nil})
		}
	} else {
		r.JSON(200, Resp{1004, "该账号已经注册", nil})
	}

}

func GetVCode(req *http.Request, c *cache.Cache, r render.Render, dbmap *gorp.DbMap) {
	req.ParseForm()
	phone := req.PostFormValue("phone")
	//TODO 校验
	//是否已经注册
	count, err := dbmap.SelectInt("SELECT COUNT(*) FROM t_local_auth WHERE phone=?", phone)
	CheckErr(err, "check phone is registed")
	if err != nil {
		r.JSON(200, Resp{1009, "获取验证码失败，服务器异常", nil})
	}
	if count == 0 {
		vCode, err := SendSMS(phone)
		if err != nil {
			r.JSON(200, Resp{1009, "获取验证码失败", nil})
		} else {
			c.Set(phone, vCode, 60*time.Second)
			r.JSON(200, Resp{0, "获取验证码成功", nil})
		}
	} else {
		r.JSON(200, Resp{1004, "该手机号已注册", nil})
	}
}

func VerifyPhone(req *http.Request, c *cache.Cache, r render.Render) {
	req.ParseForm()
	phone := req.PostFormValue("phone")
	vCode := req.PostFormValue("vcode")
	//TODO 校验
	if cacheVCode, found := c.Get(phone); found {
		if cacheVCode.(string) == vCode {
			r.JSON(200, Resp{0, "验证成功", nil})
		} else {
			r.JSON(200, Resp{1010, "输入验证码有误,请重新输入", nil})
		}
	} else {
		r.JSON(200, Resp{1011, "输入验证码无效,请重新获取验证码", nil})
	}
}

func Logout(req *http.Request, r render.Render, dbmap *gorp.DbMap) {
	//TODO
	r.JSON(200, Resp{0, "退出成功", nil})
}
