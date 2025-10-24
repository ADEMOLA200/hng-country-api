package utils

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func GenerateRandomString(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func ParseInt(s string, defaultValue int) int {
	if s == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(s)
	if err != nil {
		return defaultValue
	}
	return value
}

func ParseInt64(s string, defaultValue int64) int64 {
	if s == "" {
		return defaultValue
	}
	value, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return defaultValue
	}
	return value
}

func ParseFloat(s string, defaultValue float64) float64 {
	if s == "" {
		return defaultValue
	}
	value, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return defaultValue
	}
	return value
}

func ContainsString(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func RemoveDuplicates(slice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, item := range slice {
		if _, value := keys[item]; !value {
			keys[item] = true
			list = append(list, item)
		}
	}
	return list
}

func FormatGDP(gdp *float64) string {
	if gdp == nil {
		return "N/A"
	}

	value := *gdp
	switch {
	case value >= 1e12:
		return fmt.Sprintf("$%.2fT", value/1e12)
	case value >= 1e9:
		return fmt.Sprintf("$%.2fB", value/1e9)
	case value >= 1e6:
		return fmt.Sprintf("$%.2fM", value/1e6)
	default:
		return fmt.Sprintf("$%.2f", value)
	}
}

func FormatNumber(num int64) string {
	str := fmt.Sprintf("%d", num)
	var result strings.Builder
	count := 0

	for i := len(str) - 1; i >= 0; i-- {
		if count == 3 {
			result.WriteByte(',')
			count = 0
		}
		result.WriteByte(str[i])
		count++
	}

	runes := []rune(result.String())
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}

	return string(runes)
}

func GenerateRandomNumber(min, max int) int {
	n, _ := rand.Int(rand.Reader, big.NewInt(int64(max-min+1)))
	return min + int(n.Int64())
}

func MakeHTTPRequest(client *http.Client, url string) ([]byte, error) {
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func ParseJSONResponse(data []byte, target interface{}) error {
	return json.Unmarshal(data, target)
}

func GetTimestamp() time.Time {
	return time.Now().UTC()
}

func FormatTimestamp(t time.Time) string {
	return t.Format(time.RFC3339)
}

func ValidateRequiredFields(data map[string]interface{}, required []string) map[string]string {
	errors := make(map[string]string)

	for _, field := range required {
		if value, exists := data[field]; !exists || value == "" {
			errors[field] = "is required"
		}
	}

	return errors
}

func TruncateString(s string, length int) string {
	if len(s) <= length {
		return s
	}
	return s[:length] + "..."
}
