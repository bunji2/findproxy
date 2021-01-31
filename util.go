package main

import (
	"fmt"
	"net"
	"regexp"
	"strings"
	"time"
)

// プロキシ自動設定ファイル
// https://developer.mozilla.org/ja/docs/Web/HTTP/Proxy_servers_and_tunneling/Proxy_Auto-Configuration_(PAC)_file

/*
function FindProxyForURL(url, host) {
  // ...
}

DIRECT
	接続はプロキシを使用せずに、直接行われる

PROXY host:port
	指定されたプロキシを使用する

SOCKS host:port
	指定された SOCKS サーバーを使用する
	最近の Firefox のバージョンでは次の書き方にも対応しています。

HTTP host:port
	指定されたプロキシを使用する

HTTPS host:port
	指定された HTTPS プロキシを使用する

SOCKS4 host:port
SOCKS5 host:port
	指定された SOCKS サーバーを (指定された SOCK バージョンで) 使用する

*/

/*

isPlainHostName()
構文
isPlainHostName(host)
引数
host
URLから取り出したホスト名 (ポート番号を除いたもの)
解説
ホスト名にドメイン名が含まれていない (=ドットを含まない) 場合は true になります。

例
isPlainHostName("www.mozilla.org") // false
isPlainHostName("www") // true

*/

func isPlainHostName(host string) (r bool) {
	r = strings.Index(host, ".") < 0
	return
}

/*
dnsDomainIs()
構文
dnsDomainIs(host, domain)
引数
host
URL から取り出したホスト名
domain
ホストが所属しているか確認したいドメイン名
解説
ホスト名のドメインが一致する場合にのみ true を返します。

例
dnsDomainIs("www.mozilla.org", ".mozilla.org") // true
dnsDomainIs("www", ".mozilla.org") // false

*/

func dnsDomainIs(host, domain string) (r bool) {
	if domain == "" || domain[0] != '.' {
		return
	}
	r = strings.HasSuffix(host, domain)
	return
}

/*
localHostOrDomainIs()
構文
localHostOrDomainIs(host, hostdom)
引数
host
URLから取り出したホスト名です。
hostdom
比較対象の完全修飾ホスト名です。
解説
ホスト名が完全に指定されたホスト名と一致した場合、またはホスト名にドメイン名の部分がなく、修飾されていないホスト名が一致する場合に true を返します。

例
localHostOrDomainIs("www.mozilla.org" , "www.mozilla.org") // true (exact match)
localHostOrDomainIs("www"             , "www.mozilla.org") // true (hostname match, domain not specified)
localHostOrDomainIs("www.google.com"  , "www.mozilla.org") // false (domain name mismatch)
localHostOrDomainIs("home.mozilla.org", "www.mozilla.org") // false (hostname mismatch)
*/

func localHostOrDomainIs(host, hostdom string) (r bool) {
	r = strings.HasPrefix(hostdom, host)
	return
}

/*
isResolvable()
構文
isResolvable(host)
引数
host
URLから取り出したホスト名です。
ホスト名の解決を試みます。成功すれば true を返します。

例:
isResolvable("www.mozilla.org") // true
isInNet()
構文
*/

func isResolvable(host string) (r bool) {
	_, err := net.LookupHost(host)
	if err != nil {
		return
	}
	r = true
	return
}

/*
isInNet(host, pattern, mask)
Parameters
host
a DNS hostname, or IP address. If a hostname is passed, it will be resolved into an IP address by this function.
pattern
an IP address pattern in the dot-separated format.
mask
mask for the IP address pattern informing which parts of the IP address should be matched against. 0 means ignore, 255 means match.
True if and only if the IP address of the host matches the specified IP address pattern.

Pattern and mask specification is done the same way as for SOCKS configuration.

Examples:
function alert_eval(str) { alert(str + ' is ' + eval(str)) }
function FindProxyForURL(url, host) {
  alert_eval('isInNet(host, "63.245.213.24", "255.255.255.255")')
  // "PAC-alert: isInNet(host, "63.245.213.24", "255.255.255.255") is true"
}
*/

