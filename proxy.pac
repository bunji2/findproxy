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