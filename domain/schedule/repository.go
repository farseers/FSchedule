package schedule

import (
	"github.com/farseer-go/fs/core"
	"time"
)

type Repository interface {
	GetLock(name string, nextAt time.Time) core.ILock
}
