package pay

import (
	. "app/common"
	config "app/config"
	logger "app/logger"
	. "app/models"
	. "app/push"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	//	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/coopernurse/gorp"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	pingpp "github.com/pingplusplus/pingpp-go/pingpp"
	"github.com/pingplusplus/pingpp-go/pingpp/charge"
	"github.com/pingplusplus/pingpp-go/pingpp/transfer"
)

const (
	//	API_KEY string = `sk_test_GaLOuHePyX1KjnjzbDW5m9KG`
	APP_ID  string = `app_m5y1CSzTGCi5zjjj`
	API_KEY string = `sk_live_SOmLyDHanHSGe5irXPWfnvj5`
)

func newOrder(uid, aid,orderno, channel, clientIP, subject string, amount uint64, payType int) Order {
	return Order{
		Uid:      uid,
		OrderNo:  orderno,
		Amount:   amount,
		Channel:  channel,
		ClientIP: clientIP,
		Subject:  subject,
		Aid:      aid,
		PayType:  payType,
		Created:  JsonTime{time.Now(), true},
	}
}

func newPayRecord(uid, aid, orderno string, amount uint64) Record {
	return Record{
		Uid:     uid,
		OrderNo: orderno,
		Aid:     aid,
		Type:    2,
		Amount:  amount,
	}
}
func newWithdrawRecord(uid ,orderno string, amount uint64) Record {
	return Record{
		Uid:     uid,
		OrderNo: orderno,
		Aid:     "",
		Type:    3,
		Amount:  amount,
	}
}

func init() {
	pingpp.LogLevel = 2
	pingpp.Key = API_KEY
	fmt.Println("Go SDK Version:", pingpp.Version())
	pingpp.AcceptLanguage = "zh-CN"
	//设置商户的私钥 记得在Ping++上配置公钥
	//pingpp.AccountPrivateKey
}

func GetCharge(req *http.Request, parms martini.Params, render render.Render, dbmap *gorp.DbMap) {
	uid := req.Header.Get("uid")
	logger.Info("[Pay]","[GetCharge]","uid:", uid)
	aid := req.PostFormValue("aid")
	if aid == "" {
		render.JSON(200, Resp{1003, "aid不能为空", nil})
		return
	}

	amount, err := strconv.Atoi(req.PostFormValue("amount"))
	CheckErr("[pay]","[GetCharge]","get amount",err)
	if err != nil {
		CheckErr("[Pay]","[GetCharge]","cover amount",err)
		render.JSON(200, Resp{2001, "金额不正确", nil})
		return
	}
	if amount == 0 {
		amount = 1
	}
	payType, err := strconv.Atoi(req.PostFormValue("payType"))
	CheckErr("[pay]","[GetCharge]","get payType",err)
	if err != nil {
		render.JSON(200, Resp{2003, "订单类型不正确，支付类型不正确", nil})
		return
	}
	title, err := dbmap.SelectStr("SELECT title FROM t_activity WHERE aid=?", aid)
	CheckErr("[pay]","[GetCharge]","select title",err)
	if err != nil {
		render.JSON(200, Resp{2003, "订单类型不正确,活动不存在", nil})
		return
	}
	channel := req.PostFormValue("channel")
	if channel != "alipay" && channel != "wx" {
		render.JSON(200, Resp{2002, "支付渠道异常", nil})
		return
	}
	orderno := RandomStr(16)
	logger.Info("[Pay]","[GetCharge]","orderno :", orderno, ",channel:", channel, ",ammout:", amount)
	extra := make(map[string]interface{})
	log.Printf(req.RemoteAddr)
	ip, _, err := net.SplitHostPort(req.RemoteAddr)
	log.Print()
	if err != nil {
		logger.Info("[Pay]","[GetCharge]","userIP: [", req.RemoteAddr, "] is not IP:port")
	}
	userIP := net.ParseIP(ip)
	if userIP == nil {
		logger.Info("[Pay]","[GetCharge]","userIP: [", req.RemoteAddr, "] is not IP:port")
	}

	subject := title
	if payType == 0 {
		subject = title + "(奖赏)"
	}
	body := subject
	params := &pingpp.ChargeParams{
		Order_no:  orderno,
		App:       pingpp.App{Id: APP_ID},
		Amount:    uint64(amount),
		Channel:   channel,
		Currency:  "cny",
		Client_ip: "127.0.0.1", //userIP.String(),
		Subject:   subject,
		Body:      body,
		Extra:     extra,
	}
	ch, err := charge.New(params)
	CheckErr("[pay]","[GetCharge]","charge.New",err)
	if err != nil {
		render.JSON(200, Resp{2000, "获取支付信息失败", nil})
	} else {
		chs, _ := json.Marshal(ch)
		logger.Info("[Pay]","[GetCharge]",string(chs))
		order := newOrder(uid, aid, orderno, channel, userIP.String(), subject, uint64(amount), payType)
		err := dbmap.Insert(&order)
		CheckErr("[pay]","[GetCharge]","create order",err)
		record := newPayRecord(uid,aid,orderno, uint64(amount))
		err = dbmap.Insert(&record)
		CheckErr("[pay]","[GetCharge]","create record",err)
		var chsObj interface{}
		json.Unmarshal(chs, &chsObj)
		render.JSON(200, Resp{0, "获取charge成功", chsObj})
	}

}

