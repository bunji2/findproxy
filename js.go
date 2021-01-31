package main

import (
	"github.com/robertkrimen/otto"
)

// JSCtx JavaScript実行コンテクスト
type JSCtx struct {
	//
	Count    int
	filePath string
	vm       *otto.Otto
	script   *otto.Script
}

// NewJSCtx 新規JavaSript実行コンテクストの生成
func NewJSCtx(filePath string) (r *JSCtx, err error) {
	vm := otto.New()
	var script *otto.Script
	script, err = vm.Compile(filePath, nil)
	if err != nil {
		return
	}

	/*
		// バイト列をファイルに保存する組み込み関数
		addBuiltIn("save", func(filePath string, bb []byte) {
			saveFile(outdir, filePath, bb)
		})
	*/

	// 組み込み関数をJavaSript実行コンテクストに登録
	for name, value := range BuiltIns {
		err = vm.Set(name, value)
		if err != nil {
			return
		}
	}

	_, err = vm.Run(script)
	if err != nil {
		return
	}

	r = &JSCtx{
		filePath: filePath,
		vm:       vm,
		script:   script,
	}
	return
}

// FindProxyForURL は与えられたURLとホスト名への接続に使用すべきプロキシを返す関数
func (ctx *JSCtx) FindProxyForURL(url, host string) (r string) {
	value, err := ctx.vm.Call("FindProxyForURL", nil, url, host)
	/*
		if err != nil && err.Error()[0:15] == "ReferenceError:" {
			err = nil
		}
	*/
	if err != nil {
		return
	}

	if value.IsString() {
		r = value.String()
	}

	return
}
