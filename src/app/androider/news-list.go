package androider

import ( 
	"fmt"
	"io/ioutil"
	"net/http"
	"strings" 	 
	"time"
	"strconv"

	"appengine"
	"appengine/urlfetch"

	//"github.com/PuerkitoBio/goquery"
	//"code.google.com/p/mahonia"
	xhtml "golang.org/x/net/html"
)

const ( 
	API = "http://android-developers.blogspot.de/%d_%02d_01_archive.html" 
	ANDROID_BLOG_FORMAT = "02 January 2006 15:04:05"
)
  
func NewNewsList() (p *NewsList) {
	p = new(NewsList)
	return
}

type NewsList struct {
	 
}

func isTitle(attr []xhtml.Attribute) bool {
	for _, element := range attr {
		if element.Key == "class" && element.Val == `post-title entry-title`  {
			return true
		}
	}
	return false
}
 

func parse_xhtml(cxt appengine.Context, n *xhtml.Node, pValues *string  ) {
	link := ""
	title := ""  
	desc := ""
		
	t := time.Now()
	pubDate := t.Unix() 
	if n.Type == xhtml.ElementNode && n.Data == "h3" && isTitle(n.Attr) {
		c := n.FirstChild
		link = c.NextSibling.Attr[0].Val
		title = strings.TrimSpace(c.NextSibling.FirstChild.Data)
		datetime := ""
		if n.Parent != nil {
			if n.Parent.PrevSibling != nil {
				if n.Parent.PrevSibling.Parent != nil {
					if n.Parent.PrevSibling.Parent.Parent != nil {
						if n.Parent.PrevSibling.Parent.Parent.PrevSibling != nil {
							if n.Parent.PrevSibling.Parent.Parent.PrevSibling.PrevSibling != nil {
								if n.Parent.PrevSibling.Parent.Parent.PrevSibling.PrevSibling.FirstChild != nil {
									if n.Parent.PrevSibling.Parent.Parent.PrevSibling.PrevSibling.FirstChild.FirstChild != nil {
										datetime = n.Parent.PrevSibling.Parent.Parent.PrevSibling.PrevSibling.FirstChild.FirstChild.Data
									}
								}
							}
						}
					}
				}
			}
		}
		
 
		if datetime != "" { 
			t, _ := time.Parse(ANDROID_BLOG_FORMAT, datetime + " 15:04:05")
			pubDate = t.Unix()
		} 
		json := fmt.Sprintf(`{"title" : "%s", "desc" : "%s", "url" : "%s", "url_mobile" : "%s",  "pubDate" : %d },`,
		title, desc, link, link, pubDate) 
		*pValues += json
	} 	
	
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		parse_xhtml(cxt, c, pValues)
	}
}


func (self *NewsList) Create(cxt appengine.Context, from string, chJsonStr chan *string, chFrom chan *string) {
	client := urlfetch.Client(cxt) 
	fromDateTime := strings.Split(from, ":") 
	fromYear , _ := strconv.Atoi(fromDateTime[0])
	fromMonth , _ := strconv.Atoi(fromDateTime[1])
	if r, e := http.NewRequest("GET", fmt.Sprintf(API, fromYear, fromMonth), nil); e == nil {
		
		if resp, e := client.Do(r); e == nil {
			if resp != nil {
				defer resp.Body.Close()
			} 
			if bt, e := ioutil.ReadAll(resp.Body); e == nil { 
					jsons := ""
					source := string(bt)
				  
					pReader := strings.NewReader(source)
					doc, _ := xhtml.Parse(pReader)
					parse_xhtml(cxt, doc, &jsons)

					s := "[" //Start making a json result. 
					s+=jsons
					length := len(s)
					if length > 2 {
						s = s[:length-1] //Remove last ","
					}
					s += "]" //Stop making json
					chJsonStr <- &s
					 
					//To last year
					if fromMonth == 1 {
						fromYear-- //move to last year
						fromMonth=12 //start of month of last year.
					} else {
						fromMonth--
					}
					nextFrom := fmt.Sprintf("%d:%d", fromYear, fromMonth)
					chFrom <- &nextFrom
				 
			} else {
				chJsonStr <- nil
				chFrom <- nil
				panic(e)
			}
		} else {
			chJsonStr <- nil
			chFrom <- nil
			cxt.Errorf("Error but still going: %v", e)
		}
	} else {
		chJsonStr <- nil
		chFrom <- nil
		panic(e)
	}
}