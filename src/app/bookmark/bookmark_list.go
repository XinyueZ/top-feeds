package bookmark

import (
	"encoding/json"

	"appengine"
	"appengine/datastore"
)

func GetBookmarkList(cxt appengine.Context, ident string, ch chan *string) {
	q := datastore.NewQuery("BookmarkEntry").Filter("Ident=", ident)
	bookmarkEntries := make([]BookmarkEntry, 0)
	if _, e := q.GetAll(cxt, &bookmarkEntries); e == nil {
		json, _ := json.Marshal(&bookmarkEntries)
		s := string(json)
		ch <- &s
	} else {
		ch <- nil
	}
}
