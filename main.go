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

type Root struct {
	Temps []Temp `json:"temps"`
}

type Temp struct {
	Timestamp   time.Time       `json:"timestamp"`
	Description string          `json:"desc"`
	Water       decimal.Decimal `json:"water"`
	Air         decimal.Decimal `json:"air"`
}

// The function PersistJSON encodes a given value as JSON and writes it to a file with the specified filename.
func PersistJSON(filename string, v any) {
	var buffer bytes.Buffer
	enc := json.NewEncoder(&buffer)
	enc.SetIndent("", "    ") // pretty print
	enc.Encode(&v)
	fo, _ := os.Create(filename)
	fo.Write(buffer.Bytes())
}

func main() {
	ScrapeTemperatureAndPersist("https://www.thewave.com/live-updates/")
}

func UnmarshalJSON(filename string, v any) {
	jsonFile, _ := os.Open(filename)
	byteValue, _ := io.ReadAll(jsonFile)

	json.Unmarshal(byteValue, &v)
}

func ScrapeTemperatureAndPersist(url string) {
	d := scrapeTemperature(url)

	// File will be created if it doesn't exist already
	fn := "./temperature.json"

	var dx Root
	UnmarshalJSON(fn, &dx)

	dx.Temps = append(dx.Temps, d)

	PersistJSON(fn, dx)

}

func scrapeTemperature(url string) Temp {

	c := colly.NewCollector(colly.Debugger(&debug.LogDebugger{}))

	// do different stuff if reading local file
	if strings.HasPrefix(url, "file") {
		t := &http.Transport{}
		t.RegisterProtocol("file", http.NewFileTransport(http.Dir(".")))
		c.WithTransport(t)
	}

	// out := "Nothing"

	var out [3]string

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

		if d.Attr("class") == "flex space-x-1 font-normal mb-2" {

			d.ForEach("p.text-sm", func(_ int, p *colly.HTMLElement) {

				switch i {
				case 0:
					out[0] = strings.TrimSpace(p.Text)
				case 1:
					out[1] = strings.ReplaceAll(p.Text, "°C", "")
					out[1] = strings.TrimSpace(out[1])
				}

				i++

			})

		}

		i = 0

		if d.Attr("class") == "flex space-x-1 font-normal" {

			d.ForEach("p.text-sm", func(_ int, p *colly.HTMLElement) {

				if i == 1 {
					out[2] = strings.ReplaceAll(p.Text, "°C", "")
					out[2] = strings.TrimSpace(out[2])
				}

				i++
				// return (i > 1)

			})

		}
	})

	c.Visit(url)

	// fmt.Println(out)
	tA, _ := decimal.NewFromString(out[1])
	tW, _ := decimal.NewFromString(out[2])

	// loc, _ := time.LoadLocation("Europe/London")
	return Temp{
		Timestamp:   time.Now(), //TODO: UK Time - BST
		Description: out[0],
		Air:         tA,
		Water:       tW,
	}

}
