package topfeeds

import (
	"csdn"
	"fmt"
	"net/http"
	"oschina"
	"strconv"
)

import "appengine"

type Error string

func (e Error) Error() string {
	return string(e)
}

func init() {
	http.HandleFunc("/topfeeds", handleTopFeeds)
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
	default:
		//Ask news-list of www.oschina.net
		chOsc := make(chan *string)
		go oschina.NewNewsList().Create(cxt, page, chOsc)
		res = *(<-chOsc)

		site = "http://www.oschina.net"
		siteMobile = "http://m.oschina.net"
	}

	//Output result of all news-list of all sites.
	s := fmt.Sprintf(`{"status":%d, "site" : "%s", "site_mobile":"%s", "result":%s}`, 200, site, siteMobile, res)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, s)
}
