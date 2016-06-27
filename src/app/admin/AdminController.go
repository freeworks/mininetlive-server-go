package admin

import (
	. "app/common"
	config "app/config"
	easemob "app/easemob"
	logger "app/logger"
	. "app/models"
	"app/sessionauth"
	"app/sessions"
	upload "app/upload"
	"fmt"
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

var mDbMap *gorp.DbMap

func (admin AdminModel) String() string {
	adminString := fmt.Sprintf("[%s, %s, %d]", admin.Id, admin.Username, admin.Password)
	logger.Info(adminString)
	return adminString
}

// GetAnonymousUser should generate an anonymous user model
// for all sessions. This should be an unauthenticated 0 value struct.
func GenerateAnonymousUser() sessionauth.User {
	return &AdminModel{}
}

func SetDBMap(dbmap *gorp.DbMap) {
	mDbMap = dbmap
}

// Login will preform any actions that are required to make a user model
// officially authenticated.
func (u *AdminModel) Login() {
	// Update last login time
	// Add to logged-in user's list
	// etc ...
	logger.Info("login ....")
	u.authenticated = true
}

// Logout will preform any actions that are required to completely
// logout a user.
func (u *AdminModel) Logout() {
	// Remove from logged-in user's list
	// etc ...
	u.authenticated = false
}

func (u *AdminModel) IsAuthenticated() bool {
	return u.authenticated
}

func (u *AdminModel) UniqueId() interface{} {
	return u.Id
}

// GetById will populate a user object from a database model with
// a matching id.
func (u *AdminModel) GetById(id interface{}) error {
	logger.Info("GetById:", id)
	err := mDbMap.SelectOne(u, "SELECT * FROM t_admin WHERE id = ?", id)
	CheckErr(err, "GetById select one")
	if err != nil {
		return err
	}
	return nil
}

func Index(r render.Render) {
	logger.Debug("Index")
	r.HTML(200, "index", nil)
}

func Login(args martini.Params, req *http.Request, session sessions.Session, r render.Render, dbmap *gorp.DbMap) {
	req.ParseForm()
	if len(req.Form["username"]) > 0 && len(req.Form["password"]) > 0 {
		username := req.Form["username"][0]
		password := req.Form["password"][0]
		logger.Info("admin-login:" + username + " " + password)
		var admin AdminModel
		err := dbmap.SelectOne(&admin, "SELECT * FROM t_admin WHERE username = ? AND password = ?", username, password)
		CheckErr(err, "Login select one")
		if err != nil {
			//			r.Redirect("/login")
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
		r.JSON(500, "缺参数")
	}
}

func RedirectLogin(r render.Render) {
	r.HTML(200, "login", nil)
}

func Logout(session sessions.Session, user sessionauth.User, r render.Render) {
	sessionauth.Logout(session, user)
	r.Redirect("/")
}

func GetAdminList(r render.Render, dbmap *gorp.DbMap) {
	var admins []AdminModel
	_, err := dbmap.Select(&admins, "select * from t_admin")
	CheckErr(err, "GetAdminList select failed")
	if err == nil {
		newmap := map[string]interface{}{"admins": admins}
		r.HTML(200, "adminlist", newmap)
	} else {
		r.HTML(500, "adminlist", nil)
	}
}

//type UserMode struct {
//	User
//	accountType int
//	Phone       string
//	Plat
//}
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

func GetIncome(r render.Render) {
	r.HTML(200, "income", "")
}

func AddActivity(r render.Render) {
	r.HTML(200, "activityform", "")
}

func GetActivityList(r render.Render, dbmap *gorp.DbMap) {
	var activities []Activity
	_, err := dbmap.Select(&activities, "select * from t_activity")
	CheckErr(err, "GetActivityList select failed")
	logger.Info(activities)
	newmap := map[string]interface{}{"activities": activities}
	r.HTML(200, "activitylist", newmap)
}

func NewActivity(activity Activity, user sessionauth.User, r render.Render, c *cache.Cache, dbmap *gorp.DbMap) {
	// uid := UUID()
	err := easemob.CreateGroup(strconv.FormatInt(user.UniqueId().(int64), 10), activity.Title, activity.Title, c)
	if err != nil {
		CheckErr(err, "easemob create group error")
		r.JSON(500, "创建活动失败")
		return
	}
	logger.Info(activity.String())
	activity.VideoId = GeneraToken8()
	activity.VideoPushPath = fmt.Sprintf(config.RtmpPath, activity.VideoId)
	// activity.Uid = uid
	//奇葩
	// activity.Date = time.Unix(activity.ADate, 0).Format("2006-01-02 15:04:05")
	activity.Date = time.Unix(activity.ADate, 0)
	activity.Created = time.Now()
	CheckErr(err, "create group")
	err = dbmap.Insert(&activity)
	CheckErr(err, "NewActivity insert failed")
	if err == nil {
		r.JSON(200, "/activity")
	} else {
		r.JSON(500, "创建活动失败")
	}
}

func UpdateActivity(activity Activity, r render.Render, dbmap *gorp.DbMap) {
	logger.Info(activity.String())
	activity.VideoId = UUID()
	activity.VideoPushPath = "xxxxxx"
	//奇葩
	// activity.Date = time.Unix(activity.ADate, 0).Format("2006-01-02 15:04:05")
	activity.Date = time.Unix(activity.ADate, 0)
	activity.Updated = time.Now()
	_, err := dbmap.Update(&activity)
	CheckErr(err, "NewActivity insert failed")
	if err == nil {
		r.JSON(200, "/activity")
	} else {
		r.JSON(500, "删除活动失败")
	}
}

func DeleteActivity(args martini.Params, r render.Render, dbmap *gorp.DbMap) {
	logger.Info("DeleteActivity", args["id"])
	_, err := dbmap.Exec("DELETE from t_activity WHERE id=?", args["id"])
	CheckErr(err, "DeleteActivity delete failed")
	if err == nil {
		r.JSON(200, "删除活动成功")
	} else {
		r.JSON(500, "删除活动失败")
	}

}

func Upload(r *http.Request, render render.Render) {
	logger.Info("parsing form")
	err := r.ParseMultipartForm(100000)
	// CheckErr("upload ParseMultipartForm",err)
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
	url, err := upload.UploadToUCloudCND(filepath, "frontCover/"+head.Filename, render)
	if err == nil {
		render.JSON(200, map[string]interface{}{"status": strconv.Itoa(1), "id": strconv.Itoa(5), "url": url})
	} else {
		render.JSON(200, map[string]interface{}{"status": strconv.Itoa(0)})
	}
}

//func NewActivity(activity Activity, r render.Render, dbmap *gorp.DbMap) {
//	activity.VideoId = uuid.New()
//	activity.VideoPushPath = "xxxxxx"
//	activity.BelongUserId = 123
//	//奇葩
//	// activity.Date = time.Unix(activity.ADate, 0).Format("2006-01-02 15:04:05")
//	activity.Date = time.Unix(activity.ADate, 0)
//	activity.Created = time.Now()
//	err := dbmap.Insert(&activity)
//	CheckErr(err, "NewActivity insert failed")
//	if err != nil {
//		r.JSON(200, Resp{1100, "添加活动失败", nil})
//	} else {
//		r.JSON(200, Resp{0, "添加活动成功", activity})
//	}
//}

//func UpdateActivity(args martini.Params, activity Activity, r render.Render, dbmap *gorp.DbMap) {
//	obj, err := dbmap.Get(Activity{}, args["id"])
//	CheckErr(err, "UpdateActivity get Activity err ")
//	if err != nil {
//		r.JSON(200, Resp{1101, "更新活动失败", nil})
//	} else {
//		orgActivity := obj.(*Activity)
//		orgActivity.Title = activity.Title
//		orgActivity.Date = activity.Date
//		orgActivity.Desc = activity.Desc
//		orgActivity.Type = activity.Type
//		orgActivity.VideoType = activity.VideoType
//		orgActivity.FontCover = activity.FontCover
//		orgActivity.Updated = time.Now()
//		_, err := dbmap.Update(orgActivity)
//		CheckErr(err, "UpdateActivity  update failed")
//		if err != nil {
//			r.JSON(200, Resp{1101, "更新活动失败", nil})
//		} else {
//			r.JSON(200, Resp{0, "更新活动成功", activity})
//		}
//	}
//}

//func DeleteActivity(args martini.Params, r render.Render, dbmap *gorp.DbMap) {
//	_, err := dbmap.Exec("DELETE from t_activity WHERE id=?", args["id"])
//	CheckErr(err, "DeleteActivity delete failed")
//	if err != nil {
//		r.JSON(200, Resp{1102, "删除活动失败", nil})
//	} else {
//		r.JSON(200, Resp{0, "删除活动成功", nil})
//	}
//}

//func DeleteUser(args martini.Params, r render.Render, dbmap *gorp.DbMap) {
//	_, err := dbmap.Exec("DELETE from t_user WHERE id=?", args["id"])
//	CheckErr(err, "DeleteUser delete failed")
//	if err != nil {
//		r.JSON(200, Resp{1007, "删除用户失败", nil})
//	} else {
//		r.JSON(200, Resp{0, "删除用户成功", nil})
//	}
//}
