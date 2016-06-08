package controller

import (
		."app/common"
		."app/models"
		"github.com/martini-contrib/render"
		"github.com/coopernurse/gorp"
		"github.com/go-martini/martini"
		// "github.com/pborman/uuid"
		"time"
	)
//0    成功
//1000 未注册
//1002 已经注册但获取用户信息失败/用户信息不存在
//1003 注册失败
//1004 账号已经注册
//1005 账号/密码错误
//1006 更新账户信息失败
//1007 删除账户失败
		
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
		trans.Insert(&register.User)
		register.OAuth.Expires = time.Now().Add(time.Second * time.Duration(register.OAuth.ExpiresIn))
		register.OAuth.UserId = register.User.Id
		trans.Insert(&register.OAuth)
		err = trans.Commit()
		CheckErr(err, "RegisterOAuth trans commit ")
		if err == nil {
			r.JSON(200, Resp{0, "注册成功", map[string]interface{}{"user": register.User}})
		} else {
			r.JSON(200, Resp{1003, "注册失败", nil})
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
			r.JSON(200, Resp{1003, "注册失败", nil})
		}
	} else {
		r.JSON(200, Resp{1004, "该账号已经注册", nil})
	}

}

func UpdateUser(args martini.Params, user User, r render.Render, dbmap *gorp.DbMap) {
	obj, err := dbmap.Get(User{}, args["id"])
	CheckErr(err, "UpdateUser get failed")
	if err != nil {
		r.JSON(200, Resp{1002, "更新信息失败", nil})
	} else {
		orgUser := obj.(*User)
		if user.Name != "" {
			orgUser.Name = user.Name
		}
		if user.Avatar != "" {
			orgUser.Avatar = user.Avatar
		}
		_, err := dbmap.Update(orgUser)
		CheckErr(err, "UpdateUser update failed")
		if err != nil {
			r.JSON(200, Resp{1006, "更新信息失败", nil})
		} else {
			r.JSON(200, Resp{0, "更新信息成功", user})
		}
	}
}

func DeleteUser(args martini.Params, r render.Render, dbmap *gorp.DbMap) {
	_, err := dbmap.Exec("DELETE from t_user WHERE id=?", args["id"])
	CheckErr(err, "DeleteUser delete failed")
	if err != nil {
		r.JSON(200, Resp{1007, "删除用户失败", nil})
	} else {
		r.JSON(200, Resp{0, "删除用户成功", nil})
	}
}

func GetUser(args martini.Params, r render.Render, dbmap *gorp.DbMap) {
	var user User
	err := dbmap.SelectOne(&user, "select * from t_user where id=?", args["id"])
	CheckErr(err, "GetUser selectOne failed")
	//simple error check
	if err != nil {
		r.JSON(200, Resp{1002, "获取用户信息失败", nil})
	} else {
		r.JSON(200, Resp{0,"获取用户信息成功",user})
	}
}


