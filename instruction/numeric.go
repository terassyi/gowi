package instruction

type I32Add struct{}

func (*I32Add) Opcode() Opcode {
	return I32_ADD
}

func (*I32Add) imm() any {
	return NoImm
}

func (*I32Add) String() string {
	return "i32.add"
}

func (*I32Add) ImmString() string {
	return ""
}

type I32Sub struct{}

func (*I32Sub) Opcode() Opcode {
	return I32_SUB
}

func (*I32Sub) imm() any {
	return NoImm
}

func (*I32Sub) String() string {
	return "i32.sub"
}

func (*I32Sub) ImmString() string {
	return ""
}

type I32Mul struct{}

func (*I32Mul) Opcode() Opcode {
	return I32_MUL
}

func (*I32Mul) imm() any {
	return NoImm
}

func (*I32Mul) String() string {
	return "i32.mul"
}

func (*I32Mul) ImmString() string {
	return ""
}

type I32DivS struct{}

func (*I32DivS) Opcode() Opcode {
	return I32_DIV_S
}

func (*I32DivS) imm() any {
	return NoImm
}

func (*I32DivS) String() string {
	return "i32.div_s"
}

func (*I32DivS) ImmString() string {
	return ""
}

type I32DivU struct{}

func (*I32DivU) Opcode() Opcode {
	return I32_DIV_U
}

func (*I32DivU) imm() any {
	return NoImm
}

func (*I32DivU) String() string {
	return "i32.div_u"
}

func (*I32DivU) ImmString() string {
	return ""
}

type I64Add struct{}

func (*I64Add) Opcode() Opcode {
	return I64_ADD
}

func (*I64Add) imm() any {
	return NoImm
}

func (*I64Add) String() string {
	return "i64.add"
}

func (*I64Add) ImmString() string {
	return ""
}

type I64Sub struct{}

func (*I64Sub) Opcode() Opcode {
	return I64_SUB
}

func (*I64Sub) imm() any {
	return NoImm
}

func (*I64Sub) String() string {
	return "i64.sub"
}

func (*I64Sub) ImmString() string {
	return ""
}

type I64Mul struct{}

func (*I64Mul) Opcode() Opcode {
	return I64_MUL
}

func (*I64Mul) imm() any {
	return NoImm
}

func (*I64Mul) String() string {
	return "i64.mul"
}

func (*I64Mul) ImmString() string {
	return ""
}

type I64DivS struct{}

func (*I64DivS) Opcode() Opcode {
	return I64_DIV_S
}

func (*I64DivS) imm() any {
	return NoImm
}

func (*I64DivS) String() string {
	return "i64.div_s"
}

func (*I64DivS) ImmString() string {
	return ""
}

type I64DivU struct{}

func (*I64DivU) Opcode() Opcode {
	return I64_DIV_U
}

func (*I64DivU) imm() any {
	return NoImm
}

func (*I64DivU) String() string {
	return "i64.div_u"
}

func (*I64DivU) ImmString() string {
	return ""
}

type I32RemS struct{}

func (*I32RemS) Opcode() Opcode {
	return I32_REM_S
}

func (*I32RemS) imm() any {
	return NoImm
}

func (*I32RemS) String() string {
	return "i32.rem_s"
}

func (*I32RemS) ImmString() string {
	return ""
}

type I32RemU struct{}

func (*I32RemU) Opcode() Opcode {
	return I32_REM_U
}

func (*I32RemU) imm() any {
	return NoImm
}

func (*I32RemU) String() string {
	return "i32.rem_u"
}

func (*I32RemU) ImmString() string {
	return ""
}

type I64RemS struct{}

func (*I64RemS) Opcode() Opcode {
	return I64_REM_S
}

func (*I64RemS) imm() any {
	return NoImm
}

func (*I64RemS) String() string {
	return "i64.rem_s"
}

func (*I64RemS) ImmString() string {
	return ""
}

type I64RemU struct{}

func (*I64RemU) Opcode() Opcode {
	return I64_REM_U
}

func (*I64RemU) imm() any {
	return NoImm
}

func (*I64RemU) String() string {
	return "i64.rem_u"
}

func (*I64RemU) ImmString() string {
	return ""
}

type I32And struct{}

func (*I32And) Opcode() Opcode {
	return I32_AND
}

func (*I32And) imm() any {
	return NoImm
}

