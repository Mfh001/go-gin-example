package bank_service

import (
	"encoding/json"
	"github.com/EDDYCJY/go-gin-example/models"
	"github.com/EDDYCJY/go-gin-example/pkg/gredis"
	"github.com/EDDYCJY/go-gin-example/pkg/logging"
	"github.com/EDDYCJY/go-gin-example/pkg/setting"
	"strconv"
	"time"
)

func GetRedisKeyBank(id int) string {
	return "game_bank_card_info:" + strconv.Itoa(id)
}

func ExistBank(id int) bool {
	if id <= 0 {
		return false
	}
	key := GetRedisKeyBank(id)
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
				logging.Error("ExistBank:TryLock-failed-tryAgain")
				ttl, err := gredis.GetTTL(lock.GetKey())
				if err != nil {
					logging.Error("ExistBank:TryLock-GetTTL-failed")
				}
				time.Sleep(time.Duration(ttl/2) * time.Second)
			} else {
				break
			}
		}
	}

	info := models.Exchange{
		Id: id,
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

func GetBankParam(id int, param string) int {
	if !ExistBank(id) {
		return 0
	}
	strParam, err := gredis.HGet(GetRedisKeyBank(id), param)
	if err != nil {
		logging.Error("GetBankParam:" + strconv.Itoa(id))
		return 0
	}
	if strParam == "" {
		strParam = "0"
	}
	p, err := strconv.Atoi(strParam)
	if err != nil {
		logging.Error("GetBankParam:" + strParam)
		return 0
	}
	return p
}
func GetBankParamString(id int, param string) string {
	if !ExistBank(id) {
		return ""
	}
	strParam, err := gredis.HGet(GetRedisKeyBank(id), param)
	if err != nil {
		logging.Error("GetBankParamString:" + strconv.Itoa(id))
		return ""
	}
	return strParam
}
