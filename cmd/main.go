package main

import (
	"fmt"
	"zhihu-go/config"
	"zhihu-go/internal/dao"
	"zhihu-go/internal/handler"
	"zhihu-go/internal/model"
	"zhihu-go/internal/service"
	"zhihu-go/router"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	//初始化配置
	config.Init()

	//获取mysql配置
	mysqlConfig := config.Config.Mysql
	dsn := mysqlConfig.DSN

	//获取redis配置
	redisConfig := config.Config.Redis
	addr := redisConfig.Addr
	password := redisConfig.Password
	db1 := redisConfig.DB

	//连接到数据库
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("failed to connect database:", err)
		return
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db1,
	})

	//自动迁移
	if err := db.AutoMigrate(&model.User{}, &model.Post{}, &model.Follow{}, &model.Comment{}, &model.Like{}); err != nil {
		fmt.Println("Failed to run migrate:", err)
		return
	}

	//创建 DAO 实例
	userDAO := dao.NewUserDAO(db)
	postDAO := dao.NewPostDAO(db)
	likeDAO := dao.NewLikeDAO(db)
	followDAO := dao.NewFollowDAO(db)
	commentDAO := dao.NewCommentDAO(db)

	//创建 Service 实例，注入 DAO
	userService := service.NewUserService(userDAO)
	postService := service.NewPostService(postDAO)
	likeService := service.NewLikeService(likeDAO, rdb)
	followService := service.NewFollowService(followDAO, userDAO)
	commentService := service.NewCommentService(commentDAO, rdb)

	//创建 Handler 实例，注入 Service
	userHandler := handler.NewUserHandler(userService)
	postHandler := handler.NewPostHandler(postService)
	likeHandler := handler.NewLikeHandler(likeService)
	followHandler := handler.NewFollowHandler(followService)
	commentHandler := handler.NewCommentHandler(commentService)

	//设置路由
	r := gin.Default()

	//使用 Router 结构体
	routerInstance := router.NewRouter(userHandler, postHandler, likeHandler, followHandler, commentHandler) // 传入需要的 handler
	routerInstance.SetUp(r)

	/*// 使用配置中的管理员信息
	adminUser := config.Config.Admin.Username
	adminPass := config.Config.Admin.Password

	// 调用初始化管理员
	if err := service.InitAdmin(db, adminUser, adminPass); err != nil {
		log.Fatalf("初始化管理员失败: %v", err)
	}*/

	//启动gin服务
	r.Run(":8080")
}
