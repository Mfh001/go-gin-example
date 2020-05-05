package models

import "github.com/jinzhu/gorm"

type Check struct {
	UserId     int    `gorm:"primary_key;type:int(12);not null" form:"user_id" json:"user_id" valid:"Required;Range(1, 1000000000)"`
	GameId     string `gorm:"type:varchar(20);not null;default:''" form:"game_id" json:"game_id"  valid:"Required;MaxSize(20)"`
	GameServer int    `gorm:"type:int(5);not null;default:0" form:"game_server" json:"game_server" valid:"Required;Range(1, 512)"`
	GamePos    int    `gorm:"type:int(5);not null;default:0" form:"game_pos" json:"game_pos" valid:"Required;Range(1, 512)"`
	GameLevel  int    `gorm:"type:int(5);not null;default:0" form:"game_level" json:"game_level" valid:"Required;Range(1, 200)"`
	ImgUrl     string `gorm:"type:varchar(100);not null;default:''" form:"img_url" json:"img_url"  valid:"Required;MaxSize(100)"`
}

//insert
func (info Check) Insert() bool {
	create := db.Create(&info)
	if create.Error != nil {
		return false
	}
	return true
}

//update
func (info Check) Updates(m map[string]interface{}) bool {
	err := db.Model(&info).Updates(m).Error
	if err != nil {
		return false
	}
	return true
}

//insert and update
func (info Check) Save() bool {
	create := db.Save(&info)
	if create.Error != nil {
		return false
	}
	return true
}

//select
func (info *Check) First() (int, error) {
	err := db.First(&info).Error
	if gorm.IsRecordNotFoundError(err) {
		return 0, nil
	} else if err != nil {
		return -1, err
	}
	return 1, nil
}

//select all
func FindChecks(infos *[]Check) (bool, error) {
	err := db.Find(infos).Error
	if gorm.IsRecordNotFoundError(err) {
		*infos = []Check{}
		return true, nil
	} else if err != nil {
		*infos = []Check{}
		return false, err
	}
	return true, nil
}
