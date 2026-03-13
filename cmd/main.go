package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"
	"zhihu-go/config"
	"zhihu-go/internal/cache"
	"zhihu-go/internal/dao"
	"zhihu-go/internal/handler"
	"zhihu-go/internal/model"
	"zhihu-go/internal/service"
	"zhihu-go/pkg/bloom"
	"zhihu-go/pkg/logger"
	"zhihu-go/router"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	//初始化配置
	config.Init()

	// 2. 初始化日志（使用 config 中的日志配置）
	if err := logger.Init(config.Config.Log.Level, config.Config.Log.Output); err != nil {
		// 如果日志初始化失败，可以暂时用标准库打印并退出，或者 panic
		panic("初始化日志失败: " + err.Error())
	}
	// 确保程序退出前日志全部写入
	defer logger.L().Sync()

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
		logger.S().Fatalw("数据库连接失败", "dsn", dsn, "error", err)
		// 可以选择 panic 或 os.Exit，但为了容器不立即重启，可以让程序 sleep 一段时间再退出
		time.Sleep(10 * time.Second)
		os.Exit(1)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db1,
	})

	logger.S().Infow("数据库和 Redis 初始化完成",
		"mysql_dsn", config.Config.Mysql.DSN,
		"redis_addr", config.Config.Redis.Addr,
	)

	//自动迁移
	if err := db.AutoMigrate(&model.User{},
		&model.Post{},
		&model.Comment{},
		&model.Like{},
		&model.UserRelation{},
		&model.Timeline{}); err != nil {
		fmt.Println("Failed to run migrate:", err)
		return
	}

	//创建 DAO 实例
	userDAO := dao.NewUserDAO(db)
	postDAO := dao.NewPostDAO(db)
	likeDAO := dao.NewLikeDAO(db)
	followDAO := dao.NewFollowDAO(db)
	commentDAO := dao.NewCommentDAO(db)
	timelineDAO := dao.NewTimelineDAO(db)

	// 初始化缓存
	postCache := cache.NewPostCache(rdb)

	// 初始化布隆过滤器
	bloomFilter := bloom.NewRedisBloom(rdb)

	//创建 Service 实例，注入 DAO
	feedService := service.NewFeedService(
		followDAO,
		timelineDAO,
		postDAO,
		userDAO,
		5, // worker 数量
	)
	userService := service.NewUserService(userDAO, rdb)
	postService := service.NewPostService(postDAO, userService, postCache, bloomFilter, feedService)
	likeService := service.NewLikeService(likeDAO, rdb)
	followService := service.NewFollowService(followDAO, userDAO)
	commentService := service.NewCommentService(commentDAO, userService, rdb)

	//创建 Handler 实例，注入 Service
	userHandler := handler.NewUserHandler(userService)
	postHandler := handler.NewPostHandler(postService, userService)
	likeHandler := handler.NewLikeHandler(likeService)
	followHandler := handler.NewFollowHandler(followService)
	commentHandler := handler.NewCommentHandler(commentService, userService)
	feedHandler := handler.NewFeedHandler(feedService)

	//设置路由
	r := gin.Default()

	//使用 Router 结构体
	routerInstance := router.NewRouter(userHandler, postHandler, likeHandler,
		followHandler, commentHandler, userService, feedHandler) // 传入需要的 handler
	routerInstance.SetUp(r)

	// 使用配置中的管理员信息
	adminUser := config.Config.Admin.Username
	adminPass := config.Config.Admin.Password

	// 创建带超时的 context（5秒超时）
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel() // 确保资源释放

	// 调用初始化管理员
	if err := userService.InitAdmin(ctx, adminUser, adminPass); err != nil {
		log.Fatalf("初始化管理员失败: %v", err)
	}

	//启动gin服务
	err1 := r.Run(":8080")
	if err1 != nil {
		return
	}
}
