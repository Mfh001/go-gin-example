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
	"net/http"
	"time"
)

// @Summary 提现申请
// @Produce  json
// @Param user_id body int false "user_id"
// @Param money body string false "money"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/check [post]
// @Tags 提现
func AddExchange(c *gin.Context) {
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
