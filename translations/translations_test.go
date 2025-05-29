package translations

import (
	"net/http"
	"testing"
)

// Pinging the libretranslate
func TestLibreTranslateConnection(t *testing.T) {
	resp,err :=http.Get("http://localhost:5000")

	if err != nil{
		t.Fatalf("LibreTranslate service is off, please turn it on\n" +
		"docker start lt-container\n" +
		"or\n" +
		"docker build -t my-libretranslate .\n" +
		"docker run -d -p 5000:5000 --name lt-container my-libretranslate\n")
	}

	if resp.StatusCode != http.StatusOK{
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}
}