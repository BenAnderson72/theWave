package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/debug"
	"github.com/shopspring/decimal"
)

// const cWaterUrl string = "https://www.thewave.com/live-updates/"
// const cAirUrl string = "https://weather.com/en-GB/weather/today/l/51.54,-2.62"

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
	scrapeTemperatureAndPersist("https://www.thewave.com/live-updates/", "https://weather.com/en-GB/weather/today/l/51.54,-2.62")
}

func unmarshalJSON(filename string, v any) {
	jsonFile, _ := os.Open(filename)
	byteValue, _ := io.ReadAll(jsonFile)

	json.Unmarshal(byteValue, &v)
}

func scrapeTemperatureAndPersist(waterUrl string, airUrl string) {
	d := scrape(waterUrl, airUrl)

	// File will be created if it doesn't exist already
	fn := "./temperature.json"

	var dx Root
	unmarshalJSON(fn, &dx)

	dx.Temps = append(dx.Temps, d)

	persistJSON(fn, dx)

}

func scrape(waterUrl string, airUrl string) Temp {

	// loc, _ := time.LoadLocation("Europe/London")
	return Temp{
		Timestamp: time.Now(), //TODO: UK Time - BST
		// Description: out[desc],
		Air:   scrapeAirTemp(airUrl),
		Water: scrapeWaterTemp(waterUrl),
	}
}

func scrapeWaterTemp(waterUrl string) decimal.Decimal {

	c := colly.NewCollector(colly.Debugger(&debug.LogDebugger{}))

	// do different stuff if reading local file
	if strings.HasPrefix(waterUrl, "file") {
		t := &http.Transport{}
		t.RegisterProtocol("file", http.NewFileTransport(http.Dir(".")))
		c.WithTransport(t)
	}

	strTemp := "0"

	// c.OnRequest(func(r *colly.Request) {
	// 	fmt.Println("Visiting: ", r.URL)
	// })

	// c.OnResponse(func(r *colly.Response) { //get body
	// 	fmt.Println("Responding: ", string(r.Body))
	// })

	// c.OnError(func(r *colly.Response, err error) {
	// 	fmt.Println("Request URL: ", r.Request.URL, " failed with response: ", r, "\nError: ", err)
	// })

	// https://benjamincongdon.me/blog/2018/03/01/Scraping-the-Web-in-Golang-with-Colly-and-Goquery/

	c.OnHTML("div.flex", func(d *colly.HTMLElement) {

		i := 0

		if d.Attr("class") == "flex space-x-1 font-normal" {

			d.ForEach("p.text-sm", func(_ int, p *colly.HTMLElement) {

				if i == 1 {
					strTemp = strings.ReplaceAll(p.Text, "°C", "")
					strTemp = strings.TrimSpace(strTemp)
				}

				i++
				// return (i > 1)

			})

		}
	})

	c.Visit(waterUrl)

	tW, _ := decimal.NewFromString(strTemp)

	return tW
}

func scrapeAirTemp(url string) decimal.Decimal {

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

	strTemp := "0"

	c.OnHTML(".HourlyWeatherCard--TableWrapper--1OobO", func(d *colly.HTMLElement) {

		d.ForEachWithBreak("div", func(_ int, p *colly.HTMLElement) bool {
			// we only wnat the 0th item
			strTemp = strings.ReplaceAll(p.Text, "°", "")
			strTemp = strings.TrimSpace(strTemp)
			return false
		})

	})

	c.Visit(url)

	tA, _ := decimal.NewFromString(strTemp)

	return tA

}
