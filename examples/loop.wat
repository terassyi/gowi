(module
  (func $dummy)

  (func (export "empty")
    (loop)
    (loop $l)
  )

  (func (export "singular") (result i32)
    (loop (nop))
    (loop (result i32) (i32.const 7))
  )

  (func (export "multi") (result i32)
    (loop (call $dummy) (call $dummy) (call $dummy) (call $dummy))
    (loop (result i32) (call $dummy) (call $dummy) (i32.const 8) (call $dummy))
    (drop)
    (loop (result i32 i32 i32)
      (call $dummy) (call $dummy) (call $dummy) (i32.const 8) (call $dummy)
      (call $dummy) (call $dummy) (call $dummy) (i32.const 7) (call $dummy)
      (call $dummy) (call $dummy) (call $dummy) (i32.const 9) (call $dummy)
    )
    (drop) (drop)
  )

  (func (export "nest") (result i32)
    (loop (result i32)
      (loop (call $dummy) (block) (nop))
      (loop (result i32) (call $dummy) (i32.const 9))
    )
  )

  (func (export "deep") (result i32)
    (loop (result i32) (block (result i32)
      (loop (result i32) (block (result i32)
        (loop (result i32) (block (result i32)
          (loop (result i32) (block (result i32)
            (loop (result i32) (block (result i32)
              (loop (result i32) (block (result i32)
                (loop (result i32) (block (result i32)
                  (loop (result i32) (block (result i32)
                    (loop (result i32) (block (result i32)
                      (loop (result i32) (block (result i32)
                        (loop (result i32) (block (result i32)
                          (loop (result i32) (block (result i32)
                            (loop (result i32) (block (result i32)
                              (loop (result i32) (block (result i32)
                                (loop (result i32) (block (result i32)
                                  (loop (result i32) (block (result i32)
                                    (loop (result i32) (block (result i32)
                                      (loop (result i32) (block (result i32)
                                        (loop (result i32) (block (result i32)
                                          (loop (result i32) (block (result i32)
                                            (call $dummy) (i32.const 150)
                                          ))
                                        ))
                                      ))
                                    ))
                                  ))
                                ))
                              ))
                            ))
                          ))
                        ))
                      ))
                    ))
                  ))
                ))
              ))
            ))
          ))
        ))
      ))
    ))
  )
  (func (export "param-break") (result i32)
    (local $x i32)
    (i32.const 1)
    (loop (param i32) (result i32)
      (i32.const 4)
      (i32.add)
      (local.tee $x)
      (local.get $x)
      (i32.const 10)
      (i32.lt_u)
      (br_if 0)
    )
  )
  (func (export "params-break") (result i32)
    (local $x i32)
    (i32.const 1)
    (i32.const 2)
    (loop (param i32 i32) (result i32)
      (i32.add)
      (local.tee $x)
      (i32.const 3)
      (local.get $x)
      (i32.const 10)
      (i32.lt_u)
      (br_if 0)
      (drop)
    )
  )
  (func (export "params-id-break") (result i32)
    (local $x i32)
    (local.set $x (i32.const 0))
    (i32.const 1)
    (i32.const 2)
    (loop (param i32 i32) (result i32 i32)
      (local.set $x (i32.add (local.get $x) (i32.const 1)))
      (br_if 0 (i32.lt_u (local.get $x) (i32.const 10)))
    )
    (i32.add)
  )

  (func $fx (export "effects") (result i32)
    (local i32)
    (block
      (loop
        (local.set 0 (i32.const 1))
        (local.set 0 (i32.mul (local.get 0) (i32.const 3)))
        (local.set 0 (i32.sub (local.get 0) (i32.const 5)))
        (local.set 0 (i32.mul (local.get 0) (i32.const 7)))
        (br 1)
        (local.set 0 (i32.mul (local.get 0) (i32.const 100)))
      )
    )
    (i32.eq (local.get 0) (i32.const -14))
  )

  (func (export "while") (param i64) (result i64)
    (local i64)
    (local.set 1 (i64.const 1))
    (block
      (loop
        (br_if 1 (i64.eqz (local.get 0)))
        (local.set 1 (i64.mul (local.get 0) (local.get 1)))
        (local.set 0 (i64.sub (local.get 0) (i64.const 1)))
        (br 0)
      )
    )
    (local.get 1)
  )

  (func (export "for") (param i64) (result i64)
    (local i64 i64)
    (local.set 1 (i64.const 1))
    (local.set 2 (i64.const 2))
    (block
      (loop
        (br_if 1 (i64.gt_u (local.get 2) (local.get 0)))
        (local.set 1 (i64.mul (local.get 1) (local.get 2)))
        (local.set 2 (i64.add (local.get 2) (i64.const 1)))
        (br 0)
      )
    )
    (local.get 1)
  )

  (func (export "nesting") (param i32 i32) (result i32)
    (local i32 i32)
    (block
      (loop
        (br_if 1 (i32.eq (local.get 0) (i32.const 0)))
        (local.set 2 (local.get 1))
        (block
          (loop
            (br_if 1 (i32.eq (local.get 2) (i32.const 0)))
            (br_if 3 (i32.lt_s (local.get 2) (i32.const 0)))
            (local.set 3 (i32.add (local.get 3) (local.get 2)))
            (local.set 2 (i32.sub (local.get 2) (i32.const 2)))
            (br 0)
          )
        )
        (local.set 3 (i32.div_s (local.get 3) (local.get 0)))
        (local.set 0 (i32.sub (local.get 0) (i32.const 1)))
        (br 0)
      )
    )
    (local.get 3)
  )
)
