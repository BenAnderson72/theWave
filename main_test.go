package main

import "testing"

func Test_Scrape(t *testing.T) {
	Scrape("https://www.thewave.com/live-updates/")
}

func Test_ScrapeStatic(t *testing.T) {
	exp := "19.2"
	act := Scrape("file://./static.html")

	if act != exp {
		t.Errorf("expected %s, got %s", exp, act)
	}
}
