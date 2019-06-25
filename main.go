package main

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	fmt.Println(searchYoutube("testing"))
}
func extractAHrefFromURL(url string) string {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		panic(err)
	}
	items := doc.Find(".item-section")
	a := items.Find("a")
	// output, err := goquery.OuterHtml(a)
	// if err != nil {
	// 	panic(err)
	// }
	// println(output)
	band, ok := a.Attr("href")
	if ok {
		return band
	}
	return ("ErrorGettingHref")
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
	var queryURL = urlStart + query
	var videoURL = extractAHrefFromURL(queryURL)
	for strings.Contains(videoURL, "googleadservices") ||
		!strings.Contains(videoURL, "/watch?v=") {
		println(videoURL)
		videoURL = extractAHrefFromURL(queryURL)
	}
	return "https://www.youtube.com" + videoURL

}
