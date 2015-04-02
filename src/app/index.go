package topfeeds

import (
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
			fmt.Fprintf(w, `{"status":%d}`, 300)
		}
	}()

	args := r.URL.Query()
	typ, _ := strconv.Atoi(args["type"][0])  //Which type, 0: oschina, 1: csdn
	page, _ := strconv.Atoi(args["page"][0]) //Which page, if typ is csdn, then ignore this.

	site := ""
	siteMobile := ""
	res := ""
	switch typ {
	case 1:
		//Ask csdn:
		res = "coming soon."
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
