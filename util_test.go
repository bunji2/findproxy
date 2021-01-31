package main

import (
	"testing"
	"time"
)

func TestIsPlainHostName(t *testing.T) {
	pats := map[string]bool{
		"www.mozilla.org": false,
		"www":             true,
	}
	for host, want := range pats {
		got := isPlainHostName(host)
		if got != want {
			t.Errorf("isPlainHostName(%s) = %v; want %v", host, got, want)
		}
	}
}

func TestDnsDomainIs(t *testing.T) {
	pats := map[[2]string]bool{
		{"www.mozilla.org", ".mozilla.org"}: true,
		{"www", ".mozilla.org"}:             false,
	}
	for args, want := range pats {
		got := dnsDomainIs(args[0], args[1])
		if got != want {
			t.Errorf("dnsDomainIs(%s, %s) = %v; want %v", args[0], args[1], got, want)
		}
	}
}

func TestLocalHostOrDomainIs(t *testing.T) {
	pats := map[[2]string]bool{
		{"www.mozilla.org", "www.mozilla.org"}:  true,
		{"www", "www.mozilla.org"}:              true,
		{"www.google.com", "www.mozilla.org"}:   false,
		{"home.mozilla.org", "www.mozilla.org"}: false,
	}
	for args, want := range pats {
		got := localHostOrDomainIs(args[0], args[1])
		if got != want {
			t.Errorf("localHostOrDomainIs(%s, %s) = %v; want %v", args[0], args[1], got, want)
		}
	}
}

func TestIsResolvable(t *testing.T) {
	pats := map[string]bool{
		"www.mozilla.org": true,
		"www":             false,
	}
	for host, want := range pats {
		got := isResolvable(host)
		if got != want {
			t.Errorf("isResolvable(%s) = %v; want %v", host, got, want)
		}
	}
}

func TestIsInNet(t *testing.T) {
	pats := map[[3]string]bool{
		{"63.245.213.3", "63.245.213.24", "255.255.255.0"}:    true,
		{"63.245.213.24", "63.245.213.24", "255.255.255.255"}: true,
		{"63.245.213.3", "63.245.213.24", "255.255.255.255"}:  false,
		//[3]string{}:      true,
	}
	for args, want := range pats {
		got := isInNet(args[0], args[1], args[2])
		if got != want {
			t.Errorf("isInNet(%s, %s, %s) = %v; want %v", args[0], args[1], args[2], got, want)
		}
	}
}

func TestDnsResolve(t *testing.T) {
	pats := map[string]string{
		"www.isc2.org": "107.162.133.105",
	}
	for host, want := range pats {
		got := dnsResolve(host)
		if got != want {
			t.Errorf("dnsResolve(%s) = %v; want %v", host, got, want)
		}
	}
}

func TestConvertAddr(t *testing.T) {
	pats := map[string]uint32{
		"104.16.41.2": 1745889538,
	}
	for ipaddr, want := range pats {
		got := convertAddr(ipaddr)
		if got != want {
			t.Errorf("convertAddr(%s) = %v; want %v", ipaddr, got, want)
		}
	}
}

func TestMyIPAddress(t *testing.T) {
	want := "127.0.0.1"
	got := myIPAddress()
	if got != want {
		t.Errorf("myIPAddress() = %v; want %v", got, want)
	}
}

func TestDnsDomainLevels(t *testing.T) {
	pats := map[string]int{
		"www":             0,
		"mozilla.org":     1,
		"www.mozilla.org": 2,
	}

	for host, want := range pats {
		got := dnsDomainLevels(host)
		if got != want {
			t.Errorf("dnsDomainLevels(%s) = %v; want %v", host, got, want)
		}
	}
}

func TestShExpMatch(t *testing.T) {
	pats := map[[2]string]bool{
		{"http://home.netscape.com/people/ari/index.html", "*/ari/*"}:      true,
		{"http://home.netscape.com/people/montulli/index.html", "*/ari/*"}: false,
		{"https://www.foo.com/abc.jpg", "*.jpg"}:                           true,
		{"https://www.foo.com/abc.jpg", "https://*.jpg"}:                   true,
		{"https://www.foo.com/abc.jpg", "https://www.foo.com/???.jpg"}:     true,
		{"https://www.foo.com/abc.jpg", "https://*/abc.jpg"}:               true,
	}
	for args, want := range pats {
		got := shExpMatch(args[0], args[1])
		if got != want {
			t.Errorf("shExpMatch(%s, %s) = %v; want %v", args[0], args[1], got, want)
		}
	}
}

