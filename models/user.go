package models

import (
	"github.com/jinzhu/gorm"
)

type User struct {
	UserId      int    `json:"user_id" gorm:"primary_key;type:int(12);not null"`
	OpenId      string `json:"open_id" gorm:"unique;not null"`
	Balance     int    `json:"balance" gorm:"type:int(12);not null;default:0"`
	Margin      int    `json:"margin" gorm:"type:int(12);not null;default:0"`
	TeamCardNum int    `json:"team_card_num" form:"-" gorm:"type:int(12);not null;default:0"`
	NickName    string `json:"nick_name" gorm:"type:varchar(32);not null;default:''"`
	AvatarUrl   string `json:"avatar_url" gorm:"type:varchar(100);not null;default:''"`
	Phone       string `json:"phone" gorm:"type:varchar(11);not null;default:''"`
	Gender      int    `json:"gender" gorm:"type:int(2);not null;default:1"`
	Type        int    `json:"type" gorm:"type:int(2);not null;default:1"`
	City        string `json:"city" gorm:"type:varchar(8);not null;default:''"`
	Province    string `json:"province" gorm:"type:varchar(8);not null;default:''"`
	RegTime     int    `json:"reg_time" gorm:"type:int(12);not null;default:0"`
	CheckPass   int    `json:"check_pass" gorm:"type:int(2);not null;default:0"`
	GameId      string `gorm:"type:varchar(20);not null;default:''" form:"game_id" json:"game_id"  valid:"Required;MaxSize(20)"`
	GameServer  int    `gorm:"type:int(5);not null;default:0" form:"game_server" json:"game_server" valid:"Required;Range(1, 512)"`
	GamePos     int    `gorm:"type:int(5);not null;default:0" form:"game_pos" json:"game_pos" valid:"Required;Range(1, 512)"`
	GameLevel   int    `gorm:"type:int(5);not null;default:0" form:"game_level" json:"game_level" valid:"Required;Range(1, 200)"`
	ImgUrl      string `gorm:"type:varchar(100);not null;default:''" form:"img_url" json:"img_url"  valid:"Required;MaxSize(100)"`
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

//select
func (info *User) First() (int, error) {
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
func FindUsers(infos *[]User) (bool, error) {
	err := db.Find(infos).Error
	if gorm.IsRecordNotFoundError(err) {
		*infos = []User{}
		return true, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
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

func (info User) FindPhone() (int, error) {
	err := db.Select("user_id").Where("phone = ?", info.Phone).Take(&info).Error
	if gorm.IsRecordNotFoundError(err) {
		return -1, nil
	} else if err != nil {
		return -1, err
	}
	return info.UserId, nil
}
