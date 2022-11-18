package model

import "math"

const MinorUnitsInMajor = 100

type Balance struct {
	UserID      uint64  `json:"userId"`
	AmountMinor int     `json:"-"`
	AmountMajor float64 `json:"balance"`
}

func (t *Balance) ConvertAmountToMajor() {
	t.AmountMajor = math.Round(float64(t.AmountMinor) / MinorUnitsInMajor)
}

func (t *Balance) ConvertAmountToMinor() {
	t.AmountMinor = int(t.AmountMajor * MinorUnitsInMajor)
}

type Transaction struct {
	UserID      uint64  `json:"userId"`
	Type        string  `json:"-"`
	ServiceID   uint64  `json:"serviceId"`
	OrderID     uint64  `json:"orderId"`
	AmountMajor float64 `json:"total"`
	AmountMinor int     `json:"-"`
}

func (t *Transaction) ConvertAmountToMajor() {
	t.AmountMajor = math.Round(float64(t.AmountMinor) / MinorUnitsInMajor)
}

func (t *Transaction) ConvertAmountToMinor() {
	t.AmountMinor = int(t.AmountMajor * MinorUnitsInMajor)
}
