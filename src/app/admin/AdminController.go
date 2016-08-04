package admin

import (
	. "app/common"
	config "app/config"
	logger "app/logger"
	. "app/models"
	. "app/push"
	"app/sessionauth"
	"app/sessions"
	//	"fmt"
	upload "app/upload"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/coopernurse/gorp"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	cache "github.com/patrickmn/go-cache"
)

func Index(r render.Render, dbmap *gorp.DbMap) {
	newAmount, err := dbmap.SelectInt("SELECT sum(amount) FROM t_order Where create_time > ?", time.Now().Format("2006-01-02 00:00:00"))
	CheckErr(err, "index ")
	newOrderCount, err := dbmap.SelectInt("SELECT count(*) FROM t_order Where create_time > ?", time.Now().Format("2006-01-02 00:00:00"))
	CheckErr(err, "index ")
	newUserCount, err := dbmap.SelectInt("SELECT count(*) FROM t_user Where create_time > ?", time.Now().Format("2006-01-02 00:00:00"))
	CheckErr(err, "index ")
	newmap := map[string]interface{}{"newAmount": newAmount, "newOrderCount": newOrderCount, "newUserCount": newUserCount}
	r.JSON(200, Resp{0, "首页获取成功", newmap})
}

func PostLogin(req *http.Request, session sessions.Session, r render.Render, dbmap *gorp.DbMap) {
	req.ParseForm()
	phone := req.PostFormValue("phone")
	password := req.PostFormValue("password")
	if ValidatePhone(phone) && ValidatePassword(password) {
		logger.Info("admin-login:" + phone + " " + password)
		var admin AdminModel
		err := dbmap.SelectOne(&admin, "SELECT * FROM t_admin WHERE phone = ? AND password = ?", phone, password)
		CheckErr(err, "Login select one")
		if err != nil {
			r.JSON(200, Resp{1021, "用户名密码错误!", nil})
			return
		} else {
			err := sessionauth.AuthenticateSession(session, &admin)
			CheckErr(err, "Login AuthenticateSession")
			if err != nil {
				r.JSON(200, Resp{1022, "校验失败!", nil})
				return
			}
			logger.Info(req.URL)
			redirectParams := req.URL.Query()[sessionauth.RedirectParam]
			logger.Info("redirectParams", redirectParams)

			var redirectPath string
			if len(redirectParams) > 0 && redirectParams[0] != "null" {
				redirectPath = redirectParams[0]
			} else {
				redirectPath = "/index.html"
			}
			r.JSON(200, Resp{0, "登陆成功!", map[string]interface{}{"redirectPath": redirectPath}})
			return
		}
	} else {
		r.JSON(200, Resp{1023, "账号，密码格式错误！", nil})
	}
}

func Logout(session sessions.Session, user sessionauth.User, r render.Render) {
	r.JSON(200, Resp{0, "退出成功", nil})
}

func GetLogin(r render.Render) {
	r.HTML(200, "login", nil)
}

func GetVCode(req *http.Request, c *cache.Cache, r render.Render, dbmap *gorp.DbMap) {
	req.ParseForm()
	phone := req.PostFormValue("phone")
	if ValidatePhone(phone) {
		count, err := dbmap.SelectInt("select count(*) from t_admin where phone=?", phone)
		if err == nil && count >= 1 {
			vCode, err := SendSMS(phone)
			if err == nil {
				c.Set(phone, vCode, 60*time.Second)
				r.JSON(200, Resp{0, "验证码发送成功!", nil})
				return
			}
		}
		r.JSON(200, Resp{1016, "验证码发送失败", nil})
		return
	} else {
		r.JSON(200, Resp{1015, "手机号格式错误！", nil})
	}
}

