package translations

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mashinki/logging"
	"net/http"
)

var logger *logging.FLogger

func init() {
	var err error
	logger, err = logging.NewFLogger("log.txt")
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize logger: %v", err))
	}
}

func translateByLibreTranslate(text string, sourceLang string, targetLang string) (string, error) {

	// forming the request data
	data := map[string]interface{}{
		"q":      text,
		"source": sourceLang,
		"target": targetLang,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("error while forming JSON: %v", err)
	}

	// Sending the request
	resp, err := http.Post("http://localhost:5000/translate", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("error while sending request: %v", err)
	}
	defer resp.Body.Close()

	// Reading the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error while reading response: %v", err)
	}

	// Getting the translated text
	var result struct {
		TranslatedText string `json:"translatedText"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("error while parsing the response: %v", err)
	}

	return result.TranslatedText, nil
}

func Translate(chineseText string) string {
	if chineseText == "" {
		return chineseText
	}

	// Ch to En
	englishText, err := translateByLibreTranslate(chineseText, "zh", "en")
	if err != nil {
		logger.LogErrorF("Error while translating to English\n "+
			"Chinese:: %v\n"+
			"error: %v", chineseText, err)
		return chineseText
	}

	// En to Rus
	russianText, err := translateByLibreTranslate(englishText, "en", "ru")
	if err != nil {
		logger.LogErrorF("Error while translating ro Russian\n "+
			"English: %v\n"+
			"error: %v", englishText, err)
		return englishText
	}

	return russianText
}
