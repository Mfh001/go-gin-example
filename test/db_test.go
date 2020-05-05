package test

import (
	"github.com/EDDYCJY/go-gin-example/models"
	"testing"
)

func TestDBFind(t *testing.T) {
	user := models.User{
		UserId: 1,
	}
	user.First()
}
