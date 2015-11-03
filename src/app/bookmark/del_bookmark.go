package bookmark

import (
	"fmt"
	"encoding/json"
	"net/http" 
	"io/ioutil" 
	
	"appengine" 
	"appengine/urlfetch"
)

type BmobList struct {
	Results       []Bmob `json:"results"` 
}

type Bmob struct {
	ObjectId        string `json:"objectId"` 
}

func DelBookmark(cxt appengine.Context, ident string, body []byte, ch chan bool) {
	pBookmarkEntry := new(BookmarkEntry)
	json.Unmarshal(body, pBookmarkEntry)
	cxt.Infof("%v", pBookmarkEntry)
	where := fmt.Sprintf(DEL, ident, pBookmarkEntry.Title, pBookmarkEntry.Description,pBookmarkEntry.PubDate,pBookmarkEntry.Url,pBookmarkEntry.UrlMobile)
	if req, err := http.NewRequest("GET", URL + "?" + where, nil); err == nil {
		req.Header.Add(DB_HEADER_APP_ID, DB_APP_ID)
		req.Header.Add(DB_HEADER_API_KEY, DB_API_KEY)  
		httpClient := urlfetch.Client(cxt)
		r, err := httpClient.Do(req)
		if r != nil {
			defer r.Body.Close()
		}
		if err == nil {
			if bytes, err := ioutil.ReadAll(r.Body); err == nil { 
				cxt.Infof("GET bookmark-list successfully.")
				pRes := new(BmobList)
				json.Unmarshal(bytes, pRes) 
				cxt.Infof("%v", pRes)
				req, _ := http.NewRequest("DELETE", URL + "/" + pRes.Results[0].ObjectId, nil)
				req.Header.Add(DB_HEADER_APP_ID, DB_APP_ID)
				req.Header.Add(DB_HEADER_API_KEY, DB_API_KEY)  
				httpClient := urlfetch.Client(cxt)
				r, err := httpClient.Do(req)
				if r != nil {
					defer r.Body.Close()
				}
				if err == nil {
					cxt.Infof("DELETE entry successfully.")
					ch <- true
				} else {
					cxt.Errorf("DELETE entry error: %v.", err)
					ch <- false
				}
			} else {
				cxt.Errorf("Read bookmark-list error: %v.", err)
				ch <- false
			}
		} else {
			cxt.Errorf("Call bookmark-list error: %v.", err)
			ch <- false
		}
	} else {
		cxt.Errorf("Build bookmark-list request error: %v.", err)
		ch <- false
	}
	return
}
