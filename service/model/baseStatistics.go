package model

import (
	"admonitor-server/common/global"
	"errors"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// StatisticsString 广告投放监控数据，数字类型全部转为字符串类型
// 用于方便前端展示
type StatisticsString struct {
	ID         int64     `json:"id"`         // 主键
	IDentifier string    `json:"identifier"` // 唯一标识符，由channel_id、account_id、record_date、type用-符号连接而成，目的是为了防止重复插入
	ChannelID  int64     `json:"channelId"`  // 渠道id，指的是channel表的主键
	AccountID  int64     `json:"accountId"`  // 账号id，指的是account表中的主键
	Cost       string    `json:"cost"`       // 消耗额
	CostCny    string    `json:"costCny"`    //人民币消耗
	Click      string    `json:"click"`      // 点击数
	Imp        string    `json:"imp"`        // 展现量
	Ctr        string    `json:"ctr"`        // 点击率
	Cv         string    `json:"cv"`         // 转化数
	Cvr        string    `json:"cvr"`        // 转化率
	Cpa        string    `json:"cpa"`        // 转化成本
	Status     bool      `json:"status"`     // 状态
	Type       string    `json:"type"`       // 数据类型：newest最新数据，node节点数据，day全天数据，logic逻辑数据
	Currency   string    `json:"currency"`   // 币种
	RecordDate time.Time `json:"recordDate"` // 记录日期，09:30-09:30为一天
	CreateAt   time.Time `json:"createAt"`   // 创建时间
	UpdateAt   time.Time `json:"updateAt"`   // 更新时间
}

type StatisticsModel struct {
	*_StatisticsMgr
}

func NewStatisticsModel() *StatisticsModel {
	return &StatisticsModel{
		_StatisticsMgr: StatisticsMgr(global.DB),
	}
}

func (m *StatisticsModel) baseModel() *gorm.DB {
	return m.DB.WithContext(m.ctx).Model(Statistics{})
}

func (m *StatisticsModel) whereModel(channel_id, account_id int64) *gorm.DB {
	return m.baseModel().Where("channel_id = ? AND account_id = ?", channel_id, account_id)
}

func (m *StatisticsModel) FindOne(id int64) (*Statistics, error) {
	statistics := new(Statistics)
	result := m.baseModel().First(statistics, id)
	return statistics, result.Error
}

func (m *StatisticsModel) FindByDateAndType(channel_id, account_id int64, date time.Time, Type string) (*Statistics, error) {
	statistics := new(Statistics)
	result := m.whereModel(channel_id, account_id).Where("DATE(record_date) = ? AND type = ?", date.Local().Format("2006-01-02"), Type).Order("update_at desc").First(statistics)
	return statistics, result.Error
}

func (m *StatisticsModel) FindsByDateAndType(channel_id, account_id int64, start, end time.Time, Type string) ([]*Statistics, error) {
	var statistics []*Statistics
	result := m.whereModel(channel_id, account_id).Where("DATE(record_date) between ? AND ? AND type = ?", start.Local().Format("2006-01-02"), end.Local().Format("2006-01-02"), Type).Find(&statistics)
	return statistics, result.Error
}

func (m *StatisticsModel) FindsByAccountIds(channel_id int64, account_ids []int64, start, end time.Time, Type string) ([]*Statistics, error) {
	var statistics []*Statistics
	result := m.baseModel().Where("channel_id = ? AND account_id in (?) AND DATE(record_date) between ? AND ? AND type = ?", channel_id, account_ids, start.Local().Format("2006-01-02"), end.Local().Format("2006-01-02"), Type).Find(&statistics)
	return statistics, result.Error
}

func (m *StatisticsModel) FindByIdentifier(identifier string) (*Statistics, error) {
	statistic := new(Statistics)
	err := m.baseModel().Where("identifier = ? ", identifier).First(statistic).Error
	return statistic, err
}

func (m *StatisticsModel) ExistIdentifier(identifier string) (bool, error) {
	var count int64
	err := m.baseModel().Where("identifier = ? ", identifier).Limit(1).Count(&count).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	} else if err != nil {
		return false, err
	}
	if count == 0 {
		return false, nil
	}
	return true, nil
}