func Transfer(req *http.Request, parms martini.Params, render render.Render, dbmap *gorp.DbMap) {
	uid := req.Header.Get("uid")
	if uid == "" {
		render.JSON(200, Resp{1013, "uid不正确", nil})
		return
	}
	var user User
	err := dbmap.SelectOne(&user, "SELECT * FROM t_user WHERE uid=?", uid)
	if err != nil {
		CheckErr("[pay]","[Transfer]","select user",err)
		render.JSON(200, Resp{2001, "服务器异常，查询用户信息失败", nil})
		return
	}
	if user.Phone == "" {
		logger.Info("[Pay]","[Transfer]","还没有绑定手机")
		render.JSON(200, Resp{2006, "还没有绑定手机", nil})
		return
	}

	wxpubOpenId, err := dbmap.SelectStr("SELECT openid FROM t_wxpub WHERE phone=?", user.Phone)
	if err != nil {
		CheckErr("[pay]","[Transfer]","select openid",err)
		render.JSON(200, Resp{2001, "服务器异常，查询用户信息失败", nil})
		return
	}
	logger.Info("[Pay]","[Transfer]","wxpubOpenId ", wxpubOpenId)
	if wxpubOpenId == "" {
		logger.Info("[Pay]","还没有绑定公众账号")
		render.JSON(200, Resp{2007, "还没有绑定公众账号，请关注微网直播公众账号！", nil})
		return
	}

	amount, err := strconv.ParseInt(req.PostFormValue("amount"), 10, 64)
	realAmount := uint64(amount)
	CheckErr("[pay]","[Transfer]","parse amount",err)
	// 1.00 和 20000.00
	if err != nil {
		logger.Info("[Pay]","[Transfer]","金额错误,", realAmount, "账户余额：", user.Balance)
		render.JSON(200, Resp{2008, "金额错误，输入金额不正确！", nil})
		return
	}
	if realAmount >= uint64(user.Balance) {
		render.JSON(200, Resp{2008, "金额错误，账户余额不足！", nil})
		return
	}
	if realAmount < 100 || realAmount > 20000 {
		render.JSON(200, Resp{2008, "金额错误，必须大于1块钱，小于2万块钱", nil})
		return
	}
	extra := make(map[string]interface{})
	extra["user_name"] = user.NickName
	extra["force_check"] = false
	orderno := RandomStr(16)
	logger.Info("[Pay]","[Transfer]","orderno:", orderno, " amount:", realAmount)
	transferParams := &pingpp.TransferParams{
		App:         pingpp.App{Id: APP_ID},
		Channel:     "wx_pub",
		Order_no:    orderno,
		Amount:      realAmount,
		Currency:    "cny",
		Type:        "b2c",
		Recipient:   wxpubOpenId,
		Description: "提现",
		Extra:       extra,
	}
	transfer, err := transfer.New(transferParams)
	if err != nil {
		CheckErr("[pay]","[Transfer]","transfer.New",err)
		render.JSON(200, Resp{2006, "服务器异常，请稍后再试", nil})
		return
	}
	fr, _ := json.Marshal(transfer)
	logger.Info("[Pay]","[Transfer]",string(fr))
	var dat map[string]interface{}
	json.Unmarshal(fr, &dat)
	failure_msg := dat["failure_msg"].(string)
	if failure_msg != "" {
		render.JSON(200, Resp{2009, "服务器异常，请稍后再试", nil})
		return
	}
	_, err = dbmap.Exec("UPDATE t_user SET balance = ? WHERE uid = ?",
		user.Balance-realAmount, user.Uid)
	CheckErr("[pay]","[Transfer]","update user Transfer",err)
	record := newWithdrawRecord(user.Uid, orderno, realAmount)
	err = dbmap.Insert(&record)
	CheckErr("[pay]","[Transfer]","save record Transfer",err)
	render.JSON(200, Resp{0, "提现成功", nil})
	return
}

