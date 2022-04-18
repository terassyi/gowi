(module

  ;; recursive implementation

  (func $fib_recursive (param $N i32) (result i32)
    (if (result i32)
      (i32.lt_u (local.get $N) (i32.const 3))
      (then
        (if (result i32)
          (i32.eq (i32.const 1) (local.get $N))
          (then (i32.const 0))
        (else (i32.const 1))
        )
      )
    (else 
      (i32.add 
        (call $fib_recursive
          (i32.sub (local.get $N) (i32.const 1)))
        (call $fib_recursive
          (i32.sub (local.get $N) (i32.const 2)))
    )))
  )

  ;; iterative implementation, avoids stack overflow

  (func $fib_iterative (param $N i32) (result i32)
    (local $n1 i32)
    (local $n2 i32)
    (local $tmp i32)
    (local $i i32)
    (local.set $n1 (i32.const 1))
    (local.set $n2 (i32.const 1))
    (local.set $i (i32.const 2))


    ;; return 0 for N <= 0
    (if
      (i32.le_s (local.get $N) (i32.const 0))
      (then (return (i32.const 0)))
    )

    ;;since we normally return n2, handle n=1 case specially
    (if
      (i32.le_s (local.get $N) (i32.const 2))
      (then (return (i32.const 1)))
    )

    (loop $again
      (if
        (i32.lt_s (local.get $i) (local.get $N))
        (then
          (local.set $tmp (i32.add (local.get $n1) (local.get $n2)))
          (local.set $n1 (local.get $n2))
          (local.set $n2 (local.get $tmp))
          (local.set $i (i32.add (local.get $i) (i32.const 1)))
          (br $again)
        )
      )
    )

    (local.get $n2)
  )

  ;; export fib_iterative as the main thing, because it's the fastest
  (func $fib_rec_exp1 (result i32)
    (i32.const 0)
    (call $fib_recursive)
  )
  (func $fib_rec_exp2 (result i32)
    (i32.const 1)
    (call $fib_recursive)
  )
  (func $fib_rec_exp3 (result i32)
    (i32.const 2)
    (call $fib_recursive)
  )
  (func $fib_rec_exp4 (result i32)
    (i32.const 3)
    (call $fib_recursive)
  )
  (func $fib_rec_exp5 (result i32)
    (i32.const 20)
    (call $fib_recursive)
  )

  (export "fib" (func $fib_iterative))
  (export "fib_iterative" (func $fib_iterative))
  (export "fib_recursive" (func $fib_recursive))
  (export "fib_rec_exp1" (func $fib_rec_exp1))
  (export "fib_rec_exp2" (func $fib_rec_exp2))
  (export "fib_rec_exp3" (func $fib_rec_exp3))
  (export "fib_rec_exp4" (func $fib_rec_exp4))
  (export "fib_rec_exp5" (func $fib_rec_exp5))
)
