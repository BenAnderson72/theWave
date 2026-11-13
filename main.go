package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/debug"
	"github.com/shopspring/decimal"
)

const openMeteoAirTempURL = "https://api.open-meteo.com/v1/forecast?latitude=51.54&longitude=-2.62&current=temperature_2m&temperature_unit=celsius&timezone=Europe%2FLondon"

type Root struct {
	Temps []Temp `json:"temps"`
}

type Temp struct {
	Timestamp time.Time `json:"timestamp"`
	// Description string          `json:"desc"`
	Water decimal.Decimal `json:"water"`
	Air   decimal.Decimal `json:"air"`
}

// The function persistJSON encodes a given value as JSON and writes it to a file with the specified filename.
func persistJSON(filename string, v any) {
	var buffer bytes.Buffer
	enc := json.NewEncoder(&buffer)
	enc.SetIndent("", "    ") // pretty print
	enc.Encode(&v)
	fo, _ := os.Create(filename)
	fo.Write(buffer.Bytes())
}

func main() {
	if err := scrapeTemperatureAndPersist("https://www.thewave.com/live-updates/"); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func unmarshalJSON(filename string, v any) {
	jsonFile, _ := os.Open(filename)
	byteValue, _ := io.ReadAll(jsonFile)

	json.Unmarshal(byteValue, &v)
}

func scrapeTemperatureAndPersist(waterUrl string) error {
	d, err := scrape(waterUrl)
	if err != nil {
		return err
	}

	// File will be created if it doesn't exist already
	fn := "./temperature.json"

	var dx Root
	unmarshalJSON(fn, &dx)

	dx.Temps = append(dx.Temps, d)

	persistJSON(fn, dx)

	return nil
}

func scrape(waterUrl string) (Temp, error) {

	air, err := scrapeAirTemp()
	if err != nil {
		return Temp{}, err
	}

	water, err := scrapeWaterTemp2(waterUrl)
	if err != nil {
		return Temp{}, err
	}

	// loc, _ := time.LoadLocation("Europe/London")
	return Temp{
		Timestamp: time.Now(), //TODO: UK Time - BST
		// Description: out[desc],
		Air:   air,
		Water: water,
	}, nil
}

// func scrapeWaterTemp(waterUrl string) decimal.Decimal {

// 	c := colly.NewCollector(colly.Debugger(&debug.LogDebugger{}))

// 	// do different stuff if reading local file
// 	if strings.HasPrefix(waterUrl, "file") {
// 		t := &http.Transport{}
// 		t.RegisterProtocol("file", http.NewFileTransport(http.Dir(".")))
// 		c.WithTransport(t)
// 	}

// 	strTemp := "0"

// 	// c.OnRequest(func(r *colly.Request) {
// 	// 	fmt.Println("Visiting: ", r.URL)
// 	// })

// 	// c.OnResponse(func(r *colly.Response) { //get body
// 	// 	fmt.Println("Responding: ", string(r.Body))
// 	// })

// 	// c.OnError(func(r *colly.Response, err error) {
// 	// 	fmt.Println("Request URL: ", r.Request.URL, " failed with response: ", r, "\nError: ", err)
// 	// })

// 	// https://benjamincongdon.me/blog/2018/03/01/Scraping-the-Web-in-Golang-with-Colly-and-Goquery/

// 	c.OnHTML("div.flex", func(d *colly.HTMLElement) {

// 		i := 0

// 		if d.Attr("class") == "flex space-x-1 font-normal" {

// 			d.ForEach("p.text-sm", func(_ int, p *colly.HTMLElement) {

// 				if i == 1 {
// 					strTemp = strings.ReplaceAll(p.Text, "°C", "")
// 					strTemp = strings.TrimSpace(strTemp)
// 				}

// 				i++
// 				// return (i > 1)

// 			})

// 		}
// 	})

// 	c.Visit(waterUrl)

// 	tW, _ := decimal.NewFromString(strTemp)

// 	return tW
// }

// The Wave messed around with where the water temp was published
func scrapeWaterTemp2(waterUrl string) (decimal.Decimal, error) {

	c := colly.NewCollector(colly.Debugger(&debug.LogDebugger{}))

	// do different stuff if reading local file
	if strings.HasPrefix(waterUrl, "file") {
		t := &http.Transport{}
		t.RegisterProtocol("file", http.NewFileTransport(http.Dir(".")))
		c.WithTransport(t)
	}

	strTemp := "0"

	// bg-robins-egg-500 text-bluewood rounded-3xl flex items-center justify-between gap-4 py-2 px-4 max-w-96

	c.OnHTML("div.bg-robins-egg-500", func(d *colly.HTMLElement) {

		i := 0

		if d.Attr("class") == "bg-robins-egg-500 text-bluewood rounded-3xl flex items-center justify-between gap-4 py-2 px-4 max-w-96" {

			d.ForEachWithBreak("p.text-sm", func(_ int, p *colly.HTMLElement) bool {

				if i == 2 {
					strTemp = strings.ReplaceAll(p.Text, "°C", "")
					strTemp = strings.ReplaceAll(strTemp, "Water:", "")
					strTemp = strings.TrimSpace(strTemp)
					return false
				}

				i++
				return true

				// return (i > 1)

			})

		}
	})

	c.Visit(waterUrl)

	// tW, _ := decimal.NewFromString(strTemp)

	return decimal.NewFromString(strTemp)
}

func scrapeAirTemp() (decimal.Decimal, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	return fetchAirTemp(client, openMeteoAirTempURL)
}

func fetchAirTemp(client *http.Client, endpoint string) (decimal.Decimal, error) {
	resp, err := client.Get(endpoint)
	if err != nil {
		return decimal.Zero, fmt.Errorf("request Open-Meteo air temperature: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return decimal.Zero, fmt.Errorf("Open-Meteo returned HTTP %s", resp.Status)
	}

	var result struct {
		Current struct {
			Temperature2m *float64 `json:"temperature_2m"`
		} `json:"current"`
	}
	if err := json.NewDecoder(io.LimitReader(resp.Body, 1<<20)).Decode(&result); err != nil {
		return decimal.Zero, fmt.Errorf("decode Open-Meteo air temperature response: %w", err)
	}
	if result.Current.Temperature2m == nil {
		return decimal.Zero, fmt.Errorf("Open-Meteo response has no current temperature")
	}

	return decimal.NewFromFloat(*result.Current.Temperature2m), nil
}

func addAirTempHistoric() {

	var dx Root
	fn := "./temperature.json"
	unmarshalJSON(fn, &dx)

	for i := range dx.Temps {

		t := &dx.Temps[i]
		if t.Air.Equal(decimal.Zero) {
			t.Air = scrapeAirTempHistoric(t.Timestamp)
		}
	}

	persistJSON(fn, dx)

}

func getAirTempHistoricId(t time.Time) string {
	format := "2006-01-02 15:04:05"
	startDt, _ := time.Parse(format, "2024-05-01 00:00:00")

	if t.Month() == time.June {
		startDt, _ = time.Parse(format, "2024-06-01 00:00:00")
	}

	diff := t.Sub(startDt).Hours()

	id := int(math.Floor(diff / 6))

	return fmt.Sprintf("#ws_%d", id)

}

func scrapeAirTempHistoric(t time.Time) decimal.Decimal {

	// work out the #id
	id := getAirTempHistoricId(t)

	url := "file://./historicTempsMay.html"

	if t.Month() == time.June {
		url = "file://./historicTempsJune.html"
	}

	c := colly.NewCollector(colly.Debugger(&debug.LogDebugger{}))

	// do different stuff if reading local file
	if strings.HasPrefix(url, "file") {
		t := &http.Transport{}
		t.RegisterProtocol("file", http.NewFileTransport(http.Dir(".")))
		c.WithTransport(t)
	}

	// c.OnRequest(func(r *colly.Request) {
	// 	fmt.Println("Visiting: ", r.URL)
	// })

	// c.OnResponse(func(r *colly.Response) { //get body
	// 	fmt.Println("Responding: ", string(r.Body))
	// })

	// c.OnError(func(r *colly.Response, err error) {
	// 	fmt.Println("Request URL: ", r.Request.URL, " failed with response: ", r, "\nError: ", err)
	// })

	strTempLo := "0"
	strTempHi := "0"

	c.OnHTML(id, func(d *colly.HTMLElement) {

		fmt.Println(d.DOM.Children().Find(".tempLow").Text())

		// fmt.Println(qoquerySelection.Find(" span").Children().Text())

		d.ForEach("div", func(_ int, p *colly.HTMLElement) {

			// "tempLow low"
			// "temp low"
			if p.Attr("class") == "tempLow low" {
				strTempLo = strings.ReplaceAll(p.Text, "Lo:", "")
				strTempLo = strings.TrimSpace(strTempLo)
			}

			if p.Attr("class") == "temp low" {
				strTempHi = strings.ReplaceAll(p.Text, "Hi:", "")
				strTempHi = strings.TrimSpace(strTempHi)
			}

		})

	})

	c.Visit(url)

	tL, _ := decimal.NewFromString(strTempLo)
	tH, _ := decimal.NewFromString(strTempHi)

	// tS := decimal.Sum(tL, tH)
	// return tS.Div(decimal.NewFromInt(2))

	return decimal.Avg(tL, tH)

}
