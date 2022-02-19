package main

import (
	"github.com/terassyi/gowi/decoder"
)

func main() {
	d := decoder.New("examples/func1.wasm", 0)
	d.Decode()
}
