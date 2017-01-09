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
	logger.Info("[AuthController]","[LoginOAuth]","openId->", oauth.OpenId)
	err := dbmap.SelectOne(&oauth, "select * from t_oauth where openid=?", oauth.OpenId)
	CheckErr("[AuthController]","[LoginOAuth]","select oauth",err)
	if err != nil {
		r.JSON(200, Resp{1000, "未注册", nil})
	} else {
		count, err := dbmap.Update(&oauth)
		CheckErr("[AuthController]","[LoginOAuth]","update",err)
		logger.Info("[AuthController]","[LoginOAuth]","updated:", count)
		var user User
		logger.Info("[AuthController]","[LoginOAuth]","user->", oauth.OpenId)
		err = dbmap.SelectOne(&user, "select * from t_user where uid=?", oauth.Uid)
		CheckErr("[AuthController]","[LoginOAuth]","select user",err)
		if err != nil {
			r.JSON(200, Resp{1002, "用户资料信息不存在", nil})
		} else {
			count, _ := dbmap.SelectInt("select count(*) from t_invite_relation where uid = ?", user.Uid)
			r.JSON(200, Resp{0, "登陆成功", map[string]interface{}{"token": oauth.AccessToken, "showInvited": count == int64(0), "user": user}})
		}
	}
}

func RegisterOAuth(register OAuthUser, r render.Render, c *cache.Cache, dbmap *gorp.DbMap) {
	logger.Info("[AuthController]","[RegisterOAuth]","openId->", register.OAuth.OpenId)
	var oauth OAuth
	err := dbmap.SelectOne(&oauth, "select * from t_oauth where openid=?", register.OAuth.OpenId)
	CheckErr("[AuthController]","[RegisterOAuth]","",err)
	logger.Info("[AuthController]","[RegisterOAuth]","oauth.OpenId->", oauth.OpenId, ",register.OAuth.OpenId", register.OAuth.OpenId)
	if err != nil && oauth.OpenId != register.OAuth.OpenId {
		uid := UID()
		trans, err := dbmap.Begin()
		CheckErr("[AuthController]","[RegisterOAuth]","begin trans",err)
		register.User.InviteCode = RandomStr(6)
		register.User.Uid = uid
		register.User.Qrcode = "http://h.hiphotos.baidu.com/image/pic/item/3bf33a87e950352a5936aa0a5543fbf2b2118b59.jpg"
		logger.Info("[AuthController]","[RegisterOAuth]",register.User.String())
		trans.Insert(&register.User)
		register.OAuth.Expires = time.Now().Add(time.Second * time.Duration(register.OAuth.ExpiresIn))
		register.OAuth.Uid = register.User.Uid
		trans.Insert(&register.OAuth)
		logger.Info("[AuthController]","[RegisterOAuth]", register.OAuth.String())
		err = trans.Commit()
		CheckErr("[AuthController]","[RegisterOAuth]","trans commit ",err)
		if err == nil {
			logger.Info("[AuthController]","[RegisterOAuth]","RegisterOAuth token", register.OAuth.AccessToken)
			r.JSON(200, Resp{0, "注册成功", map[string]interface{}{"token": register.OAuth.AccessToken, "showInvited": true, "user": register.User}})
		} else {
			r.JSON(200, Resp{1003, "注册失败，服务器异常", nil})
		}
	} else {
		r.JSON(200, Resp{1004, "该账号已经注册", nil})
	}
}

func Login(localAuth LocalAuth, r render.Render, dbmap *gorp.DbMap) {
	logger.Info("[AuthController]","[Login]","phone->", localAuth.Phone, " pwd->", localAuth.Password)
	var auth LocalAuth
	err := dbmap.SelectOne(&auth, "select * from t_local_auth where phone=? and password = ?", localAuth.Phone, MD5(localAuth.Password))
	CheckErr("[AuthController]","[Login]","",err)
	if err != nil {
		r.JSON(200, Resp{1005, "账户或密码错误", nil})
	} else {
		auth.Token = GeneraToken16()
		auth.Expires = time.Now().Add(time.Hour * 24 * 30)
		count, err := dbmap.Update(&auth)
		CheckErr("[AuthController]","[Login]","update auth",err)
		logger.Info("[AuthController]","[Login]","updated", count)
		logger.Info("[AuthController]","[Login]", auth.String())
		var user User
		err = dbmap.SelectOne(&user, "select * from t_user where t_user.uid = ?", auth.Uid)
		logger.Info("[AuthController]","[Login]","get user info uid", auth.Uid)
		CheckErr("[AuthController]","[Login]","",err)
		if err != nil {
			r.JSON(200, Resp{1002, "账户错误，请重新注册", nil})
			//			dbmap.Delete(&auth)
		} else {
			r.JSON(200, Resp{0, "登陆成功", map[string]interface{}{"token": auth.Token, "showInvited": false, "user": user}})
		}
	}
}

