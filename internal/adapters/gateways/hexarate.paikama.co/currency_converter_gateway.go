package hexaratepaikamaco

import (
	"encoding/json"
	"io"
	"net/http"
	"travelWallet/internal/domain"
	"travelWallet/internal/usecase"

	"github.com/shopspring/decimal"
)

var _ usecase.CurrencyConverter = (*CurrencyConverterGateway)(nil)

type CurrencyConverterGateway struct {
	baseURL string
}

func NewConverter() *CurrencyConverterGateway {
	return &CurrencyConverterGateway{baseURL: "https://hexarate.paikama.co/api/rates/latest"}
}

// ConvertToAnother implements usecase.CurrencyConverter.
func (c *CurrencyConverterGateway) ConvertToAnother(currency domain.Currency, another string) (domain.Currency, error) {
	resp, err := http.Get(c.baseURL + "/" + currency.GetCurrencyCode() + "?target=" + another)
	if err != nil {
		panic("unhandled error")
	}

	body, _ := io.ReadAll(resp.Body)

	var responseStruct struct {
		Data struct {
			Mid float64 `json:"mid"`
		} `json:"data"`
	}

	json.Unmarshal(body, &responseStruct)

	curr, _ := domain.NewCurrency(decimal.NewFromFloat(responseStruct.Data.Mid).Mul(currency.GetAmount()), another)
	return curr, nil
}

// ConvertToMany implements usecase.CurrencyConverter.
func (c *CurrencyConverterGateway) ConvertToMany(currency domain.Currency, anothers ...string) ([]domain.Currency, error) {
	result := make([]domain.Currency, len(anothers))
	for i, v := range anothers {
		itResult, err := c.ConvertToAnother(currency, v)
		if err != nil {
			continue
		}
		result[i] = itResult
	}

	return result, nil
}
