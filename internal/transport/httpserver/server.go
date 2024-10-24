package httpserver

import (
	"context"
	"github.com/labstack/echo/v4"
)

type HttpServer struct {
	addr string
	e    *echo.Echo
}

func New(ratesHandler *RatesHandler, addr string) *HttpServer {
	e := echo.New()

	e.GET("/rate/:currency_id", ratesHandler.GetCurrencyRate)
	e.GET("/rates", ratesHandler.GetCurrenciesRates)

	return &HttpServer{
		addr: addr,
		e:    e,
	}
}

func (s *HttpServer) Start() error {
	return s.e.Start(s.addr)
}

func (s *HttpServer) Shutdown(ctx context.Context) error {
	return s.e.Shutdown(ctx)
}
