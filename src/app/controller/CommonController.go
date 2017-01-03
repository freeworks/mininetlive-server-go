package controller

import (
	. "app/common"
	logger "app/logger"
	"net/http"
	"github.com/coopernurse/gorp"
	"github.com/martini-contrib/render"
)

func GetStartConfig(req *http.Request, r render.Render, dbmap *gorp.DbMap) {
	req.ParseForm()
	platform := req.Header.Get("platform")
	logger.Info("GetStartConfig platform ", platform)
	config := make(map[string]interface{})
	config["isRelease"] = true
	config["enable"] = true
	r.JSON(200, Resp{0, "ok", config})
}