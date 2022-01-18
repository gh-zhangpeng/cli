package controller

import (
	"github.com/gin-gonic/gin"
)

func CreateUser(ctx *gin.Context) {
	var input struct {
	}
	if err := ctx.ShouldBind(&input); err != nil {
		return
	}
}

func CreateUser(ctx *gin.Context) {
	var input struct {
	}
	if err := ctx.ShouldBind(&input); err != nil {
		return
	}
}

func CreateUser(ctx *gin.Context) {
	var input struct {
	}
	if err := ctx.ShouldBind(&input); err != nil {
		return
	}
}
