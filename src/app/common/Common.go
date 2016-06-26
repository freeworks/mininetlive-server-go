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

func SendSMS(mobile string) (string, error) {
	//TODO 判断电话号码
	url := "https://sms.yunpian.com/v1/sms/send.json"
	vcode := GeneraVCode6()
	text := fmt.Sprintln()
	logger.Info(text)
	//uid := "verifyPhone"  //该条短信在您业务系统内的ID
	//callback_url
	res, err := http.Post(url, "application/x-www-form-urlencoded",
		strings.NewReader("apikey=47d1ae5bc2c8f1bf6ed14ac828200299&mobile="+mobile+"&text=【微网直播间】您的验证码是"+vcode))
	if err != nil {
		CheckErr(err, "GetSMSCode send msg")
		return "", err
	} else {
		data, _ := ioutil.ReadAll(res.Body)
		js, _ := NewJson(data)
		logger.Info(string(data))
		code, _ := js.Get("code").Int()
		if code != 0 {
			sid, _ := js.Get("result").Get("sid").Int()
			logger.Info(sid)
		} else {
			msg, _ := js.Get("msg").String()
			return "", errors.New(msg)
		}
	}
	return vcode, err
}
func GeneraVCode6() string {
	rnd := mathRand.New(mathRand.NewSource(time.Now().UnixNano()))
	vcode := fmt.Sprintf("%03v", rnd.Int31n(1000000))
	return vcode
}

func CheckErr(err error, msg string) {
	if err != nil {
		logger.Error(msg, err)
		// log.Fatalln(msg, err)
	}
}

func GeneraToken8() string {
	token, err := GenerateRandomString(8)
	CheckErr(err, "GeneraToken16")
	return token
}

func GeneraToken16() string {
	token, err := GenerateRandomString(16)
	CheckErr(err, "GeneraToken16")
	return token
}

func GeneraToken32() string {
	token, err := GenerateRandomString(32)
	CheckErr(err, "GeneraToken32")
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
