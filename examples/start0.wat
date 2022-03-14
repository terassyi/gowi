(module
  (memory (data "A"))
  (func $inc 
  )
  (func $get (result i32)
    (return (i32.add (i32.const 0) (i32.const 1)))
  )
  (func $main
  	(call $inc)
  )

  (start $main)
  (export "inc" (func $inc))
  (export "get" (func $get))
)
