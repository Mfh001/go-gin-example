package models

import (
	"github.com/jinzhu/gorm"
)

type Profit struct {
	UserId                           int `json:"user_id" gorm:"primary_key;type:int(12);not null"`
	OrderTotalTimes                  int `json:"order_total_times" gorm:"type:int(12);not null;default:0"`                    //累计订单次数
	OrderTotalTimesStatus            int `json:"order_total_times_status" gorm:"type:int(12);not null;default:0"`             //累计订单次数//累计订单领取状态
	OrderYesterdayPublishTimes       int `json:"order_yesterday_publish_times" gorm:"type:int(12);not null;default:0"`        //昨日下级总发单次数
	OrderYesterdayTakerTimes         int `json:"order_yesterday_taker_times" gorm:"type:int(12);not null;default:0"`          //昨日下级总接单次数
	OrderTodayPublishTimes           int `json:"order_today_publish_times" gorm:"type:int(12);not null;default:0"`            //当日下级总发单次数
	OrderTodayTakerTimes             int `json:"order_today_taker_times" gorm:"type:int(12);not null;default:0"`              //当日下级总接单次数
	OrderYesterdayPublishProfit      int `json:"order_yesterday_publish_profit" gorm:"type:int(12);not null;default:0"`       //昨日下级总发单收益
	OrderYesterdayTakerProfit        int `json:"order_yesterday_taker_profit" gorm:"type:int(12);not null;default:0"`         //昨日下级总接单收益
	OrderYesterdayAgentPublishProfit int `json:"order_yesterday_agent_publish_profit" gorm:"type:int(12);not null;default:0"` //昨日发单上级收益
	OrderYesterdayAgentTakerProfit   int `json:"order_yesterday_agent_taker_profit" gorm:"type:int(12);not null;default:0"`   //昨日接单上级收益
	ResetTime                        int `json:"reset_time" gorm:"type:int(12);not null;default:0"`                           //重置时间
}

//insert
func (info Profit) Insert() bool {
	create := db.Create(&info)
	if create.Error != nil {
		return false
	}
	return true
}

//update
func (info Profit) Updates(shop map[string]interface{}) bool {
	err := db.Model(&info).Updates(shop).Error
	if err != nil {
		return false
	}
	return true
}

//insert and update
func (info Profit) Save() bool {
	create := db.Save(&info)
	if create.Error != nil {
		return false
	}
	return true
}

//select
func (info *Profit) First() (int, error) {
	//err := db.Where("user_id=?",info.UserId).Find(&info)
	err := db.First(&info).Error
	if gorm.IsRecordNotFoundError(err) {
		return 0, nil
	} else if err != nil {
		return -1, err
	}
	return 1, nil
}

//select all
func FindUserProfits(infos *[]Profit) (bool, error) {
	err := db.Find(infos).Error
	if gorm.IsRecordNotFoundError(err) {
		*infos = []Profit{}
		return true, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}
