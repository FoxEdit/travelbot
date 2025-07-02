package domain

import (
	"encoding/json"
	"errors"

	"github.com/shopspring/decimal"
)

// The amount contains the minimum possible monetary unit
type Currency struct {
	currencyCode string
	amount       decimal.Decimal
}

func NewCurrency(currencyAmount decimal.Decimal, currencyCode string) (Currency, error) {
	if currencyCode == "" {
		return Currency{}, errors.New("currency code cannot be empty")
	}

	return Currency{
		currencyCode: currencyCode,
		amount:       currencyAmount,
	}, nil
}

func (c Currency) MarshalJSON() ([]byte, error) {
	type alias struct {
		CurrencyCode string          `json:"currency_code"`
		Amount       decimal.Decimal `json:"amount"`
	}

	return json.Marshal(&alias{
		CurrencyCode: c.currencyCode,
		Amount:       c.amount,
	})
}

func (c *Currency) UnmarshalJSON(data []byte) error {
	type alias struct {
		CurrencyCode string          `json:"currency_code"`
		Amount       decimal.Decimal `json:"amount"`
	}

	var a alias
	if err := json.Unmarshal(data, &a); err != nil {
		return err
	}

	newCurrency, err := NewCurrency(a.Amount, a.CurrencyCode)
	if err != nil {
		return err
	}

	*c = newCurrency
	return nil
}

func (c *Currency) GetAmount() decimal.Decimal {
	return c.amount
}

func (c *Currency) GetCurrencyCode() string {
	return c.currencyCode
}
