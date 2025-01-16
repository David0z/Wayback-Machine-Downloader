STEPS:

1. https://web.archive.org/web/timemap/json?url=Ma3ali.net&matchType=prefix&output=json
JSON RESPONSE:
```
[
  [
      "urlkey",
      "timestamp",
      "original",
      "mimetype",
      "statuscode",
      "digest",
      "length"
  ],
  [
      "net,ma3ali)/",
      "20040528102848",
      "http://www.ma3ali.net:80/",
      "text/html",
      "200",
      "3I42H3S6NNFQ2MSVX7XZKYAYSCX5QBYJ",
      "337"
  ]
]
```

2. filter by mimetype beginning with "audio"

3. if "statuscode" === "200" -> download https://web.archive.org/web/{timestamp}/{original}

else, get https://web.archive.org/cdx/search/cdx?url={original}
TEXT RESPONSE:
```
net,ma3ali)/forsan/eaemhh.mp3 20060421181949 http://ma3ali.net:80/forsan/eaemhh.MP3 audio/mpeg 200 WGFTOAVL6V37BHU3ILCFSYKWFAKEHPE2 3213174
net,ma3ali)/forsan/eaemhh.mp3 20060421182010 http://ma3ali.net:80/forsan/eaemhh.MP3 audio/mpeg 200 WGFTOAVL6V37BHU3ILCFSYKWFAKEHPE2 3213174
net,ma3ali)/forsan/eaemhh.mp3 20060614053529 http://ma3ali.net:80/forsan/eaemhh.MP3 audio/mpeg 200 WGFTOAVL6V37BHU3ILCFSYKWFAKEHPE2 3213162
net,ma3ali)/forsan/eaemhh.mp3 20060614184751 http://ma3ali.net:80/forsan/eaemhh.MP3 audio/mpeg 200 WGFTOAVL6V37BHU3ILCFSYKWFAKEHPE2 3213163
net,ma3ali)/forsan/eaemhh.mp3 20061210085305 http://ma3ali.net:80/forsan/eaemhh.MP3 audio/mpeg 200 WGFTOAVL6V37BHU3ILCFSYKWFAKEHPE2 3213161
net,ma3ali)/forsan/eaemhh.mp3 20061215002138 http://ma3ali.net:80/forsan/eaemhh.MP3 text/html 404 TGPY5X2HDTPWIEGESQF364WFWQ5L7HYN 666
net,ma3ali)/forsan/eaemhh.mp3 20061220175100 http://ma3ali.net:80/forsan/eaemhh.MP3 text/html 404 TGPY5X2HDTPWIEGESQF364WFWQ5L7HYN 667
net,ma3ali)/forsan/eaemhh.mp3 20061230155026 http://ma3ali.net:80/forsan/eaemhh.MP3 text/html 404 TGPY5X2HDTPWIEGESQF364WFWQ5L7HYN 667
net,ma3ali)/forsan/eaemhh.mp3 20061230230320 http://ma3ali.net:80/forsan/eaemhh.MP3 text/html 404 TGPY5X2HDTPWIEGESQF364WFWQ5L7HYN 666
net,ma3ali)/forsan/eaemhh.mp3 20070104160615 http://ma3ali.net:80/forsan/eaemhh.MP3 text/html 404 TGPY5X2HDTPWIEGESQF364WFWQ5L7HYN 667
net,ma3ali)/forsan/eaemhh.mp3 20070114150023 http://ma3ali.net:80/forsan/eaemhh.MP3 text/html 404 TGPY5X2HDTPWIEGESQF364WFWQ5L7HYN 690
net,ma3ali)/forsan/eaemhh.mp3 20070119135716 http://ma3ali.net:80/forsan/eaemhh.MP3 text/html 404 TGPY5X2HDTPWIEGESQF364WFWQ5L7HYN 690
net,ma3ali)/forsan/eaemhh.mp3 20070124202016 http://ma3ali.net:80/forsan/eaemhh.MP3 text/html 404 TGPY5X2HDTPWIEGESQF364WFWQ5L7HYN 689
net,ma3ali)/forsan/eaemhh.mp3 20070205223148 http://ma3ali.net:80/forsan/eaemhh.MP3 text/html 404 TGPY5X2HDTPWIEGESQF364WFWQ5L7HYN 690
net,ma3ali)/forsan/eaemhh.mp3 20070216215410 http://ma3ali.net:80/forsan/eaemhh.MP3 text/html 404 TGPY5X2HDTPWIEGESQF364WFWQ5L7HYN 689
net,ma3ali)/forsan/eaemhh.mp3 20070302115614 http://ma3ali.net:80/forsan/eaemhh.MP3 text/html 404 TGPY5X2HDTPWIEGESQF364WFWQ5L7HYN 689
```

split by space, get first with status 200
download https://web.archive.org/web/{timestamp}/{original}