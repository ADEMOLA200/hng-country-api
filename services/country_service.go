package services

import (
	"fmt"
	"time"

	"github.com/ADEMOLA200/hng-country-api/models"
	"github.com/ADEMOLA200/hng-country-api/repositories"
)

type CountryService struct {
	countryRepo     *repositories.CountryRepository
	externalAPI     *ExternalAPIService
	exchangeService *ExchangeService
	imageService    *ImageService
}

func NewCountryService() *CountryService {
	return &CountryService{
		countryRepo:     repositories.NewCountryRepository(),
		externalAPI:     NewExternalAPIService(),
		exchangeService: NewExchangeService(),
		imageService:    NewImageService(),
	}
}

func (s *CountryService) RefreshCountries() (*models.RefreshResponse, error) {
	fmt.Println("Cleaning existing duplicates...")
	if err := s.countryRepo.CleanDuplicates(); err != nil {
		fmt.Printf("Warning: Error cleaning duplicates before refresh: %v\n", err)
	}

	cachedRates := s.exchangeService.GetCachedRates()

	fmt.Println("Fetching countries data...")
	restCountries, err := s.externalAPI.FetchCountries()
	if err != nil {
		return nil, fmt.Errorf("external data source unavailable: %v", err)
	}

	var exchangeRates map[string]float64
	if cachedRates != nil {
		exchangeRates = cachedRates
		fmt.Println("Using cached exchange rates...")
	} else {
		fmt.Println("Fetching fresh exchange rates...")
		exchangeRates, err = s.externalAPI.FetchExchangeRates()
		if err != nil {
			return nil, fmt.Errorf("external data source unavailable: %v", err)
		}
		s.exchangeService.SetCachedRates(exchangeRates)
	}

	currentTime := time.Now()
	processedCount := 0
	updateCount := 0
	createCount := 0

	fmt.Printf("Processing %d countries...\n", len(restCountries))

	for _, restCountry := range restCountries {
		country := models.Country{
			Name:            restCountry.Name,
			Capital:         restCountry.Capital,
			Region:          restCountry.Region,
			Population:      restCountry.Population,
			FlagURL:         restCountry.Flag,
			LastRefreshedAt: currentTime,
		}

		currencyCode := s.externalAPI.GetCurrencyCode(restCountry.Currencies)
		country.CurrencyCode = currencyCode

		if currencyCode != "" {
			exchangeRate := s.exchangeService.GetExchangeRate(currencyCode, exchangeRates)
			country.ExchangeRate = exchangeRate
			country.EstimatedGDP = s.exchangeService.CalculateGDP(restCountry.Population, exchangeRate)
		} else {
			country.CurrencyCode = ""
			country.ExchangeRate = nil
			country.EstimatedGDP = nil
		}

		existing, err := s.countryRepo.CheckIfExists(country.Name)

		if err == nil {
			country.ID = existing.ID
			err = s.countryRepo.Update(&country)
			updateCount++
		} else {
			err = s.countryRepo.Create(&country)
			createCount++
		}

		if err != nil {
			fmt.Printf("Error processing country %s: %v\n", country.Name, err)
			continue
		}

		processedCount++
	}

	fmt.Println("Cleaning duplicates after refresh...")
	if err := s.countryRepo.CleanDuplicates(); err != nil {
		fmt.Printf("Warning: Error cleaning duplicates after refresh: %v\n", err)
	}

	actualCount, err := s.countryRepo.GetTotalCount()
	if err != nil {
		fmt.Printf("Warning: Could not get actual count: %v\n", err)
		actualCount = int64(processedCount)
	}

	fmt.Printf("Refresh completed: %d processed (%d created, %d updated), %d total in DB\n",
		processedCount, createCount, updateCount, actualCount)

	fmt.Println("Generating summary image...")
	if err := s.imageService.GenerateSummaryImage(); err != nil {
		fmt.Printf("Error generating summary image: %v\n", err)
	}

	response := &models.RefreshResponse{
		Message:         "Countries data refreshed successfully",
		TotalCountries:  int(actualCount),
		LastRefreshedAt: currentTime,
	}

	return response, nil
}

func (s *CountryService) GetAllCountries(filters map[string]string, sort string) ([]models.Country, error) {
	return s.countryRepo.FindAll(filters, sort)
}

func (s *CountryService) GetCountryByName(name string) (*models.Country, error) {
	return s.countryRepo.FindByName(name)
}

func (s *CountryService) DeleteCountryByName(name string) error {
	return s.countryRepo.DeleteByName(name)
}

func (s *CountryService) GetStatus() (*models.StatusResponse, error) {
	total, err := s.countryRepo.GetTotalCount()
	if err != nil {
		return nil, err
	}

	countries, err := s.countryRepo.GetLastRefreshTime()
	if err != nil {
		return nil, err
	}

	var lastRefresh time.Time
	if len(countries) > 0 {
		lastRefresh = countries[0].LastRefreshedAt
	}

	return &models.StatusResponse{
		TotalCountries:  int(total),
		LastRefreshedAt: lastRefresh,
	}, nil
}

func (s *CountryService) GetCountryNames() ([]string, error) {
	return s.countryRepo.GetCountryNames()
}
