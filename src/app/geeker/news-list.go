package geeker

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"appengine"
	"appengine/urlfetch"

	"github.com/PuerkitoBio/goquery"
	xhtml "golang.org/x/net/html"
)

const (
	API = "http://geek.csdn.net/service/news/get_news_list?type=hackernewsv2_new&size=20&from=%s"
	// `Format` and `Parse` use example-based layouts. Usually
	// you'll use a constant from `time` for these layouts, but
	// you can also supply custom layouts. Layouts must use the
	// reference time `Mon Jan 2 15:04:05 MST 2006` to show the
	// pattern with which to format/parse a given time/string.
	// The example time must be exactly as shown: the year 2006,
	// 15 for the hour, Monday for the day of the week, etc.
	OSC_DATE_FORMAT = "2006-01-02 15:04:05"
)

func NewNewsList() (p *NewsList) {
	p = new(NewsList)
	return
}

type NewsList struct {
	From    string `json:"from"`
	Html    string `json:"html"`
	HasMore bool   `json:"has_more"`
}

func isTitle(attr []xhtml.Attribute) bool {
	for _, element := range attr {
		if element.Key == "class" && element.Val == `\"title\"` {
			return true
		}
	}
	return false
}

func parse_xhtml(cxt appengine.Context, n *xhtml.Node, pValues *[]string) {
	if n.Type == xhtml.ElementNode && n.Data == "a" && isTitle(n.Attr) {
		link := ""
		title := ""
		for _, element := range n.Attr {
			if element.Key == "href" {
				link = element.Val
				link = strings.Replace(link, `\"`, ``, -1)
			}
		}
		c := n.FirstChild
		title = c.Data

		desc := ""
		if link != "" {
			doc, err := goquery.NewDocument(link)
			if err == nil {
				desc = doc.Find("head").Text()
				if desc != "" {
					desc = fmt.Sprintf("%s...", desc[:len(desc)-1])
				}
			}
		}
		
		cxt.Infof("desc: %s", desc)
		
		now, _ := strconv.ParseInt(time.Now().Local().Format("20060102150405"), 10, 64)
		json := fmt.Sprintf(`{"title" : "%s", "desc" : "%s", "url" : "%s", "url_mobile" : "%s",  "pubDate" : %d },`,
			title, desc, link, link, now)
		*pValues = append(*pValues, json)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		parse_xhtml(cxt, c, pValues)
	}
}

//Create a news-list and return a json-feeds to client through channel.
func (self *NewsList) Create(cxt appengine.Context, from string, chJsonStr chan *string) {
	client := urlfetch.Client(cxt)
	if r, e := http.NewRequest("GET", fmt.Sprintf(API, from), nil); e == nil {
		if resp, e := client.Do(r); e == nil {
			if resp != nil {
				defer resp.Body.Close()
			}

			pNewsList := new(NewsList)
			if bt, e := ioutil.ReadAll(resp.Body); e == nil {
				if e := json.Unmarshal(bt, pNewsList); e == nil {
					jsons := []string{}
					source := strings.TrimSpace(pNewsList.Html)
					source = strings.Replace(source, `"`, `\"`, -1)
				 
					pReader := strings.NewReader(source)
					doc, _ := xhtml.Parse(pReader)
					parse_xhtml(cxt, doc, &jsons)

					s := "[" //Start making a json result.
					for _, json := range jsons {
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
