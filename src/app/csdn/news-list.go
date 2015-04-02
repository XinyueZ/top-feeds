package csdn

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"appengine"
	"appengine/urlfetch"
)

const (
	API = "http://www.csdn.net/article/rss_lastnews"
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
						v.UrlMobile = strings.Replace(v.Url, "www", "m", -1)
						json := fmt.Sprintf(`{"title" : "%s", "desc" : "%s", "url" : "%s", "url_mobile" : "%s",  "pubDate" : "%s" },`, v.Title, v.Description, v.Url, v.UrlMobile, v.PubDate)
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
