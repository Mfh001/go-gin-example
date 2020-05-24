package cron_service

import (
	"github.com/EDDYCJY/go-gin-example/pkg/gredis"
	"github.com/EDDYCJY/go-gin-example/pkg/logging"
	"github.com/EDDYCJY/go-gin-example/service/profit_service"
	"github.com/robfig/cron"
	"strconv"
	"strings"
)

func Start() {
	//开启定时任务
	c := cron.New()
	//每小时 自动评价订单
	err := c.AddFunc("0 10 0 * * ?", func() {
		Loop()
	})
	if err != nil {
		logging.Error("init cron error")
	}
	c.Start()
	//defer c.Stop()
}

func Loop() {
	keys := gredis.GetKeys("game_profit:")
	for _, v := range keys {
		arr := strings.Split(v, ":")
		if len(arr) == 2 {
			userId, err := strconv.Atoi(arr[1])
			if err == nil && userId > 0 {
				go profit_service.CalDailyProfit(userId)
			}
		}

	}
}
