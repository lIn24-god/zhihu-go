package utils

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// Allow 检查用户某个动作是否允许执行（限流）
// rdb: Redis 客户端
// action: 动作标识
// userID: 用户ID
// limit: 时间窗口内允许的最大次数
// window: 时间窗口，单位秒
func Allow(rdb *redis.Client, action string, userID uint, limit int, window int64) (bool, error) {
	ctx := context.Background()
	key := fmt.Sprintf("rate %s, user %d", action, userID)

	count, err := rdb.Incr(ctx, key).Result()
	if err != nil {
		return false, err
	}

	if count == 1 {
		rdb.Expire(ctx, key, time.Duration(window)*time.Second)
	}

	if count > int64(limit) {
		return false, nil
	}

	return true, nil
}
