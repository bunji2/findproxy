// builtins.go
// 組み込み関数

package main

// BuiltIns は組み込み関数を格納する変数
var BuiltIns = map[string]interface{}{
	"isPlainHostName":     isPlainHostName,
	"dnsDomainIs":         dnsDomainIs,
	"localHostOrDomainIs": localHostOrDomainIs,
	"isResolvable":        isResolvable,
	"isInNet":             isInNet,
	"dnsResolve":          dnsResolve,
	"convertAddr":         convertAddr,
	"myIPAddress":         myIPAddress,
	"dnsDomainLevels":     dnsDomainLevels,
	"shExpMatch":          shExpMatch,
	"weekdayRange":        weekdayRange,
	"dateRange":           dateRange,
	"timeRange":           timeRange,

	// *ADD HERE*
}

func addBuiltIn(name string, value interface{}) {
	BuiltIns[name] = value
}
