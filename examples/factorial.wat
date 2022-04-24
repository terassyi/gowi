(module
	(func $factorial (param $n i32) (result i32)
		(if (result i32)
		(i32.lt_u (i32.const 0) (local.get $n))
		(then 
      (i32.sub (local.get $n) (i32.const 1))
      (i32.mul (call $factorial) (local.get $n))
    )
    (else (i32.const 1))
		)
	)

  (func $exp1 (result i32) 
    (i32.const 0) (call $factorial))
  (func $exp2 (result i32) 
    (i32.const 1) (call $factorial))
  (func $exp3 (result i32) 
    (i32.const 5) (call $factorial))

  (export "exp1" (func $exp1))
  (export "exp2" (func $exp2))
  (export "exp3" (func $exp3))
  (export "factorial" (func $factorial))
)
