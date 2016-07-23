package controller

import (
	. "app/common"
	config "app/config"
	logger "app/logger"
	. "app/models"
	upload "app/upload"
	"io"
	"net/http"
	"os"
	//	"strings"
	"strconv"
	"time"

	"github.com/coopernurse/gorp"
	"github.com/martini-contrib/render"
	cache "github.com/patrickmn/go-cache"
)

func GetAccountInfo(req *http.Request, r render.Render, dbmap *gorp.DbMap) {
	req.ParseForm()
	uid := req.Header.Get("uid")
	if uid == "" {
		r.JSON(200, Resp{1013, "uid不能为空", nil})
		return
	}
	var user User
	err := dbmap.SelectOne(&user, "SELECT * FROM t_user WHERE uid=?", uid)
	CheckErr(err, "GetUser selectOne failed")
	if err != nil {
		r.JSON(200, Resp{1002, "获取用户信息失败", nil})
	} else {
		r.JSON(200, Resp{0, "获取用户信息成功", user})
	}
}

func UpdateAccountNickName(req *http.Request, r render.Render, dbmap *gorp.DbMap) {
	req.ParseForm()
	name := req.PostFormValue("nickname")
	uid := req.Header.Get("uid")
	if uid == "" {
		r.JSON(200, Resp{1013, "uid不能为空", nil})
		return
	}
	if name == "" {
		r.JSON(200, Resp{1014, "昵称不能为空", nil})
		return
	}
	_, err := dbmap.Exec("UPDATE t_user SET nickname = ? WHERE uid = ?", name, uid)
	CheckErr(err, "Update nickname get failed")
	if err != nil {
		r.JSON(200, Resp{1002, "更新昵称失败,服务器异常", nil})
	} else {
		r.JSON(200, Resp{0, "更新昵称成功", map[string]string{"nickname": name}})
	}
}

func GetVCodeForUpdatePhone(req *http.Request, c *cache.Cache, r render.Render) {
	req.ParseForm()
	phone := req.PostFormValue("phone")
	//TODO 校验
	vCode, err := SendSMS(phone)
	if err != nil {
		r.JSON(200, Resp{1009, "获取验证码失败", nil})
	} else {
		c.Set(phone, vCode, 60*time.Second)
		r.JSON(200, Resp{0, "获取验证码成功", nil})
	}
}

//废弃了
func UpdateAccountPhone(req *http.Request, c *cache.Cache, r render.Render, dbmap *gorp.DbMap) {
	uid := req.Header.Get("uid")
	req.ParseForm()
	phone := req.PostFormValue("phone")
	vCode := req.PostFormValue("vcode")
	if uid == "" {
		r.JSON(200, Resp{1013, "uid不能为空", nil})
		return
	}
	if cacheVCode, found := c.Get(phone); found {
		if cacheVCode.(string) == vCode {
			_, err := dbmap.Exec("UPDATE t_user SET phone = ? WHERE uid = ?", phone, uid)
			CheckErr(err, "Update phone get failed")
			if err != nil {
				r.JSON(200, Resp{1002, "绑定手机失败，服务器异常", nil})
				return
			} else {
				r.JSON(200, Resp{0, "版定手机成功", nil})
				return
			}
		} else {
			r.JSON(200, Resp{1010, "输入验证码有误,请重新输入", nil})
			return
		}
	} else {
		r.JSON(200, Resp{1011, "验证码过期,请重新获取验证码", nil})
		return
	}
}

func UpdateAccountAvatar(req *http.Request, r render.Render, dbmap *gorp.DbMap) {
	uid := req.Header.Get("uid")
	req.ParseMultipartForm(32 << 20)
	logger.Info(req.Header)
	file, head, err := req.FormFile("file")
	CheckErr(err, "upload Fromfile")
	logger.Info(head.Filename)
	defer file.Close()
	fileName := uid + "_" + strconv.FormatInt(time.Now().Unix(), 10) + "_avatar.png"
	filepath := config.ImgDir + fileName
	os.Remove(filepath)
	fW, err := os.Create(filepath)
	CheckErr(err, "create file error")
	defer fW.Close()
	_, err = io.Copy(fW, file)
	CheckErr(err, "create file error")
	fileInfo, err := os.Stat(filepath)
	if err != nil {
		CheckErr(err, "os.Stat(filepath)")
	}
	fileSize := fileInfo.Size()
	logger.Info("fileSize:", fileSize)
	url, err := upload.UploadToUCloudCND(filepath, "avatar/"+fileName)
	if err == nil {
		_, err = dbmap.Exec("UPDATE t_user SET avatar = ? WHERE uid = ?", url, uid)
	}
	if err == nil {
		r.JSON(200, Resp{0, "头像上传成功", map[string]string{"url": url}})
	} else {
		r.JSON(200, Resp{1012, "头像上传失败", nil})
	}
}

