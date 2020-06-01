package v1

import (
	"encoding/json"
	"fmt"
	var_const "github.com/EDDYCJY/go-gin-example/const"
	"github.com/EDDYCJY/go-gin-example/models"
	"github.com/EDDYCJY/go-gin-example/pkg/app"
	"github.com/EDDYCJY/go-gin-example/pkg/e"
	"github.com/EDDYCJY/go-gin-example/pkg/gredis"
	"github.com/EDDYCJY/go-gin-example/pkg/logging"
	"github.com/EDDYCJY/go-gin-example/pkg/util"
	"github.com/EDDYCJY/go-gin-example/service/auth_service"
	"github.com/EDDYCJY/go-gin-example/service/order_service"
	"github.com/EDDYCJY/go-gin-example/service/profit_service"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"net/http"
	"strconv"
	"time"
)

// @Summary 下单
// @Produce  json
// @Param user_id body int false "user_id"
// @Param title body string false "标题"
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

// @Summary Get 获取订单列表
// @Produce  json
// @Param index body int false "index"
// @Param count body int false "count"
// @Param price_b body int false "开始价格"
// @Param price_e body int false "结束价格"
// @Param time_b body int false "最低时限"
// @Param time_e body int false "最高时限"
// @Param star_b body int false "最少星数"
// @Param star_e body int false "最多星数"
// @Param star_price_b body int false "最低每颗星平均价格"
// @Param star_price_e body int false "最高每颗星平均价格"
// @Param level_b body int false "最低段位"
// @Param level_e body int false "最高段位"
// @Param instead_type body int false "排位赛/巅峰赛"
// @Param zoom body int false "区服"
// @Param min_runes body int false "最低铭文等级"
// @Param max_runes body int false "最高铭文等级"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/order/all [get]
// @Tags 接单
func GetAllOrders(c *gin.Context) {
	appG := app.Gin{C: c}
	index := com.StrTo(c.Query("index")).MustInt()
	count := com.StrTo(c.Query("count")).MustInt()
	priceB := com.StrTo(c.Query("price_b")).MustInt()
	priceE := com.StrTo(c.Query("price_e")).MustInt()
	timeB := com.StrTo(c.Query("time_b")).MustInt()
	timeE := com.StrTo(c.Query("time_e")).MustInt()
	starB := com.StrTo(c.Query("star_b")).MustInt()
	starE := com.StrTo(c.Query("star_e")).MustInt()
	starPriceB := com.StrTo(c.Query("star_price_b")).MustInt()
	starPriceE := com.StrTo(c.Query("star_price_e")).MustInt()
	levelB := com.StrTo(c.Query("level_b")).MustInt()
	levelE := com.StrTo(c.Query("level_e")).MustInt()
	insteadType := com.StrTo(c.Query("instead_type")).MustInt()
	zoom := com.StrTo(c.Query("zoom")).MustInt()
	minRunes := com.StrTo(c.Query("min_runes")).MustInt()
	maxRunes := com.StrTo(c.Query("max_runes")).MustInt()

	var list []models.Order

	where := " "
	if priceB <= priceE && priceE > 0 {
		where = fmt.Sprintf(where+"and price >=%d and price <=%d ", priceB, priceE)
	}
	if timeB <= timeE && timeE > 0 {
		where = fmt.Sprintf(where+"and time_limit >=%d and time_limit <=%d ", timeB, timeE)
	}
	if starB <= starE && starE > 0 {
		where = fmt.Sprintf(where+"and star_num >=%d and star_num <=%d ", starB, starE)
	}
	if starPriceB <= starPriceE && starPriceE > 0 {
		where = fmt.Sprintf(where+"and star_per_price >=%d and star_per_price <=%d ", starPriceB, starPriceE)
	}
	if levelB <= levelE && levelE > 0 {
		where = fmt.Sprintf(where+"and cur_level >=%d and target_level <=%d ", levelB, levelE)
	}
	if minRunes <= maxRunes && minRunes >= 0 {
		where = fmt.Sprintf(where+"and runes_level >=%d and runes_level <=%d ", minRunes, maxRunes)
	}
	if zoom >= 0 {
		where = fmt.Sprintf(where+"and game_zone >=%d and game_zone <%d ", zoom*1000, (zoom+1)*1000)
	}
	if insteadType == 1 || insteadType == 0 {
		where = fmt.Sprintf(where+"and instead_type =%d ", insteadType)
	}

	logging.Info("where:" + where + ";index:" + c.Query("index") + ";count:" + c.Query("count"))
	order_service.GetNeedTakeOrderList(&list, where, index, count)
	appG.Response(http.StatusOK, e.SUCCESS, list)
}

