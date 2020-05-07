package order_service

import (
	"crypto/md5"
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
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

func GetRedisKeyOrder(orderId int) string {
	return "game_order:" + strconv.Itoa(orderId)
}

func getRedisKeyOrderIncr() string {
	return "max_order"
}

func IncrOrderId() (int, error) {
	return gredis.Incr(getRedisKeyOrderIncr())
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

func CreateOrder(form *models.Order) bool {
	//等级idx*1000 + star
	price := 0
	for i := 0; i < len(setting.PlatFormLevelAll); i++ {
		if setting.PlatFormLevelAll[i].Idx > form.CurLevel && setting.PlatFormLevelAll[i].Idx <= form.TargetLevel {
			price += setting.PlatFormLevelAll[i].Price
		}
	}
	//
	orderId, err := IncrOrderId()
	if err != nil {
		return false
	}
	form.OrderId = orderId
	nickName, _ := auth_service.GetUserNickName(form.UserId)
	form.NickName = nickName
	form.Price = price
	form.Status = var_const.OrderStatusAddOrder
	form.RegTime = int(time.Now().Unix())
	if !form.Save() {
		log, _ := json.Marshal(form)
		logging.Error("CreateOrder:form-" + string(log))
		return false
	}
	return true
}

func ExistOrder(orderId int) bool {
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
				logging.Error("ExistOrder:TryLock-failed-tryAgain")
				ttl, err := gredis.GetTTL(lock.GetKey())
				if err != nil {
					logging.Error("ExistOrder:TryLock-GetTTL-failed")
				}
				time.Sleep(time.Duration(ttl/2) * time.Second)
			} else {
				break
			}
		}
	}

	info := models.Order{
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

func GetOrderPrice(orderId int) (int, error) {
	if !ExistOrder(orderId) {
		return 0, fmt.Errorf("GetOrderPrice:OrderIdnoExist")
	}
	strPrice, err := gredis.HGet(GetRedisKeyOrder(orderId), "price")
	if err != nil {
		logging.Error("GetOrderPrice:" + strconv.Itoa(orderId))
		return 0, err
	}
	price, err := strconv.Atoi(strPrice)
	if err != nil {
		logging.Error("GetOrderPrice:" + strPrice)
		return 0, err
	}
	return price, nil
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

func Pay(userId int, orderId int, ip string) bool {
	totalFee, err := GetOrderPrice(orderId)
	if err != nil || totalFee == 0 {
		return false
	}

	openId, err := auth_service.GetUserOpenId(userId)
	if err != nil || openId == "" {
		return false
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
		return false
	}
	strReq := string(bytesReq)
	//wxpay的unifiedorder接口需要http body中xmldoc的根节点是<xml></xml>这种，所以这里需要replace一下
	strReq = strings.Replace(strReq, "PayOrderReq", "xml", -1)
	bytesReq = []byte(strReq)

	req, err2 := http.NewRequest("POST", "https://api.mch.weixin.qq.com/pay/unifiedorder", strings.NewReader(string(bytesReq)))
	if err2 != nil {
		return false
	}
	req.Header.Set("Content-Type", "text/xml;charset=utf-8")
	client := &http.Client{}
	resp, _ := client.Do(req)
	defer resp.Body.Close()

	body2, err3 := ioutil.ReadAll(resp.Body)
	if err3 != nil {
		return false
	}
	var resp1 WXPayResp
	err = xml.Unmarshal(body2, &resp1)
	if err != nil {
		return false
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
		dbInfo := models.Order{
			OrderId: orderId,
		}
		var m = make(map[string]interface{})
		m["trade_no"] = payOrderId
		m["status"] = var_const.OrderStatusWaitPay
		m["pay_amount"] = totalFee
		m["pay_desc"] = desc
		m["pay_ip"] = ip
		m["trade_type"] = tradeType
		m["upd_time"] = int(time.Now().Unix())
		if !dbInfo.Updates(m) {
			log, _ := json.Marshal(m)
			logging.Error("Pay:failed-" + string(log))
			return false
		}
		return true
	}
	return false
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
