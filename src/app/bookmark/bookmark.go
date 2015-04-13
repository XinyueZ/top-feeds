package bookmark

type BookmarkList struct {
	Status          int             `json:"status"`
	BookmarkEntries []BookmarkEntry `json:"result"`
}

type BookmarkEntry struct {
	Ident   string 
	Title       string `json:"title"`
	Description string `json:"desc"`
	PubDate     int    `json:"pubDate"`
	Url         string `json:"url"`
	UrlMobile   string `json:"url_mobile"`
}
