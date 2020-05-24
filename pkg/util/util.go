package util

import (
	"github.com/EDDYCJY/go-gin-example/pkg/setting"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	"regexp"
	"time"
)

// Setup Initialize the util
func Setup() {
	jwtSecret = []byte(setting.AppSetting.JwtSecret)
}

func IsPhoneNum(phone string) bool {
	reg := `^1([38][0-9]|14[579]|5[^4]|16[6]|7[1-35-8]|9[189])\d{8}$`
	rgx := regexp.MustCompile(reg)
	if rgx.MatchString(phone) {
		return true
	}
	return false
}

func SendSMSCode(phone string, code string) {
	client, _ := dysmsapi.NewClientWithAccessKey("cn-hangzhou", "LTAIuRQqDhPUYhiU", "g1pNUm72YERpIC0XQDuy8O6uCDJtpt")

	request := dysmsapi.CreateSendSmsRequest()
	request.Scheme = "https"

	request.PhoneNumbers = phone
	request.SignName = "富松"
	request.TemplateCode = "SMS_164681343"
	request.TemplateParam = "{\"code\":\"" + code + "\"}"
	_, _ = client.SendSms(request)
	//TODO log
	//response, err2 := client.SendSms(request)
	//if err2 != nil {
	//fmt.Print(err2.Error())
	//}
	//fmt.Printf("response is %#v\n", response)
}

//判断该时间戳是否属于今天
func IsToday(tick int) bool {
	currentTime := time.Now()
	startTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 0, 0, 0, 0, currentTime.Location())
	if int64(tick) >= startTime.Unix() {
		return true
	}
	return false
}
