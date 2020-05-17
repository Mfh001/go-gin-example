package v1

import (
	var_const "github.com/EDDYCJY/go-gin-example/const"
	"github.com/EDDYCJY/go-gin-example/models"
	"github.com/EDDYCJY/go-gin-example/pkg/app"
	"github.com/EDDYCJY/go-gin-example/pkg/e"
	"github.com/EDDYCJY/go-gin-example/pkg/setting"
	"github.com/EDDYCJY/go-gin-example/service/auth_service"
	"github.com/EDDYCJY/go-gin-example/service/order_service"
	"github.com/EDDYCJY/go-gin-example/service/team_service"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"net/http"
)

// @Summary 下单
// @Produce  json
// @Param owner_id body int false "owner_id"
// @Param game_type body int false "游戏"
// @Param big_zone body int false "游戏大区"
// @Param cur_level body int false "起始段位"
// @Param target_level body int false "结束段位"
// @Param need_num body int false "代练发布最多2人，用户发布只能1人"
// @Param need_pwd body int false "代练发布可以设置是否需要密码"
// @Param pwd body string false "密码"
// @Param contact body string false "contact"
// @Param qq body string false "微信"
// @Param description body string false "description"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/team [post]
// @Tags 车队
func AddTeam(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form models.Team
	)
	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}
	if !auth_service.ExistUserInfo(form.OwnerId) {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	if form.CurLevel >= form.TargetLevel || form.TargetLevel <= 0 {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	ownType, err := auth_service.GetUserType(form.OwnerId)
	if err != nil {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	if ownType != var_const.UserTypeNormal {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	form.OwnerType = ownType
	nickName, _ := auth_service.GetUserNickName(form.OwnerId)
	form.NickName = nickName
	teamId, err := team_service.IncrTeamId()
	if err != nil {
		appG.Response(http.StatusBadRequest, e.ERROR, nil)
		return
	}
	form.TeamId = teamId
	if ownType == var_const.UserTypeNormal {
		form.NeedNum = 1
		form.Num = 1
		form.UserId1 = form.OwnerId
		form.NickName1 = nickName
		//创建order订单
		order := models.Order{
			CurLevel:    form.CurLevel,
			TargetLevel: form.TargetLevel,
			UserId:      form.OwnerId,
			GameType:    form.GameType,
			GameZone:    form.BigZone,
			Contact:     form.Contact,
			Qq:          form.Qq,
			Description: form.Description,
		}
		if !order_service.CreateOrder(&order, form.TeamId) {
			appG.Response(http.StatusBadRequest, e.ERROR, nil)
			return
		}
		form.OrderId1 = order.OrderId
		form.Price = order.Price
		if !team_service.CreateTeam(&form) {
			appG.Response(http.StatusBadRequest, e.ERROR, nil)
			return
		}
		data := gin.H{}
		data["order_id"] = form.OrderId1
		appG.Response(http.StatusOK, e.SUCCESS, data)
		return
	} else {
		if !team_service.CreateTeam(&form) {
			appG.Response(http.StatusBadRequest, e.ERROR, nil)
			return
		}
		data := gin.H{}
		data["team_id"] = form.TeamId
		appG.Response(http.StatusOK, e.SUCCESS, data)
		return
	}

}

// @Summary 验证密码 用户是否可以加入车队
// @Produce  json
// @Param team_id body int false "team_id"
// @Param pwd body string false "pwd"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/team/joincheck [post]
// @Tags 订单
func JoinTeamCheckPwd(c *gin.Context) {
	var (
		appG   = app.Gin{C: c}
		teamId = com.StrTo(c.PostForm("team_id")).MustInt()
		pwd    = c.PostForm("pwd")
	)
	if teamId <= 0 {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	if team_service.GetTeamParam(teamId, "need_pwd") == 0 {
		appG.Response(http.StatusOK, e.SUCCESS, nil)
		return
	} else {
		if team_service.GetTeamParamString(teamId, "pwd") == pwd {
			appG.Response(http.StatusOK, e.SUCCESS, nil)
			return
		}
	}
	appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
	return
}

// @Summary 验证密码 用户是否可以加入车队
// @Produce  json
// @Param user_id body int false "user_id"
// @Param team_id body int false "team_id"
// @Param cur_level body int false "起始段位"
// @Param target_level body int false "结束段位"
// @Param contact body string false "contact"
// @Param qq body string false "微信"
// @Param description body string false "description"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/team/join [post]
// @Tags 订单
func JoinTeam(c *gin.Context) {
	var (
		appG        = app.Gin{C: c}
		userId      = com.StrTo(c.PostForm("user_id")).MustInt()
		teamId      = com.StrTo(c.PostForm("team_id")).MustInt()
		curLevel    = com.StrTo(c.PostForm("cur_level")).MustInt()
		targetLevel = com.StrTo(c.PostForm("target_level")).MustInt()
		contact     = c.PostForm("contact")
		qq          = c.PostForm("qq")
		description = c.PostForm("description")
	)
	if teamId <= 0 {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	if team_service.GetTeamParam(teamId, "cur_level") > curLevel {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	if team_service.GetTeamParam(teamId, "target_level") < targetLevel {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	//不能多于5颗星
	count := 0
	for i := 0; i < len(setting.PlatFormLevelAll); i++ {
		if setting.PlatFormLevelAll[i].Idx > curLevel && setting.PlatFormLevelAll[i].Idx <= targetLevel {
			count++
		}
	}
	if count > 5 {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	if team_service.GetTeamParam(teamId, "user_id1") > 0 && team_service.GetTeamParam(teamId, "user_id2") > 0 {
		appG.Response(http.StatusBadRequest, e.TEAM_FULL, nil)
		return
	}
	//创建order订单
	order := models.Order{
		CurLevel:    curLevel,
		TargetLevel: targetLevel,
		UserId:      userId,
		GameType:    team_service.GetTeamParam(teamId, "game_type"),
		GameZone:    team_service.GetTeamParam(teamId, "big_zone"),
		Contact:     contact,
		Qq:          qq,
		Description: description,
	}
	if !order_service.CreateOrder(&order, teamId) {
		appG.Response(http.StatusBadRequest, e.ERROR, nil)
		return
	}
	team := models.Team{
		TeamId: teamId,
	}
	m := make(map[string]interface{})
	if team_service.GetTeamParam(teamId, "num") == 0 {
		m["user_id1"] = userId
	}
	if team_service.GetTeamParam(teamId, "num") == 1 {
		m["user_id2"] = userId
	}

	if !team.Updates(m) {
		appG.Response(http.StatusBadRequest, e.ERROR, nil)
		return
	}
	data := gin.H{}
	data["order_id"] = order.OrderId
	appG.Response(http.StatusOK, e.SUCCESS, data)
	return
}

//解散