func isInNet(host, pattern, mask string) (r bool) {
	ipNet := &net.IPNet{
		IP:   net.ParseIP(pattern),
		Mask: net.IPMask(net.ParseIP(mask)),
	}

	addrs, err := net.LookupHost(host)
	if err != nil {
		return
	}

	for _, addr := range addrs {
		if ipNet.Contains(net.ParseIP(addr)) {
			r = true
			break
		}
	}

	// [MEMO] above implimentation is slow...

	return
}

/*
dnsResolve()
dnsResolve(host)
Parameters
host
hostname to resolve.
Resolves the given DNS hostname into an IP address, and returns it in the dot-separated format as a string.

Example
dnsResolve("www.mozilla.org"); // returns the string "104.16.41.2"
*/

func dnsResolve(host string) (r string) {
	addrs, err := net.LookupHost(host)
	if err != nil {
		return
	}
	r = addrs[0]
	return
}

/*
convert_addr()
Syntax
convert_addr(ipaddr)
Parameters
ipaddr
Any dotted address such as an IP address or mask.
Concatenates the four dot-separated bytes into one 4-byte word and converts it to decimal.

Example
convert_addr("104.16.41.2"); // returns the decimal number 1745889538
*/

func convertAddr(ipaddr string) (r uint32) {
	for _, b := range net.ParseIP(ipaddr) {
		r = r*256 + uint32(b)
	}
	return
}

/*
myIpAddress()
Syntax
myIpAddress()
Parameters
(none)

Returns the server IP address of the machine Firefox is running on, as a string in the dot-separated integer format.

myIpAddress() returns the same IP address as the server address returned by nslookup localhost on a Linux machine. It does not return the public IP address.

Example
myIpAddress() //returns the string "127.0.1.1" if you were running Firefox on that localhost
*/

func myIPAddress() (r string) {
	r = "127.0.0.1"
	return
}

/*
dnsDomainLevels()
Syntax
dnsDomainLevels(host)
Parameters
host
is the hostname from the URL.
Returns the number (integer) of DNS domain levels (number of dots) in the hostname.

Examples:
dnsDomainLevels("www");             // 0
dnsDomainLevels("mozilla.org");     // 1
dnsDomainLevels("www.mozilla.org"); // 2
*/

func dnsDomainLevels(host string) (r int) {
	a := strings.Split(host, ".")
	if len(a) <= 2 {
		r = len(a) - 1
	} else {
		r = 2
	}
	return
}

/*
shExpMatch()
Syntax
shExpMatch(str, shexp)
Parameters
str
is any string to compare (e.g. the URL, or the hostname).
shexp
is a shell expression to compare against.
Returns true if the string matches the specified shell expression.

Note that the patterns are shell glob expressions, not regular expressions. * and ? are always supported, while [characters] and [^characters] are supported by some implmentations including Firefox. This is mainly because the expression is translated to a RegExp via subsitution of [.*?]. For a reliable way to use these RegExp syntaxes, just use RegExp instead.
*/
//Examples
//shExpMatch("http://home.netscape.com/people/ari/index.html"     , "*/ari/*"); // returns true
//shExpMatch("http://home.netscape.com/people/montulli/index.html", "*/ari/*"); // returns false

func shExpMatch(url, shexp string) (r bool) {
	// 変換する順番
	// シェル   正規表現
	// .     => \.      (ピリオド)
	// *     => .*      (一文字以上の任意の文字列)
	// ?     => .       (任意の一文字)
	expr := strings.ReplaceAll(shexp, ".", "\\.")
	expr = strings.ReplaceAll(expr, "*", ".*")
	expr = strings.ReplaceAll(expr, "?", ".")
	rxp, err := regexp.Compile(expr)
	if err != nil {
		return
	}
	r = rxp.MatchString(url)
	return
}

