package models

import (
	"github.com/jinzhu/gorm"
)

type Team struct {
	TeamId      int    `json:"team_id" form:"-" gorm:"primary_key;type:int(12);not null"`
	Price       int    `json:"price" form:"-" gorm:"type:int(12);not null;default:0"`
	Status      int    `json:"status" form:"-" gorm:"type:int(12);not null;default:0"`
	OwnerId     int    `json:"owner_id" form:"owner_id" gorm:"type:int(12);not null" valid:"Required;Range(1, 1000000000)"`
	OwnerType   int    `json:"owner_type" form:"-" gorm:"type:int(2);not null;default:1"`
	NickName    string `json:"nick_name" gorm:"type:varchar(32);not null;default:''"`
	GameType    int    `json:"game_type" form:"game_type" gorm:"type:int(2);not null" valid:"Range(0, 1)"`
	BigZone     int    `json:"big_zone" form:"big_zone" gorm:"type:int(8);not null" valid:"Range(0, 10)"`
	CurLevel    int    `json:"cur_level" form:"cur_level" gorm:"type:int(8);not null" valid:"Range(0, 10000)"`
	TargetLevel int    `json:"target_level" form:"target_level" gorm:"type:int(8);not null" valid:"Required;Range(1, 10000)"`

	NeedPwd int    `json:"need_pwd" form:"need_pwd" gorm:"type:int(2);not null;default:0"`
	Pwd     string `json:"pwd" form:"pwd" gorm:"type:varchar(32);not null;default:''"`

	NeedNum      int    `json:"need_num" form:"need_num" gorm:"type:int(12);not null;default:1"`
	Num          int    `json:"num" form:"-" gorm:"type:int(12);not null;default:0"`
	UserId1      int    `json:"user_id1" form:"-" gorm:"type:int(12);not null;default:0"`
	NickName1    string `json:"nick_name1" gorm:"type:varchar(32);not null;default:''"`
	OrderStatus1 int    `json:"order_status1" form:"-" gorm:"type:int(2);not null;default:0"`
	User1PayTime int    `json:"user1_pay_time" gorm:"type:int(12);not null;default:0"`
	OrderId1     int    `json:"order_id1" form:"-" gorm:"type:int(12);not null;default:0"`
	PayAmount1   int    `json:"pay_amount1" form:"-" gorm:"type:int(12);not null;default:0"`
	UserId2      int    `json:"user_id2" form:"-" gorm:"type:int(12);not null;default:0"`
	NickName2    string `json:"nick_name2" gorm:"type:varchar(32);not null;default:''"`
	OrderStatus2 int    `json:"order_status2" form:"-" gorm:"type:int(2);not null;default:0"`
	User2PayTime int    `json:"user2_pay_time" gorm:"type:int(12);not null;default:0"`
	OrderId2     int    `json:"order_id2" form:"-" gorm:"type:int(12);not null;default:0"`
	PayAmount2   int    `json:"pay_amount2" form:"-" gorm:"type:int(12);not null;default:0"`

	Contact     string `json:"contact" form:"contact" gorm:"type:varchar(30);not null;default:''"`
	Qq          string `json:"qq" form:"qq" gorm:"type:varchar(30);not null;default:''"`
	Description string `json:"description" form:"description" gorm:"type:varchar(1000);not null;default:''"`

	TakerTradeNo       string `json:"taker_trade_no" gorm:"type:varchar(50);not null;default:''"`
	TakerPayStatus     int    `json:"taker_pay_status" form:"-" gorm:"type:int(2);not null;default:0"`
	TakerUserId        int    `json:"taker_user_id" form:"-" gorm:"type:int(12);not null;default:0"`
	TakerNickName      string `json:"taker_nick_name" gorm:"type:varchar(32);not null;default:''"`
	TakerPayAmount     int    `json:"taker_pay_amount" form:"-" gorm:"type:int(12);not null;default:0"`
	TakerPayDesc       string `json:"taker_pay_desc" form:"-" gorm:"type:varchar(100);not null;default:''"`
	TakerPayIp         string `json:"taker_pay_ip" form:"-" gorm:"type:varchar(30);not null;default:''"`
	TakerTradeType     string `json:"taker_trade_type" form:"-" gorm:"type:varchar(30);not null;default:''"`
	TakerTransactionId string `json:"taker_transaction_id" form:"-" gorm:"type:varchar(40);not null;default:''"`
	TakerTime          int    `json:"taker_time" gorm:"type:int(12);not null;default:0"`

	RefundTradeNo   string `json:"refund_trade_no" gorm:"type:varchar(50);not null;default:''"`
	RefundAmount    int    `json:"refund_amount" form:"-" gorm:"type:int(12);not null;default:0"`
	TakerRefundTime int    `json:"taker_refund_time" gorm:"type:int(12);not null;default:0"`

	IsUrgent        int    `json:"is_urgent" form:"-" gorm:"type:int(2);not null;default:0"`      //是否加急
	UrgentUserId    int    `json:"urgent_user_id" form:"-" gorm:"type:int(2);not null;default:0"` //哪个用户加急
	UrgentNickName  string `json:"urgent_nick_name" gorm:"type:varchar(32);not null;default:''"`
	UrgentTradeNo   string `json:"urgent_trade_no" gorm:"type:varchar(50);not null;default:''"`
	UrgentPayAmount int    `json:"urgent_pay_amount" form:"-" gorm:"type:int(12);not null;default:0"`
	UrgentPayTime   int    `json:"urgent_pay_time" gorm:"type:int(12);not null;default:0"`

	UrgentRefundTradeNo string `json:"urgent_refund_trade_no" gorm:"type:varchar(50);not null;default:''"`
	UrgentRefundAmount  int    `json:"urgent_refund_amount" form:"-" gorm:"type:int(12);not null;default:0"`
	UrgentRefundTime    int    `json:"urgent_refund_time" gorm:"type:int(12);not null;default:0"`

	RegTime int `json:"reg_time" gorm:"type:int(12);not null;default:0"`
	UpdTime int `json:"upd_time" gorm:"type:int(12);not null;default:0"`

	TeamCardNum int `json:"team_card_num" form:"team_card_num" gorm:"-" valid:"Range(0, 100)"`
}

