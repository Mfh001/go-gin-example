package auth_service

import (
	"encoding/json"
	var_const "github.com/EDDYCJY/go-gin-example/const"
	"github.com/EDDYCJY/go-gin-example/models"
	"github.com/EDDYCJY/go-gin-example/pkg/logging"
	"github.com/EDDYCJY/go-gin-example/pkg/setting"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type WxLoginUserInfo struct {
	Code      string `form:"code" valid:"Required;MinSize(30);MaxSize(32)"`
	NickName  string `json:"nick_name" form:"nickname" valid:"Required;MaxSize(50)"`
	AvatarUrl string `json:"avatar_url" form:"avatar_url" valid:"Required;MaxSize(300)"`
	Gender    int    `json:"gender" form:"gender" valid:"Required;Range(1,2)"`
	Province  string `json:"province" form:"province"`
	City      string `json:"city" form:"city"`
}

func (info *WxLoginUserInfo) WXLogin() (string, bool) {
	params := url.Values{}
	Url, _ := url.Parse("https://api.weixin.qq.com/sns/jscode2session")
	params.Set("appid", setting.AppSetting.WXAppID)
	params.Set("secret", setting.AppSetting.WXSecret)
	params.Set("js_code", info.Code)
	params.Set("grant_type", "authorization_code")
	//如果参数中有中文参数,这个方法会进行URLEncode
	Url.RawQuery = params.Encode()
	urlPath := Url.String()
	resp, _ := http.Get(urlPath)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	var dat map[string]interface{}
	if err := json.Unmarshal([]byte(string(body)), &dat); err == nil {
		//dat["session_key"] = "wwwwwwwwwwwwwwwwwwwwwwwwwwwwwwww"
		//dat["openid"] = "wwwwwwwwwwwwwwwwwwwwwwwwwwwwwwww"
		_, ok := dat["session_key"]

		if !ok {
			return "", false
		} else {
			str := dat["session_key"]
			v, ok := str.(string)
			if ok {
				//openid要与userid对应 没有对应的userid要立即生成
				userExist := false
				userId, err := models.GetUserIdByOpenId(dat["openid"].(string))
				if err == nil {
					if userId == 0 {
						//创建user
						userId, userExist = CreateUserInfo(dat["openid"].(string))
						if userExist == false || userId == 0 {
							return "", false
						}
					}
					dat["userid"] = userId
					userExist = true
				}
				if userExist {
					sessionKey, err := SetWXCode(v, dat)
					if sessionKey == "" || err != nil {
						return "", false
					}
					//更新用户信息DB
					userInfo := models.User{
						UserId: userId,
					}
					var dbInfo = make(map[string]interface{})
					dbInfo["nick_name"] = info.NickName
					dbInfo["avatar_url"] = info.AvatarUrl
					dbInfo["gender"] = info.Gender
					dbInfo["type"] = var_const.UserTypeNormal
					dbInfo["city"] = info.City
					dbInfo["province"] = info.Province
					if !userInfo.Updates(dbInfo) {
						logInfo, _ := json.Marshal(dbInfo)
						logging.Error("WXLogin:db-user-Updates" + string(logInfo))
						return "", false
					}
					return sessionKey, true
				}
			}
		}
	}
	return "", false
}

func CreateUserInfo(openId string) (int, bool) {
	newUserId, err := IncrUserId()
	if err != nil || newUserId == 0 {
		return 0, false
	}
	info := models.User{
		UserId:  newUserId,
		OpenId:  openId,
		RegTime: int(time.Now().Unix()),
	}
	if info.Insert() {
		return info.UserId, true
	}
	return 0, false
}
