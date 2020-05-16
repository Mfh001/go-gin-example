package v1

import (
	"github.com/EDDYCJY/go-gin-example/models"
	"github.com/EDDYCJY/go-gin-example/pkg/app"
	"github.com/EDDYCJY/go-gin-example/pkg/e"
	"github.com/EDDYCJY/go-gin-example/service/auth_service"
	"github.com/EDDYCJY/go-gin-example/service/order_service"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"net/http"
)

// @Summary Get 代练获取已接单列表
// @Produce  json
// @Param user_id body int false "user_id"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/order/takelist [get]
// @Tags 接单
func GetTakerOrders(c *gin.Context) {
	var (
		appG   = app.Gin{C: c}
		userId = com.StrTo(c.Query("user_id")).MustInt()
	)
	if userId == 0 || !auth_service.ExistUserInfo(userId) {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	//必须是代练 TODO
	var list []models.Order
	order_service.GetTakeOrderList(userId, &list)
	appG.Response(http.StatusOK, e.SUCCESS, list)
}

// @Summary Get 用户获取自己已发单列表
// @Produce  json
// @Param user_id body int false "user_id"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/order/userlist [get]
// @Tags 接单
func GetUserOrders(c *gin.Context) {
	var (
		appG   = app.Gin{C: c}
		userId = com.StrTo(c.Query("user_id")).MustInt()
	)
	if userId == 0 || !auth_service.ExistUserInfo(userId) {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	var list []models.Order
	order_service.GetUserOrderList(userId, &list)
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
	m["balance"], _ = auth_service.GetUserBalance(userId)
	m["margin"], _ = auth_service.GetUserMargin(userId)
	appG.Response(http.StatusOK, e.SUCCESS, m)
}
