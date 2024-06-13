package api

import "time"

type CompanyProfile struct {
	Symbol       string
	Description  string
	SecurityType string
	Price        float64 //current price
	Open         float64
	DailyHigh    float64
	DailyLow     float64
	YearHigh     float64    //52w high
	YearLow      float64    //52w low
	MarketCap    float64    // market cap in Millions
	DivYield     float64    //dividendYieldIndicatedAnnual
	PEratio      float64    //peExclExtraTTM
	Updated      *time.Time // last time info was updated
}

func NewCompanyProfile(symbol string, description string) *CompanyProfile {
	return &CompanyProfile{
		Symbol:      symbol,
		Description: description,
	}
}
