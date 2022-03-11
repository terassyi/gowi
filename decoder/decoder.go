package decoder

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/terassyi/gowi/structure"
	"github.com/terassyi/gowi/types"
)

const (
	WASM_EXT string = ".wasm"
)

const (
	MAJIC_NUMBER uint32 = 0x0061736d // \0asm
	WASM_VERSION uint32 = 0x1
)

var (
	InvalidFileFormat  error = errors.New("Given file is not .wasm.")
	InvalidMajicNumber error = errors.New("Invalid Majic Number.")
	InvalidWasmVersion error = errors.New("Invalid WASM version.")
)

type Decoder struct {
	path string
	mod  *mod
}

func New(path string) (*Decoder, error) {
	m, err := decode(path)
	if err != nil {
		return nil, fmt.Errorf("Decoder new: %w", err)
	}
	return &Decoder{
		path: path,
		mod:  m,
	}, nil
}

func (d *Decoder) Decode() (*structure.Module, error) {
	m, err := decode(d.path)
	if err != nil {
		return nil, err
	}
	d.mod = m
	return m.build()
}

func decode(path string) (*mod, error) {
	data, err := readWasmFile(path)
	if err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}
	buf := bytes.NewBuffer(data)
	if err := validateMajicNumber(buf); err != nil {
		return nil, fmt.Errorf("decode: majic_number: %w", err)
	}
	module := &mod{}
	version, err := decodeVersion(buf)
	if err != nil {
		return nil, fmt.Errorf("decode: version: %w", err)
	}
	module.version = version
	for buf.Len() > 0 {
		sd, err := newSectionDecoder(buf)
		if err != nil {
			return nil, fmt.Errorf("decode: %w", err)
		}
		switch sd.id {
		case CUSTOM:
			s, err := newCustom(sd.payloadData)
			if err != nil {
				return nil, fmt.Errorf("decode: %w", err)
			}
			module.custom = s
		case TYPE:
			s, err := newType(sd.payloadData)
			if err != nil {
				return nil, fmt.Errorf("decode: %w", err)
			}
			module.typ = s
		case IMPORT:
			s, err := newImport(sd.payloadData)
			if err != nil {
				return nil, fmt.Errorf("decode: %w", err)
			}
			module.imports = s
		case FUNCTION:
			s, err := newFunction(sd.payloadData)
			if err != nil {
				return nil, fmt.Errorf("decode: %w", err)
			}
			module.function = s
		case TABLE:
			s, err := newTable(sd.payloadData)
			if err != nil {
				return nil, fmt.Errorf("decode: %w", err)
			}
			module.table = s
		case MEMORY:
			s, err := newMemory(sd.payloadData)
			if err != nil {
				return nil, fmt.Errorf("decode: %w", err)
			}
			module.memory = s
		case GLOBAL:
			s, err := newGlobal(sd.payloadData)
			if err != nil {
				return nil, fmt.Errorf("decode: %w", err)
			}
			module.global = s
		case EXPORT:
			s, err := newExport(sd.payloadData)
			if err != nil {
				return nil, fmt.Errorf("decode: %w", err)
			}
			module.export = s
		case START:
			s, err := newStart(sd.payloadData)
			if err != nil {
				return nil, fmt.Errorf("decode: %w", err)
			}
			module.start = s
		case ELEMENT:
			s, err := newElement(sd.payloadData)
			if err != nil {
				return nil, fmt.Errorf("decode: %w", err)
			}
			module.element = s
		case CODE:
			s, err := newCode(sd.payloadData)
			if err != nil {
				return nil, fmt.Errorf("decode: %w", err)
			}
			module.code = s
		case DATA:
			s, err := newData(sd.payloadData)
			if err != nil {
				return nil, fmt.Errorf("decode: %w", err)
			}
			module.data = s
		default:
			return nil, InvalidSectionCode
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
	id            SectionCode // if custom section, id == 0
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
	sectionId, err := newSectionCode(id)
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
	if majic != MAJIC_NUMBER {
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

func (d *Decoder) dumpVersion() string {
	return fmt.Sprintf("%s: file format wasm 0x%x\n", d.path, d.mod.version)
}

func (d *Decoder) DumpSection() string {
	str := d.dumpVersion()
	str += "Sections:\n\n"
	if d.mod.custom != nil {
		str += fmt.Sprintf("Custom is not implemented.\n")
	}
	if d.mod.typ != nil {
		str += fmt.Sprintf("Type : count=0x%04x\n", len(d.mod.typ.entries))
	}
	if d.mod.imports != nil {
		str += fmt.Sprintf("Import: count 0x%04x\n", len(d.mod.imports.entries))
	}
	if d.mod.function != nil {
		str += fmt.Sprintf("Function : count=0x%04x\n", len(d.mod.function.types))
	}
	if d.mod.table != nil {
		str += fmt.Sprintf("Table : count 0x%04x\n", len(d.mod.table.entries))
	}
	if d.mod.memory != nil {
		str += fmt.Sprintf("Memory : count 0x%04x\n", len(d.mod.memory.entries))
	}
	if d.mod.global != nil {
		str += fmt.Sprintf("Global : count 0x%04x\n", len(d.mod.global.globals))
	}
	if d.mod.export != nil {
		str += fmt.Sprintf("Export : count 0x%04x\n", len(d.mod.export.entries))
	}
	if d.mod.start != nil {
		str += fmt.Sprintf("Start: index %d\n", d.mod.start.index)
	}
	if d.mod.element != nil {
		str += fmt.Sprintf("Element : count 0x%04x\n", len(d.mod.element.entries))
	}
	if d.mod.code != nil {
		str += fmt.Sprintf("Code : count 0x%04x\n", len(d.mod.code.bodies))
	}
	if d.mod.data != nil {
		str += fmt.Sprintf("Data : count 0x%04x\n", len(d.mod.data.entries))
	}
	return str
}

func (d *Decoder) DumpDetail() (string, error) {
	str := d.dumpVersion()
	str += "Section Details:\n\n"
	if d.mod.custom != nil {
		str += d.mod.custom.detail()
		str += "\n"
	}
	if d.mod.typ != nil {
		str += d.mod.typ.detail()
		str += "\n"
	}
	if d.mod.imports != nil {
		str += d.mod.imports.detail()
		str += "\n"
	}
	if d.mod.function != nil {
		str += d.mod.function.detail()
		str += "\n"
	}
	if d.mod.table != nil {
		str += d.mod.table.detail()
		str += "\n"
	}
	if d.mod.memory != nil {
		str += d.mod.memory.detail()
		str += "\n"
	}
	if d.mod.global != nil {
		s, err := d.mod.global.detail()
		if err != nil {
			return "", err
		}
		str += s
		str += "\n"
	}
	if d.mod.export != nil {
		str += d.mod.export.detail()
		str += "\n"
	}
	if d.mod.start != nil {
		str += d.mod.start.detail()
		str += "\n"
	}
	if d.mod.element != nil {
		s, err := d.mod.element.detail()
		if err != nil {
			return "", err
		}
		str += s
		str += "\n"
	}
	if d.mod.code != nil {
		str += d.mod.code.detail()
		str += "\n"
	}
	if d.mod.data != nil {
		s, err := d.mod.data.detail()
		if err != nil {
			return "", err
		}
		str += s
		str += "\n"
	}
	return str, nil
}
