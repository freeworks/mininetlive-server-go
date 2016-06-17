package controller

import (
	. "app/common"
	. "app/models"

	"github.com/coopernurse/gorp"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
)

func GetUser(args martini.Params, r render.Render, dbmap *gorp.DbMap) {
	var user User
	err := dbmap.SelectOne(&user, "select * from t_user where id=?", args["id"])
	CheckErr(err, "GetUser selectOne failed")
	//simple error check
	if err != nil {
		r.JSON(200, Resp{1002, "获取用户信息失败", nil})
	} else {
		r.JSON(200, Resp{0, "获取用户信息成功", user})
	}
}
