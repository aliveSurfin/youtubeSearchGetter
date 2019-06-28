/*
Package search implements searching functions for video/music resources
*/
package search

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func extractAHrefFromURL(url string) string {
	doc, err := goquery.NewDocument(url) // get youtube page from url
	if err != nil {
		panic(err)
	}
	items := doc.Find(".item-section") // first item section
	a := items.Find("a")               // find the a section that contains the href
	band, ok := a.Attr("href")         // grab the href
	if ok {
		if strings.Contains(band, "user") { // check if it's a channel
			traversal := doc.Find(".yt-lockup-content").Get(1) // go to the second item as the channel takes up the first one

			traversal = traversal.FirstChild // do down to title

			traversal = traversal.FirstChild // go down to a

			// println(traversal.Attr[0].Key, traversal.Attr[0].Val)

			if strings.Contains(traversal.Attr[0].Val, "/watch?v=") { // check we've found the correct <a> tag
				return traversal.Attr[0].Val
			}
			return ("ErrorGettingHref")

		}
		return band
	}
	return ("ErrorGettingHref")

}

//GetVideoData ... get data from video
func GetVideoData(url string) YoutubeDataRawJSON {
	// fmt.Printf("URL : ", url)
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	respData, err := ioutil.ReadAll(resp.Body)
	bodyString := string(respData)                                        // read the response into a string
	endIndex := strings.Index(bodyString, ";ytplayer.load = function()")  // find the end of the JSON
	startIndex := strings.Index(bodyString, "{};ytplayer.config = ") + 21 // find the start of the JSON
	JSONstring := bodyString[startIndex:endIndex]                         // extract the JSON as a string
	var YoutubeJSON YoutubeDataRawJSON

	json.Unmarshal([]byte(JSONstring), &YoutubeJSON) //parse JSON
	// jsonData, err := json.MarshalIndent(YoutubeJSON, "", "    ")
	// fmt.Println(string(jsonData))
	return YoutubeJSON

}

type YoutubeDataRawJSON struct {
	Attrs struct {
		ID string `json:"id"`
	} `json:"attrs"`
	Assets struct {
		Js  string `json:"js"`
		CSS string `json:"css"`
	} `json:"assets"`
	Args struct {
		Cver                          string `json:"cver"`
		Watermark                     string `json:"watermark"`
		Author                        string `json:"author"`
		C                             string `json:"c"`
		Title                         string `json:"title"`
		PlayerResponse                string `json:"player_response"`
		InnertubeAPIVersion           string `json:"innertube_api_version"`
		LoaderURL                     string `json:"loaderUrl"`
		FmtList                       string `json:"fmt_list"`
		Hl                            string `json:"hl"`
		Ucid                          string `json:"ucid"`
		CsiPageType                   string `json:"csi_page_type"`
		HostLanguage                  string `json:"host_language"`
		AdaptiveFmts                  string `json:"adaptive_fmts"`
		InnertubeAPIKey               string `json:"innertube_api_key"`
		Cr                            string `json:"cr"`
		Ssl                           string `json:"ssl"`
		Timestamp                     string `json:"timestamp"`
		InnertubeContextClientVersion string `json:"innertube_context_client_version"`
		Fexp                          string `json:"fexp"`
		VssHost                       string `json:"vss_host"`
		URLEncodedFmtStreamMap        string `json:"url_encoded_fmt_stream_map"`
		EnabledEngageTypes            string `json:"enabled_engage_types"`
		VideoID                       string `json:"video_id"`
		ShowContentThumbnail          bool   `json:"show_content_thumbnail"`
		LengthSeconds                 string `json:"length_seconds"`
		Enablejsapi                   string `json:"enablejsapi"`
		AccountPlaybackToken          string `json:"account_playback_token"`
		Enablecsi                     string `json:"enablecsi"`
		GapiHintParams                string `json:"gapi_hint_params"`
		Fflags                        string `json:"fflags"`
	} `json:"args"`
	Sts int `json:"sts"`
}

//YoutubeFirstResult ... function to return first youtube result for given query string
//query: the search string
//string: the returned youtube URL
func YoutubeFirstResult(query string) string {
	// fmt.Println(query)
	query = strings.TrimSpace(query) // remove the whitespace

	query = strings.Replace(query, " ", "+", -1) //replace spaces with +
	if len(query) < 1 {
		return "EMPTYQUERYSTRING"
	}
	query = url.QueryEscape(query)
	println(query)
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
