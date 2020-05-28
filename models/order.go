package models

import (
	"fmt"
	var_const "github.com/EDDYCJY/go-gin-example/const"
	"github.com/jinzhu/gorm"
)

type Order struct {
	OrderId       int    `json:"order_id" form:"-" gorm:"primary_key;type:int(12);not null"`
	OrderNo       string `json:"order_no" gorm:"type:varchar(50);not null;default:''"`
	Price         int    `json:"price" form:"price" gorm:"type:int(12);not null;default:0" valid:"Required;Range(0, 100000000)"`
	TimeLimit     int    `json:"time_limit" form:"time_limit" gorm:"type:int(12);not null;default:0" valid:"Required;Range(0, 1000000)"`
	TeamId        int    `json:"team_id" form:"-" gorm:"type:int(12);not null;default:0"`
	Status        int    `json:"status" form:"-" gorm:"type:int(12);not null;default:0"`
	StarNum       int    `json:"star_num" form:"-" gorm:"type:int(12);not null;default:0"`
	StarPerPrice  int    `json:"star_per_price" form:"-" gorm:"type:int(12);not null;default:0"`
	TradeNo       string `json:"trade_no" gorm:"type:varchar(50);not null;default:''"`
	UserId        int    `json:"user_id" form:"user_id" gorm:"type:int(12);not null" valid:"Required;Range(1, 1000000000)"`
	Title         string `json:"title" form:"title" gorm:"type:varchar(100);not null;default:''"`
	NickName      string `json:"nick_name" gorm:"type:varchar(32);not null;default:''"`
	GameType      int    `json:"game_type" form:"game_type" gorm:"type:int(2);not null" valid:"Range(0, 1)"`
	OrderType     int    `json:"order_type" form:"order_type" gorm:"type:int(2);not null" valid:"Range(0, 1)"`
	InsteadType   int    `json:"instead_type" form:"instead_type" gorm:"type:int(2);not null" valid:"Range(0, 2)"`
	GameZone      int    `json:"game_zone" form:"game_zone" gorm:"type:int(8);not null" valid:"Range(0, 4000)"`
	RunesLevel    int    `json:"runes_level" form:"runes_level" gorm:"type:int(4);not null" valid:"Range(0, 200)"`
	HeroNum       int    `json:"hero_num" form:"hero_num" gorm:"type:int(4);not null;default:0" valid:"Range(0, 1000)"`
	CurLevel      int    `json:"cur_level" form:"cur_level" gorm:"type:int(8);not null" valid:"Range(0, 1000000)"`
	TargetLevel   int    `json:"target_level" form:"target_level" gorm:"type:int(8);not null" valid:"Required;Range(1, 1000000)"`
	GameAcc       string `json:"game_acc" form:"game_acc" gorm:"type:varchar(60);not null;default:''"`
	GamePwd       string `json:"game_pwd" form:"game_pwd" gorm:"type:varchar(60);not null;default:''"`
	GameRole      string `json:"game_role" form:"game_role" gorm:"type:varchar(50);not null;default:''"`
	GamePhone     string `json:"game_phone" form:"game_phone" gorm:"type:varchar(15);not null;default:''"`
	Margin        int    `json:"margin" gorm:"type:int(7);not null;default:0"`
	MarginSafe    int    `json:"margin_safe" form:"margin_safe" gorm:"type:int(7);not null;default:0" valid:"Range(0, 9000000)"`
	MarginEff     int    `json:"margin_eff" form:"margin_eff" gorm:"type:int(7);not null;default:0" valid:"Range(0, 9000000)"`
	MarginArb     int    `json:"margin_arb" form:"margin_arb" gorm:"type:int(7);not null;default:0" valid:"Range(0, 9000000)"`
	AntiAddiction int    `json:"anti_addiction" form:"anti_addiction" gorm:"type:int(1);not null;default:0" valid:"Range(0, 1)"`
	DesignateHero int    `json:"designate_hero" form:"designate_hero" gorm:"type:int(1);not null;default:0" valid:"Range(0, 1)"`
	HeroName      string `json:"hero_name" form:"hero_name" gorm:"type:varchar(30);not null;default:''"`
	PayAmount     int    `json:"pay_amount" form:"-" gorm:"type:int(12);not null;default:0"`
	PayDesc       string `json:"pay_desc" form:"-" gorm:"type:varchar(100);not null;default:''"`
	PayIp         string `json:"pay_ip" form:"-" gorm:"type:varchar(30);not null;default:''"`
	TradeType     string `json:"trade_type" form:"-" gorm:"type:varchar(30);not null;default:''"`
	TransactionId string `json:"transaction_id" form:"-" gorm:"type:varchar(40);not null;default:''"`
	PayTime       int    `json:"pay_time" gorm:"type:int(12);not null;default:0"`

	ChannelType int `json:"channel_type" form:"-" gorm:"type:int(2);not null;default:0"`

	Contact     string `json:"contact" form:"contact" gorm:"type:varchar(30);not null;default:''"`
	Qq          string `json:"qq" form:"qq" gorm:"type:varchar(30);not null;default:''"`
	Description string `json:"description" form:"description" gorm:"type:varchar(1000);not null;default:''"`

	TakerTradeNo       string `json:"taker_trade_no" gorm:"type:varchar(50);not null;default:''"`
	TakerUserId        int    `json:"taker_user_id" form:"-" gorm:"type:int(12);not null;default:0"`
	TakerNickName      string `json:"taker_nick_name" gorm:"type:varchar(32);not null;default:''"`
	TakerPayAmount     int    `json:"taker_pay_amount" form:"-" gorm:"type:int(12);not null;default:0"`
	TakerPayDesc       string `json:"taker_pay_desc" form:"-" gorm:"type:varchar(100);not null;default:''"`
	TakerPayIp         string `json:"taker_pay_ip" form:"-" gorm:"type:varchar(30);not null;default:''"`
	TakerTradeType     string `json:"taker_trade_type" form:"-" gorm:"type:varchar(30);not null;default:''"`
	TakerTransactionId string `json:"taker_transaction_id" form:"-" gorm:"type:varchar(40);not null;default:''"`
	TakerTime          int    `json:"taker_time" gorm:"type:int(12);not null;default:0"`

	RefundTradeNo   string `json:"refund_trade_no" gorm:"type:varchar(50);not null;default:''"` //订单完成 保证金退还
	RefundAmount    int    `json:"refund_amount" form:"-" gorm:"type:int(12);not null;default:0"`
	TakerRefundTime int    `json:"taker_refund_time" gorm:"type:int(12);not null;default:0"`

	PayRefundTradeNo string `json:"pay_refund_trade_no" gorm:"type:varchar(50);not null;default:''"` //用户退款
	PayRefundAmount  int    `json:"pay_refund_amount" form:"-" gorm:"type:int(12);not null;default:0"`
	PayRefundTime    int    `json:"pay_refund_time" gorm:"type:int(12);not null;default:0"`

	TeamCardNum int `json:"team_card_num" form:"-" gorm:"type:int(12);not null;default:0"`
	RealPrice   int `json:"real_price" form:"-" gorm:"type:int(12);not null;default:0"`

	RegTime      int    `json:"reg_time" gorm:"type:int(12);not null;default:0"`
	UpdTime      int    `json:"upd_time" gorm:"type:int(12);not null;default:0"`
	ImgTakeUrl   string `gorm:"type:varchar(100);not null;default:''" json:"img_take_url"`
	ImgFinishUrl string `gorm:"type:varchar(100);not null;default:''" json:"img_finish_url"`
}

