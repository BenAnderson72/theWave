package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/debug"
)

type Datax struct {
	// Desc   string  `json:"desc"`
	Timestamp time.Time `json:"timestamp"`
	// Spots     []Spot    `json:"spots"`
}

func Scrape(url string) string {

	c := colly.NewCollector(colly.Debugger(&debug.LogDebugger{}))

	// do different stuff if reading local file
	if strings.HasPrefix(url, "file") {
		t := &http.Transport{}
		t.RegisterProtocol("file", http.NewFileTransport(http.Dir(".")))
		c.WithTransport(t)
	}

	out := "Nothing"

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting: ", r.URL)
	})

	c.OnResponse(func(r *colly.Response) { //get body
		// fmt.Println("Responding: ", string(r.Body))
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL: ", r.Request.URL, " failed with response: ", r, "\nError: ", err)
	})

	c.OnHTML("div.flex", func(d *colly.HTMLElement) {
		if d.Attr("class") == "flex space-x-1 font-normal" {
			d.ForEach("p.text-sm", func(_ int, p *colly.HTMLElement) {
				out += p.Text
			})

			// i := d.DOM.Children().Length()
			// fmt.Println(i)

		}
	})

	// c.OnHTML("p.text-sm", func(p *colly.HTMLElement) {
	// 	out += p.Text
	// })

	c.Visit(url)

	// p.text-sm

	// c.OnHTML("head title", func(d *colly.HTMLElement) {

	// 	out = d.Text

	// #forecast-day-1 > div > div.TableDayContainer_dayTableWrapper__347n4.MuiBox-root.mui-style-0 > div.MuiTableContainer-root.ForecastTableRows_container__7O5Tf.mui-style-kge0eu > div.ForecastTableRows_containerInner__bwn2U.MuiBox-root.mui-style-0 > div > table:nth-child(1)
	// d.ForEach("table", func(_ int, t *colly.HTMLElement) {

	// 	//*[@id="forecast-day-1"]/div/div[2]/div[2]/div[1]/div/table[1]/tbody/tr[1]/td[2]

	// 	out = t.Attr("class")

	// 	if t.Attr("class") != "table table-sm table-striped table-inverse table-tide" {

	// 		// mswDays = nil

	// 		d.ForEach("tbody", func(_ int, tb *colly.HTMLElement) {

	// 			tb.ForEach("tr", func(_ int, tr *colly.HTMLElement) {

	// 				if daysIndex > daysToCollect-1 {
	// 					// break
	// 				} else if tr.Attr("class") == "tbody-title" {
	// 				} else if tr.Attr("class") == "background-clear msw-js-tide" {
	// 				} else {
	// 					et, e := strconv.ParseInt(tr.Attr("data-timestamp"), 10, 64)
	// 					if e == nil {
	// 						t1 := time.Unix(et, 0).In(location)

	// 						// for today skip all earlier forecasts
	// 						if t1.Day() == time.Now().Day() &&
	// 							t1.Before(time.Now().Add(-3*time.Hour)) {
	// 							goto skip
	// 						}

	// 						// on change of day create the t0 object

	// 						tr.ForEach("td", func(_ int, td *colly.HTMLElement) {

	// 							if strings.HasPrefix(td.Attr("class"), "table-forecast-rating td-nowrap") {
	// 								class := td.ChildAttrs("li", "class")
	// 								var counter [2]int
	// 								// 0 = inactive 1 = active
	// 								for _, c := range class {
	// 									switch c {
	// 									case "active":
	// 										counter[1]++
	// 									case "inactive":
	// 										counter[0]++
	// 									}
	// 								}
	// 								// sw := Swell{Timestamp: t1, StarsNow: [2]int{counter[1], counter[0]}}
	// 								// swells4day = append(swells4day, sw)

	// 								if !t0.IsZero() && (t1.Day() != t0.Day()) {
	// 									daysCaught++
	// 								}

	// 								t0 = t1

	// 							}
	// 						})

	// 					skip:
	// 					}

	// 				}
	// 			})
	// 		})
	// 	}
	// })
	// })
	// c.Visit(url)

	return out
}