func TestWeekdayRange(t *testing.T) {
	// https://developer.mozilla.org/ja/docs/Web/HTTP/Proxy_servers_and_tunneling/Proxy_Auto-Configuration_(PAC)_file#weekdayRange

	jst, _ := time.LoadLocation("Asia/Tokyo")

	// 2021/1/4 12:00:00 0nsec (Monday) JST
	now := time.Date(2021, 1, 4, 12, 0, 0, 0, jst)

	r, err := subWeekdayRange(now, "MON", "FRI")
	if err != nil && !r {
		t.Errorf("subWeekdayRange(now, \"MON\", \"FRI\") = false, nil; want true")
	}

	r, err = subWeekdayRange(now, "MON", "FRI", "GMT")
	if err != nil && !r {
		t.Errorf("subWeekdayRange(now, \"MON\", \"FRI\", \"GMT\") = false, nil; want true")
	}

	// 2021/1/2 12:00:00 0nsec (Saturday) JST
	now = time.Date(2021, 1, 2, 12, 0, 0, 0, jst)

	r, err = subWeekdayRange(now, "SAT")
	if err != nil && !r {
		t.Errorf("subWeekdayRange(now, \"SAT\") = false, nil; want true")
	}

	r, err = subWeekdayRange(now, "SAT", "GMT")
	if err != nil && !r {
		t.Errorf("subWeekdayRange(now, \"SAT\", \"GMT\") = false, nil; want true")
	}

	// 2021/1/1 12:00:00 0nsec (Friday) JST
	now = time.Date(2021, 1, 1, 12, 0, 0, 0, jst)

	r, err = subWeekdayRange(now, "FRI", "MON")
	if err != nil && !r {
		t.Errorf("subWeekdayRange(now, \"FRI\", \"MON\") = false, nil; want true")
	}

	// 2021/1/4 12:00:00 0nsec (Monday) JST
	now = time.Date(2021, 1, 4, 12, 0, 0, 0, jst)

	r, err = subWeekdayRange(now, "FRI", "MON")
	if err != nil && !r {
		t.Errorf("subWeekdayRange(now, \"FRI\", \"MON\") = false, nil; want true")
	}

	r, err = subWeekdayRange(now, "FOO")
	if err == nil {
		t.Errorf("subWeekdayRange(now, \"FOO\") = _, nil; want !nil")
	}

	r, err = subWeekdayRange(now, "SUN", "MON", "FRI", "SAT")
	if err == nil {
		t.Errorf("subWeekdayRange(now, \"SUN\", \"MON\", \"FRI\", \"SAT\") = _, nil; want !nil")
	}

}

