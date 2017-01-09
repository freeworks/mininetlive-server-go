package common

import (
	logger "app/logger"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	mathRand "math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	. "github.com/bitly/go-simplejson"
	"github.com/pborman/uuid"
)


func GetLimit(req *http.Request) (int, int) {
	mPageIndex := 0
	mPageSize := 10
	pageSizes := req.URL.Query()["pageSize"]
	if len(pageSizes) > 0 && pageSizes[0] != "null" {
		pageSize, _ := strconv.Atoi(pageSizes[0])
		mPageSize = pageSize
	}
	pageIndexs := req.URL.Query()["pageIndex"]
	if len(pageIndexs) > 0 && pageIndexs[0] != "null" {
		pageIndex, _ := strconv.Atoi(pageIndexs[0])
		mPageIndex = pageIndex - 1
		if mPageIndex < 0 {
			mPageIndex = 0
		}
	}
	start := mPageIndex * mPageSize
	return start, mPageSize
}

func GetTimesampe(req *http.Request) (string, string) {
	var mBeginDate, mEndDate string
	beginDate := req.URL.Query()["beginDate"]
	if len(beginDate) > 0 && beginDate[0] != "null" {
		mBeginDate = beginDate[0]
	}
	endDate := req.URL.Query()["endDate"]
	if len(endDate) > 0 && endDate[0] != "null" {
		mEndDate = endDate[0]
	}
	return mBeginDate, mEndDate
}

func SendSMS(mobile string) (string, error) {
	//TODO 判断电话号码
	url := "https://sms.yunpian.com/v1/sms/send.json"
	vcode := GeneraVCode6()
	text := fmt.Sprintln()
	logger.Info("[Common][SendSMS]",text)
	//uid := "verifyPhone"  //该条短信在您业务系统内的ID
	//callback_url
	res, err := http.Post(url, "application/x-www-form-urlencoded",
		strings.NewReader("apikey=47d1ae5bc2c8f1bf6ed14ac828200299&mobile="+mobile+"&text=【微网直播间】您的验证码为"+vcode+""))
	if err != nil {
		CheckErr("[Common]","[SendSMS]","",err)
		return "", err
	} else {
		data, _ := ioutil.ReadAll(res.Body)
		js, _ := NewJson(data)
		logger.Info("[Common][SendSMS]",string(data))
		code, _ := js.Get("code").Int()
		if code == 0 {
			sid, _ := js.Get("result").Get("sid").Int()
			logger.Info("[Common][SendSMS]",sid)
		} else {
			msg, _ := js.Get("msg").String()
			return "", errors.New(msg)
		}
	}
	return vcode, err
}
func GeneraVCode6() string {
	rnd := mathRand.New(mathRand.NewSource(time.Now().UnixNano()))
	vcode := fmt.Sprintf("%06v", rnd.Int31n(1000000))
	return vcode
}

func RandomStr(length int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := mathRand.New(mathRand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

func CheckErr(tag, method ,msg string,err error) {
	if err != nil {
		logger.Error(tag,method, msg, err)
		// log.Fatalln(msg, err)
	}
}

func GeneraToken8() string {
	token, err := GenerateRandomString(8)
	CheckErr("[common]","[GeneraToken8]","",err)
	return token
}

func GeneraToken16() string {
	token, err := GenerateRandomString(16)
	CheckErr("[common]","[GeneraToken16]","",err )
	return token
}

func GeneraToken32() string {
	token, err := GenerateRandomString(32)
	CheckErr("[common]","[GeneraToken32]","",err)
	return token
}

func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}

func GenerateRandomString(s int) (string, error) {
	b, err := GenerateRandomBytes(s)
	return base64.URLEncoding.EncodeToString(b), err
}

func Token() string {
	now := time.Now().Unix()
	h := md5.New()
	io.WriteString(h, strconv.FormatInt(now, 10))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func Token2() string {
	b := make([]byte, 8)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

func UID() string {
	return Token2()
}

func AID() string {
	return Token2()
}

func UUID() string {
	return uuid.New()
}

func Mkdir(dir string) (e error) {
	_, er := os.Stat(dir)
	logger.Error(er)
	b := er == nil || os.IsExist(er)
	if !b {
		if err := os.MkdirAll(dir, 0777); err != nil {
			if os.IsPermission(err) {
				logger.Error("create dir error:", err.Error())
				e = err
			}
		}
	}
	return
}

func MD5(data string) string {
	t := md5.New()
	io.WriteString(t, data)
	return fmt.Sprintf("%x", t.Sum(nil))
}

func ValidatePhone(phone string) bool {
	//TODO
	return true
}

func ValidatePassword(password string) bool {
	//TODO
	return true
}
