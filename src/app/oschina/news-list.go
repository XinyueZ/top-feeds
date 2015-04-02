package oschina

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
	//catalog:1 for all
	//pageIndex: The first page.
	//pageSize:20
	API    = "http://www.oschina.net/action/api/news_list?catalog=1&pageIndex=%d&pageSize=20"
	DETAIL = "http://www.oschina.net/action/api/news_detail?id=%d"
)

func NewNewsList() (p *NewsList) {
	p = new(NewsList)
	return
}

type NewsList struct {
	XMLName     xml.Name    `xml:"oschina"`
	NewsEntries NewsEntries `xml:"newslist"`
}

type NewsEntries struct {
	Array []NewsEntry `xml:"news"`
}

type NewsEntry struct {
	Id           int    `xml:"id"`
	Title        string `xml:"title"`
	CommentCount int    `xml:"commentCount"`
	Author       string `xml:"author"`
	AuthorId     int    `xml:"authorid"`
	PubDate      string `xml:"pubDate"`
	Url          string `xml:"url"` //Might be empty then the news-type should be used to build a url
	UrlMobile    string
	NewsType     NewsType `xml:"newstype"`
	Description  string
}

type NewsType struct {
	Type int `xml:"type"`
}

type NewsDetail struct {
	XMLName  xml.Name `xml:"oschina"`
	NewsBody NewsBody `xml:"news"`
}

type NewsBody struct {
	Content string `xml:"body"`
}

//Create a news-list and return a json-feeds to client through channel.
func (self *NewsList) Create(cxt appengine.Context, page int, chJsonStr chan *string) {
	client := urlfetch.Client(cxt)
	if r, e := http.NewRequest("GET", fmt.Sprintf(API, page), nil); e == nil {
		if resp, e := client.Do(r); e == nil {
			if resp != nil {
				defer resp.Body.Close()
			}

			pNewsList := new(NewsList)
			if bytes, e := ioutil.ReadAll(resp.Body); e == nil {
				if e := xml.Unmarshal(bytes, pNewsList); e == nil {
					s := "[" //Start making a json result.
					for _, v := range pNewsList.NewsEntries.Array {
						v.Title = strings.Replace(v.Title, "\"", "'", -1)
						if v.Url == "" { //A Url might be null then we need change it self associated with its type.
							v.Url = fmt.Sprintf("http://www.oschina.net/news/%d", v.Id)
							v.UrlMobile = fmt.Sprintf("http://m.oschina.net/news/%d", v.Id)
						} else {
							v.UrlMobile = v.Url
						}

						//To fetch details and make description.
						/*
							if r, e := http.NewRequest("GET", fmt.Sprintf(DETAIL, v.Id), nil); e == nil {
								if resp, e := client.Do(r); e == nil {
									if resp != nil {
										defer resp.Body.Close()
									}

									pNewsDetail := new(NewsDetail)
									if bytes, e := ioutil.ReadAll(resp.Body); e == nil {
										//cxt.Infof("Details: %v", string(bytes))
										if e := xml.Unmarshal(bytes, pNewsDetail); e == nil {
											v.Description = pNewsDetail.NewsBody.Content[0:20]
										}
									}
								}
							}
						*/
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
