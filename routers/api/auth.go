package api

import (
	"github.com/EDDYCJY/go-gin-example/middleware/bloom_filter"
	"github.com/EDDYCJY/go-gin-example/pkg/setting"
	"net/http"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"

	"github.com/EDDYCJY/go-gin-example/pkg/app"
	"github.com/EDDYCJY/go-gin-example/pkg/e"
	"github.com/EDDYCJY/go-gin-example/pkg/util"
	"github.com/EDDYCJY/go-gin-example/service/auth_service"
)

// @Summary 微信登陆接口 发送code获取session_key
// @Param code body string false "code"
// @Param nickname body string false "nickname"
// @Param avatar_url body string false "avatar_url"
// @Param gender body int false "gender"
// @Param province body string false "province"
// @Param city body string false "city"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /wxlogin [post]
// @Tags 登陆
func WXLogin(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form auth_service.WxLoginUserInfo
	)
	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}
	data := gin.H{}
	if sessionKey, b := form.WXLogin(); sessionKey != "" && b {
		data["session_key"] = sessionKey
		appG.Response(http.StatusOK, e.SUCCESS, data)
		return
	}
	appG.Response(http.StatusBadRequest, e.ERROR, nil)
	return
}

type login struct {
	SessionKey string `valid:"Required; Length(32)"`
}

// @Summary 微信登陆接口 发送code获取session_key
// @Param session_key body string false "sessionKey"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /login [post]
// @Tags 登陆
func Login(c *gin.Context) {
	appG := app.Gin{C: c}
	valid := validation.Validation{}

	sessionKey := c.PostForm("session_key")

	a := login{SessionKey: sessionKey}
	ok, _ := valid.Valid(&a)
	if !ok {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	info, err := auth_service.GetWXCode(sessionKey)
	ok, _ = valid.Valid(info)
	if !ok {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	if !bloom_filter.Filter.Has(auth_service.GetRedisKeyUserInfo(info.UserId)) {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	token, err := util.GenerateToken(sessionKey, sessionKey)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_AUTH_TOKEN, nil)
		return
	}
	if !auth_service.ExistUserInfo(info.UserId) {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	data := gin.H{}
	userInfo, err := auth_service.GetUserInfo(info.UserId)
	if err != nil {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return

	}
	data["token"] = token
	data["user_info"] = userInfo
	data["platform"] = setting.Platform
	appG.Response(http.StatusOK, e.SUCCESS, data)
	return
}
