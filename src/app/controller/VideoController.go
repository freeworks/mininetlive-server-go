package controller

import (
	"net/http"
	. "app/common"
	logger "app/logger"
	"strconv"
)


// http://cgi.ucloud.com.cn/record_callback?filename=300000039_1462860643.m3u8&filesize=13719488&spacename=record&duration=163
func  CallbackRecordFinish(r *http.Request) {
	qs := r.URL.Query()
	filename, filesize, spacename,duration := qs.Get("filename"), qs.Get("filesize"), qs.Get("spacename"),qs.Get("duration")
	filesizeInt, err := strconv.Atoi(filesize)
	CheckErr(err,"CallbackRecordFinish : parse filesize ")
	durationInt, err := strconv.Atoi(duration)
	CheckErr(err,"CallbackRecordFinish : parse duration ")
	url := "http://"+spacename+".ufile.ucloud.com.cn/"+filename
	logger.Info("live record ",url," size:",filesizeInt," duration:",durationInt)
	//TODO 更新状态 直播结束，变成点播
}

// http://127.0.0.1/publish_start?ip=推流端IP&id=流名&node=节点IP&app=推流域名&appname=发布点
func  CallbackLiveBegin(r *http.Request) {
	qs := r.URL.Query()
	filename, filesize, spacename,duration := qs.Get("filename"), qs.Get("filesize"), qs.Get("spacename"),qs.Get("duration")
	filesizeInt, err := strconv.Atoi(filesize)
	CheckErr(err,"CallbackRecordFinish : parse filesize ")
	durationInt, err := strconv.Atoi(duration)
	CheckErr(err,"CallbackRecordFinish : parse duration ")
	url := "http://"+spacename+".ufile.ucloud.com.cn/"+filename
	logger.Info("live record ",url," size:",filesizeInt," duration:",durationInt)
	//TODO 更新状态 直播结束，变成点播
}

// http://127.0.0.1/publish_stop?ip=推流端IP&id=流名&node=节点IP&app=推流域名&appname=发布点 其中，publish_start和publish_stop为客户提供的回调cgi。
func  CallbackLiveEnd(r *http.Request) {
	qs := r.URL.Query()
	filename, filesize, spacename,duration := qs.Get("filename"), qs.Get("filesize"), qs.Get("spacename"),qs.Get("duration")
	filesizeInt, err := strconv.Atoi(filesize)
	CheckErr(err,"CallbackRecordFinish : parse filesize ")
	durationInt, err := strconv.Atoi(duration)
	CheckErr(err,"CallbackRecordFinish : parse duration ")
	url := "http://"+spacename+".ufile.ucloud.com.cn/"+filename
	logger.Info("live record ",url," size:",filesizeInt," duration:",durationInt)
	//TODO 更新状态 直播结束，变成点播
}