func Webhook(w http.ResponseWriter, r *http.Request, dbmap *gorp.DbMap) {
	if strings.ToUpper(r.Method) == "POST" {
		buf := new(bytes.Buffer)
		buf.ReadFrom(r.Body)
		//示例 - 签名在头部信息的 x-pingplusplus-signature 字段
		//signed := `BX5sToHUzPSJvAfXqhtJicsuPjt3yvq804PguzLnMruCSvZ4C7xYS4trdg1blJPh26eeK/P2QfCCHpWKedsRS3bPKkjAvugnMKs+3Zs1k+PshAiZsET4sWPGNnf1E89Kh7/2XMa1mgbXtHt7zPNC4kamTqUL/QmEVI8LJNq7C9P3LR03kK2szJDhPzkWPgRyY2YpD2eq1aCJm0bkX9mBWTZdSYFhKt3vuM1Qjp5PWXk0tN5h9dNFqpisihK7XboB81poER2SmnZ8PIslzWu2iULM7VWxmEDA70JKBJFweqLCFBHRszA8Nt3AXF0z5qe61oH1oSUmtPwNhdQQ2G5X3g==`
		signed := r.Header.Get("X-Pingplusplus-Signature")

		//示例 - 待验签的数据
		//data := `{"id":"evt_eYa58Wd44Glerl8AgfYfd1sL","created":1434368075,"livemode":true,"type":"charge.succeeded","data":{"object":{"id":"ch_bq9IHKnn6GnLzsS0swOujr4x","object":"charge","created":1434368069,"livemode":true,"paid":true,"refunded":false,"app":"app_vcPcqDeS88ixrPlu","channel":"wx","order_no":"2015d019f7cf6c0d","client_ip":"140.227.22.72","amount":100,"amount_settle":0,"currency":"cny","subject":"An Apple","body":"A Big Red Apple","extra":{},"time_paid":1434368074,"time_expire":1434455469,"time_settle":null,"transaction_no":"1014400031201506150354653857","refunds":{"object":"list","url":"/v1/charges/ch_bq9IHKnn6GnLzsS0swOujr4x/refunds","has_more":false,"data":[]},"amount_refunded":0,"failure_code":null,"failure_msg":null,"metadata":{},"credential":{},"description":null}},"object":"event","pending_webhooks":0,"request":"iar_Xc2SGjrbdmT0eeKWeCsvLhbL"}`
		data := buf.String()
		// 请从 https://dashboard.pingxx.com 获取「Ping++ 公钥」
		publicKey, err := ioutil.ReadFile("./pingpp_rsa_public_key.pem")
		CheckErr("[Pay]","[webhook]","pingpp.Verify",err)
		if err != nil {
			return
		}
		//base64解码再验证
		decodeStr, _ := base64.StdEncoding.DecodeString(signed)
		//		logger.Info(TAG,decodeStr)
		err = pingpp.Verify([]byte(data), publicKey, decodeStr)
		CheckErr("[Pay]","[webhook]","pingpp.Verify",err)
		if err != nil {
			return
		}
		logger.Info("[Pay]","[Webhook]","ping++ verify hook success")
		webhook, err := pingpp.ParseWebhooks(buf.Bytes())
		logger.Info("[Pay]","[Webhook]","webhook.Type:" + webhook.Type)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "fail")
			return
		}

		if webhook.Type == "charge.succeeded" {
			//orderno -> uid -> invited1 -> invite2->invited3
			logger.Info("[Pay]","[Webhook]",webhook.Type, "id:", webhook.Data.Object["id"],
				"transaction_no:", webhook.Data.Object["transaction_no"],
				"order_no:", webhook.Data.Object["order_no"],
				"amount:", webhook.Data.Object["amount"])
			orderNo := webhook.Data.Object["order_no"].(string)
			amount, _ := webhook.Data.Object["amount"].(json.Number).Int64()
			var order Order
			err = dbmap.SelectOne(&order, "SELECT * FROM t_order t WHERE t.no=?", orderNo)
			CheckErr("[Pay]","[webhook]","select order by no",err)
			if err == nil {
				logger.Info("[Pay]","[Webhook]","UPDATE t_record SET state = 1 WHERE orderno = " + orderNo + " AND type = 2")
				//update pay record
				_, err := dbmap.Exec("UPDATE t_record SET state = 1 WHERE orderno = ? AND type = 2", orderNo)
				CheckErr("[Pay]","[Webhook]","update pay record",err)
				var user User
				err = dbmap.SelectOne(&user, "SELECT * FROM t_user WHERE uid=?", order.Uid)
				CheckErr("[Pay]","[Webhook]","select user", err)
				var title string
				err = dbmap.SelectOne(&title, "SELECT title FROM t_activity WHERE aid=?", order.Aid)
				CheckErr("[Pay]","[Webhook]","select activity", err)
				if err == nil {
					// https://www.zhihu.com/question/29083902
					obj1, obj2, obj3 := getDistributionUsers(user.Uid, dbmap)
					logger.Info("[Pay]","[Webhook]","getDistributionUsers ", obj1, obj2, obj3)
					if obj1 != nil {
						user1 := obj1.(User)
						dividend := uint64(float64(amount) * config.DeductPercent1)
						logger.Info("[Pay]","[Webhook]","user1 dividend->", dividend)
						user1.Balance = user1.Balance + dividend
						_, err := dbmap.Exec("UPDATE t_user SET balance = ? WHERE uid = ?", user1.Balance, user1.Uid)
						CheckErr("[Pay]","[Webhook]","update user1 dividend", err)
						newDividendRecord(dbmap, user.Uid, user.NickName, user.Avatar, order.Aid, title, dividend, user1.Uid, user1.DeviceId)
					}
					if obj2 != nil {
						user2 := obj2.(User)
						dividend := uint64(float64(amount) * config.DeductPercent2)
						user2.Balance = user2.Balance + dividend
						logger.Info("[Pay]","[Webhook]"," user2 dividend->", dividend)
						_, err := dbmap.Exec("UPDATE t_user SET balance = ? WHERE uid = ?", user2.Balance, user2.Uid)
						CheckErr("[Pay]","[Webhook]","update user2 dividend", err)
						newDividendRecord(dbmap, user.Uid, user.NickName, user.Avatar, order.Aid, title, dividend, user2.Uid, user2.DeviceId)
					}

					if obj3 != nil {
						user3 := obj3.(User)
						dividend := uint64(float64(amount) * config.DeductPercent3)
						user3.Balance = user3.Balance + dividend
						logger.Info("[Pay]","[webhook]","user2 dividend->", dividend)
						_, err := dbmap.Exec("UPDATE t_user SET balance = ? WHERE uid = ?", user3.Balance, user3.Uid)
						CheckErr("[Pay]","[Webhook]","update user3 dividend", err)
						newDividendRecord(dbmap, user.Uid, user.NickName, user.Avatar, order.Aid, title, dividend, user3.Uid, user3.DeviceId)
					}
				}
			}
			w.WriteHeader(http.StatusOK)
		} else if webhook.Type == "transfer.succeeded" {
			orderNo := webhook.Data.Object["order_no"].(string)
			logger.Info("[Pay]","[webhook]","UPDATE t_record SET state = 1 WHERE orderno = " + orderNo + " AND type = 3")
			_, err := dbmap.Exec("UPDATE t_record SET state = 1 WHERE orderno = ? AND type = 3", orderNo)
			CheckErr("[Pay]","[Webhook]","update transfer record", err)
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func newDividendRecord(dbmap *gorp.DbMap, uid, nickname, avatar, aid, title string, amount uint64, ownerUid string, deviceId string) {
	r := DividendRecord{
		Uid:      uid,
		NickName: nickname,
		Avatar:   avatar,
		Aid:      aid,
		Title:    title,
		Amount:   amount,
		OwnerUid: ownerUid,
		Created:  JsonTime{time.Now(), true},
	}
	err := dbmap.Insert(&r)
	CheckErr("[Pay]","[newDividendRecord]","", err)
	if err != nil {
		PushDividend(deviceId)
	}
	return
}

func getDistributionUsers(uid string, dbmap *gorp.DbMap) (interface{}, interface{}, interface{}) {
	logger.Info("[Pay]","[getDistributionUsers]")
	var user1, user2, user3 User
	err := dbmap.SelectOne(&user1, `SELECT * FROM t_user 
									WHERE invite_code = 
									(SELECT be_invited_code FROM t_invite_relation WHERE uid = ?) `, uid)
	CheckErr("[Pay]","[getDistributionUsers]","query user1", err)
	if err != nil || isEmptyUser(user1) {
		return nil, nil, nil
	}
	err = dbmap.SelectOne(&user2, `SELECT * FROM t_user 
									WHERE invite_code = 
									(SELECT be_invited_code FROM t_invite_relation WHERE uid = ?) `, user1.Uid)
	CheckErr("[Pay]","[getDistributionUsers]","query user2", err)
	if err != nil || isEmptyUser(user2) {
		return user1, nil, nil
	}

	err = dbmap.SelectOne(&user3, `SELECT * FROM t_user 
									WHERE invite_code = 
									(SELECT be_invited_code FROM t_invite_relation WHERE uid = ?) `, user2.Uid)
	CheckErr("[Pay]","[getDistributionUsers]","query user3", err)
	if err != nil || isEmptyUser(user3) {
		logger.Info("[Pay]","[getDistributionUsers]","webhook,二级关系")
		return user1, user2, nil
	} else {
		logger.Info("[Pay]","[getDistributionUsers]","webhook,三级关系")
		return user1, user2, user3
	}
}

func isEmptyUser(user User) bool {
	return user.Uid == ""
}
