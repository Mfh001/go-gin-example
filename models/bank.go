package models

import (
	"github.com/jinzhu/gorm"
)

type BankCardInfo struct {
	UserId         int    `form:"user_id" gorm:"primary_key;type:int(12);not null"`
	BankName       string `form:"bank_name" gorm:"type:varchar(100);not null;default:''"`        //银行名称
	BankBranchName string `form:"Bank_branch_name" gorm:"type:varchar(100);not null;default:''"` //支行名称
	BankCard       string `form:"bank_card" gorm:"type:varchar(32);not null;default:''"`         //银行卡号
	UserName       string `form:"user_name" gorm:"type:varchar(32);not null;default:''"`         //开户名
	Password       string `form:"password" gorm:"type:varchar(64);not null;default:''"`          //给银行卡信息设置的密码
}

type RequestBankCardInfo struct {
	BankCardInfo BankCardInfo
	Code         string //验证码
}

func (info BankCardInfo) Insert() bool {
	create := db.Create(&info)
	if create.Error != nil {
		return false
	}
	return true
}

func (info BankCardInfo) UpdateByOneColumn(name string, val interface{}) bool {
	create := db.Model(&info).Update(name, val)
	if create.Error != nil {
		return false
	}
	return true
}

func (info BankCardInfo) Save() bool {
	create := db.Save(&info)
	if create.Error != nil {
		return false
	}
	return true
}
func (info *BankCardInfo) FindBankCardInfo() (int, error) {
	err := db.Select("bank_name, Bank_branch_name,bank_card,user_name").Where("user_id = ?", info.UserId).Take(&info).Error
	if gorm.IsRecordNotFoundError(err) {
		return -1, nil
	} else if err != nil {
		return -1, err
	}
	return info.UserId, nil
}
func (info BankCardInfo) CheckBankCardPassword() bool {
	err := db.Select("bank_name").Where("shop_id = ? AND password = ?", info.UserId, info.Password).Take(&info).Error
	if gorm.IsRecordNotFoundError(err) {
		return false
	} else if err != nil {
		return false
	}
	return true
}
