package pay_service

import (
	"crypto/md5"
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"encoding/xml"
	"fmt"
	var_const "github.com/EDDYCJY/go-gin-example/const"
	"github.com/EDDYCJY/go-gin-example/models"
	"github.com/EDDYCJY/go-gin-example/pkg/gredis"
	"github.com/EDDYCJY/go-gin-example/pkg/logging"
	"github.com/EDDYCJY/go-gin-example/pkg/setting"
	"github.com/EDDYCJY/go-gin-example/pkg/util"
	"github.com/EDDYCJY/go-gin-example/service/auth_service"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

func GetRedisKeyOrder(orderId int) string {
	return "game_pay_order:" + strconv.Itoa(orderId)
}

func getRedisKeyPayOrderIncr() string {
	return "max_pay_order"
}

func IncrOrderId() (int, error) {
	return gredis.Incr(getRedisKeyPayOrderIncr())
}

//支付订单号
func GeneratePayOrderId() string {
	t := time.Now().UnixNano()
	r := util.RandomStringNoUp(8)
	newId := "pay" + strconv.FormatInt(t, 16) + r
	return newId
}

//支付随机数
func GenerateNonceStr() string {
	t := time.Now().Format("20060102150405")
	r := util.RandomStringNoUp(4)
	newId := t + r
	return newId
}

func CreateOrder(form *models.Pay) bool {

	orderId, err := IncrOrderId()
	if err != nil {
		return false
	}
	form.OrderId = orderId
	form.NickName = auth_service.GetUserParamString(form.UserId, "nick_name")
	form.RegTime = int(time.Now().Unix())
	if !form.Save() {
		log, _ := json.Marshal(form)
		logging.Error("CreatePayOrder:form-" + string(log))
		return false
	}
	return true
}

func ExistOrder(orderId int) bool {
	if orderId <= 0 {
		return false
	}
	key := GetRedisKeyOrder(orderId)
	if gredis.Exists(key) {
		return true
	}
	//从数据库获取
	//先加锁
	lock, ok, _ := gredis.TryLock(key, setting.RedisSetting.LockTimeout)
	if lock == nil {
		return false
	}
	if !ok {
		i := 50
		for {
			if i--; i <= 0 {
				break
			}
			ok, _ := gredis.SetNX(lock.GetKey(), lock.GenerateToken(), lock.Timeout)
			if !ok {
				logging.Error("ExistPayOrder:TryLock-failed-tryAgain")
				ttl, err := gredis.GetTTL(lock.GetKey())
				if err != nil {
					logging.Error("ExistPayOrder:TryLock-GetTTL-failed")
				}
				time.Sleep(time.Duration(ttl/2) * time.Second)
			} else {
				break
			}
		}
	}

	info := models.Pay{
		OrderId: orderId,
	}
	dbRes, err := info.First()
	if dbRes == 0 {
		lock.UnLock()
		return false
	} else if dbRes == -1 {
		lock.UnLock()
		return false
	}
	//将数据放入redis
	jData, _ := json.Marshal(info)
	m := make(map[string]interface{})
	if err := json.Unmarshal(jData, &m); err != nil {
		lock.UnLock()
		return false
	}
	_, err = gredis.HMSet(key, m)
	if err != nil {
		lock.UnLock()
		return false
	}
	lock.UnLock()
	return true
}

func GetOrderParam(orderId int, param string) int {
	if !ExistOrder(orderId) {
		return 0
	}
	strParam, err := gredis.HGet(GetRedisKeyOrder(orderId), param)
	if err != nil {
		logging.Error("GetPayOrderParam:" + strconv.Itoa(orderId))
		return 0
	}
	if strParam == "" {
		strParam = "0"
	}
	p, err := strconv.Atoi(strParam)
	if err != nil {
		logging.Error("GetPayOrderParam:" + strParam)
		return 0
	}
	return p
}
func GetOrderParamString(orderId int, param string) string {
	if !ExistOrder(orderId) {
		return ""
	}
	strParam, err := gredis.HGet(GetRedisKeyOrder(orderId), param)
	if err != nil {
		logging.Error("GetPayOrderParamString:" + strconv.Itoa(orderId))
		return ""
	}
	return strParam
}

type PayOrderReq struct {
	AppId          string `xml:"appid"`
	Body           string `xml:"body"`
	MchId          string `xml:"mch_id"`
	NonceStr       string `xml:"nonce_str"`
	NotifyUrl      string `xml:"notify_url"`
	OpenId         string `xml:"openid"`
	TradeType      string `xml:"trade_type"`
	SpbillCreateIp string `xml:"spbill_create_ip"`
	TotalFee       int    `xml:"total_fee"`
	OutTradeNo     string `xml:"out_trade_no"`
	Sign           string `xml:"sign"`
}

//响应信息
type WXPayResp struct {
	ReturnCode string `xml:"return_code"`
	ReturnMsg  string `xml:"return_msg"`
	AppId      string `xml:"appid"`
	MchId      string `xml:"mch_id"`
	NonceStr   string `xml:"nonce_str"`
	Sign       string `xml:"sign"`
	ResultCode string `xml:"result_code"`
	ErrCode    string `xml:"err_code"`
	PrepayId   string `xml:"prepay_id"`
	TradeType  string `xml:"trade_type"`
	CodeUrl    string `xml:"code_url"`
}

func Pay(userId int, orderId int, ip string) (map[string]interface{}, bool) {
	//totalFee, err := GetOrderPrice(orderId)
	totalFee := GetOrderParam(orderId, "pay_amount")
	if totalFee == 0 {
		return nil, false
	}

	openId, err := auth_service.GetUserOpenId(userId)
	if err != nil || openId == "" {
		return nil, false
	}

	payOrderId := GeneratePayOrderId()
	desc := "费用说明"
	tradeType := "JSAPI"

	var payReq PayOrderReq
	payReq.AppId = var_const.WXAppID //微信开放平台我们创建出来的app的app id
	payReq.Body = desc
	payReq.MchId = var_const.WXMchID
	payReq.NonceStr = GenerateNonceStr()
	payReq.NotifyUrl = "https://www.bafangwangluo.com/pay/notify"
	payReq.OpenId = openId
	payReq.TradeType = tradeType
	payReq.SpbillCreateIp = ip
	payReq.TotalFee = totalFee
	payReq.OutTradeNo = payOrderId

	var reqMap = make(map[string]interface{}, 0)
	reqMap["appid"] = payReq.AppId                      //微信小程序appid
	reqMap["body"] = payReq.Body                        //商品描述
	reqMap["mch_id"] = payReq.MchId                     //商户号
	reqMap["nonce_str"] = payReq.NonceStr               //随机数
	reqMap["notify_url"] = payReq.NotifyUrl             //通知地址
	reqMap["out_trade_no"] = payReq.OutTradeNo          //订单号
	reqMap["openid"] = payReq.OpenId                    //openid
	reqMap["spbill_create_ip"] = payReq.SpbillCreateIp  //用户端ip   //订单生成的机器 IP
	reqMap["total_fee"] = strconv.Itoa(payReq.TotalFee) //订单总金额，单位为分
	reqMap["trade_type"] = payReq.TradeType             //trade_type=JSAPI时（即公众号支付），此参数必传，此参数为微信用户在商户对应appid下的唯一标识
	payReq.Sign = WxPayCalcSign(reqMap, var_const.WXMchKey)

	// 调用支付统一下单API
	bytesReq, err := xml.Marshal(payReq)
	if err != nil {
		return nil, false
	}
	strReq := string(bytesReq)
	//wxpay的unifiedorder接口需要http body中xmldoc的根节点是<xml></xml>这种，所以这里需要replace一下
	strReq = strings.Replace(strReq, "PayOrderReq", "xml", -1)
	bytesReq = []byte(strReq)

	req, err2 := http.NewRequest("POST", "https://api.mch.weixin.qq.com/pay/unifiedorder", strings.NewReader(string(bytesReq)))
	if err2 != nil {
		return nil, false
	}
	req.Header.Set("Content-Type", "text/xml;charset=utf-8")
	client := &http.Client{}
	resp, _ := client.Do(req)
	defer resp.Body.Close()

	body2, err3 := ioutil.ReadAll(resp.Body)
	if err3 != nil {
		return nil, false
	}
	var resp1 WXPayResp
	err = xml.Unmarshal(body2, &resp1)
	if err != nil {
		return nil, false
	}

	// 返回预付单信息
	if strings.ToUpper(resp1.ReturnCode) == "SUCCESS" && strings.ToUpper(resp1.ResultCode) == "SUCCESS" {
		// 再次签名
		var resMap = make(map[string]interface{}, 0)
		resMap["appId"] = resp1.AppId
		resMap["nonceStr"] = resp1.NonceStr                            //商品描述
		resMap["package"] = "prepay_id=" + resp1.PrepayId              //商户号
		resMap["signType"] = "MD5"                                     //签名类型
		resMap["timeStamp"] = strconv.FormatInt(time.Now().Unix(), 10) //当前时间戳

		resMap["paySign"] = WxPayCalcSign(resMap, var_const.WXMchKey)
		//保存支付订单 TODO
		dbInfo := models.Pay{
			OrderId: orderId,
		}
		var m = make(map[string]interface{})
		m["trade_no"] = payOrderId
		m["pay_desc"] = desc
		m["pay_ip"] = ip
		m["trade_type"] = tradeType
		if !dbInfo.Updates(m) {
			log, _ := json.Marshal(m)
			logging.Error("Pay:failed-" + string(log))
			return nil, false
		}
		return resMap, true
	}
	return nil, false
}

type RefundReq struct {
	AppId       string `xml:"appid"`
	MchId       string `xml:"mch_id"`
	NonceStr    string `xml:"nonce_str"`
	Sign        string `xml:"sign"`
	OutTradeNo  string `xml:"out_trade_no"`
	OutRefundNo string `xml:"out_refund_no"`
	TotalFee    int    `xml:"total_fee"`
	RefundFee   int    `xml:"refund_fee"`
	NotifyUrl   string `xml:"notify_url"`
}
type RefundResp struct {
	ReturnCode string `xml:"return_code"`
	ReturnMsg  string `xml:"return_msg"`
	ResultCode string `xml:"result_code"`
	ErrCode    string `xml:"err_code"`

	AppId    string `xml:"appid"`
	MchId    string `xml:"mch_id"`
	SubMchId string `xml:"sub_mch_id"`
	NonceStr string `xml:"nonce_str"`
	Sign     string `xml:"sign"`

	OutTradeNo  string `xml:"out_trade_no"`
	OutRefundNo string `xml:"out_refund_no"`
	RefundId    string `xml:"refund_id"`
	RefundFee   int    `xml:"refund_fee"`
	TotalFee    int    `xml:"total_fee"`
	CashFee     int    `xml:"cash_fee"`
}
type RefundNotify struct {
	Return_code    string `xml:"return_code"`
	Return_msg     string `xml:"return_msg"`
	Appid          string `xml:"appid"`
	Mch_id         string `xml:"mch_id"`
	Nonce          string `xml:"nonce_str"`
	Req_info       string `xml:"req_info"`
	Out_refund_no  string `xml:"out_refund_no"`
	Out_trade_no   string `xml:"out_trade_no"`
	Refund_fee     string `xml:"refund_fee"`
	Refund_status  string `xml:"refund_status"`
	Success_time   string `xml:"success_time"`
	Transaction_id string `xml:"transaction_id"`
}

//func Refund(orderId int) bool {
//	totalFee, err := GetOrderTakerPayAmount(orderId)
//	if err != nil || totalFee == 0 {
//		return false
//	}
//	outTradeNo, err := GetOrderTakerTradeNo(orderId)
//	if err != nil {
//		return false
//	}
//	payOrderId := GeneratePayOrderId()
//
//	var payReq RefundReq
//	payReq.AppId = var_const.WXAppID //微信开放平台我们创建出来的app的app id
//	payReq.MchId = var_const.WXMchID
//	payReq.NonceStr = GenerateNonceStr()
//	payReq.OutTradeNo = outTradeNo
//	payReq.OutRefundNo = payOrderId
//	payReq.TotalFee = totalFee
//	payReq.RefundFee = totalFee
//	payReq.NotifyUrl = "https://www.bafangwangluo.com/pay/taker/refundnotify"
//
//	var reqMap = make(map[string]interface{}, 0)
//	reqMap["appid"] = payReq.AppId        //微信小程序appid
//	reqMap["mch_id"] = payReq.MchId       //商户号
//	reqMap["nonce_str"] = payReq.NonceStr //随机数
//	reqMap["out_refund_no"] = payReq.OutRefundNo
//	reqMap["out_trade_no"] = payReq.OutTradeNo
//	reqMap["total_fee"] = payReq.TotalFee
//	reqMap["refund_fee"] = payReq.RefundFee
//	reqMap["notify_url"] = payReq.NotifyUrl
//	payReq.Sign = WxPayCalcSign(reqMap, var_const.WXMchKey)
//
//	// 调用支付统一下单API
//	bytesReq, err := xml.Marshal(payReq)
//	if err != nil {
//		return false
//	}
//	strReq := string(bytesReq)
//
//	strReq = strings.Replace(strReq, "RefundReq", "xml", -1)
//	bytesReq = []byte(strReq)
//
//	resp, err2 := KeyHttpsPost("https://api.mch.weixin.qq.com/secapi/pay/refund", "application/xml;charset=utf-8", strings.NewReader(string(bytesReq)))
//	if err2 != nil {
//		return false
//	}
//	//req.Header.Set("Content-Type", "text/xml;charset=utf-8")
//	//client := &http.Client{}
//	//resp, _ := client.Do(req)
//	defer resp.Body.Close()
//
//	body2, err3 := ioutil.ReadAll(resp.Body)
//	if err3 != nil {
//		return false
//	}
//	var resp1 RefundResp
//	err = xml.Unmarshal(body2, &resp1)
//	if err != nil {
//		return false
//	}
//	if resp1.ReturnCode == "SUCCESS" && resp1.ResultCode == "SUCCESS" && resp1.ReturnMsg == "OK" {
//		dbInfo := models.Order{
//			OrderId: orderId,
//		}
//		var m = make(map[string]interface{})
//		m["refund_trade_no"] = payReq.OutRefundNo
//		m["refund_amount"] = payReq.RefundFee
//		m["upd_time"] = int(time.Now().Unix())
//		log, _ := json.Marshal(m)
//		logging.Info("refund: begin order_id-" + strconv.Itoa(orderId) + "," + string(log))
//		if !dbInfo.Updates(m) {
//			logging.Error("refund: failed-db order_id-" + strconv.Itoa(orderId) + "," + string(log))
//			return false
//		}
//		db2Info := models.Order{
//			RefundTradeNo: payReq.OutRefundNo,
//		}
//		_, _ = db2Info.GetOrderInfoByRefundTradeNo()
//
//		return auth_service.RemoveUserMargin(db2Info.TakerUserId, db2Info.TakerPayAmount, "订单完成，取消保证金")
//	}
//	return true
//}

func KeyHttpsPost(url string, contentType string, body io.Reader) (*http.Response, error) {
	var wechatPayCert = "cert/apiclient_cert.pem"
	var wechatPayKey = "cert/apiclient_key.pem"
	var rootCa = "cert/cacert.pem"
	var tr *http.Transport
	// 微信提供的API证书,证书和证书密钥 .pem格式
	certs, err := tls.LoadX509KeyPair(wechatPayCert, wechatPayKey)
	if err != nil {
		log.Println("certs load err:", err)

	} else {
		// 微信支付HTTPS服务器证书的根证书  .pem格式
		rootCa, err := ioutil.ReadFile(rootCa)
		if err != nil {
			log.Println("err2222:", err)
		} else {
			pool := x509.NewCertPool()
			pool.AppendCertsFromPEM(rootCa)

			tr = &http.Transport{
				TLSClientConfig: &tls.Config{
					RootCAs:      pool,
					Certificates: []tls.Certificate{certs},
				},
			}
		}

	}
	client := &http.Client{Transport: tr}
	return client.Post(url, contentType, body)
}

//微信支付计算签名的函数
func WxPayCalcSign(mReq map[string]interface{}, key string) (sign string) {
	//STEP 1, 对key进行升序排序.
	sortedKeys := make([]string, 0)
	for k := range mReq {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Strings(sortedKeys)

	//STEP2, 对key=value的键值对用&连接起来，略过空值
	var signStrings string
	for _, k := range sortedKeys {
		value := fmt.Sprintf("%v", mReq[k])
		if value != "" {
			signStrings = signStrings + k + "=" + value + "&"
		}
	}

	//STEP3, 在键值对的最后加上key=API_KEY
	if key != "" {
		signStrings = signStrings + "key=" + key
	}
	//STEP4, 进行MD5签名并且将所有字符转为大写.
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(signStrings)) //
	cipherStr := md5Ctx.Sum(nil)
	upperSign := strings.ToUpper(hex.EncodeToString(cipherStr))
	return upperSign
}
