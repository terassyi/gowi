package decoder

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/terassyi/gowi/module"
	"github.com/terassyi/gowi/module/section"
	"github.com/terassyi/gowi/types"
)

const (
	WASM_EXT string = ".wasm"
)

var (
	InvalidFileFormat  error = errors.New("Given file is not .wasm.")
	InvalidMajicNumber error = errors.New("Invalid Majic Number.")
	InvalidWasmVersion error = errors.New("Invalid WASM version.")
)

type Decoder struct {
	path string
}

func New(path string) *Decoder {
	return &Decoder{
		path: path,
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

type sectionDecoder struct {
	id            uint8 // if custom section, id == 0
	payloadLength uint32
	nameLength    uint32 // present if id == 0;
	name          []byte // present if id == 0
	payloadData   []byte
}

func decodeSections(data []byte) ([]section.Section, error) {
	if err := validateMajicNumber(data); err != nil {
		return nil, fmt.Errorf("decodeSections: %w", err)
	}
	var sections []section.Section
	offset := 8
	for offset < len(data) {
		sd := sectionDecoder{}
		sd.id = data[offset]
		offset++
		p, n, err := types.DecodeVarUint32(bytes.NewBuffer(data[offset : offset+5]))
		sd.payloadLength = uint32(p)
		if err != nil {
			return nil, fmt.Errorf("decodeSections: %w", err)
		}
		offset += n
		if sd.id == 0 {
			p, n, err := types.DecodeVarUint32(bytes.NewBuffer(data[offset : offset+5]))
			if err != nil {
				return nil, fmt.Errorf("decodeSections: %w", err)
			}
			offset += n
			sd.nameLength = uint32(p)
			sd.name = data[offset : offset+int(sd.nameLength)]
			offset += int(sd.nameLength)
		}
		sd.payloadData = data[offset : offset+int(sd.payloadLength)]
		offset += int(sd.payloadLength)
		sec, err := section.New(sd.id, sd.payloadData)
		if err != nil {
			return nil, fmt.Errorf("decodeSections: %w", err)
		}
		sections = append(sections, sec)
	}
	return sections, nil
}

func validateMajicNumber(data []byte) error {
	majic := binary.BigEndian.Uint32(data[0:4])
	if majic != module.MAJIC_NUMBER {
		return fmt.Errorf("validateMajicNumber: %w", InvalidMajicNumber)
	}
	if binary.LittleEndian.Uint32(data[4:8]) != 0x1 {
		return fmt.Errorf("validateMajicNumber: %w", InvalidWasmVersion)
	}
	return nil
}
