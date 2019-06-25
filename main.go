package main

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	fmt.Printf(searchYoutube("  S  "))
}
func searchYoutube(query string) string {
	// fmt.Println(query)
	query = strings.TrimSpace(query)

	query = strings.Replace(query, " ", "+", -1)
	if len(query) < 1 {
		return "EMPTYQUERYSTRING"
	}
	// fmt.Println(query)
	var urlStart = "https://www.youtube.com/results?search_query="
	var finalURL = urlStart + query

	doc, err := goquery.NewDocument(finalURL)
	if err != nil {
		panic(err)
	}

	first := doc.Find(".item-section")
	a := first.Find("a")
	band, ok := a.Attr("href")
	if ok {
		return "https://www.youtube.com" + band
	}
	return "NOTFOUNDSTRING"

}