//insert
func (info Team) Insert() bool {
	create := db.Create(&info)
	if create.Error != nil {
		return false
	}
	return true
}

//update
func (info Team) Updates(m map[string]interface{}) bool {
	err := db.Model(&info).Updates(m).Error
	if err != nil {
		return false
	}
	return true
}

//insert and update
func (info Team) Save() bool {
	create := db.Save(&info)
	if create.Error != nil {
		return false
	}
	return true
}

//select
func (info *Team) First() (int, error) {
	err := db.First(&info).Error
	if gorm.IsRecordNotFoundError(err) {
		return 0, nil
	} else if err != nil {
		return -1, err
	}
	return 1, nil
}

//func (info *Team) GetOrderInfoByTradeNo() (bool, error) {
//	err := db.Select("user_id, order_id, status").Where("trade_no = ?", info.TradeNo).First(&info).Error
//	if err != nil {
//		return false, err
//	}
//	return true, nil
//}

func (info *Team) GetOrderInfoByTakerTradeNo() (bool, error) {
	err := db.Select("owner_type, order_status1, taker_user_id, taker_pay_amount, team_id, status").Where("taker_trade_no = ?", info.TakerTradeNo).First(&info).Error
	if err != nil {
		return false, err
	}
	return true, nil
}
func (info *Team) GetOrderInfoByRefundTradeNo() (bool, error) {
	err := db.Select("taker_user_id, taker_pay_amount, order_id, status").Where("refund_trade_no = ?", info.RefundTradeNo).First(&info).Error
	if err != nil {
		return false, err
	}
	return true, nil
}

//select all
//func GetNeedTakeOrders(infos *[]Order) (bool, error) {
//	err := db.Select("order_id, price, status, user_id, nick_name, game_type, description, order_type, instead_type, game_zone, runes_level, hero_num, cur_level, target_level, margin, anti_addiction, designate_hero, hero_name, upd_time, contact, qq").Where("status = ?", var_const.OrderStatusPaidPay).Find(&infos).Error
//	if gorm.IsRecordNotFoundError(err) {
//		*infos = []Order{}
//		return true, nil
//	} else if err != nil {
//		*infos = []Order{}
//		return false, err
//	}
//	return true, nil
//}
//func GetTakeOrders(takerId int, infos *[]Order) (bool, error) {
//	err := db.Select("*").Where("status >= ? and taker_user_id = ?", var_const.OrderStatusTakerPaid, takerId).Find(&infos).Error
//	if gorm.IsRecordNotFoundError(err) {
//		*infos = []Order{}
//		return true, nil
//	} else if err != nil {
//		*infos = []Order{}
//		return false, err
//	}
//	return true, nil
//}
//
//func GetUserOrders(userId int, infos *[]Order) (bool, error) {
//	err := db.Select("*").Where("status >= ? and user_id = ?", var_const.OrderStatusWaitPay, userId).Find(&infos).Error
//	if gorm.IsRecordNotFoundError(err) {
//		*infos = []Order{}
//		return true, nil
//	} else if err != nil {
//		*infos = []Order{}
//		return false, err
//	}
//	return true, nil
//}

// Delete
//func (info *Order)Delete() bool {
//	//err := db.Where("id = ?", id).Delete(&Tag{}).Error
//	if err := db.Delete(info).Error; err != nil {
//		return false
//	}
//	return true
//}