func Register(authUser LocalAuthUser, r render.Render, c *cache.Cache, dbmap *gorp.DbMap) {
	logger.Info("[AuthController]","[Register]","phone->", authUser.LocalAuth.Phone)
	var auth LocalAuth
	err := dbmap.SelectOne(&auth, "select * from t_local_auth where phone=?", authUser.LocalAuth.Phone)
	CheckErr("[AuthController]","[Register]","selectOne failed",err)
	logger.Info("[AuthController]","[Register]","LocalAuth", auth.String())
	if err != nil && auth.Phone == "" {
		uid := UID()
		trans, err := dbmap.Begin()
		CheckErr("[AuthController]","[Register]","begin trans failed",err)
		authUser.User.Uid = uid
		authUser.User.Phone = authUser.LocalAuth.Phone
		authUser.User.InviteCode = RandomStr(6)
		authUser.User.Qrcode = "http://h.hiphotos.baidu.com/image/pic/item/3bf33a87e950352a5936aa0a5543fbf2b2118b59.jpg"
		logger.Info("[AuthController]","[Register]",authUser.User.String())
		err = trans.Insert(&authUser.User)
		CheckErr("[AuthController]","[Register]","insert user failed",err)
		authUser.LocalAuth.Token = Token()
		authUser.LocalAuth.Uid = authUser.User.Uid
		authUser.LocalAuth.Expires = time.Now().Add(time.Hour * 24 * 30)
		authUser.LocalAuth.Password = MD5(authUser.LocalAuth.Password)
		logger.Info("[AuthController]","[Register]",authUser.LocalAuth.String())
		err = trans.Insert(&authUser.LocalAuth)
		CheckErr("[AuthController]","[Register]","insert LocalAuth failed",err)
		err = trans.Commit()
		CheckErr("[AuthController]","[Register]","commit failed",err)
		if authUser.BeInviteCode != "" {
			logger.Info("[AuthController]","[Register]","insert inviteRelationShip")
			err = dbmap.Insert(&InviteRelationship{
				Uid:           uid,
				BeInvitedCode: authUser.BeInviteCode,
				Created:       JsonTime{time.Now(), true},
			})
			CheckErr("[AuthController]","[Register]","Insert invited relationship",err)
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
	logger.Info("[AuthController]","[GetVCode]","phone->", phone)
	if phone == "" {
		r.JSON(200, Resp{1014, "手机号不能为空", nil})
		return
	}
	//是否已经注册
	count, err := dbmap.SelectInt("SELECT COUNT(*) FROM t_local_auth WHERE phone=?", phone)
	CheckErr("[AuthController]","[GetVCode]","check phone is registed",err)
	logger.Info("[AuthController]","[GetVCode]","updated", count)
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
	logger.Info("[AuthController]","[VerifyPhone]","phone->", phone, "vcode", vCode)
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
			logger.Info("[AuthController]","[VerifyPhone]","vCode err")
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
	logger.Info("[AuthController]","[Logout]", uid)
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
	logger.Info("[AuthController]","[PostInviteCode]","inviteCode:", inviteCode)
	if inviteCode != "" {
		err := dbmap.Insert(&InviteRelationship{
			Uid:           uid,
			BeInvitedCode: inviteCode,
			Created:       JsonTime{time.Now(), true},
		})
		CheckErr("[AuthController]","[PostInviteCode]", "Insert invited relationship",err)
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
	logger.Info("[AuthController]","[RestPassword]","phone->", phone, "vcode", vCode, "password", password)
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
			CheckErr("[AuthController]","[RestPassword]", "Update password",err)
			logger.Info("[AuthController]","[RestPassword]","updated ", count)
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
	logger.Info("[AuthController]","[BindPhone]","phone->", phone, "vcode", vCode)
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
					CheckErr("[AuthController]","BindPhone","dbmap begin trans",err)
					r.JSON(200, Resp{1016, "服务器异常！", nil})
					return
				}
				trans.Exec("UPDATE t_local_auth SET phone = ? WHERE uid = ?", phone, uid)
				trans.Exec("UPDATE t_user SET phone = ? WHERE uid = ?", phone, uid)
				err = trans.Commit()
				CheckErr("[AuthController]","BindPhone","update phone",err)
				if err == nil {
					r.JSON(200, Resp{0, "绑定成功", nil})
				} else {
					r.JSON(200, Resp{1016, "服务器异常！", nil})
				}
				return
			} else {
				_, err := dbmap.Exec("UPDATE t_user SET phone = ? WHERE uid = ?", phone, uid)
				CheckErr("[AuthController]","BindPhone","update phone",err)
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
	CheckErr("[AuthController]","BindPush","",err)
	if err != nil {
		r.JSON(200, Resp{1018, "绑定Push失败", nil})
	} else {
		r.JSON(200, Resp{0, "绑定Push成功", nil})
	}
}
