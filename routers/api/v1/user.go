package v1

import (
	"encoding/json"
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
