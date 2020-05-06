package auth_service

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/EDDYCJY/go-gin-example/models"
	"github.com/EDDYCJY/go-gin-example/pkg/gredis"
	"github.com/EDDYCJY/go-gin-example/pkg/logging"
	"github.com/EDDYCJY/go-gin-example/pkg/setting"
	"strconv"
	"time"
)

func getRedisKeyUserIncr() string {
	return "max_user"
}

func getRedisKeyWXCode(str string) string {
	return "sessionKey:" + str
}
func GetRedisKeyUserInfo(id int) string {
	return "game_user:" + strconv.Itoa(id)
}

func IncrUserId() (int, error) {
	return gredis.Incr(getRedisKeyUserIncr())
}

func SetWXCode(str string, data map[string]interface{}) (string, error) {
	ctx := md5.New()
	ctx.Write([]byte(str))
	cipherStr := hex.EncodeToString(ctx.Sum(nil))
	datStr, _ := json.Marshal(data)
	if err := gredis.Set(getRedisKeyWXCode(cipherStr), string(datStr), setting.ServerSetting.WXCodeExpireTime); err != nil {
		logging.Error("SetWXCode:" + cipherStr + "-" + string(datStr))
		return "", err
	}
	return cipherStr, nil
}

func GetWXCode(sessionKey string) (*models.WXCode, error) {
	data, err := gredis.Get(getRedisKeyWXCode(sessionKey))
	if err != nil {
		logging.Error("GetWXCode:" + sessionKey)
		return nil, err
	}
	var m *models.WXCode
	_ = json.Unmarshal([]byte(data), &m)
	return m, nil
}

func ExistUserInfo(userId int) bool {
	key := GetRedisKeyUserInfo(userId)
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
			i--
			if i <= 0 {
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

	info := models.User{
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

func GetUserInfo(userId int) (map[string]interface{}, error) {
	if !ExistUserInfo(userId) {
		return nil, fmt.Errorf("GetUserInfo:userIdnoExist")
	}
	fields := []string{"user_id", "nick_name", "avatar_url", "phone", "type", "check_pass",
		"game_id", "game_server", "game_pos", "game_level", "img_url"}
	data, err := gredis.HMGet(GetRedisKeyUserInfo(userId), fields...)
	if err != nil {
		logging.Error("GetUserInfo:" + strconv.Itoa(userId))
		return nil, err
	}
	var m = make(map[string]interface{})
	for i := 0; i < len(fields); i++ {
		key, _ := data[i].([]byte)
		m[fields[i]] = string(key)
	}
	return m, nil
}

func GetUserCheckPassState(userId int) (int, error) {
	if !ExistUserInfo(userId) {
		return 0, fmt.Errorf("GetUserCheckPassState:userIdnoExist")
	}
	strState, err := gredis.HGet(GetRedisKeyUserInfo(userId), "check_pass")
	if err != nil {
		logging.Error("GetUserCheckPassState:" + strconv.Itoa(userId))
		return 0, err
	}
	state, err := strconv.Atoi(strState)
	if err != nil {
		return 0, err
	}
	return state, nil
}

func GetUserNickName(userId int) (string, error) {
	if !ExistUserInfo(userId) {
		return "", fmt.Errorf("GetUserNickName:userIdnoExist")
	}
	nickName, err := gredis.HGet(GetRedisKeyUserInfo(userId), "nick_name")
	if err != nil {
		logging.Error("GetUserNickName:" + strconv.Itoa(userId))
		return "", err
	}
	return nickName, nil
}
