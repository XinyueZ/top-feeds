package bookmark

import (
	"encoding/json"

	"appengine"
	"appengine/datastore"
)

func AddBookmark(cxt appengine.Context, body []byte, ch chan bool) {
	pBookmarkEntry := new(BookmarkEntry)
	cxt.Infof("incom:%s", string(body))
	if e := json.Unmarshal(body, pBookmarkEntry); e == nil {
		cxt.Infof("incom:%s", "Unmarshal")
		if _, e := datastore.Put(cxt, datastore.NewIncompleteKey(cxt, "BookmarkEntry", nil), pBookmarkEntry); e == nil {
			cxt.Infof("incom:%s", "Put")
			ch <- true
		} else {
			ch <- false
		}
	} else {
		ch <- false
	}
}
