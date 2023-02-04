package schedule

import (
	"github.com/farseer-go/fs/core"
)

type Repository interface {
	GetLock(name string) core.ILock
}
