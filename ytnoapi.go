/*
Package youtubeapiless implements searching functions for video/music resources
*/
package youtubeapiless

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func extractAHrefFromURL(url string) (string, error) {
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

			// println(traversal.Attr[0].Key, traversal.Attr[0].Val) // print the href

			if strings.Contains(traversal.Attr[0].Val, "/watch?v=") { // check we've found the correct <a> tag
				return traversal.Attr[0].Val, nil
			}
			err := errors.New("Error getting href: unknown reason")
			return "", err

		}
		return band, err
	}
	err2 := errors.New("Error getting href: unknown reason")
	return "", err2

}

//GetVideoData ... get data from video
//url: url string
//YoutubeDataJSONStruct: struct of youtube JSON
func GetVideoData(url string) (YoutubeDataJSONStruct, error) {
	// fmt.Printf("URL : ", url)
	var finalYoutubeJSON YoutubeDataJSONStruct
	resp, err := http.Get(url)
	if err != nil {
		httpERR := errors.New("HTTP Error: " + err.Error())
		return finalYoutubeJSON, httpERR
	}
	defer resp.Body.Close()
	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		ioutilReadError := errors.New("ioutilReadError: " + err.Error())
		return finalYoutubeJSON, ioutilReadError
	}
	bodyString := string(respData)                                        // read the response into a string
	endIndex := strings.Index(bodyString, ";ytplayer.load = function()")  // find the end of the JSON
	startIndex := strings.Index(bodyString, "{};ytplayer.config = ") + 21 // find the start of the JSON

	if startIndex <= -1 || endIndex <= -1 { // not a youtube watch page ?
		invalidPageError := errors.New("invalidPageError: this page has no data associated with it ")
		return finalYoutubeJSON, invalidPageError
	}

	JSONstring := bodyString[startIndex:endIndex] // extract the JSON as a string
	var YoutubeJSON youtubeDataRawJSONStruct

	json.Unmarshal([]byte(JSONstring), &YoutubeJSON) //parse JSON
	var playerResponse PlayerResponseJSONStruct
	json.Unmarshal([]byte(YoutubeJSON.Args.PlayerResponse), &playerResponse)

	finalYoutubeJSON.Attrs = YoutubeJSON.Attrs
	finalYoutubeJSON.Assets = YoutubeJSON.Assets
	// args
	finalYoutubeJSON.Args.Cver = YoutubeJSON.Args.Cver
	finalYoutubeJSON.Args.Watermark = YoutubeJSON.Args.Watermark
	finalYoutubeJSON.Args.Author = YoutubeJSON.Args.Author
	finalYoutubeJSON.Args.C = YoutubeJSON.Args.C
	finalYoutubeJSON.Args.Title = YoutubeJSON.Args.Title
	finalYoutubeJSON.Args.PlayerResponse = playerResponse
	finalYoutubeJSON.Args.InnertubeAPIVersion = YoutubeJSON.Args.InnertubeAPIVersion
	finalYoutubeJSON.Args.LoaderURL = YoutubeJSON.Args.LoaderURL
	finalYoutubeJSON.Args.FmtList = YoutubeJSON.Args.FmtList
	finalYoutubeJSON.Args.Hl = YoutubeJSON.Args.Hl
	finalYoutubeJSON.Args.Ucid = YoutubeJSON.Args.Ucid
	finalYoutubeJSON.Args.CsiPageType = YoutubeJSON.Args.CsiPageType
	finalYoutubeJSON.Args.HostLanguage = YoutubeJSON.Args.HostLanguage
	finalYoutubeJSON.Args.AdaptiveFmts = YoutubeJSON.Args.AdaptiveFmts
	finalYoutubeJSON.Args.InnertubeAPIKey = YoutubeJSON.Args.InnertubeAPIKey
	finalYoutubeJSON.Args.Cr = YoutubeJSON.Args.Cr
	finalYoutubeJSON.Args.Ssl = YoutubeJSON.Args.Ssl
	finalYoutubeJSON.Args.Timestamp = YoutubeJSON.Args.Timestamp
	finalYoutubeJSON.Args.InnertubeContextClientVersion = YoutubeJSON.Args.InnertubeContextClientVersion
	finalYoutubeJSON.Args.Fexp = YoutubeJSON.Args.Fexp
	finalYoutubeJSON.Args.VssHost = YoutubeJSON.Args.VssHost
	finalYoutubeJSON.Args.URLEncodedFmtStreamMap = YoutubeJSON.Args.URLEncodedFmtStreamMap
	finalYoutubeJSON.Args.EnabledEngageTypes = YoutubeJSON.Args.EnabledEngageTypes
	finalYoutubeJSON.Args.VideoID = YoutubeJSON.Args.VideoID
	finalYoutubeJSON.Args.ShowContentThumbnail = YoutubeJSON.Args.ShowContentThumbnail
	finalYoutubeJSON.Args.LengthSeconds = YoutubeJSON.Args.LengthSeconds
	finalYoutubeJSON.Args.Enablejsapi = YoutubeJSON.Args.Enablejsapi
	finalYoutubeJSON.Args.AccountPlaybackToken = YoutubeJSON.Args.AccountPlaybackToken
	finalYoutubeJSON.Args.Enablecsi = YoutubeJSON.Args.Enablecsi
	finalYoutubeJSON.Args.GapiHintParams = YoutubeJSON.Args.GapiHintParams
	finalYoutubeJSON.Args.Fflags = YoutubeJSON.Args.Fflags
	finalYoutubeJSON.Sts = YoutubeJSON.Sts
	// https://stackoverflow.com/questions/21268000/unmarshaling-nested-json-objects-in-golang
	// reflectFinalYoutubeJSON := reflect.ValueOf(finalYoutubeJSON)

	return finalYoutubeJSON, nil

}

