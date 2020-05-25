package v1

import (
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	var_const "github.com/EDDYCJY/go-gin-example/const"
	"github.com/EDDYCJY/go-gin-example/models"
	"github.com/EDDYCJY/go-gin-example/pkg/app"
	"github.com/EDDYCJY/go-gin-example/pkg/e"
	"github.com/EDDYCJY/go-gin-example/pkg/logging"
	"github.com/EDDYCJY/go-gin-example/pkg/util"
	"github.com/EDDYCJY/go-gin-example/service/auth_service"
	"github.com/EDDYCJY/go-gin-example/service/order_service"
	"github.com/EDDYCJY/go-gin-example/service/team_service"
	"github.com/gin-gonic/gin"
	"github.com/nanjishidu/gomini/gocrypto"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// @Summary 微信下单 获取发起微信支付所需的数据
// @Produce  json
// @Param user_id body int false "user_id"
// @Param order_id body int false "order_id"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/pay [post]
// @Tags 微信支付
func WxPay(c *gin.Context) {
	appG := app.Gin{C: c}

	userId, err := strconv.Atoi(c.PostForm("user_id"))
	if err != nil {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	if !auth_service.ExistUserInfo(userId) {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	orderId, err := strconv.Atoi(c.PostForm("order_id"))
	if err != nil {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	if !order_service.ExistOrder(orderId) {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	totalFee, err := order_service.GetOrderPrice(orderId)
	if err != nil || totalFee == 0 {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	d, ok := order_service.Pay(userId, orderId, c.ClientIP())
	if !ok {
		appG.Response(http.StatusBadRequest, e.ERROR, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, d)
	return
}

// 支付回调接口
func WxNotify(c *gin.Context) {
	//logging.Info("--------pay")
	var resMap = make(map[string]interface{}, 0)
	resMap["return_code"] = "SUCCESS"
	resMap["return_msg"] = "OK"

	valueXml, _ := ioutil.ReadAll(c.Request.Body) //获取post的数据
	//logging.Info("--------pay:" + string(valueXml))
	values := util.Xml2Map(string(valueXml))

	if retCode, ok := values["result_code"]; retCode != "SUCCESS" || !ok {
		resMap["return_code"] = "FAIL"
		resMap["return_msg"] = "result_code错误"
		resStr := util.Map2Xml(resMap)
		c.JSON(http.StatusOK, resStr)
		return
	}
	if _, ok := values["out_trade_no"]; !ok {
		resMap["return_code"] = "FAIL"
		resMap["return_msg"] = "out_trade_no错误"
		resStr := util.Map2Xml(resMap)
		c.JSON(http.StatusOK, resStr)
		return
	}
	if _, ok := values["sign"]; !ok {
		resMap["return_code"] = "FAIL"
		resMap["return_msg"] = "sign错误"
		resStr := util.Map2Xml(resMap)
		c.JSON(http.StatusOK, resStr)
		return
	}

	//微信提交过来的签名
	postSign := values["sign"]
	delete(values, "sign")
	//根据提交过来的值，和我的商户支付秘钥，生成的签名
	userSign := order_service.WxPayCalcSign(values, var_const.WXMchKey)
	//验证提交过来的签名是否正确
	if userSign != postSign {
		resMap["return_code"] = "FAIL"
		resMap["return_msg"] = "sign错误"
		resStr := util.Map2Xml(resMap)
		c.JSON(http.StatusOK, resStr)
		return
	}

	//判断订单存在
	payOrderId := values["out_trade_no"].(string)
	info := models.Order{
		TradeNo: payOrderId,
	}
	_, err := info.GetOrderInfoByTradeNo()
	if err != nil {
		resMap["return_code"] = "FAIL"
		resMap["return_msg"] = "out_trade_no错误"
		resStr := util.Map2Xml(resMap)
		c.JSON(http.StatusOK, resStr)
		return
	}
	//logging.Info("--------info.Status")
	if info.Status != var_const.OrderStatusWaitPay {
		//logging.Info("--------info.!Status")
		resMap["return_code"] = "FAIL"
		resMap["return_msg"] = "out_trade_no错误"
		resStr := util.Map2Xml(resMap)
		c.JSON(http.StatusOK, resStr)
		return
	}

	//
	//order信息更新
	//保存支付订单 TODO
	dbInfo := models.Order{
		OrderId: info.OrderId,
	}
	var m = make(map[string]interface{})
	m["status"] = var_const.OrderStatusPaidPay
	m["transaction_id"] = resMap["transaction_id"]
	m["pay_time"] = int(time.Now().Unix())
	if !dbInfo.Updates(m) {
		log, _ := json.Marshal(m)
		logging.Error("WxNotify:db-failed-" + string(log))
		resMap["return_code"] = "FAIL"
		resMap["return_msg"] = "out_trade_no错误"
		resStr := util.Map2Xml(resMap)
		c.JSON(http.StatusOK, resStr)
		return
	}
	//判断是否是车队订单
	//teamId, _ := order_service.GetOrderTeamId(info.OrderId)
	//if teamId > 0 {
	//	team_service.SetTeamOrderFinished(teamId, info.OrderId)
	//
	//	uId, _ := order_service.GetOrderUserId(info.OrderId)
	//	var m2 = make(map[string]interface{})
	//	team := models.Team{
	//		TeamId: teamId,
	//	}
	//	if team_service.GetTeamParam(teamId, "owner_type") == var_const.UserTypeNormal && team_service.GetTeamParam(teamId, "owner_id") == uId {
	//		m2["status"] = var_const.TeamCanShow
	//		if uId == team_service.GetTeamParam(teamId, "user_id1") {
	//			m2["nick_name1"], _ = auth_service.GetUserNickName(uId)
	//			m2["order_status1"] = 1
	//			m2["user1_pay_time"] = int(time.Now().Unix())
	//		}
	//	} else {
	//		if uId == team_service.GetTeamParam(teamId, "user_id1") {
	//			m2["nick_name1"], _ = auth_service.GetUserNickName(uId)
	//			m2["order_status1"] = 1
	//			m2["user1_pay_time"] = int(time.Now().Unix())
	//		}
	//		if uId == team_service.GetTeamParam(teamId, "user_id2") {
	//			m2["nick_name2"], _ = auth_service.GetUserNickName(uId)
	//			m2["order_status2"] = 1
	//			m2["user2_pay_time"] = int(time.Now().Unix())
	//		}
	//	}
	//	team.Updates(m2)
	//}

	//logging.Info("--------Updates")
	resMap["return_code"] = "SUCCESS"
	resMap["return_msg"] = "OK"
	//logging.Info("WxNotify:SUCCESS-")
	resStr := util.Map2Xml(resMap)
	logging.Info("WxNotify:SUCCESS-" + resStr)
	c.JSON(http.StatusOK, resStr)
	return
}

// @Summary 代练交平台押金
// @Produce  json
// @Param user_id body int false "user_id"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/pay/deposit [post]
// @Tags 微信支付
func DepositWxPay(c *gin.Context) {
	appG := app.Gin{C: c}

	userId, err := strconv.Atoi(c.PostForm("user_id"))
	if err != nil {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	if !auth_service.ExistUserInfo(userId) || auth_service.IsUserTypeInstead(userId) {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	if auth_service.GetUserParam(userId, "deposit") > 0 {
		appG.Response(http.StatusBadRequest, e.NO_DEPOSIT, nil)
		return
	}

	d, ok := order_service.DepositPay(userId, c.ClientIP())
	if !ok {
		appG.Response(http.StatusBadRequest, e.ERROR, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, d)
	return
}

// 支付回调接口
func DepositWxNotify(c *gin.Context) {
	//logging.Info("--------pay")
	var resMap = make(map[string]interface{}, 0)
	resMap["return_code"] = "SUCCESS"
	resMap["return_msg"] = "OK"

	valueXml, _ := ioutil.ReadAll(c.Request.Body) //获取post的数据
	//logging.Info("--------pay:" + string(valueXml))
	values := util.Xml2Map(string(valueXml))

	if retCode, ok := values["result_code"]; retCode != "SUCCESS" || !ok {
		resMap["return_code"] = "FAIL"
		resMap["return_msg"] = "result_code错误"
		resStr := util.Map2Xml(resMap)
		c.JSON(http.StatusOK, resStr)
		return
	}
	if _, ok := values["out_trade_no"]; !ok {
		resMap["return_code"] = "FAIL"
		resMap["return_msg"] = "out_trade_no错误"
		resStr := util.Map2Xml(resMap)
		c.JSON(http.StatusOK, resStr)
		return
	}
	if _, ok := values["sign"]; !ok {
		resMap["return_code"] = "FAIL"
		resMap["return_msg"] = "sign错误"
		resStr := util.Map2Xml(resMap)
		c.JSON(http.StatusOK, resStr)
		return
	}

	//微信提交过来的签名
	postSign := values["sign"]
	delete(values, "sign")
	//根据提交过来的值，和我的商户支付秘钥，生成的签名
	userSign := order_service.WxPayCalcSign(values, var_const.WXMchKey)
	//验证提交过来的签名是否正确
	if userSign != postSign {
		resMap["return_code"] = "FAIL"
		resMap["return_msg"] = "sign错误"
		//logging.Info("--------pay-userSign")
		resStr := util.Map2Xml(resMap)
		c.JSON(http.StatusOK, resStr)
		return
	}

	//判断订单存在
	payOrderId := values["out_trade_no"].(string)
	info := models.User{
		DepositTradeNo: payOrderId,
	}
	_, err := info.GetUserInfoByDepositTradeNo()
	if err != nil {
		resMap["return_code"] = "FAIL"
		resMap["return_msg"] = "out_trade_no错误"
		resStr := util.Map2Xml(resMap)
		c.JSON(http.StatusOK, resStr)
		return
	}
	if info.Deposit > 0 {
		resMap["return_code"] = "FAIL"
		resMap["return_msg"] = "out_trade_no错误"
		resStr := util.Map2Xml(resMap)
		c.JSON(http.StatusOK, resStr)
		return
	}

	//
	//order信息更新
	//保存支付订单 TODO
	dbInfo := models.User{
		UserId: info.UserId,
	}
	var m = make(map[string]interface{})
	m["deposit"] = var_const.Deposit
	m["deposit_time"] = int(time.Now().Unix())
	m["type"] = var_const.UserTypeInstead
	m["check_pass"] = var_const.CheckPass
	if !dbInfo.Updates(m) {
		log, _ := json.Marshal(m)
		logging.Error("DepositWxNotify:db-failed-" + string(log))
		resMap["return_code"] = "FAIL"
		resMap["return_msg"] = "out_trade_no错误"
		resStr := util.Map2Xml(resMap)
		c.JSON(http.StatusOK, resStr)
		return
	}

	resMap["return_code"] = "SUCCESS"
	resMap["return_msg"] = "OK"
	resStr := util.Map2Xml(resMap)
	c.JSON(http.StatusOK, resStr)
	return
}

// @Summary 微信接单 获取发起微信支付所需的数据
// @Produce  json
// @Param user_id body int false "user_id"
// @Param order_id body int false "order_id"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/pay/taker [post]
// @Tags 微信支付
func TakerWxPay(c *gin.Context) {
	appG := app.Gin{C: c}

	userId, err := strconv.Atoi(c.PostForm("user_id"))
	if err != nil {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	if !auth_service.ExistUserInfo(userId) || !auth_service.IsUserTypeInstead(userId) {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	if auth_service.GetUserParam(userId, "deposit") <= 0 {
		appG.Response(http.StatusBadRequest, e.NO_DEPOSIT, nil)
		return
	}
	orderId, err := strconv.Atoi(c.PostForm("order_id"))
	teamId, _ := order_service.GetOrderTeamId(orderId)
	if teamId > 0 {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	if err != nil {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	if !order_service.ExistOrder(orderId) {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	status := order_service.GetOrderParam(orderId, "status")
	if status != var_const.OrderStatusPaidPay {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	uId := order_service.GetOrderParam(orderId, "user_id")
	if uId == userId {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	d, ok := order_service.TakerPay(userId, orderId, c.ClientIP())
	if !ok {
		appG.Response(http.StatusBadRequest, e.ERROR, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, d)
	return
}

// 支付回调接口
func TakerWxNotify(c *gin.Context) {
	//logging.Info("--------pay")
	var resMap = make(map[string]interface{}, 0)
	resMap["return_code"] = "SUCCESS"
	resMap["return_msg"] = "OK"

	valueXml, _ := ioutil.ReadAll(c.Request.Body) //获取post的数据
	//logging.Info("--------pay:" + string(valueXml))
	values := util.Xml2Map(string(valueXml))

	if retCode, ok := values["result_code"]; retCode != "SUCCESS" || !ok {
		resMap["return_code"] = "FAIL"
		resMap["return_msg"] = "result_code错误"
		resStr := util.Map2Xml(resMap)
		c.JSON(http.StatusOK, resStr)
		return
	}
	if _, ok := values["out_trade_no"]; !ok {
		resMap["return_code"] = "FAIL"
		resMap["return_msg"] = "out_trade_no错误"
		resStr := util.Map2Xml(resMap)
		c.JSON(http.StatusOK, resStr)
		return
	}
	if _, ok := values["sign"]; !ok {
		resMap["return_code"] = "FAIL"
		resMap["return_msg"] = "sign错误"
		resStr := util.Map2Xml(resMap)
		c.JSON(http.StatusOK, resStr)
		return
	}

	//微信提交过来的签名
	postSign := values["sign"]
	delete(values, "sign")
	//根据提交过来的值，和我的商户支付秘钥，生成的签名
	userSign := order_service.WxPayCalcSign(values, var_const.WXMchKey)
	//验证提交过来的签名是否正确
	if userSign != postSign {
		resMap["return_code"] = "FAIL"
		resMap["return_msg"] = "sign错误"
		//logging.Info("--------pay-userSign")
		resStr := util.Map2Xml(resMap)
		c.JSON(http.StatusOK, resStr)
		return
	}

	//判断订单存在
	payOrderId := values["out_trade_no"].(string)
	info := models.Order{
		TakerTradeNo: payOrderId,
	}
	_, err := info.GetOrderInfoByTakerTradeNo()
	if err != nil {
		resMap["return_code"] = "FAIL"
		resMap["return_msg"] = "out_trade_no错误"
		resStr := util.Map2Xml(resMap)
		c.JSON(http.StatusOK, resStr)
		return
	}
	if info.Status != var_const.OrderStatusTakerWaitPay {
		resMap["return_code"] = "FAIL"
		resMap["return_msg"] = "out_trade_no错误"
		resStr := util.Map2Xml(resMap)
		c.JSON(http.StatusOK, resStr)
		return
	}

	//
	//order信息更新
	//保存支付订单 TODO
	dbInfo := models.Order{
		OrderId: info.OrderId,
	}
	var m = make(map[string]interface{})
	m["status"] = var_const.OrderStatusTakerPaid
	m["taker_transaction_id"] = resMap["transaction_id"]
	m["taker_time"] = int(time.Now().Unix())
	if !dbInfo.Updates(m) {
		log, _ := json.Marshal(m)
		logging.Error("WxNotify:db-failed-" + string(log))
		resMap["return_code"] = "FAIL"
		resMap["return_msg"] = "out_trade_no错误"
		resStr := util.Map2Xml(resMap)
		c.JSON(http.StatusOK, resStr)
		return
	}

	if !auth_service.AddUserMargin(info.TakerUserId, info.TakerPayAmount) {
		resMap["return_code"] = "FAIL"
		resMap["return_msg"] = "out_trade_no错误"
		resStr := util.Map2Xml(resMap)
		c.JSON(http.StatusOK, resStr)
		return
	}

	resMap["return_code"] = "SUCCESS"
	resMap["return_msg"] = "OK"
	resStr := util.Map2Xml(resMap)
	c.JSON(http.StatusOK, resStr)
	return
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

//微信退款回调
func WxRefundCallback(c *gin.Context) {
	body, _ := ioutil.ReadAll(c.Request.Body)
	var mr RefundNotify
	_ = xml.Unmarshal(body, &mr)
	key := "t7v5TMsxhW6VH2f231NaB1BGL33CRjt3"
	b, _ := base64.StdEncoding.DecodeString(mr.Req_info)
	gocrypto.SetAesKey(strings.ToLower(gocrypto.Md5(key)))
	plaintext, _ := gocrypto.AesECBDecrypt(b)

	var mnr RefundNotify
	_ = xml.Unmarshal(plaintext, &mnr)

	if mnr.Refund_status == "SUCCESS" {
		//f(mnr.Out_trade_no)
		dbInfo := models.Order{
			RefundTradeNo: mnr.Out_refund_no,
		}
		_, _ = dbInfo.GetOrderInfoByRefundTradeNo()
		if dbInfo.Status != var_const.OrderStatusConfirmFinished {
			c.JSON(http.StatusOK, nil)
			return
		}
		var m = make(map[string]interface{})
		m["status"] = var_const.OrderStatusRefundFinished
		m["upd_time"] = int(time.Now().Unix())
		if !dbInfo.Updates(m) {
			log, _ := json.Marshal(m)
			logging.Error("WxRefundCallback:failed-" + string(log))
			c.JSON(http.StatusOK, nil)
			return
		}
		logging.Info("refund finish RefundTradeNo:" + mnr.Out_refund_no)
		resStr := "<xml><return_code>SUCCESS</return_code><return_msg>OK</return_msg></xml>"

		c.JSON(http.StatusOK, resStr)
		return
	} else {
		c.JSON(http.StatusOK, nil)
		return
	}
}

// @Summary 微信下单 获取发起微信支付所需的数据
// @Produce  json
// @Param user_id body int false "user_id"
// @Param order_id body int false "team_id"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/teampay [post]
// @Tags 微信支付
func TeamWxPay(c *gin.Context) {
	appG := app.Gin{C: c}

	userId, err := strconv.Atoi(c.PostForm("user_id"))
	if err != nil {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	if !auth_service.ExistUserInfo(userId) {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	teamId, err := strconv.Atoi(c.PostForm("team_id"))
	if err != nil {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	if !team_service.ExistTeam(teamId) {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	staus := team_service.GetTeamParam(teamId, "status")
	if staus >= var_const.TeamWorking {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	takePayStaus := team_service.GetTeamParam(teamId, "taker_pay_status")
	if takePayStaus == 1 {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	//如果是接用户的车队
	if team_service.GetTeamParam(teamId, "owner_type") == var_const.UserTypeNormal {
		takerId := team_service.GetTeamParam(teamId, "taker_user_id")
		if takerId != 0 && takerId != userId {
			appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
			return
		}
	}
	if team_service.GetTeamParam(teamId, "owner_type") == var_const.UserTypeInstead {
		if team_service.GetTeamParam(teamId, "owner_id") != userId {
			appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
			return
		}
	}

	d, ok := team_service.Pay(userId, teamId, c.ClientIP())
	if !ok {
		appG.Response(http.StatusBadRequest, e.ERROR, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, d)
	return
}

func TeamWxNotify(c *gin.Context) {
	//logging.Info("--------pay")
	var resMap = make(map[string]interface{}, 0)
	resMap["return_code"] = "SUCCESS"
	resMap["return_msg"] = "OK"

	valueXml, _ := ioutil.ReadAll(c.Request.Body) //获取post的数据
	//logging.Info("--------pay:" + string(valueXml))
	values := util.Xml2Map(string(valueXml))

	if retCode, ok := values["result_code"]; retCode != "SUCCESS" || !ok {
		resMap["return_code"] = "FAIL"
		resMap["return_msg"] = "result_code错误"
		resStr := util.Map2Xml(resMap)
		c.JSON(http.StatusOK, resStr)
		return
	}
	if _, ok := values["out_trade_no"]; !ok {
		resMap["return_code"] = "FAIL"
		resMap["return_msg"] = "out_trade_no错误"
		resStr := util.Map2Xml(resMap)
		c.JSON(http.StatusOK, resStr)
		return
	}
	if _, ok := values["sign"]; !ok {
		resMap["return_code"] = "FAIL"
		resMap["return_msg"] = "sign错误"
		resStr := util.Map2Xml(resMap)
		c.JSON(http.StatusOK, resStr)
		return
	}

	//微信提交过来的签名
	postSign := values["sign"]
	delete(values, "sign")
	//根据提交过来的值，和我的商户支付秘钥，生成的签名
	userSign := order_service.WxPayCalcSign(values, var_const.WXMchKey)
	//验证提交过来的签名是否正确
	if userSign != postSign {
		resMap["return_code"] = "FAIL"
		resMap["return_msg"] = "sign错误"
		//logging.Info("--------pay-userSign")
		resStr := util.Map2Xml(resMap)
		c.JSON(http.StatusOK, resStr)
		return
	}

	//判断订单存在
	payOrderId := values["out_trade_no"].(string)
	info := models.Team{
		TakerTradeNo: payOrderId,
	}
	_, err := info.GetOrderInfoByTakerTradeNo()
	if err != nil {
		resMap["return_code"] = "FAIL"
		resMap["return_msg"] = "out_trade_no错误"
		resStr := util.Map2Xml(resMap)
		c.JSON(http.StatusOK, resStr)
		return
	}
	if info.TakerPayStatus != 0 {
		resMap["return_code"] = "FAIL"
		resMap["return_msg"] = "out_trade_no错误"
		resStr := util.Map2Xml(resMap)
		c.JSON(http.StatusOK, resStr)
		return
	}

	//
	//order信息更新
	//保存支付订单 TODO
	dbInfo := models.Team{
		TeamId: info.TeamId,
	}
	var m = make(map[string]interface{})
	m["taker_pay_status"] = 1
	if info.OwnerType == var_const.UserTypeNormal && info.OrderStatus1 == 1 {
		m["status"] = var_const.TeamWorking
		//设置订单为已接单 TODO
		dbInfo := models.Order{
			OrderId: team_service.GetTeamParam(info.TeamId, "order_id1"),
		}
		var m = make(map[string]interface{})
		m["status"] = var_const.OrderStatusTakerPaid
		m["taker_user_id"] = team_service.GetTeamParam(info.TeamId, "taker_user_id")
		m["taker_time"] = int(time.Now().Unix())
		if !dbInfo.Updates(m) {
			log, _ := json.Marshal(m)
			logging.Error("TeamWxNotify:db-failed-" + string(log))
			return
		}
	}
	if info.OwnerType == var_const.UserTypeInstead {
		m["status"] = var_const.TeamCanShow
	}
	m["taker_time"] = int(time.Now().Unix())
	if !dbInfo.Updates(m) {
		log, _ := json.Marshal(m)
		logging.Error("WxNotify:db-failed-" + string(log))
		resMap["return_code"] = "FAIL"
		resMap["return_msg"] = "out_trade_no错误"
		resStr := util.Map2Xml(resMap)
		c.JSON(http.StatusOK, resStr)
		return
	}

	if !auth_service.AddUserMargin(info.TakerUserId, info.TakerPayAmount) {
		resMap["return_code"] = "FAIL"
		resMap["return_msg"] = "out_trade_no错误"
		resStr := util.Map2Xml(resMap)
		c.JSON(http.StatusOK, resStr)
		return
	}

	resMap["return_code"] = "SUCCESS"
	resMap["return_msg"] = "OK"
	resStr := util.Map2Xml(resMap)
	c.JSON(http.StatusOK, resStr)
	return
}

func TeamUrgentWxNotify(c *gin.Context) {
	//logging.Info("--------pay")
	var resMap = make(map[string]interface{}, 0)
	resMap["return_code"] = "SUCCESS"
	resMap["return_msg"] = "OK"

	valueXml, _ := ioutil.ReadAll(c.Request.Body) //获取post的数据
	//logging.Info("--------pay:" + string(valueXml))
	values := util.Xml2Map(string(valueXml))

	if retCode, ok := values["result_code"]; retCode != "SUCCESS" || !ok {
		resMap["return_code"] = "FAIL"
		resMap["return_msg"] = "result_code错误"
		resStr := util.Map2Xml(resMap)
		c.JSON(http.StatusOK, resStr)
		return
	}
	if _, ok := values["out_trade_no"]; !ok {
		resMap["return_code"] = "FAIL"
		resMap["return_msg"] = "out_trade_no错误"
		resStr := util.Map2Xml(resMap)
		c.JSON(http.StatusOK, resStr)
		return
	}
	if _, ok := values["sign"]; !ok {
		resMap["return_code"] = "FAIL"
		resMap["return_msg"] = "sign错误"
		resStr := util.Map2Xml(resMap)
		c.JSON(http.StatusOK, resStr)
		return
	}

	//微信提交过来的签名
	postSign := values["sign"]
	delete(values, "sign")
	//根据提交过来的值，和我的商户支付秘钥，生成的签名
	userSign := order_service.WxPayCalcSign(values, var_const.WXMchKey)
	//验证提交过来的签名是否正确
	if userSign != postSign {
		resMap["return_code"] = "FAIL"
		resMap["return_msg"] = "sign错误"
		//logging.Info("--------pay-userSign")
		resStr := util.Map2Xml(resMap)
		c.JSON(http.StatusOK, resStr)
		return
	}

	//判断订单存在
	payOrderId := values["out_trade_no"].(string)
	info := models.Team{
		UrgentTradeNo: payOrderId,
	}
	_, err := info.GetTeamInfoByUrgentTradeNo()
	if err != nil {
		resMap["return_code"] = "FAIL"
		resMap["return_msg"] = "out_trade_no错误"
		resStr := util.Map2Xml(resMap)
		c.JSON(http.StatusOK, resStr)
		return
	}
	if info.UrgentPayStatus != 0 {
		resMap["return_code"] = "FAIL"
		resMap["return_msg"] = "out_trade_no错误"
		resStr := util.Map2Xml(resMap)
		c.JSON(http.StatusOK, resStr)
		return
	}

	//
	//order信息更新
	//保存支付订单 TODO
	dbInfo := models.Team{
		TeamId: info.TeamId,
	}
	var m = make(map[string]interface{})
	m["urgent_pay_status"] = 1
	m["urgent_pay_time"] = int(time.Now().Unix())
	if !dbInfo.Updates(m) {
		log, _ := json.Marshal(m)
		logging.Error("WxNotify:db-failed-" + string(log))
		resMap["return_code"] = "FAIL"
		resMap["return_msg"] = "out_trade_no错误"
		resStr := util.Map2Xml(resMap)
		c.JSON(http.StatusOK, resStr)
		return
	}

	resMap["return_code"] = "SUCCESS"
	resMap["return_msg"] = "OK"
	resStr := util.Map2Xml(resMap)
	c.JSON(http.StatusOK, resStr)
	return
}

//加急退款回调
func UrgentRefundCallback(c *gin.Context) {
	body, _ := ioutil.ReadAll(c.Request.Body)
	var mr RefundNotify
	_ = xml.Unmarshal(body, &mr)
	key := "t7v5TMsxhW6VH2f231NaB1BGL33CRjt3"
	b, _ := base64.StdEncoding.DecodeString(mr.Req_info)
	gocrypto.SetAesKey(strings.ToLower(gocrypto.Md5(key)))
	plaintext, _ := gocrypto.AesECBDecrypt(b)

	var mnr RefundNotify
	_ = xml.Unmarshal(plaintext, &mnr)

	if mnr.Refund_status == "SUCCESS" {
		//f(mnr.Out_trade_no)
		dbInfo := models.Team{
			UrgentRefundTradeNo: mnr.Out_refund_no,
		}
		_, _ = dbInfo.GetTeamInfoByUrgentRefundTradeNo()
		if dbInfo.UrgentTradeNo != mnr.Out_trade_no {
			c.JSON(http.StatusOK, nil)
			return
		}
		var m = make(map[string]interface{})
		m["urgent_refund_trade_no"] = ""
		m["urgent_refund_amount"] = 0
		m["urgent_refund_time"] = 0
		m["urgent_refund_status"] = 0
		m["is_urgent"] = 0
		m["urgent_user_id"] = 0
		m["urgent_nick_name"] = ""
		m["urgent_trade_no"] = ""
		m["urgent_pay_amount"] = 0
		m["urgent_pay_time"] = 0
		m["urgent_pay_status"] = 0
		if !dbInfo.Updates(m) {
			log, _ := json.Marshal(m)
			logging.Error("WxRefundCallback:failed-" + string(log))
			c.JSON(http.StatusOK, nil)
			return
		}
		resStr := "<xml><return_code>SUCCESS</return_code><return_msg>OK</return_msg></xml>"

		c.JSON(http.StatusOK, resStr)
		return
	} else {
		c.JSON(http.StatusOK, nil)
		return
	}
}
