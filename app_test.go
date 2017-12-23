package main

import (
    "net/http"
    "net/http/httptest"
	"testing"
	"./core"
)

func makeRequest(url string, configFile string) (int, string) {
	config := core.GetConfigFromJSON(configFile)
	addHandlers()

	req, _ := http.NewRequest("GET", url, nil)
	w := httptest.NewRecorder()
	core.MainHandler(config).ServeHTTP(w, req)

	result := w.Body.String()
	return w.Code, result
}

func makeRequestForSingleLanguage(url string) (int, string)  {
	return makeRequest(url, "./config_test_single.json")
}

func makeRequestForMultiLanguage(url string) (int, string)  {
	return makeRequest(url, "./config.json")
}

func TestExistingPage(t *testing.T) {
	code, result := makeRequestForSingleLanguage("/anasayfa")
	if code != http.StatusOK {
		t.Errorf("HttpStatus is not 200, %s", code)
	}
	if result != "Anasayfa" {
		t.Errorf("Unexpected body: %v", result)
	}
}

func TestNonExistentPage(t *testing.T) {
	code, _ := makeRequestForSingleLanguage("/home")
	if code != http.StatusNotFound {
		t.Errorf("HttpStatus is not 404, %s", code)
	}
}

func TestRedirection(t *testing.T) {
	code, _ := makeRequestForMultiLanguage("/home")
	if code != 301 {
		t.Errorf("HttpStatus is not 301, %s", code)
	}
}

func TestMultiDefaultLanguageHome(t *testing.T) {
	code, result := makeRequestForMultiLanguage("/")
	if code != http.StatusOK {
		t.Errorf("HttpStatus is not OK, %s", code)
	}
	if result != "Home" {
		t.Errorf("Default Language's home page not here ?")
	}
}

func TestMultiLanguagePage(t *testing.T) {
	code, result := makeRequestForMultiLanguage("/tr/anasayfa")
	if code != http.StatusOK {
		t.Errorf("HttpStatus is not OK, %s", code)
	}
	if result != "Anasayfa" {
		t.Errorf("Anasayfa ?")
	}
}

func TestMultiLanguageNonExistentPage(t *testing.T) {
	code, _ := makeRequestForMultiLanguage("/tr/babasayfa")
	if code != http.StatusNotFound {
		t.Errorf("HttpStatus is not 404 %s", code)
	}
}