package admin

import (
	. "app/common"
	. "app/models"
	"log"
	"net/http"
	"time"

	"app/sessionauth"
	"app/sessions"

	"github.com/coopernurse/gorp"
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
	err := mDbMap.SelectOne(u, "SELECT * FROM t_admin WHERE id = $1", id)
	if err != nil {
		return err
	}
	return nil
}

func Index(r render.Render) {
	r.HTML(200, "index", nil)
}

func Login(req *http.Request, session sessions.Session, r render.Render, dbmap *gorp.DbMap) {
	req.ParseForm()
	if len(req.Form["username"]) > 0 && len(req.Form["password"]) > 0 {
		username := req.Form["username"][0]
		password := req.Form["password"][0]
		var admin AdminModel
		err := dbmap.SelectOne(&admin, "SELECT * FROM t_admin WHERE username = $1 AND password $2", username, password)
		if err != nil {
			r.Redirect("/login")
			return
		} else {
			err := sessionauth.AuthenticateSession(session, &admin)
			if err != nil {
				r.JSON(500, err)
			}
			r.Redirect("/")
			return
		}
	} else {
		r.JSON(200, "缺参数")
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
	r.HTML(200, "activity", "")
}

func NewActivity(activity Activity, r render.Render, dbmap *gorp.DbMap) {
	log.Println("test")
	activity.VideoId = uuid.New()
	activity.VideoPushPath = "xxxxxx"
	activity.BelongUserId = 123
	//奇葩
	// activity.Date = time.Unix(activity.ADate, 0).Format("2006-01-02 15:04:05")
	activity.Date = time.Unix(activity.ADate, 0)
	activity.Created = time.Now()
	err := dbmap.Insert(&activity)
	CheckErr(err, "NewActivity insert failed")
	if err != nil {
		r.HTML(200, "adminlist", "")
	} else {
		r.HTML(200, "adminlist", "")
	}
}
