package controller

import (
	. "app/models"
	"net/http"

	"github.com/coopernurse/gorp"
	"github.com/martini-contrib/render"
)

func PostInviteCode(req *http.Request, r render.Render, dbmap *gorp.DbMap) {
	r.JSON(200, Resp{0, "提交成功", nil})
}
