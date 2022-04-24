[![Test main](https://github.com/terassyi/gowi/actions/workflows/test.yml/badge.svg?branch=main)](https://github.com/terassyi/gowi/actions/workflows/test.yml)

# GOWI

Gowi is the Web Assembly interpreter written in Go version 1.18.

## About
Gowi is my first Web Assembly interpreter written in Go.
I only implemented limited features now.
Please see features section.

## Features
[*] Control flow instructions
[*] Integer instructions
[] Float instructions
[] Global values
[] Import some functions
[] Implement instruction validator
[] WASI

## Try
You can try Gowi with docker container or go version 1.18.

### Build
```shell
$ git clone https://github.com/terassyi/gowi.git
$ cd gowi
$ docker run -it --rm -v `pwd`:/gowi gowi:latest bash
$ cd gowi
$ go build .
```

### Test
You can run all test with this command.
```shell
$ go test -v ./...
```

### Run
You can run WASM binary file with gowi.
There are some examples in `examples/`.

For example, We run [examples/fibonacci.wasm](https://github.com/terassyi/gowi/blob/main/examples/fibonacci.wasm) compiled from [examples/fibonacci.wat](https://github.com/terassyi/gowi/blob/main/examples/fibonacci.wat).

#### dump
You can get many information of the target WASM binary.
```shell
$ ./gowi dump -x examples/fibonacci.wasm
WASM file: examples/fibonacci.wasm

examples/fibonacci.wasm: file format wasm 0x1
Section Details:

Type[2]:
 - type[0] (i32) -> (i32)
 - type[1] () -> (i32)

Func[12]:
 - func[0] sig=0
 - func[1] sig=0
 - func[2] sig=1
 - func[3] sig=1
 - func[4] sig=1
 - func[5] sig=1
 - func[6] sig=1
 - func[7] sig=1
 - func[8] sig=1
 - func[9] sig=1
 - func[10] sig=1
 - func[11] sig=1

Export[13]:
 - func[0] <fib> -> index=1
 - func[1] <fib_iterative> -> index=1
 - func[2] <fib_recursive> -> index=0
 - func[3] <fib_rec_exp1> -> index=2
 - func[4] <fib_rec_exp2> -> index=3
 - func[5] <fib_rec_exp3> -> index=4
 - func[6] <fib_rec_exp4> -> index=5
 - func[7] <fib_rec_exp5> -> index=6
 - func[8] <fib_iter_exp1> -> index=7
 - func[9] <fib_iter_exp2> -> index=8
 - func[10] <fib_iter_exp3> -> index=9
 - func[11] <fib_iter_exp4> -> index=10
 - func[12] <fib_iter_exp5> -> index=11

Code[12]:
 - func[0] instruction size=24
 - func[1] instruction size=42
 - func[2] instruction size=3
 - func[3] instruction size=3
 - func[4] instruction size=3
 - func[5] instruction size=3
 - func[6] instruction size=3
 - func[7] instruction size=3
 - func[8] instruction size=3
 - func[9] instruction size=3
 - func[10] instruction size=3
 - func[11] instruction size=3

```

#### Execute
First, you have to find a function you want to run.
You can find the list of all functions in the target binary by running `exec` command with `--list-all-exports` flag.
```shell
$ ./gowi exec examples/fibonacci.wasm --list-all-exports
WASM fule: examples/fibonacci.wasm

List all exported functions
        fib(i32) -> (i32)
        fib_iterative(i32) -> (i32)
        fib_recursive(i32) -> (i32)
        fib_rec_exp1() -> (i32)
        fib_rec_exp2() -> (i32)
        fib_rec_exp3() -> (i32)
        fib_rec_exp4() -> (i32)
        fib_rec_exp5() -> (i32)
        fib_iter_exp1() -> (i32)
        fib_iter_exp2() -> (i32)
        fib_iter_exp3() -> (i32)
        fib_iter_exp4() -> (i32)
        fib_iter_exp5() -> (i32)


```

Next, you can run the target function like bellow.
```shell
$ ./gowi exec examples/fibonacci.wasm --invoke fib --args 10

  fib(10) = (34)
```

And you can trace instructions you run with `--debug 1`.
```shell
$ ./gowi exec examples/fibonacci.wasm --invoke fib --args 10 --debug 1


Invoke fib
--------------------
    i32.const 0x0
    set_local $1
    i32.const 0x1
    set_local $2
    i32.const 0x2
    set_local $4
    get_local $0
    i32.const 0x1
    i32.le_s
    if empty
    end
    get_local $0
    i32.const 0x2
    i32.le_s
    if empty
    end
    loop empty
      get_local $4
      get_local $0
      i32.lt_s
      if empty
        get_local $1
        get_local $2
        i32.add
        set_local $3
        get_local $2
        set_local $1
        get_local $3
        set_local $2
        get_local $4
        i32.const 0x1
        i32.add
        set_local $4
        br 1
      get_local $4
      get_local $0
      i32.lt_s
      if empty
        get_local $1
        get_local $2
        i32.add
        set_local $3
        get_local $2
        set_local $1
        get_local $3
        set_local $2
        get_local $4
        i32.const 0x1
        i32.add
        set_local $4
        br 1
      get_local $4
      get_local $0
      i32.lt_s
      if empty
        get_local $1
        get_local $2
        i32.add
        set_local $3
        get_local $2
        set_local $1
        get_local $3
        set_local $2
        get_local $4
        i32.const 0x1
        i32.add
        set_local $4
        br 1
      get_local $4
      get_local $0
      i32.lt_s
      if empty
        get_local $1
        get_local $2
        i32.add
        set_local $3
        get_local $2
        set_local $1
        get_local $3
        set_local $2
        get_local $4
        i32.const 0x1
        i32.add
        set_local $4
        br 1
      get_local $4
      get_local $0
      i32.lt_s
      if empty
        get_local $1
        get_local $2
        i32.add
        set_local $3
        get_local $2
        set_local $1
        get_local $3
        set_local $2
        get_local $4
        i32.const 0x1
        i32.add
        set_local $4
        br 1
      get_local $4
      get_local $0
      i32.lt_s
      if empty
        get_local $1
        get_local $2
        i32.add
        set_local $3
        get_local $2
        set_local $1
        get_local $3
        set_local $2
        get_local $4
        i32.const 0x1
        i32.add
        set_local $4
        br 1
      get_local $4
      get_local $0
      i32.lt_s
      if empty
        get_local $1
        get_local $2
        i32.add
        set_local $3
        get_local $2
        set_local $1
        get_local $3
        set_local $2
        get_local $4
        i32.const 0x1
        i32.add
        set_local $4
        br 1
      get_local $4
      get_local $0
      i32.lt_s
      if empty
        get_local $1
        get_local $2
        i32.add
        set_local $3
        get_local $2
        set_local $1
        get_local $3
        set_local $2
        get_local $4
        i32.const 0x1
        i32.add
        set_local $4
        br 1
      get_local $4
      get_local $0
      i32.lt_s
      if empty
      end
    end
    get_local $2
  end
Execution Result = (34)

  fib(10) = (34)

```

## Future works
I will implement insufficient features listed in [Features](#features).

## License
Gowi is under the MIT License: See [LICENSE](https://github.com/terassyi/gowi/blob/main/LICENSE) file.
