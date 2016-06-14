package admin

import (
	. "app/common"
	. "app/models"
	"app/sessionauth"
	"app/sessions"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/coopernurse/gorp"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"github.com/pborman/uuid"
)

var mDbMap *gorp.DbMap

// AdminModel can be any struct that represents a user in my system
type AdminModel struct {
	Id            int64     `form:"id" db:"id"`
	Username      string    `form:"name" db:"username"`
	Password      string    `form:"password" db:"password"`
	Avatar        string    `form:"avatar"  db:"avatar"`
	Updated       time.Time `db:"update_time"`
	Created       time.Time `db:"create_time"`
	authenticated bool      `form:"-" db:"-"`
}

func (admin AdminModel) String() string {
	return fmt.Sprintf("[%s, %s, %d]", admin.Id, admin.Username, admin.Password)
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
	log.Println(id)
	err := mDbMap.SelectOne(u, "SELECT * FROM t_admin WHERE id = ?", id)
	CheckErr(err, "GetById select one")
	if err != nil {
		return err
	}
	return nil
}

func Index(r render.Render) {
	log.Println("Index")
	r.HTML(200, "index", nil)
}

func Login(args martini.Params, req *http.Request, session sessions.Session, r render.Render, dbmap *gorp.DbMap) {
	req.ParseForm()
	if len(req.Form["username"]) > 0 && len(req.Form["password"]) > 0 {
		username := req.Form["username"][0]
		password := req.Form["password"][0]
		log.Println("admin-login:" + username + " " + password)
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
			log.Println(req.URL)
			redirectParams := req.URL.Query()[sessionauth.RedirectParam]
			log.Println(redirectParams)
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

func GetAdminList(r render.Render) {
	r.HTML(200, "adminlist", "")
}

func GetUserList(r render.Render) {
	r.HTML(200, "userlist", "")
}

func GetIncome(r render.Render) {
	r.HTML(200, "income", "")
}

func AddActivity(r render.Render) {
	r.HTML(200, "activityform", "")
}

func GetActivityList(req *http.Request, r render.Render, dbmap *gorp.DbMap) {
	var activities []Activity
	_, err := dbmap.Select(&activities, "select * from t_activity")
	CheckErr(err, "GetActivityList select failed")
	//	if err != nil {
	//		r.JSON(200, Resp{1104, "查询活动失败", nil})
	//	} else {
	//		r.JSON(200, Resp{0, "查询活动成功", activities})
	//	}
	log.Print(activities)
	newmap := map[string]interface{}{"activities": activities}
	r.HTML(200, "activitylist", newmap)
}

func NewActivity(activity Activity, r render.Render, dbmap *gorp.DbMap) {
	log.Println(activity.String())
	activity.VideoId = uuid.New()
	activity.VideoPushPath = "xxxxxx"
	activity.BelongUserId = 1
	//奇葩
	// activity.Date = time.Unix(activity.ADate, 0).Format("2006-01-02 15:04:05")
	activity.Date = time.Unix(activity.ADate, 0)
	activity.Created = time.Now()
	activity.Updated = time.Now()
	err := dbmap.Insert(&activity)
	CheckErr(err, "NewActivity insert failed")
	if err == nil {
		r.JSON(200, "/activity")
	} else {
		r.JSON(500, "删除活动失败")
	}
}

func DeleteActivity(args martini.Params, r render.Render, dbmap *gorp.DbMap) {
	log.Println("DeleteActivity", args["id"])
	_, err := dbmap.Exec("DELETE from t_activity WHERE id=?", args["id"])
	CheckErr(err, "DeleteActivity delete failed")
	if err == nil {
		r.JSON(200, "删除活动成功")
	} else {
		r.JSON(500, "删除活动失败")
	}

}
