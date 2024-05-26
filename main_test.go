package main

import (
	"regexp"
	"testing"

	"github.com/shopspring/decimal"
)

func Test_Scrape(t *testing.T) {

	act := scrapeTemperature("https://www.thewave.com/live-updates/")

	matched, _ := regexp.MatchString(`^.\d{1,2}(\.\d)*$`, act.Water.String())

	if !matched {
		t.Errorf("expected decimal Water temp, got |%s|", act.Water)
	}

	matched, _ = regexp.MatchString(`^.\d{1,2}(\.\d)*$`, act.Air.String())

	if !matched {
		t.Errorf("expected decimal Air temp, got |%s|", act.Air)
	}

}

func Test_ScrapeStatic(t *testing.T) {

	act := scrapeTemperature("file://./static.html")

	tW, _ := decimal.NewFromString("19.3")
	tA, _ := decimal.NewFromString("19")
	desc := "Broken Clouds &"

	if !decimal.Decimal.Equal(act.Water, tW) {
		t.Errorf("expected %s, got %s", tW, act.Water)
	}

	if !decimal.Decimal.Equal(act.Air, tA) {
		t.Errorf("expected %s, got %s", tA, act.Air)
	}

	if act.Description != desc {
		t.Errorf("expected %s, got %s", desc, act)
	}

}

func Test_ScrapeTemperatureAndPersist(t *testing.T) {
	ScrapeTemperatureAndPersist("https://www.thewave.com/live-updates/")
}
