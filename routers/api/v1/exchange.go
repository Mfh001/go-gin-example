package v1

import (
	"encoding/json"
	var_const "github.com/EDDYCJY/go-gin-example/const"
	"github.com/EDDYCJY/go-gin-example/models"
	"github.com/EDDYCJY/go-gin-example/pkg/app"
	"github.com/EDDYCJY/go-gin-example/pkg/e"
	"github.com/EDDYCJY/go-gin-example/pkg/logging"
	"github.com/EDDYCJY/go-gin-example/service/auth_service"
	"github.com/EDDYCJY/go-gin-example/service/bank_service"
	"github.com/EDDYCJY/go-gin-example/service/exchange_service"
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
	if !auth_service.ExistUserInfo(form.UserId) || !auth_service.IsUserTypeInstead(form.UserId) {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	//TODO 是否绑定银行卡
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
	_, _ = models.GetNeedExchanges(&list, index, count)
	appG.Response(http.StatusOK, e.SUCCESS, list)
}

// @Summary Get 管理员审核提现
// @Produce  json
// @Param id body int false "提现id"
// @Param state body int false "state ：-1拒绝 1通过"
// @Param remarks body string false "备注"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /exchange/check [post]
// @Tags 提现
func ExchangeCheck(c *gin.Context) {
	appG := app.Gin{C: c}
	id := com.StrTo(c.PostForm("id")).MustInt()
	state := com.StrTo(c.PostForm("state")).MustInt()
	if state != -1 && state != 1 {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	if !exchange_service.ExistExchange(id) {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	if exchange_service.GetExchangeParam(id, "status") != 0 {
		appG.Response(http.StatusBadRequest, e.ERROR, nil)
		return
	}
	exchange := models.Exchange{
		Id: id,
	}

	m := make(map[string]interface{})
	m["status"] = state
	m["upd_time"] = int(time.Now().Unix())

	if !exchange.Updates(m) {
		logging.Error("ExchangeCheck-db: id-" + c.PostForm("id") + ",state-" + c.PostForm("state"))
		appG.Response(http.StatusBadRequest, e.ERROR, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

// @Summary Get 获取银行卡
// @Produce  json
// @Param user_id body int false "user_id"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /exchange/bank [get]
// @Tags 提现
func GetExchangeBank(c *gin.Context) {
	appG := app.Gin{C: c}
	userId := com.StrTo(c.Query("user_id")).MustInt()

	if !bank_service.ExistBank(userId) {
		appG.Response(http.StatusBadRequest, e.ERROR, nil)
		return
	}
	m := make(map[string]interface{})
	m["bank_name"] = bank_service.GetBankParamString(userId, "bank_name")
	m["Bank_branch_name"] = bank_service.GetBankParamString(userId, "Bank_branch_name")
	m["bank_card"] = bank_service.GetBankParamString(userId, "bank_card")
	m["user_name"] = bank_service.GetBankParamString(userId, "user_name")
	appG.Response(http.StatusOK, e.SUCCESS, m)
}
