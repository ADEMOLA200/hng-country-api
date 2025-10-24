package models

import (
	"time"
)

type Country struct {
	ID              uint      `json:"id" gorm:"primaryKey"`
	Name            string    `json:"name" gorm:"type:varchar(255);uniqueIndex;not null"`
	Capital         string    `json:"capital" gorm:"type:varchar(255)"`
	Region          string    `json:"region" gorm:"type:varchar(255)"`
	Population      int64     `json:"population" gorm:"not null"`
	CurrencyCode    string    `json:"currency_code" gorm:"type:varchar(10);not null"`
	ExchangeRate    *float64  `json:"exchange_rate"`
	EstimatedGDP    *float64  `json:"estimated_gdp"`
	FlagURL         string    `json:"flag_url" gorm:"type:text"`
	LastRefreshedAt time.Time `json:"last_refreshed_at"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type CountryResponse struct {
	ID              uint      `json:"id"`
	Name            string    `json:"name"`
	Capital         string    `json:"capital"`
	Region          string    `json:"region"`
	Population      int64     `json:"population"`
	CurrencyCode    string    `json:"currency_code"`
	ExchangeRate    *float64  `json:"exchange_rate"`
	EstimatedGDP    *float64  `json:"estimated_gdp"`
	FlagURL         string    `json:"flag_url"`
	LastRefreshedAt time.Time `json:"last_refreshed_at"`
}

type RefreshResponse struct {
	Message         string    `json:"message"`
	TotalCountries  int       `json:"total_countries"`
	LastRefreshedAt time.Time `json:"last_refreshed_at"`
}

type StatusResponse struct {
	TotalCountries  int       `json:"total_countries"`
	LastRefreshedAt time.Time `json:"last_refreshed_at"`
}

type ErrorResponse struct {
	Error   string            `json:"error"`
	Details map[string]string `json:"details,omitempty"`
}