func UpdatePassword(req *http.Request, c *cache.Cache, r render.Render, dbmap *gorp.DbMap) {
	req.ParseForm()
	phone := req.PostFormValue("phone")
	password := req.PostFormValue("password")
	vCode := req.PostFormValue("vCode")
	if ValidatePhone(phone) && ValidatePassword(password) {
		if cacheVCode, found := c.Get(phone); found {
			if cacheVCode.(string) == vCode {
				_, err := dbmap.Exec("UPDATE t_admin SET password=? WHERE phone=?", password, phone)
				CheckErr(err, "update password")
				if err == nil {
					r.JSON(200, Resp{0, "密码更新成功！", nil})
					return
				} else {
					r.JSON(200, Resp{1014, "系统内部错误！", nil})
					return
				}
			} else {
				r.JSON(200, Resp{1013, "验证码错误，请重新输入！", nil})
			}
		} else {
			r.JSON(200, Resp{1012, "验证码过期,请重新获取验证码！", nil})
		}
	} else {
		r.JSON(200, Resp{1011, "账号，密码格式错误！", nil})
	}
}

func GetOrderList(req *http.Request, r render.Render, dbmap *gorp.DbMap) {
	queryType := req.URL.Query()["type"]
	start, size := GetLimit(req)
	if len(queryType) > 0 && queryType[0] == "filter" {
		beginDate, endDate := GetTimesampe(req)
		logger.Info("FilterOrderList", beginDate, endDate)
		var orders []Order
		var mErr error
		var mTotalCount, mTotalPageCount int
		if beginDate == "" && endDate == "" {
			_, err := dbmap.Select(&orders, "SELECT * FROM t_order LIMIT ?,?", start, size)
			totalCount, err := dbmap.SelectInt("select count(*) from t_order")
			CheckErr(err, "GetOrderList")
			m := int(totalCount) % size
			totalPageCount := int(totalCount) / size
			if m != 0 {
				totalPageCount = totalPageCount + 1
			}
			mTotalCount = int(totalCount)
			mTotalPageCount = totalPageCount
			mErr = err
		} else if beginDate != "" && endDate != "" {
			_, err := dbmap.Select(&orders, "SELECT * FROM t_order WHERE create_time >= ? AND create_time <= ? LIMIT ?,?",
				beginDate,
				endDate, start, size)
			totalCount, err := dbmap.SelectInt("SELECT count(*) FROM t_order WHERE create_time >= ? AND create_time <= ? ", beginDate,
				endDate)
			CheckErr(err, "GetOrderList")
			m := int(totalCount) % size
			totalPageCount := int(totalCount) / size
			if m != 0 {
				totalPageCount = totalPageCount + 1
			}
			mTotalCount = int(totalCount)
			mTotalPageCount = totalPageCount
			mErr = err
		} else if beginDate == "" {
			_, err := dbmap.Select(&orders, "SELECT * FROM t_order WHERE  create_time <= ? LIMIT ?,?",
				endDate, start, size)
			totalCount, err := dbmap.SelectInt("SELECT count(*) FROM t_order WHERE  create_time <= ? ", endDate)
			CheckErr(err, "GetOrderList")
			m := int(totalCount) % size
			totalPageCount := int(totalCount) / size
			if m != 0 {
				totalPageCount = totalPageCount + 1
			}
			mTotalCount = int(totalCount)
			mTotalPageCount = totalPageCount
			mErr = err
		} else if endDate == "" {
			_, err := dbmap.Select(&orders, "SELECT * FROM t_order WHERE  create_time >= ? LIMIT ?,?",
				beginDate, start, size)
			totalCount, err := dbmap.SelectInt("SELECT count(*) FROM t_order WHERE  create_time >= ? ", beginDate)
			CheckErr(err, "GetOrderList")
			m := int(totalCount) % size
			totalPageCount := int(totalCount) / size
			if m != 0 {
				totalPageCount = totalPageCount + 1
			}
			mTotalCount = int(totalCount)
			mTotalPageCount = totalPageCount
			mErr = err
		}
		CheckErr(mErr, "GetOrderList")
		if mErr == nil {
			newmap := map[string]interface{}{
				"totalCount":     mTotalCount,
				"totalPageCount": mTotalPageCount,
				"orderList":      orders}
			r.JSON(200, Resp{0, "订单列表查询成功", newmap})
		} else {
			r.JSON(200, Resp{1009, "订单列表查询失败", nil})
		}
	} else {
		var orders []Order
		_, err := dbmap.Select(&orders, "SELECT * FROM t_order LIMIT ?,?", start, size)
		CheckErr(err, "GetOrderList")
		totalCount, err := dbmap.SelectInt("select count(*) from t_order")
		CheckErr(err, "GetOrderList")
		m := int(totalCount) % size
		totalPageCount := int(totalCount) / size
		if m != 0 {
			totalPageCount = totalPageCount + 1
		}
		newmap := map[string]interface{}{
			"totalCount":     totalCount,
			"totalPageCount": totalPageCount,
			"orderList":      orders}
		if err == nil {
			r.JSON(200, Resp{0, "获取订单列表查询成功", newmap})
		} else {
			r.JSON(200, Resp{1010, "获取订单列表失败", nil})
		}
	}
}

