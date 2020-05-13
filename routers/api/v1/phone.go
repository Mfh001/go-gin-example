package v1

import (
	var_const "github.com/EDDYCJY/go-gin-example/const"
	"github.com/EDDYCJY/go-gin-example/models"
	"github.com/EDDYCJY/go-gin-example/pkg/app"
	"github.com/EDDYCJY/go-gin-example/pkg/e"
	"github.com/EDDYCJY/go-gin-example/pkg/gredis"
	"github.com/EDDYCJY/go-gin-example/pkg/util"
	"github.com/EDDYCJY/go-gin-example/service/auth_service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// @Summary 绑定手机号 确认绑定的接口
// @Produce  json
// @Param user_id body int false "user_id"
// @Param   phone     body    string     false        "phone"
// @Param   code     body    string     false        "code"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/phone/bind [post]
// @Tags 绑定手机号
func BindPhone(c *gin.Context) {
	appG := app.Gin{C: c}

	userId, err := strconv.Atoi(c.PostForm("user_id"))
	if err != nil {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	if !auth_service.ExistUserInfo(userId) {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	phone := c.PostForm("phone")
	if !util.IsPhoneNum(phone) {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	newCode := c.PostForm("code")
	if newCode == "" {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	//判断手机号是否已存在
	userPhone, err := auth_service.GetUserPhone(userId)
	if userPhone != "" || err != nil {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	user := models.User{
		Phone: phone,
	}
	usId, err := user.FindPhone()
	if err != nil || usId != -1 {
		appG.Response(http.StatusBadRequest, e.ERROR_EXIST_PHONE, nil)
		return
	}

	code, err := gredis.Get(auth_service.GetRedisKeySmsCode(phone))
	if err != nil {
		appG.Response(http.StatusBadRequest, e.ERROR, nil)
		return
	}
	if code != newCode { //验证码无效 请重新获取
		appG.Response(http.StatusBadRequest, e.SMSCODE_ERROR, nil)
		return
	} else { //验证通过
		user := models.User{
			UserId: userId,
		}
		dbInfo := make(map[string]interface{})
		dbInfo["phone"] = phone
		if !user.Updates(dbInfo) {
			appG.Response(http.StatusBadRequest, e.ERROR, nil)
			return
		}
		appG.Response(http.StatusOK, e.SUCCESS, nil)
		return
	}
}

// @Summary 获取要绑定的手机号的验证码 绑定手机时使用
// @Accept  application/json; charset=utf-8
// @Param user_id body int false "user_id"
// @Param   phone     body    string     false        "phone"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/phone/code [get]
// @Tags 绑定手机号
func GetPhoneRegCode(c *gin.Context) {
	appG := app.Gin{C: c}
	phone := c.Query("phone")
	if !util.IsPhoneNum(phone) {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	userId, err := strconv.Atoi(c.Query("user_id"))
	if err != nil {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	if userId == 0 || !auth_service.ExistUserInfo(userId) {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	//判断手机号是否已存在
	userPhone, err := auth_service.GetUserPhone(userId)
	if userPhone != "" || err != nil {
		appG.Response(http.StatusBadRequest, e.ERROR_EXIST_PHONE, nil)
		return
	}
	user := models.User{
		Phone: phone,
	}
	usId, err := user.FindPhone()
	if err != nil || usId != -1 {
		appG.Response(http.StatusBadRequest, e.ERROR_EXIST_PHONE, nil)
		return
	}

	code, err := gredis.Get(auth_service.GetRedisKeySmsCode(phone))
	//if err != nil {
	//	appG.Response(http.StatusBadRequest, e.ERROR, nil)
	//	return
	//}
	if code == "" { //从api重新获取验证码
		newCode := util.RandomStringNoLetter(6)
		_ = gredis.Set(auth_service.GetRedisKeySmsCode(phone), newCode, var_const.SMSCodeExpireTime)

		util.SendSMSCode(phone, newCode)
		appG.Response(http.StatusOK, e.SUCCESS, nil)
		return
	} else { //已获取到验证码 并且未过期
		_ = gredis.Set(auth_service.GetRedisKeySmsCode(phone), code, var_const.SMSCodeExpireTime)
		util.SendSMSCode(phone, code)
		appG.Response(http.StatusOK, e.SUCCESS, nil)
		return
	}
}

// @Summary 获取已绑定手机号的验证码
// @Accept  application/json; charset=utf-8
// @Param user_id body int false "user_id"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/phone/code2 [get]
// @Tags 手机验证码
func GetPhoneCodeNoPhone(c *gin.Context) {
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
	//判断手机号是否已存在
	userPhone, err := auth_service.GetUserPhone(userId)
	if userPhone == "" || err != nil {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	if !util.IsPhoneNum(userPhone) {
		appG.Response(http.StatusBadRequest, e.ERROR, nil)
		return
	}
	code, err := gredis.Get(auth_service.GetRedisKeySmsCode(userPhone))
	//if err != nil {
	//	appG.Response(http.StatusBadRequest, e.ERROR, nil)
	//	return
	//}
	if code == "" { //从api重新获取验证码
		newCode := util.RandomStringNoLetter(6)
		_ = gredis.Set(auth_service.GetRedisKeySmsCode(userPhone), newCode, var_const.SMSCodeExpireTime)

		util.SendSMSCode(userPhone, newCode)
		appG.Response(http.StatusOK, e.SUCCESS, nil)
		return
	} else { //已获取到验证码 并且未过期
		_ = gredis.Set(auth_service.GetRedisKeySmsCode(userPhone), code, var_const.SMSCodeExpireTime)
		util.SendSMSCode(userPhone, code)
		appG.Response(http.StatusOK, e.SUCCESS, nil)
		return
	}
}