//insert
func (info Order) Insert() bool {
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
func (info Order) Updates(m map[string]interface{}) bool {
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
func (info Order) Save() bool {
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
func (info *Order) First() (int, error) {
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

func (info *Order) GetOrderInfoByTradeNo() (bool, error) {
	err := db.Select("user_id, order_id, status").Where("trade_no = ?", info.TradeNo).First(&info).Error
	if err != nil {
		return false, err
	}
	return true, nil
}

func (info *Order) GetOrderInfoByTakerTradeNo() (bool, error) {
	err := db.Select("taker_user_id, taker_pay_amount, order_id, status").Where("taker_trade_no = ?", info.TakerTradeNo).First(&info).Error
	if err != nil {
		return false, err
	}
	return true, nil
}
func (info *Order) GetOrderInfoByRefundTradeNo() (bool, error) {
	err := db.Select("taker_user_id, taker_pay_amount, order_id, status").Where("refund_trade_no = ?", info.RefundTradeNo).First(&info).Error
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

//select all
func GetNeedTakeOrders(infos *[]Order, where string, index int, count int) (bool, error) {
	err := db.Select("order_no, star_num, title, time_limit, star_per_price, channel_type, order_id, price, status, user_id, nick_name, game_type, description, order_type, instead_type, game_zone, runes_level, hero_num, cur_level, target_level, margin_safe, margin_eff, margin_arb, anti_addiction, designate_hero, hero_name, upd_time, contact, qq").Where("status = ? and team_id = 0 "+where, var_const.OrderStatusPaidPay).Limit(count).Offset(index).Find(&infos).Error
	if gorm.IsRecordNotFoundError(err) {
		*infos = []Order{}
		return true, nil
	} else if err != nil {
		*infos = []Order{}
		return false, err
	}
	return true, nil
}
func GetTakeOrders(takerId int, infos *[]Order, index int, count int) (bool, error) {
	err := db.Select("order_no, img_take_url, img_finish_url, star_num, title, time_limit, star_per_price, channel_type, order_id, price, status, user_id, nick_name, game_type, description, order_type, instead_type, game_zone, runes_level, hero_num, cur_level, target_level, margin_safe, margin_eff, margin_arb, anti_addiction, designate_hero, hero_name, upd_time, contact, qq").Where("status >= ? and taker_user_id = ?", var_const.OrderStatusTakerPaid, takerId).Limit(count).Offset(index).Find(&infos).Error
	if gorm.IsRecordNotFoundError(err) {
		*infos = []Order{}
		return true, nil
	} else if err != nil {
		*infos = []Order{}
		return false, err
	}
	return true, nil
}

func GetUserOrders(userId int, infos *[]Order, index int, count int) (bool, error) {
	err := db.Select("order_no, img_take_url, img_finish_url, star_num, title, time_limit, star_per_price, channel_type, order_id, price, status, user_id, nick_name, game_type, description, order_type, instead_type, game_zone, runes_level, hero_num, cur_level, target_level, margin_safe, margin_eff, margin_arb, anti_addiction, designate_hero, hero_name, upd_time, contact, qq").Where("status >= ? and user_id = ?", var_const.OrderStatusWaitPay, userId).Limit(count).Offset(index).Find(&infos).Error
	if gorm.IsRecordNotFoundError(err) {
		*infos = []Order{}
		return true, nil
	} else if err != nil {
		*infos = []Order{}
		return false, err
	}
	return true, nil
}

// Delete
//func (info *Order)Delete() bool {
//	//err := db.Where("id = ?", id).Delete(&Tag{}).Error
//	if err := db.Delete(info).Error; err != nil {
//		return false
//	}
//	return true
//}
