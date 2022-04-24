
(module
  ;; Typing

  (func (export "type-local-i32") (result i32) (local i32) (local.get 0))
  (func (export "type-local-i64") (result i64) (local i64) (local.get 0))
  (func $param-and-local-get-param (export "param-and-local-get-param") (param i32) (result i32) (local i32 i32) (local.get 0))
  (func $param-and-local-get-local0 (export "param-and-local-get-local0") (param i32) (result i32) (local i32 i32) (local.get 1))
  (func $param-and-local-get-local1 (export "param-and-local-get-local1") (param i32) (result i32) (local i32 i32) (i32.const 1) (local.set 2) (local.get 2))

  (func (export "use-param-and-local-get-param") (result i32) (i32.const 0xeb) (call $param-and-local-get-param))
  (func (export "use-param-and-local-get-local0") (result i32) (i32.const 0xeb) (call $param-and-local-get-local0))
  (func (export "use-param-and-local-get-local1") (result i32) (i32.const 0xeb) (call $param-and-local-get-local1))
)
