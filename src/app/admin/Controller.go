package admin

import (
	. "app/common"
	. "app/models"
	"log"
	"net/http"
	"time"

	"github.com/coopernurse/gorp"
	"github.com/martini-contrib/render"
	"github.com/pborman/uuid"
)

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

func Index(r render.Render) {
	r.HTML(200, "index", "")
}

func Login(r render.Render) {
	r.HTML(200, "index", "")
}
func Logout(r render.Render) {
	r.HTML(200, "signin", "")
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