/*
weekdayRange()
Syntax
weekdayRange(wd1, wd2, [gmt])
Note: (Before Firefox 49) wd1 must be less than wd2 if you want the function to evaluate these parameters as a range. See the warning below.

Parameters
wd1 and wd2
One of the ordered weekday strings:
"SUN"|"MON"|"TUE"|"WED"|"THU"|"FRI"|"SAT"
gmt
Is either the string "GMT" or is left out.
Only the first parameter is mandatory. Either the second, the third, or both may be left out.

If only one parameter is present, the function returns a value of true on the weekday that the parameter represents. If the string "GMT" is specified as a second parameter, times are taken to be in GMT. Otherwise, they are assumed to be in the local timezone.

If both wd1 and wd1 are defined, the condition is true if the current weekday is in between those two ordered weekdays. Bounds are inclusive, but the bounds are ordered. If the "GMT" parameter is specified, times are taken to be in GMT. Otherwise, the local timezone is used.

The order of the days matters; Before Firefox 49, weekdayRange("SUN", "SAT") will always evaluate to true. Now weekdayRange("WED", "SUN") will only evaluate true if the current day is Wednesday or Sunday.

Examples
weekdayRange("MON", "FRI");        // returns true Monday through Friday (local timezone)
weekdayRange("MON", "FRI", "GMT"); // returns true Monday through Friday (GMT timezone)
weekdayRange("SAT");               // returns true on Saturdays local time
weekdayRange("SAT", "GMT");        // returns true on Saturdays GMT time
weekdayRange("FRI", "MON");        // returns true Friday and Monday only (note, order does matter!)
*/

func weekdayRange(params ...string) (r bool) {
	var err error
	r, err = subWeekdayRange(time.Now(), params...)
	if err != nil {
		panic(err)
	}
	return
}

func subWeekdayRange(now time.Time, params ...string) (r bool, err error) {
	if len(params) > 0 && params[len(params)-1] == "GMT" {
		now = now.UTC()
		params = params[0 : len(params)-1]
	}

	w := int(now.Weekday()) // Sunday=0, Monday=1, ...

	switch len(params) {
	case 1:
		if subIsWeekday(params[0]) {
			r = (w == subWeekdayNumber(params[0]))
		} else {
			err = fmt.Errorf("abnormal parameter")
		}
	case 2:
		if subIsWeekday(params[0]) && subIsWeekday(params[1]) {
			w1 := subWeekdayNumber(params[0])
			w2 := subWeekdayNumber(params[1])
			r = (w1 == w2 && w == w1) || (w1 < w2 && w1 <= w && w <= w2) || (w1 > w2 && (w == w1 || w == w2))
		} else {
			err = fmt.Errorf("abnormal parameter")
		}
	default:
		err = fmt.Errorf("abnormal parameter")
	}

	return
}

var weekdayStrs = []string{
	"SUN", "MON", "TUE", "WED", "THU", "FRI", "SAT",
}

func subWeekdayNumber(wd string) (r int) {
	for i, s := range weekdayStrs {
		if s == wd {
			r = i
			return
		}
	}
	r = -1
	return
}

func subIsWeekday(wd string) (r bool) {
	r = subWeekdayNumber(wd) >= 0
	return
}

