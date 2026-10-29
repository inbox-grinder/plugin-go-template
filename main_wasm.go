//go:build wasm

package main

import "unsafe"

// inputBuf holds the bytes allocated by alloc() for the host to write into.
var inputBuf []byte

// outputBuf retains the last output from process() until the host reads it.
var outputBuf []byte

func main() {}

// alloc allocates size bytes and returns a pointer into WASM linear memory
// for the host to write the input JSON into.
//
//export alloc
func alloc(size int32) int32 {
	inputBuf = make([]byte, size)
	if len(inputBuf) == 0 {
		return 0
	}
	return int32(uintptr(unsafe.Pointer(&inputBuf[0])))
}

// process reads the input JSON written at ptr/length, calls run(), and returns
// the output as a packed int64: (output_ptr << 32) | output_len.
// Returns 0 on error.
//
//export process
func process(ptr int32, length int32) int64 {
	input := inputBuf[:length]
	output, err := run(input)
	if err != nil {
		return 0
	}
	outputBuf = output
	if len(outputBuf) == 0 {
		return 0
	}
	outPtr := int64(uintptr(unsafe.Pointer(&outputBuf[0])))
	outLen := int64(len(outputBuf))
	return (outPtr << 32) | outLen
}
