package services

import (
	"fmt"
	"math/rand"
	"time"
)

type ExchangeService struct {
	cachedRates map[string]float64
	lastFetch   time.Time
	cacheTTL    time.Duration
}

func NewExchangeService() *ExchangeService {
	return &ExchangeService{
		cachedRates: make(map[string]float64),
		cacheTTL:    30 * time.Minute,
	}
}

func (es *ExchangeService) GetExchangeRate(currencyCode string, rates map[string]float64) *float64 {
	if currencyCode == "" {
		return nil
	}

	if rate, exists := rates[currencyCode]; exists {
		return &rate
	}

	if rate, exists := rates[currencyCode+"D"]; exists {
		return &rate
	}

	return nil
}

func (es *ExchangeService) CalculateGDP(population int64, exchangeRate *float64) *float64 {
	if exchangeRate == nil || *exchangeRate == 0 {
		return nil
	}

	rand.Seed(time.Now().UnixNano() + int64(population))
	multiplier := 1000.0 + rand.Float64()*1000.0

	populationFloat := float64(population)
	gdp := (populationFloat * multiplier) / *exchangeRate

	return &gdp
}

func (es *ExchangeService) ValidateCurrencyData(currencyCode string, rates map[string]float64) (bool, string) {
	if currencyCode == "" {
		return false, "currency code is empty"
	}

	if len(currencyCode) > 10 {
		return false, "currency code too long"
	}

	if _, exists := rates[currencyCode]; !exists {
		return false, fmt.Sprintf("currency %s not found in exchange rates", currencyCode)
	}

	return true, ""
}

func (es *ExchangeService) GetSupportedCurrencies(rates map[string]float64) []string {
	currencies := make([]string, 0, len(rates))
	for currency := range rates {
		currencies = append(currencies, currency)
	}
	return currencies
}

func (es *ExchangeService) IsCacheValid() bool {
	return time.Since(es.lastFetch) < es.cacheTTL && len(es.cachedRates) > 0
}

func (es *ExchangeService) SetCachedRates(rates map[string]float64) {
	es.cachedRates = rates
	es.lastFetch = time.Now()
}

func (es *ExchangeService) GetCachedRates() map[string]float64 {
	if es.IsCacheValid() {
		return es.cachedRates
	}
	return nil
}
