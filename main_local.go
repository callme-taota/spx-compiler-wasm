//go:build darwin
// +build darwin

package main

import "spx-compiler-wasm/internal/spx"

func main() {
	spx.StartSPXTypesAnalyser("test.spx", `
onStart => {
	flag := true
	for flag {
		onMsg "die", => {
			flag = false
		}
		glide -877, 180, 3
		setXYpos -240, 180
	}
}
`)
}
