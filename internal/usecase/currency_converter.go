package usecase

import "travelWallet/internal/domain"

type CurrencyConverterUsecase struct {
	converter CurrencyConverter
}

func NewCurrencyConverterService(converter CurrencyConverter) *CurrencyConverterUsecase {
	return &CurrencyConverterUsecase{converter: converter}
}

func (c *CurrencyConverterUsecase) ConvertFromBaseToMany(base domain.Currency, anothers ...string) ([]domain.Currency, error) {
	currs, err := c.converter.ConvertToMany(base, anothers...)
	if err != nil {
		return nil, err
	}

	return currs, err
}

func (c *CurrencyConverterUsecase) ConvertFromBaseToSingle(base domain.Currency, another string) (domain.Currency, error) {
	curr, err := c.converter.ConvertToAnother(base, another)
	if err != nil {
		return domain.Currency{}, err
	}

	return curr, err
}
