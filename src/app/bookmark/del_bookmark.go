package bookmark

import (
	"encoding/json"

	"appengine"
	"appengine/datastore"
)

func DelBookmark(cxt appengine.Context, body []byte, ch chan bool) {
	pBookmarkEntry := new(BookmarkEntry)
	if e := json.Unmarshal(body, pBookmarkEntry); e == nil {
		q := datastore.NewQuery("BookmarkEntry").Filter("Title=", pBookmarkEntry.Title).Filter("Description=", pBookmarkEntry.Description).Filter("PubDate=", pBookmarkEntry.PubDate).Filter("Url=", pBookmarkEntry.Url).Filter("UrlMobile=", pBookmarkEntry.UrlMobile)
		pBookmarkEntries := make([]BookmarkEntry, 0)
		if keys, e := q.GetAll(cxt, &pBookmarkEntries); e == nil {
			if datastore.DeleteMulti(cxt, keys) == nil {
				ch <- true
			} else {
				ch <- false
			}
		} else {
			ch <- false
		}
	} else {
		ch <- false
	}
}
