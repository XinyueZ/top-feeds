package topfeeds

import (
	"csdn"
	"techug"
	"oschina"
	"geeker"
	"androider" //Deprecated
	"blogger"
	"bookmark"
	
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"encoding/json"
)

import "appengine"

type Error string

func (e Error) Error() string {
	return string(e)
}

type Blog struct {
		Id           string    `json:"id"`
}

func init() {
	http.HandleFunc("/topfeeds", handleTopFeeds) 
	http.HandleFunc("/bookmark", handleAddBookmark)
	http.HandleFunc("/bookmarkList", handleBookmarkList)
	http.HandleFunc("/removeBookmark", handleRemoveBookmark)
	//http.HandleFunc("/bookmarkTransfer", handleBookmarkListTransfer)
}

func handleTopFeeds(w http.ResponseWriter, r *http.Request) {
	cxt := appengine.NewContext(r)

	page := 0
	from := "0"
	pBlog :=  new(Blog)
	
	//Error-handling anyway.
	defer func() {
		if err := recover(); err != nil {
			cxt.Errorf("handleTopFeeds: %v", err)
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(w, `{"status":%d}`, 300)
		}
	}()

	args := r.URL.Query()

	typ := 0 
	if len(args["type"]) > 0 {
		t, _ := strconv.Atoi(args["type"][0])
		typ = t
		if typ == 5 {
			//For blogs, to get id of blog
			if bys, e := ioutil.ReadAll(r.Body); e == nil { 
				if e := json.Unmarshal(bys, pBlog); e == nil {
					//Ignore error
				}
			} else {
				//Ignore error
			}
		}
	} 

	if len(args["page"]) > 0 {
		if typ != 3 &&  typ != 4 &&  typ != 5 {
			p, _ := strconv.Atoi(args["page"][0]) //Which page, if typ is csdn(1), techug.com(2), when geekers(3), android(4),blogger(5)  ignore this.
			page = p
		} else {
			from = args["page"][0] //For geekers(3), android(4) , blogger(5)the "from" value for starting page.
		}
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
	case 3:
		//Ask geeker:
		chGeeker := make(chan *string)
		chFrom :=  make(chan *string)
		go geeker.NewNewsList().Create(cxt, from, chGeeker, chFrom)
		res = *(<-chGeeker)
		from =  *(<-chFrom)
		
		site = "http://geek.csdn.net"
		siteMobile = "http://geek.csdn.net"
	case 4: //Deprecated
		//Ask android-developer blog:
		chAndroider:= make(chan *string)
		chFrom :=  make(chan *string)
		go androider.NewNewsList().Create(cxt, from, chAndroider, chFrom)
		res = *(<-chAndroider)
		from =  *(<-chFrom)
		
		site = "http://android-developers.blogspot.de/"
		siteMobile = "http://android-developers.blogspot.de/"
	case 5:
		//Ask a blogger:
		chBlogger:= make(chan *string)
		chFrom :=  make(chan *string)
		go blogger.NewNewsList().Create(cxt, pBlog.Id, from, chBlogger, chFrom)
		res = *(<-chBlogger)
		from =  *(<-chFrom)
		
		site = "http://www.blogger.com"
		siteMobile = "http://www.blogger.com"
	default:
		//Ask news-list of www.oschina.net
		chOsc := make(chan *string)
		go oschina.NewNewsList().Create(cxt, page, chOsc)
		res = *(<-chOsc)

		site = "http://www.oschina.net"
		siteMobile = "http://m.oschina.net"
	}

	//Output result of all news-list of all sites.
	s := fmt.Sprintf(`{"status":%d, "page_index" : %d, "from" : "%s", "site" : "%s", "site_mobile":"%s", "result":%s}`, 200, page, from, site, siteMobile, res)
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
		s = fmt.Sprintf(`{"status":%d, %s}`, 200, *p)
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

func handleBookmarkListTransfer(w http.ResponseWriter, r *http.Request) {
	cxt := appengine.NewContext(r)
	ch := make(chan int)
	go bookmark.Transfer(cxt, ch)
	<-ch
	fmt.Fprintf(w, "Done, DB transfer." )
}


