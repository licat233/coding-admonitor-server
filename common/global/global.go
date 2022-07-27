package global

import (
	"admonitor-server/common/timer"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	LOG   *zap.SugaredLogger
	DB    *gorm.DB
	TIMER timer.Timer
	REDIS *redis.Client
)
