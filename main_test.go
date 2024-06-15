package main

import (
	"regexp"
	"testing"

	"github.com/shopspring/decimal"
)

func Test_scrapeWater(t *testing.T) {

	act := scrapeWaterTemp("https://www.thewave.com/live-updates/")

	matched, _ := regexp.MatchString(`^.\d{1,2}(\.\d)*$`, act.String())

	if !matched {
		t.Errorf("expected decimal Water temp, got |%s|", act)
	}

	// matched, _ = regexp.MatchString(`^.\d{1,2}(\.\d)*$`, act.Air.String())

	// if !matched {
	// 	t.Errorf("expected decimal Air temp, got |%s|", act.Air)
	// }

}

func Test_scrapeWaterStatic(t *testing.T) {

	act := scrapeWaterTemp("file://./staticWater.html")

	tW, _ := decimal.NewFromString("19.3")
	// tA, _ := decimal.NewFromString("19")
	// desc := "Broken Clouds &"

	if !decimal.Decimal.Equal(act, tW) {
		t.Errorf("expected %s, got %s", tW, act)
	}

	// if !decimal.Decimal.Equal(act.Air, tA) {
	// 	t.Errorf("expected %s, got %s", tA, act.Air)
	// }

	// if act.Description != desc {
	// 	t.Errorf("expected %s, got %s", desc, act)
	// }

}

func Test_scrapeStatic(t *testing.T) {

	act := scrape("file://./staticWater.html", "file://./staticAir.html")

	tW, _ := decimal.NewFromString("19.3")
	tA, _ := decimal.NewFromString("13")
	// desc := "Broken Clouds &"

	if !decimal.Decimal.Equal(act.Water, tW) {
		t.Errorf("expected %s, got %s", tW, act)
	}

	if !decimal.Decimal.Equal(act.Air, tA) {
		t.Errorf("expected %s, got %s", tA, act.Air)
	}

	// if act.Description != desc {
	// 	t.Errorf("expected %s, got %s", desc, act)
	// }

}

func Test_scrapeAir(t *testing.T) {

	act := scrapeAirTemp("https://weather.com/en-GB/weather/today/l/51.54,-2.62")

	matched, _ := regexp.MatchString(`^.\d{1,2}(\.\d)*$`, act.String())

	if !matched {
		t.Errorf("expected decimal Water temp, got |%s|", act)
	}

}

func Test_scrapeAirTempStatic(t *testing.T) {

	act := scrapeAirTemp("file://./staticAir.html")

	exp, _ := decimal.NewFromString("13")
	// tA, _ := decimal.NewFromString("19")
	// desc := "Broken Clouds &"

	if !decimal.Decimal.Equal(act, exp) {
		t.Errorf("expected %s, got %s", exp, act)
	}

	// if !decimal.Decimal.Equal(act.Air, tA) {
	// 	t.Errorf("expected %s, got %s", tA, act.Air)
	// }

	// if act.Description != desc {
	// 	t.Errorf("expected %s, got %s", desc, act)
	// }

}

func Test_scrapeTemperatureAndPersist(t *testing.T) {
	scrapeTemperatureAndPersist("https://www.thewave.com/live-updates/", "https://weather.com/en-GB/weather/today/l/51.54,-2.62")
}

func Test_scrape(t *testing.T) {

	act := scrape("https://www.thewave.com/live-updates/", "https://weather.com/en-GB/weather/today/l/51.54,-2.62")

	matched, _ := regexp.MatchString(`^.\d{1,2}(\.\d)*$`, act.Water.String())

	if !matched {
		t.Errorf("expected decimal Water temp, got |%s|", act.Water)
	}

	matched, _ = regexp.MatchString(`^.\d{1,2}(\.\d)*$`, act.Air.String())

	if !matched {
		t.Errorf("expected decimal Water temp, got |%s|", act.Air)
	}

}
