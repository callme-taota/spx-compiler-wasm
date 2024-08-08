//go:build wasm
// +build wasm

package main

import "syscall/js"

func main() {
	c := make(chan struct{})
	// register js function
	jsFuncRegister()
	<-c
}

type jsFuncName string
type jsFuncBody func(this js.Value, p []js.Value) interface{}
type jsFuncList map[jsFuncName]jsFuncBody

var jsFuncs = jsFuncList{
	"getInlayHints":      getInlayHints,
	"getDiagnostics":     getDiagnostics,
	"getCompletionItems": getCompletionItems,
	"getDefinition":      getDefinition,
	"getTypes":           getTypes,
}

func jsFuncRegister() {
	for key, f := range jsFuncs {
		js.Global().Set(string(key), js.FuncOf(f))
	}
}
