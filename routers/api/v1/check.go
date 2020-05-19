package v1

import (
	"encoding/json"
	var_const "github.com/EDDYCJY/go-gin-example/const"
	"github.com/EDDYCJY/go-gin-example/models"
	"github.com/EDDYCJY/go-gin-example/pkg/app"
	"github.com/EDDYCJY/go-gin-example/pkg/e"
	"github.com/EDDYCJY/go-gin-example/pkg/logging"
	"github.com/EDDYCJY/go-gin-example/service/auth_service"
	"github.com/EDDYCJY/go-gin-example/service/check_service"
	"github.com/EDDYCJY/go-gin-example/service/order_service"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"net/http"
	"time"
)

// @Summary 代练提交或更新段位审核
// @Produce  json
// @Param user_id body int false "user_id"
// @Param game_id body string false "game_id"
// @Param game_server body int false "game_server"
// @Param game_pos body int false "game_pos"
// @Param game_level body string false "game_level"
// @Param img_url body string false "img_url"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/check [post]
// @Tags 审核
func AddCheck(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form models.Check
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
	form.RegTime = int(time.Now().Unix())
	if !form.Save() {
		log, _ := json.Marshal(form)
		logging.Error("AddCheck:form-" + string(log))
		appG.Response(http.StatusBadRequest, e.ERROR, nil)
		return
	}
	//更新用户审核信息DB
	userInfo := models.User{
		UserId: form.UserId,
	}
	var dbInfo = make(map[string]interface{})
	dbInfo["check_pass"] = var_const.CheckNeed
	if !userInfo.Updates(dbInfo) {
		logInfo, _ := json.Marshal(dbInfo)
		logging.Error("AddCheck:db-check-Updates" + string(logInfo))
		appG.Response(http.StatusBadRequest, e.ERROR, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
	return
}

// @Summary 用户获取提交的审核信息
// @Produce  json
// @Param user_id body int false "user_id"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/check [get]
// @Tags 审核
func GetCheckInfo(c *gin.Context) {
	var (
		appG   = app.Gin{C: c}
		userId = com.StrTo(c.Query("user_id")).MustInt()
	)
	if userId == 0 || !auth_service.ExistUserInfo(userId) || !check_service.ExistUserCheck(userId) {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	data, err := check_service.GetUserCheckInfo(userId)
	if err != nil {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, data)
}

// @Summary Get 管理员获取审核列表
// @Produce  json
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /check/admin [get]
// @Tags 审核
func GetAdminChecks(c *gin.Context) {
	appG := app.Gin{C: c}
	var list []models.Check
	order_service.Refund(42)
	check_service.GetCheckList(&list)
	appG.Response(http.StatusOK, e.SUCCESS, list)
}

// @Summary 管理员进行审核
// @Produce  json
// @Param state body int false "State -1/1"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /check/admin/{user_id} [put]
// @Tags 审核
func AdminCheck(c *gin.Context) {
	var (
		appG   = app.Gin{C: c}
		userId = com.StrTo(c.Param("id")).MustInt()
		state  = com.StrTo(c.Query("state")).MustInt()
	)
	if userId == 0 || !auth_service.ExistUserInfo(userId) || !check_service.ExistUserCheck(userId) {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	if state != var_const.CheckRefuse && state != var_const.CheckPass {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	userState, err := auth_service.GetUserCheckPassState(userId)
	if err != nil {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	if userState == var_const.CheckPass {
		appG.Response(http.StatusBadRequest, e.ERROR_CHECK_PASSED, nil)
		return
	}

	//更新用户审核信息DB
	userInfo := models.User{
		UserId: userId,
	}
	var dbInfo = make(map[string]interface{})
	dbInfo["check_pass"] = state
	if state == var_const.CheckPass {
		data, err := check_service.GetUserCheckInfo(userId)
		if err != nil {
			appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
			return
		}
		dbInfo["type"] = var_const.UserTypeInstead
		dbInfo["game_id"] = data["game_id"]
		dbInfo["game_server"] = data["game_server"]
		dbInfo["game_pos"] = data["game_pos"]
		dbInfo["game_level"] = data["game_level"]
		dbInfo["img_url"] = data["img_url"]
	}
	if !userInfo.Updates(dbInfo) {
		logInfo, _ := json.Marshal(dbInfo)
		logging.Error("AddCheck:db-check-Updates" + string(logInfo))
		appG.Response(http.StatusBadRequest, e.ERROR, nil)
		return
	}

	//删除审核数据
	checkInfo := models.Check{
		UserId: userId,
	}
	if !checkInfo.Delete() {
		appG.Response(http.StatusBadRequest, e.ERROR, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
	return
}
