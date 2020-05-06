package check_service

import (
	"encoding/json"
	"fmt"
	"github.com/EDDYCJY/go-gin-example/models"
	"github.com/EDDYCJY/go-gin-example/pkg/gredis"
	"github.com/EDDYCJY/go-gin-example/pkg/logging"
	"github.com/EDDYCJY/go-gin-example/pkg/setting"
	"strconv"
	"time"
)

func GetRedisKeyUserCheck(userId int) string {
	return "game_check:" + strconv.Itoa(userId)
}

func ExistUserCheck(userId int) bool {
	key := GetRedisKeyUserCheck(userId)
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
				logging.Error("ExistUserInfo:TryLock-failed-tryAgain")
				ttl, err := gredis.GetTTL(lock.GetKey())
				if err != nil {
					logging.Error("ExistUserInfo:TryLock-GetTTL-failed")
				}
				time.Sleep(time.Duration(ttl/2) * time.Second)
			} else {
				break
			}
		}
	}

	info := models.Check{
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

func GetUserCheckInfo(userId int) (map[string]interface{}, error) {
	if !ExistUserCheck(userId) {
		return nil, fmt.Errorf("GetUserCheckInfo:userIdnoExist")
	}
	fields := []string{"game_id", "game_server", "game_pos", "game_level", "img_url"}
	data, err := gredis.HMGet(GetRedisKeyUserCheck(userId), fields...)
	if err != nil {
		logging.Error("GetUserCheckInfo:" + strconv.Itoa(userId))
		return nil, err
	}
	var m = make(map[string]interface{})
	for i := 0; i < len(fields); i++ {
		key, _ := data[i].([]byte)
		m[fields[i]] = string(key)
	}
	return m, nil
}

func GetCheckList(checks *[]models.Check) {
	_, err := models.FindChecks(checks)
	if err != nil {
		logging.Error("GetCheckList:db-FindChecks")
	}
	return
}
