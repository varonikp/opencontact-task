package services

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/varonikp/opencontact-task/internal/domain"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"time"
)

type ratesRepo interface {
	Insert(ctx context.Context, currency *domain.CurrencyRate) (*domain.CurrencyRate, error)
}

type RatesService struct {
	httpClient *http.Client
	r          ratesRepo
}

type NBRBRate struct {
	CurID           int     `json:"Cur_ID"`
	Date            string  `json:"Date"`
	CurAbbreviation string  `json:"Cur_Abbreviation"`
	CurScale        int     `json:"Cur_Scale"`
	CurName         string  `json:"Cur_Name"`
	CurOfficialRate float64 `json:"Cur_OfficialRate"`
}

func NewRatesService(httpClient *http.Client, r ratesRepo) *RatesService {
	return &RatesService{
		httpClient: httpClient,
		r:          r,
	}
}

func (s *RatesService) Start(ctx context.Context) error {
	if inserted, err := s.collectRates(ctx); err != nil {
		return fmt.Errorf("failed to collect rates")
	} else {
		slog.Debug("successfully collected rates", slog.Int("count", inserted))
	}

	// TODO: change time ticker to go-cron
	go func() {
		ticker := time.NewTicker(24 * time.Hour)
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				inserted, err := s.collectRates(ctx)
				if err != nil {
					slog.Error("failed to collect rates", slog.String("error", err.Error()))
					continue
				}

				slog.Debug("successfully collected rates", slog.Int("count", inserted))
			}
		}
	}()

	return nil
}

func (s *RatesService) collectRates(ctx context.Context) (int, error) {
	rates, err := s.fetchRates(ctx)
	if err != nil {
		return -1, fmt.Errorf("failed to fetch rates: %w", err)
	}

	var inserted int
	// TODO: batch insert
	for _, rate := range rates {
		domainRate := domain.NewCurrencyRate(domain.NewCurrencyRateData{
			CurrencyID:   rate.CurID,
			Name:         rate.CurName,
			Abbreviation: rate.CurAbbreviation,
			Rate:         rate.CurOfficialRate,
		})
		_, err := s.r.Insert(ctx, domainRate)

		if err != nil {
			slog.Warn("error happened while inserting currency rate", slog.Any("rate", domainRate), slog.String("error", err.Error()))
			continue
		}

		inserted++
	}

	return inserted, nil
}

func (s *RatesService) fetchRates(ctx context.Context) ([]*NBRBRate, error) {
	uri := url.URL{
		Scheme: "https",
		Host:   "api.nbrb.by",
		Path:   "exrates/rates",
		RawQuery: url.Values{
			"periodicity": {"0"},
		}.Encode(),
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request with context: %w", err)
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned unexpected response: %s", resp.Status)
	}

	var jsonResp []*NBRBRate
	if err := json.NewDecoder(resp.Body).Decode(&jsonResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return jsonResp, nil
}
