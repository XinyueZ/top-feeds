Top Feeds  
==============
Getting top-news of IT from comminuty [CSDN](http://www.csdn.net)  and
[OSCHINA](http://www.oschina)


#API

####Url: http://top-feeds-90308.appspot.com/topfeeds

####Parameters

Var     |  Value
--------|---------
type    | 1 ([CSDN](http://www.csdn.net)) or other([OSCHINA](http://www.oschina.net)), default is other.
page    | The page-index >= 0. But it works only when type == 0 ([OSCHINA](http://www.oschina.net))


####Result

######List

Var      | Type     | Comment
---------|---------|---------
status        |int   |200 is OK, 300 or other is error.
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
pubDate        |string    |News published date.

####Example

```
{
  "status": 200,
  "site": "http://www.oschina.net",
  "site_mobile": "http://m.oschina.net",
  "result": [
              {
                "title": "Java HeartBeat 0.3 发布，应用服务器心跳检测",
                "desc": "",
                "url": "http://www.oschina.net/news/61151",
                "url_mobile": "http://m.oschina.net/news/61151",
                "pubDate": "2015-04-02 18:46:25"
              },
              {
                "title": "JSUtil 1.1.2 开始支持存储过程调用啦！",
                "desc": "",
                "url": "http://www.oschina.net/news/61149",
                "url_mobile": "http://m.oschina.net/news/61149",
                "pubDate": "2015-04-02 15:03:14"
              },
              ......
        ]
}

```

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
