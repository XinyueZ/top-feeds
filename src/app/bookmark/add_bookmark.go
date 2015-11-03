package bookmark

import (
	"fmt"
	"encoding/json"
	"net/http"
	"bytes"
	
	"appengine" 
	"appengine/urlfetch"
	 
)

func AddBookmark(cxt appengine.Context, ident string, body []byte, ch chan bool) {
	pBookmarkEntry := new(BookmarkEntry)
	pBookmarkEntry.Ident = ident
	json.Unmarshal(body, pBookmarkEntry);
	entryJson, _ := json.Marshal(*pBookmarkEntry) 
	
	where := fmt.Sprintf(WHERE, ident)
	if req, err := http.NewRequest("POST", URL + "?" + where,  bytes.NewBuffer(entryJson)); err == nil {
		req.Header.Add(DB_HEADER_APP_ID, DB_APP_ID)
		req.Header.Add(DB_HEADER_API_KEY, DB_API_KEY)
		req.Header.Add(CONTENT_TYPE, API_RESTYPE)
		httpClient := urlfetch.Client(cxt)
		r, err := httpClient.Do(req)
		if r != nil {
			defer r.Body.Close()
		}
		if err == nil {
			cxt.Infof("POST entry successfully.")
			ch<-true
		} else {
			cxt.Errorf("POST entry error: %v.", err) 
			ch<-false
		}
	} else {
		cxt.Errorf("Build request error: %v.", err) 
		ch<-false
	}
	 
} 
