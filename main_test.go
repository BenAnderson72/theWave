package main

import (
	"testing"
	"time"

	"github.com/shopspring/decimal"
)

// func Test_scrapeWater(t *testing.T) {

// 	act := scrapeWaterTemp("https://www.thewave.com/live-updates/")

// 	matched, _ := regexp.MatchString(`^.\d{1,2}(\.\d)*$`, act.String())

// 	if !matched {
// 		t.Errorf("expected decimal Water temp, got |%s|", act)
// 	}

// }

func Test_scrapeWater2(t *testing.T) {

	act, err := scrapeWaterTemp2("https://www.thewave.com")

	if err != nil {
		t.Error(err)
	}

	if act.GreaterThan(high_water_temp) || act.LessThan(low_water_temp) {
		t.Errorf("Unlikely temp %s", act)
	}

	// matched, _ := regexp.MatchString(`^.\d{1,2}(\.\d)*$`, act.String())

	// if !matched {
	// 	t.Errorf("expected decimal Water temp, got |%s|", act)
	// }

}

// func Test_scrapeWaterStatic(t *testing.T) {

// 	act := scrapeWaterTemp("file://./staticWater.html")

// 	tW, _ := decimal.NewFromString("19.3")
// 	// tA, _ := decimal.NewFromString("19")
// 	// desc := "Broken Clouds &"

// 	if !decimal.Decimal.Equal(act, tW) {
// 		t.Errorf("expected %s, got %s", tW, act)
// 	}

// 	// if !decimal.Decimal.Equal(act.Air, tA) {
// 	// 	t.Errorf("expected %s, got %s", tA, act.Air)
// 	// }

// 	// if act.Description != desc {
// 	// 	t.Errorf("expected %s, got %s", desc, act)
// 	// }

// }

func Test_scrapeWater2Static(t *testing.T) {

	act, err := scrapeWaterTemp2("file://./staticWater2.html")

	if err != nil {
		t.Error(err)
	}

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

var high_air_temp, _ = decimal.NewFromString("40")
var low_air_temp, _ = decimal.NewFromString("-20")
var high_water_temp, _ = decimal.NewFromString("40")
var low_water_temp, _ = decimal.NewFromString("-20")

func Test_scrapeAir(t *testing.T) {

	act, err := scrapeAirTemp("https://weather.com/en-GB/weather/today/l/51.54,-2.62")

	if err != nil {
		t.Error(err)
	}

	if act.GreaterThan(high_air_temp) || act.LessThan(low_air_temp) {
		t.Errorf("Unlikely temp %s", act)
	}
}

func Test_scrapeAirTempStatic(t *testing.T) {

	act, err := scrapeAirTemp("file://./staticAir.html")
	if err != nil {
		t.Error(err)
	}

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

func Test_scrapeAirTempStatic2(t *testing.T) {

	act, err := scrapeAirTemp("file://./staticAir2.html")

	if err != nil {
		t.Error(err)
	}

	exp, _ := decimal.NewFromString("9")

	if !decimal.Decimal.Equal(act, exp) {
		t.Errorf("expected %s, got %s", exp, act)
	}
}

func Test_scrapeTemperatureAndPersist(t *testing.T) {
	scrapeTemperatureAndPersist("https://www.thewave.com/live-updates/", "https://weather.com/en-GB/weather/today/l/51.54,-2.62")
}

func Test_scrape(t *testing.T) {

	act := scrape("https://www.thewave.com/live-updates/", "https://weather.com/en-GB/weather/today/l/51.54,-2.62")

	if act.Air.GreaterThan(high_air_temp) || act.Air.LessThan(low_air_temp) {
		t.Errorf("Unlikely temp %s", act)
	}

	if act.Water.GreaterThan(high_water_temp) || act.Water.LessThan(low_water_temp) {
		t.Errorf("Unlikely temp %s", act)
	}

}

func Test_scrape2(t *testing.T) {

	act := scrape("https://www.thewave.com", "https://weather.com/en-GB/weather/today/l/51.54,-2.62")

	if act.Air.GreaterThan(high_air_temp) || act.Air.LessThan(low_air_temp) {
		t.Errorf("Unlikely temp %s", act)
	}

	if act.Water.GreaterThan(high_water_temp) || act.Water.LessThan(low_water_temp) {
		t.Errorf("Unlikely temp %s", act)
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
