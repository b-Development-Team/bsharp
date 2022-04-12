# B#
Welcome to the B# programming language! Check out the [docs](docs/docs.md) to learn how to code in B#, and keep on reading to learn how to run it on your system!

## Installation
To install or update B#, run
```sh
go install github.com/b-Development-Team/bsharp@latest
```

## Usage
To run a file or multiple files, do
```sh
bsharp run file.bsp <file2.bsp> <file3.bsp>
```

To time how long it takes to run a B# program, use the `--time` or `-t` flag.
```
bsharp run file.bsp -t
```

To compile a file,  do
```
bsharp build file.bsp <file2.bsp> <file3.bsp> -o file
```
This will produce a compiled executable.

The C source code that is generated can be viewed using
```
bsharp build file.bsp <file2.bsp> <file3.bsp> -o file.c
```

To time how long it takes to compile a B# program, use the `--time` or `-t` flag.
```
bsharp build file.bsp -o file -t
```