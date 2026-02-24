package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"zhihu-go/internal/model"

	"github.com/redis/go-redis/v9"
)

// PostCache 定义文章缓存接口
type PostCache interface {
	Key(postID uint) string
	Get(ctx context.Context, postID uint) (*model.Post, error)
	Set(ctx context.Context, post *model.Post) error
}

// postCache Redis 实现
type postCache struct {
	client  *redis.Client
	prefix  string        //缓存key前缀
	baseTTL time.Duration //默认过期时间
}

// NewPostCache 创建文章缓存实例
func NewPostCache(client *redis.Client) PostCache {
	return &postCache{
		client:  client,
		prefix:  "post:",
		baseTTL: time.Hour,
	}
}

// Key 制造key
func (c *postCache) Key(postID uint) string {
	return fmt.Sprintf("%s%d", c.prefix, postID)
}

// Get 从缓存获取文章， 返回nil表示不存在
func (c *postCache) Get(ctx context.Context, postID uint) (*model.Post, error) {
	data, err := c.client.Get(ctx, c.Key(postID)).Bytes()
	if errors.Is(err, redis.Nil) {
		return nil, nil //缓存不存在
	}
	if err != nil {
		return nil, err
	}
	var post model.Post
	if err := json.Unmarshal(data, &post); err != nil {
		return nil, err
	}

	return &post, nil
}

// Set 将文章存入缓存
func (c *postCache) Set(ctx context.Context, post *model.Post) error {
	data, err := json.Marshal(post)
	if err != nil {
		return err
	}

	// 随机偏移量：基础 TTL 的 ±20%
	offset := time.Duration(rand.Intn(int(c.baseTTL/5))) - c.baseTTL/10
	ttl := c.baseTTL + offset
	
	return c.client.Set(ctx, c.Key(post.ID), data, ttl).Err()
}
