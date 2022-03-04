package decoder

import (
	"bytes"
	"encoding/binary"
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

func (d *Decoder) Decode() (*module.Module, error) {
	m, err := decode(d.path)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func decode(path string) (*module.Module, error) {
	data, err := readWasmFile(path)
	if err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}
	buf := bytes.NewBuffer(data)
	if err := validateMajicNumber(buf); err != nil {
		return nil, fmt.Errorf("decode: majic_number: %w", err)
	}
	module := &module.Module{}
	version, err := decodeVersion(buf)
	if err != nil {
		return nil, fmt.Errorf("decode: version: %w", err)
	}
	module.Version = version
	for buf.Len() > 0 {
		sd, err := newSectionDecoder(buf)
		if err != nil {
			return nil, fmt.Errorf("decode: %w", err)
		}
		switch sd.id {
		case section.CUSTOM:
			s, err := section.NewCustom(sd.payloadData)
			if err != nil {
				return nil, fmt.Errorf("decode: %w", err)
			}
			module.Custom = s
		case section.TYPE:
			s, err := section.NewType(sd.payloadData)
			if err != nil {
				return nil, fmt.Errorf("decode: %w", err)
			}
			module.Type = s
		case section.IMPORT:
			s, err := section.NewImport(sd.payloadData)
			if err != nil {
				return nil, fmt.Errorf("decode: %w", err)
			}
			module.Import = s
		case section.FUNCTION:
			s, err := section.NewFunction(sd.payloadData)
			if err != nil {
				return nil, fmt.Errorf("decode: %w", err)
			}
			module.Function = s
		case section.TABLE:
			s, err := section.NewTable(sd.payloadData)
			if err != nil {
				return nil, fmt.Errorf("decode: %w", err)
			}
			module.Table = s
		case section.MEMORY:
			s, err := section.NewMemory(sd.payloadData)
			if err != nil {
				return nil, fmt.Errorf("decode: %w", err)
			}
			module.Memory = s
		case section.GLOBAL:
			s, err := section.NewGlobal(sd.payloadData)
			if err != nil {
				return nil, fmt.Errorf("decode: %w", err)
			}
			module.Global = s
		case section.EXPORT:
			s, err := section.NewExport(sd.payloadData)
			if err != nil {
				return nil, fmt.Errorf("decode: %w", err)
			}
			module.Export = s
		case section.START:
			s, err := section.NewStart(sd.payloadData)
			if err != nil {
				return nil, fmt.Errorf("decode: %w", err)
			}
			module.Start = s
		case section.ELEMENT:
			s, err := section.NewElement(sd.payloadData)
			if err != nil {
				return nil, fmt.Errorf("decode: %w", err)
			}
			module.Element = s
		case section.CODE:
			s, err := section.NewCode(sd.payloadData)
			if err != nil {
				return nil, fmt.Errorf("decode: %w", err)
			}
			module.Code = s
		case section.DATA:
			s, err := section.NewData(sd.payloadData)
			if err != nil {
				return nil, fmt.Errorf("decode: %w", err)
			}
			module.Data = s
		default:
			return nil, section.InvalidSectionCode
		}
	}
	return module, nil
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
	id            section.SectionCode // if custom section, id == 0
	payloadLength uint32
	nameLength    uint32 // present if id == 0;
	name          []byte // present if id == 0
	payloadData   []byte
}

func newSectionDecoder(buf *bytes.Buffer) (*sectionDecoder, error) {
	sd := &sectionDecoder{}
	id, err := buf.ReadByte()
	if err != nil {
		return nil, fmt.Errorf("newSectionDecoder: decode id: %w", err)
	}
	sectionId, err := section.NewSectionCode(id)
	if err != nil {
		return nil, fmt.Errorf("newSectionDecoder: decode id: %w", err)
	}
	sd.id = sectionId
	payloadLength, _, err := types.DecodeVarUint32(buf)
	if err != nil {
		return nil, fmt.Errorf("newSectionDecoder: decode payload_length: %w", err)
	}
	sd.payloadLength = uint32(payloadLength)
	if id == byte(0x00) {
		nameLength, _, err := types.DecodeVarUint32(buf)
		if err != nil {
			return nil, fmt.Errorf("newSectionDecoder: decode name_length: %w", err)
		}
		sd.nameLength = uint32(nameLength)
		name := make([]byte, int(nameLength))
		if _, err := buf.Read(name); err != nil {
			return nil, fmt.Errorf("newSectionDecoder: decode name: %w", err)
		}
		sd.name = name
	}
	data := make([]byte, int(sd.payloadLength))
	if _, err := buf.Read(data); err != nil {
		return nil, fmt.Errorf("newSectionDecoder: decode payload_data: %w", err)
	}
	sd.payloadData = data
	return sd, nil
}

func validateMajicNumber(buf *bytes.Buffer) error {
	b := make([]byte, 4)
	if _, err := buf.Read(b); err != nil {
		return fmt.Errorf("validateMajicNumber: read: %w", err)
	}
	majic := binary.BigEndian.Uint32(b)
	if majic != module.MAJIC_NUMBER {
		return fmt.Errorf("validateMajicNumber: %w", InvalidMajicNumber)
	}
	return nil
}

func decodeVersion(buf *bytes.Buffer) (uint32, error) {
	b := make([]byte, 4)
	if _, err := buf.Read(b); err != nil {
		return 0, fmt.Errorf("decodeVersion: read: %w", err)
	}
	return binary.LittleEndian.Uint32(b), nil
}

func HexDump(file string) ([]byte, error) {
	return readWasmFile(file)
}
