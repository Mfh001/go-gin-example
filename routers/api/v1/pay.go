package v1

import (
	"encoding/json"
	var_const "github.com/EDDYCJY/go-gin-example/const"
	"github.com/EDDYCJY/go-gin-example/models"
	"github.com/EDDYCJY/go-gin-example/pkg/app"
	"github.com/EDDYCJY/go-gin-example/pkg/e"
	"github.com/EDDYCJY/go-gin-example/pkg/logging"
	"github.com/EDDYCJY/go-gin-example/pkg/util"
	"github.com/EDDYCJY/go-gin-example/service/auth_service"
	"github.com/EDDYCJY/go-gin-example/service/order_service"
	"github.com/EDDYCJY/go-gin-example/service/pay_service"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

// @Summary 支付下单
// @Produce  json
// @Param user_id body int false "user_id"
// @Param pay_amount body int false "价格 单位分"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/pay/order [post]
// @Tags 微信支付
func AddPayOrder(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form models.Pay
	)
	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}
	if form.PayAmount <= 0 {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	if !auth_service.ExistUserInfo(form.UserId) {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	if !pay_service.CreateOrder(&form) {
		appG.Response(http.StatusBadRequest, e.ERROR, nil)
		return
	}
	d, ok := pay_service.Pay(form.UserId, form.OrderId, c.ClientIP())
	if !ok {
		appG.Response(http.StatusBadRequest, e.ERROR, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, d)
	return
}

//// @Summary 余额充值 获取发起微信支付所需的数据
//// @Produce  json
//// @Param user_id body int false "user_id"
//// @Param order_id body int false "order_id"
//// @Success 200 {object} app.Response
//// @Failure 500 {object} app.Response
//// @Router /api/v1/pay [post]
//// @Tags 微信支付
//func Pay(c *gin.Context) {
//	appG := app.Gin{C: c}
//
//	userId, err := strconv.Atoi(c.PostForm("user_id"))
//	if err != nil {
//		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
//		return
//	}
//	if !auth_service.ExistUserInfo(userId) {
//		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
//		return
//	}
//	orderId, err := strconv.Atoi(c.PostForm("order_id"))
//	if err != nil {
//		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
//		return
//	}
//	if !pay_service.ExistOrder(orderId) {
//		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
//		return
//	}
//
//	d, ok := pay_service.Pay(userId, orderId, c.ClientIP())
//	if !ok {
//		appG.Response(http.StatusBadRequest, e.ERROR, nil)
//		return
//	}
//	appG.Response(http.StatusOK, e.SUCCESS, d)
//	return
//}

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
	info := models.Pay{
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
	if info.Status != var_const.PayOrderStatusBegin {
		//logging.Info("--------info.!Status")
		resMap["return_code"] = "FAIL"
		resMap["return_msg"] = "out_trade_no错误"
		resStr := util.Map2Xml(resMap)
		c.JSON(http.StatusOK, resStr)
		return
	}

	//
	dbInfo := models.Pay{
		OrderId: info.OrderId,
	}
	var m = make(map[string]interface{})
	m["status"] = var_const.PayOrderStatusFinish
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

	//加余额
	auth_service.AddUserBalance(info.UserId, info.PayAmount, "充值")

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
		appG.Response(http.StatusBadRequest, e.HAS_DEPOSIT, nil)
		return
	}
	ok := order_service.DepositPay(userId)
	if ok > 0 {
		if ok == 1 {
			appG.Response(http.StatusBadRequest, e.MONEY_NO_ENOUGH, nil)
			return
		} else {
			appG.Response(http.StatusBadRequest, e.ERROR, nil)
			return
		}
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
	return
}
