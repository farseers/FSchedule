package context

import (
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/redis"
)

// RedisContextIns Redis实例
var RedisContextIns *RedisContext

type RedisContext struct {
	redis.IClient `inject:"default"` // 使用farseer.yaml的Redis.default配置节点，并自动注入
}

// InitRedisContext 初始化上下文
func InitRedisContext() {
	RedisContextIns = container.ResolveIns(&RedisContext{})
}
