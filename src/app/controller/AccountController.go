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
	logger.Info("[AccountController]", "[GetAccountInfo]", "uid", uid)
	if uid == "" {
		r.JSON(200, Resp{1013, "uid不能为空", nil})
		return
	}
	var user User
	err := dbmap.SelectOne(&user, "SELECT * FROM t_user WHERE uid=?", uid)
	CheckErr("[AccountController]", "[GetAccountInfo]", "GetUser selectOne failed", err)
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
	logger.Info("[AccountController]", "[UpdateAccountNickName]", "uid->", uid, ",name->"+name)
	if uid == "" {
		r.JSON(200, Resp{1013, "uid不能为空", nil})
		return
	}
	if name == "" {
		r.JSON(200, Resp{1014, "昵称不能为空", nil})
		return
	}
	_, err := dbmap.Exec("UPDATE t_user SET nickname = ? WHERE uid = ?", name, uid)
	CheckErr("[AccountController]", "[UpdateAccountNickName]", "Update nickname get failed", err)
	if err != nil {
		r.JSON(200, Resp{1002, "更新昵称失败,服务器异常", nil})
	} else {
		r.JSON(200, Resp{0, "更新昵称成功", map[string]string{"nickname": name}})
	}
}

func GetVCodeForUpdatePhone(req *http.Request, c *cache.Cache, r render.Render) {
	req.ParseForm()
	phone := req.PostFormValue("phone")
	logger.Info("[AccountController]", "[GetVCodeForUpdatePhone]", "phone->", phone)
	if phone == "" {
		r.JSON(200, Resp{1014, "手机号不能为空", nil})
		return
	}
	vCode, err := SendSMS(phone)
	logger.Info("[AccountController]", "GetVCodeForUpdatePhone", "SendSMS vCode ->", vCode)
	if err != nil {
		r.JSON(200, Resp{1009, "获取验证码失败", nil})
	} else {
		c.Set(phone, vCode, 60*time.Second)
		r.JSON(200, Resp{0, "获取验证码成功", nil})
	}
}

