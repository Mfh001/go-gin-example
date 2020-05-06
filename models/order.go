package models

import "github.com/jinzhu/gorm"

type Order struct {
	OrderId     int    `json:"order_id" form:"-" gorm:"primary_key;type:int(12);not null"`
	Price       int    `json:"price" form:"-" gorm:"type:int(12);not null"`
	Status      int    `json:"status" form:"-" gorm:"type:int(12);not null"`
	UserId      int    `json:"user_id" form:"user_id" gorm:"type:int(12);not null" valid:"Required;Range(1, 1000000000)"`
	NickName    string `json:"nick_name" gorm:"type:varchar(32);not null;default:''"`
	GameType    int    `json:"game_type" form:"game_type" gorm:"type:int(2);not null" valid:"Range(0, 1)"`
	OrderType   int    `json:"order_type" form:"order_type" gorm:"type:int(2);not null" valid:"Range(0, 1)"`
	InsteadType int    `json:"instead_type" form:"instead_type" gorm:"type:int(2);not null" valid:"Range(0, 2)"`
	GameZone    int    `json:"game_zone" form:"game_zone" gorm:"type:int(8);not null" valid:"Range(0, 4000)"`
	RunesLevel  int    `json:"runes_level" form:"runes_level" gorm:"type:int(4);not null" valid:"Range(0, 200)"`
	HeroNum     int    `json:"hero_num" form:"hero_num" gorm:"type:int(4);not null;default:0" valid:"Range(0, 1000)"`
	CurLevel    int    `json:"cur_level" form:"cur_level" gorm:"type:int(8);not null" valid:"Range(0, 10000)"`
	TargetLevel int    `json:"target_level" form:"target_level" gorm:"type:int(8);not null" valid:"Required;Range(1, 10000)"`
	RegTime     int    `json:"reg_time" gorm:"type:int(12);not null;default:0"`
	UpdTime     int    `json:"upd_time" gorm:"type:int(12);not null;default:0"`
}

//insert
func (info Order) Insert() bool {
	create := db.Create(&info)
	if create.Error != nil {
		return false
	}
	return true
}

//update
func (info Order) Updates(m map[string]interface{}) bool {
	err := db.Model(&info).Updates(m).Error
	if err != nil {
		return false
	}
	return true
}

//insert and update
func (info Order) Save() bool {
	create := db.Save(&info)
	if create.Error != nil {
		return false
	}
	return true
}

//select
func (info *Order) First() (int, error) {
	err := db.First(&info).Error
	if gorm.IsRecordNotFoundError(err) {
		return 0, nil
	} else if err != nil {
		return -1, err
	}
	return 1, nil
}

// Delete
//func (info *Order)Delete() bool {
//	//err := db.Where("id = ?", id).Delete(&Tag{}).Error
//	if err := db.Delete(info).Error; err != nil {
//		return false
//	}
//	return true
//}
