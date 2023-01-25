package client

type Repository interface {
	// Save 保存客户端信息
	Save(do DomainObject)
}
