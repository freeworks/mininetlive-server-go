package controller

import (
	. "app/common"
	. "app/models"

	"github.com/coopernurse/gorp"
	"github.com/martini-contrib/render"
)

//TODO 分页

//0 成功
//1300 获取播放记录失败
//1301 获取预约记录失败
//1302 获取支付记录失败

func GetPlayRecords(r render.Render, dbmap *gorp.DbMap) {
	var playRecords []PlayRecord
	_, err := dbmap.Select(&playRecords, "SELECT * FROM t_play_record ORDER BY create_time")
	CheckErr(err, "GetPlayRecords failed")
	if err != nil {
		r.JSON(200, Resp{1300, "获取播放记录成功", nil})
	} else {
		r.JSON(200, Resp{0, "获取播放记录失败", nil})
	}
	//	log.Println("All rows:")
	//	for x, p := range posts {
	//		log.Printf("    %d: %v\n", x, p)
	//	}
}

func GetAppointmentRecords(r render.Render, dbmap *gorp.DbMap) {
	var appointmentRecords []AppointmentRecord
	_, err := dbmap.Select(&appointmentRecords, "SELECT * FROM t_appointment_record ORDER BY create_time")
	CheckErr(err, "GetPlayRecords failed")
	if err != nil {
		r.JSON(200, Resp{1301, "获取预约记录成功", nil})
	} else {
		r.JSON(200, Resp{0, "获取预约记录失败", nil})
	}
}

func GetPayRecords(r render.Render, dbmap *gorp.DbMap) {
	var payRecords []PayRecord
	_, err := dbmap.Select(&payRecords, "SELECT * FROM t_pay_record ORDER BY create_time")
	CheckErr(err, "GetPayRecords failed")
	if err != nil {
		r.JSON(200, Resp{1302, "获取支付记录成功", nil})
	} else {
		r.JSON(200, Resp{0, "获取支付记录失败", nil})
	}
}
