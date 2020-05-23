package v1

import (
	"github.com/EDDYCJY/go-gin-example/models"
	"github.com/EDDYCJY/go-gin-example/pkg/app"
	"github.com/EDDYCJY/go-gin-example/pkg/e"
	"github.com/EDDYCJY/go-gin-example/pkg/logging"
	"github.com/EDDYCJY/go-gin-example/service/auth_service"
	"github.com/EDDYCJY/go-gin-example/service/order_service"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"net/http"
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
	if userId == 0 || !auth_service.ExistUserInfo(userId) {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	//必须是代练 TODO
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
