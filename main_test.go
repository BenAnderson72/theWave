package main

import (
	"regexp"
	"testing"
	"time"

	"github.com/shopspring/decimal"
)

func Test_scrapeWater(t *testing.T) {

	act := scrapeWaterTemp("https://www.thewave.com/live-updates/")

	matched, _ := regexp.MatchString(`^.\d{1,2}(\.\d)*$`, act.String())

	if !matched {
		t.Errorf("expected decimal Water temp, got |%s|", act)
	}

}

func Test_scrapeWater2(t *testing.T) {

	act := scrapeWaterTemp2("https://www.thewave.com")

	matched, _ := regexp.MatchString(`^.\d{1,2}(\.\d)*$`, act.String())

	if !matched {
		t.Errorf("expected decimal Water temp, got |%s|", act)
	}

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

func Test_scrapeWater2Static(t *testing.T) {

	act := scrapeWaterTemp2("file://./staticWater2.html")

	exp := decimal.NewFromInt(20)

	if !decimal.Decimal.Equal(act, exp) {
		t.Errorf("expected %s, got %s", exp, act)
	}

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

func Test_scrape2(t *testing.T) {

	act := scrape("https://www.thewave.com", "https://weather.com/en-GB/weather/today/l/51.54,-2.62")

	matched, _ := regexp.MatchString(`^.\d{1,2}(\.\d)*$`, act.Water.String())

	if !matched {
		t.Errorf("expected decimal Water temp, got |%s|", act.Water)
	}

	matched, _ = regexp.MatchString(`^.\d{1,2}(\.\d)*$`, act.Air.String())

	if !matched {
		t.Errorf("expected decimal Water temp, got |%s|", act.Air)
	}

}

func Test_getAirTempHistoricIdMay(t *testing.T) {

	format := "2006-01-02 15:04:05"
	tm, _ := time.Parse(format, "2024-05-23 00:05:19")

	exp := "#ws_88"
	act := getAirTempHistoricId(tm)

	if act != exp {
		t.Errorf("expected %s, got %s", exp, act)
	}

	// ####

	tm, _ = time.Parse(format, "2024-05-22 21:56:04")

	exp = "#ws_87"
	act = getAirTempHistoricId(tm)

	if act != exp {
		t.Errorf("expected %s, got %s", exp, act)
	}

	// ####

	tm, _ = time.Parse(format, "2024-05-25 18:01:38")

	exp = "#ws_99"
	act = getAirTempHistoricId(tm)

	if act != exp {
		t.Errorf("expected %s, got %s", exp, act)
	}

}

func Test_getAirTempHistoricIdJune(t *testing.T) {

	format := "2006-01-02 15:04:05"
	tm, _ := time.Parse(format, "2024-06-23 00:05:19")

	exp := "#ws_88"
	act := getAirTempHistoricId(tm)

	if act != exp {
		t.Errorf("expected %s, got %s", exp, act)
	}

	// ####

	tm, _ = time.Parse(format, "2024-06-22 21:56:04")

	exp = "#ws_87"
	act = getAirTempHistoricId(tm)

	if act != exp {
		t.Errorf("expected %s, got %s", exp, act)
	}

	// ####

	tm, _ = time.Parse(format, "2024-06-25 18:01:38")

	exp = "#ws_99"
	act = getAirTempHistoricId(tm)

	if act != exp {
		t.Errorf("expected %s, got %s", exp, act)
	}

}

func Test_scrapeAirTempHistoric1(t *testing.T) {

	format := "2006-01-02 15:04:05"
	tm, _ := time.Parse(format, "2024-05-23 00:05:19")

	exp := decimal.NewFromFloat(11.5)
	act := scrapeAirTempHistoric(tm)

	if !decimal.Decimal.Equal(act, exp) {
		t.Errorf("expected |%s|, got |%s|", exp, act)
	}
}

func Test_scrapeAirTempHistoric2(t *testing.T) {

	// "2024-06-14T15:02:32.469839948Z"

	format := "2006-01-02 15:04:05"
	tm, _ := time.Parse(format, "2024-06-15 15:02:32")

	exp := decimal.NewFromFloat(11.5)
	act := scrapeAirTempHistoric(tm)

	if !decimal.Decimal.Equal(act, exp) {
		t.Errorf("expected |%s|, got |%s|", exp, act)
	}
}

func Test_addAirTempHistoric(t *testing.T) {

	addAirTempHistoric()
}
