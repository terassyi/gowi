(module
  (memory (data "A"))
  (func $inc
    (i32.store8
      (i32.const 0)
      (i32.add
        (i32.load8_u (i32.const 0))
        (i32.const 1)
      )
    )
  )
  (func $get (result i32)
    (return (i32.load8_u (i32.const 0)))
  )
  (func $main
    (call $inc)
    (call $inc)
    (call $inc)
  )

  (start $main)
  (export "inc" (func $inc))
  (export "get" (func $get))
)
