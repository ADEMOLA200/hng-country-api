package services

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"
)

type RestCountry struct {
	Name       string     `json:"name"`
	Capital    string     `json:"capital"`
	Region     string     `json:"region"`
	Population int64      `json:"population"`
	Flag       string     `json:"flag"`
	Currencies []Currency `json:"currencies"`
}

type Currency struct {
	Code string `json:"code"`
}

type ExchangeRates struct {
	Rates map[string]float64 `json:"rates"`
}

type ExternalAPIService struct {
	client *http.Client
}

func NewExternalAPIService() *ExternalAPIService {
	return &ExternalAPIService{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (s *ExternalAPIService) FetchCountries() ([]RestCountry, error) {
	resp, err := s.client.Get("https://restcountries.com/v2/all?fields=name,capital,region,population,flag,currencies")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch countries: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("countries API returned status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	var countries []RestCountry
	err = json.Unmarshal(body, &countries)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal countries data: %v", err)
	}

	return countries, nil
}

func (s *ExternalAPIService) FetchExchangeRates() (map[string]float64, error) {
	resp, err := s.client.Get("https://open.er-api.com/v6/latest/USD")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch exchange rates: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("exchange rates API returned status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	var ratesResponse ExchangeRates
	err = json.Unmarshal(body, &ratesResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal exchange rates data: %v", err)
	}

	return ratesResponse.Rates, nil
}

func (s *ExternalAPIService) CalculateEstimatedGDP(population int64, exchangeRate *float64) *float64 {
	if exchangeRate == nil || *exchangeRate == 0 {
		return nil
	}

	rand.Seed(time.Now().UnixNano())
	multiplier := rand.Float64()*1000 + 1000

	populationFloat := float64(population)
	gdp := (populationFloat * multiplier) / *exchangeRate

	return &gdp
}

func (s *ExternalAPIService) GetCurrencyCode(currencies []Currency) string {
	if len(currencies) == 0 {
		return ""
	}
	return currencies[0].Code
}
