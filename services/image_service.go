package services

import (
	"fmt"
	"image"
	"image/color"
	"os"
	"time"

	"github.com/ADEMOLA200/hng-country-api/repositories"
	"github.com/fogleman/gg"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"
)

type ImageService struct {
	countryRepo *repositories.CountryRepository
}

func NewImageService() *ImageService {
	return &ImageService{
		countryRepo: repositories.NewCountryRepository(),
	}
}

func (s *ImageService) GenerateSummaryImage() error {
	if err := os.MkdirAll("cache", 0755); err != nil {
		return fmt.Errorf("failed to create cache directory: %v", err)
	}

	total, err := s.countryRepo.GetTotalCount()
	if err != nil {
		return err
	}

	topCountries, err := s.countryRepo.GetTopCountriesByGDP(5)
	if err != nil {
		return err
	}

	countries, err := s.countryRepo.GetLastRefreshTime()
	if err != nil {
		return err
	}

	var lastRefresh time.Time
	if len(countries) > 0 {
		lastRefresh = countries[0].LastRefreshedAt
	}

	const width = 800
	const height = 600

	dc := gg.NewContext(width, height)

	dc.SetColor(color.White)
	dc.Clear()

	font, err := s.loadFont(20)
	if err != nil {
		return err
	}
	dc.SetFontFace(font)

	dc.SetColor(color.Black)

	// Draw title
	dc.DrawString("Countries Summary", 50, 50)

	// Draw total countries
	dc.DrawString(fmt.Sprintf("Total Countries: %d", total), 50, 100)

	// Draw last refresh time
	dc.DrawString(fmt.Sprintf("Last Refreshed: %s", lastRefresh.Format("2006-01-02 15:04:05")), 50, 140)

	// Draw top 5 countries by GDP
	dc.DrawString("Top 5 Countries by GDP:", 50, 200)

	yPos := 240
	titleFont, _ := s.loadFont(16)
	dc.SetFontFace(titleFont)

	for i, country := range topCountries {
		gdp := "N/A"
		if country.EstimatedGDP != nil {
			gdp = fmt.Sprintf("$%.2f", *country.EstimatedGDP)
		}
		text := fmt.Sprintf("%d. %s: %s", i+1, country.Name, gdp)
		dc.DrawString(text, 70, float64(yPos))
		yPos += 40
	}

	if err := dc.SavePNG("cache/summary.png"); err != nil {
		return fmt.Errorf("failed to save image: %v", err)
	}

	return nil
}

func (s *ImageService) loadFont(size float64) (font.Face, error) {
	font, err := opentype.Parse(goregular.TTF)
	if err != nil {
		return nil, err
	}

	face, err := opentype.NewFace(font, &opentype.FaceOptions{
		Size: size,
		DPI:  72,
	})
	if err != nil {
		return nil, err
	}

	return face, nil
}

func (s *ImageService) GetSummaryImage() (image.Image, error) {
	file, err := os.Open("cache/summary.png")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	return img, nil
}
