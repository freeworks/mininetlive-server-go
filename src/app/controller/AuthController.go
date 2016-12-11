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
	logger.Info("LoginOAuth....openId->", oauth.OpenId)
	err := dbmap.SelectOne(&oauth, "select * from t_oauth where openid=?", oauth.OpenId)
	CheckErr(err, "LoginOAuth select oauth")
	if err != nil {
		r.JSON(200, Resp{1000, "未注册", nil})
	} else {
		count, err := dbmap.Update(&oauth)
		CheckErr(err, "LoginOAuth update")
		logger.Info("LoginOAuth updated:", count)
		var user User
		logger.Info("LoginOAuth....u->", oauth.OpenId)
		err = dbmap.SelectOne(&user, "select * from t_user where uid=?", oauth.Uid)
		CheckErr(err, "LoginOAuth select user ")
		if err != nil {
			r.JSON(200, Resp{1002, "用户资料信息不存在", nil})
		} else {
			count, _ := dbmap.SelectInt("select count(*) from t_invite_relation where uid = ?", user.Uid)
			r.JSON(200, Resp{0, "登陆成功", map[string]interface{}{"token": oauth.AccessToken, "showInvited": count == int64(0), "user": user}})
		}
	}
}

func RegisterOAuth(register OAuthUser, r render.Render, c *cache.Cache, dbmap *gorp.DbMap) {
	logger.Info("RegisterOAuth....openId->", register.OAuth.OpenId)
	var oauth OAuth
	err := dbmap.SelectOne(&oauth, "select * from t_oauth where openid=?", register.OAuth.OpenId)
	CheckErr(err, "RegisterOAuth")
	logger.Info("RegisterOAuth....oauth.OpenId->", oauth.OpenId, ",register.OAuth.OpenId", register.OAuth.OpenId)
	if err != nil && oauth.OpenId != register.OAuth.OpenId {
		uid := UID()
		trans, err := dbmap.Begin()
		CheckErr(err, "RegisterOAuth begin trans")
		register.User.InviteCode = RandomStr(6)
		register.User.Uid = uid
		register.User.Qrcode = "http://h.hiphotos.baidu.com/image/pic/item/3bf33a87e950352a5936aa0a5543fbf2b2118b59.jpg"
		logger.Info("RegisterOAuth ", register.User.String())
		trans.Insert(&register.User)
		register.OAuth.Expires = time.Now().Add(time.Second * time.Duration(register.OAuth.ExpiresIn))
		register.OAuth.Uid = register.User.Uid
		trans.Insert(&register.OAuth)
		logger.Info("RegisterOAuth ", register.OAuth.String())
		err = trans.Commit()
		CheckErr(err, "RegisterOAuth trans commit ")
		if err == nil {
			logger.Info("RegisterOAuth token", register.OAuth.AccessToken)
			r.JSON(200, Resp{0, "注册成功", map[string]interface{}{"token": register.OAuth.AccessToken, "showInvited": true, "user": register.User}})
		} else {
			r.JSON(200, Resp{1003, "注册失败，服务器异常", nil})
		}
	} else {
		r.JSON(200, Resp{1004, "该账号已经注册", nil})
	}
}

func Login(localAuth LocalAuth, r render.Render, dbmap *gorp.DbMap) {
	logger.Info("Login....phone->", localAuth.Phone, " pwd->", localAuth.Password)
	var auth LocalAuth
	err := dbmap.SelectOne(&auth, "select * from t_local_auth where phone=? and password = ?", localAuth.Phone, localAuth.Password)
	CheckErr(err, "Login")
	if err != nil {
		r.JSON(200, Resp{1005, "账户或密码错误", nil})
	} else {
		auth.Token = GeneraToken16()
		auth.Expires = time.Now().Add(time.Hour * 24 * 30)
		count, err := dbmap.Update(&auth)
		CheckErr(err, "Login update auth")
		logger.Info("Login Updated", count)
		logger.Info("Login", auth.String())
		var user User
		err = dbmap.SelectOne(&user, "select * from t_user where t_user.uid = ?", auth.Uid)
		logger.Info("Login get user info uid", auth.Uid)
		CheckErr(err, "Login")
		if err != nil {
			r.JSON(200, Resp{1002, "账户错误，请重新注册", nil})
			//			dbmap.Delete(&auth)
		} else {
			r.JSON(200, Resp{0, "登陆成功", map[string]interface{}{"token": auth.Token, "showInvited": false, "user": user}})
		}
	}
}

