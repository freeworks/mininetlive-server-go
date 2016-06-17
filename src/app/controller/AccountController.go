package controller

import (
	. "app/common"
	. "app/models"
	"time"

	"github.com/coopernurse/gorp"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
)

func GetAccountInfo(args martini.Params, r render.Render, dbmap *gorp.DbMap) {
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

func UpdateAccountName(args martini.Params, user User, r render.Render, dbmap *gorp.DbMap) {
	obj, err := dbmap.Get(User{}, args["id"])
	CheckErr(err, "UpdateUser get failed")
	if err != nil {
		r.JSON(200, Resp{1002, "更新信息失败", nil})
	} else {
		orgUser := obj.(*User)
		if user.Name != "" {
			orgUser.Name = user.Name
		}
		if user.Avatar != "" {
			orgUser.Avatar = user.Avatar
		}
		orgUser.Updated = time.Now()
		_, err := dbmap.Update(orgUser)
		CheckErr(err, "UpdateUser update failed")
		if err != nil {
			r.JSON(200, Resp{1006, "更新信息失败", nil})
		} else {
			r.JSON(200, Resp{0, "更新信息成功", user})
		}
	}
}

func UpdateAccountPhone() {

}

func GetPlayRecordList(r render.Render, dbmap *gorp.DbMap) {
	var playRecords []PlayRecord
	_, err := dbmap.Select(&playRecords, "SELECT * FROM t_play_record ORDER BY create_time")
	CheckErr(err, "GetPlayRecords failed")
	if err != nil {
		r.JSON(200, Resp{1300, "获取播放记录成功", nil})
	} else {
		r.JSON(200, Resp{0, "获取播放记录失败", nil})
	}
}

func GetAppointmentRecordList(r render.Render, dbmap *gorp.DbMap) {
	var appointmentRecords []AppointmentRecord
	_, err := dbmap.Select(&appointmentRecords, "SELECT * FROM t_appointment_record ORDER BY create_time")
	CheckErr(err, "GetPlayRecords failed")
	if err != nil {
		r.JSON(200, Resp{1301, "获取预约记录成功", nil})
	} else {
		r.JSON(200, Resp{0, "获取预约记录失败", nil})
	}
}

func GetPayRecordList(r render.Render, dbmap *gorp.DbMap) {
	var payRecords []PayRecord
	_, err := dbmap.Select(&payRecords, "SELECT * FROM t_pay_record ORDER BY create_time")
	CheckErr(err, "GetPayRecords failed")
	if err != nil {
		r.JSON(200, Resp{1302, "获取支付记录成功", nil})
	} else {
		r.JSON(200, Resp{0, "获取支付记录失败", nil})
	}
}
