package controller

import (
	. "app/common"
	logger "app/logger"
	. "app/models"
	. "app/push"
	"net/http"
	"strconv"

	"github.com/coopernurse/gorp"
)

// http://cgi.ucloud.com.cn/record_callback?filename=300000039_1462860643.m3u8&filesize=13719488&spacename=record&duration=163
func CallbackRecordFinish(r *http.Request) {
	//TODO
	qs := r.URL.Query()
	logger.Info(qs)
	filename, filesize, spacename, duration := qs.Get("filename"), qs.Get("filesize"), qs.Get("spacename"), qs.Get("duration")
	filesizeInt, err := strconv.Atoi(filesize)
	CheckErr(err, "CallbackRecordFinish : parse filesize ")
	durationInt, err := strconv.Atoi(duration)
	CheckErr(err, "CallbackRecordFinish : parse duration ")
	url := "http://" + spacename + ".ufile.ucloud.com.cn/" + filename
	logger.Info("live record ", url, " size:", filesizeInt, " duration:", durationInt)

}

//[info [map[ip:[116.204.87.129] node:[211.162.55.57] id:[7RUTrhBiMwY=] app:[publish.weiwanglive.com] appname:[mininetlive]]]]
func CallbackLiveBegin(r *http.Request, dbmap *gorp.DbMap) {
	qs := r.URL.Query()
	logger.Info(qs)
	aid := qs.Get("id")
	var activity QActivity
	err := dbmap.SelectOne(&activity, "select * from t_activity where aid =?", aid)
	CheckErr(err, "CallbackLiveBegin - GetActivity select failed")
	if err == nil {
		dbmap.Exec("UPDATE t_activity SET activity_state = 1 WHERE aid = ?", aid)
		PushLiveBegin(activity.Aid, activity.Title)
	}
}

func CallbackLiveEnd(r *http.Request, dbmap *gorp.DbMap) {
	qs := r.URL.Query()
	logger.Info(qs)
	aid := qs.Get("id")
	var activity QActivity
	err := dbmap.SelectOne(&activity, "select * from t_activity where aid =?", aid)
	CheckErr(err, "CallbackLiveEnd - GetActivity select failed")
	if err == nil {
		dbmap.Exec("UPDATE t_activity SET activity_state = 2 WHERE aid = ?", aid)
		PushLiveEnd(activity.Aid, activity.Title)
	}
}
