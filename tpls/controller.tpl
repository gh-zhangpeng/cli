
func {{.FuncName}}(ctx *gin.Context) {
	var input struct {
	}
	if err := ctx.ShouldBind(&input); err != nil {
		return
	}
}
