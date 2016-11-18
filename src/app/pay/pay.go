package pay

import (
	. "app/common"
	config "app/config"
	logger "app/logger"
	. "app/models"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
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

func newOrder(uid string, orderno, channel, clientIP, subject, aid string, amount uint64, payType int) Order {
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

func newPayRecord(uid string, aid string, orderno string) Record {
	return Record{
		Uid:     uid,
		OrderNo: orderno,
		Aid:     aid,
		Type:    2,
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
	logger.Info("uid:", uid)
	amount, err := strconv.Atoi(req.PostFormValue("amount"))
	CheckErr(err, "get amount")
	if err != nil {
		render.JSON(200, Resp{2001, "金额不正确", nil})
		return
	}
	if amount == 0 {
		amount = 1
	}
	payType, err := strconv.Atoi(req.PostFormValue("payType"))
	CheckErr(err, "get payType")
	if err != nil {
		render.JSON(200, Resp{2003, "订单类型不正确，支付类型不正确", nil})
		return
	}
	aid := req.PostFormValue("aid")
	title, err := dbmap.SelectStr("SELECT title FROM t_activity WHERE aid=?", aid)
	if err != nil {
		render.JSON(200, Resp{2003, "订单类型不正确,活动不存在", nil})
		return
	}
	channel := req.PostFormValue("channel")
	if channel != "alipay" && channel != "wx" {
		render.JSON(200, Resp{2002, "支付渠道异常", nil})
		return
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	orderno := r.Intn(999999999999999)
	logger.Info("ordderno :", orderno, ",channel:", channel, ",ammout:", amount)
	extra := make(map[string]interface{})
	log.Printf(req.RemoteAddr)
	ip, _, err := net.SplitHostPort(req.RemoteAddr)
	log.Print()
	if err != nil {
		log.Print("userIP: [", req.RemoteAddr, "] is not IP:port")
	}
	userIP := net.ParseIP(ip)
	if userIP == nil {
		log.Print("userIP: [", req.RemoteAddr, "] is not IP:port")
	}

	subject := title
	if payType == 0 {
		subject = title + "(奖赏)"
	}
	body := "sjfdlfjlsdjflsdfjdslf"
	params := &pingpp.ChargeParams{
		Order_no:  strconv.Itoa(orderno),
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
	if err != nil {
		logger.Info(err)
		render.JSON(200, Resp{2000, "获取支付信息失败", nil})
	} else {
		chs, _ := json.Marshal(ch)
		logger.Info(string(chs))
		order := newOrder(uid, strconv.Itoa(orderno), channel, userIP.String(), subject, aid, uint64(amount), payType)
		err := dbmap.Insert(&order)
		record := newPayRecord(aid, uid, strconv.Itoa(orderno))
		err = dbmap.Insert(&record)
		CheckErr(err, "create order")
		var chsObj interface{}
		json.Unmarshal(chs, &chsObj)
		render.JSON(200, Resp{0, "获取charge成功", chsObj})
	}

}

func Transfer(req *http.Request, parms martini.Params, render render.Render, dbmap *gorp.DbMap) {
	uid := req.Header.Get("uid")
	if uid == "" {
		render.JSON(200, Resp{1013, "uidb不正确", nil})
		return
	}
	var oauth OAuth
	err := dbmap.SelectOne(&oauth, "SELECT * FROM t_oauth WHERE uid=?", uid)
	if err != nil {
		render.JSON(200, Resp{1013, "服务器异常，查询用户信息失败", nil})
		return
	}
	if oauth.Plat != "Wechat" {
		render.JSON(200, Resp{2005, "用户还没有开通微信", nil})
		return
	}
	amount, err := strconv.ParseInt(req.PostFormValue("amount"), 10, 64)
	amount2 := uint64(amount)
	if err != nil {
		logger.Info(err)
		render.JSON(200, Resp{2005, "金额不正确", nil})
		return
	}

	//TODO 校验金额类型，以及用户余额是否可以提现

	openId := oauth.OpenId
	extra := make(map[string]interface{})
	extra["user_name"] = "user name"
	extra["force_check"] = false
	//这里是随便设置的随机数作为订单号，仅作示例，该方法可能产生相同订单号，商户请自行生成订单号，不要纠结该方法。
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	orderno := r.Intn(999999999999999)
	transferParams := &pingpp.TransferParams{
		App:         pingpp.App{Id: APP_ID},
		Channel:     "wx_pub",
		Order_no:    strconv.Itoa(orderno),
		Amount:      amount2,
		Currency:    "cny",
		Type:        "b2c",
		Recipient:   openId,
		Description: "Your Description",
		Extra:       extra,
	}
	transfer, err := transfer.New(transferParams)
	if err != nil {
		log.Fatal(err)
		render.JSON(200, Resp{2006, "提现失败", nil})
		return
	}
	logger.Info(transfer)
	fr, _ := json.Marshal(transfer)
	logger.Info(string(fr))
	var dat map[string]interface{}
	json.Unmarshal(fr, &dat)
	failure_msg := dat["failure_msg"].(string)
	if failure_msg != `` {
		render.JSON(200, Resp{2007, "未绑定公众号", nil})
		return
	}
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
		CheckErr(err, "pingpp.Verify")
		if err != nil {
			return
		}
		//base64解码再验证
		decodeStr, _ := base64.StdEncoding.DecodeString(signed)
		logger.Info(decodeStr)
		err = pingpp.Verify([]byte(data), publicKey, decodeStr)
		CheckErr(err, "webhook pingpp.Verify")
		if err != nil {
			return
		}
		logger.Info("webhook ping++ verify hook success")
		webhook, err := pingpp.ParseWebhooks(buf.Bytes())
		logger.Info("webhook webhook.Type:" + webhook.Type)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "fail")
			return
		}

		if webhook.Type == "charge.succeeded" {
			//orderno -> uid -> invited1 -> invite2->invited3
			logger.Info("webhook ", webhook.Type, "id:", webhook.Data.Object["id"],
				"transaction_no:", webhook.Data.Object["transaction_no"],
				"order_no:", webhook.Data.Object["order_no"],
				"amount:", webhook.Data.Object["amount"])
			orderNo := webhook.Data.Object["order_no"].(string)
			amount, _ := webhook.Data.Object["amount"].(json.Number).Int64()
			var order Order
			err = dbmap.SelectOne(&order, "SELECT * FROM t_order t WHERE t.no=?", orderNo)
			CheckErr(err, "webhook select order by no ")
			if err == nil {
				var user User
				err = dbmap.SelectOne(&user, "SELECT * FROM t_user WHERE uid=?", order.Uid)
				CheckErr(err, "webhook select user")
				var title string
				err = dbmap.SelectOne(&user, "SELECT title FROM t_activity WHERE aid=?", order.Aid)
				CheckErr(err, "webhook select activity")
				if err == nil {
					// https://www.zhihu.com/question/29083902
					obj1, obj2, obj3 := getDistributionUsers(user.Uid, dbmap)
					if obj1 != nil {
						user1 := obj1.(User)
						dividend := int(float64(amount) * config.DeductPercent1)
						user1.Balance = user1.Balance + dividend
						_, err := dbmap.Update(user1)
						CheckErr(err, "update user1 dividend")
						newDividendRecord(dbmap, user.Uid, user.NickName, user.Avatar, order.Aid, title, dividend, user1.Uid)
					}
					if obj2 != nil {
						user2 := obj2.(User)
						dividend := int(float64(amount) * config.DeductPercent2)
						user2.Balance = user2.Balance + dividend
						_, err = dbmap.Update(user2)
						CheckErr(err, "update user2 dividend")
						newDividendRecord(dbmap, user.Uid, user.NickName, user.Avatar, order.Aid, title, dividend, user2.Uid)
					}

					if obj3 != nil {
						user3 := obj3.(User)
						dividend := int(float64(amount) * config.DeductPercent3)
						user3.Balance = user3.Balance + dividend
						_, err = dbmap.Update(user3)
						CheckErr(err, "update user3 dividend")
						newDividendRecord(dbmap, user.Uid, user.NickName, user.Avatar, order.Aid, title, dividend, user3.Uid)
					}
				}
			}
			w.WriteHeader(http.StatusOK)
		} else if webhook.Type == "transfer.succeeded" {
			//TODO
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func newDividendRecord(dbmap *gorp.DbMap, uid, nickname, avatar, aid, title string, amount int, ownerUid string) {
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
	CheckErr(err, "newDividendRecord")
	return
}

func getDistributionUsers(uid string, dbmap *gorp.DbMap) (interface{}, interface{}, interface{}) {
	logger.Info("getDistributionUsers...")
	var user1, user2, user3 User
	err := dbmap.SelectOne(&user1, `SELECT * FROM t_user 
									WHERE invite_code = 
									(SELECT be_invited_code FROM t_invite_relation WHERE uid = ?) `, uid)
	CheckErr(err, "webhook getDistributionUsers query user1")
	if err != nil {
		return nil, nil, nil
	}
	err = dbmap.SelectOne(&user2, `SELECT * FROM t_user 
									WHERE invite_code = 
									(SELECT be_invited_code FROM t_invite_relation WHERE uid = ?) `, user1.Uid)
	CheckErr(err, "webhook getDistributionUsers query user2")
	if err != nil {
		return user1, nil, nil
	}

	err = dbmap.SelectOne(&user3, `SELECT * FROM t_user 
									WHERE invite_code = 
									(SELECT be_invited_code FROM t_invite_relation WHERE uid = ?) `, user2.Uid)
	CheckErr(err, "webhook getDistributionUsers query user3")
	if err != nil {
		return user1, user2, user3
	} else {
		return user1, user2, nil
	}
}

// {
//     "id": "evt_eYa58Wd44Glerl8AgfYfd1sL",
//     "created": 1434368075,
//     "livemode": true,
//     "type": "charge.succeeded",
//     "data": {
//         "object": {
//             "id": "ch_bq9IHKnn6GnLzsS0swOujr4x",
//             "object": "charge",
//             "created": 1434368069,
//             "livemode": true,
//             "paid": true,
//             "refunded": false,
//             "app": "app_vcPcqDeS88ixrPlu",
//             "channel": "wx",
//             "order_no": "2015d019f7cf6c0d",
//             "client_ip": "140.227.22.72",
//             "amount": 100,
//             "amount_settle": 0,
//             "currency": "cny",
//             "subject": "An Apple",
//             "body": "A Big Red Apple",
//             "extra": { },
//             "time_paid": 1434368074,
//             "time_expire": 1434455469,
//             "time_settle": null,
//             "transaction_no": "1014400031201506150354653857",
//             "refunds": {
//                 "object": "list",
//                 "url": "/v1/charges/ch_bq9IHKnn6GnLzsS0swOujr4x/refunds",
//                 "has_more": false,
//                 "data": [ ]
//             },
//             "amount_refunded": 0,
//             "failure_code": null,
//             "failure_msg": null,
//             "metadata": { },
//             "credential": { },
//             "description": null
//         }
//     },
//     "object": "event",
//     "pending_webhooks": 0,
//     "request": "iar_Xc2SGjrbdmT0eeKWeCsvLhbL"
// }
