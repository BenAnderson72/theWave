package main

import (
	"regexp"
	"testing"
)

func Test_ScrapeFull(t *testing.T) {

	act := scrapeTemperatureFull("https://www.thewave.com/live-updates/")

	matched, _ := regexp.MatchString(`^.\d{1,2}\.\d°C$`, act)

	if !matched {
		t.Errorf("expected something ending with °C, got |%s|", act)
	}

}

func Test_ScrapeTrim(t *testing.T) {

	act := scrapeTemperatureTrim("https://www.thewave.com/live-updates/")

	matched, _ := regexp.MatchString(`^.\d{1,2}\.\d°C$`, act)

	if matched {
		t.Errorf("expected something ending without °C, got |%s|", act)
	}

}

func Test_ScrapeStaticFull(t *testing.T) {
	exp := "19.3°C"
	act := scrapeTemperatureFull("file://./static.html")

	if act != exp {
		t.Errorf("expected %s, got %s", exp, act)
	}
}

func Test_ScrapeStaticTrim(t *testing.T) {
	exp := "19.3"
	act := scrapeTemperatureTrim("file://./static.html")

	if act != exp {
		t.Errorf("expected %s, got %s", exp, act)
	}
}

func Test_ScrapeTemperatureAndPersist(t *testing.T) {
	ScrapeTemperatureAndPersist("https://www.thewave.com/live-updates/")
}
