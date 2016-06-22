package controller

//0    成功
//1000 未注册
//1002 已经注册但获取用户信息失败/用户信息不存在
//1003 注册失败
//1004 账号已经注册
//1005 账号/密码错误
//1006 更新账户信息失败
//1007 删除账户失败
//1008 获取用户信息失败
//1009 获取验证码失败
//1010 验证码输入错误
//1011 验证码过期
//1012 头像上传失败
//1013 uid不能为空


//0 成功
//1100 创建活动失败
//1101 更新活动失败
//1102 删除活动失败
//1103 获取活动失败
//1104 获取活动列表失败
//1105 预约活动失败
//1106 取消预约活动失败

//0 成功
//1300 获取播放记录失败
//1301 获取预约记录失败
//1302 获取支付记录失败


//2000 获取charge失败




//func CancelAppointmentActivity(args martini.Params, r render.Render, dbmap *gorp.DbMap) {
//	var record AppointmentRecord
//	userId := 1 //session from Id
//	err := dbmap.SelectOne(&record, "SELECT * FROM t_appointment_record WHERE activity_id = ? AND user_id = ?",
//		args["activityId"], userId)
//	CheckErr(err, "CancelAppointmentActivity selectOne failed")
//	if err != nil {
//		r.JSON(200, Resp{1106, "取消预约活动失败", nil})
//	} else {
//		record.State = 2
//		_, err := dbmap.Update(record)
//		CheckErr(err, "CancelAppointmentActivity update failed")
//		if err != nil {
//			r.JSON(200, Resp{1106, "取消预约活动失败", nil})
//		} else {
//			r.JSON(200, Resp{0, "更新活动成功", nil})
//		}
//	}
//}

//func PayActivity(req *http.Request, args martini.Params, r render.Render, dbmap *gorp.DbMap) {
//	req.ParseForm()
//	var record PayRecord
//	record.ActivityId, _ = strconv.Atoi(args["id"])
//	record.UserId = 1 //TODO session 取id
//	record.Amount, _ = strconv.Atoi(req.Form["amount"][0])
//	record.Type, _ = strconv.Atoi(req.Form["type"][0])
//	record.Created = time.Now()
//	//TODO校验
//	err := dbmap.Insert(&record)
//	CheckErr(err, "PayActivity insert failed")
//	if err != nil {
//		r.JSON(200, Resp{1105, "支付失败", nil})
//	} else {
//		r.JSON(200, Resp{0, "支付成功", nil})
//	}
//}
