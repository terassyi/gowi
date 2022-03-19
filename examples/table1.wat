(module
  (table 2 funcref)
  (func $f1 (result i32)
    i32.const 42)
  (func $f2 (result i32)
  i32.const 13)
  (func $f3 (result i32)
  i32.const 100
  )
  (elem (i32.const 0) $f1 $f2 $f3)
  (type $return_i32 (func (result i32)))
  (func (export "callByIndex") (param $i i32) (result i32)
    local.get $i
    call_indirect (type $return_i32))
)
