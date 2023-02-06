package schedule

import (
	"github.com/farseer-go/fs/core"
)

type Repository interface {
	NewLock(name string) core.ILock
}
