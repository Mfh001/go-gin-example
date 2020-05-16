package v1

import (
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
	"github.com/gin-gonic/gin"
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
	logging.Info("--------pay")
	var resMap = make(map[string]interface{}, 0)
	resMap["return_code"] = "SUCCESS"
	resMap["return_msg"] = "OK"

	valueXml, _ := ioutil.ReadAll(c.Request.Body) //获取post的数据
	logging.Info("--------pay:" + string(valueXml))
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
	logging.Info("--------info.Status")
	if info.Status != var_const.OrderStatusWaitPay {
		logging.Info("--------info.!Status")
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
	m["upd_time"] = int(time.Now().Unix())
	if !dbInfo.Updates(m) {
		log, _ := json.Marshal(m)
		logging.Error("WxNotify:db-failed-" + string(log))
		resMap["return_code"] = "FAIL"
		resMap["return_msg"] = "out_trade_no错误"
		resStr := util.Map2Xml(resMap)
		c.JSON(http.StatusOK, resStr)
		return
	}
	logging.Info("--------Updates")
	strResp := "<xml><return_code><![CDATA[SUCCESS]]></return_code><return_msg><![CDATA[OK]]></return_msg></xml>"
	c.JSON(http.StatusOK, strResp)
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
	if !auth_service.ExistUserInfo(userId) {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	//判断是代练 TODO
	//if !check_service.ExistUserCheck(userId) {
	//	appG.Response(http.StatusBadRequest, e.CHECK_NO_PASS, nil)
	//	return
	//}
	orderId, err := strconv.Atoi(c.PostForm("order_id"))
	if err != nil {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	if !order_service.ExistOrder(orderId) {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	status, err := order_service.GetOrderStatus(orderId)
	if err != nil || status != var_const.OrderStatusPaidPay {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	//不是代练不可接单 TODO

	d, ok := order_service.TakerPay(userId, orderId, c.ClientIP())
	if !ok {
		appG.Response(http.StatusBadRequest, e.ERROR, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, d)
	return
}

type WXPayNotifyResp struct {
	Return_code string `xml:"return_code"`
	Return_msg  string `xml:"return_msg"`
}

// 支付回调接口
func TakerWxNotify(c *gin.Context) {
	logging.Info("--------pay")
	//var resMap = make(map[string]interface{}, 0)
	var resp WXPayNotifyResp
	resp.Return_code = "SUCCESS"
	resp.Return_msg = "OK"

	valueXml, _ := ioutil.ReadAll(c.Request.Body) //获取post的数据
	logging.Info("--------pay:" + string(valueXml))
	values := util.Xml2Map(string(valueXml))

	if retCode, ok := values["result_code"]; retCode != "SUCCESS" || !ok {
		resp.Return_code = "FAIL"
		resp.Return_msg = "result_code错误"

		bytes, _ := xml.Marshal(resp)
		strResp := strings.Replace(string(bytes), "WXPayNotifyResp", "xml", -1)

		c.JSON(http.StatusOK, strResp)
		return
	}
	if _, ok := values["out_trade_no"]; !ok {
		resp.Return_code = "FAIL"
		resp.Return_msg = "result_code错误"

		bytes, _ := xml.Marshal(resp)
		strResp := strings.Replace(string(bytes), "WXPayNotifyResp", "xml", -1)

		c.JSON(http.StatusOK, strResp)
		return
	}
	if _, ok := values["sign"]; !ok {
		resp.Return_code = "FAIL"
		resp.Return_msg = "result_code错误"

		bytes, _ := xml.Marshal(resp)
		strResp := strings.Replace(string(bytes), "WXPayNotifyResp", "xml", -1)

		c.JSON(http.StatusOK, strResp)
		return
	}

	//微信提交过来的签名
	postSign := values["sign"]
	delete(values, "sign")
	//根据提交过来的值，和我的商户支付秘钥，生成的签名
	userSign := order_service.WxPayCalcSign(values, var_const.WXMchKey)
	//验证提交过来的签名是否正确
	if userSign != postSign {
		resp.Return_code = "FAIL"
		resp.Return_msg = "result_code错误"

		bytes, _ := xml.Marshal(resp)
		strResp := strings.Replace(string(bytes), "WXPayNotifyResp", "xml", -1)

		c.JSON(http.StatusOK, strResp)
		return
	}

	//判断订单存在
	payOrderId := values["out_trade_no"].(string)
	info := models.Order{
		TakerTradeNo: payOrderId,
	}
	_, err := info.GetOrderInfoByTakerTradeNo()
	if err != nil {
		resp.Return_code = "FAIL"
		resp.Return_msg = "result_code错误"

		bytes, _ := xml.Marshal(resp)
		strResp := strings.Replace(string(bytes), "WXPayNotifyResp", "xml", -1)

		c.JSON(http.StatusOK, strResp)
		return
	}
	if info.Status != var_const.OrderStatusTakerWaitPay {
		resp.Return_code = "FAIL"
		resp.Return_msg = "result_code错误"

		bytes, _ := xml.Marshal(resp)
		strResp := strings.Replace(string(bytes), "WXPayNotifyResp", "xml", -1)

		c.JSON(http.StatusOK, strResp)
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
	m["upd_time"] = int(time.Now().Unix())
	if !dbInfo.Updates(m) {
		log, _ := json.Marshal(m)
		logging.Error("WxNotify:db-failed-" + string(log))
		resp.Return_code = "FAIL"
		resp.Return_msg = "result_code错误"

		bytes, _ := xml.Marshal(resp)
		strResp := strings.Replace(string(bytes), "WXPayNotifyResp", "xml", -1)

		c.JSON(http.StatusOK, strResp)
		return
	}

	strResp := "<xml><return_code><![CDATA[SUCCESS]]></return_code><return_msg><![CDATA[OK]]></return_msg></xml>"
	c.JSON(http.StatusOK, strResp)
	return
}