func (*I32And) String() string {
	return "i32.and"
}

func (*I32And) ImmString() string {
	return ""
}

type I64And struct{}

func (*I64And) Opcode() Opcode {
	return I64_AND
}

func (*I64And) imm() any {
	return NoImm
}

func (*I64And) String() string {
	return "i64.and"
}

func (*I64And) ImmString() string {
	return ""
}

type I32Or struct{}

func (*I32Or) Opcode() Opcode {
	return I32_OR
}

func (*I32Or) imm() any {
	return NoImm
}

func (*I32Or) String() string {
	return "i32.or"
}

func (*I32Or) ImmString() string {
	return ""
}

type I64Or struct{}

func (*I64Or) Opcode() Opcode {
	return I64_OR
}

func (*I64Or) imm() any {
	return NoImm
}

func (*I64Or) String() string {
	return "i64.or"
}

func (*I64Or) ImmString() string {
	return ""
}

type I32Xor struct{}

func (*I32Xor) Opcode() Opcode {
	return I32_XOR
}

func (*I32Xor) imm() any {
	return NoImm
}

func (*I32Xor) String() string {
	return "i32.xor"
}

func (*I32Xor) ImmString() string {
	return ""
}

type I64Xor struct{}

func (*I64Xor) Opcode() Opcode {
	return I64_XOR
}

func (*I64Xor) imm() any {
	return NoImm
}

func (*I64Xor) String() string {
	return "i64.xor"
}

func (*I64Xor) ImmString() string {
	return ""
}

type I32Shl struct{}

func (*I32Shl) Opcode() Opcode {
	return I32_SHL
}

func (*I32Shl) imm() any {
	return NoImm
}

func (*I32Shl) String() string {
	return "i32.shl"
}

func (*I32Shl) ImmString() string {
	return ""
}

type I64Shl struct{}

func (*I64Shl) Opcode() Opcode {
	return I64_SHL
}

func (*I64Shl) imm() any {
	return NoImm
}

func (*I64Shl) String() string {
	return "i64.shl"
}

func (*I64Shl) ImmString() string {
	return ""
}

type I32ShrU struct{}

func (*I32ShrU) Opcode() Opcode {
	return I32_SHR_U
}

func (*I32ShrU) imm() any {
	return NoImm
}

func (*I32ShrU) String() string {
	return "i32.shr_u"
}

func (*I32ShrU) ImmString() string {
	return ""
}

type I64ShrU struct{}

func (*I64ShrU) Opcode() Opcode {
	return I64_SHR_U
}

func (*I64ShrU) imm() any {
	return NoImm
}

func (*I64ShrU) String() string {
	return "i64.shr_u"
}

func (*I64ShrU) ImmString() string {
	return ""
}

type I32ShrS struct{}

func (*I32ShrS) Opcode() Opcode {
	return I32_SHR_S
}

func (*I32ShrS) imm() any {
	return NoImm
}

func (*I32ShrS) String() string {
	return "i32.shr_s"
}

func (*I32ShrS) ImmString() string {
	return ""
}

type I64ShrS struct{}

func (*I64ShrS) Opcode() Opcode {
	return I64_SHR_S
}

func (*I64ShrS) imm() any {
	return NoImm
}

func (*I64ShrS) String() string {
	return "i64.shr_s"
}

func (*I64ShrS) ImmString() string {
	return ""
}

type I32RotR struct{}

func (*I32RotR) Opcode() Opcode {
	return I32_ROTR
}

func (*I32RotR) imm() any {
	return NoImm
}

func (*I32RotR) String() string {
	return "i32.rotr"
}

func (*I32RotR) ImmString() string {
	return ""
}

type I64RotR struct{}

func (*I64RotR) Opcode() Opcode {
	return I64_ROTR
}

func (*I64RotR) imm() any {
	return NoImm
}

func (*I64RotR) String() string {
	return "i64.rotr"
}

func (*I64RotR) ImmString() string {
	return ""
}

type I32RotL struct{}

func (*I32RotL) Opcode() Opcode {
	return I32_ROTL
}

func (*I32RotL) imm() any {
	return NoImm
}

func (*I32RotL) String() string {
	return "i32.rotl"
}

func (*I32RotL) ImmString() string {
	return ""
}

type I64RotL struct{}

func (*I64RotL) Opcode() Opcode {
	return I64_ROTL
}

