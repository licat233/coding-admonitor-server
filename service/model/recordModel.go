package model

import (
	"admonitor-server/common/errorx"
	"admonitor-server/common/global"
	"context"
	"errors"
	"time"
)

type RecordModel struct {
	*_BaseMgr
}

func NewRecordModel(tablename string) *RecordModel {
	if global.DB == nil {
		global.LOG.Error(errorx.New("实例化RecordModel失败", errors.New("global.DB is nil")))
		panic(errors.New("global.DB is nil"))
	}
	if len(tablename) == 0 {
		global.LOG.Error(errorx.New("实例化RecordModel失败", errors.New("tablename is nil")))
		panic(errors.New("tablename is nil"))
	}
	ctx, cancel := context.WithCancel(context.Background())
	return &RecordModel{
		_BaseMgr: &_BaseMgr{DB: global.DB.Table(tablename), isRelated: globalIsRelated, ctx: ctx, cancel: cancel, timeout: -1},
	}
}

func (r *RecordModel) FindOne(id int64) (*Record, error) {
	record := new(Record)
	result := r.DB.WithContext(r.ctx).Model(Record{}).First(record, id)
	return record, result.Error
}

func (r *RecordModel) FindByDate(account string, date time.Time) (*Record, error) {
	record := new(Record)
	result := r.DB.WithContext(r.ctx).Model(Record{}).Where("account = ? AND DATE(record_date) = ?", account, date.Local().Format("2006-01-02")).First(record)
	return record, result.Error
}

func (r *RecordModel) FindByDateAndType(account string, date time.Time, Type string) (*Record, error) {
	record := new(Record)
	result := r.DB.WithContext(r.ctx).Model(Record{}).Where("account = ? AND DATE(record_date) = ? AND type = ?", account, date.Local().Format("2006-01-02"), Type).First(record)
	return record, result.Error
}

func (r *RecordModel) FindsByDateAndType(account string, start, end time.Time, Type string) ([]Record, error) {
	var records []Record
	result := r.DB.WithContext(r.ctx).Model(Record{}).Where("account = ? AND  DATE(record_date) between ? AND ? AND type = ?", account, start.Local().Format("2006-01-02"), end.Local().Format("2006-01-02"), Type).Find(&records)
	return records, result.Error
}

func (r *RecordModel) InsertOne(record *Record) error {
	result := r.DB.WithContext(r.ctx).Model(Record{}).Omit("id", "create_at", "update_at").Create(record)
	return result.Error
}

func (r *RecordModel) InsertAll(records []*Record) error {
	result := r.DB.WithContext(r.ctx).Model(Record{}).Omit("id", "create_at", "update_at").Create(records)
	return result.Error
}

// func (r *RecordModel) Save(record *Record) error {
// 	result := r.DB.WithContext(r.ctx).Model(Record{}).Omit("id", "create_at", "update_at").Save(record)
// 	return result.Error
// }

func (r *RecordModel) UpdateByDateAndType(record *Record, date time.Time, Type string) error {
	account := record.Account
	result := r.DB.WithContext(r.ctx).Model(Record{}).Omit("id", "create_at", "update_at").Where("account = ? AND DATE(record_date) = ? AND type = ?", account, date.Local().Format("2006-01-02"), Type).Save(record)
	return result.Error
}

func (r *RecordModel) Delete(record *Record) error {
	result := r.DB.WithContext(r.ctx).Model(Record{}).Delete(record)
	return result.Error
}

func (r *RecordModel) DeleteById(id int64) error {
	record := new(Record)
	result := r.DB.WithContext(r.ctx).Model(Record{}).Delete(record, id)
	return result.Error
}

func (r *RecordModel) DeleteByBetweenDate(account string, start, end time.Time, Type string) error {
	record := new(Record)
	result := r.DB.WithContext(r.ctx).Model(Record{}).Where("account = ? AND DATE(record_date) between ? AND ? AND type = ?", account, start.Local().Format("2006-01-02"), end.Local().Format("2006-01-02"), Type).Delete(record)
	return result.Error
}

func (r *RecordModel) DeleteAll(ids ...int) error {
	var records []Record
	result := r.DB.WithContext(r.ctx).Model(Record{}).Delete(&records, ids)
	return result.Error
}
