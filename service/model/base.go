package model

import (
	"errors"
	"fmt"
	"time"
)

const (
	NEWEST = "newest" //最新数据
	NODE   = "node"   //节点数据
	DAY    = "day"    //全天数据
	LOGIC  = "logic"  //逻辑数据

)

var ErrNotFound = errors.New("not found record")

//SetStatisticsIdentifier 设置statistics唯一标识符
func SetStatisticsIdentifier(statistics *Statistics) {
	statistics.IDentifier = fmt.Sprintf("%d_%d_%s_%s", statistics.ChannelID, statistics.AccountID, statistics.RecordDate.Format("2006-01-02"), statistics.Type)
}

//GenStatisticsIdentifier 生成statistics唯一标识符
func GenStatisticsIdentifier(channel_id, account_id int64, date time.Time, Type string) string {
	return fmt.Sprintf("%d_%d_%s_%s", channel_id, account_id, date.Format("2006-01-02"), Type)
}
