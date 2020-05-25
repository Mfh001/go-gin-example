package var_const

const (
	//用户类型
	UserTypeNormal    = 1 //普通用户
	UserTypeInstead   = 2 //代练用户
	UserTypeAccompany = 3 //陪练用户
	UserTypeAdmin     = 4 //管理员

	WXAppID  = "wx431aa31b8f263177"
	WXSecret = "7d9deb35e916e28974a30a45b6dd618e"
	//商户号
	WXMchID = "1554737721"
	//商户密钥
	WXMchKey = "t7v5TMsxhW6VH2f231NaB1BGL33CRjt3"

	//订单状态
	OrderStatusAddOrder                 = 0 //订单 未支付 状态
	OrderStatusWaitPay                  = 1
	OrderStatusPaidPay                  = 2  //订单 已支付下单 状态
	OrderStatusTakerWaitPay             = 3  //订单 接单准备支付
	OrderStatusTakerPaid                = 4  //订单 接单已支付
	OrderStatusTakerFinishedNeedConfirm = 5  //订单 代练已完成，请求确认 状态
	OrderStatusConfirmFinished          = 6  //订单 确认完成 状态
	OrderStatusRefundFinished           = 7  //订单 保证金已退 状态
	OrderStatusCancel                   = -1 //订单 取消 状态

	TeamCreate              = 0 //创建车队未支付
	TeamCanShow             = 1 //车队可展示出来，供用户加入或者代练接单
	TeamWorking             = 2 //已经发车
	TeamFinishedNeedConfirm = 3 //代练完成请求确认
	TeamConfirmed           = 4 //已确认完成

	TeamCardPrice      = 100
	TeamCardMax        = 5   //每次最多使用数量
	TeamTakerMargin    = 1   //车队代练支付的保证金
	TeamUrgentPrice    = 1   //加急费用
	RunesAddPriceLevel = 120 //铭文等级低于该等级，附加费用

	ChannelTypePlatform = 1     //平台频道
	OrderRate           = 10    //订单抽成费率%
	OrderNeedRate       = 3000  //订单金额大于等于 需要抽成
	OrderNeedRateMax    = 30000 //订单金额大于 不需要抽成
	ExchangeMinMoney    = 10000 //每次最少提现100元
	ExchangeRate        = 1     //订单抽成费率%

	AccessTokenExpireTime = 1.8 * 60 * 60
	Deposit               = 8000 //保证金80元

	OrderTotalTimesStatus100   = 1 //100已领取
	OrderTotalTimesStatus500   = 2 //500已领取
	OrderTotalTimesStatus1000  = 3 //1000已领取
	OrderTotalTimesStatus2000  = 4 //2000已领取
	OrderTotalTimesStatus10000 = 5 //10000已领取

	//满意度
	OrderSatisfied    = 1 //满意
	OrderOrdinary     = 2 //一般
	OrderDissatisfied = 3 //不满意

	//预定的时间段设置
	ShopBookPeriodIdxBegin = 1 //
	ShopBookPeriodIdxEnd   = 5 //

	//预定的桌位设置
	ShopBookSeatTypeBegin = 1 //
	ShopBookSeatTypeEnd   = 2 //

	//用户下单类型范围
	OrderTypeBegin = 1
	OrderTypeEnd   = 5
	//下单类型
	OrderTypeScanNoPay = 3

	//锁超时时间
	LockExpireTime = 2
	//订单自动评论时间
	OrderCommentTime = 2 * 60 * 60

	//审核状态
	CheckRefuse = -1
	CheckNeed   = 0
	CheckPass   = 1

	SMSCodeExpireTime = 300
)
