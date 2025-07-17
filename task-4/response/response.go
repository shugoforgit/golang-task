package response

import "github.com/gin-gonic/gin"

// 统一响应结构
type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

// 全局响应处理器
var Handler = &ResponseHandlerStruct{}

type ResponseHandlerStruct struct{}

// 成功响应
func (r *ResponseHandlerStruct) Success(c *gin.Context, data interface{}) {
	c.JSON(CodeSuccess, Response{
		Code: CodeSuccess,
		Msg:  "成功",
		Data: data,
	})
}

// 成功响应（无数据）
func (r *ResponseHandlerStruct) SuccessMsg(c *gin.Context, msg string) {
	c.JSON(CodeSuccess, Response{
		Code: CodeSuccess,
		Msg:  msg,
	})
}

// 错误响应
func (r *ResponseHandlerStruct) Error(c *gin.Context, code int, err error) {
	c.JSON(code, Response{
		Code: code,
		Msg:  err.Error(),
	})
}

// 常用错误响应函数
func (r *ResponseHandlerStruct) BadRequest(c *gin.Context, err error) {
	r.Error(c, CodeBadRequest, err)
}

func (r *ResponseHandlerStruct) Unauthorized(c *gin.Context, err error) {
	r.Error(c, CodeUnauthorized, err)
}

func (r *ResponseHandlerStruct) Forbidden(c *gin.Context, err error) {
	r.Error(c, CodeForbidden, err)
}

func (r *ResponseHandlerStruct) NotFound(c *gin.Context, err error) {
	r.Error(c, CodeNotFound, err)
}

func (r *ResponseHandlerStruct) InternalError(c *gin.Context, err error) {
	r.Error(c, CodeInternalError, err)
}

func (r *ResponseHandlerStruct) DatabaseError(c *gin.Context, err error) {
	r.Error(c, CodeDatabaseError, err)
}

func (r *ResponseHandlerStruct) ErrorMsg(c *gin.Context, code int, msg string) {
	c.JSON(code, Response{
		Code: code,
		Msg:  msg,
	})
}
