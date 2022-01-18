
type {{.FuncName}}Output struct {
}
func {{.FuncName}}(ctx *gin.Context) ({{.FuncName}}Output, error) {
    output := {{.FuncName}}Output{}
    return output, nil
}
