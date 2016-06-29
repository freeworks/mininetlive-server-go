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
	"log"
	"github.com/coopernurse/gorp"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	cache "github.com/patrickmn/go-cache"
)

var mDbMap *gorp.DbMap

func (admin AdminModel) String() string {
	adminString := fmt.Sprintf("[%s, %s, %d]", admin.Id, admin.NickName, admin.Password)
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
	u.Authenticated = true
}

// Logout will preform any actions that are required to completely
// logout a user.
func (u *AdminModel) Logout() {
	// Remove from logged-in user's list
	// etc ...
	u.Authenticated = false
}

func (u *AdminModel) IsAuthenticated() bool {
	return u.Authenticated
}

func (u *AdminModel) UniqueId() interface{} {
	return u.Uid
}

// GetById will populate a user object from a database model with
// a matching id.
func (u *AdminModel) GetById(uid interface{}) error {
	logger.Info("GetById:", uid)
	err := mDbMap.SelectOne(u, "SELECT * FROM t_admin WHERE uid = ?", uid)
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
	phone := req.PostFormValue("phone")
	password := req.PostFormValue("password")
	if validatePhone(phone) && validatePassword(password) {
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
	}else{
		r.JSON(500, "账号或密码错误！")
	}
}

func validatePhone(phone string) bool {
	//TODO
	return true
}

func validatePassword(password string) bool {
	//TODO
	return true
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

func NewActivity(activity NActivity, user sessionauth.User, r render.Render, c *cache.Cache, dbmap *gorp.DbMap) {
	dbmap.TraceOn("[gorp]", log.New(os.Stdout, "NewActivity:", log.Lmicroseconds))
	logger.Info("NewActivity ")
	// uid := user.UniqueId().(string)
	uid :="1e046709049d59b5"
	groupId,err := easemob.CreateGroup(uid, activity.Title, c)
	if err != nil {
		CheckErr(err, "easemob create group error")
		r.JSON(500, "创建活动失败")
		return
	}
	activity.Aid = AID()
	activity.Uid = uid
	activity.Date = time.Unix(activity.ADate, 0)
	activity.Created = time.Now()
	activity.GroupId = groupId
	activity.StreamId = GeneraToken8()
	activity.LivePushPath = generatePushPath(activity.StreamId,activity.IsRecord,"")
	logger.Info("info ",activity.String())
	err = dbmap.Insert(&activity)
	CheckErr(err, "NewActivity insert failed")
	if err == nil {
		r.JSON(200, "/activity")
	} else {
		r.JSON(500, "创建活动失败")
	}
	dbmap.TraceOff()
}

func generatePushPath(streamId string,record bool,filename string) string {
	pushPath := "rtmp://域名/接入点/"+streamId+"?record="+strconv.FormatBool(record)
	if filename != "" {
		pushPath = pushPath+ "&filename="+filename 
	}
	logger.Info("GeneratePushPath :",pushPath)

	return  pushPath
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
