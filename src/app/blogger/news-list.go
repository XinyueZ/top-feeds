package blogger

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time" 

	"appengine"
	"appengine/urlfetch"
)

const ( 
	//DEV
	KEY = "AIzaSyDpy_nLxxUU3n7CLG5ozCqMEzWnMrxvTSA"
	//LIVE 
	//KEY = "AIzaSyCeZuL95RHFZG3Q6THzPZYgncP_Ewd2x5Q"
	API_FIRST_PAGE = "https://www.googleapis.com/blogger/v3/blogs/%s/posts?key=%s"
	API_NEXT_PAGE = "https://www.googleapis.com/blogger/v3/blogs/%s/posts?key=%s&pageToken=%s"
	BLOGGER_DATE_FORMAT = "2006-01-02T15:04:05"
)

func NewNewsList() (p *NewsList) {
	p = new(NewsList)
	return
}

type NewsList struct {
	Next 			string   `json:"nextPageToken"`
	NewsEntries []NewsEntry  `json:"items"` 
}
 
type NewsEntry struct { 
	Title        string `json:"title"` 
	Updated      string `json:"updated"` 
	Url          string `json:"url"` 
	UrlMobile    string
}

 

//Create a news-list and return a json-feeds to client through channel.
func (self *NewsList) Create(cxt appengine.Context, bloggerId string, from string, chJsonStr chan *string, chFrom chan *string) {
	client := urlfetch.Client(cxt)
	api := ""
	if from == "0" {
		api = fmt.Sprintf(API_FIRST_PAGE, bloggerId, KEY )
	} else {
		api = fmt.Sprintf(API_NEXT_PAGE, bloggerId, KEY, from)
	}
	if r, e := http.NewRequest("GET", api, nil); e == nil {
		if resp, e := client.Do(r); e == nil {
			if resp != nil {
				defer resp.Body.Close()
			}

			pNewsList := new(NewsList)
			if bytes, e := ioutil.ReadAll(resp.Body); e == nil {
				if e := json.Unmarshal(bytes, pNewsList); e == nil {
					s := "[" //Start making a json result.
					for _, v := range pNewsList.NewsEntries {
						v.Title = strings.Replace(v.Title, "\"", "'", -1)
						v.Title = strings.Replace(v.Title, "\\", ",", -1)
						v.Title = strings.Replace(v.Title, "%", "％", -1)
						v.Title = strings.Replace(v.Title, "\n", "", -1)
						v.Title = strings.Replace(v.Title, "\t", "", -1)  
						v.UrlMobile = fmt.Sprintf("%s??m=1", v.Url)
  
						lastBin := strings.LastIndex( v.Updated, "-" )
						v.Updated =  v.Updated[0:lastBin]
						t, _ := time.Parse(BLOGGER_DATE_FORMAT, v.Updated)
						v.Updated = t.String()
						pubDate := t.Unix()
						json := fmt.Sprintf(`{"title" : "%s", "desc" : "%s", "url" : "%s", "url_mobile" : "%s",  "pubDate" : %d },`, v.Title, "", v.Url, v.UrlMobile, pubDate)
						//cxt.Infof("Json: %s", json)
						s += json
					}
					length := len(s)
					if length > 2 {
						s = s[:length-1] //Remove last ","
					}
					s += "]" //Stop making json
					chJsonStr <- &s
					chFrom <- &pNewsList.Next
				} else {
					chJsonStr <- nil	
					chFrom <- nil
					cxt.Errorf("Error but still going: %v", e)
				}
			} else {
				chJsonStr <- nil
				chFrom <- nil
				panic(e)
			}
		} else {
			chJsonStr <- nil
			chFrom <- nil
			cxt.Errorf("Error but still going: %v", e)
		}
	} else {
		chJsonStr <- nil
		chFrom <- nil
		panic(e)
	}
}
