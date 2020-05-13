package v1

import (
	"github.com/EDDYCJY/go-gin-example/models"
	"github.com/EDDYCJY/go-gin-example/pkg/app"
	"github.com/EDDYCJY/go-gin-example/pkg/e"
	"github.com/EDDYCJY/go-gin-example/pkg/gredis"
	"github.com/EDDYCJY/go-gin-example/service/auth_service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// @Summary 绑定银行卡
// @Produce  json
// @Param   json     body    models.RequestBankCardInfo     true        "请求的json结构"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/bink/bind [post]
// @Tags 银行卡
func BindBankCard(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form models.RequestBankCardInfo
	)
	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}
	if form.BankCardInfo.UserId == 0 || !auth_service.ExistUserInfo(form.BankCardInfo.UserId) {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	userPhone, err := auth_service.GetUserPhone(form.BankCardInfo.UserId)
	if userPhone == "" || err != nil {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	code, err := gredis.Get(auth_service.GetRedisKeySmsCode(userPhone))
	if err != nil {
		appG.Response(http.StatusBadRequest, e.ERROR, nil)
		return
	}

	if code != form.Code { //验证码无效 请重新获取
		appG.Response(http.StatusBadRequest, e.SMSCODE_ERROR, nil)
		return
	}

	form.BankCardInfo.Insert()
	appG.Response(http.StatusOK, e.SUCCESS, nil)
	return
}

// @Summary 获取银行卡信息
// @Accept  application/json; charset=utf-8
// @Param user_id body int false "user_id"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/bank [get]
// @Tags 银行卡
func GetBankCardInfo(c *gin.Context) {
	appG := app.Gin{C: c}
	userId, err := strconv.Atoi(c.Query("user_id"))
	if err != nil {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	if userId == 0 || !auth_service.ExistUserInfo(userId) {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	//shop := model.ShopInfo{
	//	ShopId: shopId,
	//}
	//phone, _ := shop.FindPhoneByShopId()
	//code, err2 := db_redis.GetRedisDB(config.RedisSessionDb).Get("2" + phone)
	//if err2 != nil {
	//	server_handler.RequestCallBack(context, 402, "系统繁忙", nil)
	//	return
	//}
	//if code != reqCode { //验证码无效 请重新获取
	//	server_handler.RequestCallBack(context, 403, "验证码无效 请重新获取", nil)
	//	return
	//}
	//card := model.BankCardInfo{
	//	ShopId: shopId,
	//}
	//card.Password = utils.MD5(password)
	//if !card.CheckBankCardPassword() {
	//	server_handler.RequestCallBack(context, 404, "密码错误", nil)
	//	return
	//}

	bankInfo := models.BankCardInfo{
		UserId: userId,
	}

	if _, err := bankInfo.FindBankCardInfo(); err != nil {
		appG.Response(http.StatusBadRequest, e.ERROR, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, bankInfo)
	return
}
