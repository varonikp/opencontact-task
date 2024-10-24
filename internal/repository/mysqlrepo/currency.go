package mysqlrepo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/varonikp/opencontact-task/internal/domain"
	"time"
)

const (
	paginationLimit = 10
)

type DBTX interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

type CurrencyRatesRepository struct {
	db DBTX
}

func NewCurrencyRepository(db DBTX) *CurrencyRatesRepository {
	return &CurrencyRatesRepository{
		db: db,
	}
}

// Insert returns currency with inserted unique id
func (r *CurrencyRatesRepository) Insert(ctx context.Context, currency *domain.CurrencyRate) (*domain.CurrencyRate, error) {
	op := "mysqlrepo.CurrencyRatesRepository.Insert"

	query := "INSERT INTO rates(currency_id, name, abbreviation, rate, inserted_at) VALUES (?, ?, ?, ?, ?)"

	result, err := r.db.ExecContext(ctx, query, currency.CurrencyID(), currency.Name(), currency.Abbreviation(), currency.Rate(), time.Now())
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	lastInsertedID, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("last inserted id: %w", err)
	}

	currency.WithUniqueID(int(lastInsertedID))

	return currency, nil
}

// GetCurrenciesRates fetch all rates with pagination
func (r *CurrencyRatesRepository) GetCurrenciesRates(ctx context.Context, page int) ([]*domain.CurrencyRate, error) {
	op := "mysqlrepo.CurrencyRatesRepository.GetByCurrencyID"

	query := "SELECT * FROM rates ORDER BY inserted_at, currency_id LIMIT ? OFFSET ?"

	offset := (page - 1) * paginationLimit

	rows, err := r.db.QueryContext(ctx, query, paginationLimit, offset)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("%s: %w", op, domain.ErrNotFound)
	}
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var out []*domain.CurrencyRate

	for rows.Next() {
		var (
			uniqueID     int
			currencyID   int
			name         string
			abbreviation string
			rate         float64
			insertedAt   time.Time
		)

		err = rows.Scan(
			&uniqueID, &currencyID, &name, &abbreviation, &rate, &insertedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		out = append(out, domain.NewCurrencyRate(domain.NewCurrencyRateData{
			UniqueID:     uniqueID,
			CurrencyID:   currencyID,
			Name:         name,
			Abbreviation: abbreviation,
			Rate:         rate,
			InsertedAt:   insertedAt,
		}))
	}

	return out, nil
}

// GetCurrencyRate get currency rate at certain time
func (r *CurrencyRatesRepository) GetCurrencyRate(ctx context.Context, currencyID int, date time.Time) (*domain.CurrencyRate, error) {
	op := "mysqlrepo.CurrencyRatesRepository.GetCurrencyRate"

	query := "SELECT * FROM rates WHERE currency_id = ? AND inserted_at = ? ORDER BY inserted_at"

	row := r.db.QueryRowContext(ctx, query, currencyID, date)

	var (
		uniqueID     int
		name         string
		abbreviation string
		rate         float64
		insertedAt   time.Time
	)

	err := row.Scan(
		&uniqueID, &currencyID, &name, &abbreviation, &rate, &insertedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("%s: %w", op, domain.ErrNotFound)
	}
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return domain.NewCurrencyRate(domain.NewCurrencyRateData{
		UniqueID:     uniqueID,
		CurrencyID:   currencyID,
		Name:         name,
		Abbreviation: abbreviation,
		Rate:         rate,
		InsertedAt:   insertedAt,
	}), nil
}
