package admin

import (
	. "app/common"
	. "app/models"
	//	config "app/config"
	easemob "app/easemob"
	logger "app/logger"
	"app/sessionauth"
	"app/sessions"
	//	"fmt"
	//	"io"
	"net/http"
	//	"os"
	"strconv"
	"time"

	"github.com/coopernurse/gorp"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	cache "github.com/patrickmn/go-cache"
)

func Index(r render.Render) {
	logger.Debug("Index")
	r.HTML(200, "index", nil)
}

func PostLogin(req *http.Request, session sessions.Session, r render.Render, dbmap *gorp.DbMap) {
	phone := req.PostFormValue("phone")
	password := req.PostFormValue("password")
	if ValidatePhone(phone) && ValidatePassword(password) {
		logger.Info("admin-login:" + phone + " " + password)
		var admin AdminModel
		err := dbmap.SelectOne(&admin, "SELECT * FROM t_admin WHERE phone = ? AND password = ?", phone, password)
		CheckErr(err, "Login select one")
		if err != nil {
			r.JSON(401, "用户名密码错误")
			return
		} else {
			err := sessionauth.AuthenticateSession(session, &admin)
			CheckErr(err, "Login AuthenticateSession")
			if err != nil {
				r.JSON(406, err)
			}
			logger.Info(req.URL)
			redirectParams := req.URL.Query()[sessionauth.RedirectParam]
			logger.Info("redirectParams", redirectParams)
			var redirectPath string
			if len(redirectParams) > 0 && redirectParams[0] != "null" {
				redirectPath = redirectParams[0]
			} else {
				redirectPath = "/"
			}
			r.JSON(200, redirectPath)
			return
		}
	} else {
		r.JSON(406, "账号，密码格式错误！")
	}
}

func Logout(session sessions.Session, user sessionauth.User, r render.Render) {
	sessionauth.Logout(session, user)
	r.Redirect("/")
}

func GetLogin(r render.Render) {
	r.HTML(200, "login", nil)
}

func GetVCode(req *http.Request, c *cache.Cache, r render.Render, dbmap *gorp.DbMap) {
	phone := req.PostFormValue("phone")
	if ValidatePhone(phone) {
		count, err := dbmap.SelectInt("select count(*) from t_admin where phone=?", phone)
		if err == nil && count >= 1 {
			vCode, err := SendSMS(phone)
			if err == nil {
				c.Set(phone, vCode, 60*time.Second)
				r.JSON(200, "验证码发送成功!")
				return
			}
		}
		r.JSON(500, "验证码发送失败!")
		return
	} else {
		r.JSON(406, "手机号格式错误！")
	}
}

func UpdatePassword(req *http.Request, c *cache.Cache, r render.Render, dbmap *gorp.DbMap) {
	phone := req.PostFormValue("phone")
	password := req.PostFormValue("password")
	vCode := req.PostFormValue("vCode")
	if ValidatePhone(phone) && ValidatePassword(password) {
		if cacheVCode, found := c.Get(phone); found {
			if cacheVCode.(string) == vCode {
				_, err := dbmap.Exec("UPDATE t_admin SET password=? WHERE phone=?", password, phone)
				CheckErr(err, "update password")
				if err == nil {
					r.JSON(200, "密码更新成功！")
					return
				} else {
					r.JSON(500, "系统内部错误！")
					return
				}
			} else {
				r.JSON(406, "验证码错误，请重新输入！")
			}
		} else {
			r.JSON(406, "验证码过期,请重新获取验证码！")
		}
	} else {
		r.JSON(406, "账号，密码格式错误！")
	}
}

func GetOrderList(req *http.Request, r render.Render, dbmap *gorp.DbMap) {
	start, size := GetLimit(req)
	var orders []Order
	_, err := dbmap.Select(&orders, "SELECT * FROM t_order LIMIT ?,?", start, size)
	CheckErr(err, "GetOrderList")
	if err == nil {
		//	newmap := map[string]interface{}{"orders": orders}
		//	r.HTML(200, "xxxxx", newmap)
		r.JSON(200, orders)
	} else {
		r.HTML(500, "服务器异常", nil)
	}
}

func FilterOrderList(req *http.Request, r render.Render, dbmap *gorp.DbMap) {
	//	start, size := GetTimeDatae(req)
	//	t1, e := time.Parse(
	//		time.RFC3339,
	//		"2012-11-01T22:08:41+00:00")
	//	p(t1)
	//	var orders []Order
	//	_, err := dbmap.Select(&orders, "SELECT * FROM t_order LIMIT ?,?", start, size)
	//	CheckErr(err, "GetOrderList")
	//	if err == nil {
	//		//	newmap := map[string]interface{}{"orders": orders}
	//		//	r.HTML(200, "xxxxx", newmap)
	//		r.JSON(200, orders)
	//	} else {
	//		r.HTML(500, "服务器异常", nil)
	//	}
}

func GetOrderChat() {

}

func GetIncomChart() {

}

