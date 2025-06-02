package parser

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	envhandler "mashinki/envHandler"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

// makeRequest makes http request
//
// mode: 0 - just get response
//
// mode: 1 - get response into GBK encoding
func makeRequest(targetUrl string, mode int) (string, error) {
	proxyURL, err := url.Parse(envhandler.GetEnv("PROXY"))
	if err != nil {
		panic("error while patsing proxy url: " + err.Error())
	}

	// initializing http client
	client := &http.Client{
		Transport: &http.Transport{
			Proxy:             http.ProxyURL(proxyURL),
			IdleConnTimeout:   30 * time.Second,
			DisableKeepAlives: false,
		},
		Timeout: 15 * time.Second,
	}

	req, err := http.NewRequest("GET", targetUrl, nil)
	if err != nil {
		return "", fmt.Errorf("error while creating request: %v", err)
	}

	// Setting up headers
	switch mode {
	case 0:
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
		req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
		req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")
		req.Header.Set("Cache-Control", "no-cache")
		req.Header.Set("Connection", "keep-alive")
		req.Header.Set("Pragma", "no-cache")
		req.Header.Set("Referer", "https://www.che168.com")
	case 1:
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
		req.Header.Set("Accept", "*/*")
		req.Header.Set("Referer", "https://www.che168.com")
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error while sending request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("server returned status %d instead of 200 OK", resp.StatusCode)
	}

	reader := transform.NewReader(resp.Body, simplifiedchinese.GBK.NewDecoder())
	body, err := io.ReadAll(reader)
	if err != nil {
		return "", fmt.Errorf("error while reading response: %v", err)
	}
	return string(body), nil

	// switch mode {
	// case 0:
	// 	body, err := io.ReadAll(resp.Body)
	// 	if err != nil {
	// 		return "", fmt.Errorf("error while reading response: %v", err)
	// 	}
	// 	return string(body), nil
	// case 1:
	// 	reader := transform.NewReader(resp.Body, simplifiedchinese.GBK.NewDecoder())
	// 	body, err := io.ReadAll(reader)
	// 	if err != nil {
	// 		return "", fmt.Errorf("error while reading response: %v", err)
	// 	}
	// 	return string(body), nil
	// }
	// return "", nil
}
