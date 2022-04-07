
(module
	(func (export "if_func")
		(if (i32.const 0) 
			(then (
				if (i32.const 1) 
				(then (nop))
			)))
	)
)
