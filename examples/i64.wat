(module
  (func (export "add") (param $x i64) (param $y i64) (result i64) (i64.add (local.get $x) (local.get $y)))
  (func (export "sub") (param $x i64) (param $y i64) (result i64) (i64.sub (local.get $x) (local.get $y)))
  (func (export "mul") (param $x i64) (param $y i64) (result i64) (i64.mul (local.get $x) (local.get $y)))
  (func (export "div_s") (param $x i64) (param $y i64) (result i64) (i64.div_s (local.get $x) (local.get $y)))
  (func (export "div_u") (param $x i64) (param $y i64) (result i64) (i64.div_u (local.get $x) (local.get $y)))
  (func (export "rem_s") (param $x i64) (param $y i64) (result i64) (i64.rem_s (local.get $x) (local.get $y)))
  (func (export "rem_u") (param $x i64) (param $y i64) (result i64) (i64.rem_u (local.get $x) (local.get $y)))
  (func (export "and") (param $x i64) (param $y i64) (result i64) (i64.and (local.get $x) (local.get $y)))
  (func (export "or") (param $x i64) (param $y i64) (result i64) (i64.or (local.get $x) (local.get $y)))
  (func (export "xor") (param $x i64) (param $y i64) (result i64) (i64.xor (local.get $x) (local.get $y)))
  (func (export "shl") (param $x i64) (param $y i64) (result i64) (i64.shl (local.get $x) (local.get $y)))
;;  (func (export "shr_s") (param $x i64) (param $y i64) (result i64) (i64.shr_s (local.get $x) (local.get $y)))
;;  (func (export "shr_u") (param $x i64) (param $y i64) (result i64) (i64.shr_u (local.get $x) (local.get $y)))
;;  (func (export "rotl") (param $x i64) (param $y i64) (result i64) (i64.rotl (local.get $x) (local.get $y)))
;;  (func (export "rotr") (param $x i64) (param $y i64) (result i64) (i64.rotr (local.get $x) (local.get $y)))
;;  (func (export "clz") (param $x i64) (result i64) (i64.clz (local.get $x)))
;;  (func (export "ctz") (param $x i64) (result i64) (i64.ctz (local.get $x)))
;;  (func (export "popcnt") (param $x i64) (result i64) (i64.popcnt (local.get $x)))
;;  (func (export "extend8_s") (param $x i64) (result i64) (i64.extend8_s (local.get $x)))
;;  (func (export "extend16_s") (param $x i64) (result i64) (i64.extend16_s (local.get $x)))
;;  (func (export "extend32_s") (param $x i64) (result i64) (i64.extend32_s (local.get $x)))
;;  (func (export "eqz") (param $x i64) (result i32) (i64.eqz (local.get $x)))
;;  (func (export "eq") (param $x i64) (param $y i64) (result i32) (i64.eq (local.get $x) (local.get $y)))
;;  (func (export "ne") (param $x i64) (param $y i64) (result i32) (i64.ne (local.get $x) (local.get $y)))
;;  (func (export "lt_s") (param $x i64) (param $y i64) (result i32) (i64.lt_s (local.get $x) (local.get $y)))
;;  (func (export "lt_u") (param $x i64) (param $y i64) (result i32) (i64.lt_u (local.get $x) (local.get $y)))
;;  (func (export "le_s") (param $x i64) (param $y i64) (result i32) (i64.le_s (local.get $x) (local.get $y)))
;;  (func (export "le_u") (param $x i64) (param $y i64) (result i32) (i64.le_u (local.get $x) (local.get $y)))
;;  (func (export "gt_s") (param $x i64) (param $y i64) (result i32) (i64.gt_s (local.get $x) (local.get $y)))
;;  (func (export "gt_u") (param $x i64) (param $y i64) (result i32) (i64.gt_u (local.get $x) (local.get $y)))
;;  (func (export "ge_s") (param $x i64) (param $y i64) (result i32) (i64.ge_s (local.get $x) (local.get $y)))
;;  (func (export "ge_u") (param $x i64) (param $y i64) (result i32) (i64.ge_u (local.get $x) (local.get $y)))
)
