package v1

import (
	"encoding/json"
	var_const "github.com/EDDYCJY/go-gin-example/const"
	"github.com/EDDYCJY/go-gin-example/models"
	"github.com/EDDYCJY/go-gin-example/pkg/app"
	"github.com/EDDYCJY/go-gin-example/pkg/e"
	"github.com/EDDYCJY/go-gin-example/pkg/logging"
	"github.com/EDDYCJY/go-gin-example/service/auth_service"
	"github.com/EDDYCJY/go-gin-example/service/order_service"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"net/http"
	"strconv"
	"time"
)

// @Summary 下单
// @Produce  json
// @Param user_id body int false "user_id"
// @Param game_type body int false "游戏"
// @Param price body int false "价格"
// @Param time_limit body int false "时限"
// @Param order_type body int false "订单类型"
// @Param instead_type body int false "代练类型"
// @Param game_zone body int false "游戏区服"
// @Param runes_level body int false "铭文等级"
// @Param hero_num body int false "英雄数量"
// @Param cur_level body int false "当前段位"
// @Param target_level body int false "目标段位"
// @Param game_acc body string false "游戏账号"
// @Param game_pwd body string false "游戏密码"
// @Param game_role body string false "游戏角色名"
// @Param game_phone body string false "验证手机"
// @Param margin body int false "保证金"
// @Param anti_addiction body int false "有防沉迷"
// @Param designate_hero body int false "有指定英雄"
// @Param hero_name body string false "指定英雄"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/order [post]
// @Tags 下单
func AddOrder(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form models.Order
	)
	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}
	if !auth_service.ExistUserInfo(form.UserId) {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	if form.CurLevel >= form.TargetLevel || form.TargetLevel <= 0 {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	if !order_service.CreateOrder(&form, 0, 0) {
		appG.Response(http.StatusBadRequest, e.ERROR, nil)
		return
	}
	data := gin.H{}
	data["order_id"] = form.OrderId
	appG.Response(http.StatusOK, e.SUCCESS, data)
	return
}

// @Summary 接单
// @Produce  json
// @Param user_id body int false "user_id"
// @Param order_id body int false "order_id"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/order/take [post]
// @Tags 接单
func TakeOrder(c *gin.Context) {

}

// @Summary 绑定上级
// @Produce  json
// @Param user_id body int false "user_id"
// @Param agent_id body int false "agent_id"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/agent/bind [post]
// @Tags 代理
func BindAgent(c *gin.Context) {

}

// @Summary Get 获取订单列表
// @Produce  json
// @Param index body int false "index"
// @Param count body int false "count"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/order/all [get]
// @Tags 接单
func GetAllOrders(c *gin.Context) {
	appG := app.Gin{C: c}
	index := com.StrTo(c.Query("index")).MustInt()
	count := com.StrTo(c.Query("count")).MustInt()
	var list []models.Order
	order_service.GetNeedTakeOrderList(&list, index, count)
	appG.Response(http.StatusOK, e.SUCCESS, list)
}

// @Summary 完成订单
// @Produce  json
// @Param user_id body int false "user_id"
// @Param order_id body int false "order_id"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/order/finish [post]
// @Tags 订单
func FinishOrder(c *gin.Context) {
	var (
		appG    = app.Gin{C: c}
		userId  = com.StrTo(c.PostForm("user_id")).MustInt()
		orderId = com.StrTo(c.PostForm("order_id")).MustInt()
	)
	if userId == 0 || !auth_service.ExistUserInfo(userId) {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	//必须是代练 TODO
	//if !check_service.ExistUserCheck(userId) {
	//	appG.Response(http.StatusBadRequest, e.CHECK_NO_PASS, nil)
	//	return
	//}
	status := order_service.GetOrderParam(orderId, "status")
	if status != var_const.OrderStatusTakerPaid {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	takerId := order_service.GetOrderParam(orderId, "taker_user_id")
	if takerId != userId {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	dbInfo := models.Order{
		OrderId: orderId,
	}
	var m = make(map[string]interface{})
	m["status"] = var_const.OrderStatusTakerFinishedNeedConfirm
	m["upd_time"] = int(time.Now().Unix())
	logging.Info("FinishOrder: begin order_id-" + strconv.Itoa(orderId))
	if !dbInfo.Updates(m) {
		log, _ := json.Marshal(m)
		logging.Error("TakerFinish:failed-" + string(log))
		appG.Response(http.StatusBadRequest, e.ERROR, nil)
		return
	}
	logging.Info("FinishOrder: end")
	appG.Response(http.StatusOK, e.SUCCESS, nil)
	return
}

// @Summary 完成订单
// @Produce  json
// @Param user_id body int false "user_id"
// @Param order_id body int false "order_id"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/order/confirm [post]
// @Tags 订单
func ConfirmOrder(c *gin.Context) {
	var (
		appG    = app.Gin{C: c}
		userId  = com.StrTo(c.PostForm("user_id")).MustInt()
		orderId = com.StrTo(c.PostForm("order_id")).MustInt()
	)
	if userId == 0 || !auth_service.ExistUserInfo(userId) {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	//if !check_service.ExistUserCheck(userId){
	//	appG.Response(http.StatusBadRequest, e.CHECK_NO_PASS, nil)
	//	return
	//}
	status := order_service.GetOrderParam(orderId, "status")
	if status != var_const.OrderStatusTakerFinishedNeedConfirm {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	uId := order_service.GetOrderParam(orderId, "user_id")
	if uId != userId {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	dbInfo := models.Order{
		OrderId: orderId,
	}
	logging.Info("ConfirmOrder: begin order_id-" + strconv.Itoa(orderId))
	var m = make(map[string]interface{})
	m["status"] = var_const.OrderStatusConfirmFinished
	m["upd_time"] = int(time.Now().Unix())
	if !dbInfo.Updates(m) {
		log, _ := json.Marshal(m)
		logging.Error("ConfirmOrder:db1-failed-" + string(log))
		appG.Response(http.StatusBadRequest, e.ERROR, nil)
		return
	}
	//收益
	takerId := order_service.GetOrderParam(orderId, "taker_user_id")
	add := order_service.GetOrderParam(orderId, "price")
	logging.Info("ConfirmOrder: price-" + strconv.Itoa(add))
	if add >= var_const.OrderNeedRate {
		add = add * (100 - var_const.OrderRate) / 100
	}
	if !auth_service.AddUserBalance(takerId, add) {
		appG.Response(http.StatusBadRequest, e.ERROR, nil)
		return
	}
	logging.Info("ConfirmOrder: end")
	//order_service.Refund(13)
	order_service.Refund(orderId)

	appG.Response(http.StatusOK, e.SUCCESS, nil)
	return
}
