package main

import (
	"fmt"
	"net/url"
)

func process(scriptFilePath string, urls []string) (err error) {

	var ctx *JSCtx
	ctx, err = NewJSCtx(scriptFilePath)
	if err != nil {
		return
	}

	for _, urlStr := range urls {
		u, e := url.Parse(urlStr)
		if e != nil {
			err = e
			break
		}
		fmt.Println(urlStr, "=>", ctx.FindProxyForURL(urlStr, u.Hostname()))
	}

	return
}
