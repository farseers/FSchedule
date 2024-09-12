package client

import (
	"sync"
)

var clientList = sync.Map{}

// 统计客户端数量
func GetClientCount() int {
	count := 0
	clientList.Range(func(key, value any) bool {
		count++
		return true
	})
	return count
}
