# 蓝山GO组寒假考核项目一：知乎社区

## 实现功能

### 基础功能

1. 基本的用户系统：用户的注册与登录，用户的基本信息，基本的用户鉴权，密码的加密，用户个人信息的更改
2. 发布文章与评论：发布/删除/更新/获取文章，对文章进行评论，对文章进行点赞
3. 关注功能：关注与取消关注，关注列表与关注者列表
4. 内容的搜索：利用MySQL的全文索引功能实现搜索功能
5. ⽂章状态(草稿，已删除，未删除)
6. 管理员禁言
7. 防刷机制(点赞频率，评论频率)

### 进阶功能

1. 加入缓存策略：对基本的⽂章内容等数据进⾏缓存，并采取合理的缓存策略。使⽤go内置的singleflight、布隆过滤器等来应对缓存击穿、缓存雪崩这些问题
2. 接口文档：使用ApiFox保存接口文档: https://s.apifox.cn/c4dc979c-f66d-4aa6-8d65-5e936429c08e
3. 配置管理：使用viper加载如MySQL相关的配置信息
4. 日志管理：通过zap日志库集成日志
5. docker：通过编写dockerFile将项目打包成镜像，并推送到dockerhub上
6. feed流（接受关注的人的动态）：采用推模式（写扩散）实现关注动态推送。用户发布文章后，异步将文章推送给所有粉丝的收件箱（timeline 表），粉丝查看关注动态时直接读取自己的收件箱，分页返回。
---
## 环境依赖
1. Gin: 轻量级 Web 框架，用于处理 HTTP 请求

2. GORM: ORM 工具，连接并操作 MySQL 数据库

3. Go-Redis: Redis 客户端，用作缓存存储

4. Jwt-Go: JWT 组件，用于身份验证和授权

5. Viper: 配置管理库，支持多格式配置文件读取

6. Crypto: 加密库（golang.org/x/crypto），用于密码哈希等安全操作

7. Zap: 高性能结构化日志库，用于快速记录和输出日志信息 

8. RedisBloom: 布隆过滤器模块

---
## 目录结构

| 目录 | 说明 |
|------|------|
| `cmd/` | 项目入口，main.go 所在，负责依赖组装 |
| `config/` | 配置加载逻辑（viper 等） |
| `internal/cache/` | 缓存层接口及 Redis 实现（文章、用户等缓存） |
| `internal/dao/` | 数据访问层，定义接口并通过 GORM 操作数据库 |
| `internal/dto/` | 数据传输对象，定义与前端交互的请求/响应结构体 |
| `internal/handler/` | HTTP 处理层，接收请求、调用 service 并返回响应 |
| `internal/middleware/` | Gin 中间件（JWT 鉴权、管理员检查、限流等） |
| `internal/model/` | 数据库模型（GORM 模型定义） |
| `internal/service/` | 业务逻辑层，组合 DAO、缓存等，实现核心功能 |
| `pkg/` | 公共工具包，可被外部项目复用（加密、JWT、限流、统一响应等） |
| `router/` | 路由注册，集中管理所有路由及中间件 |
| `Dockerfile` | Docker 构建文件 |
| `docker-compose.yaml` | 本地开发环境编排（MySQL、Redis） |
| `go.mod` | 环境依赖 |
| `README` |

## 快速开始（示例）
1. 克隆项目:  

`git clone https://github.com/yourname/zhihu-go.git`  

2. 配置:  

`复制 config/config.example.yaml 为 config/config.yaml，修改数据库、Redis 等配置。`  

3. 启动依赖:  

`docker-compose up -d mysql redis（需安装 Docker）`  

4. 运行项目:  

`go run cmd/main.go`  

5. 访问接口:

`接口文档见 Apifox 链接，默认地址 http://localhost:8080`

