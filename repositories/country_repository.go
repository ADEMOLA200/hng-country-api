package repositories

import (
	"strings"

	"github.com/ADEMOLA200/hng-country-api/models"
	"gorm.io/gorm"
)

type CountryRepository struct {
	db *gorm.DB
}

func NewCountryRepository() *CountryRepository {
	return &CountryRepository{db: DB}
}

func (r *CountryRepository) Create(country *models.Country) error {
	return r.db.Create(country).Error
}

func (r *CountryRepository) Update(country *models.Country) error {
	return r.db.Save(country).Error
}

func (r *CountryRepository) FindByName(name string) (*models.Country, error) {
	var country models.Country

	cleanName := r.normalizeCountryName(name)

	err := r.db.Where("LOWER(name) = LOWER(?)", cleanName).First(&country).Error
	if err == nil {
		return &country, nil
	}

	if err == gorm.ErrRecordNotFound {
		variations := r.getCountryNameVariations(cleanName)
		for _, variation := range variations {
			err = r.db.Where("LOWER(name) = LOWER(?)", variation).First(&country).Error
			if err == nil {
				return &country, nil
			}
		}

		err = r.db.Where("LOWER(name) LIKE LOWER(?)", "%"+cleanName+"%").First(&country).Error
		if err == nil {
			return &country, nil
		}
	}

	return nil, err
}

func (r *CountryRepository) FindAll(filters map[string]string, sort string) ([]models.Country, error) {
	var countries []models.Country
	query := r.db.Model(&models.Country{})

	if region, exists := filters["region"]; exists {
		query = query.Where("LOWER(region) = LOWER(?)", region)
	}
	if currency, exists := filters["currency"]; exists {
		query = query.Where("LOWER(currency_code) = LOWER(?)", currency)
	}

	switch sort {
	case "gdp_asc":
		query = query.Order("estimated_gdp ASC")
	case "gdp_desc":
		query = query.Order("estimated_gdp DESC")
	case "population_asc":
		query = query.Order("population ASC")
	case "population_desc":
		query = query.Order("population DESC")
	case "name_asc":
		query = query.Order("name ASC")
	case "name_desc":
		query = query.Order("name DESC")
	default:
		query = query.Order("name ASC")
	}

	err := query.Find(&countries).Error
	return countries, err
}

func (r *CountryRepository) DeleteByName(name string) error {
	return r.db.Where("LOWER(name) = LOWER(?)", name).Delete(&models.Country{}).Error
}

func (r *CountryRepository) GetTotalCount() (int64, error) {
	var count int64
	err := r.db.Model(&models.Country{}).Count(&count).Error
	return count, err
}

func (r *CountryRepository) GetLastRefreshTime() ([]models.Country, error) {
	var countries []models.Country
	err := r.db.Order("last_refreshed_at DESC").Limit(1).Find(&countries).Error
	return countries, err
}

func (r *CountryRepository) Upsert(country *models.Country) error {
	var existing models.Country
	err := r.db.Where("LOWER(name) = LOWER(?)", country.Name).First(&existing).Error

	if err == gorm.ErrRecordNotFound {
		return r.Create(country)
	} else if err != nil {
		return err
	}

	country.ID = existing.ID
	return r.Update(country)
}

func (r *CountryRepository) GetTopCountriesByGDP(limit int) ([]models.Country, error) {
	var countries []models.Country
	err := r.db.Where("estimated_gdp IS NOT NULL").Order("estimated_gdp DESC").Limit(limit).Find(&countries).Error
	return countries, err
}

func (r *CountryRepository) CleanDuplicates() error {
	result := r.db.Exec(`
		DELETE c1 FROM countries c1
		INNER JOIN countries c2 
		WHERE 
			c1.id < c2.id AND 
			LOWER(c1.name) = LOWER(c2.name)
	`)

	if result.Error != nil {
		return result.Error
	}

	result = r.db.Exec(`
		DELETE c1 FROM countries c1
		LEFT JOIN (
			SELECT MAX(id) as max_id, LOWER(name) as name_lower
			FROM countries 
			GROUP BY name_lower
		) c2 ON c1.id = c2.max_id
		WHERE c2.max_id IS NULL
	`)

	return result.Error
}

func (r *CountryRepository) CheckIfExists(name string) (*models.Country, error) {
	var country models.Country
	err := r.db.Where("LOWER(name) = LOWER(?)", name).First(&country).Error
	if err != nil {
		return nil, err
	}
	return &country, nil
}

func (r *CountryRepository) normalizeCountryName(name string) string {
	cleanName := strings.TrimSpace(name)

	replacements := map[string]string{
		", State of":                "",
		" (the)":                    "",
		", Republic of":             "",
		" Republic":                 "",
		" (Bolivarian Republic of)": "",
		" (Islamic Republic of)":    "",
		" (Plurinational State of)": "",
		" Democratic Republic":      "",
		" People's Republic":        "",
	}

	for old, new := range replacements {
		cleanName = strings.Replace(cleanName, old, new, -1)
	}

	return strings.TrimSpace(cleanName)
}

func (r *CountryRepository) getCountryNameVariations(name string) []string {
	variations := []string{name}
	cleanName := r.normalizeCountryName(name)

	if cleanName != name {
		variations = append(variations, cleanName)
	}

	commonVariations := map[string][]string{
		"Palestine, State of":                {"Palestine"},
		"United States of America":           {"United States", "USA", "US"},
		"United Kingdom":                     {"UK", "Great Britain"},
		"Russian Federation":                 {"Russia"},
		"Iran (Islamic Republic of)":         {"Iran"},
		"Venezuela (Bolivarian Republic of)": {"Venezuela"},
		"Bolivia (Plurinational State of)":   {"Bolivia"},
		"Viet Nam":                           {"Vietnam"},
		"Brunei Darussalam":                  {"Brunei"},
		"Congo (the)":                        {"Congo"},
		"Democratic Republic of the Congo":   {"DR Congo", "Democratic Republic of Congo"},
		"Tanzania, United Republic of":       {"Tanzania"},
		"Syrian Arab Republic":               {"Syria"},
		"Lao People's Democratic Republic":   {"Laos"},
		"Republic of Moldova":                {"Moldova"},
		"Republic of the Congo":              {"Congo Republic"},
		"North Macedonia":                    {"Macedonia"},
		"Eswatini":                           {"Swaziland"},
		"Czech Republic":                     {"Czechia"},
		"Timor-Leste":                        {"East Timor"},
		"Myanmar":                            {"Burma"},
		"Cabo Verde":                         {"Cape Verde"},
		"Côte d'Ivoire":                      {"Ivory Coast"},
	}

	if vars, exists := commonVariations[name]; exists {
		variations = append(variations, vars...)
	}

	return variations
}

func (r *CountryRepository) GetCountryNames() ([]string, error) {
	var names []string
	err := r.db.Model(&models.Country{}).Pluck("name", &names).Error
	return names, err
}
