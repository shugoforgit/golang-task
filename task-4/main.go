package main

import (
	"errors"
	"fmt"
	"task-4/logger"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"task-4/response"
)

// * 数据库表结构
type User struct {
	ID        int       `gorm:"primaryKey"`
	Username  string    `gorm:"type:varchar(100);not null;unique"`
	Password  string    `gorm:"type:varchar(100);not null"`
	Email     string    `gorm:"type:varchar(100);not null"`
	PostCount int       `gorm:"default:0"`
	Posts     []Post    `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Comments  []Comment `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Post struct {
	ID            int       `gorm:"primaryKey"`
	Title         string    `gorm:"type:varchar(100);not null"`
	Content       string    `gorm:"type:text;not null"`
	UserID        int       `gorm:"not null"`
	Comments      []Comment `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE"`
	CommentStatus string    `gorm:"type:varchar(20);default:'无评论'"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type Comment struct {
	ID        int    `gorm:"primaryKey"`
	UserID    int    `gorm:"not null"`
	Content   string `gorm:"type:text;not null"`
	PostID    int    `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

var dsn = "root:shugoformysql@tcp(127.0.0.1:3306)/blog_db?charset=utf8mb4&parseTime=True&loc=Local"
var db *gorm.DB
var jwtkey = []byte("CQYHSfndtfB8ZtGZjPb7ceD3R57CcX9Y")

func init() {
	// 初始化日志
	if err := logger.Init(); err != nil {
		panic("初始化日志失败: " + err.Error())
	}
	logger.Info("日志系统初始化成功")

	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Fatal("数据库连接失败: %v", err)
		panic("数据库连接失败: " + err.Error())
	}
	logger.Info("数据库连接成功")

	db.AutoMigrate(&User{}, &Post{}, &Comment{})
	logger.Info("数据库表迁移完成")
}

func main() {
	router := gin.Default()
	registerRoutes(router)
	router.Run(":8080")
}

func registerRoutes(router *gin.Engine) {
	router.POST("/register", registerUser)
	router.POST("/login", loginUser)

	auth := router.Group("/auth", JWTAuthMiddleware())
	// * 文章操作
	auth.POST("/create-post", createPost) // * 需要认证
	auth.POST("/update-post", updatePost)
	auth.POST("/delete-post", deletePost)
	router.GET("/posts", getPosts)
	router.GET("/posts/:id", getPost)
	// * 评论操作
	auth.POST("/create-comment", createComment)
	router.GET("/comments/:post_id", getComments)
}

func registerUser(c *gin.Context) {
	type UserRegister struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	var userInfo UserRegister
	if err := c.ShouldBindJSON(&userInfo); err != nil {
		logger.Error("用户注册参数解析失败: %v", err)
		response.Handler.BadRequest(c, response.ErrMissingParams)
		return
	}

	// * 检查用户是否存在
	var userExist User
	res := db.Select("username").Where("username = ?", userInfo.Username).First(&userExist).Error
	if !errors.Is(res, gorm.ErrRecordNotFound) {
		logger.Error("用户注册失败，用户名已存在: %s", userInfo.Username)
		response.Handler.BadRequest(c, response.ErrUserAlreadyExists)
		return
	}

	passHash, err := HashPassword(userInfo.Password)
	if err != nil {
		logger.Error("密码加密失败: %v", err)
		response.Handler.BadRequest(c, response.ErrUserPasswordIncorrect)
		return
	}

	var user User
	user.Username = userInfo.Username
	user.Password = passHash
	user.Email = userInfo.Email

	res = db.Create(&user).Error
	if res != nil {
		logger.Error("用户创建失败: %v", user.Username)
		response.Handler.DatabaseError(c, response.ErrDatabaseError)
		return
	}

	response.Handler.SuccessMsg(c, "注册成功")
}

func loginUser(c *gin.Context) {
	type UserLogin struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	var userInfo UserLogin
	if err := c.ShouldBindJSON(&userInfo); err != nil {
		logger.Error("用户登录参数解析失败: %v %v", userInfo.Username, err)
		response.Handler.BadRequest(c, response.ErrMissingParams)
		return
	}

	var user User
	db.Where("username = ?", userInfo.Username).First(&user)
	if user.ID == 0 {
		logger.Error("用户登录失败，用户不存在: %s", userInfo.Username)
		response.Handler.BadRequest(c, response.ErrUserNotFound)
		return
	}

	if !CheckPasswordHash(userInfo.Password, user.Password) {
		logger.Error("用户登录失败，密码错误: username=%s", userInfo.Username)
		response.Handler.BadRequest(c, response.ErrUserPasswordIncorrect)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       user.ID,
		"username": user.Username,
		"iat":      time.Now().Unix(),
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(jwtkey)
	if err != nil {
		logger.Error("JWT生成失败: %v", err)
		response.Handler.BadRequest(c, response.ErrInvalidToken)
		return
	}

	response.Handler.Success(c, gin.H{
		"token": tokenString,
	})
}

// * 创建文章
func createPost(c *gin.Context) {
	type PostCreate struct {
		UserID  int    `json:"user_id"`
		Title   string `json:"title"`
		Content string `json:"content"`
	}

	var postInfo PostCreate
	if err := c.ShouldBindJSON(&postInfo); err != nil {
		logger.Error("创建文章参数解析失败: %v", err)
		response.Handler.BadRequest(c, response.ErrMissingParams)
		return
	}

	var post Post
	post.Title = postInfo.Title
	post.Content = postInfo.Content
	post.UserID = c.GetInt("user_id")

	res := db.Create(&post).Error
	if res != nil {
		logger.Error("db 创建文章失败: %v", postInfo)
		response.Handler.DatabaseError(c, response.ErrDatabaseError)
		return
	}

	response.Handler.SuccessMsg(c, "创建成功")
}

// * 读所有文章
func getPosts(c *gin.Context) {
	var posts []Post
	res := db.Find(&posts).Error
	if res != nil {
		logger.Error("db 读所有文章失败: %v", res)
		response.Handler.DatabaseError(c, response.ErrDatabaseError)
		return
	}

	response.Handler.Success(c, posts)
}

// * 读单个文章
func getPost(c *gin.Context) {
	id := c.Param("id")
	var post Post
	res := db.Where("id = ?", id).First(&post).Error
	if res != nil {
		logger.Error("db 读单个文章失败: %v %v", res, id)
		response.Handler.BadRequest(c, response.ErrPostNotFound)
		return
	}

	response.Handler.Success(c, post)
}

// * 更新文章
func updatePost(c *gin.Context) {
	type PostUpdate struct {
		ID      int    `json:"id"`
		Title   string `json:"title"`
		Content string `json:"content"`
	}

	var postInfo PostUpdate
	if err := c.ShouldBindJSON(&postInfo); err != nil {
		logger.Error("更新文章参数解析失败: %v", err)
		response.Handler.BadRequest(c, response.ErrMissingParams)
		return
	}

	var post Post
	res := db.Where("id = ?", postInfo.ID).First(&post).Error
	if res != nil {
		logger.Error("db 更新文章失败: %v %v", res, postInfo)
		response.Handler.DatabaseError(c, response.ErrPostUpdateError)
		return
	}

	if post.UserID != c.GetInt("user_id") {
		logger.Error("更新文章失败, 用户无权限: %v tokenuserid: %v postuserid: %v", postInfo, c.GetInt("user_id"), post.UserID)
		response.Handler.Forbidden(c, response.ErrUserNotAuthorized)
		return
	}

	post.Title = postInfo.Title
	post.Content = postInfo.Content

	res = db.Save(&post).Error
	if res != nil {
		logger.Error("db 更新文章失败: %v %v", res, postInfo)
		response.Handler.DatabaseError(c, response.ErrDatabaseError)
		return
	}

	response.Handler.SuccessMsg(c, "更新成功")
}

// * 删除文章
func deletePost(c *gin.Context) {
	type DeletePost struct {
		ID int `json:"id"`
	}

	var postInfo DeletePost
	if err := c.ShouldBindJSON(&postInfo); err != nil {
		logger.Error("删除文章参数解析失败: %v", err)
		response.Handler.BadRequest(c, response.ErrMissingParams)
		return
	}

	res := db.Where("id = ? and user_id = ?", postInfo.ID, c.GetInt("user_id")).Delete(&Post{}).Error
	if res != nil {
		logger.Error("db 删除文章失败: %v %v", res, postInfo)
		response.Handler.DatabaseError(c, response.ErrDatabaseError)
		return
	}

	response.Handler.SuccessMsg(c, "删除成功")
}

// * 创建评论
func createComment(c *gin.Context) {
	type CommentCreate struct {
		PostID  int    `json:"post_id"`
		Content string `json:"content"`
	}

	var commentInfo CommentCreate
	if err := c.ShouldBindJSON(&commentInfo); err != nil {
		logger.Error("创建评论参数解析失败: %v", err)
		response.Handler.BadRequest(c, response.ErrMissingParams)
		return
	}

	var comment Comment
	comment.PostID = commentInfo.PostID
	comment.Content = commentInfo.Content
	comment.UserID = c.GetInt("user_id")

	res := db.Create(&comment).Error
	if res != nil {
		logger.Error("db 创建评论失败: %v %v", res, commentInfo)
		response.Handler.DatabaseError(c, response.ErrCommentCreateError)
		return
	}

	response.Handler.SuccessMsg(c, "创建成功")
}

// * 读所有评论
func getComments(c *gin.Context) {
	postID := c.Param("post_id")
	var comments []Comment
	res := db.Where("post_id = ?", postID).Find(&comments).Error
	if res != nil {
		logger.Error("db 读所有评论失败: %v %v", res, postID)
		response.Handler.DatabaseError(c, response.ErrDatabaseError)
		return
	}

	response.Handler.Success(c, comments)
}

// * 密码加密
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	return string(bytes), err
}

// * 密码验证
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			logger.Error("JWT认证失败, 缺少Authorization头")
			response.Handler.Unauthorized(c, response.ErrUserNotLogin)
			c.Abort()
			return
		}

		tokenString, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected method: %s", token.Header["alg"])
			}

			return jwtkey, nil
		})

		if err != nil || !tokenString.Valid {
			logger.Error("JWT认证失败, token无效: %v", err)
			response.Handler.Unauthorized(c, response.ErrUserNotLogin)
			c.Abort()
			return
		}

		if claims, ok := tokenString.Claims.(jwt.MapClaims); ok && tokenString.Valid {
			if id, ok := claims["id"].(float64); ok {
				c.Set("user_id", int(id))
			}
		}
	}
}

func (p *Post) AfterCreate(tx *gorm.DB) error {
	return tx.Model(&User{}).Where("id = ?", p.UserID).UpdateColumn("post_count", gorm.Expr("post_count + 1")).Error
}

func (p *Post) AfterDelete(tx *gorm.DB) error {
	return tx.Model(&User{}).Where("id = ?", p.UserID).UpdateColumn("post_count", gorm.Expr("post_count - 1")).Error
}

func (c *Comment) AfterDelete(tx *gorm.DB) error {
	var count int64
	err := tx.Model(&Comment{}).Where("post_id = ?", c.PostID).Count(&count).Error
	if err != nil {
		return err
	}
	if count == 0 {
		return tx.Model(&Post{}).Where("id = ?", c.PostID).Update("comment_status", "无评论").Error
	}
	return nil
}