func Register(authUser LocalAuthUser, r render.Render, c *cache.Cache, dbmap *gorp.DbMap) {
	logger.Info("Register....phone->", authUser.LocalAuth.Phone)
	var auth LocalAuth
	err := dbmap.SelectOne(&auth, "select * from t_local_auth where phone=?", authUser.LocalAuth.Phone)
	CheckErr(err, "Register selectOne failed")
	logger.Info("Register LocalAuth", auth.String())
	if err != nil && auth.Phone == "" {
		uid := UID()
		trans, err := dbmap.Begin()
		CheckErr(err, "Register begin trans failed")
		authUser.User.Uid = uid
		authUser.User.Phone = authUser.LocalAuth.Phone
		authUser.User.InviteCode = RandomStr(6)
		authUser.User.Qrcode = "http://h.hiphotos.baidu.com/image/pic/item/3bf33a87e950352a5936aa0a5543fbf2b2118b59.jpg"
		logger.Info(authUser.User.String())
		err = trans.Insert(&authUser.User)
		CheckErr(err, "Register insert user failed")
		authUser.LocalAuth.Token = Token()
		authUser.LocalAuth.Uid = authUser.User.Uid
		authUser.LocalAuth.Expires = time.Now().Add(time.Hour * 24 * 30)
		logger.Info(authUser.LocalAuth.String())
		err = trans.Insert(&authUser.LocalAuth)
		CheckErr(err, "Register insert LocalAuth failed")
		err = trans.Commit()
		CheckErr(err, "Register commit failed")
		if authUser.BeInviteCode != "" {
			logger.Info("Register insert inviteRelationShip")
			err = dbmap.Insert(&InviteRelationship{
				Uid:           uid,
				BeInvitedCode: authUser.BeInviteCode,
				Created:       JsonTime{time.Now(), true},
			})
			CheckErr(err, "Insert invited relationship ")
		}
		if err == nil {
			r.JSON(200, Resp{0, "注册成功", map[string]interface{}{"token": authUser.LocalAuth.Token, "showInvited": false, "user": authUser.User}})
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
	logger.Info("GetVCode  phone->", phone)
	if phone == "" {
		r.JSON(200, Resp{1014, "手机号不能为空", nil})
		return
	}
	//是否已经注册
	count, err := dbmap.SelectInt("SELECT COUNT(*) FROM t_local_auth WHERE phone=?", phone)
	CheckErr(err, "check phone is registed")
	logger.Info("GetVCode updated", count)
	if err != nil {
		r.JSON(200, Resp{1009, "获取验证码失败，服务器异常", nil})
		return
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
	logger.Info("VerifyPhone  phone->", phone, "vcode", vCode)
	if phone == "" {
		r.JSON(200, Resp{1014, "手机号不能为空", nil})
		return
	}
	if vCode == "" {
		r.JSON(200, Resp{1014, "验证码不能为空", nil})
		return
	}
	if cacheVCode, found := c.Get(phone); found {
		if cacheVCode.(string) == vCode {
			r.JSON(200, Resp{0, "验证成功", nil})
		} else {
			logger.Info("")
			r.JSON(200, Resp{1010, "输入验证码有误,请重新输入", nil})
		}
	} else {
		r.JSON(200, Resp{1011, "输入验证码无效,请重新获取验证码", nil})
	}
}

func Logout(req *http.Request, r render.Render, dbmap *gorp.DbMap) {
	uid := req.Header.Get("uid")
	if uid == "" {
		r.JSON(200, Resp{1013, "uid不能为空", nil})
		return
	}
	logger.Info("logout:", uid)
	r.JSON(200, Resp{0, "退出成功", nil})
}

func PostInviteCode(req *http.Request, r render.Render, dbmap *gorp.DbMap) {
	uid := req.Header.Get("uid")
	if uid == "" {
		r.JSON(200, Resp{1013, "uid不能为空", nil})
		return
	}
	req.ParseForm()
	inviteCode := req.PostFormValue("inviteCode")
	logger.Info("inviteCode:", inviteCode)
	if inviteCode != "" {
		err := dbmap.Insert(&InviteRelationship{
			Uid:           uid,
			BeInvitedCode: inviteCode,
			Created:       JsonTime{time.Now(), true},
		})
		CheckErr(err, "Insert invited relationship ")
	}
	r.JSON(200, Resp{0, "提交成功", nil})
}

func RestPassword(req *http.Request, r render.Render, dbmap *gorp.DbMap, c *cache.Cache) {
	uid := req.Header.Get("uid")
	if uid == "" {
		r.JSON(200, Resp{1013, "uid不能为空", nil})
		return
	}
	req.ParseForm()
	phone := req.PostFormValue("phone")
	vCode := req.PostFormValue("vcode")
	password := req.PostFormValue("password")
	logger.Info("VerifyPhone phone->", phone, "vcode", vCode, "password", password)
	if phone == "" {
		r.JSON(200, Resp{1014, "手机号不能为空", nil})
		return
	}
	if vCode == "" {
		r.JSON(200, Resp{1014, "验证码不能为空", nil})
		return
	}
	if password == "" {
		r.JSON(200, Resp{1014, "密码不能为空", nil})
		return
	}
	if cacheVCode, found := c.Get(phone); found {
		if cacheVCode.(string) == vCode {
			count, err := dbmap.Exec("UPDATE t_local_auth SET password = ? WHERE phone = ? AND uid = ?", password, phone, uid)
			CheckErr(err, "Update password ")
			logger.Info("RestPassword updated ", count)
			if err != nil {
				r.JSON(200, Resp{1015, "重置密码失败，服务器异常", nil})
				return
			} else {
				r.JSON(200, Resp{0, "重置密码成功", nil})
				return
			}
		} else {
			r.JSON(200, Resp{1010, "输入验证码有误,请重新输入", nil})
		}
	} else {
		r.JSON(200, Resp{1011, "输入验证码无效,请重新获取验证码", nil})
	}
}

func BindPhone(req *http.Request, r render.Render, dbmap *gorp.DbMap, c *cache.Cache) {
	uid := req.Header.Get("uid")
	if uid == "" {
		r.JSON(200, Resp{1013, "uid不能为空", nil})
		return
	}
	req.ParseForm()
	phone := req.PostFormValue("phone")
	vCode := req.PostFormValue("vcode")
	logger.Info("BindPhone phone->", phone, "vcode", vCode)
	if phone == "" {
		r.JSON(200, Resp{1014, "手机号不能为空", nil})
		return
	}
	if vCode == "" {
		r.JSON(200, Resp{1014, "验证码不能为空", nil})
		return
	}
	if cacheVCode, found := c.Get(phone); found {
		if cacheVCode.(string) == vCode {
			_, err := dbmap.SelectInt("select count(*) from t_oauth")
			if err == nil {
				trans, err := dbmap.Begin()
				if err != nil {
					CheckErr(err, "dbmap begin trans")
					r.JSON(200, Resp{1016, "服务器异常！", nil})
					return
				}
				trans.Exec("UPDATE t_local_auth SET phone = ? WHERE uid = ?", phone, uid)
				trans.Exec("UPDATE t_user SET phone = ? WHERE uid = ?", phone, uid)
				err = trans.Commit()
				CheckErr(err, "update phone ")
				if err == nil {
					r.JSON(200, Resp{0, "绑定成功", nil})
				} else {
					r.JSON(200, Resp{1016, "服务器异常！", nil})
				}
				return
			} else {
				_, err := dbmap.Exec("UPDATE t_user SET phone = ? WHERE uid = ?", phone, uid)
				CheckErr(err, "update phone ")
				if err == nil {
					r.JSON(200, Resp{0, "绑定成功", nil})
				} else {
					r.JSON(200, Resp{1016, "服务器异常！", nil})

				}
				return
			}
		} else {
			r.JSON(200, Resp{1010, "输入验证码有误,请重新输入", nil})
		}
	} else {
		r.JSON(200, Resp{1011, "输入验证码无效,请重新获取验证码", nil})
	}
}

func BindPush(req *http.Request, r render.Render, dbmap *gorp.DbMap) {
	uid := req.Header.Get("uid")
	if uid == "" {
		r.JSON(200, Resp{1013, "uid不能为空", nil})
		return
	}
	req.ParseForm()
	deviceId := req.PostFormValue("deviceId")
	if deviceId == "" {
		r.JSON(200, Resp{1017, "deviceId 为空", nil})
		return
	}
	_, err := dbmap.Exec("UPDATE t_user SET device_id= ? WHERE uid = ?", deviceId, uid)
	CheckErr(err, "BindPush")
	if err != nil {
		r.JSON(200, Resp{1018, "绑定Push失败", nil})
	} else {
		r.JSON(200, Resp{0, "绑定Push成功", nil})
	}
}
