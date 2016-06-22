package pay

import (
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
	"strconv"
	"strings"
	"time"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	pingpp "github.com/pingplusplus/pingpp-go/pingpp"
	"github.com/pingplusplus/pingpp-go/pingpp/charge"
	"github.com/pingplusplus/pingpp-go/pingpp/redEnvelope"
)

const (
	API_KEY string = `sk_test_GaLOuHePyX1KjnjzbDW5m9KG`
	APP_ID  string = `app_m5y1CSzTGCi5zjjj`
	//	API_KEY    string = `sk_live_SOmLyDHanHSGe5irXPWfnvj5`
	WX_OPEN_ID string = `wx36d2981a085f6370`
)

func init() {
	pingpp.LogLevel = 2
	pingpp.Key = API_KEY
	fmt.Println("Go SDK Version:", pingpp.Version())
	pingpp.AcceptLanguage = "zh-CN"
	//设置商户的私钥 记得在Ping++上配置公钥
	//pingpp.AccountPrivateKey
}

func GetCharge(req *http.Request, parms martini.Params, render render.Render) {
	amount, err := strconv.Atoi(req.FormValue("amount"))
	if err != nil {
		render.JSON(200, Resp{2001, "支付金额不正确", nil})
		return
	}
	if amount == 0 {
		amount = 1
	}
	channel := req.FormValue("channel")
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
	if err != nil {
		log.Print("userIP: [", req.RemoteAddr, "] is not IP:port")
	}
	userIP := net.ParseIP(ip)
	if userIP == nil {
		log.Print("userIP: [", req.RemoteAddr, "] is not IP:port")
	}

	params := &pingpp.ChargeParams{
		Order_no:  strconv.Itoa(orderno),
		App:       pingpp.App{Id: APP_ID},
		Amount:    uint64(amount),
		Channel:   channel,
		Currency:  "cny",
		Client_ip: userIP.String(),
		Subject:   "Your Subject",
		Body:      "Your Body",
		Extra:     extra,
	}
	ch, err := charge.New(params)
	if err != nil {
		logger.Info(err)
		render.JSON(200, Resp{2000, "获取支付信息失败", nil})
	} else {
		chs, _ := json.Marshal(ch)
		//fmt.Fprintln(w, string(chs))
		//TODO 创建订单
		render.JSON(200, Resp{0, "获取charge成功", chs})
	}
}

type RedEnvelopeModel struct {
	Nickname    string
	SendName    string
	Subject     string
	Body        string
	Description string
	Amount      uint64
}

func Withdraw(req *http.Request, redEnvelope RedEnvelopeModel, render render.Render) {
	//TODO 是否绑定微信,
	//TODO 是否绑定电话
	result, err := NewRedEnvelope(redEnvelope)
	if err == nil {
		render.JSON(200, Resp{2003, "提现成功", result})
	} else {
		logger.Info(err)
		render.JSON(200, Resp{0, "提现失败", nil})
	}
}

func NewRedEnvelope(redEnvelopeModel RedEnvelopeModel) (interface{}, error) {
	extra := make(map[string]interface{})
	extra["nick_name"] = redEnvelopeModel.Nickname
	extra["send_name"] = redEnvelopeModel.SendName
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	orderno := r.Intn(999999999999999)
	redenvelopeParams := &pingpp.RedEnvelopeParams{
		App:         pingpp.App{Id: APP_ID},
		Channel:     "wx_pub",
		Order_no:    strconv.Itoa(orderno),
		Amount:      redEnvelopeModel.Amount,
		Currency:    "cny",
		Recipient:   WX_OPEN_ID,
		Subject:     redEnvelopeModel.Subject,
		Body:        redEnvelopeModel.Body,
		Description: redEnvelopeModel.Description,
		Extra:       extra,
	}
	redEnvelope, err := redEnvelope.New(redenvelopeParams)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	logger.Info(redEnvelope)
	redString, _ := json.Marshal(redEnvelope)
	return redString, nil
}

func Webhook(w http.ResponseWriter, r *http.Request) {
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
		publicKey, err := ioutil.ReadFile("pingpp_rsa_public_key.pem")
		if err != nil {
			fmt.Errorf("read failure: %v", err)
		}
		//base64解码再验证
		decodeStr, _ := base64.StdEncoding.DecodeString(signed)
		errs := pingpp.Verify([]byte(data), publicKey, decodeStr)
		if errs != nil {
			fmt.Println(errs)
			return
		} else {
			fmt.Println("success")
		}

		webhook, err := pingpp.ParseWebhooks(buf.Bytes())
		//fmt.Println(webhook.Type)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "fail")
			return
		}

		if webhook.Type == "charge.succeeded" {
			// TODO your code for charge
			w.WriteHeader(http.StatusOK)
		} else if webhook.Type == "refund.succeeded" {
			// TODO your code for refund
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}