func TestDateRange(t *testing.T) {
	// https://developer.mozilla.org/ja/docs/Web/HTTP/Proxy_servers_and_tunneling/Proxy_Auto-Configuration_(PAC)_file#dateRange

	jst, _ := time.LoadLocation("Asia/Tokyo")

	// 2021/1/1 12:30:35 0nsec local
	now := time.Date(2021, 1, 1, 12, 30, 35, 0, jst)

	r, err := subDateRange(now, 1)
	if err != nil && !r {
		t.Errorf("subDateRange(now, 1) = false, nil; want true")
	}
	r, err = subDateRange(now, 1, "GMT")
	if err != nil && !r {
		t.Errorf("subDateRange(now, 1, \"GMT\") = false, nil; want true")
	}
	r, err = subDateRange(now, 1, 15)
	if err != nil && !r {
		t.Errorf("subDateRange(now, 1, 15) = false, nil; want true")
	}

	// 2021/12/24 12:30:35 0nsec local
	now = time.Date(2021, 12, 24, 12, 30, 35, 0, jst)

	r, err = subDateRange(now, 24, "DEC")
	if err != nil && !r {
		t.Errorf("subDateRange(now, 24, \"DEC\") = false, nil; want true")
	}

	// 2021/2/24 12:30:35 0nsec local
	now = time.Date(2021, 2, 24, 12, 30, 35, 0, jst)

	r, err = subDateRange(now, "JAN", "MAR")
	if err != nil && !r {
		t.Errorf("subDateRange(now, \"JAN\", \"MAR\") = false, nil; want true")
	}

	// 2021/7/24 12:30:35 0nsec local
	now = time.Date(2021, 7, 24, 12, 30, 35, 0, jst)

	r, err = subDateRange(now, 1, "JUN", 15, "AUG")
	if err != nil && !r {
		t.Errorf("subDateRange(now, 1, \"JUN\", 15, \"AUG\") = false, nil; want true")
	}

	// 1995/7/24 12:30:35 0nsec local
	now = time.Date(1995, 7, 24, 12, 30, 35, 0, jst)

	r, err = subDateRange(now, 1, "JUN", 1995, 15, "AUG", 1995)
	if err != nil && !r {
		t.Errorf("subDateRange(now, 1, \"JUN\", 1995, 15, \"AUG\", 1995) = false, nil; want true")
	}

	// 1996/1/24 12:30:35 0nsec local
	now = time.Date(1996, 1, 24, 12, 30, 35, 0, jst)

	r, err = subDateRange(now, "OCT", 1995, "MAR", 1996)
	if err != nil && !r {
		t.Errorf("subDateRange(now, \"OCT\", 1995, \"MAR\", 1996) = false, nil; want true")
	}

	// 1995/1/24 12:30:35 0nsec local
	now = time.Date(1995, 1, 24, 12, 30, 35, 0, jst)

	r, err = subDateRange(now, 1995)
	if err != nil && !r {
		t.Errorf("subDateRange(now, 1995) = false, nil; want true")
	}

	// 1996/1/24 12:30:35 0nsec local
	now = time.Date(1996, 1, 24, 12, 30, 35, 0, jst)

	r, err = subDateRange(now, 1995, 1997)
	if err != nil && !r {
		t.Errorf("subDateRange(now, 1995) = false, nil; want true")
	}

	r, err = subDateRange(now, -1)
	if err == nil {
		t.Errorf("subDateRange(now, -1) = _, nil; want !nil")
	}

	r, err = subDateRange(now, 32)
	if err == nil {
		t.Errorf("subDateRange(now, 32) = _, nil; want !nil")
	}

	r, err = subDateRange(now, "FOO")
	if err == nil {
		t.Errorf("subDateRange(now, \"FOO\") = _, nil; want !nil")
	}

	r, err = subDateRange(now, 3.14)
	if err == nil {
		t.Errorf("subDateRange(now, 3.14) = _, nil; want !nil")
	}

	r, err = subDateRange(now, true)
	if err == nil {
		t.Errorf("subDateRange(now, true) = _, nil; want !nil")
	}

	r, err = subDateRange(now, 2, 2, 2021, 1, 1, 2020)
	if err == nil {
		t.Errorf("subDateRange(now, 2, 2, 2021, 1, 1, 2020) = _, nil; want !nil")
	}

	r, err = subDateRange(now)
	if err == nil {
		t.Errorf("subDateRange(now) = _, nil; want !nil")
	}

	r, err = subDateRange(now, 1, 2, 3, 4, 5, 6, 7)
	if err == nil {
		t.Errorf("subDateRange(now, 1, 2, 3, 4, 5, 6, 7) = _, nil; want !nil")
	}

}

