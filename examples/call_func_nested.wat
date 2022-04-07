(module
	(func $rootFunc (export "rootFunc") (param $a i32) (param $b i32) (result i32)
		local.get $a
		local.get $a
		call $childFunc1
		local.get $b
		local.get $b
		call $childFunc2
		i32.add
	)
	(func $childFunc1 (param $a i32) (param $b i32) (result i32)
		local.get $a
		local.get $b
		call $add
	)
	(func $childFunc2 (param $a i32) (param $b i32) (result i32)
		local.get $a
		local.get $b
		call $add
	)
	(func $add (param $a i32) (param $b i32) (result i32)
		local.get $a
		local.get $b
		i32.add
	)
)
