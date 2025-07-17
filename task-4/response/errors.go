package response

import "errors"

// HTTP 状态码常量
const (
	CodeSuccess            = 200
	CodeBadRequest         = 400
	CodeUnauthorized       = 401
	CodeForbidden          = 403
	CodeNotFound           = 404
	CodePostNotFound       = 405
	CodePostUpdateError    = 406
	CodeCommentCreateError = 407
	CodeInternalError      = 500
	CodeDatabaseError      = 501
)

var (
	ErrUserNotFound          = errors.New("用户不存在")
	ErrUserAlreadyExists     = errors.New("用户已存在")
	ErrUserPasswordIncorrect = errors.New("密码错误")
	ErrUserNotLogin          = errors.New("用户未登录")
	ErrUserNotAuthorized     = errors.New("无权限")
	ErrMissingParams         = errors.New("缺少参数")
	ErrInvalidToken          = errors.New("无效的token")
	ErrDatabaseError         = errors.New("数据库错误")
	ErrPostNotFound          = errors.New("文章不存在")
	ErrPostUpdateError       = errors.New("文章更新失败")
	ErrCommentCreateError    = errors.New("评论创建失败")
)