func GetPlayRecordList(req *http.Request, r render.Render, dbmap *gorp.DbMap) {
	uid := req.Header.Get("uid")
	var playRecords []QueryPlayRecord
	sql := `SELECT r.id,r.aid,r.uid,r.create_time,a.title,a.front_cover,a.date,a.play_count,u.nickname 
		FROM t_record  r LEFT JOIN t_activity  a  ON r.aid = a.aid  LEFT JOIN t_user u ON a.uid=u.uid 
		WHERE r.type = 1 AND r.uid=? 
		ORDER BY create_time `
	_, err := dbmap.Select(&playRecords, sql, uid)
	CheckErr(err, "GetPlayRecords failed")
	if err != nil {
		r.JSON(200, Resp{1300, "获取播放记录失败", nil})
	} else {
		r.JSON(200, Resp{0, "获取播放记录成功", playRecords})
	}
}

func GetAppointmentRecordList(req *http.Request, r render.Render, dbmap *gorp.DbMap) {
	req.ParseForm()
	uid := req.Header.Get("uid")
	var appointmentRecords []QueryAppointmentRecord
	sql := `SELECT r.id,r.aid,r.uid,r.create_time,a.title,a.activity_state,a.front_cover,a.date,u.nickname
	        FROM t_record  r LEFT JOIN t_activity  a  ON r.aid = a.aid  LEFT JOIN t_user u ON a.uid=u.uid 
		    WHERE r.type = 0 AND r.uid=?
	        ORDER BY create_time`
	_, err := dbmap.Select(&appointmentRecords, sql, uid)
	CheckErr(err, "GetPlayRecords failed")
	if err != nil {
		r.JSON(200, Resp{1301, "获取预约记录失败", nil})
	} else {
		r.JSON(200, Resp{0, "获取预约记录成功", appointmentRecords})
	}
}

func GetPayRecordList(req *http.Request, r render.Render, dbmap *gorp.DbMap) {
	uid := req.Header.Get("uid")
	var payRecords []QueryPayRecord
	sql := `SELECT r.id,r.aid,r.uid,r.create_time,a.front_cover,a.title,a.date,a.activity_type,a.activity_state,u.nickname,o.channel,o.amount 
		    FROM t_record  r LEFT JOIN t_activity  a  ON r.aid = a.aid  LEFT JOIN t_user u ON a.uid=u.uid LEFT JOIN t_order o ON r.orderno=o.no
		    WHERE r.type = 2 AND r.uid=? AND o.type = 1
		    ORDER BY create_time`
	_, err := dbmap.Select(&payRecords, sql, uid)
	CheckErr(err, "GetPayRecords failed")
	if err != nil {
		r.JSON(200, Resp{1302, "获取支付记录失败", nil})
	} else {
		r.JSON(200, Resp{0, "获取支付记录成功", payRecords})
	}
}

func GetWithdrawRecordList(req *http.Request, r render.Render, dbmap *gorp.DbMap) {
	//	uid := req.Header.Get("uid")
	withdrawRecords := [2]QueryWithdrawRecord{
		QueryWithdrawRecord{100, JsonTime{time.Now(), true}},
		QueryWithdrawRecord{100, JsonTime{time.Now(), true}}}
	r.JSON(200, Resp{0, "获取支付记录成功", withdrawRecords})
	//	if err != nil {
	//		r.JSON(200, Resp{1302, "获取支付记录失败", nil})
	//	} else {
	//		r.JSON(200, Resp{0, "获取支付记录成功", withdrawRecords})
	//	}
}
