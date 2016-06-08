package controller

import (
		."app/common"
		."app/models"
		"github.com/martini-contrib/render"
		"github.com/coopernurse/gorp"
		"github.com/go-martini/martini"
		"github.com/pborman/uuid"
		"time"
	)
		
func GetActivityList(r render.Render, dbmap *gorp.DbMap) {
	var activities []Activity
	_, err := dbmap.Select(&activities, "select * from t_activity")
	CheckErr(err, "GetActivityList select failed")
	if err != nil {
		r.JSON(200, Resp{2002, "查询活动失败", nil})
	} else {
		r.JSON(200, Resp{0, "查询活动成功", activities})
	}
}

func GetActivity(args martini.Params, r render.Render, dbmap *gorp.DbMap) {
	var activity Activity
	err := dbmap.SelectOne(&activity, "select * from t_activity where id =?", args["id"])
	CheckErr(err, "GetActivity select failed")
	if err != nil {
		r.JSON(200, Resp{2003, "活动不存在", nil})
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
	err := dbmap.Insert(&activity)
	CheckErr(err, "NewActivity insert failed")
	if err != nil {
		r.JSON(200, Resp{2001, "添加活动失败", nil})
	} else {
		r.JSON(200, Resp{0, "添加活动成功", activity})
	}
}

func UpdateActivity(args martini.Params, activity Activity, r render.Render, dbmap *gorp.DbMap) {
	obj, err := dbmap.Get(Activity{}, args["id"])
	CheckErr(err,"UpdateActivity get Activity err ")
	if err != nil {
		r.JSON(200, Resp{2004, "更新活动失败", nil})
	} else {
		orgActivity := obj.(*Activity)
		orgActivity.Title = activity.Title
		orgActivity.Date = activity.Date
		orgActivity.Desc = activity.Desc
		orgActivity.Type = activity.Type
		orgActivity.VideoType = activity.VideoType
		orgActivity.FontCover = activity.FontCover
		_, err := dbmap.Update(orgActivity)
		CheckErr(err, "UpdateActivity  update failed")
		if err != nil {
			r.JSON(200, Resp{2004, "更新活动失败", nil})
		} else {
			r.JSON(200, Resp{0, "更新活动成功", activity})
		}
	}
}

func DeleteActivity(args martini.Params, r render.Render, dbmap *gorp.DbMap) {
	_, err := dbmap.Exec("DELETE from t_activity WHERE id=?", args["id"])
	CheckErr(err, "DeleteActivity delete failed")
	if err != nil {
		r.JSON(200, Resp{2005, "删除活动失败", nil})
	} else {
		r.JSON(200, Resp{0, "删除活动成功", nil})
	}
}
