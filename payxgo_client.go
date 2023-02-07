package payxgo_client

import (
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"payxgo_client_go/payxgo_util"
	"reflect"
	"strings"
	"time"
)

type Config struct {
	// 币种
	Currency string `json:"currency"` // 币种
	// 支付金额
	Amount float64 `json:"amount"` // 支付金额
	// 支付通道
	Vendor string `json:"vendor"` // 支付通道
	// 订单号
	OrderNum string `json:"orderNum"` // 订单号
	// 异步回调
	IpnUrl string `json:"ipnUrl"` // 支付结果回调, 异步通知的url
}

type action int

const (
	Update action = iota + 1
	Pay
)

type payxgoClient struct {
	c          *Config
	requestId  string
	t          int64
	randomStr  string
	accessAddr string
	secretKey  string
	action     action
	cookie     string
	accessKey  string
}

// 刷新时获取二维码时接收到的cookie value原样传入
func New(urlPath, secretKey, accessKey string, action action, c *Config, cookieValue ...string) *payxgoClient {
	var cookie = ""
	if len(cookieValue) > 0 {
		cookie = cookieValue[0]
	}
	return &payxgoClient{
		c:          c,
		accessAddr: urlPath,
		secretKey:  secretKey,
		action:     action,
		cookie:     cookie,
		accessKey:  accessKey,
	}
}

func (n *payxgoClient) checkParams() bool {
	v := reflect.ValueOf(n.c)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	t := v.Type()
	if t.Kind() != reflect.Struct {
		log.Println(payxgo_util.TypeError.Error())
		return false
	}

	for i := 0; i < t.NumField(); i++ {
		fmt.Println(t.Field(i).Name, v.Field(i).String())
		if v.Field(i).IsZero() {
			log.Println(payxgo_util.NewError(1020, fmt.Sprintf("%s字段未初始化", t.Field(i).Name)))
			return false
		}
	}
	return true
}

// 设置请求头
func (n *payxgoClient) setHeader(req *http.Request) {
	if req == nil {
		return
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Request-ID", base64.StdEncoding.EncodeToString([]byte(n.requestId)))
	req.Header.Set("t", fmt.Sprint(n.t))
	if n.action == Update {
		req.Header.Set("Cookie", n.cookie)
	}
}

// 签名
func (n *payxgoClient) sign(params map[string]interface{}) string {
	if params == nil {
		params = make(map[string]interface{})
	}
	n.t = time.Now().Unix()
	params["t"] = n.t
	params["accessKey"] = n.accessKey
	n.randomStr = payxgo_util.Xid()
	fmt.Println("before:", n.randomStr)
	n.requestId = payxgo_util.Sign(n.randomStr, params)
	n.randomStr = base64.StdEncoding.EncodeToString(payxgo_util.RsaEncrypt([]byte(n.randomStr), []byte(n.secretKey)))
	params["randomStr"] = n.randomStr
	fmt.Println("after:", n.randomStr, n.requestId, n.t)
	delete(params, "t")
	buf, err := json.Marshal(params)
	if err != nil {
		log.Println(payxgo_util.NewError(2005, err.Error()).Error())
		return ""
	}
	return string(buf)
}

// 设置请求参数
func (n *payxgoClient) setParams() io.Reader {
	if n.action == Pay {
		if !n.checkParams() {
			return nil
		}
	}
	var r = make(map[string]interface{})
	reader, err := json.Marshal(n.c)
	if err != nil {
		return nil
	}

	err = json.Unmarshal(reader, &r)
	if err != nil {
		return nil
	}

	return strings.NewReader(n.sign(r))
}

// 操作处理
func (n *payxgoClient) PayAction() (content string, cookieValue string, err error) {
	return n.request()
}

// 请求api
func (n *payxgoClient) request() (string, string, error) {
	t := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: false,
		},
	}
	client := http.Client{
		Transport: t,
	}

	var cookieValue string

	req, err := http.NewRequest("POST", n.accessAddr, n.setParams())
	if err != nil {
		return "", cookieValue, err
	}

	n.setHeader(req)

	res, err := client.Do(req)
	if err != nil {
		return "", cookieValue, err
	}
	defer res.Body.Close()

	buf, err := io.ReadAll(res.Body)
	if err != nil {
		return "", cookieValue, err
	}
	cookieValue = res.Header.Get("Cookie")
	return string(buf), cookieValue, nil
}