func GetOrderChat(r render.Render, dbmap *gorp.DbMap) {
	//beginDate, endDate := GetTimesampe(req)
	sql := "SELECT DATE_FORMAT(create_time,'%Y/%m/%d') date,count(id) count FROM t_order GROUP BY date"
	var result []Graph
	_, err := dbmap.Select(&result, sql)
	if err == nil {
		r.JSON(200, Resp{0, "成功", result})
	} else {
		r.JSON(200, Resp{1017, "失败", nil})
	}
}

func GetIncomChart(r render.Render, dbmap *gorp.DbMap) {
	//beginDate, endDate := GetTimesampe(req)
	sql := "SELECT DATE_FORMAT(create_time,'%Y/%m/%d') date,SUM(amount) count FROM t_order GROUP BY date"
	var result []Graph
	_, err := dbmap.Select(&result, sql)
	if err == nil {
		r.JSON(200, Resp{0, "成功", result})
	} else {
		r.JSON(200, Resp{1018, "失败", nil})
	}
}

func GetUserList(req *http.Request, r render.Render, dbmap *gorp.DbMap) {
	start, size := GetLimit(req)
	sql := `SELECT u.uid,u.nickname,u.gender,u.avatar,u.balance,u.create_time,auth.plat,"" as phone,
			(SELECT COUNT(*) FROM t_invite_relation WHERE be_invited_code = u.invite_code) as inviteCount
			FROM t_user u JOIN t_oauth auth ON u.uid = auth.uid
			UNION ALL 
			SELECT u.uid,u.nickname,u.gender,u.avatar,u.balance,u.create_time,"" as plat, auth.phone,
			(SELECT COUNT(*) FROM t_invite_relation WHERE be_invited_code = u.invite_code) as inviteCount
			FROM t_user u JOIN t_local_auth auth ON u.uid = auth.uid 
			LIMIT ?,?`
	var userList []QUserModel
	_, err := dbmap.Select(&userList, sql, start, size)
	totalCount, err := dbmap.SelectInt("select count(*) from t_user")
	m := int(totalCount) % size
	totalPageCount := int(totalCount) / size
	if m != 0 {
		totalPageCount = totalPageCount + 1
	}
	newmap := map[string]interface{}{
		"totalCount":     totalCount,
		"totalPageCount": totalPageCount,
		"userList":       userList}
	if err == nil {
		r.JSON(200, Resp{0, "获取用户列表成功", newmap})
	} else {
		r.JSON(200, Resp{1008, "获取用户列表失败", nil})
	}
}

func GetActivityList(req *http.Request, r render.Render, dbmap *gorp.DbMap) {
	start, size := GetLimit(req)
	var activities []QActivity
	_, err := dbmap.Select(&activities, "SELECT * FROM t_activity ORDER BY create_time DESC LIMIT ?,? ", start, size)
	CheckErr(err, "GetActivityList")
	logger.Info(activities)

	totalCount, err := dbmap.SelectInt("select count(*) from t_activity")
	m := int(totalCount) % size
	totalPageCount := int(totalCount) / size
	if m != 0 {
		totalPageCount = totalPageCount + 1
	}
	newmap := map[string]interface{}{
		"totalCount":     totalCount,
		"totalPageCount": totalPageCount,
		"activityList":   activities}

	if err == nil {
		r.JSON(200, Resp{0, "获取活动列表成功", newmap})
	} else {
		r.JSON(200, Resp{1007, "获取活动列表失败", nil})
	}
	//	r.HTML(200, "activitylist", newmap)
}