// @Summary Get 获取订单列表
// @Produce  json
// @Param index body int false "index"
// @Param count body int false "count"
// @Param instead_type body int false "排位赛/巅峰赛"
// @Param zoom body int false "区服"
// @Param min_runes body int false "最低铭文等级"
// @Param max_runes body int false "最高铭文等级"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/order/sortall [get]
// @Tags 接单
func GetAllOrdersB(c *gin.Context) {
	appG := app.Gin{C: c}
	index := com.StrTo(c.Query("index")).MustInt()
	count := com.StrTo(c.Query("count")).MustInt()
	insteadType := com.StrTo(c.Query("instead_type")).MustInt()
	zoom := com.StrTo(c.Query("zoom")).MustInt()
	minRunes := com.StrTo(c.Query("min_runes")).MustInt()
	maxRunes := com.StrTo(c.Query("max_runes")).MustInt()

	var list []models.Order

	where := " "
	if minRunes <= maxRunes && minRunes >= 0 {
		where = fmt.Sprintf(where+"and runes_level >=%d and runes_level <=%d ", minRunes, maxRunes)
	}
	if zoom >= 0 {
		where = fmt.Sprintf(where+"and game_zone >=%d and game_zone <%d ", zoom*1000, (zoom+1)*1000)
	}
	if insteadType == 1 || insteadType == 0 {
		where = fmt.Sprintf(where+"and instead_type =%d ", insteadType)
	}
	logging.Info("where:" + where + ";index:" + c.Query("index") + ";count:" + c.Query("count"))
	order_service.GetNeedTakeOrderList(&list, where, index, count)
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
	if userId == 0 || !auth_service.ExistUserInfo(userId) || !auth_service.IsUserTypeInstead(userId) {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
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
	takerId := order_service.GetOrderParam(orderId, "taker_user_id")
	if takerId <= 0 || !auth_service.ExistUserInfo(takerId) {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	add := order_service.GetOrderParam(orderId, "price")
	//收益
	logging.Info("ConfirmOrder: price-" + strconv.Itoa(add) + " order_id" + strconv.Itoa(orderId))
	if add >= var_const.OrderNeedRate && add < var_const.OrderNeedRateMax {
		add = add * (100 - var_const.OrderRate) / 100
	}
	if !auth_service.AddUserBalance(takerId, add, "ConfirmOrder") {
		appG.Response(http.StatusBadRequest, e.ERROR, nil)
		return
	}
	logging.Info("ConfirmOrder: end")
	order_service.Refund(orderId)

	//代理逻辑计算
	price := order_service.GetOrderParam(orderId, "price")
	if price >= var_const.OrderNeedRate && price < var_const.OrderNeedRateMax {
		if !profit_service.ExistProfit(userId) {
			profit := models.Profit{
				UserId: userId,
			}
			if !profit.Insert() {
				logging.Error("ConfirmOrder:insert Profit -failed-" + c.PostForm("user_id"))
			}
		}
		if !profit_service.ExistProfit(takerId) {
			profit := models.Profit{
				UserId: takerId,
			}
			if !profit.Insert() {
				logging.Error("ConfirmOrder:insert takerId Profit -failed-" + strconv.Itoa(takerId))
			}
		}
		{
			agentId := auth_service.GetUserParam(userId, "agent_id")
			if agentId > 0 && !profit_service.ExistProfit(agentId) {
				profit := models.Profit{
					UserId: agentId,
				}
				if !profit.Insert() {
					logging.Error("ConfirmOrder:insert takerId Profit -failed-" + strconv.Itoa(agentId))
				}
			}
		}
		{
			agentId := auth_service.GetUserParam(takerId, "agent_id")
			if agentId > 0 && !profit_service.ExistProfit(agentId) {
				profit := models.Profit{
					UserId: agentId,
				}
				if !profit.Insert() {
					logging.Error("ConfirmOrder:insert takerId Profit -failed-" + strconv.Itoa(agentId))
				}
			}
		}

		//--累计订单统计
		{
			userOrderTotalTimes := profit_service.GetProfitParam(userId, "order_total_times")
			userOrderTotalTimes++
			userProfit := models.Profit{
				UserId: userId,
			}
			m := make(map[string]interface{})
			m["order_total_times"] = userOrderTotalTimes
			userProfit.Updates(m)
		}
		{
			takerOrderTotalTimes := profit_service.GetProfitParam(takerId, "order_total_times")
			takerOrderTotalTimes++
			userProfit := models.Profit{
				UserId: takerId,
			}
			m := make(map[string]interface{})
			m["order_total_times"] = takerOrderTotalTimes
			userProfit.Updates(m)
		}
		{
			//今天下级发单统计
			agentId := auth_service.GetUserParam(userId, "agent_id")
			if agentId > 0 && auth_service.ExistUserInfo(agentId) {
				resetTime := profit_service.GetProfitParam(agentId, "reset_time")
				if !util.IsToday(resetTime) {
					//计算返利 重置当日数据
					profit_service.CalDailyProfit(agentId)
				}
				//今天发单统计
				orderTodayPublishTimes := profit_service.GetProfitParam(agentId, "order_today_publish_times")
				orderTodayPublishTimes++
				userProfit := models.Profit{
					UserId: agentId,
				}
				m := make(map[string]interface{})
				m["order_today_publish_times"] = orderTodayPublishTimes
				//m["reset_time"] = int(time.Now().Unix())
				userProfit.Updates(m)
			}
		}
		{
			//今天下级接单统计
			agentId := auth_service.GetUserParam(takerId, "agent_id")
			if agentId > 0 && auth_service.ExistUserInfo(agentId) {
				resetTime := profit_service.GetProfitParam(agentId, "reset_time")
				if !util.IsToday(resetTime) {
					//计算返利 重置当日数据
					profit_service.CalDailyProfit(agentId)
				}
				//今天发单统计
				orderTodayTakerTimes := profit_service.GetProfitParam(agentId, "order_today_taker_times")
				orderTodayTakerTimes++
				userProfit := models.Profit{
					UserId: agentId,
				}
				m := make(map[string]interface{})
				m["order_today_taker_times"] = orderTodayTakerTimes
				//m["reset_time"] = int(time.Now().Unix())
				userProfit.Updates(m)
			}
		}
	}

	//

	appG.Response(http.StatusOK, e.SUCCESS, nil)
	return
}

// @Summary 上传订单截图
// @Produce  json
// @Param user_id body int false "user_id"
// @Param order_id body int false "order_id"
// @Param img_url body string false "img_url"
// @Param img_type body int false "img_type 1接单的图片 2完成订单的图片"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/order/updorderimg [post]
// @Tags 订单
func UpdateOrderImg(c *gin.Context) {
	var (
		appG    = app.Gin{C: c}
		userId  = com.StrTo(c.PostForm("user_id")).MustInt()
		orderId = com.StrTo(c.PostForm("order_id")).MustInt()
		imgUrl  = c.PostForm("img_url")
		imgType = com.StrTo(c.PostForm("img_type")).MustInt()
	)
	if userId == 0 || !auth_service.ExistUserInfo(userId) || !auth_service.IsUserTypeInstead(userId) || !order_service.ExistOrder(orderId) {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	dbInfo := models.Order{
		OrderId: orderId,
	}
	logging.Info("ConfirmOrder: begin order_id-" + strconv.Itoa(orderId))
	var m = make(map[string]interface{})
	if imgType == 1 {
		m["img_take_url"] = imgUrl
	} else {
		m["img_finish_url"] = imgUrl
	}

	if !dbInfo.Updates(m) {
		log, _ := json.Marshal(m)
		logging.Error("ConfirmOrder:db1-failed-" + string(log))
		appG.Response(http.StatusBadRequest, e.ERROR, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
	return
}

// @Summary 取消订单
// @Produce  json
// @Param user_id body int false "user_id"
// @Param order_id body int false "order_id"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/order/cancel [post]
// @Tags 订单
func CancelOrder(c *gin.Context) {
	var (
		appG    = app.Gin{C: c}
		userId  = com.StrTo(c.PostForm("user_id")).MustInt()
		orderId = com.StrTo(c.PostForm("order_id")).MustInt()
	)
	if userId == 0 || !auth_service.ExistUserInfo(userId) {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	if userId != order_service.GetOrderParam(orderId, "user_id") {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	status := order_service.GetOrderParam(orderId, "status")
	if status != var_const.OrderStatusPaidPay && status != var_const.OrderStatusWaitPay {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	order_service.OrderCancelRefund(orderId)
	//dbInfo := models.Order{
	//	OrderId: orderId,
	//}
	//var m = make(map[string]interface{})
	//m["status"] = var_const.OrderStatusTakerFinishedNeedConfirm
	//m["upd_time"] = int(time.Now().Unix())
	//logging.Info("FinishOrder: begin order_id-" + strconv.Itoa(orderId))
	//if !dbInfo.Updates(m) {
	//	log, _ := json.Marshal(m)
	//	logging.Error("TakerFinish:failed-" + string(log))
	//	appG.Response(http.StatusBadRequest, e.ERROR, nil)
	//	return
	//}
	//logging.Info("FinishOrder: end")
	appG.Response(http.StatusOK, e.SUCCESS, nil)
	return
}

type Message struct {
	UserId        int    `json:"user_id"`
	SendId        int    `json:"send_id"`
	NickName      string `json:"nick_name"`
	TakerUserId   int    `json:"taker_user_id"`
	TakerNickName string `json:"taker_nick_name"`
	OrderId       int    `json:"order_id"`
	OrderNo       string `json:"order_no"`
	OrderTitle    string `json:"order_title"`
	Time          int    `json:"time"`
	Msg           string `json:"msg"`
}

// @Summary 订单留言
// @Produce  json
// @Param user_id body int false "user_id"
// @Param order_id body int false "order_id"
// @Param message body string false "message"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/order/message [post]
// @Tags 订单留言
func MessageOrder(c *gin.Context) {
	var (
		appG    = app.Gin{C: c}
		userId  = com.StrTo(c.PostForm("user_id")).MustInt()
		orderId = com.StrTo(c.PostForm("order_id")).MustInt()
		message = c.PostForm("message")
	)
	if userId == 0 || !auth_service.ExistUserInfo(userId) {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	u := order_service.GetOrderParam(orderId, "user_id")
	tu := order_service.GetOrderParam(orderId, "taker_user_id")
	if userId != u && userId != tu {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	msgUserId := u
	if userId == u {
		msgUserId = tu
	}
	//status := order_service.GetOrderParam(orderId, "status")
	//if status != var_const.OrderStatusPaidPay {
	//	appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
	//	return
	//}
	msg := Message{
		UserId:        u,
		SendId:        userId,
		NickName:      auth_service.GetUserParamString(u, "nick_name"),
		TakerUserId:   tu,
		TakerNickName: auth_service.GetUserParamString(tu, "nick_name"),
		OrderId:       orderId,
		OrderNo:       order_service.GetOrderParamString(orderId, "order_no"),
		OrderTitle:    order_service.GetOrderParamString(orderId, "title"),
		Time:          int(time.Now().Unix()),
		Msg:           message,
	}
	jsonData, _ := json.Marshal(msg)
	_ = gredis.LPush(order_service.GetRedisKeyMessageOrder(orderId), string(jsonData))
	_ = gredis.LPush(order_service.GetRedisKeyMessageUser(msgUserId), string(jsonData))
	_, _ = gredis.HSet(order_service.GetRedisKeyMessageNoRead(), strconv.Itoa(msgUserId), "1")
	appG.Response(http.StatusOK, e.SUCCESS, nil)
	return
}

// @Summary 获取某一订单的消息/留言
// @Produce  json
// @Param user_id body int false "order_id"
// @Param index body int false "index"
// @Param count body int false "count"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/order/getmessage [post]
// @Tags 订单留言
func GetOrderMessage(c *gin.Context) {
	var (
		appG    = app.Gin{C: c}
		orderId = com.StrTo(c.PostForm("order_id")).MustInt()
		index   = com.StrTo(c.PostForm("index")).MustInt()
		count   = com.StrTo(c.PostForm("count")).MustInt()
	)
	if orderId == 0 || !order_service.ExistOrder(orderId) {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	res, _ := gredis.LRange(order_service.GetRedisKeyMessageOrder(orderId), index, index+count)
	appG.Response(http.StatusOK, e.SUCCESS, res)
	return
}

// @Summary 代练请求撤销订单
// @Produce  json
// @Param user_id body int false "user_id"
// @Param order_id body int false "order_id"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/order/undorequest [post]
// @Tags 订单
func UndoRequestOrder(c *gin.Context) {
	var (
		appG    = app.Gin{C: c}
		userId  = com.StrTo(c.PostForm("user_id")).MustInt()
		orderId = com.StrTo(c.PostForm("order_id")).MustInt()
	)
	if userId == 0 || !auth_service.ExistUserInfo(userId) {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	if userId != order_service.GetOrderParam(orderId, "taker_user_id") {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	status := order_service.GetOrderParam(orderId, "status")
	if status != var_const.OrderStatusTakerPaid {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	dbInfo := models.Order{
		OrderId: orderId,
	}
	var m = make(map[string]interface{})
	m["status"] = var_const.OrderStatusUndoRequest
	m["upd_time"] = int(time.Now().Unix())
	logging.Info("UndoRequestOrder: begin order_id-" + strconv.Itoa(orderId))
	if !dbInfo.Updates(m) {
		log, _ := json.Marshal(m)
		logging.Error("UndoRequestOrder:failed-" + string(log))
		appG.Response(http.StatusBadRequest, e.ERROR, nil)
		return
	}
	logging.Info("UndoRequestOrder: end")
	appG.Response(http.StatusOK, e.SUCCESS, nil)
	return
}

// @Summary 用户回应代练的撤销订单请求
// @Produce  json
// @Param user_id body int false "user_id"
// @Param order_id body int false "order_id"
// @Param agree body int false "agree"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/order/undo [post]
// @Tags 订单
func UndoOrder(c *gin.Context) {
	var (
		appG    = app.Gin{C: c}
		userId  = com.StrTo(c.PostForm("user_id")).MustInt()
		orderId = com.StrTo(c.PostForm("order_id")).MustInt()
		agree   = com.StrTo(c.PostForm("agree")).MustInt()
	)
	if userId == 0 || !auth_service.ExistUserInfo(userId) {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	if userId != order_service.GetOrderParam(orderId, "user_id") {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	status := order_service.GetOrderParam(orderId, "status")
	if status != var_const.OrderStatusUndoRequest {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	dbInfo := models.Order{
		OrderId: orderId,
	}
	var m = make(map[string]interface{})
	if agree == 1 {
		m["status"] = var_const.OrderStatusUndoRequest
	} else {
		m["status"] = var_const.OrderStatusTakerPaid
	}
	m["upd_time"] = int(time.Now().Unix())
	logging.Info("UndoRequestOrder: begin order_id-" + strconv.Itoa(orderId))
	if !dbInfo.Updates(m) {
		log, _ := json.Marshal(m)
		logging.Error("UndoRequestOrder:failed-" + string(log))
		appG.Response(http.StatusBadRequest, e.ERROR, nil)
		return
	}
	if agree == 1 {
		order_service.OrderUndoRefundUser(orderId)
		order_service.OrderUndoRefundTaker(orderId)
	}
	logging.Info("UndoRequestOrder: end")
	appG.Response(http.StatusOK, e.SUCCESS, nil)
	return
}

// @Summary 订单加时
// @Produce  json
// @Param user_id body int false "user_id"
// @Param order_id body int false "order_id"
// @Param time body int false "time"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/order/addtime [post]
// @Tags 订单
func AddTimeOrder(c *gin.Context) {
	var (
		appG    = app.Gin{C: c}
		userId  = com.StrTo(c.PostForm("user_id")).MustInt()
		orderId = com.StrTo(c.PostForm("order_id")).MustInt()
		addTime = com.StrTo(c.PostForm("time")).MustInt()
	)
	if userId == 0 || !auth_service.ExistUserInfo(userId) {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	if userId != order_service.GetOrderParam(orderId, "user_id") {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	//status := order_service.GetOrderParam(orderId, "status")
	//if status != var_const.OrderStatusUndoRequest {
	//	appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
	//	return
	//}
	dbInfo := models.Order{
		OrderId: orderId,
	}
	var m = make(map[string]interface{})
	m["time_limit"] = order_service.GetOrderParam(orderId, "time_limit") + addTime
	m["upd_time"] = int(time.Now().Unix())
	logging.Info("AddTimeOrder: begin order_id-" + strconv.Itoa(orderId))
	if !dbInfo.Updates(m) {
		log, _ := json.Marshal(m)
		logging.Error("AddTimeOrder:failed-" + string(log))
		appG.Response(http.StatusBadRequest, e.ERROR, nil)
		return
	}
	logging.Info("AddTimeOrder: end")
	appG.Response(http.StatusOK, e.SUCCESS, nil)
	return
}

// @Summary 获取订单信息
// @Produce  json
// @Param user_id body int false "user_id"
// @Param order_id body int false "order_id"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/order/info [post]
// @Tags 订单
func GetOrderInfo(c *gin.Context) {
	var (
		appG    = app.Gin{C: c}
		userId  = com.StrTo(c.PostForm("user_id")).MustInt()
		orderId = com.StrTo(c.PostForm("order_id")).MustInt()
	)
	if userId == 0 || !auth_service.ExistUserInfo(userId) {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	//if userId != order_service.GetOrderParam(orderId, "user_id") {
	//	appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
	//	return
	//}
	//status := order_service.GetOrderParam(orderId, "status")
	//if status != var_const.OrderStatusUndoRequest {
	//	appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
	//	return
	//}
	dbInfo := models.Order{
		OrderId: orderId,
	}
	_, _ = dbInfo.First()
	appG.Response(http.StatusOK, e.SUCCESS, dbInfo)
	return
}

// @Summary 申请仲裁
// @Produce  json
// @Param user_id body int false "user_id"
// @Param order_id body int false "order_id"
// @Param msg body string false "msg"
// @Param imgurl body string false "imgurl"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/order/adjudgerequest [post]
// @Tags 订单
func AdjudgeRequestOrder(c *gin.Context) {
	var (
		appG    = app.Gin{C: c}
		userId  = com.StrTo(c.PostForm("user_id")).MustInt()
		orderId = com.StrTo(c.PostForm("order_id")).MustInt()
		msg     = c.PostForm("msg")
		imgUrl  = c.PostForm("imgurl")
	)
	if userId == 0 || !auth_service.ExistUserInfo(userId) {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	if userId != order_service.GetOrderParam(orderId, "taker_user_id") && userId != order_service.GetOrderParam(orderId, "user_id") {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	status := order_service.GetOrderParam(orderId, "status")
	if status != var_const.OrderStatusTakerPaid {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	dbInfo := models.Order{
		OrderId: orderId,
	}
	var m = make(map[string]interface{})
	m["status"] = var_const.OrderStatusAdjudgeRequest
	m["adjudge_user_id"] = userId
	m["adjudge_msg"] = msg
	m["img_adjudge_url"] = imgUrl
	m["upd_time"] = int(time.Now().Unix())
	logging.Info("AdjudgeRequestOrder: begin order_id-" + strconv.Itoa(orderId))
	if !dbInfo.Updates(m) {
		log, _ := json.Marshal(m)
		logging.Error("AdjudgeRequestOrder:failed-" + string(log))
		appG.Response(http.StatusBadRequest, e.ERROR, nil)
		return
	}
	logging.Info("AdjudgeRequestOrder: end")
	appG.Response(http.StatusOK, e.SUCCESS, nil)
	return
}

// @Summary 客服退回订单的部分金额
// @Produce  json
// @Param order_id body int false "order_id"
// @Param money body int false "money 金额单位是分"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /admin/order/refund/user [post]
// @Tags 客服
func AdminRefundPay(c *gin.Context) {
	var (
		appG    = app.Gin{C: c}
		orderId = com.StrTo(c.PostForm("order_id")).MustInt()
		money   = com.StrTo(c.PostForm("money")).MustInt()
	)
	if orderId == 0 || !order_service.ExistOrder(orderId) {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	if money <= 0 || money > order_service.GetOrderParam(orderId, "pay_amount") {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	status := order_service.GetOrderParam(orderId, "status")
	if status != var_const.OrderStatusAdjudgeRequest {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	if order_service.GetOrderParamString(orderId, "pay_refund_trade_no") != "" {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	logging.Info("AdminRefundPay: begin order_id-" + strconv.Itoa(orderId))
	order_service.AdminRefundUser(orderId, money)
	logging.Info("AdminRefundPay: end")
	appG.Response(http.StatusOK, e.SUCCESS, nil)
	return
}

// @Summary 客服退回代练的部分保证金
// @Produce  json
// @Param order_id body int false "order_id"
// @Param money body int false "money 金额单位是分"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /admin/order/refund/taker [post]
// @Tags 客服
func AdminRefundTaker(c *gin.Context) {
	var (
		appG    = app.Gin{C: c}
		orderId = com.StrTo(c.PostForm("order_id")).MustInt()
		money   = com.StrTo(c.PostForm("money")).MustInt()
	)
	if orderId == 0 || !order_service.ExistOrder(orderId) {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	if money <= 0 || money > order_service.GetOrderParam(orderId, "taker_pay_amount") {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	status := order_service.GetOrderParam(orderId, "status")
	if status != var_const.OrderStatusAdjudgeRequest {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	if order_service.GetOrderParamString(orderId, "refund_trade_no") != "" {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	logging.Info("AdminRefundTaker: begin order_id-" + strconv.Itoa(orderId))
	order_service.AdminRefundTaker(orderId, money)
	logging.Info("AdminRefundTaker: end")
	appG.Response(http.StatusOK, e.SUCCESS, nil)
	return
}

// @Summary 客服添加用户余额
// @Produce  json
// @Param user_id body int false "user_id"
// @Param money body int false "money 金额单位是分"
// @Param msg body int false "备注"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /admin/money/add [post]
// @Tags 客服
func AdminAddBalance(c *gin.Context) {
	var (
		appG   = app.Gin{C: c}
		userId = com.StrTo(c.PostForm("user_id")).MustInt()
		money  = com.StrTo(c.PostForm("money")).MustInt()
		msg    = c.PostForm("msg")
	)
	if userId == 0 || !auth_service.ExistUserInfo(userId) {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	if money <= 0 {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	logging.Info("AdminAddBalance: begin userId-" + strconv.Itoa(money))
	if !auth_service.AddUserBalance(userId, money, "客服添加余额:"+msg) {
		appG.Response(http.StatusBadRequest, e.ERROR, nil)
		return
	}
	logging.Info("AdminAddBalance: end")
	appG.Response(http.StatusOK, e.SUCCESS, nil)
	return
}

// @Summary 客服扣除用户余额
// @Produce  json
// @Param user_id body int false "user_id"
// @Param money body int false "money 金额单位是分"
// @Param msg body int false "备注"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /admin/money/remove [post]
// @Tags 客服
func AdminRemBalance(c *gin.Context) {
	var (
		appG   = app.Gin{C: c}
		userId = com.StrTo(c.PostForm("user_id")).MustInt()
		money  = com.StrTo(c.PostForm("money")).MustInt()
		msg    = c.PostForm("msg")
	)
	if userId == 0 || !auth_service.ExistUserInfo(userId) {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	if money <= 0 {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	logging.Info("AdminRemBalance: begin userId-" + strconv.Itoa(money))
	if !auth_service.RemoveUserBalance(userId, money, "客服扣除余额:"+msg) {
		appG.Response(http.StatusBadRequest, e.ERROR, nil)
		return
	}
	logging.Info("AdminRemBalance: end")
	appG.Response(http.StatusOK, e.SUCCESS, nil)
	return
}

// @Summary 客服扣除代练保证金
// @Produce  json
// @Param user_id body int false "user_id"
// @Param money body int false "money 金额单位是分"
// @Param msg body int false "备注"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /admin/margin/remove [post]
// @Tags 客服
func AdminRemMargin(c *gin.Context) {
	var (
		appG   = app.Gin{C: c}
		userId = com.StrTo(c.PostForm("user_id")).MustInt()
		money  = com.StrTo(c.PostForm("money")).MustInt()
		msg    = c.PostForm("msg")
	)
	if userId == 0 || !auth_service.ExistUserInfo(userId) {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	if money <= 0 {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	logging.Info("AdminRemMargin: begin userId-" + strconv.Itoa(money))
	if !auth_service.RemoveUserMargin(userId, money, "客服扣除保证金:"+msg) {
		appG.Response(http.StatusBadRequest, e.ERROR, nil)
		return
	}
	logging.Info("AdminRemMargin: end")
	appG.Response(http.StatusOK, e.SUCCESS, nil)
	return
}

// @Summary Get 客服获取需要仲裁列表
// @Produce  json
// @Param index body int false "index"
// @Param count body int false "count"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /admin/order/adjudge/all [post]
// @Tags 客服
func GetAdminAdjudgeList(c *gin.Context) {
	appG := app.Gin{C: c}
	index := com.StrTo(c.PostForm("index")).MustInt()
	count := com.StrTo(c.PostForm("count")).MustInt()
	var list []models.Order
	order_service.GetNeedAdjudgeOrderList(&list, "", index, count)
	m := make(map[string]interface{})
	m["list"] = list
	m["count"] = len(list)
	appG.Response(http.StatusOK, e.SUCCESS, m)
}
