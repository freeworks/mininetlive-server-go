package pay

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	pingpp "github.com/pingplusplus/pingpp-go/pingpp"
	"github.com/pingplusplus/pingpp-go/pingpp/charge"
	"github.com/pingplusplus/pingpp-go/pingpp/redEnvelope"
)

const (
	API_KEY    string = `sk_test_ibbTe5jLGCi5rzfH4OqPW9KC`
	WX_OPEN_ID string = ``
)

func init() {
	pingpp.LogLevel = 2
	pingpp.Key = API_KEY
	fmt.Println("Go SDK Version:", pingpp.Version())
	pingpp.AcceptLanguage = "zh-CN"
	//设置商户的私钥 记得在Ping++上配置公钥
	//pingpp.AccountPrivateKey
}

func GetCharge(w http.ResponseWriter, r *http.Request) {
	if strings.ToUpper(r.Method) == "GET" {
		http.NotFound(w, r)
	} else if strings.ToUpper(r.Method) == "POST" {
		var chargeParams pingpp.ChargeParams
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")
		defer r.Body.Close()
		buf := new(bytes.Buffer)
		buf.ReadFrom(r.Body)

		json.Unmarshal(buf.Bytes(), &chargeParams)
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		orderno := r.Intn(999999999999999) //TODO 活动号 活动id_随机数
		extra := make(map[string]interface{})
		switch strings.ToLower(chargeParams.Channel) {
		case "alipay_wap":
			extra["cancel_url"] = "http://www.yourdomain.com/cancel"
			extra["success_url"] = "http://www.yourdomain.com/success"
		case "wx_pub":
			extra["open_id"] = WX_OPEN_ID
		case "wx_pub_qr":
			extra["product_id"] = "your_productid"
		}
		params := &pingpp.ChargeParams{
			Order_no:  strconv.Itoa(orderno),
			App:       pingpp.App{Id: "app_1Gqj58ynP0mHeX1q"},
			Amount:    chargeParams.Amount,
			Channel:   strings.ToLower(chargeParams.Channel),
			Currency:  "cny",
			Client_ip: "127.0.0.1",
			Subject:   "Your Subject",
			Body:      "Your Body",
			Extra:     extra,
		}
		//返回的第一个参数是 charge 对象，你需要将其转换成 json 给客户端，或者客户端接收后转换。
		ch, err := charge.New(params)
		if err != nil {
			errs, _ := json.Marshal(err)
			fmt.Fprint(w, string(errs))
		} else {
			chs, _ := json.Marshal(ch)
			fmt.Fprintln(w, string(chs))
		}
	}
}

func NewRedEnvelope() {
	extra := make(map[string]interface{})
	extra["nick_name"] = "Nick Name"
	extra["send_name"] = "Send Name"
	//这里是随便设置的随机数作为订单号，仅作示例，该方法可能产生相同订单号，商户请自行生成，不要纠结该方法。
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	orderno := r.Intn(999999999999999)

	redenvelopeParams := &pingpp.RedEnvelopeParams{
		App:         pingpp.App{Id: "app_1Gqj58ynP0mHeX1q"},
		Channel:     "wx_pub",
		Order_no:    strconv.Itoa(orderno),
		Amount:      100,
		Currency:    "cny",
		Recipient:   "youropenid",
		Subject:     "Your Subject",
		Body:        "Your Body",
		Description: "Your Description",
		Extra:       extra,
	}
	redEnvelope, err := redEnvelope.New(redenvelopeParams)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(redEnvelope)
	redString, _ := json.Marshal(redEnvelope)
	fmt.Println(string(redString))
}

func GetRedEnvelope() {
	redEnvelope, err := redEnvelope.Get("red_id")
	if err != nil {
		log.Fatal(err)
	}
	restring, _ := json.Marshal(redEnvelope)
	log.Printf("%v\n", string(restring))
}

func ListRedEnvelope() {
	params := &pingpp.RedEnvelopeListParams{}
	params.Filters.AddFilter("limit", "", "2")
	//设置是不是只需要之前设置的 limit 这一个查询参数
	params.Single = true
	i := redEnvelope.List(params)
	for i.Next() {
		c := i.RedEnvelope()
		ch, _ := json.Marshal(c)
		fmt.Println(string(ch))
	}
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
