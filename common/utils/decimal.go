package utils

import (
	"admonitor-server/service/model"
	"strings"

	"github.com/shopspring/decimal"
)

func CtrCalculation(click, imp int64) float64 {
	if imp == 0 {
		return 0
	}
	clicks := decimal.NewFromInt(click * 100)
	imps := decimal.NewFromInt(imp)
	ctr, _ := clicks.Div(imps).Round(2).Float64()
	return ctr
}

func CvrCalculation(cv, click int64) float64 {
	if click == 0 {
		return 0
	}
	cvs := decimal.NewFromInt(cv * 100)
	clicks := decimal.NewFromInt(click)
	cvr, _ := cvs.Div(clicks).Round(2).Float64()
	return cvr
}

func CpaCalculation(cost float64, cv int64) float64 {
	if cv == 0 {
		return 0
	}
	costs := decimal.NewFromFloat(cost)
	cvs := decimal.NewFromInt(cv)
	cpa, _ := costs.Div(cvs).Round(2).Float64()
	return cpa
}

func CostCalculation(newCost, oldCost float64) float64 {
	newCosts := decimal.NewFromFloat(newCost)
	oldCosts := decimal.NewFromFloat(oldCost)
	cost, _ := newCosts.Sub(oldCosts).Round(2).Float64()
	return cost
}

func StringToFloat64(str string) float64 {
	if str == "" {
		return 0
	}
	d, _ := decimal.NewFromString(str)
	f, _ := d.Float64()
	return f
}

func StringToInt64(str string) int64 {
	if str == "" {
		return 0
	}
	d, _ := decimal.NewFromString(str)
	return d.IntPart()
}

//CalculateStatistics 计算两个statistic合并的数据
func CalculateStatistics(target, source *model.Statistics) {
	cost, _ := decimal.NewFromFloat(target.Cost).Add(decimal.NewFromFloat(source.Cost)).Round(2).Float64()
	imp := target.Imp + source.Imp
	click := target.Click + source.Click
	cv := target.Cv + source.Cv
	ctr := CtrCalculation(click, imp)
	cvr := CvrCalculation(cv, click)
	cpa := CpaCalculation(cost, cv)
	target.IDentifier = source.IDentifier
	target.AccountID = source.AccountID
	target.ChannelID = source.ChannelID
	target.Cost = cost
	target.Imp = imp
	target.Click = click
	target.Cv = cv
	target.Ctr = ctr
	target.Cvr = cvr
	target.Cpa = cpa
	target.Type = source.Type
	target.Currency = source.Currency
}

//CalculateTotalStatistics 计算总的统计数据
func CalculateTotalStatistics(lists []*model.Statistics) *model.Statistics {
	if len(lists) == 0 {
		return nil
	}
	target := &model.Statistics{}
	cost := decimal.NewFromFloat(0)
	for _, source := range lists {
		cost = cost.Add(decimal.NewFromFloat(source.Cost))
		target.Imp += source.Imp
		target.Click += source.Click
		target.Cv += source.Cv
	}
	target.Cost, _ = cost.Round(2).Float64()
	target.Ctr = CtrCalculation(target.Click, target.Imp)
	target.Cvr = CvrCalculation(target.Cv, target.Click)
	target.Cpa = CpaCalculation(target.Cost, target.Cv)
	firstRecord := lists[0]
	target.IDentifier = firstRecord.IDentifier
	target.AccountID = firstRecord.AccountID
	target.ChannelID = firstRecord.ChannelID
	target.Type = firstRecord.Type
	target.Currency = firstRecord.Currency
	target.RecordDate = firstRecord.RecordDate
	target.CreateAt = firstRecord.CreateAt
	target.UpdateAt = firstRecord.UpdateAt
	return target
}

//StatisticConvert statistic数据转换
func StatisticConvert(statistic *model.Statistics) *model.StatisticsString {
	if statistic == nil {
		return nil
	}
	cost := decimal.NewFromFloat(statistic.Cost).Round(2).String()
	var cny string
	currency := strings.ToUpper(statistic.Currency)
	if currency == "TWD" {
		cny = decimal.NewFromFloat(statistic.Cost).Div(decimal.NewFromFloat(4.5)).Round(2).String()
	} else if currency == "CNY" {
		cny = cost
	} else {
		cny = "暂不支持该币种的换算"
	}
	res := &model.StatisticsString{
		ID:         0,
		IDentifier: statistic.IDentifier,
		ChannelID:  statistic.ChannelID,
		AccountID:  statistic.AccountID,
		Cost:       cost,
		CostCny:    cny,
		Click:      decimal.NewFromInt(statistic.Click).String(),
		Imp:        decimal.NewFromInt(statistic.Imp).String(),
		Ctr:        decimal.NewFromFloat(statistic.Ctr).Round(2).String() + "%",
		Cv:         decimal.NewFromInt(statistic.Cv).String(),
		Cvr:        decimal.NewFromFloat(statistic.Cvr).Round(2).String() + "%",
		Cpa:        decimal.NewFromFloat(statistic.Cpa).Round(2).String(),
		Status:     statistic.Status,
		Type:       statistic.Type,
		Currency:   statistic.Currency,
		RecordDate: statistic.RecordDate,
		CreateAt:   statistic.CreateAt,
		UpdateAt:   statistic.UpdateAt,
	}
	return res
}
