package models

import "time"

type ExchangeRate struct {
	Source      string    `json:"source"`
	Destination string    `json:"destination"`
	Date        time.Time `json:"date"`
	Rate        *float64  `json:"rate"`
}

type Currency struct {
	Id   int
	Code string
}

func (Currency) TableName() string {
	return "currencies_codes"
}
