(module
	(func $f0 (export "f0") (param $a i32) (param $b i32) (result i32)
		local.get $a
		local.get $b
		i32.add
	)
	(func (export "f1") (param $a i32) (param $b i32) (result i32)
		local.get $a
		local.get $b
		call $f0
		i32.const 1
		i32.add
	)
)
