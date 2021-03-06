package routers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	_ "github.com/EDDYCJY/go-gin-example/docs"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"

	"github.com/EDDYCJY/go-gin-example/middleware/jwt"
	"github.com/EDDYCJY/go-gin-example/pkg/export"
	"github.com/EDDYCJY/go-gin-example/pkg/qrcode"
	"github.com/EDDYCJY/go-gin-example/pkg/upload"
	"github.com/EDDYCJY/go-gin-example/routers/api"
	"github.com/EDDYCJY/go-gin-example/routers/api/v1"
)

// InitRouter initialize routing information
func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.StaticFS("/export", http.Dir(export.GetExcelFullPath()))
	r.StaticFS("/upload/images", http.Dir(upload.GetImageFullPath()))
	r.StaticFS("/qrcode", http.Dir(qrcode.GetQrCodeFullPath()))

	r.POST("/wxlogin", api.WXLogin)
	r.POST("/login", api.Login)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.POST("/upload", api.UploadImage)
	r.Any("/pay/notify", v1.WxNotify)
	//r.Any("/pay/taker/notify", v1.TakerWxNotify)
	//r.Any("/pay/taker/refundnotify", v1.WxRefundCallback)
	//r.Any("/pay/teamnotify", v1.TeamWxNotify)
	//r.Any("/pay/urgentteamnotify", v1.TeamUrgentWxNotify)
	//r.Any("/pay/team/urgentrefundnotify", v1.UrgentRefundCallback)
	//r.Any("/pay/taker/depositnotify", v1.DepositWxNotify)
	//r.Any("/pay/order/refundnotify", v1.OrderCancelRefundCallback)
	//r.Any("/pay/order/undo/userrefundnotify", v1.OrderUndoUserRefundCallback)
	//r.Any("/pay/order/undo/takerrefundnotify", v1.OrderUndoTakerRefundCallback)
	//r.Any("/pay/order/admin/userrefundnotify", v1.AdminRefundUserCallback)
	//r.Any("/pay/order/admin/takerrefundnotify", v1.AdminRefundTakerCallback)

	//管理员获取审核列表
	r.GET("/check/admin", v1.GetAdminChecks)
	r.GET("/exchange/all", v1.GetAdminExchanges)
	//管理员进行审核
	r.PUT("/check/admin/:id", v1.AdminCheck)
	r.GET("/exchange/check", v1.ExchangeCheck)
	r.POST("/exchange/check", v1.ExchangeCheckPost)
	r.GET("/exchange/bank", v1.GetExchangeBank)
	r.POST("/admin/order/refund/user", v1.AdminRefundPay)
	//r.POST("/admin/order/refund/taker", v1.AdminRefundTaker)
	r.POST("/admin/money/add", v1.AdminAddBalance)
	r.POST("/admin/money/remove", v1.AdminRemBalance)
	r.POST("/admin/margin/remove", v1.AdminRemMargin)
	r.POST("/admin/order/adjudge/all", v1.GetAdminAdjudgeList)
	r.POST("/admin/user/all", v1.GetAdminUserList)
	r.POST("/admin/user/canpublish", v1.AdminSetUserCanPublish)

	apiV1 := r.Group("/api/v1")
	apiV1.Use(jwt.JWT())
	{
		//代练提交段位审核
		apiV1.POST("/check", v1.AddCheck)
		//用户获取提交的审核信息
		apiV1.GET("/check", v1.GetCheckInfo)

		apiV1.GET("/orderaward", v1.GetUserTotalOrderTimesAward)

		//订单
		apiV1.POST("/order", v1.AddOrder)

		apiV1.GET("/order/all", v1.GetAllOrders)
		apiV1.GET("/order/sortall", v1.GetAllOrdersB)
		apiV1.GET("/order/takelist", v1.GetTakerOrders)
		apiV1.GET("/order/userlist", v1.GetUserOrders)
		apiV1.POST("/order/message", v1.MessageOrder)
		apiV1.POST("/order/getmessage", v1.GetOrderMessage)
		apiV1.POST("/order/undorequest", v1.UndoRequestOrder)
		apiV1.POST("/order/undo", v1.UndoOrder)
		apiV1.POST("/order/addtime", v1.AddTimeOrder)
		apiV1.POST("/order/info", v1.GetOrderInfo)
		apiV1.POST("/order/adjudgerequest", v1.AdjudgeRequestOrder)

		//完成订单
		apiV1.POST("/order/finish", v1.FinishOrder)
		//确认完成订单
		apiV1.POST("/order/confirm", v1.ConfirmOrder)
		apiV1.POST("/order/updorderimg", v1.UpdateOrderImg)

		apiV1.POST("/order/cancel", v1.CancelOrder)
		apiV1.POST("/order/canceltake", v1.CancelTakeOrder)
		//pay
		apiV1.POST("/pay/deposit", v1.DepositWxPay)
		apiV1.POST("/pay/order", v1.AddPayOrder)
		//apiV1.POST("/pay", v1.Pay)

		apiV1.POST("/order/take", v1.TakeOrder)

		//绑定手机号
		apiV1.POST("/phone/bind", v1.BindPhone)
		apiV1.GET("/phone/code", v1.GetPhoneRegCode)
		apiV1.GET("/phone/code2", v1.GetPhoneCodeNoPhone)

		//绑定银行卡
		apiV1.POST("/bank/bind", v1.BindBankCard)
		apiV1.GET("/bank", v1.GetBankCardInfo)

		apiV1.GET("/balance", v1.GetUserBalance)
		apiV1.GET("/totalordertimes", v1.GetUserTotalOrderTimes)

		apiV1.POST("/user/getmessage", v1.GetUserMessage)
		apiV1.POST("/user/msgnoread", v1.GetUserMessageNoRead)

		//车队
		//apiV1.POST("/team", v1.AddTeam)
		apiV1.POST("/team/joincheck", v1.JoinTeamCheckPwd)
		//apiV1.POST("/team/join", v1.JoinTeam)
		//apiV1.POST("/teampay", v1.TeamWxPay)
		apiV1.POST("/team/urgent", v1.Urgent)
		apiV1.GET("/team/list", v1.GetAllTeams)
		apiV1.POST("/team/cancelurgent", v1.CancelUrgent)

		//代理
		apiV1.POST("/agent/bind", v1.BindAgent)
		apiV1.GET("/agent/profit", v1.GetAgentProfit)
		//获取二维码
		apiV1.POST("/qrcode", v1.QRcodeGet)

		//提现
		apiV1.POST("/exchange", v1.AddExchange)

		////获取标签列表
		//apiV1.GET("/tags", v1.GetTags)
		////新建标签
		//apiV1.POST("/tags", v1.AddTag)
		////更新指定标签
		//apiV1.PUT("/tags/:id", v1.EditTag)
		////删除指定标签
		//apiV1.DELETE("/tags/:id", v1.DeleteTag)
		////导出标签
		//r.POST("/tags/export", v1.ExportTag)
		////导入标签
		//r.POST("/check/import", v1.ImportTag)
		//
		////获取文章列表
		//apiV1.GET("/articles", v1.GetArticles)
		////获取指定文章
		//apiV1.GET("/articles/:id", v1.GetArticle)
		////新建文章
		//apiV1.POST("/articles", v1.AddArticle)
		////更新指定文章
		//apiV1.PUT("/articles/:id", v1.EditArticle)
		////删除指定文章
		//apiV1.DELETE("/articles/:id", v1.DeleteArticle)
		////生成文章海报
		//apiV1.POST("/articles/poster/generate", v1.GenerateArticlePoster)
	}

	return r
}
