package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
)

type Pay struct {
	OrderId   int    `json:"order_id" form:"-" gorm:"primary_key;type:int(12);not null"`
	TradeNo   string `json:"trade_no" gorm:"type:varchar(50);not null;default:''"`
	UserId    int    `json:"user_id" form:"user_id" gorm:"type:int(12);not null" valid:"Required;Range(1, 1000000000)"`
	NickName  string `json:"nick_name" gorm:"type:varchar(32);not null;default:''"`
	PayAmount int    `json:"pay_amount" form:"pay_amount" gorm:"type:int(12);not null;default:0"` //用户支付的订单费用 与订单价格相同 单位分
	PayDesc   string `json:"pay_desc" form:"-" gorm:"type:varchar(100);not null;default:''"`
	PayIp     string `json:"pay_ip" form:"-" gorm:"type:varchar(30);not null;default:''"`
	TradeType string `json:"trade_type" form:"-" gorm:"type:varchar(30);not null;default:''"`
	Status    int    `json:"status" form:"-" gorm:"type:int(12);not null;default:0"`
	PayTime   int    `json:"pay_time" gorm:"type:int(12);not null;default:0"`
	RegTime   int    `json:"reg_time" gorm:"type:int(12);not null;default:0"`
}

//insert
func (info Pay) Insert() bool {
	if info.OrderId <= 0 {
		return false
	}
	create := db.Create(&info)
	if create.Error != nil {
		return false
	}
	return true
}

//update
func (info Pay) Updates(m map[string]interface{}) bool {
	if info.OrderId <= 0 {
		return false
	}
	err := db.Model(&info).Updates(m).Error
	if err != nil {
		return false
	}
	return true
}

//insert and update
func (info Pay) Save() bool {
	if info.OrderId <= 0 {
		return false
	}
	create := db.Save(&info)
	if create.Error != nil {
		return false
	}
	return true
}

//select
func (info *Pay) First() (int, error) {
	if info.OrderId <= 0 {
		return -1, fmt.Errorf("First:ProfitnoExist")
	}
	err := db.First(&info).Error
	if gorm.IsRecordNotFoundError(err) {
		return 0, nil
	} else if err != nil {
		return -1, err
	}
	return 1, nil
}

func (info *Pay) GetOrderInfoByTradeNo() (bool, error) {
	err := db.Select("user_id, pay_amount, order_id, status").Where("trade_no = ?", info.TradeNo).First(&info).Error
	if err != nil {
		return false, err
	}
	return true, nil
}