/*
dateRange()
Syntax
dateRange(<day> | <month> | <year>, [gmt])  // ambiguity is resolved by assuming year is greater than 31
dateRange(<day1>, <day2>, [gmt])
dateRange(<month1>, <month2>, [gmt])
dateRange(<year1>, <year2>, [gmt])
dateRange(<day1>, <month1>, <day2>, <month2>, [gmt])
dateRange(<month1>, <year1>, <month2>, <year2>, [gmt])
dateRange(<day1>, <month1>, <year1>, <day2>, <month2>, <year2>, [gmt])
Note: (Before Firefox 49) day1 must be less than day2, month1 must be less than month2, and year1 must be less than year2 if you want the function to evaluate these parameters as a range. See the warning below.

Parameters
day
Is the ordered day of the month between 1 and 31 (as an integer).
1|2|3|4|5|6|7|8|9|10|11|12|13|14|15|16|17|18|19|20|21|22|23|24|25|26|27|28|29|30|31
month
Is one of the ordered month strings below.
"JAN"|"FEB"|"MAR"|"APR"|"MAY"|"JUN"|"JUL"|"AUG"|"SEP"|"OCT"|"NOV"|"DEC"
year
Is the ordered full year integer number. For example, 2016 (not 16).
gmt
Is either the string "GMT", which makes time comparison occur in GMT timezone, or is left out. If left unspecified, times are taken to be in the local timezone.
If only a single value is specified (from each category: day, month, year), the function returns a true value only on days that match that specification. If both values are specified, the result is true between those times, including bounds, but the bounds are ordered.

The order of the days, months, and years matter; Before Firefox 49, dateRange("JAN", "DEC") will always evaluate to true. Now dateRange("DEC", "JAN") will only evaluate true if the current month is December or January.

Examples
dateRange(1);            // returns true on the first day of each month, local timezone
dateRange(1, "GMT")      // returns true on the first day of each month, GMT timezone
dateRange(1, 15);        // returns true on the first half of each month
dateRange(24, "DEC");    // returns true on 24th of December each year
dateRange("JAN", "MAR"); // returns true on the first quarter of the year

dateRange(1, "JUN", 15, "AUG");
// returns true from June 1st until August 15th, each year
// (including June 1st and August 15th)

dateRange(1, "JUN", 1995, 15, "AUG", 1995);
// returns true from June 1st, 1995, until August 15th, same year

dateRange("OCT", 1995, "MAR", 1996);
// returns true from October 1995 until March 1996
// (including the entire month of October 1995 and March 1996)

dateRange(1995);
// returns true during the entire year of 1995

dateRange(1995, 1997);
// returns true from beginning of year 1995 until the end of year 1997
*/

func dateRange(params ...interface{}) bool {
	r, err := subDateRange(time.Now(), params...)
	if err != nil {
		panic(err)
	}
	return r
}

