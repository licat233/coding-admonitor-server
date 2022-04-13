package utils

import "github.com/shopspring/decimal"

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