func GetActivity(params martini.Params, r render.Render, dbmap *gorp.DbMap) {
	id := params["id"]
	var activity QActivity
	err := dbmap.SelectOne(&activity, "SELECT * FROM t_activity WHERE aid = ?", id)
	CheckErr(err, "GetActivity")
	logger.Info(activity)
	if err == nil {
		r.JSON(200, Resp{0, "获取活动成功", activity})
	} else {
		r.JSON(200, Resp{1006, "获取活动失败", nil})
	}
}

func DeleteActivity(params martini.Params, r render.Render, dbmap *gorp.DbMap) {
	_, err := dbmap.Exec("DELETE from t_activity WHERE aid=?", params["id"])
	CheckErr(err, "DeleteActivity delete failed")
	if err == nil {
		r.JSON(200, Resp{0, "删除活动成功", nil})
	} else {
		r.JSON(200, Resp{1005, "删除活动失败", nil})
	}
}

func Upload(r *http.Request, render render.Render) {
	logger.Info("Upload")
	err := r.ParseMultipartForm(100000)
	if err != nil {
		render.JSON(500, "server err")
	}
	file, head, err := r.FormFile("file")
	CheckErr(err, "upload Fromfile")
	logger.Info(head.Filename)
	defer file.Close()
	err = Mkdir(config.ImgDir)
	CheckErr(err, "create dir error")
	filepath := config.ImgDir + head.Filename
	fW, err := os.Create(filepath)
	CheckErr(err, "create file error")
	defer fW.Close()
	_, err = io.Copy(fW, file)
	CheckErr(err, "copy file error")
	url, err := upload.UploadToUCloudCND(filepath, "frontCover/"+head.Filename)
	if err == nil {
		render.JSON(200, Resp{0, "图片上传成功！", map[string]interface{}{"url": url}})
	} else {
		render.JSON(200, Resp{1004, "图片上传失败！", nil})
	}
}

func NewActivity(activity NActivity, user sessionauth.User, r render.Render, c *cache.Cache, dbmap *gorp.DbMap) {
	logger.Info("NewActivity ")
	uid := user.UniqueId().(string)
	activity.Aid = AID()
	activity.Uid = uid
	t, err := time.Parse("2006-01-02 15:04", activity.DateString)
	activity.Date = JsonTime{t, true}
	activity.StreamId = GeneraToken8()
	activity.LivePushPath = generatePushPath(activity.StreamId, activity.IsRecord, "")
	activity.LivePullPath = generatePullPath(activity.StreamId)
	logger.Info("info ", activity.String())
	err = dbmap.Insert(&activity)
	CheckErr(err, "NewActivity insert failed")
	if err == nil {
		newmap := map[string]interface{}{"id": activity.Aid, "livePushPath": activity.LivePushPath}
		go PushNewActivity(activity.Title, activity.FrontCover)
		r.JSON(200, Resp{0, "创建活动成功!", newmap})
	} else {
		r.JSON(200, Resp{1002, "创建活动失败", nil})
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
		r.JSON(200, Resp{1003, "更新活动失败", nil})
	} else {
		orgActivity := obj.(*Activity)
		orgActivity.Title = activity.Title
		t, err := time.Parse("2006-01-02 15:04", activity.DateString)
		activity.Date = JsonTime{t, true}
		orgActivity.Date = activity.Date
		orgActivity.Desc = activity.Desc
		orgActivity.ActivityType = activity.ActivityType
		orgActivity.StreamType = activity.StreamType
		orgActivity.FrontCover = activity.FrontCover
		_, err = dbmap.Update(orgActivity)
		CheckErr(err, "UpdateActivity  update failed")
		if err != nil {
			r.JSON(200, Resp{1004, "更新活动失败", nil})
		} else {
			r.JSON(200, Resp{0, "更新活动成功", activity})
		}
	}
}
