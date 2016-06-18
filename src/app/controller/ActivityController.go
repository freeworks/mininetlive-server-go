package controller

import (
	. "app/common"
	. "app/models"
	"net/http"
	"time"

	"github.com/coopernurse/gorp"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
)

func AppointmentActivity(req *http.Request, r render.Render, dbmap *gorp.DbMap) {
	req.ParseForm()
	var record AppointmentRecord
	record.Aid = req.PostFormValue("aid")
	record.Uid = req.PostFormValue("uid")
	record.Created = time.Now()
	err := dbmap.Insert(&record)
	CheckErr(err, "AppointmentActivity insert failed")
	if err != nil {
		r.JSON(200, Resp{1105, "添加活动失败", nil})
	} else {
		r.JSON(200, Resp{0, "预约成功", nil})
	}
}

func PlayActivity(req *http.Request, r render.Render, dbmap *gorp.DbMap) {
	req.ParseForm()
	var record PlayRecord
	record.Aid = req.PostFormValue("aid")
	record.Uid = req.PostFormValue("uid")
	record.Created = time.Now()
	err := dbmap.Insert(&record)
	CheckErr(err, "PayActivity insert failed")
	r.JSON(200, Resp{0, "ok", nil})
}

//TODO 前10个
func GetActivityList(req *http.Request, r render.Render, dbmap *gorp.DbMap) {
	var activities []Activity
	_, err := dbmap.Select(&activities, "select * from t_activity")
	CheckErr(err, "GetActivityList select failed")
	if err != nil {
		r.JSON(200, Resp{1104, "查询活动失败", nil})
	} else {
		r.JSON(200, Resp{0, "查询活动成功", activities})
	}
}

//TODO 前10个
func GetMoreActivityList(req *http.Request, r render.Render, dbmap *gorp.DbMap) {
	var activities []Activity
	_, err := dbmap.Select(&activities, "select * from t_activity")
	CheckErr(err, "GetActivityList select failed")
	if err != nil {
		r.JSON(200, Resp{1104, "查询活动失败", nil})
	} else {
		r.JSON(200, Resp{0, "查询活动成功", activities})
	}
}

func GetActivityDetail(args martini.Params, r render.Render, dbmap *gorp.DbMap) {
	var activity Activity
	err := dbmap.SelectOne(&activity, "select * from t_activity where aid =?", args["id"])
	CheckErr(err, "GetActivity select failed")
	if err != nil {
		r.JSON(200, Resp{1103, "活动不存在", nil})
	} else {
		r.JSON(200, Resp{0, "查询活动成功", activity})
	}
}
