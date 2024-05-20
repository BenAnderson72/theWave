package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/debug"
)

type Datax struct {
	Temps []Data `json:"temps"`
}

type Data struct {
	Timestamp time.Time `json:"timestamp"`
	Temp      float64   `json:"temp"`
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

func UnmarshalJSON(filename string, v any) {
	jsonFile, _ := os.Open(filename)
	byteValue, _ := io.ReadAll(jsonFile)

	json.Unmarshal(byteValue, &v)
}

func ScrapeTemperatureAndPersist(url string) {
	out := scrapeTemperatureTrim(url)

	fout, _ := strconv.ParseFloat(out, 32)

	fn := "./temperature.json"

	var dx Datax
	UnmarshalJSON(fn, &dx)

	d := Data{
		Timestamp: time.Now(),
		Temp:      fout,
	}

	dx.Temps = append(dx.Temps, d)

	PersistJSON(fn, dx)

	// temps := dx.Temps

}

func scrapeTemperatureTrim(url string) string {

	out := scrapeTemperatureFull(url)
	return strings.ReplaceAll(out, "°C", "")
}

func scrapeTemperatureFull(url string) string {

	c := colly.NewCollector(colly.Debugger(&debug.LogDebugger{}))

	// do different stuff if reading local file
	if strings.HasPrefix(url, "file") {
		t := &http.Transport{}
		t.RegisterProtocol("file", http.NewFileTransport(http.Dir(".")))
		c.WithTransport(t)
	}

	out := "Nothing"

	// c.OnRequest(func(r *colly.Request) {
	// 	fmt.Println("Visiting: ", r.URL)
	// })

	// c.OnResponse(func(r *colly.Response) { //get body
	// 	fmt.Println("Responding: ", string(r.Body))
	// })

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL: ", r.Request.URL, " failed with response: ", r, "\nError: ", err)
	})

	// https://benjamincongdon.me/blog/2018/03/01/Scraping-the-Web-in-Golang-with-Colly-and-Goquery/

	i := 0

	c.OnHTML("div.flex", func(d *colly.HTMLElement) {
		if d.Attr("class") == "flex space-x-1 font-normal" {
			// i := d.DOM.Children().Length()
			// fmt.Println(i)

			// x := d.DOM.Children()

			d.ForEach("p.text-sm", func(_ int, p *colly.HTMLElement) {

				if i == 1 {
					out = strings.TrimSpace(p.Text)
				}
				i++

			})

		}
	})

	c.Visit(url)

	return out
}
