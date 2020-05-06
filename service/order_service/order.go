package order_service

import (
	"encoding/json"
	var_const "github.com/EDDYCJY/go-gin-example/const"
	"github.com/EDDYCJY/go-gin-example/models"
	"github.com/EDDYCJY/go-gin-example/pkg/gredis"
	"github.com/EDDYCJY/go-gin-example/pkg/logging"
	"github.com/EDDYCJY/go-gin-example/pkg/setting"
	"github.com/EDDYCJY/go-gin-example/service/auth_service"
	"strconv"
	"time"
)

func GetRedisKeyOrder(id int) string {
	return "game_order:" + strconv.Itoa(id)
}

func getRedisKeyOrderIncr() string {
	return "max_order"
}

func IncrOrderId() (int, error) {
	return gredis.Incr(getRedisKeyOrderIncr())
}

func CreateOrder(form *models.Order) bool {
	//等级idx*1000 + star
	price := 0
	for i := 0; i < len(setting.PlatFormLevelAll); i++ {
		if setting.PlatFormLevelAll[i].Idx > form.CurLevel && setting.PlatFormLevelAll[i].Idx <= form.TargetLevel {
			price += setting.PlatFormLevelAll[i].Price
		}
	}
	//
	orderId, err := IncrOrderId()
	if err != nil {
		return false
	}
	form.OrderId = orderId
	nickName, _ := auth_service.GetUserNickName(form.UserId)
	form.NickName = nickName
	form.Price = price
	form.Status = var_const.OrderStatusWaitPay
	form.RegTime = int(time.Now().Unix())
	if !form.Save() {
		log, _ := json.Marshal(form)
		logging.Error("CreateOrder:form-" + string(log))
		return false
	}
	return true
}
