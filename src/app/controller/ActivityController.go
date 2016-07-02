package controller

import (
	. "app/common"
	easemob "app/easemob"
	logger "app/logger"
	. "app/models"
	"net/http"
	"strings"

	"github.com/coopernurse/gorp"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	cache "github.com/patrickmn/go-cache"
)

const (
	PageSize int = 9
)

func AppointmentActivity(req *http.Request, r render.Render, dbmap *gorp.DbMap) {
	uid := req.Header.Get("uid")
	req.ParseForm()
	aid := req.PostFormValue("aid")
	if aid == "" {
		r.JSON(200, Resp{1105, "添加活动失败,aid不能为空", nil})
	}
	var record AppointmentRecord
	record.Aid = aid
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
	var recomendActivities []QActivity
	var activities []QActivity
	_, err := dbmap.Select(&recomendActivities, "SELECT * FROM t_activity WHERE is_recommend = 1 ORDER BY activity_state ASC, create_time DESC")
	CheckErr(err, "get recomend list")
	if err == nil {
		for _, activity := range recomendActivities {
			queryState(activity, uid, dbmap)
		}
	}
	_, err = dbmap.Select(&activities, "SELECT * FROM t_activity WHERE is_recommend = 0 ORDER BY activity_state ASC, create_time DESC LIMIT ? ", PageSize+1)
	CheckErr(err, "get Activity List")
	if err == nil {
		for _, activity := range activities {
			queryState(activity, uid, dbmap)
		}
	}
	if err != nil {
		r.JSON(200, Resp{1104, "查询活动失败", nil})
	} else {
		var hasmore bool
		logger.Info(len(activities))
		if len(activities) > PageSize {
			hasmore = true
			activities = activities[:PageSize]
		} else {
			hasmore = false
		}
		r.JSON(200, Resp{0, "查询活动成功", map[string]interface{}{
			"hasmore": hasmore, "recommend": recomendActivities, "general": activities}})
	}
}

func GetMoreActivityList(req *http.Request, params martini.Params, r render.Render, dbmap *gorp.DbMap) {
	uid := req.Header.Get("uid")
	lastAid := params["lastAid"]
	var activity QActivity
	err := dbmap.SelectOne(&activity, "SELECT * FROM t_activity WHERE aid = ? ", lastAid)
	logger.Info("GetMoreActivityList..", activity.Created)
	var activities []QActivity
	_, err = dbmap.Select(&activities, "SELECT * FROM t_activity WHERE create_time < ? AND is_recommend = 0 ORDER BY activity_state ASC, create_time DESC  LIMIT ?", activity.Created, PageSize+1)
	CheckErr(err, "GetActivityList select failed")
	if err != nil {
		r.JSON(200, Resp{1104, "查询活动失败", nil})
	} else {
		for _, activity := range activities {
			queryState(activity, uid, dbmap)
		}
		var hasmore bool
		if len(activities) > PageSize {
			hasmore = true
			activities = activities[:PageSize]
		} else {
			hasmore = false
		}
		r.JSON(200, Resp{0, "查询活动成功", map[string]interface{}{"hasmore": hasmore, "general": activities}})
	}
}

func GetLiveActivityList(req *http.Request, r render.Render, dbmap *gorp.DbMap) {
	uid := req.Header.Get("uid")
	var activities []QActivity
	_, err := dbmap.Select(&activities, "SELECT * FROM t_activity WHERE activity_state = 1 AND stream_type = 0 ORDER BY create_time DESC")
	CheckErr(err, "GetLiveActivityList select failed")
	if err != nil {
		r.JSON(200, Resp{1104, "查询活动失败", nil})
	} else {
		for _, activity := range activities {
			queryState(activity, uid, dbmap)
		}
		r.JSON(200, Resp{0, "查询活动成功", activities})
	}
}

func queryState(activity QActivity, uid string, dbmap *gorp.DbMap) {
	if uid != "" {
		count, err := dbmap.SelectInt("select count(*) from t_pay_record where aid = ? and uid = ?", activity.Aid, uid)
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
	var activity QActivity
	err := dbmap.SelectOne(&activity, "select * from t_activity where aid =?", args["id"])
	CheckErr(err, "GetActivity select failed")
	if err != nil {
		r.JSON(200, Resp{1103, "活动不存在", nil})
	} else {
		r.JSON(200, Resp{0, "查询活动成功", activity})
	}
}

func JoinGroup(req *http.Request, r render.Render, c *cache.Cache) {
	uid := req.Header.Get("uid")
	req.ParseForm()
	groupId := req.PostFormValue("groupId")
	if uid == "" || groupId == "" {
		logger.Info("JoinGroup", "uid or groupId is ''")
	} else {
		easemob.JoinGroup(groupId, uid, c)
	}
	r.JSON(200, Resp{0, "成功", nil})
}

func LeaveGroup(req *http.Request, r render.Render, c *cache.Cache) {
	uid := req.Header.Get("uid")
	req.ParseForm()
	groupId := req.PostFormValue("groupId")
	if uid == "" || groupId == "" {
		logger.Info("LeaveGroup", "uid or groupId is ''")
	} else {
		easemob.LeaveGroup(groupId, uid, c)
	}
	r.JSON(200, Resp{0, "成功", nil})
}

func GetLiveActivityMemberCount(req *http.Request, r render.Render, c *cache.Cache, dbmap *gorp.DbMap) {
	uid := req.Header.Get("uid")
	req.ParseForm()
	groupId := req.PostFormValue("groupId")
	if uid == "" || groupId == "" {
		logger.Info("LeaveGroup", "uid or groupId is ''")
		r.JSON(200, Resp{0, "缺少参数", nil})
	} else {
		count, err := easemob.GetGroupMemberCount(groupId, c)

		if err == nil {
			r.JSON(200, Resp{0, "获取在线成员信息成功", map[string]int{"count": count}})
		} else {
			r.JSON(200, Resp{1402, "获取在线成员信息失败", nil})
		}
	}
}

func GetLiveActivityMemberList(req *http.Request, r render.Render, c *cache.Cache, dbmap *gorp.DbMap) {
	uid := req.Header.Get("uid")
	req.ParseForm()
	groupId := req.PostFormValue("groupId")
	if uid == "" || groupId == "" {
		logger.Info("LeaveGroup", "uid or groupId is ''")
		r.JSON(200, Resp{0, "缺少参数", nil})
	} else {
		uids, err := easemob.GetGroupMemberList(groupId, c)
		if err == nil {
			var users []User
			_, err = dbmap.Select(&users, "SELECT * FROM t_user WHERE uid IN ('"+strings.Join(uids, "','")+"')")
			if err == nil {
				r.JSON(200, Resp{0, "获取在线成员信息成功", users})
				return
			}
		}
		r.JSON(200, Resp{1402, "获取在线成员信息失败", nil})
	}
}
