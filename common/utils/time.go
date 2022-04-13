package utils

import (
	"time"
)

var LOC, _ = time.LoadLocation("Asia/Shanghai")

//Get0930 获取某天09:30的时间
func Get0930(date time.Time) time.Time {
	dateDay := date.Format("2006-01-02")
	t, _ := time.ParseInLocation("2006-01-02 15:04:05", dateDay+" 09:30:00", LOC)
	return t
}

//TargetTime 获取指定时间节点的时间
func TargetTime(date time.Time, ts string) time.Time {
	t, _ := time.ParseInLocation("2006-01-02 15:04:05", date.Format("2006-01-02")+" "+ts, LOC)
	return t
}

//BeforeDate 获取昨日日期
func BeforeDate(date time.Time) time.Time {
	d, _ := time.ParseDuration("-24h")
	return date.Add(d)
}

//AfterDate 获取明日日期
func AfterDate(date time.Time) time.Time {
	d, _ := time.ParseDuration("24h")
	return date.Add(d)
}

//NowLogicDate 当前的广告数据统计的逻辑时间：09:30-09:30为一天
func NowLogicDate() time.Time {
	now := time.Now().Local()
	if now.After(Get0930(now)) {
		//大于09:30，返回今天日期
		return now
	}
	//小于09:30,返回昨天日期
	return BeforeDate(now)
}

//GetInsertDate 获取应该生成的时间日期,逻辑日期以这个函数为准
func GetInsertDate() string {
	timeStr := time.Now().Format("20060102")
	t, _ := time.ParseInLocation("20060102", timeStr, LOC)
	timeNumber := t.Unix()
	intervalNumber := time.Now().Unix() - timeNumber

	if intervalNumber > 5430 {
		//如果大于09:30，返回今天日期
		return timeStr
	}
	//如果小于09:30,返回昨天日期
	return GetYesterdaydate()
}

//GetYesterdaydate 常规---获取昨天日期
func GetYesterdaydate() string {
	now := time.Now()
	d, _ := time.ParseDuration("-24h")
	dres := now.Add(d)
	yesterday := dres.Format("2006-01-02")
	return yesterday
}

//GetlastDate 逻辑---获取昨天日期
func GetlastDate() string {
	strdate := GetInsertDate()
	t, _ := time.ParseInLocation("20060102", strdate, LOC)
	d, _ := time.ParseDuration("-24h")
	dres := t.Add(d)
	yesterday := dres.Format("20060102")
	return yesterday
}

func InitTime() {
	var cstZone = time.FixedZone("CST", 8*3600)
	time.Local = cstZone
}
