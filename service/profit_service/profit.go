package profit_service

import (
	"encoding/json"
	"github.com/EDDYCJY/go-gin-example/models"
	"github.com/EDDYCJY/go-gin-example/pkg/gredis"
	"github.com/EDDYCJY/go-gin-example/pkg/logging"
	"github.com/EDDYCJY/go-gin-example/pkg/setting"
	"github.com/EDDYCJY/go-gin-example/pkg/util"
	"github.com/EDDYCJY/go-gin-example/service/auth_service"
	"strconv"
	"time"
)

func GetRedisKeyProfit(userId int) string {
	return "game_profit:" + strconv.Itoa(userId)
}

func ExistProfit(userId int) bool {
	key := GetRedisKeyProfit(userId)
	if gredis.Exists(key) {
		return true
	}
	//从数据库获取
	//先加锁
	lock, ok, _ := gredis.TryLock(key, setting.RedisSetting.LockTimeout)
	if lock == nil {
		return false
	}
	if !ok {
		i := 50
		for {
			if i--; i <= 0 {
				break
			}
			ok, _ := gredis.SetNX(lock.GetKey(), lock.GenerateToken(), lock.Timeout)
			if !ok {
				logging.Error("ExistProfit:TryLock-failed-tryAgain")
				ttl, err := gredis.GetTTL(lock.GetKey())
				if err != nil {
					logging.Error("ExistProfit:TryLock-GetTTL-failed")
				}
				time.Sleep(time.Duration(ttl/2) * time.Second)
			} else {
				break
			}
		}
	}

	info := models.Profit{
		UserId: userId,
	}
	dbRes, err := info.First()
	if dbRes == 0 {
		lock.UnLock()
		return false
	} else if dbRes == -1 {
		lock.UnLock()
		return false
	}
	//将数据放入redis
	jData, _ := json.Marshal(info)
	m := make(map[string]interface{})
	if err := json.Unmarshal(jData, &m); err != nil {
		lock.UnLock()
		return false
	}
	_, err = gredis.HMSet(key, m)
	if err != nil {
		lock.UnLock()
		return false
	}
	lock.UnLock()
	return true
}

func GetProfitParam(userId int, param string) int {
	if !ExistProfit(userId) {
		return 0
	}
	strParam, err := gredis.HGet(GetRedisKeyProfit(userId), param)
	if err != nil {
		logging.Error("GetProfitParam:" + strconv.Itoa(userId))
		return 0
	}
	if strParam == "" {
		strParam = "0"
	}
	p, err := strconv.Atoi(strParam)
	if err != nil {
		logging.Error("GetProfitParam:" + strParam)
		return 0
	}
	return p
}
func GetOrderParamString(userId int, param string) string {
	if !ExistProfit(userId) {
		return ""
	}
	strParam, err := gredis.HGet(GetRedisKeyProfit(userId), param)
	if err != nil {
		logging.Error("GetOrderParamString:" + strconv.Itoa(userId))
		return ""
	}
	return strParam
}

func getOrderTodayPublishProfit(orderTodayPublishTimes int) (int, int) {
	if orderTodayPublishTimes < 50 {
		return 0, 0
	}
	if orderTodayPublishTimes < 100 {
		return 3000, 200
	}
	if orderTodayPublishTimes < 300 {
		return 6000, 400
	}
	if orderTodayPublishTimes < 500 {
		return 24000, 1200
	}
	if orderTodayPublishTimes < 700 {
		return 35000, 2000
	}
	if orderTodayPublishTimes < 1000 {
		return 56000, 2800
	}
	if orderTodayPublishTimes < 5000 {
		return 80000, 4000
	}
	if orderTodayPublishTimes < 10000 {
		return 400000, 20000
	}
	if orderTodayPublishTimes < 20000 {
		return 800000, 40000
	} else {
		return 1600000, 80000
	}
}

func getOrderTodayTakerProfit(orderTodayTakerTimes int) (int, int) {
	if orderTodayTakerTimes < 50 {
		return 0, 0
	}
	if orderTodayTakerTimes < 100 {
		return 3000, 200
	}
	if orderTodayTakerTimes < 300 {
		return 6000, 400
	}
	if orderTodayTakerTimes < 500 {
		return 24000, 1200
	}
	if orderTodayTakerTimes < 700 {
		return 35000, 2000
	}
	if orderTodayTakerTimes < 1000 {
		return 56000, 2800
	}
	if orderTodayTakerTimes < 5000 {
		return 80000, 4000
	}
	if orderTodayTakerTimes < 10000 {
		return 400000, 20000
	} else {
		return 800000, 40000
	}
}

//每日返利计算
func CalDailyProfit(userId int) bool {
	if !ExistProfit(userId) {
		return false
	}
	resetTime := GetProfitParam(userId, "reset_time")
	if !util.IsToday(resetTime) {
		orderTodayPublishTimes := GetProfitParam(userId, "order_today_publish_times")
		pProfit, pAgentProfit := getOrderTodayPublishProfit(orderTodayPublishTimes)
		orderTodayTakerTimes := GetProfitParam(userId, "order_today_taker_times")
		tProfit, tAgentProfit := getOrderTodayTakerProfit(orderTodayTakerTimes)
		if pProfit > 0 || tProfit > 0 {
			//addmoney
			auth_service.AddUserBalance(userId, pProfit+tProfit, "CalDailyProfit")
		}
		//上级addmoney
		agentId := auth_service.GetUserParam(userId, "agent_id")
		if pAgentProfit+tAgentProfit > 0 {
			if agentId > 0 && auth_service.ExistUserInfo(agentId) {
				auth_service.AddUserBalance(agentId, pAgentProfit+tAgentProfit, "CalDailyProfitAgent")
			}
		}

		profit := models.Profit{
			UserId: userId,
		}
		m := make(map[string]interface{})
		m["order_yesterday_publish_profit"] = pProfit
		m["order_yesterday_taker_profit"] = tProfit
		m["order_yesterday_agent_publish_profit"] = pAgentProfit
		m["order_yesterday_agent_taker_profit"] = tAgentProfit
		if agentId == 0 {
			m["order_yesterday_agent_publish_profit"] = 0
			m["order_yesterday_agent_taker_profit"] = 0
		}
		m["reset_time"] = int(time.Now().Unix())
		m["order_yesterday_publish_times"] = orderTodayPublishTimes
		m["order_yesterday_taker_times"] = orderTodayTakerTimes
		m["order_today_publish_times"] = 0
		m["order_today_taker_times"] = 0
		log, _ := json.Marshal(m)
		if !profit.Updates(m) {
			logging.Info("CalDailyProfit failed: user_id-" + strconv.Itoa(userId) + string(log))
			return false
		}
		logging.Info("CalDailyProfit success: user_id-" + strconv.Itoa(userId) + string(log))
		return true
	}
	return false
}
