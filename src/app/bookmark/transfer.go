package bookmark

import (
	"encoding/json" 
	"net/http"
	"bytes"
	
	"appengine"
	"appengine/datastore"
	"appengine/urlfetch"
)

 
func getBookmarkList(cxt appengine.Context  ) (bookmarkEntries []BookmarkEntry){
	q := datastore.NewQuery("BookmarkEntry")
	bookmarkEntries = make([]BookmarkEntry, 0)
	if _, e := q.GetAll(cxt, &bookmarkEntries); e != nil {
		cxt.Errorf("Get list entries error.")
	} 
	return
}


func Transfer(cxt appengine.Context,   ch chan int) {  	
	bookmarkEntries := getBookmarkList(cxt  )
	for _, v := range bookmarkEntries {
		entryJson, _ := json.Marshal(v) 
		if req, err := http.NewRequest("POST", URL,  bytes.NewBuffer(entryJson)); err == nil {
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
			} else {
				cxt.Errorf("POST entry error: %v.", err) 
			}
		} else {
			cxt.Errorf("Build request error: %v.", err) 
		}
	}
	ch <- 0
}