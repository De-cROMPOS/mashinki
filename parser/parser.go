package parser

import (
	"encoding/json"
	"fmt"
	"mashinki/translations"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// TODO: add getter to mob version
func getCarId(url string) (string, error) {

	parts := strings.Split(url, "/")

	for _, part := range parts {
		if strings.Contains(part, ".html") {
			return strings.TrimSuffix(part, ".html"), nil
		}
	}

	return "", fmt.Errorf("car id not found from url: %s", url)
}

func getCarConfig(url string, CI *carInfo) error {

	var err error
	CI.CarId, err = getCarId(url)
	if err != nil {
		return err
	}

	carInfoUrl := fmt.Sprintf("https://www.che168.com/CarConfig/CarConfig.html?infoid=%s", CI.CarId)

	resp, err := makeRequest(carInfoUrl, 1)
	if err != nil {
		return err
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(resp))
	if err != nil {
		return fmt.Errorf("eror while parsing html: %v", err)
	}

	// Getting spec id
	CI.SpecID = doc.Find("#CarSpecid").AttrOr("value", "")

	// Getting car price
	priceStr := doc.Find("#car_price").AttrOr("value", "")
	priceConv, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		return err
	}
	CI.Price = priceConv * 10_000

	// getting full car name
	fullName := doc.Find(".source-info-con h3 a").Text()
	CI.FullName = translations.Translate(fullName)

	// getting mileage and year
	infoText := doc.Find(".source-info-con p").First().Text()
	if idx := strings.Index(infoText, "／"); idx != -1 {
		CI.Milage = translations.Translate(strings.TrimSpace(infoText[:idx]))

		// Getting year
		rest := infoText[idx+3:] // skipping the first symbol
		if idx2 := strings.Index(rest, "／"); idx2 != -1 {
			CI.Year = strings.TrimSpace(rest[:idx2])
		}
	}

	return nil
}

func getCarSpecInfo(CI *carInfo) error {

	carSpecUrl := fmt.Sprintf("https://cacheapigo.che168.com/CarProduct/GetParam.ashx?specid=%s&callback=configTitle", CI.SpecID)
	specs, err := makeRequest(carSpecUrl, 1)
	if err != nil {
		return err
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
					CI.Drive = translations.Translate(value) // TODO: hardcoded drive types
				}

				// Searching for fuel type
				if strings.Contains(name, "燃料形式") {
					CI.FuelType = translations.Translate(value) // TODO: hardcoded value
				}
			}
		}
	}
	return nil
}

func GetCarInfo(url string) (carInfo, error) {

	carInformation := &carInfo{}

	// getting full name, price, year, mileage
	err := getCarConfig(url, carInformation)
	if err != nil {
		return carInfo{}, err
	}

	if carInformation.SpecID != "" {
		// getting car power, engine size, drive, fuel type
		err = getCarSpecInfo(carInformation)
		if err != nil {
			return carInfo{}, err
		}
		return *carInformation, nil
	}

	return carInfo{}, nil
}
