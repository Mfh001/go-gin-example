package v1

import (
	"encoding/json"
	var_const "github.com/EDDYCJY/go-gin-example/const"
	"github.com/EDDYCJY/go-gin-example/models"
	"github.com/EDDYCJY/go-gin-example/pkg/app"
	"github.com/EDDYCJY/go-gin-example/pkg/e"
	"github.com/EDDYCJY/go-gin-example/pkg/gredis"
	"github.com/EDDYCJY/go-gin-example/pkg/logging"
	"github.com/EDDYCJY/go-gin-example/service/auth_service"
	"github.com/EDDYCJY/go-gin-example/service/order_service"
	"github.com/EDDYCJY/go-gin-example/service/profit_service"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

// @Summary Get 代练获取已接单列表
// @Produce  json
// @Param user_id body int false "user_id"
// @Param index body int false "index"
// @Param count body int false "count"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/order/takelist [get]
// @Tags 接单
func GetTakerOrders(c *gin.Context) {
	var (
		appG   = app.Gin{C: c}
		userId = com.StrTo(c.Query("user_id")).MustInt()
		index  = com.StrTo(c.Query("index")).MustInt()
		count  = com.StrTo(c.Query("count")).MustInt()
	)
	if userId == 0 || !auth_service.ExistUserInfo(userId) || !auth_service.IsUserTypeInstead(userId) {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	var list []models.Order
	order_service.GetTakeOrderList(userId, &list, index, count)
	appG.Response(http.StatusOK, e.SUCCESS, list)
}

// @Summary Get 用户获取自己已发单列表
// @Produce  json
// @Param user_id body int false "user_id"
// @Param index body int false "index"
// @Param count body int false "count"

// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/order/userlist [get]
// @Tags 接单
func GetUserOrders(c *gin.Context) {
	var (
		appG   = app.Gin{C: c}
		userId = com.StrTo(c.Query("user_id")).MustInt()
		index  = com.StrTo(c.Query("index")).MustInt()
		count  = com.StrTo(c.Query("count")).MustInt()
	)
	if userId == 0 || !auth_service.ExistUserInfo(userId) {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	var list []models.Order
	order_service.GetUserOrderList(userId, &list, index, count)
	appG.Response(http.StatusOK, e.SUCCESS, list)
}

// @Summary Get 用户获取钱包余额
// @Produce  json
// @Param user_id body int false "user_id"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/balance [get]
// @Tags 接单
func GetUserBalance(c *gin.Context) {
	var (
		appG   = app.Gin{C: c}
		userId = com.StrTo(c.Query("user_id")).MustInt()
	)
	if userId == 0 || !auth_service.ExistUserInfo(userId) {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	m := make(map[string]interface{})
	m["balance"] = auth_service.GetUserParam(userId, "balance")
	m["margin"] = auth_service.GetUserParam(userId, "margin")
	appG.Response(http.StatusOK, e.SUCCESS, m)
}

// @Summary Get 用户获取累计订单次数和领取状态
// @Produce  json
// @Param user_id body int false "user_id"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/totalordertimes [get]
// @Tags 用户
func GetUserTotalOrderTimes(c *gin.Context) {
	var (
		appG   = app.Gin{C: c}
		userId = com.StrTo(c.Query("user_id")).MustInt()
	)
	if userId <= 0 || !auth_service.ExistUserInfo(userId) || !profit_service.ExistProfit(userId) {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	m := make(map[string]interface{})
	m["order_total_times"] = profit_service.GetProfitParam(userId, "order_total_times")
	m["order_total_times_status"] = profit_service.GetProfitParam(userId, "order_total_times_status")
	appG.Response(http.StatusOK, e.SUCCESS, m)
}

// @Summary Get 领取累计订单奖励
// @Produce  json
// @Param user_id body int false "user_id"
// @Param award_id body int false "award_id： 1 2 3 4 5"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/orderaward [get]
// @Tags 用户
func GetUserTotalOrderTimesAward(c *gin.Context) {
	var (
		appG    = app.Gin{C: c}
		userId  = com.StrTo(c.Query("user_id")).MustInt()
		awardId = com.StrTo(c.Query("award_id")).MustInt()
	)
	if userId <= 0 || !auth_service.ExistUserInfo(userId) || !profit_service.ExistProfit(userId) {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	orderTotalTimesStatus := profit_service.GetProfitParam(userId, "order_total_times_status")
	orderTotalTimes := profit_service.GetProfitParam(userId, "order_total_times")
	profit := models.Profit{
		UserId: userId,
	}
	m := make(map[string]interface{})
	if awardId == 1 {
		if orderTotalTimesStatus != 0 {
			appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
			return
		}
		if orderTotalTimes < 100 {
			appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
			return
		}
		auth_service.AddUserBalance(userId, 3000, "累计100笔")
		m["order_total_times_status"] = var_const.OrderTotalTimesStatus100
		m["order_total_times"] = 0
	} else if awardId == 2 {
		if orderTotalTimesStatus != var_const.OrderTotalTimesStatus100 {
			appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
			return
		}
		if orderTotalTimes < 500 {
			appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
			return
		}
		auth_service.AddUserBalance(userId, 15000, "累计500笔")
		m["order_total_times_status"] = var_const.OrderTotalTimesStatus500
		m["order_total_times"] = 0
	} else if awardId == 3 {
		if orderTotalTimesStatus != var_const.OrderTotalTimesStatus500 {
			appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
			return
		}
		if orderTotalTimes < 1000 {
			appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
			return
		}
		auth_service.AddUserBalance(userId, 30000, "累计1000笔")
		m["order_total_times_status"] = var_const.OrderTotalTimesStatus1000
		m["order_total_times"] = 0
	} else if awardId == 4 {
		if orderTotalTimesStatus != var_const.OrderTotalTimesStatus1000 {
			appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
			return
		}
		if orderTotalTimes < 2000 {
			appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
			return
		}
		auth_service.AddUserBalance(userId, 60000, "累计2000笔")
		m["order_total_times_status"] = var_const.OrderTotalTimesStatus2000
		m["order_total_times"] = 0
	} else if awardId == 5 {
		if orderTotalTimesStatus != var_const.OrderTotalTimesStatus2000 {
			appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
			return
		}
		if orderTotalTimes < 10000 {
			appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
			return
		}
		auth_service.AddUserBalance(userId, 300000, "累计10000笔")
		m["order_total_times"] = 0
		m["order_total_times_status"] = var_const.OrderTotalTimesStatus2000
	} else {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	logging.Info("GetUserTotalOrderTimesAward user_id:" + c.Query("user_id") + " award_id:" + c.Query("award_id"))
	if !profit.Updates(m) {
		appG.Response(http.StatusBadRequest, e.ERROR, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, m)
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
	appG := app.Gin{C: c}
	userId := com.StrTo(c.PostForm("user_id")).MustInt()
	agentId := com.StrTo(c.PostForm("agent_id")).MustInt()
	if userId == agentId {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	if !auth_service.ExistUserInfo(userId) || !auth_service.ExistUserInfo(agentId) {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	if auth_service.GetUserParam(userId, "agent_id") != 0 {
		appG.Response(http.StatusBadRequest, e.ERROR, nil)
		return
	}
	agentParentId := auth_service.GetUserParam(agentId, "agent_id")
	user := models.User{
		UserId: userId,
	}
	m := make(map[string]interface{})
	m["agent_id"] = agentId
	if agentParentId > 0 {
		m["agent_parent_id"] = agentParentId
	}
	if !user.Updates(m) {
		logging.Error("BindAgent-db: user_id-" + c.PostForm("user_id") + ",agent_id-" + c.PostForm("agent_id"))
		appG.Response(http.StatusBadRequest, e.ERROR, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
	return
}

type GetQRcodeInfo struct {
	Scene string `json:"scene" validate:"required"`
	Page  string `json:"page"`
	Width int    `json:"width"`
}

// @Summary 获取二维码
// @Produce  json
// @Param scene body string false "scene"
// @Param page body string false "page"
// @Param width body int false "width"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/qrcode [post]
// @Tags 二维码
func QRcodeGet(c *gin.Context) {
	var (
		appG  = app.Gin{C: c}
		scene = c.PostForm("scene")
		page  = c.PostForm("page")
		width = com.StrTo(c.Query("width")).MustInt()
	)
	var accessToken string
	if accessToken = auth_service.GetAccessToken(); accessToken == "" {
		token, retErr := auth_service.UpdateAccessToken()
		if retErr != 0 {
			appG.Response(http.StatusBadRequest, e.ERROR, nil)
			return
		}
		accessToken = token
	}
	url := "https://api.weixin.qq.com/wxa/getwxacodeunlimit?access_token=" + accessToken //请求地址
	contentType := "application/json"
	//参数，多个用&隔开

	getQRcodeInfo := GetQRcodeInfo{
		Scene: scene,
		Page:  page,
		Width: width,
	}
	jsonData, e2 := json.Marshal(getQRcodeInfo)
	if e2 != nil {
		appG.Response(http.StatusBadRequest, e.ERROR, nil)
		return
	}
	data := strings.NewReader(string(jsonData))
	resp, err2 := http.Post(url, contentType, data)
	if err2 != nil {
		appG.Response(http.StatusBadRequest, e.ERROR, nil)
		return
	}
	defer resp.Body.Close()
	body, err3 := ioutil.ReadAll(resp.Body)
	if err3 == nil {
		var dat map[string]interface{}
		if err := json.Unmarshal([]byte(string(body)), &dat); err == nil {
			appG.Response(http.StatusBadRequest, e.ERROR, nil)
			return
		} else {
			appG.Response(http.StatusOK, e.SUCCESS, body)
			return
		}

	} else {
		appG.Response(http.StatusBadRequest, e.ERROR, nil)
		return
	}
}

// @Summary Get 获取昨天团队发单接单收益
// @Produce  json
// @Param user_id body int false "user_id"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/agent/profit [get]
// @Tags 用户信息
func GetAgentProfit(c *gin.Context) {
	var (
		appG   = app.Gin{C: c}
		userId = com.StrTo(c.Query("user_id")).MustInt()
	)
	if userId == 0 || !auth_service.ExistUserInfo(userId) {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	pProfit := profit_service.GetProfitParam(userId, "order_yesterday_publish_profit")
	tProfit := profit_service.GetProfitParam(userId, "order_yesterday_taker_profit")

	keys := gredis.GetKeys("game_user:")
	for _, v := range keys {
		fields := []string{"user_id", "agent_id"}
		data, err := gredis.HMGet(v, fields...)
		if err != nil || len(data) != 2 || data[0] == nil || data[1] == nil {
			continue
		}
		agentId, err := strconv.Atoi(string(data[1].([]byte)))
		if err != nil {
			continue
		}
		if agentId == userId {
			childId, err := strconv.Atoi(string(data[0].([]byte)))
			if err != nil {
				continue
			}
			pProfit += profit_service.GetProfitParam(childId, "order_yesterday_agent_publish_profit")
			tProfit += profit_service.GetProfitParam(childId, "order_yesterday_agent_taker_profit")
		}
	}

	m := make(map[string]interface{})
	m["pProfit"] = pProfit
	m["tProfit"] = tProfit
	appG.Response(http.StatusOK, e.SUCCESS, m)
}

// @Summary 获取用户的订单消息/留言
// @Produce  json
// @Param user_id body int false "user_id"
// @Param index body int false "index"
// @Param count body int false "count"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/user/getmessage [post]
// @Tags 订单留言
func GetUserMessage(c *gin.Context) {
	var (
		appG   = app.Gin{C: c}
		userId = com.StrTo(c.PostForm("user_id")).MustInt()
		index  = com.StrTo(c.PostForm("index")).MustInt()
		count  = com.StrTo(c.PostForm("count")).MustInt()
	)
	if userId == 0 || !auth_service.ExistUserInfo(userId) {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	res, _ := gredis.LRange(order_service.GetRedisKeyMessageUser(userId), index, index+count)
	_, _ = gredis.HSet(order_service.GetRedisKeyMessageNoRead(), c.PostForm("user_id"), "0")
	appG.Response(http.StatusOK, e.SUCCESS, res)
	return
}

// @Summary 是否有未读订单消息/留言
// @Produce  json
// @Param user_id body int false "user_id"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/user/msgnoread [post]
// @Tags 订单留言
func GetUserMessageNoRead(c *gin.Context) {
	var (
		appG   = app.Gin{C: c}
		userId = com.StrTo(c.PostForm("user_id")).MustInt()
	)
	if userId == 0 || !auth_service.ExistUserInfo(userId) {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	res, _ := gredis.HGet(order_service.GetRedisKeyMessageNoRead(), c.PostForm("user_id"))
	appG.Response(http.StatusOK, e.SUCCESS, res)
	return
}

// @Summary Get 客服获取用户列表
// @Produce  json
// @Param index body int false "index"
// @Param count body int false "count"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /admin/user/all [post]
// @Tags 客服
func GetAdminUserList(c *gin.Context) {
	appG := app.Gin{C: c}
	index := com.StrTo(c.PostForm("index")).MustInt()
	count := com.StrTo(c.PostForm("count")).MustInt()
	var list []models.User
	auth_service.GetUserList(&list, "", index, count)
	m := make(map[string]interface{})
	m["list"] = list
	m["count"] = len(list)
	appG.Response(http.StatusOK, e.SUCCESS, m)
}

// @Summary 客服 设置/取消 用户发布订单权限
// @Produce  json
// @Param user_id body int false "user_id"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /admin/user/canpublish [post]
// @Tags 客服
func AdminSetUserCanPublish(c *gin.Context) {
	var (
		appG   = app.Gin{C: c}
		userId = com.StrTo(c.PostForm("user_id")).MustInt()
	)
	if userId == 0 || !auth_service.ExistUserInfo(userId) {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	canPublish := auth_service.GetUserParam(userId, "can_publish")
	user := models.User{
		UserId: userId,
	}
	if canPublish == 1 {
		canPublish = 0
	} else {
		canPublish = 1
	}
	m := make(map[string]interface{})
	m["can_publish"] = canPublish
	if !user.Updates(m) {
		logging.Error("AdminSetUserCanPublish-db: user_id-" + c.PostForm("user_id"))
		appG.Response(http.StatusBadRequest, e.ERROR, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
	return
}
