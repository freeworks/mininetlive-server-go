package controller

import (
	. "app/common"
	logger "app/logger"
	. "app/models"
	"bytes"
	"net/http"
	"strconv"
	"strings"

	"github.com/coopernurse/gorp"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
)

func AppointmentActivity(req *http.Request, r render.Render, dbmap *gorp.DbMap) {
	uid := req.Header.Get("uid")
	req.ParseForm()
	var record AppointmentRecord
	record.Aid = req.PostFormValue("aid")
	record.Uid = uid
	err := dbmap.Insert(&record)
	CheckErr(err, "AppointmentActivity insert failed")
	if err != nil {
		r.JSON(200, Resp{1105, "添加活动失败", nil})
	} else {
		r.JSON(200, Resp{0, "预约成功", nil})
	}
}

func PlayActivity(req *http.Request, r render.Render, dbmap *gorp.DbMap) {
	uid := req.Header.Get("uid")
	req.ParseForm()
	var record PlayRecord
	record.Aid = req.PostFormValue("aid")
	record.Uid = uid
	err := dbmap.Insert(&record)
	CheckErr(err, "PayActivity insert failed")
	r.JSON(200, Resp{0, "ok", nil})
}

func GetHomeList(req *http.Request, r render.Render, dbmap *gorp.DbMap) {
	uid := req.Header.Get("uid")
	var recomendIds []string
	_, err := dbmap.Select(&recomendIds, "SELECT aid FROM t_recomend ORDER BY update_time DESC")
	CheckErr(err, "get aid from recomend")

	var recomendActivities []Activity
	var activities []Activity
	logger.Info(recomendIds)
	if len(recomendIds) == 0 {
		_, err = dbmap.Select(&recomendActivities, "SELECT * FROM t_activity WHERE aid ORDER BY update_time DESC")
		CheckErr(err, "get recomend list")
		_, err = dbmap.Select(&activities, "SELECT * FROM t_activity ORDER BY update_time DESC")
		CheckErr(err, "get Activity List")
	} else {
		var buffer bytes.Buffer
		values := make(map[string]string)
		for index, recomendId := range recomendIds {
			key := "id" + strconv.Itoa(index)
			buffer.WriteString(":" + key + ",")
			values[key] = recomendId
		}
		condition := strings.TrimRight(buffer.String(), ",")
		sql := "SELECT * FROM t_activity WHERE aid IN (" + condition + ") ORDER BY update_time DESC"
		logger.Info(sql)
		logger.Info(values)
		_, err = dbmap.Select(&recomendActivities, sql, values)
		CheckErr(err, "get recomend list")
		sql = "SELECT * FROM t_activity WHERE aid NOT IN (" + condition + ") ORDER BY update_time DESC"
		_, err = dbmap.Select(&activities, sql, values)
		CheckErr(err, "get Activity List")
	}
	if err != nil {
		r.JSON(200, Resp{1104, "查询活动失败", nil})
	} else {
		for _, activity := range activities {
			queryCount(activity, uid, dbmap)
		}
		r.JSON(200, Resp{0, "查询活动成功", map[string]interface{}{"recommend": recomendActivities, "general": activities}})
	}
}

func GetMoreActivityList(req *http.Request, params martini.Params, r render.Render, dbmap *gorp.DbMap) {
	uid := req.Header.Get("uid")
	lastAid := params["lastAid"]
	lastId, err := dbmap.SelectInt("SELECT id FROM t_activity WHERE aid = ? ", lastAid)
	var activities []Activity
	_, err = dbmap.Select(&activities, "SELECT * FROM t_activity WHERE id > ? ORDER BY create_time DESC LIMIT 10", lastId)
	CheckErr(err, "GetActivityList select failed")
	if err != nil {
		r.JSON(200, Resp{1104, "查询活动失败", nil})
	} else {
		for _, activity := range activities {
			queryCount(activity, uid, dbmap)
		}
		r.JSON(200, Resp{0, "查询活动成功", activities})
	}
}

func GetLiveActivityList(req *http.Request, r render.Render, dbmap *gorp.DbMap) {
	uid := req.Header.Get("uid")
	var activities []Activity
	_, err := dbmap.Select(&activities, "SELECT * FROM t_activity WHERE activity_state = 1 AND video_type = 0 ORDER BY create_time DESC LIMIT 10")
	CheckErr(err, "GetLiveActivityList select failed")
	if err != nil {
		r.JSON(200, Resp{1104, "查询活动失败", nil})
	} else {
		for _, activity := range activities {
			queryCount(activity, uid, dbmap)
		}
		r.JSON(200, Resp{0, "查询活动成功", activities})
	}
}

func queryCount(activity Activity, uid string, dbmap *gorp.DbMap) {
	count, err := dbmap.SelectInt("select count(*) from t_play_record where aid = ? ", activity.Aid)
	CheckErr(err, "get play count")
	activity.PlayCount = int(count)
	count, err = dbmap.SelectInt("select count(*) from t_appointment_record where aid = ? ", activity.Aid)
	CheckErr(err, "get appointment count")
	activity.PlayCount = int(count)
	if uid != "" {
		count, err = dbmap.SelectInt("select count(*) from t_pay_record where aid = ? and uid = ?", activity.Aid, uid)
		CheckErr(err, "get appointment count")
		if count == 0 {
			activity.PayState = 0
		} else {
			activity.PayState = 1
		}

		count, err = dbmap.SelectInt("select count(*) from t_appointment_record where aid = ? and uid = ?", activity.Aid, uid)
		CheckErr(err, "get appointment count")
		if count == 0 {
			activity.AppointState = 0
		} else {
			activity.AppointState = 1
		}
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
