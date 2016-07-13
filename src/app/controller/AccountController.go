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

//Fixme 注册上传头像的时候还没有Uid
func UploadAccountAvatar(req *http.Request, r render.Render) {
	err := req.ParseMultipartForm(100000)
	CheckErr(err, "upload ParseMultipartForm")
	if err != nil {
		r.JSON(500, "server err")
		return
	}
	file, head, err := req.FormFile("file")
	CheckErr(err, "upload Fromfile")
	logger.Info(head.Filename)
	defer file.Close()
	//TODO bug
	fileName := GeneraToken16() + "_avatar.png"
	filepath := config.ImgDir + fileName
	fW, err := os.Create(filepath)
	CheckErr(err, "create file error")
	defer fW.Close()
	_, err = io.Copy(fW, file)
	CheckErr(err, "create file error")
	url, err := upload.UploadToUCloudCND(filepath, fileName, r)
	if err == nil {
		r.JSON(200, Resp{0, "头像上传成功", map[string]string{"url": url}})
	} else {
		r.JSON(200, Resp{1012, "头像上传失败", nil})
	}
}

type QueryPlayRecord struct {
	Record
	FrontCover string   `db:"front_cover" json:"frontCover"`
	Title      string   `db:"title" json:"title"`
	NickName   string   `db:"nickname" json:"nickname"`
	PlayCount  int      `db:"play_count" json:"playCount"`
	Date       JsonTime `db:"date" json:"date"`
}

type QueryPayRecord struct {
	Record
	FrontCover    string   `db:"front_cover" json:"frontCover"`
	Title         string   `db:"title" json:"title"`
	NickName      string   `db:"nickname" json:"nickname"`
	ActivityType  int      `db:"activity_type" json:"activityType"`
	ActivityState int      `db:"activity_state" json:"activityState"`
	Amount        int      `db:"amount" json:"amount"`
	Channel       string   `db:"channel" json:"channel"`
	Date          JsonTime `db:"date" json:"date"`
}

type QueryAppointmentRecord struct {
	Record
	FrontCover    string   `db:"front_cover" json:"frontCover"`
	Title         string   `db:"title" json:"title"`
	NickName      string   `db:"nickname" json:"nickname"`
	ActivityState int      `db:"activity_state" json:"activityState"`
	Date          JsonTime `db:"date" json:"date"`
}

type QueryWithdrawRecord struct {
	Amount  int      `db:"amount" json:"amount"`
	Created JsonTime `db:"create_time" json:"createTime"`
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
