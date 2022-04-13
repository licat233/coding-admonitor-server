package global

import (
	"monitor-data/common/timer"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	LOG   *zap.SugaredLogger
	DB    *gorm.DB
	TIMER timer.Timer
)
