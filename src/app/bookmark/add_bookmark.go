package bookmark

import (
	"encoding/json"

	"appengine"
	"appengine/datastore"
)

func AddBookmark(cxt appengine.Context, ident string, body []byte, ch chan bool) {
	pBookmarkEntry := new(BookmarkEntry)
	pBookmarkEntry.Ident = ident
	if e := json.Unmarshal(body, pBookmarkEntry); e == nil {
		if _, e := datastore.Put(cxt, datastore.NewIncompleteKey(cxt, "BookmarkEntry", nil), pBookmarkEntry); e == nil {
			ch <- true
		} else {
			ch <- false
		}
	} else {
		ch <- false
	}
}
