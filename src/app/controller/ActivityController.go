package controller

import (
	. "app/common"
	logger "app/logger"
	. "app/models"
	"net/http"

	"github.com/coopernurse/gorp"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	cache "github.com/patrickmn/go-cache"
)

const (
	PageSize int = 10
)

func AppointmentActivity(req *http.Request, r render.Render, dbmap *gorp.DbMap) {
	uid := req.Header.Get("uid")
	req.ParseForm()
	aid := req.PostFormValue("aid")
	logger.Info("[ActivityController]", "[AppointmentActivity]", "uid->", uid, "aid->", aid)
	if uid == "" {
		r.JSON(200, Resp{1013, "uid不能为空", nil})
		return
	}
	if aid == "" {
		r.JSON(200, Resp{1105, "添加活动失败,aid不能为空", nil})
		return
	}
	var record Record
	record.Aid = aid
	record.Uid = uid
	record.Type = 0
	err := dbmap.Insert(&record)
	CheckErr("[ActivityController]", "[AppointmentActivity]", "insert failed", err)
	if err != nil {
		r.JSON(200, Resp{1105, "添加活动失败", nil})
	} else {
		r.JSON(200, Resp{0, "预约成功", nil})
	}
}

func PlayActivity(req *http.Request, r render.Render, dbmap *gorp.DbMap) {
	uid := req.Header.Get("uid")
	req.ParseForm()
	aid := req.PostFormValue("aid")
	logger.Info("[ActivityController]", "[PlayActivity]", "uid->", uid, "aid->", aid)
	if aid == "" {
		r.JSON(200, Resp{1105, "添加活动失败,aid不能为空", nil})
		return
	}
	var record Record
	record.Aid = aid
	record.Uid = uid
	record.Type = 1
	err := dbmap.Insert(&record)
	CheckErr("[ActivityController]", "[PlayActivity]", "insert failed", err)
	r.JSON(200, Resp{0, "ok", nil})
}

