package controller

import (
	. "app/common"
	. "app/models"
	"net/http"
	"strconv"
	"time"

	"github.com/coopernurse/gorp"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"github.com/pborman/uuid"
)

//0 成功
//1100 创建活动失败
//1101 更新活动失败
//1102 删除活动失败
//1103 获取活动失败
//1104 获取活动列表失败
//1105 预约活动失败
//1106 取消预约活动失败

func AppointmentActivity(args martini.Params, r render.Render, dbmap *gorp.DbMap) {
	var record AppointmentRecord
	record.ActivityId, _ = strconv.Atoi(args["id"])
	record.UserId = 1 //TODO session 取id
	record.Created = time.Now()
	err := dbmap.Insert(&record)
	CheckErr(err, "AppointmentActivity insert failed")
	if err != nil {
		r.JSON(200, Resp{1105, "添加活动失败", nil})
	} else {
		r.JSON(200, Resp{0, "预约成功", nil})
	}
}

func CancelAppointmentActivity(args martini.Params, r render.Render, dbmap *gorp.DbMap) {
	var record AppointmentRecord
	userId := 1 //session from Id
	err := dbmap.SelectOne(&record, "SELECT * FROM t_appointment_record WHERE activity_id = ? AND user_id = ?",
		args["activityId"], userId)
	CheckErr(err, "CancelAppointmentActivity selectOne failed")
	if err != nil {
		r.JSON(200, Resp{1106, "取消预约活动失败", nil})
	} else {
		record.State = 2
		_, err := dbmap.Update(record)
		CheckErr(err, "CancelAppointmentActivity update failed")
		if err != nil {
			r.JSON(200, Resp{1106, "取消预约活动失败", nil})
		} else {
			r.JSON(200, Resp{0, "更新活动成功", nil})
		}
	}
}

func PayActivity(req *http.Request, args martini.Params, r render.Render, dbmap *gorp.DbMap) {
	req.ParseForm()
	var record PayRecord
	record.ActivityId, _ = strconv.Atoi(args["id"])
	record.UserId = 1 //TODO session 取id
	record.Amount, _ = strconv.Atoi(req.Form["amount"][0])
	record.Type, _ = strconv.Atoi(req.Form["type"][0])
	record.Created = time.Now()
	//TODO校验
	err := dbmap.Insert(&record)
	CheckErr(err, "PayActivity insert failed")
	if err != nil {
		r.JSON(200, Resp{1105, "支付失败", nil})
	} else {
		r.JSON(200, Resp{0, "支付成功", nil})
	}
}

func PlayActivity(req *http.Request, args martini.Params, activity Activity, r render.Render, dbmap *gorp.DbMap) {
	req.ParseForm()
	var record PlayRecord
	record.ActivityId, _ = strconv.Atoi(args["id"])
	record.UserId = 1 //TODO session 取id
	record.Type, _ = strconv.Atoi(req.Form["type"][0])
	record.Created = time.Now()
	err := dbmap.Insert(&record)
	CheckErr(err, "PayActivity insert failed")
	if err != nil {
		r.JSON(200, Resp{1105, "支付失败", nil})
	} else {
		r.JSON(200, Resp{0, "支付成功", nil})
	}
}

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

func GetActivity(args martini.Params, r render.Render, dbmap *gorp.DbMap) {
	var activity Activity
	err := dbmap.SelectOne(&activity, "select * from t_activity where id =?", args["id"])
	CheckErr(err, "GetActivity select failed")
	if err != nil {
		r.JSON(200, Resp{1103, "活动不存在", nil})
	} else {
		r.JSON(200, Resp{0, "查询活动成功", activity})
	}
}

func NewActivity(activity Activity, r render.Render, dbmap *gorp.DbMap) {
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
		r.JSON(200, Resp{1100, "添加活动失败", nil})
	} else {
		r.JSON(200, Resp{0, "添加活动成功", activity})
	}
}

func UpdateActivity(args martini.Params, activity Activity, r render.Render, dbmap *gorp.DbMap) {
	obj, err := dbmap.Get(Activity{}, args["id"])
	CheckErr(err, "UpdateActivity get Activity err ")
	if err != nil {
		r.JSON(200, Resp{1101, "更新活动失败", nil})
	} else {
		orgActivity := obj.(*Activity)
		orgActivity.Title = activity.Title
		orgActivity.Date = activity.Date
		orgActivity.Desc = activity.Desc
		orgActivity.Type = activity.Type
		orgActivity.VideoType = activity.VideoType
		orgActivity.FontCover = activity.FontCover
		orgActivity.Updated = time.Now()
		_, err := dbmap.Update(orgActivity)
		CheckErr(err, "UpdateActivity  update failed")
		if err != nil {
			r.JSON(200, Resp{1101, "更新活动失败", nil})
		} else {
			r.JSON(200, Resp{0, "更新活动成功", activity})
		}
	}
}

func DeleteActivity(args martini.Params, r render.Render, dbmap *gorp.DbMap) {
	_, err := dbmap.Exec("DELETE from t_activity WHERE id=?", args["id"])
	CheckErr(err, "DeleteActivity delete failed")
	if err != nil {
		r.JSON(200, Resp{1102, "删除活动失败", nil})
	} else {
		r.JSON(200, Resp{0, "删除活动成功", nil})
	}
}
