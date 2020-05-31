package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
)

type BalanceLog struct {
	Id         int    `form:"-" json:"id" gorm:"AUTO_INCREMENT;primary_key;"`
	UserId     int    `form:"user_id" json:"user_id" gorm:"type:int(12);not null"`
	NickName   string `json:"nick_name" gorm:"type:varchar(32);not null;default:''"`
	TotalMoney int    `form:"-" json:"total_money" gorm:"type:int(12);not null;default:0"`
	Money      int    `form:"money" json:"money" gorm:"type:int(12);not null;default:0"`
	LogType    int    `form:"-" json:"log_type" gorm:"type:int(2);not null;default:0"`
	Status     int    `form:"-" json:"status" gorm:"type:int(2);not null;default:0"`
	Remarks    string `json:"remarks" gorm:"type:varchar(32);not null;default:''"`
	RegTime    int    `json:"reg_time" gorm:"type:int(12);not null;default:0"`
}

//insert
func (info BalanceLog) Insert() bool {
	if info.UserId <= 0 {
		return false
	}
	create := db.Create(&info)
	if create.Error != nil {
		return false
	}
	return true
}

//update

func (info BalanceLog) Updates(m map[string]interface{}) bool {
	if info.Id <= 0 {
		return false
	}
	err := db.Model(&info).Updates(m).Error
	if err != nil {
		return false
	}
	return true
}

//insert and update
func (info BalanceLog) Save() bool {
	if info.UserId <= 0 || info.Id <= 0 {
		return false
	}
	create := db.Save(&info)
	if create.Error != nil {
		return false
	}
	return true
}

//select
func (info *BalanceLog) First() (int, error) {
	if info.UserId <= 0 || info.Id <= 0 {
		return -1, fmt.Errorf("GetUserInfo:userIdnoExist")
	}
	err := db.First(&info).Error
	if gorm.IsRecordNotFoundError(err) {
		return 0, nil
	} else if err != nil {
		return -1, err
	}
	return 1, nil
}

//select all
func FindBalanceLogs(infos *[]BalanceLog) (bool, error) {
	err := db.Find(infos).Error
	if gorm.IsRecordNotFoundError(err) {
		*infos = []BalanceLog{}
		return true, nil
	} else if err != nil {
		*infos = []BalanceLog{}
		return false, err
	}
	return true, nil
}

// Delete
func (info *BalanceLog) Delete() bool {
	//err := db.Where("id = ?", id).Delete(&Tag{}).Error
	if err := db.Delete(info).Error; err != nil {
		return false
	}
	return true
}

func GetUserBalanceLogs(userId int, infos *[]BalanceLog, index int, count int) (bool, error) {
	err := db.Select("*").Where("user_id = ?", userId).Limit(count).Offset(index).Find(&infos).Error
	if gorm.IsRecordNotFoundError(err) {
		*infos = []BalanceLog{}
		return true, nil
	} else if err != nil {
		*infos = []BalanceLog{}
		return false, err
	}
	return true, nil
}
