package domain

import "time"

type CurrencyRate struct {
	uniqueId     int
	currencyId   int
	name         string
	abbreviation string
	rate         float64
	insertedAt   time.Time
}

type NewCurrencyRateData struct {
	UniqueID     int
	CurrencyID   int
	Name         string
	Abbreviation string
	Rate         float64
	InsertedAt   time.Time
}

func NewCurrencyRate(data NewCurrencyRateData) *CurrencyRate {
	return &CurrencyRate{
		uniqueId:     data.UniqueID,
		currencyId:   data.CurrencyID,
		name:         data.Name,
		abbreviation: data.Abbreviation,
		rate:         data.Rate,
		insertedAt:   data.InsertedAt,
	}
}

func (c *CurrencyRate) WithUniqueID(id int) {
	c.uniqueId = id
}

func (c *CurrencyRate) UniqueID() int {
	return c.uniqueId
}

func (c *CurrencyRate) CurrencyID() int {
	return c.currencyId
}

func (c *CurrencyRate) Name() string {
	return c.name
}

func (c *CurrencyRate) Abbreviation() string {
	return c.abbreviation
}

func (c *CurrencyRate) Rate() float64 {
	return c.rate
}

func (c *CurrencyRate) InsertedAt() time.Time {
	return c.insertedAt
}
