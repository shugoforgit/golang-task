# 博客系统 API

基于 Go + Gin + GORM + MySQL 的博客系统后端 API，支持用户注册登录、文章管理、评论功能，使用 JWT 进行身份认证。

## 功能特性

- 🔐 JWT 身份认证
- 👤 用户注册/登录
- 📝 文章 CRUD 操作
- 💬 评论功能
- 📊 用户文章统计
- 📝 统一日志记录
- 🎯 统一响应格式
- 🔒 权限控制

## 技术栈

- **语言**: Go 1.22.2
- **Web框架**: Gin
- **ORM**: GORM
- **数据库**: MySQL
- **认证**: JWT
- **密码加密**: bcrypt

## 环境要求

- Go 1.22.2 或更高版本
- MySQL 8.0 或更高版本
- Git

## 快速开始

### 1. 克隆项目

```bash
git clone <repository-url>
cd task-4
```

### 2. 安装依赖

```bash
go mod download
```

### 3. 配置数据库

#### 3.1 创建数据库

```sql
CREATE DATABASE blog_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

#### 3.2 配置数据库连接

编辑 `main.go` 文件中的数据库连接字符串：

```go
var dsn = "用户名:密码@tcp(127.0.0.1:3306)/blog_db?charset=utf8mb4&parseTime=True&loc=Local"
```

请将 `用户名` 和 `密码` 替换为您的 MySQL 用户名和密码。

### 4. 运行项目

```bash
# 编译并运行
go run .

# 或者先编译再运行
go build -o main .
./main
```

服务将在 `http://localhost:8080` 启动。

### 5. 测试API

项目提供了完整的API测试用例，您可以选择以下任一方式测试：

