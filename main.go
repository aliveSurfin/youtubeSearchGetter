package main

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/robertkrimen/otto"
)

func main() {
	var test = searchYoutube("testing")
	fmt.Println(test)
}
func extractAHrefFromURL(url string) string {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		panic(err)
	}
	items := doc.Find(".item-section")
	a := items.Find("a")
	band, ok := a.Attr("href")
	if ok {
		return band
	}
	return ("ErrorGettingHref")
}
func getVideoData(url string) {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		panic(err)
	}
	vm := otto.New()
	_, err = vm.Run("var window = {};")
	_, err = vm.Run(doc.Text())
	// println(doc.Text())
	// _, err = vm.Run(`var test = document.getElementsByTagName("video")[0]`)
	// _, err = vm.Run(`var test = "not a test"`)
	// window, err := vm.Get("window")
	_, err = vm.Run(`var test = document.getElementsByTagName("video")[0]`)
	value, err := vm.Get("test")
	output, err := value.ToString()
	println(output)

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
