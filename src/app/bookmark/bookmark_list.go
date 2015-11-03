package bookmark

import ( 
	"fmt"
	"net/http" 
	"io/ioutil" 
	"strings" 
	
	"appengine" 
	"appengine/urlfetch"
)

 
 
func GetBookmarkList(cxt appengine.Context, ident string, ch chan *string) { 
	where := fmt.Sprintf(WHERE, ident)
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
				s := string(bytes) 
				s = strings.Replace(s, "\"results\":", "\"result\":",  -1)
				ss := s[1:len(s)-1]
				ch <- &ss
			} else {
				cxt.Errorf("Read bookmark-list error: %v.", err)
				ch <- nil
			}
		} else {
			cxt.Errorf("Call bookmark-list error: %v.", err)
			ch <- nil
		}
	} else {
		cxt.Errorf("Build bookmark-list request error: %v.", err)
		ch <- nil
	}
	return
}
