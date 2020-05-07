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

	apiV1 := r.Group("/api/v1")
	apiV1.Use(jwt.JWT())
	{
		//代练提交段位审核
		apiV1.POST("/check", v1.AddCheck)
		//用户获取提交的审核信息
		apiV1.GET("/check", v1.GetCheckInfo)
		//管理员获取审核列表
		apiV1.GET("/check/admin", v1.GetAdminChecks)
		//管理员进行审核
		apiV1.PUT("/check/admin/:id", v1.AdminCheck)

		//下单
		apiV1.POST("/order", v1.AddOrder)

		//pay
		apiV1.POST("/pay", v1.WxPay)

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
