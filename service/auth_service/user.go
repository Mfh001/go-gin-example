package auth_service

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	var_const "github.com/EDDYCJY/go-gin-example/const"
	"github.com/EDDYCJY/go-gin-example/models"
	"github.com/EDDYCJY/go-gin-example/pkg/gredis"
	"github.com/EDDYCJY/go-gin-example/pkg/logging"
	"github.com/EDDYCJY/go-gin-example/pkg/setting"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

func getRedisKeyUserIncr() string {
	return "max_user"
}

func getRedisKeyWXCode(str string) string {
	return "sessionKey:" + str
}
func getRedisKeyAccessToken() string {
	return "accessToken:" + setting.AppSetting.WXAppID
}
func GetRedisKeyUserInfo(id int) string {
	return "game_user:" + strconv.Itoa(id)
}
func GetRedisKeySmsCode(phone string) string {
	return "code:" + phone
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

func GetAccessToken() string {
	res, err := gredis.Get(getRedisKeyAccessToken())
	if err != nil {
		return ""
	}
	return res
}

func SetAccessToken(token string) bool {
	_ = gredis.Set(getRedisKeyAccessToken(), token, var_const.AccessTokenExpireTime)
	return true
}

func UpdateAccessToken() (string, int) {
	params := url.Values{}
	Url, _ := url.Parse("https://api.weixin.qq.com/cgi-bin/token")
	params.Set("appid", setting.AppSetting.WXAppID)
	params.Set("secret", setting.AppSetting.WXSecret)
	params.Set("grant_type", "client_credential")
	//如果参数中有中文参数,这个方法会进行URLEncode
	Url.RawQuery = params.Encode()
	urlPath := Url.String()
	resp, _ := http.Get(urlPath)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	var dat map[string]interface{}
	if err := json.Unmarshal([]byte(string(body)), &dat); err == nil {
		_, ok := dat["access_token"]

		if !ok {
			return dat["errmsg"].(string), dat["errcode"].(int)
		} else {
			SetAccessToken(dat["access_token"].(string))
			return dat["access_token"].(string), 0
		}
	}
	return "", 1
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

func UpdateUserInfoToRedis(userId int) bool {
	logging.Info("UpdateUserInfoToRedis:" + strconv.Itoa(userId))
	key := GetRedisKeyUserInfo(userId)
	info := models.User{
		UserId: userId,
	}
	dbRes, err := info.First()
	if dbRes == 0 {
		return false
	} else if dbRes == -1 {
		return false
	}
	//将数据放入redis
	jData, _ := json.Marshal(info)
	m := make(map[string]interface{})
	if err := json.Unmarshal(jData, &m); err != nil {
		return false
	}
	_, err = gredis.HMSet(key, m)
	if err != nil {
		return false
	}
	return true
}

func GetUserInfo(userId int) (map[string]interface{}, error) {
	if !ExistUserInfo(userId) {
		return nil, fmt.Errorf("GetUserInfo:userIdnoExist")
	}
	fields := []string{"user_id", "nick_name", "avatar_url", "phone", "type", "check_pass",
		"game_id", "game_server", "game_pos", "game_level", "img_url", "balance", "agent_id"}
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

func GetUserOpenId(userId int) (string, error) {
	if !ExistUserInfo(userId) {
		return "", fmt.Errorf("GetUserOpenId:userIdnoExist")
	}
	openId, err := gredis.HGet(GetRedisKeyUserInfo(userId), "open_id")
	if err != nil {
		logging.Error("GetUserOpenId:" + strconv.Itoa(userId))
		return "", err
	}
	return openId, nil
}
func GetUserPhone(userId int) (string, error) {
	if !ExistUserInfo(userId) {
		return "", fmt.Errorf("GetUserPhone:userIdnoExist")
	}
	phone, err := gredis.HGet(GetRedisKeyUserInfo(userId), "phone")
	if err != nil {
		logging.Error("GetUserOpenId:" + strconv.Itoa(userId))
		return "", err
	}
	return phone, nil
}

func AddUserMargin(userId int, amount int) bool {
	userInfo := models.User{
		UserId: userId,
	}
	margin := GetUserParam(userId, "margin")
	logging.Info("AddUserMargin:userId-" + strconv.Itoa(userId) + ",amount-" + strconv.Itoa(amount) + ",margin-" + strconv.Itoa(margin))
	margin += amount
	var db2Info = make(map[string]interface{})
	db2Info["margin"] = margin

	if !userInfo.Updates(db2Info) {
		log, _ := json.Marshal(db2Info)
		logging.Error("AddUserMargin:db-userInfo-failed-" + string(log))
		return false
	}
	logging.Info("AddUserMargin Success:userId-" + strconv.Itoa(userId) + ",amount-" + strconv.Itoa(amount))
	return true
}

func RemoveUserMargin(userId int, amount int) bool {
	userInfo := models.User{
		UserId: userId,
	}
	margin := GetUserParam(userId, "margin")
	logging.Info("RemoveUserMargin:userId-" + strconv.Itoa(userId) + ",amount-" + strconv.Itoa(amount) + ",margin-" + strconv.Itoa(margin))
	margin -= amount
	if margin < 0 {
		margin = 0
	}
	var db2Info = make(map[string]interface{})
	db2Info["margin"] = margin

	if !userInfo.Updates(db2Info) {
		log, _ := json.Marshal(db2Info)
		logging.Error("RemoveUserMargin:db-userInfo-failed-" + string(log))
		return false
	}
	logging.Info("RemoveUserMargin Success:userId-" + strconv.Itoa(userId) + ",amount-" + strconv.Itoa(amount))
	return true
}
func AddUserBalance(userId int, amount int, flag string) bool {
	userInfo := models.User{
		UserId: userId,
	}
	margin := GetUserParam(userId, "balance")
	logging.Info("AddUserBalance:userId-" + strconv.Itoa(userId) + ",amount-" + strconv.Itoa(amount) + ",balance-" + strconv.Itoa(margin) + ",flag-" + flag)
	margin += amount
	var db2Info = make(map[string]interface{})
	db2Info["balance"] = margin

	if !userInfo.Updates(db2Info) {
		log, _ := json.Marshal(db2Info)
		logging.Error("AddUserBalance:db-userInfo-failed-" + string(log))
		return false
	}
	logging.Info("AddUserBalance Success:userId-" + strconv.Itoa(userId) + ",amount-" + strconv.Itoa(amount) + ",flag-" + flag)
	return true
}

func RemoveUserBalance(userId int, amount int) bool {
	userInfo := models.User{
		UserId: userId,
	}
	margin := GetUserParam(userId, "balance")
	logging.Info("RemoveUserBalance:userId-" + strconv.Itoa(userId) + ",amount-" + strconv.Itoa(amount) + ",balance-" + strconv.Itoa(margin))
	margin -= amount
	if margin < 0 {
		margin = 0
	}
	var db2Info = make(map[string]interface{})
	db2Info["balance"] = margin

	if !userInfo.Updates(db2Info) {
		log, _ := json.Marshal(db2Info)
		logging.Error("RemoveUserBalance:db-userInfo-failed-" + string(log))
		return false
	}
	logging.Info("RemoveUserBalance Success:userId-" + strconv.Itoa(userId) + ",amount-" + strconv.Itoa(amount))
	return true
}

func GetUserParam(userId int, param string) int {
	if !ExistUserInfo(userId) {
		return 0
	}
	strParam, err := gredis.HGet(GetRedisKeyUserInfo(userId), param)
	if err != nil {
		logging.Error("GetUserParam:" + param + ":" + strconv.Itoa(userId))
		UpdateUserInfoToRedis(userId)
		strParam, _ = gredis.HGet(GetRedisKeyUserInfo(userId), param)
		//return 0
	}
	if strParam == "" {
		strParam = "0"
	}
	p, err := strconv.Atoi(strParam)
	if err != nil {
		logging.Error("GetTeamParam:" + strParam)
		return 0
	}
	return p
}
func GetUserParamString(userId int, param string) string {
	if !ExistUserInfo(userId) {
		return ""
	}
	strParam, err := gredis.HGet(GetRedisKeyUserInfo(userId), param)
	if err != nil {
		logging.Error("GetUserParamString:" + param + ":" + strconv.Itoa(userId))
		UpdateUserInfoToRedis(userId)
		strParam, _ = gredis.HGet(GetRedisKeyUserInfo(userId), param)
		//return ""
	}
	return strParam
}
