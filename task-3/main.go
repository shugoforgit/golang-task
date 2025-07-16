package main

import (
	"errors"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var dsn = "root:shugoformysql@tcp(127.0.0.1:3306)/testdb?charset=utf8mb4&parseTime=True&loc=Local"

func main() {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}

	// * SQL语句练习1：基本CRUD操作
	task1_insert(db)
	task1_select(db)
	task1_update(db)
	task1_delete(db)

	// * SQL语句练习2：事务语句
	task2_transaction(db)

	// * Sqlx入门1：使用SQL扩展库进行查询
	task3_sqlx_select()

	// * Sqlx入门2：实现类型安全映射
	task4_sqlx_type_safe_mapping()

	// * 进阶gorm1：通过结构体映射创建数据库表
	task5_gorm_auto_migrate(db)

	// * 进阶gorm2：关联查询
	task6_gorm_associations(db)

	// * 进阶gorm3：钩子函数
	task7_gorm_hooks(db)
}

// * SQL语句练习1：基本CRUD操作--------------------------------------
type Student struct {
	Name  string
	Age   int
	Grade string
}

func task1_insert(db *gorm.DB) {
	db.Create(&Student{
		Name:  "张三",
		Age:   20,
		Grade: "三年级",
	})
}

func task1_select(db *gorm.DB) {
	var students []Student
	db.Where("age > ?", 18).Find(&students)
	fmt.Println(students)
}

func task1_update(db *gorm.DB) {
	db.Model(&Student{}).Where("name = ?", "张三").Update("grade", "四年级")
}

func task1_delete(db *gorm.DB) {
	db.Where("age < ?", 15).Delete(&Student{})
}

// * SQL语句练习2：事务语句--------------------------------------
type Accounts struct {
	ID      int
	Balance int
}
type Transactions struct {
	ID              int
	From_account_id int
	To_account_id   int
	Amount          int
}

func task2_transaction(db *gorm.DB) {
	var accounts1 Accounts
	var accounts2 Accounts
	err := db.Transaction(func(tx *gorm.DB) error {
		res := tx.Where("id = ?", 1).Find(&accounts1)
		if res.Error != nil {
			return errors.New("find account1 failed")
		}

		if accounts1.Balance < 100 {
			tx.Rollback()
			return errors.New("balance is not enough")
		}

		res = tx.Where("id = ?", 2).Find(&accounts2)
		if res.Error != nil {
			return errors.New("find account2 failed")
		}

		res1 := tx.Model(&Accounts{}).Where("id = ?", 1).Update("balance", accounts1.Balance-100)
		res2 := tx.Model(&Accounts{}).Where("id = ?", 2).Update("balance", accounts2.Balance+100)
		if res1.Error != nil || res2.Error != nil {
			tx.Rollback()
			return errors.New("update balance failed")
		}

		res3 := tx.Create(&Transactions{
			From_account_id: 1,
			To_account_id:   2,
			Amount:          100,
		})
		if res3.Error != nil {
			tx.Rollback()
			return errors.New("create transaction failed")
		}

		return nil
	})

	if err != nil {
		fmt.Println(err)
	}
}

// * Sqlx入门1：使用SQL扩展库进行查询----------------------
type Employees struct {
	ID         int
	Name       string
	Department string
	Salary     int
}

func task3_sqlx_select() {
	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}

	var employees []Employees
	err = db.Select(&employees, "select * from employees where department = ?", "技术部")
	if err != nil {
		panic("failed to query: " + err.Error())
	}
	fmt.Println(employees)

	var maxSalary Employees
	err = db.Get(&maxSalary, "select * from employees where salary = (select max(salary) from employees)")
	if err != nil {
		panic("failed to query: " + err.Error())
	}
	fmt.Println("maxSalary: ", maxSalary)
}

// * Sqlx入门2：实现类型安全映射
type Books struct {
	ID     int     `db:"id"`
	Title  string  `db:"title"`
	Author string  `db:"author"`
	Price  float64 `db:"price"`
}

func task4_sqlx_type_safe_mapping() {
	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}

	var books []Books
	err = db.Select(&books, "select * from books where price > ?", 50)
	if err != nil {
		panic("failed to query: " + err.Error())
	}
	fmt.Println(books)
}

// * 进阶gorm1：模型定义
type User struct {
	ID        int    `gorm:"primaryKey"`
	Name      string `gorm:"type:varchar(100);not null"`
	Email     string `gorm:"type:varchar(100);not null"`
	PostCount int    `gorm:"default:0"`
	Posts     []Post `gorm:"foreignKey:UserID"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Post struct {
	ID            int       `gorm:"primaryKey"`
	Title         string    `gorm:"type:varchar(100);not null"`
	Content       string    `gorm:"type:text;not null"`
	UserID        int       `gorm:"not null"`
	Comments      []Comment `gorm:"foreignKey:PostID"`
	CommentStatus string    `gorm:"type:varchar(20);default:'无评论'"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type Comment struct {
	ID        int    `gorm:"primaryKey"`
	Content   string `gorm:"type:text;not null"`
	PostID    int    `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func task5_gorm_auto_migrate(db *gorm.DB) {
	err := db.AutoMigrate(&User{}, &Post{}, &Comment{})
	if err != nil {
		panic("failed to migrate database: " + err.Error())
	}
}

// * 进阶gorm2：关联查询
func task6_gorm_associations(db *gorm.DB) {
	var user User
	db.Preload("Posts").Preload("Posts.Comments").Where("id = ?", 1).First(&user)
	fmt.Println("用户发布的所有文章和评论: ", user)

	var postWithMostComments Post
	err := db.Raw(`
		SELECT p.*, COUNT(c.id) as comment_count 
		FROM posts p 
		LEFT JOIN comments c ON p.id = c.post_id 
		GROUP BY p.id 
		ORDER BY comment_count DESC 
		LIMIT 1
	`).Scan(&postWithMostComments).Error
	if err != nil {
		panic("failed to query: " + err.Error())
	}
	fmt.Println("评论最多的文章: ", postWithMostComments)
}

// * 进阶gorm3：钩子函数
func (p *Post) AfterCreate(tx *gorm.DB) error {
	result := tx.Model(&User{}).Where("id = ?", p.UserID).UpdateColumn("post_count", gorm.Expr("post_count + 1"))
	return result.Error
}

func (p *Post) AfterDelete(tx *gorm.DB) error {
	result := tx.Model(&User{}).Where("id = ?", p.UserID).UpdateColumn("post_count", gorm.Expr("post_count - 1"))
	return result.Error
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

func task7_gorm_hooks(db *gorm.DB) {
	// * 创建文章
	var post Post
	post.Title = "test"
	post.Content = "test"
	post.UserID = 1
	res := db.Create(&post)
	if res.Error != nil {
		panic("failed to create post: " + res.Error.Error())
	}
	fmt.Println("post created: ", post)

	// * 删除文章
	res = db.Delete(&post)
	if res.Error != nil {
		panic("failed to delete post: " + res.Error.Error())
	}
	fmt.Println("post deleted successfully")

	// * 删除评论
	var comment Comment
	comment.PostID = 1
	res = db.Where("post_id = ?", comment.PostID).Delete(&comment)
	if res.Error != nil {
		panic("failed to delete comment: " + res.Error.Error())
	}
	fmt.Println("comment deleted successfully")
}
