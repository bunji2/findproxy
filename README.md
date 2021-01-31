# findproxy

Go implementation of FindProxyForURL

## Usage

```
C:\work> findproxy.exe
findproxy.exe proxy.pac url...
```

## Proxy.pac

[Proxy Auto Configuration file](https://developer.mozilla.org/ja/docs/Web/HTTP/Proxy_servers_and_tunneling/Proxy_Auto-Configuration_(PAC)_file)

## Sample

```javascript
// proxy.pac
function FindProxyForURL(url, host) {
    if (dnsDomainIs(host, ".foo.co.jp")) {
        return "PROXY proxy1:8000"
    }
    if (shExpMatch(host, "*.com")) {
        return "PROXY proxy2:8080";
    }
    if (isInNet(host, "192.168.1.0", "255.255.255.0")) {
        return "PROXY 192.168.3.2:8000";
    }
    return "DIRECT";

}
```

```
C:\work> findproxy.exe proxy.pac http://hogehoge/hoge http://hoge.com/hoge http://192.168.1.45/ http://www.foo.co.jp/
http://hogehoge/hoge => DIRECT
http://hoge.com/hoge => PROXY proxy2:8080
http://192.168.1.45/ => PROXY 192.168.3.2:8000
http://www.foo.co.jp/ => PROXY proxy1:8000
```