func UpdateAccountPhone(req *http.Request, c *cache.Cache, r render.Render, dbmap *gorp.DbMap) {
	uid := req.Header.Get("uid")
	req.ParseForm()
	phone := req.PostFormValue("phone")
	vCode := req.PostFormValue("vcode")
	logger.Info("[AccountController]", "[UpdateAccountPhone]", " uid", uid, "phone ->", phone, "vCode", vCode)
	if uid == "" {
		r.JSON(200, Resp{1013, "uid不能为空", nil})
		return
	}
	if phone == "" {
		r.JSON(200, Resp{1014, "手机号不能为空", nil})
		return
	}
	if vCode == "" {
		r.JSON(200, Resp{1014, "验证码不能为空", nil})
		return
	}
	if cacheVCode, found := c.Get(phone); found {
		if cacheVCode.(string) == vCode {
			trans, err := dbmap.Begin()
			_, err = dbmap.Exec("UPDATE t_user SET phone = ? WHERE uid = ?", phone, uid)
			CheckErr("[AccountController]", "[UpdateAccountPhone]", "Update phone t_user failed", err)
			_, err = dbmap.Exec("UPDATE t_local_auth SET phone = ? WHERE uid = ? ", phone, uid)
			CheckErr("[AccountController]", "[UpdateAccountPhone]", "Update phone t_local_auth failed", err)
			err = trans.Commit()
			CheckErr("[AccountController]", "[UpdateAccountPhone]", "update phone ", err)
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
	logger.Info("[AccountController]", "[UpdateAccountAvatar]", "uid", uid)
	if uid == "" {
		r.JSON(200, Resp{1013, "uid不能为空", nil})
		return
	}
	req.ParseMultipartForm(32 << 20)
	logger.Info("[AccountController]", "[UpdateAccountAvatar]", req.Header)
	file, head, err := req.FormFile("file")
	CheckErr("[AccountController]", "[UpdateAccountAvatar]", "upload Fromfile", err)
	logger.Info("[AccountController]", "[UpdateAccountAvatar]", head.Filename)
	defer file.Close()
	fileName := uid + "_" + strconv.FormatInt(time.Now().Unix(), 10) + "_avatar.png"
	filepath := config.ImgDir + fileName
	os.Remove(filepath)
	fW, err := os.Create(filepath)
	CheckErr("[AccountController]", "[UpdateAccountAvatar]", "create file", err)
	defer fW.Close()
	_, err = io.Copy(fW, file)
	CheckErr("[AccountController]", "[UpdateAccountAvatar]", "copy file", err)
	fileInfo, err := os.Stat(filepath)
	if err != nil {
		CheckErr("[AccountController]", "[UpdateAccountAvatar]", "os.Stat(filepath)", err)
	}
	fileSize := fileInfo.Size()
	logger.Info("[AccountController]", "[UpdateAccountAvatar]", "fileSize:", fileSize)
	url, err := upload.UploadImageFile(filepath, "avatar/"+fileName)
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
	logger.Info("[AccountController]", "[GetPlayRecordList]", "uid", uid)
	if uid == "" {
		r.JSON(200, Resp{1013, "uid不能为空", nil})
		return
	}
	var playRecords []QueryPlayRecord
	sql := `SELECT t_ar.id,t_ar.aid,t_ar.uid,t_ar.create_time,t_ar.title,t_ar.front_cover,t_ar.date,t_ar.play_count,t_ar.type,u.nickname
			FROM (
				SELECT r.id,r.aid,r.uid ,r.create_time,a.title,a.front_cover,a.date,a.play_count,r.type
			    FROM t_record  r LEFT JOIN t_activity  a  ON r.aid = a.aid
			    ) t_ar JOIN t_user u ON t_ar.uid = u.uid  WHERE t_ar.type = 1 AND u.uid=?
			ORDER BY t_ar.create_time DESC`
	_, err := dbmap.Select(&playRecords, sql, uid)
	CheckErr("[AccountController]", "[GetPlayRecordList]", "", err)
	if err != nil {
		r.JSON(200, Resp{1300, "获取播放记录失败", nil})
	} else {
		r.JSON(200, Resp{0, "获取播放记录成功", playRecords})
	}
}

func GetAppointmentRecordList(req *http.Request, r render.Render, dbmap *gorp.DbMap) {
	req.ParseForm()
	uid := req.Header.Get("uid")
	logger.Info("[AccountController]", "[GetAppointmentRecordList]", "uid", uid)
	if uid == "" {
		r.JSON(200, Resp{1013, "uid不能为空", nil})
		return
	}
	var appointmentRecords []QueryAppointmentRecord
	sql := `SELECT t_ar.id,t_ar.aid,t_ar.uid,t_ar.create_time,t_ar.title,t_ar.activity_state,t_ar.front_cover,t_ar.date,t_ar.type,u.nickname
			FROM (
				SELECT r.id,r.aid,r.uid ,r.create_time,a.title,a.activity_state,a.front_cover,a.date,r.type
			    FROM t_record  r LEFT JOIN t_activity  a  ON r.aid = a.aid
			    ) t_ar JOIN t_user u ON t_ar.uid = u.uid  WHERE t_ar.type = 0 AND u.uid=?
			    
			ORDER BY t_ar.create_time DESC`
	_, err := dbmap.Select(&appointmentRecords, sql, uid)
	CheckErr("[AccountController]", "[GetAppointmentRecordList]", "", err)
	if err != nil {
		r.JSON(200, Resp{1301, "获取预约记录失败", nil})
	} else {
		r.JSON(200, Resp{0, "获取预约记录成功", appointmentRecords})
	}
}

func GetPayRecordList(req *http.Request, r render.Render, dbmap *gorp.DbMap) {
	uid := req.Header.Get("uid")
	logger.Info("[AccountController]", "[GetPayRecordList]", "uid", uid)
	if uid == "" {
		r.JSON(200, Resp{1013, "uid不能为空", nil})
		return
	}
	var payRecords []QueryPayRecord
	sql := `SELECT r.id,r.aid,r.uid,r.create_time,a.front_cover,a.title,a.date,a.activity_type,a.activity_state,u.nickname,o.channel,r.amount 
		    FROM t_record  r LEFT JOIN t_activity  a  ON r.aid = a.aid  LEFT JOIN t_user u ON a.uid=u.uid LEFT JOIN t_order o ON r.orderno=o.no
		    WHERE r.type = 2 AND r.uid=? AND o.type = 1
		    ORDER BY create_time DESC`
	_, err := dbmap.Select(&payRecords, sql, uid)
	CheckErr("[AccountController]", "[GetPayRecordList]", "", err)
	if err != nil {
		r.JSON(200, Resp{1302, "获取支付记录失败", nil})
	} else {
		r.JSON(200, Resp{0, "获取支付记录成功", payRecords})
	}
}

func GetBalance(req *http.Request, r render.Render, dbmap *gorp.DbMap) {
	uid := req.Header.Get("uid")
	logger.Info("[AccountController]", "[GetBalance]", "uid", uid)
	if uid == "" {
		r.JSON(200, Resp{1013, "uid不能为空", nil})
		return
	}
	balance, err := dbmap.SelectInt("select balance from t_user where uid = ?", uid)
	CheckErr("[AccountController]", "[GetBalance]", "", err)
	if err != nil {
		r.JSON(200, Resp{1303, "获取余额失败", nil})
	} else {
		r.JSON(200, Resp{0, "获取余额失败成功", map[string]int64{"balance": balance}})
	}

}

func GetWithdrawRecordList(req *http.Request, r render.Render, dbmap *gorp.DbMap) {
	uid := req.Header.Get("uid")
	logger.Info("[AccountController]", "[GetWithdrawRecordList]", "uid", uid)
	if uid == "" {
		r.JSON(200, Resp{1013, "uid不能为空", nil})
		return
	}

	var withdrawRecords []QueryWithdrawRecord
	sql := `SELECT amount,state,create_time FROM t_record WHERE type = 3 AND uid = ? ORDER BY create_time DESC`
	_, err := dbmap.Select(&withdrawRecords, sql, uid)
	CheckErr("[AccountController]", "[GetWithdrawRecordList]", "GetPayRecords failed", err)
	if err != nil {
		r.JSON(200, Resp{1302, "获取提现记录失败", nil})
	} else {
		r.JSON(200, Resp{0, "获取提现记录成功", withdrawRecords})
	}

}

func GetDividendRecordList(req *http.Request, r render.Render, dbmap *gorp.DbMap) {
	uid := req.Header.Get("uid")
	logger.Info("[AccountController]", "[GetDividendRecordList]", "uid", uid)
	if uid == "" {
		r.JSON(200, Resp{1013, "uid不能为空", nil})
		return
	}
	var dividendRecords []DividendRecord
	_, err := dbmap.Select(&dividendRecords, "SElECT * FROM t_dividend_record WHERE owner_uid = ?", uid)
	CheckErr("[AccountController]", "[GetDividendRecordList]", "", err)
	if err != nil {
		r.JSON(200, Resp{1304, "获取奖励列表失败", nil})
	} else {
		r.JSON(200, Resp{0, "获取奖励列表成功！", dividendRecords})
	}
}
