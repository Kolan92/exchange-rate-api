package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type ExchangeRate struct {
	Source      string           `json:"source" binding:"required" example:"USD"`
	Destination string           `json:"destination" binding:"required" example:"CHF"`
	Date        time.Time        `json:"date" binding:"required" example:"2022-05-01T00:00:00.00Z"`
	Rate        *decimal.Decimal `json:"rate" example:"1.0456"`
}

type Currency struct {
	Id   int
	Code string
}

func (Currency) TableName() string {
	return "currencies_codes"
}

type DbExchangeRate struct {
	Source      int `gorm:"column:source_currency_id"`
	Destination int `gorm:"column:destination_currency_id"`
	Date        time.Time
	Rate        *decimal.Decimal
}

func (DbExchangeRate) TableName() string {
	return "exchange_rates"
}
