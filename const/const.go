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
	OrderStatusAddOrder = 0 //订单 未支付 状态
	OrderStatusWaitPay  = 1
	OrderStatusPaidPay  = 2  //订单 已支付 状态
	OrderStatusSentPay  = 2  //订单 已送达 待评论 状态
	OrderStatusFinished = 3  //订单 已完成 状态
	OrderStatusCancel   = -1 //订单 取消 状态

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
)