func subDateRange(now time.Time, params ...interface{}) (r bool, err error) {

	if len(params) < 1 {
		err = fmt.Errorf("no parameter")
		return
	}

	loc := time.Local
	stmp, ok := params[len(params)-1].(string)
	if ok && stmp == "GMT" {
		params = params[0 : len(params)-1]
		now = now.UTC()
		loc = time.UTC
	}

	var t1, t2 time.Time

	switch len(params) {
	case 1:
		if subIsDay(params[0]) { // day1
			d, _ := params[0].(int)
			t1 = time.Date(now.Year(), now.Month(), d, 0, 0, 0, 0, loc)
			t2 = time.Date(now.Year(), now.Month(), d+1, 0, 0, 0, 0, loc)
		} else if subIsMonth(params[0]) { // month1
			m := subMonthNumber(params[0])
			t1 = time.Date(now.Year(), time.Month(m), 1, 0, 0, 0, 0, loc)
			t2 = time.Date(now.Year(), time.Month(m+1), 1, 0, 0, 0, 0, loc)
		} else if subIsYear(params[0]) { // year1
			y, _ := params[0].(int)
			t1 = time.Date(y, 1, 1, 0, 0, 0, 0, loc)
			t2 = time.Date(y+1, 1, 1, 0, 0, 0, 0, loc)
		} else {
			err = fmt.Errorf("abnormal parameter")
		}
	case 2:
		if subIsDay(params[0]) && subIsDay(params[1]) { // day1, day2
			d1, _ := params[0].(int)
			d2, _ := params[1].(int)
			t1 = time.Date(now.Year(), now.Month(), d1, 0, 0, 0, 0, loc)
			t2 = time.Date(now.Year(), now.Month(), d2+1, 0, 0, 0, 0, loc)
		} else if subIsDay(params[0]) && subIsMonth(params[1]) { // day1, month1
			d1, _ := params[0].(int)
			m1 := subMonthNumber(params[1])
			t1 = time.Date(now.Year(), time.Month(m1), d1, 0, 0, 0, 0, loc)
			t2 = time.Date(now.Year(), time.Month(m1), d1+1, 0, 0, 0, 0, loc)
		} else if subIsMonth(params[0]) && subIsMonth(params[1]) { // month1, month2
			m1 := subMonthNumber(params[0])
			m2 := subMonthNumber(params[1])
			t1 = time.Date(now.Year(), time.Month(m1), 1, 0, 0, 0, 0, loc)
			t2 = time.Date(now.Year(), time.Month(m2+1), 1, 0, 0, 0, 0, loc)
		} else if subIsYear(params[0]) && subIsYear(params[1]) { // year1, year2
			y1, _ := params[0].(int)
			y2, _ := params[1].(int)
			t1 = time.Date(y1, 1, 1, 0, 0, 0, 0, loc)
			t2 = time.Date(y2+1, 1, 1, 0, 0, 0, 0, loc)
		} else {
			err = fmt.Errorf("abnormal parameter")
		}
	case 4:
		if subIsDay(params[0]) && subIsMonth(params[1]) && subIsDay(params[2]) && subIsMonth(params[3]) { // day1, month1, day2, month2
			d1, _ := params[0].(int)
			m1 := subMonthNumber(params[1])
			d2, _ := params[2].(int)
			m2 := subMonthNumber(params[3])
			t1 = time.Date(now.Year(), time.Month(m1), d1, 0, 0, 0, 0, loc)
			t2 = time.Date(now.Year(), time.Month(m2), d2+1, 0, 0, 0, 0, loc)
		} else if subIsMonth(params[0]) && subIsYear(params[1]) && subIsMonth(params[2]) && subIsYear(params[3]) { // month1, year1, month2, year2
			m1 := subMonthNumber(params[0])
			y1, _ := params[1].(int)
			m2 := subMonthNumber(params[2])
			y2, _ := params[3].(int)
			t1 = time.Date(y1, time.Month(m1), 1, 0, 0, 0, 0, loc)
			t2 = time.Date(y2, time.Month(m2+1), 1, 0, 0, 0, 0, loc)
		} else {
			err = fmt.Errorf("abnormal parameter")
		}
	case 6: // day1, month1, year1, day2, monath2, year2
		if subIsDay(params[0]) && subIsMonth(params[1]) && subIsYear(params[2]) && subIsDay(params[3]) && subIsMonth(params[4]) && subIsYear(params[5]) {
			d1, _ := params[0].(int)
			m1 := subMonthNumber(params[1])
			y1, _ := params[2].(int)
			d2, _ := params[3].(int)
			m2 := subMonthNumber(params[4])
			y2, _ := params[5].(int)
			t1 = time.Date(y1, time.Month(m1), d1, 0, 0, 0, 0, loc)
			t2 = time.Date(y2, time.Month(m2), d2+1, 0, 0, 0, 0, loc)
		} else {
			err = fmt.Errorf("abnormal parameter")
		}
	default:
		err = fmt.Errorf("abnormal parameter")
	}

	if err != nil {
		return
	}

	if t1.Unix() >= t2.Unix() {
		err = fmt.Errorf("abnormal order of times")
	}

	if err != nil {
		return
	}

	r = t1.Unix() <= now.Unix() && now.Unix() < t2.Unix()

	return
}

func subIsDay(t interface{}) (r bool) {
	n, ok := t.(int)
	if !ok {
		return
	}
	r = n >= 1 && n <= 31
	return
}

var monthStrs = []string{
	"JAN", "FEB", "MAR", "APR", "MAY", "JUN", "JUL", "AUG", "SEP", "OCT", "NOV", "DEC",
}

func subMonthNumber(t interface{}) (r int) {
	s, ok := t.(string)
	if !ok {
		return
	}
	for i, m := range monthStrs {
		if m == s {
			r = i + 1
			break
		}
	}
	//fmt.Println("[debug]", t, "==>", r)
	return
}

func subIsMonth(t interface{}) (r bool) {
	r = subMonthNumber(t) > 0
	return
}

func subIsYear(t interface{}) (r bool) {
	n, ok := t.(int)
	if !ok {
		return
	}
	r = n >= 1970
	return
}

