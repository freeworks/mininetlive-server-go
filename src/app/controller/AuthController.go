package controller

import (
	. "app/common"
	. "app/models"
	"time"

	"github.com/coopernurse/gorp"
	"github.com/martini-contrib/render"
)

//登陆
func LoginOAuth(oauth OAuth, r render.Render, dbmap *gorp.DbMap) {
	err := dbmap.SelectOne(&oauth, "select * from t_oauth where openid=?", oauth.OpenId)
	CheckErr(err, "LoginOAuth selectOne failed")
	if err != nil {
		r.JSON(200, Resp{1000, "未注册", nil})
	} else {
		dbmap.Update(&oauth)
		obj, err := dbmap.Get(User{}, oauth.UserId)
		CheckErr(err, "LoginOAuth get user failed")
		if obj == nil {
			r.JSON(200, Resp{1002, "用户资料信息不存在", nil})
		}
		user := obj.(*User)
		r.JSON(200, Resp{0, "登陆成功", map[string]interface{}{"token": oauth.AccessToken, "user": user}})
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
		register.User.Updated = time.Now()
		register.User.Created = time.Now()
		trans.Insert(&register.User)
		register.OAuth.Expires = time.Now().Add(time.Second * time.Duration(register.OAuth.ExpiresIn))
		register.OAuth.UserId = register.User.Id
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
		auth.Token = GeneraToken()
		auth.Expires = time.Now().Add(time.Hour * 24 * 30)
		_, err := dbmap.Update(&auth)
		CheckErr(err, "Login update auth")
		obj, err := dbmap.Get(User{}, auth.UserId)
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
		authUser.LocalAuth.UserId = authUser.User.Id
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

func Logout() {
	//TODO
}
