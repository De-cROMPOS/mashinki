package parser

import (
	"encoding/json"
	"fmt"
	"mashinki/logging"
	"mashinki/translations"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// Map of known drive types
var driveTypes = map[string]string{
	"中置四驱": "Полный привод",
	"前置后驱": "Переднемоторный задний привод",
	"前置前驱": "Передний привод",
	"中置后驱": "Среднемоторный задний привод",
	"后置后驱": "Задний привод",
}

// parseMileage extracts numeric value and converts according to units
func parseMileage(mileageStr string) string {
	// Remove spaces
	mileageStr = strings.TrimSpace(mileageStr)

	// Find number in the string
	re := regexp.MustCompile(`[\d.]+`)
	numbers := re.FindString(mileageStr)
	if numbers == "" {
		return mileageStr
	}

	value, err := strconv.ParseFloat(numbers, 64)
	if err != nil {
		return mileageStr
	}

	// Convert to kilometers
	totalKm := value * 10000

	return fmt.Sprintf("%.0f км", totalKm)
}

// getCarId extracts car ID from the page URL
func getCarId(urlStr string) (string, error) {
	// Check if it's a mobile URL
	if strings.Contains(urlStr, "m.che168.com") {
		// Parse URL to get query parameters
		parsed, err := url.Parse(urlStr)
		if err == nil {
			values := parsed.Query()
			if id := values.Get("infoid"); id != "" {
				return id, nil
			}
		}
	}

	// If not mobile or failed to parse mobile URL
	parts := strings.Split(urlStr, "/")
	for _, part := range parts {
		if strings.Contains(part, ".html") {
			return part[:strings.Index(part, ".html")], nil
		}
	}

	return "", fmt.Errorf("car id not found in url: %s", urlStr)
}

// getCarConfig retrieves basic car information: name, price, year, mileage
func getCarConfig(url string, CI *carInfo) error {
	var err error
	CI.CarId, err = getCarId(url)
	if err != nil {
		return fmt.Errorf("failed to get car ID: %v", err)
	}

	carInfoUrl := fmt.Sprintf("https://www.che168.com/CarConfig/CarConfig.html?infoid=%s", CI.CarId)

	resp, err := makeRequest(carInfoUrl, 0)
	if err != nil {
		return fmt.Errorf("failed to make request: %v", err)
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(resp))
	if err != nil {
		return fmt.Errorf("error while parsing html: %v", err)
	}

	// Getting spec id
	CI.SpecID = doc.Find("#CarSpecid").AttrOr("value", "")

	// Getting car price
	priceStr := doc.Find("#car_price").AttrOr("value", "0")
	priceConv, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		return fmt.Errorf("failed to parse price: %v", err)
	}
	CI.Price = priceConv * 10_000

	// getting full car name
	fullName := doc.Find(".source-info-con h3 a").Text()
	CI.FullName = translations.Translate(fullName)

	// getting mileage and year
	infoText := doc.Find(".source-info-con p").First().Text()
	if idx := strings.Index(infoText, "／"); idx != -1 {
		CI.Milage = parseMileage(strings.TrimSpace(infoText[:idx]))

		// Getting year
		rest := infoText[idx+3:] // skipping the first symbol
		if idx2 := strings.Index(rest, "／"); idx2 != -1 {
			CI.Year = strings.TrimSpace(rest[:idx2])
			if CI.Year == "未上牌" {
				CI.Year = "Еще не ставился на учет"
			}
		}
	}

	return nil
}

// getCarSpecInfo retrieves technical specifications: power, engine size, drive type, fuel type
func getCarSpecInfo(CI *carInfo) error {
	carSpecUrl := fmt.Sprintf("https://cacheapigo.che168.com/CarProduct/GetParam.ashx?specid=%s&callback=configTitle", CI.SpecID)
	specs, err := makeRequest(carSpecUrl, 1)
	if err != nil {
		return fmt.Errorf("failed to get specs: %v", err)
	}

	start := strings.Index(specs, "(")
	end := strings.LastIndex(specs, ")")
	if start == -1 || end == -1 {
		return fmt.Errorf("invalid response format")
	}
	jsonStr := specs[start+1 : end]

	// Decoding JSON
	specsJSON := &SpecResponse{}
	if err := json.Unmarshal([]byte(jsonStr), &specsJSON); err != nil {
		return fmt.Errorf("error while parsing JSON: %v", err)
	}

	// Searching for power and engine size in characteristics
	if specsJSON != nil && len(specsJSON.Result.ParamTypeItems) > 0 {
		for _, group := range specsJSON.Result.ParamTypeItems {
			for _, param := range group.ParamItems {
				name := param.Name
				value := param.Value

				// Searching for power by (kW), taking maximum value
				if strings.Contains(name, "(kW)") {
					if intVal, err := strconv.Atoi(value); err == nil {
						if intVal > CI.Power {
							CI.Power = intVal
						}
					}
				}

				// Searching for engine size by (mL)
				if strings.Contains(name, "(mL)") {
					if intVal, err := strconv.Atoi(value); err == nil {
						CI.EngineSize = intVal
					}
				}

				// Searching for drive type
				if strings.Contains(name, "驱动方式") {
					knownType, exists := driveTypes[value]
					if exists {
						CI.Drive = knownType
					} else {
						logging.TranslationsLogger.LogErrorF("Unknown type of drive: %v", value)
						CI.Drive = translations.Translate(value)
					}
				}

				// Searching for fuel type
				if strings.Contains(name, "燃料形式") {
					CI.FuelType = translations.Translate(value)
				}
			}
		}
	}
	return nil
}

// GetCarInfo retrieves complete car information by URL.
// Returns a structure with car information or an error if something went wrong.
func GetCarInfo(url string) (carInfo, error) {
	carInformation := &carInfo{}

	// getting full name, price, year, mileage
	err := getCarConfig(url, carInformation)
	if err != nil {
		return carInfo{}, fmt.Errorf("failed to get car config: %v", err)
	}

	if carInformation.SpecID != "" {
		// getting car power, engine size, drive, fuel type
		err = getCarSpecInfo(carInformation)
		if err != nil {
			return carInfo{}, fmt.Errorf("failed to get car specs: %v", err)
		}
		return *carInformation, nil
	}

	return carInfo{}, nil
}
