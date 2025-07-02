package domain

import (
	"encoding/json"
	"errors"
	"fmt"
	"slices"

	"github.com/shopspring/decimal"
)

type Wallet struct {
	id                string
	baseCurrency      Currency
	foreignCurrencies []string
}

func NewWallet(ID string, baseCurrency string) (*Wallet, error) {
	initialBaseCurrency, err := NewCurrency(decimal.Zero, baseCurrency)
	if err != nil {
		return nil, err
	}

	return &Wallet{baseCurrency: initialBaseCurrency, id: ID, foreignCurrencies: []string{}}, nil
}

func (w *Wallet) MarshalJSON() ([]byte, error) {
	type alias struct {
		ID                string   `json:"id"`
		BaseCurrency      Currency `json:"base_currency"`
		ForeignCurrencies []string `json:"foreign_currencies"`
	}

	return json.Marshal(&alias{
		ID:                w.id,
		BaseCurrency:      w.baseCurrency,
		ForeignCurrencies: w.foreignCurrencies,
	})
}

func (w *Wallet) UnmarshalJSON(data []byte) error {
	type alias struct {
		ID                string   `json:"id"`
		BaseCurrency      Currency `json:"base_currency"`
		ForeignCurrencies []string `json:"foreign_currencies"`
	}

	var a alias
	if err := json.Unmarshal(data, &a); err != nil {
		return err
	}

	if a.ID == "" {
		return errors.New("wallet id is missing in json")
	}

	w.id = a.ID
	w.baseCurrency = a.BaseCurrency

	if a.ForeignCurrencies == nil {
		w.foreignCurrencies = []string{}
	} else {
		w.foreignCurrencies = a.ForeignCurrencies
	}

	return nil
}

func (w *Wallet) GetID() string {
	return w.id
}

func (w *Wallet) GetBaseCurrencyCode() string {
	return w.baseCurrency.currencyCode
}

func (w *Wallet) ChangeBaseCurrency(currency string) {
	w.baseCurrency.currencyCode = currency
}

func (w *Wallet) AddForeignCurrency(currency string) {
	if !slices.Contains(w.foreignCurrencies, currency) {
		w.foreignCurrencies = append(w.foreignCurrencies, currency)
	}
}

func (w *Wallet) RemoveForeignCurrency(currency string) {
	if slices.Contains(w.foreignCurrencies, currency) {
		// slice order doesn't matter
		removeIndex := slices.Index(w.foreignCurrencies, currency)
		w.foreignCurrencies[removeIndex] = w.foreignCurrencies[len(w.foreignCurrencies)-1]
		w.foreignCurrencies = w.foreignCurrencies[:len(w.foreignCurrencies)-1]
	}
}

func (w *Wallet) GetAllForeignCurrencies() []string {
	return w.foreignCurrencies
}

func (w *Wallet) RemoveAllForeignCurrencies() {
	w.foreignCurrencies = []string{}
}

func (w *Wallet) Deposit(amount Currency) error {
	if amount.currencyCode != w.baseCurrency.currencyCode {
		return fmt.Errorf("cannot deposit currency %s to a %s wallet", amount.currencyCode, w.baseCurrency.currencyCode)
	}

	w.baseCurrency.amount = w.baseCurrency.amount.Add(amount.amount)
	return nil
}

func (w *Wallet) Withdraw(amount Currency) error {
	if amount.currencyCode != w.baseCurrency.currencyCode {
		return fmt.Errorf("cannot withdraw currency %s from a %s wallet", amount.currencyCode, w.baseCurrency.currencyCode)
	}

	if w.baseCurrency.amount.Compare(amount.amount) == -1 {
		return errors.New("insufficient funds")
	}

	w.baseCurrency.amount = w.baseCurrency.amount.Sub(amount.amount)
	return nil
}

func (w *Wallet) GetRawBalance() decimal.Decimal {
	return w.baseCurrency.amount
}
