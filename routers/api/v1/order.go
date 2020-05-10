package v1

import (
	"github.com/EDDYCJY/go-gin-example/models"
	"github.com/EDDYCJY/go-gin-example/pkg/app"
	"github.com/EDDYCJY/go-gin-example/pkg/e"
	"github.com/EDDYCJY/go-gin-example/service/auth_service"
	"github.com/EDDYCJY/go-gin-example/service/order_service"
	"github.com/gin-gonic/gin"
	"net/http"
)

// @Summary 下单
// @Produce  json
// @Param user_id body int false "user_id"
// @Param game_type body int false "游戏"
// @Param order_type body int false "订单类型"
// @Param instead_type body int false "代练类型"
// @Param game_zone body int false "游戏区服"
// @Param runes_level body int false "铭文等级"
// @Param hero_num body int false "英雄数量"
// @Param cur_level body int false "当前段位"
// @Param target_level body int false "目标段位"
// @Param game_account body string false "游戏账号"
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
	if !order_service.CreateOrder(&form) {
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
	if !order_service.CreateOrder(&form) {
		appG.Response(http.StatusBadRequest, e.ERROR, nil)
		return
	}
	data := gin.H{}
	data["order_id"] = form.OrderId
	appG.Response(http.StatusOK, e.SUCCESS, data)
	return
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
	if !order_service.CreateOrder(&form) {
		appG.Response(http.StatusBadRequest, e.ERROR, nil)
		return
	}
	data := gin.H{}
	data["order_id"] = form.OrderId
	appG.Response(http.StatusOK, e.SUCCESS, data)
	return
}

// @Summary Get 获取订单列表
// @Produce  json
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/order/all [get]
// @Tags 审核
func GetAllOrders(c *gin.Context) {
	appG := app.Gin{C: c}
	var list []models.Order
	order_service.GetNeedTakeOrderList(&list)
	appG.Response(http.StatusOK, e.SUCCESS, list)
}