//YoutubeFirstResult ... function to return first youtube result for given query string
//query: the search string
//string: the returned youtube URL
func YoutubeFirstResult(query string) (string, error) {
	// fmt.Println(query)
	query = strings.TrimSpace(query) // remove the whitespace

	query = strings.Replace(query, " ", "+", -1) //replace spaces with +
	if len(query) < 1 {
		err := errors.New("Error: empty query string")
		return "", err
	}
	query = url.QueryEscape(query)
	println(query)
	// fmt.Println(query)
	var urlStart = "https://www.youtube.com/results?search_query="
	var queryURL = urlStart + query
	var videoURL, err = extractAHrefFromURL(queryURL)
	if err != nil {
		panic(err)
	}
	for strings.Contains(videoURL, "googleadservices") ||
		!strings.Contains(videoURL, "/watch?v=") {
		println(videoURL)
		videoURL, err = extractAHrefFromURL(queryURL)
		if err != nil {
			panic(err)
		}
	}
	return "https://www.youtube.com" + videoURL, nil

}

//PlayerResponseJSONStruct ... struct to hold PlayerResponse JSON
type PlayerResponseJSONStruct struct {
	PlayabilityStatus struct {
		Status          string `json:"status"`
		PlayableInEmbed bool/* replaced from bool */ `json:"playableInEmbed"`
	} `json:"playabilityStatus"`
	StreamingData struct {
		ExpiresInSeconds string `json:"expiresInSeconds"` // seconds until player expires // counts down from 6hours
	} `json:"streamingData"`
	PlaybackTracking struct { // these urls seem to interact with youtube internal resources
		VideostatsPlaybackURL struct {
			BaseURL string `json:"baseUrl"`
		} `json:"videostatsPlaybackUrl"`
		VideostatsDelayplayURL struct {
			BaseURL string `json:"baseUrl"`
		} `json:"videostatsDelayplayUrl"`
		VideostatsWatchtimeURL struct {
			BaseURL string `json:"baseUrl"`
		} `json:"videostatsWatchtimeUrl"`
		PtrackingURL struct {
			BaseURL string `json:"baseUrl"`
		} `json:"ptrackingUrl"`
		QoeURL struct {
			BaseURL string `json:"baseUrl"`
		} `json:"qoeUrl"`
		SetAwesomeURL struct {
			BaseURL                 string `json:"baseUrl"`
			ElapsedMediaTimeSeconds int    `json:"elapsedMediaTimeSeconds"`
		} `json:"setAwesomeUrl"`
		AtrURL struct {
			BaseURL                 string `json:"baseUrl"`
			ElapsedMediaTimeSeconds int    `json:"elapsedMediaTimeSeconds"`
		} `json:"atrUrl"`
	} `json:"playbackTracking"`
	VideoDetails struct {
		VideoID          string                                               `json:"videoId"`       // video ID
		Title            string                                               `json:"title"`         // title of video
		LengthSeconds    string                                               `json:"lengthSeconds"` // length of video in seconds
		ChannelID        string                                               `json:"channelId"`     // id of authors channel
		IsOwnerViewing   bool/* replaced from bool */ `json:"isOwnerViewing"` //
		ShortDescription string                                               `json:"shortDescription"` // the initial description seen on load of page
		IsCrawlable      bool/* replaced from bool */ `json:"isCrawlable"`
		Thumbnail        struct { // thumbnails
			Thumbnails []struct { // ascending size of thumbnails array
				URL    string `json:"url"`    // url
				Width  int    `json:"width"`  // width
				Height int    `json:"height"` // height
			} `json:"thumbnails"`
		} `json:"thumbnail"`
		UseCipher         bool/* replaced from bool */ `json:"useCipher"`
		AverageRating     float64 `json:"averageRating"` // ratio of likes to dislikes // 0 - 5
		AllowRatings      bool/* replaced from bool */ `json:"allowRatings"`
		ViewCount         string `json:"viewCount"` // view count integer
		Author            string `json:"author"`    // name of channel
		IsPrivate         bool/* replaced from bool */ `json:"isPrivate"`
		IsUnpluggedCorpus bool/* replaced from bool */ `json:"isUnpluggedCorpus"`
		IsLiveContent     bool/* replaced from bool */ `json:"isLiveContent"`
	} `json:"videoDetails"`
	Annotations []struct { // array of annotation data for the video
		PlayerAnnotationsUrlsRenderer struct {
			InvideoURL         string `json:"invideoUrl"` // the url of the xml file
			LoadPolicy         string `json:"loadPolicy"` // variable of when to load this annotation // "always"
			AllowInPlaceSwitch bool/* replaced from bool */ `json:"allowInPlaceSwitch"`
		} `json:"playerAnnotationsUrlsRenderer"`
	} `json:"annotations"`
	PlayerConfig struct {
		AudioConfig struct {
			LoudnessDb           float64 `json:"loudnessDb"`
			PerceptualLoudnessDb float64 `json:"perceptualLoudnessDb"`
		} `json:"audioConfig"`
		StreamSelectionConfig struct {
			MaxBitrate string `json:"maxBitrate"`
		} `json:"streamSelectionConfig"`
		MediaCommonConfig struct {
			DynamicReadaheadConfig struct {
				MaxReadAheadMediaTimeMs int `json:"maxReadAheadMediaTimeMs"`
				MinReadAheadMediaTimeMs int `json:"minReadAheadMediaTimeMs"`
				ReadAheadGrowthRateMs   int `json:"readAheadGrowthRateMs"`
			} `json:"dynamicReadaheadConfig"`
		} `json:"mediaCommonConfig"`
	} `json:"playerConfig"`
	Storyboards struct {
		PlayerStoryboardSpecRenderer struct {
			Spec string `json:"spec"` // url that needs permissions
		} `json:"playerStoryboardSpecRenderer"`
	} `json:"storyboards"`
	TrackingParams string `json:"trackingParams"` // "CAMQu2kiEwi8oJfiia7jAhWLFuEKHeI2BMEo-B0="
	Attestation    struct {
		PlayerAttestationRenderer struct {
			Challenge string `json:"challenge"` // "a=4&b=xuefcODfdrmRciBf-FVDifvfDhk&c=1562889681&d=1&e=Nq5LMGtBmis&c3a=21&c1a=1&c6a=1&hh=mzJdpHwqpqKntapm46BHsoHnpFwE3K4t6C5GJ33U-as"
		} `json:"playerAttestationRenderer"`
	} `json:"attestation"`
	Messages []struct {
		MealbarPromoRenderer struct {
			MessageTexts []struct {
				Runs []struct {
					Text string `json:"text"`
				} `json:"runs"`
			} `json:"messageTexts"`
			ActionButton struct {
				ButtonRenderer struct {
					Style string `json:"style"`
					Size  string `json:"size"`
					Text  struct {
						Runs []struct {
							Text string `json:"text"`
						} `json:"runs"`
					} `json:"text"`
					NavigationEndpoint struct {
						ClickTrackingParams string `json:"clickTrackingParams"`
						URLEndpoint         struct {
							URL    string `json:"url"`
							Target string `json:"target"`
						} `json:"urlEndpoint"`
					} `json:"navigationEndpoint"`
					VisualElement struct {
						UIType int `json:"uiType"`
					} `json:"visualElement"`
					TrackingParams string `json:"trackingParams"`
				} `json:"buttonRenderer"`
			} `json:"actionButton"`
			DismissButton struct {
				ButtonRenderer struct {
					Style string `json:"style"`
					Size  string `json:"size"`
					Text  struct {
						Runs []struct {
							Text string `json:"text"`
						} `json:"runs"`
					} `json:"text"`
					ServiceEndpoint struct {
						ClickTrackingParams string `json:"clickTrackingParams"`
						FeedbackEndpoint    struct {
							FeedbackToken string `json:"feedbackToken"`
							UIActions     struct {
								HideEnclosingContainer bool /* replaced from bool */ `json:"hideEnclosingContainer"`
							} `json:"uiActions"`
						} `json:"feedbackEndpoint"`
					} `json:"serviceEndpoint"`
					VisualElement struct {
						UIType int `json:"uiType"`
					} `json:"visualElement"`
					TrackingParams string `json:"trackingParams"`
				} `json:"buttonRenderer"`
			} `json:"dismissButton"`
			TriggerCondition string `json:"triggerCondition"`
			Style            string `json:"style"`
			VisualElement    struct {
				UIType      int `json:"uiType"`
				YoutubeData struct {
					YpcAnnotation struct {
						ItemTypeYpc    string `json:"itemTypeYpc"`
						ItemExternalID string `json:"itemExternalId"`
					} `json:"ypcAnnotation"`
					PromotionData struct {
						PromotionID      string `json:"promotionId"`
						CrossPromoRanked bool/* replaced from bool */ `json:"crossPromoRanked"`
					} `json:"promotionData"`
				} `json:"youtubeData"`
			} `json:"visualElement"`
			TrackingParams      string `json:"trackingParams"`
			ImpressionEndpoints []struct {
				ClickTrackingParams string `json:"clickTrackingParams"`
				FeedbackEndpoint    struct {
					FeedbackToken string `json:"feedbackToken"`
					UIActions     struct {
						HideEnclosingContainer bool /* replaced from bool */ `json:"hideEnclosingContainer"`
					} `json:"uiActions"`
				} `json:"feedbackEndpoint"`
			} `json:"impressionEndpoints"`
			IsVisible    bool/* replaced from bool */ `json:"isVisible"`
			MessageTitle struct {
				Runs []struct {
					Text string `json:"text"`
				} `json:"runs"`
			} `json:"messageTitle"`
		} `json:"mealbarPromoRenderer"`
	} `json:"messages"`
	AdPlacements []struct {
		AdPlacementRenderer struct {
			Config struct {
				AdPlacementConfig struct {
					Kind         string `json:"kind"`
					AdTimeOffset struct {
						OffsetStartMilliseconds string `json:"offsetStartMilliseconds"`
						OffsetEndMilliseconds   string `json:"offsetEndMilliseconds"`
					} `json:"adTimeOffset"`
					HideCueRangeMarker bool/* replaced from bool */ `json:"hideCueRangeMarker"`
				} `json:"adPlacementConfig"`
			} `json:"config"`
			Renderer struct {
				AdBreakServiceRenderer struct {
					GetAdBreakURL string `json:"getAdBreakUrl"`
				} `json:"adBreakServiceRenderer"`
			} `json:"renderer"`
			TrackingParams string `json:"trackingParams"`
		} `json:"adPlacementRenderer"`
	} `json:"adPlacements"`
	AdSafetyReason struct {
		ApmUserPreference struct {
		} `json:"apmUserPreference"`
	} `json:"adSafetyReason"`
}

