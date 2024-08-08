package main

import (
	"spx-compiler-wasm/internal/spx"
	"syscall/js"
)

func getInlayHints(this js.Value, p []js.Value) interface{} {
	return nil
}

func getDiagnostics(this js.Value, p []js.Value) interface{} {
	return nil
}

func getCompletionItems(this js.Value, p []js.Value) interface{} {
	return nil
}

func getDefinition(this js.Value, p []js.Value) interface{} {
	return nil
}

func getTypes(this js.Value, p []js.Value) interface{} {
	fileName := p[0].String()
	fileCode := p[1].String()
	return spx.StartSPXTypesAnalyser(fileName, fileCode)

}
