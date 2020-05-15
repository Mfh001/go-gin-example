package test

import (
	"github.com/EDDYCJY/go-gin-example/service/order_service"
	"testing"
)

func TestDBFind(t *testing.T) {
	//user := models.User{
	//	UserId: 1,
	//}
	//user.First()
	order_service.Refund(13)
}