func GetHomeList(req *http.Request, r render.Render, dbmap *gorp.DbMap) {
	logger.Info("[ActivityController]", "[GetHomeList]")
	var recomendActivities []QActivity
	var activities []QActivity
	_, err := dbmap.Select(&recomendActivities, `SELECT * FROM t_activity t
	WHERE is_recommend = 1 ORDER BY activity_state ASC, create_time DESC`)
	CheckErr("[ActivityController]", "[GetHomeList]", "get recomend list", err)
	_, err = dbmap.Select(&activities, `SELECT * FROM t_activity t 
	WHERE is_recommend = 0 ORDER BY activity_state ASC, create_time DESC LIMIT ?`, PageSize+1)
	CheckErr("[ActivityController]", "[GetHomeList]", "get Activity List", err)
	if err != nil {
		r.JSON(200, Resp{1104, "查询活动失败", nil})
	} else {
		var hasmore bool
		logger.Info("[ActivityController]", "[GetHomeList]", "activity count->", len(activities), " PageSize->", PageSize)
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
	lastAid := params["lastAid"]
	if lastAid == "" {
		r.JSON(200, Resp{1105, "添加活动失败,aid不能为空", nil})
		return
	}
	var activity QActivity
	err := dbmap.SelectOne(&activity, "SELECT * FROM t_activity WHERE aid = ? ", lastAid)
	logger.Info("[ActivityController]", "[GetMoreActivityList]", activity.Created)
	var activities []QActivity
	_, err = dbmap.Select(&activities, `SELECT * FROM t_activity t 
	WHERE t.create_time < ? AND t.activity_state >= ? AND t.is_recommend = 0 ORDER BY t.activity_state ASC, t.create_time DESC  LIMIT ?`, activity.Created, activity.ActivityState, PageSize+1)
	CheckErr("[ActivityController]", "[GetMoreActivityList]", "GetActivityList select failed", err)
	if err != nil {
		r.JSON(200, Resp{1104, "查询活动失败", nil})
	} else {
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
	var activities []QActivity
	_, err := dbmap.Select(&activities, `SELECT * FROM t_activity t WHERE t.stream_type = 0 and t.activity_state = 2 ORDER BY activity_state DESC, t.create_time DESC`)
	CheckErr("[ActivityController]", "[GetLiveActivityList]", "select failed", err)
	if err != nil {
		r.JSON(200, Resp{1104, "查询活动失败", nil})
	} else {
		r.JSON(200, Resp{0, "查询活动成功", activities})
	}
}

func GetActivityDetail(req *http.Request, params martini.Params, r render.Render, dbmap *gorp.DbMap) {
	var activity QActivity
	uid := req.Header.Get("uid")
	logger.Info("[ActivityController]", "[GetActivityDetail]", params)
	aid := params["aid"]
	if aid == "" {
		r.JSON(200, Resp{1105, "添加活动失败,aid不能为空", nil})
		return
	}
	err := dbmap.SelectOne(&activity, `SELECT *, (SELECT  count(*) > 0 FROM t_record r WHERE r.type = 2 AND r.state = 1 AND r.uid =? AND  r.aid=?)  AS pay_state, 
	(SELECT count(*) FROM t_record WHERE type = 0 AND uid = ? AND  aid= ?)  AS appoint_state  
	FROM t_activity t WHERE t.aid = ?`, uid, aid, uid, aid, aid)
	CheckErr("[ActivityController]", "[GetActivityDetail]", "select failed", err)
	if err != nil {
		r.JSON(200, Resp{1103, "活动不存在", nil})
	} else {
		r.JSON(200, Resp{0, "查询活动成功", activity})
	}
}

func JoinGroup(req *http.Request, r render.Render, dbmap *gorp.DbMap) {
	uid := req.Header.Get("uid")
	req.ParseForm()
	aid := req.PostFormValue("aid")
	if uid == "" || aid == "" {
		logger.Info("[ActivityController]", "[JoinGroup]", "uid or aid is ''")
		r.JSON(200, Resp{0, "加入失败！", nil})
		return
	} else {
		_, err := dbmap.Exec(`INSERT INTO t_activity_user_online VALUE(?,?,now())`, aid, uid)
		CheckErr("[ActivityController]", "JoinGroup", "", err)
	}
	r.JSON(200, Resp{0, "加入成功", nil})
}

func LeaveGroup(req *http.Request, r render.Render, dbmap *gorp.DbMap) {
	uid := req.Header.Get("uid")
	req.ParseForm()
	aid := req.PostFormValue("aid")
	if uid == "" || aid == "" {
		logger.Info("[ActivityController]", "[LeaveGroup]", "uid or aid is ''")
		r.JSON(200, Resp{0, "离开失败！", nil})
		return
	} else {
		_, err := dbmap.Exec(`DELETE FROM t_activity_user_online WHERE aid = ? AND uid = ?`, aid, uid)
		CheckErr("[ActivityController]", "[LeaveGroup]", "", err)
	}
	r.JSON(200, Resp{0, "离开成功", nil})
}

func GetLiveActivityMemberCount(req *http.Request, r render.Render, c *cache.Cache, dbmap *gorp.DbMap) {
	uid := req.Header.Get("uid")
	req.ParseForm()
	aid := req.PostFormValue("aid")
	if uid == "" || aid == "" {
		logger.Info("[ActivityController]", "[GetLiveActivityMemberCount]", "uid or aid is ''")
		r.JSON(200, Resp{0, "缺少参数uid和aid不能为空", nil})
	} else {
		count, err := dbmap.SelectInt("SELECT COUNT(*) FROM t_activity_user_online WHERE aid = ?", aid)
		CheckErr("[ActivityController]", "[GetLiveActivityMemberCount]", "", err)
		if err == nil {
			r.JSON(200, Resp{0, "获取在线成员信息成功", map[string]int{"count": int(count)}})
		} else {
			r.JSON(200, Resp{1402, "获取在线成员信息失败", nil})
		}
	}
}

func GetLiveActivityMemberList(req *http.Request, r render.Render, c *cache.Cache, dbmap *gorp.DbMap) {
	uid := req.Header.Get("uid")
	//	query := req.URL.Query()
	//	var aid string
	//	if len(query["aid"]) > 0 {
	//		aid = query["aid"][0]
	//	}
	req.ParseForm()
	aid := req.PostFormValue("aid")
	if uid == "" || aid == "" {
		logger.Info("[ActivityController]", "[GetLiveActivityMemberList]", "uid or aid is ''")
		r.JSON(200, Resp{0, "缺少参数uid和aid不能为空", nil})
		return
	} else {
		var users []OnlineUser
		_, err := dbmap.Select(&users, `SELECT o.uid,u.avatar,u.nickname FROM t_activity_user_online o LEFT JOIN t_user u ON o.uid = u.uid WHERE o.aid = ?`, aid)
		CheckErr("[ActivityController]", "[GetLiveActivityMemberList]", "", err)
		if err == nil {
			r.JSON(200, Resp{0, "获取在线成员信息成功", users})
		} else {
			r.JSON(200, Resp{1402, "获取在线成员信息失败", nil})
		}
	}
}

func GetSharePage(params martini.Params, r render.Render, c *cache.Cache, dbmap *gorp.DbMap) {
	platform := params["platform"]
	logger.Info("[ActivityController]", "[GetSharePage]", "platform", platform)
	var activity QActivity
	err := dbmap.SelectOne(&activity, "select * from t_activity where aid =?", params["id"])
	CheckErr("[ActivityController]", "[GetSharePage]", "", err)
	//TODO apple 下载地址 https://itunes.apple.com/cn/app/qq/id444934666
	url := "http://a.app.qq.com/o/simple.jsp?pkgname=com.kouchen.mininetlive"
	if err == nil {
		var result struct {
			DownloadUrl string      `json:"downloadUrl"`
			Activity    interface{} `json:"activity"`
		}
		result.DownloadUrl = url
		result.Activity = activity
		CheckErr("[ActivityController]", "[GetSharePage]", "", err)
		r.JSON(200, Resp{0, "获取成功", result})
	} else {
		r.JSON(200, Resp{1103, "获取在线成员信息失败", nil})
	}
}
