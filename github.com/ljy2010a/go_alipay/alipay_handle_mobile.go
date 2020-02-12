package alipay

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"time"
)

/**
移动端rsa 处理
*/
func AlipayAppRsaSign(amr AlipayMobileRequest, alipayConfig *AlipayConfig) string {
	p := Kvpairs{
		Kvpair{`app_id`, alipayConfig.AppId},
		Kvpair{`method`, "alipay.trade.app.pay"},
		Kvpair{`charset`, alipayConfig.Input_charset},
		Kvpair{`sign_type`, alipayConfig.Sign_type},
		Kvpair{`timestamp`, time.Now().Format("2006-01-02 15:04:05")},
		Kvpair{`version`, "1.0"},
		Kvpair{`notify_url`, alipayConfig.Notify_url},
	}
	bizMap := map[string]string{
		`body`:         amr.Body,
		`subject`:      amr.Subject,
		`out_trade_no`: amr.OutTradeNo,
		`total_amount`: fmt.Sprintf("%.2f", amr.TotalFee),
		`seller_id`:    amr.Seller_id,
		`product_code`: "QUICK_MSECURITY_PAY",
	}
	if len(amr.AppAuthToken) > 0 {
		bizMap[`auth_token`] = amr.AppAuthToken
	}
	if len(amr.EnablePayChannels) > 0 {
		bizMap[`enable_pay_channels`] = amr.EnablePayChannels
	}
	bizContent, _ := json.Marshal(bizMap)
	p = append(p, Kvpair{"biz_content", string(bizContent)})
	RsaSignApp(&p, alipayConfig)
	return createLinkstringAllUrlencode(&p)
}

/**
移动端rsa 处理
*/
func AlipayMobileRsaSign(amr AlipayMobileRequest, alipayConfig *AlipayConfig) string {
	p := Kvpairs{
		Kvpair{`_input_charset`, alipayConfig.Input_charset},
		Kvpair{`partner`, alipayConfig.Partner},
		Kvpair{`payment_type`, alipayConfig.Payment_type},
		Kvpair{`notify_url`, alipayConfig.Notify_url},
		Kvpair{`service`, alipayConfig.Service},
		Kvpair{`seller_id`, alipayConfig.Seller_id},
		Kvpair{`out_trade_no`, amr.OutTradeNo},
		Kvpair{`subject`, amr.Subject},
		Kvpair{`total_fee`, fmt.Sprintf("%.2f", amr.TotalFee)},
		Kvpair{`body`, amr.Body},
		// Kvpair{`it_b_pay`, `15m`},
	}
	RsaSign(&p, alipayConfig)
	return createLinkstringUrlencode(&p)

}

/**
 * 针对notify_url验证消息是否是支付宝发出的合法消息
 * @return 验证结果
 */
func VerifyMobileNotify(r *http.Request, alipayConfig *AlipayConfig) error {

	log.Println("VerifyMobileNotify begin")

	signErr := verifySign(r.PostForm, alipayConfig)
	if signErr != nil {
		return signErr
	}
	log.Println("VerifyMobileNotify verifySign success")

	notify_id := r.FormValue("notify_id")
	//获取支付宝远程服务器ATN结果（验证是否是支付宝发来的消息）
	responseTxt, err := getResponse(notify_id, alipayConfig)
	if err != nil {
		return err
	}
	log.Println("VerifyMobileNotify responseTxt is: %v", responseTxt)

	reg := regexp.MustCompile(`true`)
	if 0 == len(reg.FindAllString(responseTxt, -1)) {
		log.Println("responseTxt verify fail ")
		return fmt.Errorf("responseTxt is wrong")
	}
	log.Println("VerifyMobileNotify responseTxt verify success ")
	return nil
}

func VerifySign(r *http.Request, alipayConfig *AlipayConfig) error {
	return verifySign(r.PostForm, alipayConfig)
}

func RsaSignApp(para *Kvpairs, alipayConfig *AlipayConfig) string {
	buildRequestParaForPost(para, alipayConfig)
	// buildRequestMysign(para, config)
	return createLinkstringForPost(para)
}

func RsaSign(para *Kvpairs, alipayConfig *AlipayConfig) string {
	buildRequestPara(para, alipayConfig)
	// buildRequestMysign(para, config)
	return createLinkstringUrlencode(para)
}
