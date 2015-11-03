Top Feeds  
==============
Recent IT information and news from China, some tech-blogs hosted by Blogs from [Blogger](http://www.blogger.com) .

From comminuty 

[CSDN](http://www.csdn.net)

[Geek news](http://geek.csdn.net/) 

[Techug](http://www.techug.com)

[OSCHINA](http://www.oschina) 
 
Blogs from [Blogger](http://www.blogger.com) 

#API

##List of all feeds

####Url: /topfeeds

####Method: ```GET```

####Parameters

Var     |  Value
--------|---------
type    | 1 ([CSDN](http://www.csdn.net)), 2 ([Techug](http://www.techug.com)) , 3 ([Geek news](http://geek.csdn.net/)), 5 ([Blogger](http://www.blogger.com)) or other([OSCHINA](http://www.oschina.net)), default is other.
page    | The page-index >= 0 or page from. But it works only when type == 0 ([OSCHINA](http://www.oschina.net)), 3 ([Geek news](http://geek.csdn.net/))


####Example

[OSCHINA](http://www.oschina): http://top-feeds-90308.appspot.com/topfeeds?type=0&page=0

[Techug](http://www.oschina): http://top-feeds-90308.appspot.com/topfeeds?type=2

[CSDN](http://www.csdn.net) : http://top-feeds-90308.appspot.com/topfeeds?type=1

[Geek news](http://geek.csdn.net/): http://top-feeds-90308.appspot.com/topfeeds?type=3&page=6:40577

Blogs from [Blogger](http://www.blogger.com): http://top-feeds2-91415.appspot.com/topfeeds?type=5&page=CgkIChjgqtaE9ikQ64mAu4jExuBd

####Result

######List

Var      | Type     | Comment
---------|---------|---------
status        |int   |200 is OK, 300 or other is error.
from  |string  |Start Page from for [Geek news](http://geek.csdn.net/) and Blogs from [Blogger](http://www.blogger.com. Works for type=3 and 5
page_index        |int   |Page index of list. Only works for [OSCHINA](http://www.oschina), type!=1,2,3.
site        |string   |The host of feeds.
site_mobile        |string   |The mobile-host of feeds.
result        |array    |The list of news-entry.

######News-entry

Var      | Type     | Comment
---------|---------|---------
title        |string   |News title.  
desc        |string   |Description of news, might be empty.
url        |string   |News location.
url_mobile        |string    |News location on mobile-host.
pubDate        |long    |News published date in timestamp.

####Example

```
{
  "status": 200,
  "from": "6:40577",
  "page_index": 0,
  "site": "http://www.oschina.net",
  "site_mobile": "http://m.oschina.net",
  "result": [
              {
                "title": "Java HeartBeat 0.3 发布，应用服务器心跳检测",
                "desc": "",
                "url": "http://www.oschina.net/news/61151",
                "url_mobile": "http://m.oschina.net/news/61151",
                "pubDate": 234523452345
              },
              {
                "title": "JSUtil 1.1.2 开始支持存储过程调用啦！",
                "desc": "",
                "url": "http://www.oschina.net/news/61149",
                "url_mobile": "http://m.oschina.net/news/61149",
                "pubDate": 234523452345
              },
              ......
        ]
}

```

##Add to and remove from bookmark

####Url: /bookmark , /removeBookmark

####Method: ```GET```

####Parameters

Var     |  Value
--------|---------
ident    | An identifier of sender.(Device Id or UUID etc)

####Body

A json object of News-entry.


Var      | Type     | Comment
---------|---------|---------
title        |string   |News title.  
desc        |string   |Description of news, might be empty.
url        |string   |News location.
url_mobile        |string    |News location on mobile-host.
pubDate        |long    |News published date in timestamp.

####Result

Var      | Type     | Comment
---------|---------|---------
status        |int   |200 is OK, 300 or other is error.

####Example

```
{
  "status": 200
}

```


##List bookmark

####Url: /bookmarkList

####Method: ```GET```

####Parameters

Var     |  Value
--------|---------
ident    | An identifier of sender.(Device Id or UUID etc)

####Result

######List

Var      | Type     | Comment
---------|---------|---------
status        |int   |200 is OK, 300 or other is error.
result        |array    |The list of news-entry.

######News-entry

Var      | Type     | Comment
---------|---------|---------
title        |string   |News title.  
desc        |string   |Description of news, might be empty.
url        |string   |News location.
url_mobile        |string    |News location on mobile-host.
pubDate        |long    |News published date in timestamp.

####Example

```
{
  "status": 200,
  "result": [
              {
                "title": "Java HeartBeat 0.3 发布，应用服务器心跳检测",
                "desc": "",
                "url": "http://www.oschina.net/news/61151",
                "url_mobile": "http://m.oschina.net/news/61151",
                "pubDate": 234523452345
              },
              {
                "title": "JSUtil 1.1.2 开始支持存储过程调用啦！",
                "desc": "",
                "url": "http://www.oschina.net/news/61149",
                "url_mobile": "http://m.oschina.net/news/61149",
                "pubDate": 234523452345
              },
              ......
        ]
}

```



#SDK for Android

There's a SDK based on this API.

https://github.com/XinyueZ/top-feeds-client

Example App:

[![https://play.google.com/store/apps/details?id=com.topfeeds4j.sample](https://dl.dropbox.com/s/phrg0387osr3riz/images.jpeg)](https://play.google.com/store/apps/details?id=com.topfeeds4j.sample)

#License

```
                    The MIT License (MIT)

            Copyright (c) 2015 Chris Xinyue Zhao

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

```
