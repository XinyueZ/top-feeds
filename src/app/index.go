package topfeeds

import (
	"csdn"
	"techug"
	"oschina"

	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"bookmark"
)

import "appengine"

type Error string

func (e Error) Error() string {
	return string(e)
}

func init() {
	http.HandleFunc("/topfeeds", handleTopFeeds)

	http.HandleFunc("/bookmark", handleAddBookmark)
	http.HandleFunc("/bookmarkList", handleBookmarkList)
	http.HandleFunc("/removeBookmark", handleRemoveBookmark)
}

func handleTopFeeds(w http.ResponseWriter, r *http.Request) {
	cxt := appengine.NewContext(r)

	//Error-handling anyway.
	defer func() {
		if err := recover(); err != nil {
			cxt.Errorf("handleTopFeeds: %v", err)
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(w, `{"status":%d}`, 300)
		}
	}()

	args := r.URL.Query()

	typ := 0 //Which type, 0: oschina, 1: csdn
	if len(args["type"]) > 0 {
		t, _ := strconv.Atoi(args["type"][0])
		typ = t
	}

	page := 0
	if len(args["page"]) > 0 {
		p, _ := strconv.Atoi(args["page"][0]) //Which page, if typ is csdn(1), then ignore this.
		page = p
	}

	site := ""
	siteMobile := ""
	res := ""
	switch typ {
	case 1:
		//Ask csdn:
		chCsdn := make(chan *string)
		go csdn.NewNewsList().Create(cxt, chCsdn)
		res = *(<-chCsdn)

		site = "http://www.csdn.net"
		siteMobile = "http://m.csdn.net"
	case 2:
		//Ask techug:
		chTechug := make(chan *string)
		go techug.NewNewsList().Create(cxt, chTechug)
		res = *(<-chTechug)

		site = "http://www.techug.com"
		siteMobile = "http://www.techug.com"
	default:
		//Ask news-list of www.oschina.net
		chOsc := make(chan *string)
		go oschina.NewNewsList().Create(cxt, page, chOsc)
		res = *(<-chOsc)

		site = "http://www.oschina.net"
		siteMobile = "http://m.oschina.net"
	}

	//Output result of all news-list of all sites.
	s := fmt.Sprintf(`{"status":%d, "page_index" : %d, "site" : "%s", "site_mobile":"%s", "result":%s}`, 200, page, site, siteMobile, res)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "public, max-age=900")
	fmt.Fprintf(w, s)
}

func handleAddBookmark(w http.ResponseWriter, r *http.Request) {
	args := r.URL.Query()
	ident := args["ident"][0]

	var s string
	if bytes, e := ioutil.ReadAll(r.Body); e == nil {
		cxt := appengine.NewContext(r)
		ch := make(chan bool)
		go bookmark.AddBookmark(cxt, ident, bytes, ch)
		if <-ch {
			s = fmt.Sprintf(`{"status":%d}`, 200)
		} else {
			s = fmt.Sprintf(`{"status":%d}`, 300)
		}
	} else {
		s = fmt.Sprintf(`{"status":%d}`, 300)
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, s)
}

func handleBookmarkList(w http.ResponseWriter, r *http.Request) {
	args := r.URL.Query()
	ident := args["ident"][0]

	cxt := appengine.NewContext(r)
	var s string
	ch := make(chan *string)
	go bookmark.GetBookmarkList(cxt, ident, ch)
	p := <-ch
	if p != nil {
		s = fmt.Sprintf(`{"status":%d, result:%s}`, 200, *p)
	} else {
		s = fmt.Sprintf(`{"status":%d}`, 300)
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, s)
}

func handleRemoveBookmark(w http.ResponseWriter, r *http.Request) {
	args := r.URL.Query()
	ident := args["ident"][0]

	var s string
	if bytes, e := ioutil.ReadAll(r.Body); e == nil {
		cxt := appengine.NewContext(r)
		ch := make(chan bool)
		go bookmark.DelBookmark(cxt, ident, bytes, ch)
		if <-ch {
			s = fmt.Sprintf(`{"status":%d}`, 200)
		} else {
			s = fmt.Sprintf(`{"status":%d}`, 300)
		}
	} else {
		s = fmt.Sprintf(`{"status":%d}`, 300)
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, s)
}
