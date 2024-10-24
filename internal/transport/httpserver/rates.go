package httpserver

import (
	"context"
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/varonikp/opencontact-task/internal/domain"
	"net/http"
	"time"
)

type ratesRepo interface {
	GetCurrenciesRates(ctx context.Context, page int) ([]*domain.CurrencyRate, error)
	GetCurrencyRate(ctx context.Context, currencyID int, date time.Time) (*domain.CurrencyRate, error)
}

type RatesHandler struct {
	ctx       context.Context
	ratesRepo ratesRepo
}

func NewRatesHandler(ctx context.Context, ratesRepo ratesRepo) *RatesHandler {
	return &RatesHandler{
		ctx:       ctx,
		ratesRepo: ratesRepo,
	}
}

type GetCurrencyRateQuery struct {
	CurrencyID int  `param:"currency_id"`
	Date       Date `query:"date"`
}

type GetCurrencyRateResponse struct {
	Data struct {
		Name         string  `json:"name"`
		Abbreviation string  `json:"abbreviation"`
		Rate         float64 `json:"rate"`
	}
}

func (h *RatesHandler) GetCurrencyRate(c echo.Context) error {

	var query GetCurrencyRateQuery
	if err := c.Bind(&query); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	timeoutCtx, cancel := context.WithTimeout(h.ctx, 5*time.Second)
	rate, err := h.ratesRepo.GetCurrencyRate(timeoutCtx, query.CurrencyID, time.Time(query.Date))
	cancel()

	if err != nil && errors.Is(err, domain.ErrNotFound) {
		return c.String(http.StatusNotFound, "not found")
	}

	if err != nil {
		return c.String(http.StatusInternalServerError, "internal server error")
	}

	return c.JSON(http.StatusOK, GetCurrencyRateResponse{
		Data: struct {
			Name         string  `json:"name"`
			Abbreviation string  `json:"abbreviation"`
			Rate         float64 `json:"rate"`
		}{
			Name:         rate.Name(),
			Abbreviation: rate.Abbreviation(),
			Rate:         rate.Rate(),
		},
	})
}

type GetCurrenciesRatesQuery struct {
	Page int `query:"page"`
}

type GetCurrenciesRateData struct {
	Name         string    `json:"name"`
	Abbreviation string    `json:"abbreviation"`
	Rate         float64   `json:"rate"`
	Date         time.Time `json:"date"`
}

type GetCurrenciesRatesResponse struct {
	CurrentPage int
	TotalPages  int
	Count       int
	Data        []GetCurrenciesRateData
}

func (h *RatesHandler) GetCurrenciesRates(c echo.Context) error {

	var query GetCurrenciesRatesQuery
	if err := c.Bind(&query); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	if query.Page == 0 {
		return c.String(http.StatusBadRequest, "page must be greater than 0")
	}

	timeoutCtx, cancel := context.WithTimeout(h.ctx, 25*time.Second)
	rates, err := h.ratesRepo.GetCurrenciesRates(timeoutCtx, query.Page)
	cancel()

	if err != nil {
		return c.String(http.StatusInternalServerError, "internal server error")
	}

	respRates := make([]GetCurrenciesRateData, 0, len(rates))
	for _, rate := range rates {
		respRates = append(respRates, GetCurrenciesRateData{
			Name:         rate.Name(),
			Abbreviation: rate.Abbreviation(),
			Rate:         rate.Rate(),
			Date:         rate.InsertedAt(),
		})
	}

	return c.JSON(http.StatusOK, GetCurrenciesRatesResponse{
		CurrentPage: query.Page,
		TotalPages:  0,
		Count:       0,
		Data:        respRates,
	})
}
