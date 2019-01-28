# sURL

URL Shortener used by [levpay](https://github.com/levpay)

## Example

``` shell
curl -i 'http://127.0.0.1:8080/' -H 'Content-Type: application/json; charset=UTF-8' --data-binary '{"url":"https://www.google.com/"}' --compressed
HTTP/1.1 201 Created
Content-Type: application/json
Date: Mon, 28 Jan 2019 13:29:10 GMT
Content-Length: 50

{"url":"https://www.google.com/","short":"XBvHFJ"}%
```

### redirect

``` shell
curl -i http://127.0.0.1:8080/XBvHFJ
HTTP/1.1 301 Moved Permanently
Content-Type: text/html; charset=utf-8
Location: https://www.google.com/
Date: Mon, 28 Jan 2019 13:30:26 GMT
Content-Length: 58

<a href="https://www.google.com/">Moved Permanently</a>.
```
