(module
  (global $g (import "js" "global") (mut i32))
  (func (export "getGlobal1") (result i32)
    (global.get $g)
  )
  (func (export "incGlobal")
    (global.set $g)
      (i32.get (global.get $g) (i32.const 1)))
)
