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
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"net/http"
)

// @Summary 代练提交或更新段位审核
// @Produce  json
// @Param user_id body int false "user_id"
// @Param game_id body string false "game_id"
// @Param game_server body int false "game_server"
// @Param game_pos body int false "game_pos"
// @Param game_level body int false "game_level"
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

// @Summary Get 管理员获取审核列表
// @Produce  json
// @Param token query string false "token"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/admin/check [get]
// @Tags 审核
func GetAdminChecks(c *gin.Context) {
	appG := app.Gin{C: c}
	var list []models.Check
	check_service.GetCheckList(&list)
	appG.Response(http.StatusOK, e.SUCCESS, list)
}

// @Summary 管理员进行审核
// @Produce  json
// @Param user_id path int true "user_id"
// @Param state body int false "State -1/1"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/admin/check/{user_id} [put]
// @Tags 审核
func AdminCheck(c *gin.Context) {
	var (
		appG   = app.Gin{C: c}
		userId = com.StrTo(c.Param("user_id")).MustInt()
		state  = com.StrTo(c.Param("state")).MustInt()
	)
	if userId == 0 || !auth_service.ExistUserInfo(userId) {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	if state != var_const.CheckRefuse && state != var_const.CheckPass {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	//httpCode, errCode := app.BindAndValid(c, &form)
	//if errCode != e.SUCCESS {
	//	appG.Response(httpCode, errCode, nil)
	//	return
	//}
	//
	//tagService := tag_service.Tag{
	//	ID:         form.ID,
	//	Name:       form.Name,
	//	ModifiedBy: form.ModifiedBy,
	//	State:      form.State,
	//}
	//
	//exists, err := tagService.ExistByID()
	//if err != nil {
	//	appG.Response(http.StatusInternalServerError, e.ERROR_EXIST_TAG_FAIL, nil)
	//	return
	//}
	//
	//if !exists {
	//	appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_TAG, nil)
	//	return
	//}
	//
	//err = tagService.Edit()
	//if err != nil {
	//	appG.Response(http.StatusInternalServerError, e.ERROR_EDIT_TAG_FAIL, nil)
	//	return
	//}
	//
	//appG.Response(http.StatusOK, e.SUCCESS, nil)
}
