package controller

import (
	"demo/models"
	"demo/resp"

	"github.com/gin-gonic/gin"
)

// User ...
func User(ctx *gin.Context) {
	user := GetUserWithGin(ctx)
	resp.Ok(ctx, user)
}

// GetUserWithGin 通过gin的context取userinfo
func GetUserWithGin(ctx *gin.Context) *models.User {
	return ctx.MustGet("user").(*models.User)
}
