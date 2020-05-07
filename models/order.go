package models

import "github.com/jinzhu/gorm"

type Order struct {
	OrderId       int    `json:"order_id" form:"-" gorm:"primary_key;type:int(12);not null"`
	Price         int    `json:"price" form:"-" gorm:"type:int(12);not null"`
	Status        int    `json:"status" form:"-" gorm:"type:int(12);not null"`
	TradeNo       string `json:"trade_no" gorm:"type:varchar(50);not null;default:''"`
	UserId        int    `json:"user_id" form:"user_id" gorm:"type:int(12);not null" valid:"Required;Range(1, 1000000000)"`
	NickName      string `json:"nick_name" gorm:"type:varchar(32);not null;default:''"`
	GameType      int    `json:"game_type" form:"game_type" gorm:"type:int(2);not null" valid:"Range(0, 1)"`
	OrderType     int    `json:"order_type" form:"order_type" gorm:"type:int(2);not null" valid:"Range(0, 1)"`
	InsteadType   int    `json:"instead_type" form:"instead_type" gorm:"type:int(2);not null" valid:"Range(0, 2)"`
	GameZone      int    `json:"game_zone" form:"game_zone" gorm:"type:int(8);not null" valid:"Range(0, 4000)"`
	RunesLevel    int    `json:"runes_level" form:"runes_level" gorm:"type:int(4);not null" valid:"Range(0, 200)"`
	HeroNum       int    `json:"hero_num" form:"hero_num" gorm:"type:int(4);not null;default:0" valid:"Range(0, 1000)"`
	CurLevel      int    `json:"cur_level" form:"cur_level" gorm:"type:int(8);not null" valid:"Range(0, 10000)"`
	TargetLevel   int    `json:"target_level" form:"target_level" gorm:"type:int(8);not null" valid:"Required;Range(1, 10000)"`
	GameAcc       string `json:"game_acc" form:"game_acc" gorm:"type:varchar(60);not null;default:''"`
	GamePwd       string `json:"game_pwd" form:"game_pwd" gorm:"type:varchar(60);not null;default:''"`
	GameRole      string `json:"game_role" form:"game_role" gorm:"type:varchar(50);not null;default:''"`
	GamePhone     string `json:"game_phone" form:"game_phone" gorm:"type:varchar(15);not null;default:''"`
	Margin        int    `json:"margin" form:"margin" gorm:"type:int(7);not null;default:0" valid:"Range(0, 9000000)"`
	AntiAddiction int    `json:"anti_addiction" form:"anti_addiction" gorm:"type:int(1);not null;default:0" valid:"Range(0, 1)"`
	DesignateHero int    `json:"designate_hero" form:"designate_hero" gorm:"type:int(1);not null;default:0" valid:"Range(0, 1)"`
	HeroName      string `json:"hero_name" form:"hero_name" gorm:"type:varchar(30);not null;default:''"`
	PayAmount     int    `json:"pay_amount" form:"-" gorm:"type:int(12);not null;default:0"`
	PayDesc       string `json:"pay_desc" form:"-" gorm:"type:varchar(100);not null;default:''"`
	PayIp         string `json:"pay_ip" form:"-" gorm:"type:varchar(30);not null;default:''"`
	TradeType     string `json:"trade_type" form:"-" gorm:"type:varchar(30);not null;default:''"`
	RegTime       int    `json:"reg_time" gorm:"type:int(12);not null;default:0"`
	UpdTime       int    `json:"upd_time" gorm:"type:int(12);not null;default:0"`
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

func (info *Order) GetOrderInfoByTradeNo() (bool, error) {
	err := db.Select("user_id, order_id, status").Where("trade_no = ?", info.TradeNo).Take(&info).Error
	if err != nil {
		return false, err
	}
	return true, nil
}

func GetUserIdAndStatusByOrderId(orderId int) (bool, error) {
	var article Article
	err := db.Select("user_id, status").Where("order_id = ?", orderId).First(&article).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}

	if article.ID > 0 {
		return true, nil
	}

	return false, nil
}

// Delete
//func (info *Order)Delete() bool {
//	//err := db.Where("id = ?", id).Delete(&Tag{}).Error
//	if err := db.Delete(info).Error; err != nil {
//		return false
//	}
//	return true
//}
