package decoder

import (
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

const (
	WASM_EXT string = ".wasm"
)

var (
	InvalidFileFormat error = errors.New("Given file is not .wasm.")
)

type Decoder struct {
	path  string
	flags int
}

func New(path string, flags int) *Decoder {
	return &Decoder{
		path:  path,
		flags: flags,
	}
}

func (d *Decoder) Decode() (any, error) {
	if err := decode(d.path); err != nil {
		return nil, err
	}
	return nil, nil
}

func decode(path string) error {
	data, err := readWasmFile(path)
	if err != nil {
		return fmt.Errorf("decode: %w", err)
	}
	fmt.Println(hex.Dump(data))
	return nil
}

func validateExt(path string) error {
	if filepath.Ext(path) != WASM_EXT {
		return InvalidFileFormat
	}
	return nil
}

func readWasmFile(path string) ([]byte, error) {
	if err := validateExt(path); err != nil {
		return nil, fmt.Errorf("readWasmFile: %w", err)
	}
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("readWasmFile: %w", err)
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("readWasmFile: %w", err)
	}
	return data, nil
}
