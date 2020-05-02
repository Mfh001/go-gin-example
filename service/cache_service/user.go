package cache_service

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"github.com/EDDYCJY/go-gin-example/models"
	"github.com/EDDYCJY/go-gin-example/pkg/gredis"
	"github.com/EDDYCJY/go-gin-example/pkg/logging"
	"github.com/EDDYCJY/go-gin-example/pkg/setting"
	"strconv"
)

func getRedisKeyUserIncr() string {
	return "max_user"
}

func getRedisKeyWXCode(str string) string {
	return "sessionKey:" + str
}
func getRedisKeyUserInfo(id int) string {
	return "userInfo:" + strconv.Itoa(id)
}

func IncrUserId() (int, error) {
	return gredis.Incr(getRedisKeyUserIncr())
}

func SetWXCode(str string, data map[string]interface{}) (string, error) {
	ctx := md5.New()
	ctx.Write([]byte(str))
	cipherStr := hex.EncodeToString(ctx.Sum(nil))
	if err := gredis.Set(getRedisKeyWXCode(cipherStr), data, setting.ServerSetting.WXCodeExpireTime); err != nil {
		datStr, _ := json.Marshal(data)
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
	json.Unmarshal(data, &m)
	return m, nil
}

func GetUserInfo(userId int) (map[string]interface{}, error) {
	fields := []string{"user_id", "nick_name", "avatar_url", "phone"}
	data, err := gredis.HMGet(getRedisKeyUserInfo(userId), fields...)
	if err != nil {
		logging.Error("GetUserInfo:" + strconv.Itoa(userId))
		return nil, err
	}
	var m = make(map[string]interface{})
	if len(fields) == len(data) {
		for i := 0; i < len(fields); i++ {
			m[fields[i]] = data[i]
		}
		return m, nil
	}
	return nil, err
}
