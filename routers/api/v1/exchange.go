package v1

import (
	"encoding/json"
	var_const "github.com/EDDYCJY/go-gin-example/const"
	"github.com/EDDYCJY/go-gin-example/models"
	"github.com/EDDYCJY/go-gin-example/pkg/app"
	"github.com/EDDYCJY/go-gin-example/pkg/e"
	"github.com/EDDYCJY/go-gin-example/pkg/logging"
	"github.com/EDDYCJY/go-gin-example/service/auth_service"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"net/http"
	"time"
)

// @Summary 提现申请
// @Produce  json
// @Param user_id body int false "user_id"
// @Param money body string false "money"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/exchange [post]
// @Tags 提现
func AddExchange(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form models.Exchange
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
	if auth_service.GetUserParam(form.UserId, "type") != var_const.UserTypeInstead {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	balance := auth_service.GetUserParam(form.UserId, "balance")
	if balance < form.Money || form.Money < var_const.ExchangeMinMoney {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	if !auth_service.RemoveUserBalance(form.UserId, form.Money) {
		appG.Response(http.StatusBadRequest, e.ERROR, nil)
		return
	}
	form.RealMoney = form.Money * (100 - var_const.ExchangeRate) / 100
	form.Rate = var_const.ExchangeRate
	form.NickName = auth_service.GetUserParamString(form.UserId, "nick_name")
	form.RegTime = int(time.Now().Unix())
	if !form.Insert() {
		log, _ := json.Marshal(form)
		logging.Error("AddExchange:form-" + string(log))
		appG.Response(http.StatusBadRequest, e.ERROR, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
	return
}

// @Summary Get 管理员获取提现审核列表
// @Produce  json
// @Param index body int false "index"
// @Param count body int false "count"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /exchange/all [get]
// @Tags 提现
func GetAdminExchanges(c *gin.Context) {
	appG := app.Gin{C: c}
	index := com.StrTo(c.Query("index")).MustInt()
	count := com.StrTo(c.Query("count")).MustInt()
	var list []models.Exchange
	models.GetNeedExchanges(&list, index, count)
	appG.Response(http.StatusOK, e.SUCCESS, list)
}