func (*I64RotL) imm() any {
	return NoImm
}

func (*I64RotL) String() string {
	return "i64.rotl"
}

func (*I64RotL) ImmString() string {
	return ""
}

type I32Eq struct{}

func (*I32Eq) Opcode() Opcode {
	return I32_EQ
}

func (*I32Eq) imm() any {
	return NoImm
}

func (*I32Eq) String() string {
	return "i32.eq"
}

func (*I32Eq) ImmString() string {
	return ""
}

type I64Eq struct{}

func (*I64Eq) Opcode() Opcode {
	return I64_EQ
}

func (*I64Eq) imm() any {
	return NoImm
}

func (*I64Eq) String() string {
	return "i64.eq"
}

func (*I64Eq) ImmString() string {
	return ""
}

type I32Ne struct{}

func (*I32Ne) Opcode() Opcode {
	return I32_NE
}

func (*I32Ne) imm() any {
	return NoImm
}

func (*I32Ne) String() string {
	return "i32.ne"
}

func (*I32Ne) ImmString() string {
	return ""
}

type I64Ne struct{}

func (*I64Ne) Opcode() Opcode {
	return I64_NE
}

func (*I64Ne) imm() any {
	return NoImm
}

func (*I64Ne) String() string {
	return "i64.ne"
}

func (*I64Ne) ImmString() string {
	return ""
}

type I32LtU struct{}

func (*I32LtU) Opcode() Opcode {
	return I32_LT_U
}

func (*I32LtU) imm() any {
	return NoImm
}

func (*I32LtU) String() string {
	return "i32.lt_u"
}

func (*I32LtU) ImmString() string {
	return ""
}

type I64LtU struct{}

func (*I64LtU) Opcode() Opcode {
	return I64_LT_U
}

func (*I64LtU) imm() any {
	return NoImm
}

func (*I64LtU) String() string {
	return "i64.lt_u"
}

func (*I64LtU) ImmString() string {
	return ""
}

type I32LtS struct{}

func (*I32LtS) Opcode() Opcode {
	return I32_LT_S
}

func (*I32LtS) imm() any {
	return NoImm
}

func (*I32LtS) String() string {
	return "i32.lt_s"
}

func (*I32LtS) ImmString() string {
	return ""
}

type I64LtS struct{}

func (*I64LtS) Opcode() Opcode {
	return I64_LT_S
}

func (*I64LtS) imm() any {
	return NoImm
}

func (*I64LtS) String() string {
	return "i64.lt_s"
}

func (*I64LtS) ImmString() string {
	return ""
}

type I32GtS struct{}

func (*I32GtS) Opcode() Opcode {
	return I32_GT_S
}

func (*I32GtS) imm() any {
	return NoImm
}

func (*I32GtS) String() string {
	return "i32.gt_s"
}

func (*I32GtS) ImmString() string {
	return ""
}

type I32GtU struct{}

func (*I32GtU) Opcode() Opcode {
	return I32_GT_U
}

func (*I32GtU) imm() any {
	return NoImm
}

func (*I32GtU) String() string {
	return "i32.gt_u"
}

func (*I32GtU) ImmString() string {
	return ""
}

type I64GtS struct{}

func (*I64GtS) Opcode() Opcode {
	return I64_GT_S
}

func (*I64GtS) imm() any {
	return NoImm
}

func (*I64GtS) String() string {
	return "i64.gt_s"
}

func (*I64GtS) ImmString() string {
	return ""
}

type I64GtU struct{}

func (*I64GtU) Opcode() Opcode {
	return I64_GT_U
}

func (*I64GtU) imm() any {
	return NoImm
}

func (*I64GtU) String() string {
	return "i64.gt_u"
}

func (*I64GtU) ImmString() string {
	return ""
}

type I32LeU struct{}

func (*I32LeU) Opcode() Opcode {
	return I32_LE_U
}

func (*I32LeU) imm() any {
	return NoImm
}

func (*I32LeU) String() string {
	return "i32.le_u"
}

func (*I32LeU) ImmString() string {
	return ""
}

type I64LeU struct{}

func (*I64LeU) Opcode() Opcode {
	return I64_LE_U
}

func (*I64LeU) imm() any {
	return NoImm
}

func (*I64LeU) String() string {
	return "i64.le_u"
}

