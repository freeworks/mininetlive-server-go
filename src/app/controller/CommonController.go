package controller

import (
	. "app/common"
	. "app/models"
	"net/http"

	"github.com/coopernurse/gorp"
	"github.com/martini-contrib/render"
)

func PostInviteCode(req *http.Request, r render.Render, dbmap *gorp.DbMap) {
	req.ParseForm()
	uid := req.PostFormValue("phone")
	inviteCode := req.PostFormValue("inviteCode")
	var beInvitedUser User
	err := dbmap.SelectOne(&beInvitedUser, "SELECT * FROM t_user WHERE invite_code=?", inviteCode)
	CheckErr(err, "Login selectOne failed")
	if err == nil {
		if beInvitedUser.BeInvitedUid == "" {
			_, err := dbmap.Exec("UPDATE t_user SET be_invited_uid = ? WHERE uid = ?", beInvitedUser.Uid, uid)
			CheckErr(err, "Update phone get failed")
		}
		r.JSON(200, Resp{0, "提交成功", nil})
	} else {
		r.JSON(200, Resp{0, "邀请码不存在！", nil})
	}
}