/*
https://developer.mozilla.org/ja/docs/Web/HTTP/Proxy_servers_and_tunneling/Proxy_Auto-Configuration_(PAC)_file#timeRange
timeRange()
Syntax
// The full range of expansions is analogous to dateRange.
timeRange(<hour1>, <min1>, <sec1>, <hour2>, <min2>, <sec2>, [gmt])
Note: (Before Firefox 49) the category hour1, min1, sec1 must be less than the category hour2, min2, sec2 if you want the function to evaluate these parameters as a range. See the warning below.

Parameters
hour
Is the hour from 0 to 23. (0 is midnight, 23 is 11 pm.)
min
Minutes from 0 to 59.
sec
Seconds from 0 to 59.
gmt
Either the string "GMT" for GMT timezone, or not specified, for local timezone.
If only a single value is specified (from each category: hour, minute, second), the function returns a true value only at times that match that specification. If both values are specified, the result is true between those times, including bounds, but the bounds are ordered.

The order of the hour, minute, second matter; Before Firefox 49, timeRange(0, 23) will always evaluate to true. Now timeRange(23, 0) will only evaluate true if the current hour is 23:00 or midnight.

Examples
timerange(12);                // returns true from noon to 1pm
timerange(12, 13);            // returns true from noon to 1pm
timerange(12, "GMT");         // returns true from noon to 1pm, in GMT timezone
timerange(9, 17);             // returns true from 9am to 5pm
timerange(8, 30, 17, 00);     // returns true from 8:30am to 5:00pm
timerange(0, 0, 0, 0, 0, 30); // returns true between midnight and 30 seconds past midnight

*/

func timeRange(params ...interface{}) bool {
	r, err := subTimeRange(time.Now(), params...)
	if err != nil {
		panic(err)
	}
	return r
}

func subTimeRange(now time.Time, params ...interface{}) (r bool, err error) {
	var nums []int
	loc := time.Local

	for i, p := range params {
		switch p.(type) {
		case int:
			tmp, _ := p.(int)
			if tmp < 0 || tmp > 59 {
				err = fmt.Errorf("abnormal number: %d", tmp)
				break
			}
			nums = append(nums, tmp)
		case string:
			if i > 0 && i == len(params)-1 && p.(string) == "GMT" {
				now = now.UTC()
				loc = time.UTC
			} else {
				err = fmt.Errorf("abnormal occurence of string: %s", p.(string))
			}
		default:
			err = fmt.Errorf("abnormal parameter")
		}
	}

	if err != nil {
		return
	}

	if len(nums) < 1 {
		err = fmt.Errorf("empty parameter")
		return
	}

	if nums[0] > 23 {
		err = fmt.Errorf("abnormal hour number: %d", nums[0])
		return
	}

	var t1, t2 time.Time
	switch len(nums) {
	case 1: // Hour1
		t1 = time.Date(now.Year(), now.Month(), now.Day(), nums[0], 0, 0, 0, loc)
		t2 = time.Date(now.Year(), now.Month(), now.Day(), nums[0]+1, 0, 0, 0, loc)
	case 2: // Hour1, Hour2
		if nums[1] > 23 {
			err = fmt.Errorf("abnormal hour number: %d", nums[1])
		} else {
			t1 = time.Date(now.Year(), now.Month(), now.Day(), nums[0], 0, 0, 0, loc)
			t2 = time.Date(now.Year(), now.Month(), now.Day(), nums[1], 0, 1, 0, loc)
		}
	case 4: // Hour1, Min1, Hour2, Min2
		if nums[2] > 23 {
			err = fmt.Errorf("abnormal hour number: %d", nums[2])
		} else {
			t1 = time.Date(now.Year(), now.Month(), now.Day(), nums[0], nums[1], 0, 0, loc)
			t2 = time.Date(now.Year(), now.Month(), now.Day(), nums[2], nums[3], 1, 0, loc)
		}
	case 6: // Hour1, Min1, Sec1, Hour2, Min2, Sec2
		if nums[3] > 23 {
			err = fmt.Errorf("abnormal hour number: %d", nums[3])
		} else {
			t1 = time.Date(now.Year(), now.Month(), now.Day(), nums[0], nums[1], nums[2], 0, loc)
			t2 = time.Date(now.Year(), now.Month(), now.Day(), nums[3], nums[4], nums[5]+1, 0, loc)
		}
	default:
		err = fmt.Errorf("abnormal number of parameters")
	}

	if t1.Unix() >= t2.Unix() {
		err = fmt.Errorf("abnormal order of times")
	}

	if err != nil {
		return
	}

	r = t1.Unix() <= now.Unix() && now.Unix() < t2.Unix()
	return
}
