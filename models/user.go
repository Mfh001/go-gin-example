package models

import (
	"github.com/jinzhu/gorm"
)

type User struct {
	UserId    int    `gorm:"primary_key;type:int(12);not null"`
	OpenId    string `gorm:"unique;not null"`
	NickName  string `gorm:"type:varchar(32);not null;default:''"`
	AvatarUrl string `gorm:"type:varchar(100);not null;default:''"`
	Phone     string `gorm:"type:varchar(11);not null;default:''"`
	Gender    int    `gorm:"type:int(2);not null;default:1"`
	Type      int    `gorm:"type:int(2);not null;default:1"`
	City      string `gorm:"type:varchar(8);not null;default:''"`
	Province  string `gorm:"type:varchar(8);not null;default:''"`
	RegTime   string `gorm:"type:varchar(25);not null;default:''"`
}

type WXCode struct {
	OpenId     string `valid:"Required;MaxSize(100)" json:"openid"`
	UserId     int    `valid:"Required;Max(100000000)" json:"userid"`
	SessionKey string `valid:"Required;MaxSize(100)" json:"session_key"`
}

//insert
func (info User) Insert() bool {
	create := db.Create(&info)
	if create.Error != nil {
		return false
	}
	return true
}

//update
func (info User) Updates(shop map[string]interface{}) bool {
	err := db.Model(&info).Updates(shop).Error
	if err != nil {
		return false
	}
	return true
}

// GetUserIdByOpenId
func GetUserIdByOpenId(openId string) (int, error) {
	var user User
	err := db.Select("user_id").Where("open_id = ? ", openId).First(&user).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return 0, err
	}

	return user.UserId, nil
}