func GetUserList(req *http.Request, r render.Render, dbmap *gorp.DbMap) {
	start, size := GetLimit(req)
	sql := `SELECT u.uid,u.easemob_uuid,u.nickname,u.gender,u.avatar,u.balance,u.create_time,auth.plat,"" as phone  
			FROM t_user u LEFT JOIN t_oauth auth ON u.uid = auth.uid  
			WHERE auth.plat != "" 
			LIMIT ?,? 
			UNION ALL 
			SELECT u.uid,u.easemob_uuid,u.nickname,u.gender,u.avatar,u.balance,u.create_time,"" as plat, auth.phone  
			FROM t_user u LEFT JOIN t_local_auth auth ON u.uid = auth.uid 
			WHERE auth.phone != "" 
			LIMIT ?,? `
	var userList []QUserModel
	_, err := dbmap.Select(&userList, sql, start, size, start, size)
	if err == nil {
		r.HTML(200, "userlist", userList)
	} else {
		logger.Info(err)
		r.HTML(500, "服务器异常", nil)
	}
}

func GetActivityList(req *http.Request, r render.Render, dbmap *gorp.DbMap) {
	start, size := GetLimit(req)
	var activities []Activity
	_, err := dbmap.Select(&activities, "SELECT * FROM t_activity LIMIT ?,?", start, size)
	CheckErr(err, "GetActivityList")
	logger.Info(activities)
	newmap := map[string]interface{}{"activities": activities}
	r.HTML(200, "activitylist", newmap)
}

func GetActivity(params martini.Params, r render.Render, dbmap *gorp.DbMap) {
	id := params["id"]
	var activity Activity
	err := dbmap.SelectOne(&activity, "SELECT * FROM t_activity WHERE aid = ?", id)
	CheckErr(err, "GetActivity")
	logger.Info(activity)
	//	newmap := map[string]interface{}{"activity": activity}
	//	r.HTML(200, "activitylist", newmap)
	r.JSON(200, activity)
}

func DeleteActivity(params martini.Params, r render.Render, dbmap *gorp.DbMap) {
	_, err := dbmap.Exec("DELETE from t_activity WHERE id=?", params["id"])
	CheckErr(err, "DeleteActivity delete failed")
	if err == nil {
		r.JSON(200, "删除活动成功")
	} else {
		r.JSON(500, "删除活动失败")
	}
}

func NewActivity(activity NActivity, user sessionauth.User, r render.Render, c *cache.Cache, dbmap *gorp.DbMap) {
	logger.Info("NewActivity ")
	// uid := user.UniqueId().(string)
	uid := "1e046709049d59b5"
	groupId, err := easemob.CreateGroup(uid, activity.Title, c)
	if err != nil {
		CheckErr(err, "easemob create group error")
		r.JSON(500, "创建活动失败")
		return
	}
	activity.Aid = AID()
	activity.Uid = uid
	activity.Date = JsonTime{time.Unix(activity.ADate, 0), true}
	activity.GroupId = groupId
	activity.StreamId = GeneraToken8()
	activity.LivePushPath = generatePushPath(activity.StreamId, activity.IsRecord, "")
	activity.LivePullPath = generatePullPath(activity.StreamId)
	logger.Info("info ", activity.String())
	err = dbmap.Insert(&activity)
	CheckErr(err, "NewActivity insert failed")
	if err == nil {
		r.JSON(200, "/activity")
	} else {
		//TODO 删除环信id
		r.JSON(500, "创建活动失败")
	}
}

func generatePushPath(streamId string, record bool, filename string) string {
	pushPath := "rtmp://publish.weiwanglive.com/mininetlive/" + streamId + "?record=" + strconv.FormatBool(record)
	if filename != "" {
		pushPath = pushPath + "&filename=" + filename
	}
	logger.Info("GeneratePushPath :", pushPath)
	return pushPath
}

func generatePullPath(streamId string) string {
	pullPath := "rtmp://rtmp.weiwanglive.com/mininetlive/" + streamId
	logger.Info("GeneratePullPath :", pullPath)
	return pullPath
}

func UpdateActivity(params martini.Params, activity NActivity, r render.Render, dbmap *gorp.DbMap) {
	obj, err := dbmap.Get(Activity{}, params["id"])
	CheckErr(err, "UpdateActivity get Activity err ")
	if err != nil {
		r.JSON(200, Resp{1101, "更新活动失败", nil})
	} else {
		orgActivity := obj.(*Activity)
		orgActivity.Title = activity.Title
		orgActivity.Date = activity.Date
		orgActivity.Desc = activity.Desc
		orgActivity.ActivityType = activity.ActivityType
		orgActivity.StreamType = activity.StreamType
		orgActivity.FrontCover = activity.FrontCover
		_, err := dbmap.Update(orgActivity)
		CheckErr(err, "UpdateActivity  update failed")
		if err != nil {
			r.JSON(200, Resp{1101, "更新活动失败", nil})
		} else {
			r.JSON(200, Resp{0, "更新活动成功", activity})
		}
	}
}
