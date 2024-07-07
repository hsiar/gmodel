package corm

import (
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"testing"
)

type User struct {
	//gorm.Model
	ID   uint   `gorm:"primaryKey;AUTO_INCREMENT" json:"id"`
	Name string `json:"name"`
	//CategoryID uint      `json:"categoryId"`
	//Category   *Category `json:"category"`

	Products []Product
}

type Product struct {
	Id     int64  `gorm:"primaryKey;AUTO_INCREMENT" json:"id"`
	UserId int64  `json:"user_id"`
	Name   string `json:"name"`
}

func InitClient() (client *Client) {
	dns := "root:123456@tcp(127.0.0.1:3306)/gorse_demo?charset=utf8&parseTime=True&loc=Local"
	client = NewClient().WithConfig(&Config{
		Dsn:         dns,
		TablePrefix: "z_",
	}).Open()
	_ = client.Db.AutoMigrate(&User{})
	_ = client.Db.AutoMigrate(&Product{})
	return
}

func TestNewGormDao(t *testing.T) {
	client := InitClient()
	err := client.Db.AutoMigrate(&User{})
	hlog.Debug(err)
}

func TestInsert(t *testing.T) {
	//InitTestEnv(1)
	client := InitClient()
	user := &User{Name: "dd"}
	client.Db.Create(user)
}

func TestHasMany(t *testing.T) {
	client := InitClient()
	user := &User{Name: "cc"}

	err := client.Db.Preload("Products").First(user).Error
	hlog.Debug(err)

}