func (*I64LeU) ImmString() string {
	return ""
}

type I32LeS struct{}

func (*I32LeS) Opcode() Opcode {
	return I32_LE_S
}

func (*I32LeS) imm() any {
	return NoImm
}

func (*I32LeS) String() string {
	return "i32.le_s"
}

func (*I32LeS) ImmString() string {
	return ""
}

type I64LeS struct{}

func (*I64LeS) Opcode() Opcode {
	return I64_LE_S
}

func (*I64LeS) imm() any {
	return NoImm
}

func (*I64LeS) String() string {
	return "i64.le_s"
}

func (*I64LeS) ImmString() string {
	return ""
}

type I32GeS struct{}

func (*I32GeS) Opcode() Opcode {
	return I32_GE_S
}

func (*I32GeS) imm() any {
	return NoImm
}

func (*I32GeS) String() string {
	return "i32.ge_s"
}

func (*I32GeS) ImmString() string {
	return ""
}

type I32GeU struct{}

func (*I32GeU) Opcode() Opcode {
	return I32_GE_U
}

func (*I32GeU) imm() any {
	return NoImm
}

func (*I32GeU) String() string {
	return "i32.ge_u"
}

func (*I32GeU) ImmString() string {
	return ""
}

type I64GeS struct{}

func (*I64GeS) Opcode() Opcode {
	return I64_GE_S
}

func (*I64GeS) imm() any {
	return NoImm
}

func (*I64GeS) String() string {
	return "i64.ge_s"
}

func (*I64GeS) ImmString() string {
	return ""
}

type I64GeU struct{}

func (*I64GeU) Opcode() Opcode {
	return I64_GE_U
}

func (*I64GeU) imm() any {
	return NoImm
}

func (*I64GeU) String() string {
	return "i64.ge_u"
}

func (*I64GeU) ImmString() string {
	return ""
}

type I32Eqz struct{}

func (*I32Eqz) Opcode() Opcode {
	return I32_EQZ
}

func (*I32Eqz) imm() any {
	return NoImm
}

func (*I32Eqz) String() string {
	return "i32.eqz"
}

func (*I32Eqz) ImmString() string {
	return ""
}

type I64Eqz struct{}

func (*I64Eqz) Opcode() Opcode {
	return I64_EQZ
}

func (*I64Eqz) imm() any {
	return NoImm
}

func (*I64Eqz) String() string {
	return "i64.eqz"
}

func (*I64Eqz) ImmString() string {
	return ""
}

type I32Clz struct{}

func (*I32Clz) Opcode() Opcode {
	return I32_CLZ
}

func (*I32Clz) imm() any {
	return NoImm
}

func (*I32Clz) String() string {
	return "i32.clz"
}

func (*I32Clz) ImmString() string {
	return ""
}

type I64Clz struct{}

func (*I64Clz) Opcode() Opcode {
	return I64_CLZ
}

func (*I64Clz) imm() any {
	return NoImm
}

func (*I64Clz) String() string {
	return "i64.clz"
}

func (*I64Clz) ImmString() string {
	return ""
}

type I32Ctz struct{}

func (*I32Ctz) Opcode() Opcode {
	return I32_CTZ
}

func (*I32Ctz) imm() any {
	return NoImm
}

func (*I32Ctz) String() string {
	return "i32.ctz"
}

func (*I32Ctz) ImmString() string {
	return ""
}

type I64Ctz struct{}

func (*I64Ctz) Opcode() Opcode {
	return I64_CTZ
}

func (*I64Ctz) imm() any {
	return NoImm
}

func (*I64Ctz) String() string {
	return "i64.ctz"
}

func (*I64Ctz) ImmString() string {
	return ""
}

type I32Popcnt struct{}

func (*I32Popcnt) Opcode() Opcode {
	return I32_POPCNT
}

func (*I32Popcnt) imm() any {
	return NoImm
}

func (*I32Popcnt) String() string {
	return "i32.popcnt"
}

func (*I32Popcnt) ImmString() string {
	return ""
}

type I64Popcnt struct{}

func (*I64Popcnt) Opcode() Opcode {
	return I64_POPCNT
}

func (*I64Popcnt) imm() any {
	return NoImm
}

func (*I64Popcnt) String() string {
	return "i64.popcnt"
}

func (*I64Popcnt) ImmString() string {
	return ""
}