// raw json struct without playerResponse parsed
type youtubeDataRawJSONStruct struct {
	Attrs struct {
		ID string `json:"id"`
	} `json:"attrs"`
	Assets struct {
		Js  string `json:"js"`
		CSS string `json:"css"`
	} `json:"assets"`
	Args struct {
		Cver                          string `json:"cver"`                  // version number // "1.20190711"
		Watermark                     string `json:"watermark"`             // 2 values
		Author                        string `json:"author"`                // name of authors channel
		C                             string `json:"c"`                     // content version ? / "web"
		Title                         string `json:"title"`                 // title of video
		PlayerResponse                string `json:"player_response"`       // SEEMS TO BE JSON
		InnertubeAPIVersion           string `json:"innertube_api_version"` // general version number // "v1"
		LoaderURL                     string `json:"loaderUrl"`             // ACTUAL URL FULL HTTPS
		FmtList                       string `json:"fmt_list"`
		Hl                            string `json:"hl"`                               // location ? host language ? // "en_GB"
		Ucid                          string `json:"ucid"`                             // user channel id // "UCtWuB1D_E3mcyYThA9iKggQ"
		CsiPageType                   string `json:"csi_page_type"`                    // type of page // "watch" // IF NOT VIDEO THIS IS AN EMPTY STRING
		HostLanguage                  string `json:"host_language"`                    // duplicate field ? // different format for reason ? // "en-GB"
		AdaptiveFmts                  string `json:"adaptive_fmts"`                    // long string // unusable ?
		InnertubeAPIKey               string `json:"innertube_api_key"`                // "AIzaSyAO_FJ2SlqU8Q4STEHLGCilw_Y9_11qcW8"
		Cr                            string `json:"cr"`                               // content reigon ? // "GB"
		Ssl                           string `json:"ssl"`                              //empty
		Timestamp                     string `json:"timestamp"`                        // an int meaning timestamp ? // "1562887566"
		InnertubeContextClientVersion string `json:"innertube_context_client_version"` // duplicate field // `json:"cver"` // "1.20190711"
		Fexp                          string `json:"fexp"`                             // many comma sperated integers // eg: "23735271"
		VssHost                       string `json:"vss_host"`                         // resource host ? // "s.youtube.com"
		URLEncodedFmtStreamMap        string `json:"url_encoded_fmt_stream_map"`       // long string
		EnabledEngageTypes            string `json:"enabled_engage_types"`             // comma seperated 1 digit numbers // "3,6,4,5,17,1"
		VideoID                       string `json:"video_id"`                         //video id
		ShowContentThumbnail          bool/* replaced from bool */ `json:"show_content_thumbnail"`
		LengthSeconds                 string `json:"length_seconds"`         // length of video in seconds
		Enablejsapi                   string `json:"enablejsapi"`            // 1 digit number // "1"
		AccountPlaybackToken          string `json:"account_playback_token"` // long hex?
		Enablecsi                     string `json:"enablecsi"`              /// 1 digit number // "1"
		GapiHintParams                string `json:"gapi_hint_params"`       // some kind of google api url
		Fflags                        string `json:"fflags"`                 // flags for function ?
	} `json:"args"`
	Sts int `json:"sts"`
}

