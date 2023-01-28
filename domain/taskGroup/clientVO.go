package taskGroup

import "fmt"

// ClientVO 客户端
type ClientVO struct {
	Id   int64  // 客户端Id
	Name string // 客户端名称
	Ip   string // 客户端IP
	Port int    // 客户端IP
}

// Endpoint 获取终结点
func (receiver ClientVO) Endpoint() string {
	return fmt.Sprintf("%s:%d", receiver.Ip, receiver.Port)
}
