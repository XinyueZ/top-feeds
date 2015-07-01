package techug

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"appengine"
	"appengine/urlfetch"
)

const (
	API = "http://www.techug.com/feed"
	// `Format` and `Parse` use example-based layouts. Usually
	// you'll use a constant from `time` for these layouts, but
	// you can also supply custom layouts. Layouts must use the
	// reference time `Mon Jan 2 15:04:05 MST 2006` to show the
	// pattern with which to format/parse a given time/string.
	// The example time must be exactly as shown: the year 2006,
	// 15 for the hour, Monday for the day of the week, etc.
	TECHUG_DATE_FORMAT = "Mon, 02 Jan 2006 15:04:05 +0000"
)

func NewNewsList() (p *NewsList) {
	p = new(NewsList)
	return
}

type NewsList struct {
	XMLName xml.Name `xml:"rss"`
	Channel Channel  `xml:"channel"`
}


type Channel struct {
	NewsEntries []NewsEntry `xml:"item"`
}

type NewsEntry struct {
	Title       string `xml:"title"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
	Url         string `xml:"link"` //Might be empty then the news-type should be used to build a url
	UrlMobile   string
}


//Create a news-list and return a json-feeds to client through channel.
func (self *NewsList) Create(cxt appengine.Context, chJsonStr chan *string) {
	client := urlfetch.Client(cxt)
	if r, e := http.NewRequest("GET", API, nil); e == nil {
		if resp, e := client.Do(r); e == nil {
			if resp != nil {
				defer resp.Body.Close()
			}

			pNewsList := new(NewsList)
			if bytes, e := ioutil.ReadAll(resp.Body); e == nil {
				if e := xml.Unmarshal(bytes, pNewsList); e == nil {
					s := "[" //Start making a json result.
					for _, v := range pNewsList.Channel.NewsEntries {
						v.Title = strings.Replace(v.Title, "\"", "'", -1)
						v.Title = strings.Replace(v.Title, "%", "％", -1)
            v.Title = strings.Replace(v.Title, "\\", ",", -1)
						v.Description = strings.Replace(v.Description, "\"", "'", -1)
						v.Description = strings.Replace(v.Description, "%", "％", -1)
            v.Description = strings.Replace(v.Description, " ", "", -1)
            v.Description = strings.Replace(v.Description, "\n", "", -1)
            v.Description = strings.Replace(v.Description, "\t", "", -1)
						v.UrlMobile = v.Url//strings.Replace(v.Url, "www", "m", -1)

						loc, _ := time.LoadLocation("Asia/Shanghai")
						t, _ := time.ParseInLocation(TECHUG_DATE_FORMAT, v.PubDate, loc)
						v.PubDate = t.String()
						pubDate := t.Unix()
						json := fmt.Sprintf(`{"title" : "%s", "desc" : "%s", "url" : "%s", "url_mobile" : "%s",  "pubDate" : %d },`, v.Title, v.Description, v.Url, v.UrlMobile, pubDate)
						//cxt.Infof("Json: %s", json)
						s += json
					}
					length := len(s)
					if length > 2 {
						s = s[:length-1] //Remove last ","
					}
					s += "]" //Stop making json
					chJsonStr <- &s
				} else {
					chJsonStr <- nil
					cxt.Errorf("Error but still going: %v", e)
				}
			} else {
				chJsonStr <- nil
				panic(e)
			}
		} else {
			chJsonStr <- nil
			cxt.Errorf("Error but still going: %v", e)
		}
	} else {
		chJsonStr <- nil
		panic(e)
	}
}