#### 使用 Apifox
1. 下载并安装 [Apifox](https://apifox.com/), 或者使用网页版(需要安装插件才能调用本地接口)
2. 导入项目根目录下的 `个人博客.Apifox.json` 文件
3. 运行注册和登录接口获取JWT Token
4. 使用其他接口进行测试

#### 使用 Postman
1. 下载并安装 [Postman](https://www.postman.com/)
2. 导入项目根目录下的 `个人博客.postman.json` 文件
3. 运行注册和登录接口获取JWT Token
4. 使用其他接口进行测试

#### 使用 curl 命令行
```bash
# 注册用户
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"password123","email":"test@example.com"}'

# 登录获取Token
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"password123"}'

# 使用Token创建文章
curl -X POST http://localhost:8080/auth/create-post \
  -H "Content-Type: application/json" \
  -H "Authorization: YOUR_JWT_TOKEN" \
  -d '{"title":"测试文章","content":"这是测试内容"}'
```

## 项目结构

```
task-4/
├── main.go                    # 主程序入口
├── go.mod                     # Go 模块文件
├── go.sum                     # 依赖校验文件
├── README.md                  # 项目文档
├── 个人博客.Apifox.json       # Apifox 测试用例
├── 个人博客.postman.json      # Postman 测试用例
├── logger/                    # 日志模块
│   └── logger.go             # 日志工具
├── response/                  # 响应处理模块
│   ├── response.go           # 统一响应格式
│   └── errors.go             # 错误常量定义
└── log/                      # 日志文件目录
    ├── info.log              # 信息日志
    └── error.log             # 错误日志
```

## API 接口文档

### 基础信息

- **Base URL**: `http://localhost:8080`
- **Content-Type**: `application/json`
- **认证方式**: JWT Token (在请求头中添加 `Authorization: <token>`)

### 测试工具

项目提供了完整的API测试用例，可以直接导入到测试工具中使用：

#### Apifox 导入
1. 打开 [Apifox](https://apifox.com/), 或者使用网页版(需要安装插件才能调用本地接口)
2. 创建新项目或打开现有项目
3. 点击 "导入" → "从文件导入"
4. 选择项目根目录下的 `个人博客.Apifox.json` 文件
5. 导入后即可使用所有预配置的API测试用例

#### Postman 导入
1. 打开 [Postman](https://www.postman.com/)
2. 点击 "Import" 按钮
3. 选择 "File" 标签页
4. 选择项目根目录下的 `个人博客.postman.json` 文件
5. 点击 "Import" 完成导入
6. 导入后可在 Collections 中找到 "个人博客" 集合

#### 使用说明
- 测试用例已预配置所有必要的请求参数
- JWT Token 需要在登录后手动更新到环境变量中
- 建议先运行注册和登录接口获取Token
- 所有需要认证的接口都会自动使用环境变量中的Token

### 统一响应格式

```json
{
  "code": 200,
  "msg": "成功",
  "data": {}
}
```

### 1. 用户注册

**接口**: `POST /register`

**请求体**:
```json
{
  "username": "testuser",
  "password": "password123",
  "email": "test@example.com"
}
```

**响应**:
```json
{
  "code": 200,
  "msg": "注册成功"
}
```

### 2. 用户登录

**接口**: `POST /login`

**请求体**:
```json
{
  "username": "testuser",
  "password": "password123"
}
```

**响应**:
```json
{
  "code": 200,
  "msg": "成功",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

### 3. 创建文章

**接口**: `POST /auth/create-post`

**认证**: 需要 JWT Token

**请求体**:
```json
{
  "title": "我的第一篇文章",
  "content": "这是文章内容..."
}
```

**响应**:
```json
{
  "code": 200,
  "msg": "创建成功"
}
```

### 4. 获取所有文章

**接口**: `GET /posts`

**响应**:
```json
{
  "code": 200,
  "msg": "成功",
  "data": [
    {
      "id": 1,
      "title": "我的第一篇文章",
      "content": "这是文章内容...",
      "user_id": 1,
      "comment_status": "无评论",
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  ]
}
```

### 5. 获取单个文章

**接口**: `GET /posts/:id`

**响应**:
```json
{
  "code": 200,
  "msg": "成功",
  "data": {
    "id": 1,
    "title": "我的第一篇文章",
    "content": "这是文章内容...",
    "user_id": 1,
    "comment_status": "无评论",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

### 6. 更新文章

**接口**: `POST /auth/update-post`

**认证**: 需要 JWT Token

**请求体**:
```json
{
  "id": 1,
  "title": "更新后的标题",
  "content": "更新后的内容..."
}
```

**响应**:
```json
{
  "code": 200,
  "msg": "更新成功"
}
```

### 7. 删除文章

**接口**: `POST /auth/delete-post`

**认证**: 需要 JWT Token

**请求体**:
```json
{
  "id": 1
}
```

**响应**:
```json
{
  "code": 200,
  "msg": "删除成功"
}
```

### 8. 创建评论

**接口**: `POST /auth/create-comment`

**认证**: 需要 JWT Token

**请求体**:
```json
{
  "post_id": 1,
  "content": "这是一条评论"
}
```

**响应**:
```json
{
  "code": 200,
  "msg": "创建成功"
}
```

### 9. 获取文章评论

**接口**: `GET /comments/:post_id`

**响应**:
```json
{
  "code": 200,
  "msg": "成功",
  "data": [
    {
      "id": 1,
      "user_id": 1,
      "content": "这是一条评论",
      "post_id": 1,
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  ]
}
```

## 错误码说明

| 错误码 | 说明 |
|--------|------|
| 200 | 成功 |
| 400 | 请求参数错误 |
| 401 | 未授权/未登录 |
| 403 | 无权限 |
| 404 | 资源不存在 |
| 500 | 服务器内部错误 |
| 501 | 数据库错误 |

## 数据库表结构

### users 表
- `id`: 主键
- `username`: 用户名（唯一）
- `password`: 密码（加密存储）
- `email`: 邮箱
- `post_count`: 文章数量
- `created_at`: 创建时间
- `updated_at`: 更新时间

### posts 表
- `id`: 主键
- `title`: 文章标题
- `content`: 文章内容
- `user_id`: 作者ID（外键）
- `comment_status`: 评论状态
- `created_at`: 创建时间
- `updated_at`: 更新时间

### comments 表
- `id`: 主键
- `user_id`: 评论者ID（外键）
- `content`: 评论内容
- `post_id`: 文章ID（外键，级联删除）
- `created_at`: 创建时间
- `updated_at`: 更新时间

## 日志系统

项目使用自定义日志系统，日志文件存储在 `log/` 目录下：

- `info.log`: 信息日志
- `error.log`: 错误日志

日志按日期自动分割，便于问题排查。

## 安全特性

1. **密码加密**: 使用 bcrypt 进行密码哈希
2. **JWT 认证**: 基于 JWT 的身份认证
3. **权限控制**: 用户只能操作自己的文章
4. **参数验证**: 严格的请求参数验证
5. **SQL 注入防护**: 使用 GORM 参数化查询

## 开发说明

### 添加新的错误码

在 `response/errors.go` 中添加新的错误常量：

```go
const (
    CodeNewError = 502
)

var (
    ErrNewError = errors.New("新错误描述")
)
```

### 添加新的响应方法

在 `response/response.go` 中添加新的响应方法：

```go
func (r *ResponseHandlerStruct) NewError(c *gin.Context, err error) {
    r.Error(c, CodeNewError, err)
}
```

## 故障排除

### 1. 数据库连接失败

- 检查 MySQL 服务是否启动
- 验证数据库连接字符串是否正确
- 确认数据库用户权限

### 2. JWT 认证失败

- 检查请求头是否包含 `Authorization` 字段
- 验证 JWT Token 是否有效
- 确认 Token 是否过期

### 3. 外键约束错误

- 确保数据库表已正确创建
- 检查外键约束是否正确设置
- 验证级联删除是否生效

## 许可证

MIT License

## 贡献

欢迎提交 Issue 和 Pull Request！ 