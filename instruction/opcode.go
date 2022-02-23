package instruction

type Opcode uint8

const (
	// Control flow operators
	UNREACHABLE Opcode = 0x00
	NOP         Opcode = 0x01
	BLOCK       Opcode = 0x02
	LOOP        Opcode = 0x03
	IF          Opcode = 0x04
	ELSE        Opcode = 0x05
	END         Opcode = 0x0b
	BR          Opcode = 0x0c
	BR_IF       Opcode = 0x0d
	BR_TABLE    Opcode = 0x0e
	RETURN      Opcode = 0x0f

	// Call operators
	CALL          Opcode = 0x10
	CALL_INDIRECT Opcode = 0x11

	// Parametic operators
	DROP   Opcode = 0x1a
	SELECT Opcode = 0x1b

	// Variable access
	GET_LOCAL  Opcode = 0x20
	SET_LOCAL  Opcode = 0x21
	TEE_LOCAL  Opcode = 0x22
	GET_GLOBAL Opcode = 0x23
	SET_GLOBAL Opcode = 0x24

	// Memory-related operators
	I32_LOAD       Opcode = 0x28
	I64_LOAD       Opcode = 0x29
	F32_LOAD       Opcode = 0x2a
	F64_LOAD       Opcode = 0x2b
	I32_LOAD8_S    Opcode = 0x2c
	I32_LOAD8_U    Opcode = 0x2d
	I32_LOAD16_S   Opcode = 0x2e
	I32_LOAD16_U   Opcode = 0x2f
	I64_LOAD8_S    Opcode = 0x30
	I64_LOAD8_U    Opcode = 0x31
	I64_LOAD16_S   Opcode = 0x32
	I64_LOAD16_U   Opcode = 0x33
	I64_LOAD32_S   Opcode = 0x34
	I64_LOAD32_U   Opcode = 0x35
	I32_STORE      Opcode = 0x36
	I64_STORE      Opcode = 0x37
	F32_STORE      Opcode = 0x38
	F64_STORE      Opcode = 0x39
	I32_STORE8     Opcode = 0x3a
	I32_STORE16    Opcode = 0x3b
	I64_STORE8     Opcode = 0x3c
	I64_STORE16    Opcode = 0x3d
	I64_STORE32    Opcode = 0x3e
	CURRENT_MEMORY Opcode = 0x3f
	GROW_MEMORY    Opcode = 0x40

	// Constants
	I32_CONST Opcode = 0x41
	I64_CONST Opcode = 0x42
	F32_CONST Opcode = 0x43
	F64_CONST Opcode = 0x44

	// Comparison operators
	I32_EQZ  Opcode = 0x45
	I32_EQ   Opcode = 0x46
	I32_NE   Opcode = 0x47
	I32_LT_S Opcode = 0x48
	I32_LT_U Opcode = 0x49
	I32_GT_S Opcode = 0x4a
	I32_GT_U Opcode = 0x4b
	I32_LE_S Opcode = 0x4c
	I32_LE_U Opcode = 0x4d
	I32_GE_S Opcode = 0x4e
	I32_GE_U Opcode = 0x4f
	I64_EQZ  Opcode = 0x50
	I64_EQ   Opcode = 0x51
	I64_NE   Opcode = 0x52
	I64_LT_S Opcode = 0x53
	I64_LT_U Opcode = 0x54
	I64_GT_S Opcode = 0x55
	I64_GT_U Opcode = 0x56
	I64_LE_S Opcode = 0x57
	I64_LE_U Opcode = 0x58
	I64_GE_S Opcode = 0x59
	I64_GE_U Opcode = 0x5a
	F32_EQ   Opcode = 0x5b
	F32_NE   Opcode = 0x5c
	F32_LT   Opcode = 0x5d
	F32_GT   Opcode = 0x5e
	F32_LE   Opcode = 0x5f
	F32_GE   Opcode = 0x60
	F64_EQ   Opcode = 0x61
	F64_NE   Opcode = 0x62
	F64_LT   Opcode = 0x63
	F64_GT   Opcode = 0x64
	F64_LE   Opcode = 0x65
	F64_GE   Opcode = 0x66

	// Numeric operators
	I32_CLZ      Opcode = 0x67
	I32_CTZ      Opcode = 0x68
	I32_POPCNT   Opcode = 0x69
	I32_ADD      Opcode = 0x6a
	I32_SUB      Opcode = 0x6b
	I32_MUL      Opcode = 0x6c
	I32_DIV_S    Opcode = 0x6d
	I32_DIV_U    Opcode = 0x6e
	I32_REM_S    Opcode = 0x6f
	I32_REM_U    Opcode = 0x70
	I32_AND      Opcode = 0x71
	I32_OR       Opcode = 0x72
	I32_XOR      Opcode = 0x73
	I32_SHL      Opcode = 0x74
	I32_SHR_S    Opcode = 0x75
	I32_SHR_U    Opcode = 0x76
	I32_ROTL     Opcode = 0x77
	I32_ROTR     Opcode = 0x78
	I64_CLZ      Opcode = 0x79
	I64_CTZ      Opcode = 0x7a
	I64_POPCNT   Opcode = 0x7b
	I64_ADD      Opcode = 0x7c
	I64_SUB      Opcode = 0x7d
	I64_MUL      Opcode = 0x7e
	I64_DIV_S    Opcode = 0x7f
	I64_DIV_U    Opcode = 0x80
	I64_REM_S    Opcode = 0x81
	I64_REM_U    Opcode = 0x82
	I64_AND      Opcode = 0x83
	I64_OR       Opcode = 0x84
	I64_XOR      Opcode = 0x85
	I64_SHL      Opcode = 0x86
	I64_SHR_S    Opcode = 0x87
	I64_SHR_U    Opcode = 0x88
	I64_ROTL     Opcode = 0x89
	I64_ROTR     Opcode = 0x8a
	F32_ABS      Opcode = 0x8b
	F32_NEG      Opcode = 0x8c
	F32_CEIL     Opcode = 0x8d
	F32_FLOOR    Opcode = 0x8e
	F32_TRUNC    Opcode = 0x8f
	F32_NEAREST  Opcode = 0x90
	F32_SQRT     Opcode = 0x91
	F32_ADD      Opcode = 0x92
	F32_SUB      Opcode = 0x93
	F32_MUL      Opcode = 0x94
	F32_DIV      Opcode = 0x95
	F32_MIN      Opcode = 0x96
	F32_MAX      Opcode = 0x97
	F32_COPYSIGN Opcode = 0x98
	F64_ABS      Opcode = 0x99
	F64_NEG      Opcode = 0x9a
	F64_CEIL     Opcode = 0x9b
	F64_FLOOR    Opcode = 0x9c
	F64_TRUNC    Opcode = 0x9d
	F64_NEAREST  Opcode = 0x9e
	F64_SQRT     Opcode = 0x9f
	F64_ADD      Opcode = 0xa0
	F64_SUB      Opcode = 0xa1
	F64_MUL      Opcode = 0xa2
	F64_DIV      Opcode = 0xa3
	F64_MIN      Opcode = 0xa4
	F64_MAX      Opcode = 0xa5
	F64_COPYSIGN Opcode = 0xa6

	// Conversions
	I32_WRAP_I64        Opcode = 0xa7
	I32_TRUNC_S_F32     Opcode = 0xa8
	I32_TRUNC_U_F32     Opcode = 0xa9
	I32_TRUNC_S_F64     Opcode = 0xaa
	I32_TRUNC_U_F64     Opcode = 0xab
	I64_EXTEND_S_I32    Opcode = 0xac
	I64_EXTEND_U_I32    Opcode = 0xad
	I64_TRUNC_S_F32     Opcode = 0xae
	I64_TRUNC_U_F32     Opcode = 0xaf
	I64_TRUNC_S_F64     Opcode = 0xb0
	I64_TRUNC_U_F64     Opcode = 0xb1
	F32_CONVERT_S_I32   Opcode = 0xb2
	F32_CONVERT_U_I32   Opcode = 0xb3
	F32_CONVERT_S_I64   Opcode = 0xb4
	F32_CONVERT_U_I64   Opcode = 0xb5
	F32_DEMOTE_F64      Opcode = 0xb6
	F64_CONVERT_S_I32   Opcode = 0xb7
	F64_CONVERT_U_I32   Opcode = 0xb8
	F64_CONVERT_S_I64   Opcode = 0xb9
	F64_CONVERT_U_I64   Opcode = 0xba
	F64_PROMOTE_F32     Opcode = 0xbb
	I32_TRUNC_SAT_F32_S Opcode = 0xfc // 0x00
	I32_TRUNC_SAT_F32_U Opcode = 0xfc // 0x01
	I32_TRUNC_SAT_F64_S Opcode = 0xfc // 0x02
	I32_TRUNC_SAT_F64_U Opcode = 0xfc // 0x03
	I64_TRUNC_SAT_F32_S Opcode = 0xfc // 0x04
	I64_TRUNC_SAT_F32_U Opcode = 0xfc // 0x05
	I64_TRUNC_SAT_F64_S Opcode = 0xfc // 0x06
	I64_TRUNC_SAT_F64_U Opcode = 0xfc // 0x07

	// Reinterpretations
	I32_REINTERPRET_F32 Opcode = 0xbc
	I64_REINTERPRET_F64 Opcode = 0xbd
	F32_REINTERPRET_I32 Opcode = 0xbe
	F64_REINTERPRET_I64 Opcode = 0xbf
)
