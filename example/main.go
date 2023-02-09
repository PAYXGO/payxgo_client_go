package main

import (
	"fmt"
	"log"
	"time"

	payxgo_client "github.com/PAYXGO/payxgo_client_go"
)

func main() {
	// 获取支付链接和支付信息
	config := payxgo_client.New("https://www.saleoner.com",
		"MEYwDQYJKoZIhvcNAQEBBQADNQAwMgIrAbQndM7qH5JSuXoXzV+pdmEy19002H6g5bW9F6jpf9NXleAy1eYTX+/hwwIDAQAB",
		"000fe36e62ff81b8f69bbcecdc154f539a55d32207e11fe5403669846604428d1b1988d4e918e585d896ed",
		payxgo_client.Pay, &payxgo_client.Config{
			Currency: "USD",
			Amount:   0.02,
			Vendor:   "alipay",
			OrderNum: "altsjdglkiuaiia5",
			IpnUrl:   "https://www.saleoner.com/go/v1/ipnCallback",
		})
	// 返回的cookie在刷新的时候需要携带
	result, cookie, err := config.PayAction()
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(result)
	fmt.Println(cookie)

	// 刷新支付链接  刷新请求1分钟一次
	conf := payxgo_client.New("https://www.saleoner.com",
		"MEYwDQYJKoZIhvcNAQEBBQADNQAwMgIrAbQndM7qH5JSuXoXzV+pdmEy19002H6g5bW9F6jpf9NXleAy1eYTX+/hwwIDAQAB",
		"000fe36e62ff81b8f69bbcecdc154f539a55d32207e11fe5403669846604428d1b1988d4e918e585d896ed",
		payxgo_client.Update, nil, cookie)
	result, cookie, err = conf.PayAction()
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(result)
	fmt.Println(cookie)
	time.Sleep(5 * time.Second)
}