//YoutubeDataJSONStruct ... struct holding the FULL Youtube JSON
type YoutubeDataJSONStruct struct {
	Attrs struct {
		ID string `json:"id"`
	} `json:"attrs"`
	Assets struct {
		Js  string `json:"js"`
		CSS string `json:"css"`
	} `json:"assets"`
	Args struct {
		Cver                          string                   `json:"cver"`      // version number // "1.20190711"
		Watermark                     string                   `json:"watermark"` // 2 values
		Author                        string                   `json:"author"`    // name of authors channel
		C                             string                   `json:"c"`         // content version ? / "web"
		Title                         string                   `json:"title"`     // title of video
		PlayerResponse                PlayerResponseJSONStruct `json:"player_response"`
		InnertubeAPIVersion           string                   `json:"innertube_api_version"` // general version number // "v1"
		LoaderURL                     string                   `json:"loaderUrl"`             // ACTUAL URL FULL HTTPS
		FmtList                       string                   `json:"fmt_list"`
		Hl                            string                   `json:"hl"`                               // location ? host language ? // "en_GB"
		Ucid                          string                   `json:"ucid"`                             // user channel id // "UCtWuB1D_E3mcyYThA9iKggQ"
		CsiPageType                   string                   `json:"csi_page_type"`                    // type of page // "watch" // IF NOT VIDEO THIS IS AN EMPTY STRING
		HostLanguage                  string                   `json:"host_language"`                    // duplicate field ? // different format for reason ? // "en-GB"
		AdaptiveFmts                  string                   `json:"adaptive_fmts"`                    // long string // unusable ?
		InnertubeAPIKey               string                   `json:"innertube_api_key"`                // "AIzaSyAO_FJ2SlqU8Q4STEHLGCilw_Y9_11qcW8"
		Cr                            string                   `json:"cr"`                               // content reigon ? // "GB"
		Ssl                           string                   `json:"ssl"`                              //empty
		Timestamp                     string                   `json:"timestamp"`                        // an int meaning timestamp ? // "1562887566"
		InnertubeContextClientVersion string                   `json:"innertube_context_client_version"` // duplicate field // `json:"cver"` // "1.20190711"
		Fexp                          string                   `json:"fexp"`                             // many comma sperated integers // eg: "23735271"
		VssHost                       string                   `json:"vss_host"`                         // resource host ? // "s.youtube.com"
		URLEncodedFmtStreamMap        string                   `json:"url_encoded_fmt_stream_map"`       // long string
		EnabledEngageTypes            string                   `json:"enabled_engage_types"`             // comma seperated 1 digit numbers // "3,6,4,5,17,1"
		VideoID                       string                   `json:"video_id"`                         //video id
		ShowContentThumbnail          bool/* replaced from bool */ `json:"show_content_thumbnail"`
		LengthSeconds                 string `json:"length_seconds"`         // length of video in seconds
		Enablejsapi                   string `json:"enablejsapi"`            // 1 digit number // "1"
		AccountPlaybackToken          string `json:"account_playback_token"` // long hex?
		Enablecsi                     string `json:"enablecsi"`              /// 1 digit number // "1"
		GapiHintParams                string `json:"gapi_hint_params"`       // some kind of google api url
		Fflags                        string `json:"fflags"`                 // flags for function ?
	} `json:"args"`
	Sts int `json:"sts"`
}