func (m *StatisticsModel) FindsByIdentifier(identifier string) ([]*Statistics, error) {
	statistics := []*Statistics{}
	err := m.baseModel().Where("identifier = ? ", identifier).Find(statistics).Error
	return statistics, err
}

func (m *StatisticsModel) InsertOne(statistics *Statistics) error {
	err := m.baseModel().Omit("id", "create_at", "update_at").Create(statistics).Error
	// if errors.Is(err, gorm.ErrPreloadNotAllowed) {
	// 	global.LOG.Warn("不允许重复插入")
	// 	return nil
	// }
	return err
}

func (m *StatisticsModel) InsertAll(statisticss []*Statistics) error {
	return m.baseModel().Omit("id", "create_at", "update_at").Create(statisticss).Error
}

/**更新 */

func (m *StatisticsModel) UpdateOne(statistics *Statistics) error {
	result := m.baseModel().Where("id =? ", statistics.ID).Omit("id", "create_at", "update_at").Updates(statistics)
	return result.Error
}

func (m *StatisticsModel) UpdateByIdentifier(statistics *Statistics) error {
	result := m.baseModel().Where("identifier =? ", statistics.IDentifier).Omit("identifier", "create_at", "update_at").Updates(statistics)
	return result.Error
}

func (m *StatisticsModel) Delete(statistics *Statistics) error {
	result := m.baseModel().Delete(statistics)
	return result.Error
}

func (m *StatisticsModel) DeleteById(id int64) error {
	statistics := new(Statistics)
	result := m.baseModel().Delete(statistics, id)
	return result.Error
}

func (m *StatisticsModel) DeleteByBetweenDate(channel_id, account_id int64, start, end time.Time, Type string) error {
	statistics := new(Statistics)
	result := m.whereModel(channel_id, account_id).Where("DATE(record_date) between ? AND ? AND type = ?", start.Local().Format("2006-01-02"), end.Local().Format("2006-01-02"), Type).Delete(statistics)
	return result.Error
}

func (m *StatisticsModel) DeleteAll(channel_id, account_id int64, ids ...int) error {
	var statisticss []Statistics
	result := m.whereModel(channel_id, account_id).Delete(&statisticss, ids)
	return result.Error
}

//Replace 替换单个数据，不存在则创建
func (m *StatisticsModel) Replace(statistic *Statistics) error {
	if statistic == nil {
		return nil
	}
	if statistic.IDentifier == "" {
		SetStatisticsIdentifier(statistic)
	}
	// //查询是否存在
	// exist, err := m.ExistIdentifier(statistic.IDentifier)
	// if err != nil {
	// 	global.LOG.Error("mysql查询statistics数据失败", zap.Error(err))
	// 	return err
	// }
	// //如果存在
	// if exist {
	// }

	find, err := m.FindByIdentifier(statistic.IDentifier)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = m.InsertOne(statistic)
		if errors.Is(err, gorm.ErrPreloadNotAllowed) {
			global.LOG.Warn("不允许插入相同identifer的statistic", zap.Error(err))
			goto UPDATE
		}
		if err != nil {
			global.LOG.Error("mysql插入statistics数据失败", zap.Error(err))
		}
		return err
	} else if err != nil {
		global.LOG.Error("mysql查询statistics数据失败", zap.Error(err))
		return err
	}
UPDATE:
	statistic.ID = find.ID
	if err = m.UpdateOne(statistic); err != nil {
		global.LOG.Error("mysql更新statistics数据失败", zap.Error(err))
	}
	return err
}

//ReplaceAll 替换所有，存在则更新，不存在则插入
func (m *StatisticsModel) ReplaceAll(statistics []*Statistics) (err error) {
	if len(statistics) == 0 {
		return
	}

	for _, statistic := range statistics {
		if err = m.Replace(statistic); err != nil {
			return
		}
	}
	return
}
