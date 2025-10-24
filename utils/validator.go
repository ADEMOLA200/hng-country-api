package utils

import (
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

type Validator struct{}

func NewValidator() *Validator {
	return &Validator{}
}

func (v *Validator) ValidateCountryName(name string) (bool, string) {
	if strings.TrimSpace(name) == "" {
		return false, "country name is required"
	}

	if len(name) > 255 {
		return false, "country name too long"
	}

	for _, char := range name {
		if !unicode.IsLetter(char) && !unicode.IsSpace(char) && char != '-' && char != '.' && char != ',' {
			return false, "country name contains invalid characters"
		}
	}

	return true, ""
}

func (v *Validator) ValidatePopulation(population int64) (bool, string) {
	if population < 0 {
		return false, "population cannot be negative"
	}

	if population > 20000000000 { // 20 billion
		return false, "population value too large"
	}

	return true, ""
}

func (v *Validator) ValidateCurrencyCode(currency string) (bool, string) {
	if currency == "" {
		return true, "" // Empty is allowed per requirements
	}

	if len(currency) > 10 {
		return false, "currency code too long"
	}

	matched, _ := regexp.MatchString(`^[A-Za-z]{1,10}$`, currency)
	if !matched {
		return false, "currency code should contain only letters"
	}

	return true, ""
}

func (v *Validator) ValidateExchangeRate(rate *float64) (bool, string) {
	if rate == nil {
		return true, "" // Nil is allowed
	}

	if *rate < 0 {
		return false, "exchange rate cannot be negative"
	}

	if *rate > 1000000 { // 1 million
		return false, "exchange rate too large"
	}

	return true, ""
}

func (v *Validator) ValidateRegion(region string) (bool, string) {
	if region == "" {
		return true, ""
	}

	if len(region) > 255 {
		return false, "region name too long"
	}

	validRegions := []string{
		"Africa", "Americas", "Asia", "Europe", "Oceania",
		"Antarctic", "Polar", "North America", "South America",
		"Central America", "Caribbean", "Middle East", "Southeast Asia",
	}

	regionLower := strings.ToLower(region)
	for _, validRegion := range validRegions {
		if strings.ToLower(validRegion) == regionLower {
			return true, ""
		}
	}

	return true, "" // Allow custom regions as well
}

func (v *Validator) ValidateURL(url string) (bool, string) {
	if url == "" {
		return true, ""
	}

	if len(url) > 1000 {
		return false, "URL too long"
	}

	urlRegex := regexp.MustCompile(`^(https?|ftp)://[^\s/$.?#].[^\s]*$`)
	if !urlRegex.MatchString(url) {
		return false, "invalid URL format"
	}

	return true, ""
}

func (v *Validator) ValidateSortParameter(sort string) (bool, string) {
	if sort == "" {
		return true, ""
	}

	validSorts := []string{
		"gdp_asc", "gdp_desc",
		"population_asc", "population_desc",
		"name_asc", "name_desc",
	}

	for _, validSort := range validSorts {
		if sort == validSort {
			return true, ""
		}
	}

	return false, "invalid sort parameter"
}

func (v *Validator) ValidateCountryData(data map[string]interface{}) map[string]string {
	errors := make(map[string]string)

	if name, exists := data["name"]; !exists || name == "" {
		errors["name"] = "is required"
	} else if nameStr, ok := name.(string); ok {
		if valid, msg := v.ValidateCountryName(nameStr); !valid {
			errors["name"] = msg
		}
	}

	if population, exists := data["population"]; !exists {
		errors["population"] = "is required"
	} else {
		switch pop := population.(type) {
		case int64:
			if valid, msg := v.ValidatePopulation(pop); !valid {
				errors["population"] = msg
			}
		case float64:
			if valid, msg := v.ValidatePopulation(int64(pop)); !valid {
				errors["population"] = msg
			}
		case string:
			if popInt, err := strconv.ParseInt(pop, 10, 64); err == nil {
				if valid, msg := v.ValidatePopulation(popInt); !valid {
					errors["population"] = msg
				}
			} else {
				errors["population"] = "must be a valid number"
			}
		default:
			errors["population"] = "must be a number"
		}
	}

	if currency, exists := data["currency_code"]; exists {
		if currencyStr, ok := currency.(string); ok {
			if valid, msg := v.ValidateCurrencyCode(currencyStr); !valid {
				errors["currency_code"] = msg
			}
		}
	}

	if region, exists := data["region"]; exists {
		if regionStr, ok := region.(string); ok {
			if valid, msg := v.ValidateRegion(regionStr); !valid {
				errors["region"] = msg
			}
		}
	}

	if flagURL, exists := data["flag_url"]; exists {
		if flagStr, ok := flagURL.(string); ok {
			if valid, msg := v.ValidateURL(flagStr); !valid {
				errors["flag_url"] = msg
			}
		}
	}

	return errors
}

func (v *Validator) SanitizeString(input string) string {
	input = strings.ReplaceAll(input, "\x00", "")

	input = strings.TrimSpace(input)

	return input
}

func (v *Validator) IsValidUUID(uuid string) bool {
	uuidRegex := regexp.MustCompile(`^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}$`)
	return uuidRegex.MatchString(uuid)
}
