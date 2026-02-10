package main

import (
	"fmt"
	"zhihu-go/config"
	"zhihu-go/internal/model"
	"zhihu-go/router"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	//初始化配置
	config.Init()

	//获取mysql配置
	mysqlConfig := config.Config.Mysql
	dsn := mysqlConfig.DSN

	//连接到数据库
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("failed to connect database:", err)
		return
	}

	//自动迁移
	if err := db.AutoMigrate(&model.User{}, &model.Post{}, &model.Follow{}, &model.Comment{}, &model.Like{}); err != nil {
		fmt.Println("Failed to run migrate:", err)
		return
	}

	//初始化路由并传递数据库连接
	r := router.SetUpRouter(gin.Default(), db)

	//启动gin服务
	r.Run(":8080")

}