func TestTimeRange(t *testing.T) {
	// https://developer.mozilla.org/ja/docs/Web/HTTP/Proxy_servers_and_tunneling/Proxy_Auto-Configuration_(PAC)_file#timeRange

	jst, _ := time.LoadLocation("Asia/Tokyo")

	// 2021/1/3 12:30:35 0nsec local
	now := time.Date(2021, 1, 3, 12, 30, 35, 0, jst)

	r, err := subTimeRange(now, 12)
	if err != nil && !r {
		t.Errorf("subTimeRange(now, 12) = false, nil; want true")
	}
	r, err = subTimeRange(now, 12, 13)
	if err != nil && !r {
		t.Errorf("subTimeRange(now, 12, 13) = false, nil; want true")
	}
	r, err = subTimeRange(now, 3, "GMT")
	if err != nil && !r {
		t.Errorf("subTimeRange(now, 3, \"GMT\") = false, nil; want true")
	}
	r, err = subTimeRange(now, 9, 17)
	if err != nil && !r {
		t.Errorf("subTimeRange(now, 12, 13) = false, nil; want true")
	}

	// 2021/1/4 0:0:15 0nsec local
	now = time.Date(2021, 1, 4, 0, 0, 15, 0, jst)

	r, err = subTimeRange(now, 0, 0, 0, 0, 0, 30)
	if err != nil && !r {
		t.Errorf("subTimeRange(now, 0, 0, 0, 0, 0, 30) = false, nil; want true")
	}

	r, err = subTimeRange(now, -1)
	if err == nil {
		t.Errorf("subTimeRange(now, -1) = _, nil; want !nil")
	}
	r, err = subTimeRange(now, 24)
	if err == nil {
		t.Errorf("subTimeRange(now, 24) = _, nil; want !nil")
	}
	r, err = subTimeRange(now, 0, -1)
	if err == nil {
		t.Errorf("subTimeRange(now, 0, -1) = _, nil; want !nil")
	}
	r, err = subTimeRange(now, 0, 24)
	if err == nil {
		t.Errorf("subTimeRange(now, 0, 24) = _, nil; want !nil")
	}

	r, err = subTimeRange(now, 0, -1, 1, 1)
	if err == nil {
		t.Errorf("subTimeRange(now, 0, -1, 1, 1) = _, nil; want !nil")
	}
	r, err = subTimeRange(now, 0, 60, 1, 1)
	if err == nil {
		t.Errorf("subTimeRange(now, 0, 60, 1, 1) = _, nil; want !nil")
	}

	r, err = subTimeRange(now, 0, 0, 1, -1)
	if err == nil {
		t.Errorf("subTimeRange(now, 0, 0, 1, -1) = _, nil; want !nil")
	}
	r, err = subTimeRange(now, 0, 0, 1, 60)
	if err == nil {
		t.Errorf("subTimeRange(now, 0, 0, 1, 60) = _, nil; want !nil")
	}

	r, err = subTimeRange(now, 0, 0, -1, 1, 1, 1)
	if err == nil {
		t.Errorf("subTimeRange(now, 0, 0, -1, 1, 1, 1) = _, nil; want !nil")
	}
	r, err = subTimeRange(now, 0, 0, 60, 1, 1, 1)
	if err == nil {
		t.Errorf("subTimeRange(now, 0, 0, 60, 1, 1, 1) = _, nil; want !nil")
	}

	r, err = subTimeRange(now, 0, 0, 0, 1, 1, -1)
	if err == nil {
		t.Errorf("subTimeRange(now, 0, 0, 0, 1, 1, -1) = _, nil; want !nil")
	}
	r, err = subTimeRange(now, 0, 0, 0, 1, 1, 60)
	if err == nil {
		t.Errorf("subTimeRange(now, 0, 0, 0, 1, 1, 60) = _, nil; want !nil")
	}

	r, err = subTimeRange(now, "1")
	if err == nil {
		t.Errorf("subTimeRange(now, \"1\") = _, nil; want !nil")
	}

	r, err = subTimeRange(now, 1, "UTC")
	if err == nil {
		t.Errorf("subTimeRange(now, 1, \"UTC\") = _, nil; want !nil")
	}

	r, err = subTimeRange(now, 1.0)
	if err == nil {
		t.Errorf("subTimeRange(now, 1.0) = _, nil; want !nil")
	}

	r, err = subTimeRange(now, 1, 2, 3)
	if err == nil {
		t.Errorf("subTimeRange(now, 1, 2, 3) = _, nil; want !nil")
	}

	r, err = subTimeRange(now)
	if err == nil {
		t.Errorf("subTimeRange(now) = _, nil; want !nil")
	}

	r, err = subTimeRange(now, 1, 2, 3, 4, 5, 6, 7)
	if err == nil {
		t.Errorf("subTimeRange(now, 1, 2, 3, 4, 5, 6, 7) = _, nil; want !nil")
	}

}
